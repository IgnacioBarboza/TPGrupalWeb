package main

import(
	"context"
	"testing"
	"time"
	"os" // leer env
	"fmt"
	db "WebGrupal/db/sqlc"
	"database/sql"
	_ "github.com/lib/pq" //Driver postgresql
)

func OpenDB(t *testing.T) *sql.DB {

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
		t.Fatalf("no se pudo abrir la conexión: %v", err)
	}

	if err := conn.Ping(); err != nil {
		conn.Close()
		t.Fatalf("no se pudo conectar a la base de datos: %v", err)
	}

	t.Cleanup(func() {
		conn.Close()
	})

	return conn
}

func TestCRUDFilm(t *testing.T) {
	connection := OpenDB(t)

	ctx := context.Background()
	repository := db.New(connection)

	fechaEstreno := time.Date(    // 2026 - 01 - 01  00hs 00min 00sec UTC
		2026,
		time.January,
		1,
		0, 0, 0, 0,
		time.UTC,
	)
	pelicula_prueba := db.CreateFilmParams{
		Filmid          :0,
		Filmtitle       :"Pelicula Testeo",
		Filmreleasedate :fechaEstreno,
		Filmduration    :120,
	}
	pelicula, err := repository.CreateFilm(ctx, pelicula_prueba)
	if err != nil {
		t.Fatalf("No se pudo crear la pelicula por: %v",err)
	}

	id := pelicula.Filmid

	t.Run("editar film", func(t *testing.T) {
		pelicula_actualizada := db.UpdateFilmParams{
			Filmid          : id,
			Filmtitle       : "Pelicula Editada",
			Filmreleasedate : pelicula_prueba.Filmreleasedate,
			Filmduration    : pelicula_prueba.Filmduration,
		}
		err = repository.UpdateFilm(ctx,pelicula_actualizada)
		if err != nil {
			t.Fatalf("No se puedo editar la pelicula por: %v",err)
		}		
	})
	
	t.Run("obtener film", func(t *testing.T) {
		pelicula_obtenida, err := repository.GetFilm(ctx,id)
		if err != nil{
			t.Fatalf("No se pudo obtener la pelicula por: %v",err)
		}

		/*
		t.Logf(
			"esperada: id=%d fecha=%s titulo=%q duracion=%d",
			id,
			"2026-01-01",
			"Pelicula Editada",
			pelicula_prueba.Filmduration,
		)

		t.Logf(
			"Película obtenida: ID=%d, fecha=%s, duración=%d, título=%q",
			pelicula_obtenida.Filmid,
			pelicula_obtenida.Filmreleasedate.Format("2006-01-02"),
			pelicula_obtenida.Filmduration,
			pelicula_obtenida.Filmtitle,
		)
		*/
		if pelicula_obtenida.Filmreleasedate.Format("2006-01-02") != "2026-01-01" {
			t.Errorf("fecha incorrecta")
		}

		if pelicula_obtenida.Filmtitle != "Pelicula Editada" {
			t.Errorf("titulo incorrecto")
		}

		if pelicula_obtenida.Filmduration != pelicula_prueba.Filmduration {
			t.Errorf("duracion incorrecta")
		}
	})


	t.Run("borrar film", func(t *testing.T) {
		err = repository.DeleteFilm(ctx,id)
		if err != nil {
			t.Fatalf("Hubo un error en el delete de la tabla Film por: %v",err)
		}
	})
}