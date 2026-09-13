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

func TestCRUDUser(t *testing.T) {
	connection := OpenDB(t)
	ctx := context.Background()
	repository := db.New(connection)

	var id_user int32
	var err error

	dateTesting := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)

	t.Run("Crear User", func(t *testing.T) {
		u1 := db.CreateUserParams{
			Userid:             0,
			Username:           "Test Username",
			Userpasswordhashed: "hashed_pass",
			Userbirthday:       dateTesting,
			Usergenre:          "M",
		}

		u, err := repository.CreateUser(ctx, u1)
		if err != nil {
			t.Fatalf("Error al crear un usuario: %v", err)
		}
		id_user = u.Userid
	})

	t.Run("Actualizar User", func(t *testing.T) {
		u_update := db.UpdateUserParams{
			Userid:             id_user,
			Username:           "Test Username Updated",
			Userpasswordhashed: "new_hashed_pass",
			Userbirthday:       dateTesting,
			Usergenre:          "F",
		}

		err = repository.UpdateUser(ctx, u_update)
		if err != nil {
			t.Fatalf("Error al intentar actualizar el usuario: %v", err)
		}
	})

	t.Run("Obtener User y comparar", func(t *testing.T) {
		u, err := repository.GetUser(ctx, id_user)
		if err != nil {
			t.Fatalf("Error al obtener el usuario: %v", err)
		}

		if u.Username != "Test Username Updated" {
			t.Fatalf("El dato obtenido no es el mismo del update.")
		}
	})

	t.Run("Borrar User", func(t *testing.T) {
		err = repository.DeleteUser(ctx, id_user)
		if err != nil {
			t.Fatalf("Error al borrar el usuario: %v", err)
		}
	})
}

func TestCRUDGenre(t *testing.T) {
	connection := OpenDB(t)
	ctx := context.Background()
	repository := db.New(connection)

	var id_genre int32
	var err error

	t.Run("Crear Genre", func(t *testing.T) {
		g1 := db.CreateGenreParams{
			Genreid:   0,
			Genrename: "Sci-Fi",
		}

		g, err := repository.CreateGenre(ctx, g1)
		if err != nil {
			t.Fatalf("Error al crear un genero: %v", err)
		}
		id_genre = g.Genreid
	})

	t.Run("Actualizar Genre", func(t *testing.T) {
		g_update := db.UpdateGenreParams{
			Genreid:   id_genre,
			Genrename: "Science Fiction",
		}

		err = repository.UpdateGenre(ctx, g_update)
		if err != nil {
			t.Fatalf("Error al actualizar genero: %v", err)
		}
	})

	t.Run("Obtener Genre y comparar", func(t *testing.T) {
		g, err := repository.GetGenre(ctx, id_genre)
		if err != nil {
			t.Fatalf("Error al obtener genero: %v", err)
		}

		if g.Genrename != "Science Fiction" {
			t.Fatalf("El genero obtenido no coincide con la actualizacion.")
		}
	})

	t.Run("Borrar Genre", func(t *testing.T) {
		err = repository.DeleteGenre(ctx, id_genre)
		if err != nil {
			t.Fatalf("Error al borrar genero: %v", err)
		}
	})
}
func TestCRUDStatus(t *testing.T) {
	connection := OpenDB(t)
	ctx := context.Background()
	repository := db.New(connection)

	var id_status int32
	var err error

	t.Run("Crear Status", func(t *testing.T) {
		s1 := db.CreateStatusParams{
			Statusid:   0,
			Statusname: "Pendiente",
		}

		s, err := repository.CreateStatus(ctx, s1)
		if err != nil {
			t.Fatalf("Error al crear Status: %v", err)
		}
		id_status = s.Statusid
	})

	t.Run("Actualizar Status", func(t *testing.T) {
		s_update := db.UpdateStatusParams{
			Statusid:   id_status,
			Statusname: "Visto",
		}

		err = repository.UpdateStatus(ctx, s_update)
		if err != nil {
			t.Fatalf("Error al actualizar Status: %v", err)
		}
	})

	t.Run("Obtener Status y comparar", func(t *testing.T) {
		s, err := repository.GetStatus(ctx, id_status)
		if err != nil {
			t.Fatalf("Error al obtener Status: %v", err)
		}

		if s.Statusname != "Visto" {
			t.Fatalf("El status obtenido no coincide con la actualizacion.")
		}
	})

	t.Run("Borrar Status", func(t *testing.T) {
		err = repository.DeleteStatus(ctx, id_status)
		if err != nil {
			t.Fatalf("Error al borrar Status: %v", err)
		}
	})
}

func TestCRUDCastMember(t *testing.T) {
	connection := OpenDB(t)
	ctx := context.Background()
	repository := db.New(connection)

	var id_castmember int32
	var err error

	t.Run("Crear CastMember", func(t *testing.T) {
		c1 := db.CreateCastMemberParams{
			Castmemberid:     0,
			Castmembername:   "Actor Original",
			Castmemberrole:   "A",
			Castmembergender: "M",
		}

		c, err := repository.CreateCastMember(ctx, c1)
		if err != nil {
			t.Fatalf("Error al crear CastMember: %v", err)
		}
		id_castmember = c.Castmemberid
	})

	t.Run("Actualizar CastMember", func(t *testing.T) {
		c_update := db.UpdateCastMemberParams{
			Castmemberid:     id_castmember,
			Castmembername:   "Actor Editado",
			Castmemberrole:   "D",
			Castmembergender: "M",
		}

		err = repository.UpdateCastMember(ctx, c_update)
		if err != nil {
			t.Fatalf("Error al actualizar CastMember: %v", err)
		}
	})

	t.Run("Obtener CastMember y comparar", func(t *testing.T) {
		c, err := repository.GetCastMember(ctx, id_castmember)
		if err != nil {
			t.Fatalf("Error al obtener CastMember: %v", err)
		}

		if c.Castmembername != "Actor Editado" || c.Castmemberrole != "D" {
			t.Fatalf("Los datos obtenidos no coinciden con la actualizacion.")
		}
	})

	t.Run("Borrar CastMember", func(t *testing.T) {
		err = repository.DeleteCastMember(ctx, id_castmember)
		if err != nil {
			t.Fatalf("Error al borrar CastMember: %v", err)
		}
	})
}