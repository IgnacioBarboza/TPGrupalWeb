package main

import (
	db "WebGrupal/db/sqlc"
	film "WebGrupal/internal/film"
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

// Handler global para poder ser usado en los tests de integración
var filmHandler *film.FilmHandler

//var reviewHandler *film.ReviewHandler

func initializeServices() {
	conn := OpenDB()
	queries := db.New(conn)

	// Creamos repositorio, servicio y handler para la entidad Film
	filmRepo := film.NewSQLCRepository(queries)
	filmSvc := film.NewFilmService(filmRepo)
	filmHandler = film.NewFilmHandler(filmSvc)

	//reviewRepo := review.NewSQLCRepository(queries)
	//reviewSvc := review.NewFilmService(filmRepo)
	//reviewHandler = review.NewFilmHandler(filmSvc)
}

func main() {
	initializeServices()
	mux := http.NewServeMux()
	filmHandler.Register(mux)

	log.Println("Escuchando en :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
