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

// Handler globales
var filmHandler *film.FilmHandler
var reviewHandler *Review.ReviewHandler
var userHandler *User.UserHandler
var statusHandler *Status.StatusHandler
var castMemberHandler *CastMember.CastMemberHandler
var genreHandler *Genre.GenreHandler
var filmCastHandler *FilmCast.FilmCastHandler
var filmGenreHandler *FilmGenre.FilmGenreHandler
var filmUserStatusHandler *FilmUserStatus.FilmUserStatusHandler

func initializeServices() {
	conn := OpenDB()
	queries := db.New(conn)

	// Creamos repositorio, servicio y handler para la entidad Film
	filmRepo := film.NewSQLCRepository(queries)
	filmSvc := film.NewFilmService(filmRepo)
	filmHandler = film.NewFilmHandler(filmSvc)

	reviewRepo := Review.NewSQLCRepository(queries)
	reviewSvc := Review.NewReviewService(reviewRepo)
	reviewHandler = Review.NewReviewHandler(reviewSvc)

	userRepo := User.NewSQLCRepository(queries)
	userSvc := User.NewUserService(userRepo)
	userHandler = User.NewUserHandler(userSvc)

	statusRepo := Status.NewSQLCRepository(queries)
	statusSvc := Status.NewStatusService(statusRepo)
	statusHandler = Status.NewStatusHandler(statusSvc)

	castMemberRepo := CastMember.NewSQLCRepository(queries)
	castMemberSvc := CastMember.NewCastMemberService(castMemberRepo)
	castMemberHandler = CastMember.NewCastMemberHandler(castMemberSvc)

	genreRepo := Genre.NewSQLCRepository(queries)
	genreSvc := Genre.NewGenreService(genreRepo)
	genreHandler = Genre.NewGenreHandler(genreSvc)

	filmCastRepo := FilmCast.NewSQLCRepository(queries)
	filmCastSvc := FilmCast.NewFilmCastService(filmCastRepo)
	filmCastHandler = FilmCast.NewFilmCastHandler(filmCastSvc)

	filmGenreRepo := FilmGenre.NewSQLCRepository(queries)
	filmGenreSvc := FilmGenre.NewFilmGenreService(filmGenreRepo)
	filmGenreHandler = FilmGenre.NewFilmGenreHandler(filmGenreSvc)

	filmUserStatusRepo := FilmUserStatus.NewSQLCRepository(queries)
	filmUserStatusSvc := FilmUserStatus.NewFilmUserStatusService(filmUserStatusRepo)
	filmUserStatusHandler = FilmUserStatus.NewFilmUserStatusHandler(filmUserStatusSvc)
}

func main() {
	initializeServices()
	mux := http.NewServeMux()
	filmHandler.Register(mux)
	reviewHandler.Register(mux)
	userHandler.Register(mux)
	statusHandler.Register(mux)
	castMemberHandler.Register(mux)
	genreHandler.Register(mux)
	filmCastHandler.Register(mux)
	filmGenreHandler.Register(mux)
	filmUserStatusHandler.Register(mux)

	log.Println("Escuchando en :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
