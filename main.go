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
	// RENDERIZAR TODA LA PAGINA

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

	// Un handler por entidad: atiende el POST del formulario y el GET de la lista
	filmHandler := film.NewFilmHandler(filmSvc)
	genreHandler := Genre.NewGenreHandler(genreSvc)
	castMemberHandler := CastMember.NewCastMemberHandler(castMemberSvc)
	statusHandler := Status.NewStatusHandler(statusSvc)
	userHandler := User.NewUserHandler(userSvc)
	reviewHandler := Review.NewReviewHandler(reviewSvc)
	filmCastHandler := FilmCast.NewFilmCastHandler(filmCastSvc)
	filmGenreHandler := FilmGenre.NewFilmGenreHandler(filmGenreSvc)
	filmUserStatusHandler := FilmUserStatus.NewFilmUserStatusHandler(filmUserStatusSvc)

	// Formularios: POST /<entidades> crea el registro y redirige a "/" (PRG)
	mux.HandleFunc("POST /films", filmHandler.CreateFromForm)
	mux.HandleFunc("POST /genres", genreHandler.CreateFromForm)
	mux.HandleFunc("POST /castmembers", castMemberHandler.CreateFromForm)
	mux.HandleFunc("POST /statuses", statusHandler.CreateFromForm)
	mux.HandleFunc("POST /users", userHandler.CreateFromForm)
	mux.HandleFunc("POST /reviews", reviewHandler.CreateFromForm)
	mux.HandleFunc("POST /filmcasts", filmCastHandler.CreateFromForm)
	mux.HandleFunc("POST /filmgenres", filmGenreHandler.CreateFromForm)
	mux.HandleFunc("POST /filmuserstatuses", filmUserStatusHandler.CreateFromForm)

	// Listas: GET /<entidades> devuelve solo la lista de esa entidad (fragmento HTML)
	mux.HandleFunc("GET /films", filmHandler.ListFilms)
	mux.HandleFunc("GET /genres", genreHandler.ListGenres)
	mux.HandleFunc("GET /castmembers", castMemberHandler.ListCastMembers)
	mux.HandleFunc("GET /statuses", statusHandler.ListStatuses)
	mux.HandleFunc("GET /users", userHandler.ListUsers)
	mux.HandleFunc("GET /reviews", reviewHandler.ListReviews)
	mux.HandleFunc("GET /filmcasts", filmCastHandler.ListFilmCasts)
	mux.HandleFunc("GET /filmgenres", filmGenreHandler.ListFilmGenres)
	mux.HandleFunc("GET /filmuserstatuses", filmUserStatusHandler.ListFilmUserStatuses)

	// Borrado: DELETE /<entidades>/{id} responde 200 vacío. Las relaciones se identifican con dos IDs
	mux.HandleFunc("DELETE /films/{id}", filmHandler.Delete)
	mux.HandleFunc("DELETE /genres/{id}", genreHandler.Delete)
	mux.HandleFunc("DELETE /castmembers/{id}", castMemberHandler.Delete)
	mux.HandleFunc("DELETE /statuses/{id}", statusHandler.Delete)
	mux.HandleFunc("DELETE /users/{id}", userHandler.Delete)
	mux.HandleFunc("DELETE /reviews/{filmid}/{userid}", reviewHandler.Delete)
	mux.HandleFunc("DELETE /filmcasts/{filmid}/{castmemberid}", filmCastHandler.Delete)
	mux.HandleFunc("DELETE /filmgenres/{genreid}/{filmid}", filmGenreHandler.Delete)
	mux.HandleFunc("DELETE /filmuserstatuses/{filmid}/{userid}", filmUserStatusHandler.Delete)

	// Edición (click to edit): GET .../{id} devuelve la fila, GET .../{id}/edit la fila como formulario
	// y PUT .../{id} guarda y devuelve la fila actualizada
	mux.HandleFunc("GET /films/{id}", filmHandler.Item)
	mux.HandleFunc("GET /films/{id}/edit", filmHandler.EditForm)
	mux.HandleFunc("PUT /films/{id}", filmHandler.Update)
	mux.HandleFunc("GET /genres/{id}", genreHandler.Item)
	mux.HandleFunc("GET /genres/{id}/edit", genreHandler.EditForm)
	mux.HandleFunc("PUT /genres/{id}", genreHandler.Update)
	mux.HandleFunc("GET /castmembers/{id}", castMemberHandler.Item)
	mux.HandleFunc("GET /castmembers/{id}/edit", castMemberHandler.EditForm)
	mux.HandleFunc("PUT /castmembers/{id}", castMemberHandler.Update)
	mux.HandleFunc("GET /statuses/{id}", statusHandler.Item)
	mux.HandleFunc("GET /statuses/{id}/edit", statusHandler.EditForm)
	mux.HandleFunc("PUT /statuses/{id}", statusHandler.Update)
	mux.HandleFunc("GET /users/{id}", userHandler.Item)
	mux.HandleFunc("GET /users/{id}/edit", userHandler.EditForm)
	mux.HandleFunc("PUT /users/{id}", userHandler.Update)
	mux.HandleFunc("GET /reviews/{filmid}/{userid}", reviewHandler.Item)
	mux.HandleFunc("GET /reviews/{filmid}/{userid}/edit", reviewHandler.EditForm)
	mux.HandleFunc("PUT /reviews/{filmid}/{userid}", reviewHandler.Update)
	mux.HandleFunc("GET /filmuserstatuses/{filmid}/{userid}", filmUserStatusHandler.Item)
	mux.HandleFunc("GET /filmuserstatuses/{filmid}/{userid}/edit", filmUserStatusHandler.EditForm)
	mux.HandleFunc("PUT /filmuserstatuses/{filmid}/{userid}", filmUserStatusHandler.Update)
	// En FilmCast y FilmGenre la fila es solo el par de IDs: la URL trae el par actual y el formulario el nuevo
	mux.HandleFunc("GET /filmcasts/{filmid}/{castmemberid}", filmCastHandler.Item)
	mux.HandleFunc("GET /filmcasts/{filmid}/{castmemberid}/edit", filmCastHandler.EditForm)
	mux.HandleFunc("PUT /filmcasts/{filmid}/{castmemberid}", filmCastHandler.Update)
	mux.HandleFunc("GET /filmgenres/{genreid}/{filmid}", filmGenreHandler.Item)
	mux.HandleFunc("GET /filmgenres/{genreid}/{filmid}/edit", filmGenreHandler.EditForm)
	mux.HandleFunc("PUT /filmgenres/{genreid}/{filmid}", filmGenreHandler.Update)

	log.Println("Servidor escuchando en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
