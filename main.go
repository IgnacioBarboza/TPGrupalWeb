package main

import (
	db "WebGrupal/db/sqlc"
	CastMember "WebGrupal/internal/castmember"
	film "WebGrupal/internal/film"
	FilmCast "WebGrupal/internal/filmcast"
	FilmGenre "WebGrupal/internal/filmgenre"
	FilmUserStatus "WebGrupal/internal/filmuserstatus"
	Genre "WebGrupal/internal/genre"
	Review "WebGrupal/internal/review"
	Status "WebGrupal/internal/status"
	User "WebGrupal/internal/user"
	views "WebGrupal/views"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func OpenDB() *sql.DB {

	// Seteo variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := "disable"

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host,
		port,
		user,
		password,
		dbname,
		sslmode,
	)

	conn, err := sql.Open("postgres", dsn)

	if err != nil {
		log.Fatalf("Error al abrir la BD: %v", err)
	}

	if err := conn.Ping(); err != nil {
		conn.Close()
		log.Fatalf("No se pudo conectar a la base de datos: %v", err)
	}

	return conn
}

// Servicios y Handlers globales
var filmSvc *film.FilmService
var reviewSvc *Review.ReviewService
var userSvc *User.UserService
var statusSvc *Status.StatusService
var castMemberSvc *CastMember.CastMemberService
var genreSvc *Genre.GenreService
var filmCastSvc *FilmCast.FilmCastService
var filmGenreSvc *FilmGenre.FilmGenreService
var filmUserStatusSvc *FilmUserStatus.FilmUserStatusService

func initializeServices() {
	conn := OpenDB()
	queries := db.New(conn)

	// Creamos repositorio, servicio y handler para la entidad Film
	filmRepo := film.NewSQLCRepository(queries)
	filmSvc = film.NewFilmService(filmRepo)

	reviewRepo := Review.NewSQLCRepository(queries)
	reviewSvc = Review.NewReviewService(reviewRepo)

	userRepo := User.NewSQLCRepository(queries)
	userSvc = User.NewUserService(userRepo)

	statusRepo := Status.NewSQLCRepository(queries)
	statusSvc = Status.NewStatusService(statusRepo)

	castMemberRepo := CastMember.NewSQLCRepository(queries)
	castMemberSvc = CastMember.NewCastMemberService(castMemberRepo)

	genreRepo := Genre.NewSQLCRepository(queries)
	genreSvc = Genre.NewGenreService(genreRepo)

	filmCastRepo := FilmCast.NewSQLCRepository(queries)
	filmCastSvc = FilmCast.NewFilmCastService(filmCastRepo)

	filmGenreRepo := FilmGenre.NewSQLCRepository(queries)
	filmGenreSvc = FilmGenre.NewFilmGenreService(filmGenreRepo)

	filmUserStatusRepo := FilmUserStatus.NewSQLCRepository(queries)
	filmUserStatusSvc = FilmUserStatus.NewFilmUserStatusService(filmUserStatusRepo)
}

func getPageData(ctx context.Context) (views.PageData, error) {
	data := views.PageData{}

	films, err := filmSvc.GetFilms(ctx)
	if err != nil {
		return data, err
	}
	data.Films = films

	reviews, err := reviewSvc.GetReviews(ctx)
	if err != nil {
		return data, err
	}
	data.Reviews = reviews

	users, err := userSvc.GetUsers(ctx)
	if err != nil {
		return data, err
	}
	data.Users = users

	statuses, err := statusSvc.GetStatuses(ctx)
	if err != nil {
		return data, err
	}
	data.Statuses = statuses

	castMembers, err := castMemberSvc.GetCastMembers(ctx)
	if err != nil {
		return data, err
	}
	data.CastMembers = castMembers

	genres, err := genreSvc.GetGenres(ctx)
	if err != nil {
		return data, err
	}
	data.Genres = genres

	filmCasts, err := filmCastSvc.GetFilmCasts(ctx)
	if err != nil {
		return data, err
	}
	data.FilmCasts = filmCasts

	filmGenres, err := filmGenreSvc.GetFilmGenres(ctx)
	if err != nil {
		return data, err
	}
	data.FilmGenres = filmGenres

	filmUserStatuses, err := filmUserStatusSvc.GetFilmUserStatuses(ctx)
	if err != nil {
		return data, err
	}
	data.FilmUserStatuses = filmUserStatuses

	return data, nil
}

func HandlerPrincipal(w http.ResponseWriter, r *http.Request) {

	// Obtiene los datos de la página desde la base de datos
	data, err := getPageData(r.Context())
	if err != nil {
		log.Printf("Error obteniendo datos: %v", err)
		http.Error(w, "Error al obtener los datos", http.StatusInternalServerError)
		return
	}

	// Renderiza la página principal con los datos obtenidos
	page := views.Index(data)

	err = page.Render(r.Context(), w)
	if err != nil {
		log.Printf("Error renderizando la página: %v", err)
	}
}

func main() {
	initializeServices()
	mux := http.NewServeMux()

	// Archivos estáticos (CSS, imágenes, etc.) bajo /static/
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

	// Página principal: solo la ruta exacta "/"
	mux.HandleFunc("GET /{$}", HandlerPrincipal)

	// Formularios: POST /<entidades> crea el registro y redirige a "/" (PRG)
	mux.HandleFunc("POST /films", film.NewFilmHandler(filmSvc).CreateFromForm)
	mux.HandleFunc("POST /genres", Genre.NewGenreHandler(genreSvc).CreateFromForm)
	mux.HandleFunc("POST /castmembers", CastMember.NewCastMemberHandler(castMemberSvc).CreateFromForm)
	mux.HandleFunc("POST /statuses", Status.NewStatusHandler(statusSvc).CreateFromForm)
	mux.HandleFunc("POST /users", User.NewUserHandler(userSvc).CreateFromForm)
	mux.HandleFunc("POST /reviews", Review.NewReviewHandler(reviewSvc).CreateFromForm)
	mux.HandleFunc("POST /filmcasts", FilmCast.NewFilmCastHandler(filmCastSvc).CreateFromForm)
	mux.HandleFunc("POST /filmgenres", FilmGenre.NewFilmGenreHandler(filmGenreSvc).CreateFromForm)
	mux.HandleFunc("POST /filmuserstatuses", FilmUserStatus.NewFilmUserStatusHandler(filmUserStatusSvc).CreateFromForm)

	log.Println("Servidor escuchando en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
