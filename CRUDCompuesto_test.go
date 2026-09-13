package main

import(
	"context"
	"testing"
	"time"
	db "WebGrupal/db/sqlc"
	_ "github.com/lib/pq" //Driver postgresql
)
func TestCRUDReview(t *testing.T) {
	connection := OpenDB(t)

	ctx := context.Background()
	repository := db.New(connection)

	var id_film int32
	var id_user int32

	var err error

	t.Run("Crear Dependencias",func (t *testing.T){
		dateTesting := time.Date(    // 2026 - 01 - 01  00hs 00min 00sec UTC
			2026,
			time.January,
			1,
			0, 0, 0, 0,
			time.UTC,
		)

		u1 := db.CreateUserParams{
			Userid : 0,
			Username: "Test Username",
			Userpasswordhashed: "Test Password",
			Userbirthday: dateTesting,
			Usergenre: "M",
		}

		u,err :=repository.CreateUser(ctx,u1)

		if err != nil {
			t.Fatalf("Error al crear un usuario: %v",err)
		}

		id_user = u.Userid

		f1 := db.CreateFilmParams{
			Filmid:0,
			Filmtitle:"Pelicula Testing",
			Filmreleasedate:dateTesting,
			Filmduration: 120,
		}

		f,err:= repository.CreateFilm(ctx,f1)

		if err != nil {
			t.Fatalf("Error al crear una pelicula: %v",err)
		}

		id_film = f.Filmid
	})

	t.Run("Crear Review", func(t *testing.T) {
		r1 := db.CreateReviewParams{
			Reviewstars:5,
			Reviewdescription:"Esta pelicula es la mejor del mundo.",
			Reviewuserid:id_user,
			Reviewfilmid:id_film,
		}

		_,err = repository.CreateReview(ctx,r1)
		if err != nil {
			t.Fatalf("Error al crear una review: %v",err)
		}
	})

	t.Run("Actualizar Review", func(t *testing.T) {
		review_update := db.UpdateReviewParams{
			Reviewstars:5,
			Reviewdescription:"Esta pelicula actualizada es la mejor del mundo.",
			Reviewuserid:id_user,
			Reviewfilmid:id_film,
		}

		err = repository.UpdateReview(ctx,review_update)

		if err!= nil {
			t.Fatalf("Error al intentar actualizar un review: %v",err)
		}
	})

	t.Run("Obtener Review y comparar", func(t *testing.T) {
		review_get := db.GetReviewParams{
			Reviewfilmid: id_film,
			Reviewuserid: id_user,
		}

		r,err := repository.GetReview(ctx,review_get)

		if err != nil {
			t.Fatalf("Error al obtener el review: %v",err)
		}

		if r.Reviewdescription != "Esta pelicula actualizada es la mejor del mundo."{
			t.Fatalf("El dato obtenido no es el mismo del update.")
		}
	})

	t.Run("Borrado del Review", func(t *testing.T){
		review_delete := db.DeleteReviewParams{
			Reviewfilmid: id_film,
			Reviewuserid: id_user,
		}

		err = repository.DeleteReview(ctx,review_delete)

		if err != nil {
			t.Fatalf("Error al borrar el review: %v",err)
		}
	})

	t.Run("Borrado de las dependencias", func(t *testing.T){
		err = repository.DeleteFilm(ctx,id_film)

		if err != nil {
			t.Fatalf("Error al borrar una pelicula: %v",err)
		}

		err = repository.DeleteUser(ctx,id_user)

		if err != nil {
			t.Fatalf("Error al borrar un usuario: %v",err)
		}
	})
}

func TestCRUDFilmUserStatus(t *testing.T) {
	connection := OpenDB(t)
	ctx := context.Background()
	repository := db.New(connection)

	var id_film, id_user, id_status_1, id_status_2 int32
	var err error

	dateTesting := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)

	t.Run("Crear Dependencias", func(t *testing.T) {
		// 1. Crear User
		u, err := repository.CreateUser(ctx, db.CreateUserParams{
			Username: "Status Tester", 
			Userbirthday: dateTesting, 
			Usergenre: "M",
		})
		if err != nil { 
			t.Fatalf("Error User: %v", err) 
		}
		id_user = u.Userid

		// 2. Crear Film
		f, err := repository.CreateFilm(ctx, db.CreateFilmParams{
			Filmtitle: "Film Status", 
			Filmreleasedate: dateTesting, 
			Filmduration: 120,
		})
		if err != nil { 
			t.Fatalf("Error Film: %v", err) 
		}
		id_film = f.Filmid

		// 3. Crear Status 1 y 2 (para poder probar el Update)
		s1, err := repository.CreateStatus(ctx, db.CreateStatusParams{
			Statusid: 1,
			Statusname: "Viendo",
		})

		if err != nil { 
			t.Fatalf("Error Status 1: %v", err) 
		}
		id_status_1 = s1.Statusid

		s2, err := repository.CreateStatus(ctx, db.CreateStatusParams{
			Statusid: 2,
			Statusname: "Completado",
		})
		if err != nil { 
			t.Fatalf("Error Status 2: %v", err) 
		}
		id_status_2 = s2.Statusid
	})

	t.Run("Crear FilmUserStatus", func(t *testing.T) {
		fus := db.CreateFilmUserStatusParams{
			StatusStatusid: id_status_1,
			FilmFilmid:     id_film,
			UserUserid:     id_user,
		}

		_, err = repository.CreateFilmUserStatus(ctx, fus)
		if err != nil {
			t.Fatalf("Error al crear FilmUserStatus: %v", err)
		}
	})

	t.Run("Actualizar FilmUserStatus", func(t *testing.T) {
		fus_update := db.UpdateFilmUserStatusParams{
			FilmFilmid:     id_film,
			UserUserid:     id_user,
			StatusStatusid: id_status_2, 
		}

		err = repository.UpdateFilmUserStatus(ctx, fus_update)
		if err != nil {
			t.Fatalf("Error al actualizar FilmUserStatus: %v", err)
		}
	})

	t.Run("Obtener FilmUserStatus", func(t *testing.T) {
		fus_get := db.GetFilmUserStatusParams{
			FilmFilmid: id_film,
			UserUserid: id_user,
		}

		result, err := repository.GetFilmUserStatus(ctx, fus_get)
		if err != nil {
			t.Fatalf("Error al obtener FilmUserStatus: %v", err)
		}

		if result.StatusStatusid != id_status_2 {
			t.Fatalf("El estado no se actualizó correctamente")
		}
	})

	t.Run("Borrado general", func(t *testing.T) {
		// 1. Borrar la relación
		err = repository.DeleteFilmUserStatus(ctx, db.DeleteFilmUserStatusParams{
			FilmFilmid: id_film, 
			UserUserid: id_user,
		})
		if err != nil { 
			t.Fatalf("Error al borrar FilmUserStatus: %v", err) 
		}

		// 2. Borrar dependencias
		_ = repository.DeleteFilm(ctx, id_film)
		_ = repository.DeleteUser(ctx, id_user)
		_ = repository.DeleteStatus(ctx, id_status_1)
		_ = repository.DeleteStatus(ctx, id_status_2)
	})
}
func TestCRUDFilmCast(t *testing.T) {
	connection := OpenDB(t)
	ctx := context.Background()
	repository := db.New(connection)

	var id_film, id_castmember int32
	var err error

	dateTesting := time.Date(
		2026,
		time.January, 1, 0, 0, 0, 0, time.UTC)

	t.Run("Crear Dependencias (Film y CastMember)", func(t *testing.T) {
		f, err := repository.CreateFilm(ctx, db.CreateFilmParams{
			Filmtitle: "Pelicula FilmCast", 
			Filmreleasedate: dateTesting, 
			Filmduration: 120,
		})
		if err != nil { 
			t.Fatalf("Error Film: %v", err) 
		}
		id_film = f.Filmid

		c, err := repository.CreateCastMember(ctx, db.CreateCastMemberParams{
			Castmembername: "Actor Test", 
			Castmemberrole: "A", 
			Castmembergender: "M",
		})
		if err != nil { t.Fatalf("Error CastMember: %v", err) }
		id_castmember = c.Castmemberid
	})

	t.Run("Crear FilmCast", func(t *testing.T) {
		fc := db.CreateFilmCastParams{
			Filmid:       id_film,
			Castmemberid: id_castmember,
		}

		_, err = repository.CreateFilmCast(ctx, fc)
		if err != nil {
			t.Fatalf("Error al crear FilmCast: %v", err)
		}
	})

	t.Run("Obtener FilmCast", func(t *testing.T) {
		fc_get := db.GetFilmCastParams{
			Filmid:       id_film,
			Castmemberid: id_castmember,
		}

		_, err := repository.GetFilmCast(ctx, fc_get)
		if err != nil {
			t.Fatalf("Error al obtener FilmCast: %v", err)
		}
	})

	t.Run("Borrado general (FilmCast y Dependencias)", func(t *testing.T) {
		err = repository.DeleteFilmCast(ctx, db.DeleteFilmCastParams{
			Filmid: id_film, 
			Castmemberid: id_castmember,
		})
		if err != nil { 
			t.Fatalf("Error al borrar FilmCast: %v", err) 
		}

		_ = repository.DeleteFilm(ctx, id_film)
		_ = repository.DeleteCastMember(ctx, id_castmember)
	})
}

func TestCRUDFilmGenre(t *testing.T) {
	connection := OpenDB(t)
	ctx := context.Background()
	repository := db.New(connection)

	var id_film, id_genre int32
	var err error

	dateTesting := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)

	t.Run("Crear Dependencias (Film y Genre)", func(t *testing.T) {
		f, err := repository.CreateFilm(ctx, db.CreateFilmParams{
			Filmtitle: "Pelicula FilmGenre", 
			Filmreleasedate: dateTesting, 
			Filmduration: 120,
		})
		if err != nil { 
			t.Fatalf("Error Film: %v", err) 
		}
		id_film = f.Filmid

		g, err := repository.CreateGenre(ctx, db.CreateGenreParams{
			Genrename: "Comedia Test",
		})
		if err != nil { 
			t.Fatalf("Error Genre: %v", err) 
		}
		id_genre = g.Genreid
	})

	t.Run("Crear FilmGenre", func(t *testing.T) {
		fg := db.CreateFilmGenreParams{
			GenreGenreid: id_genre,
			FilmFilmid:   id_film,
		}

		_, err = repository.CreateFilmGenre(ctx, fg)
		if err != nil {
			t.Fatalf("Error al crear FilmGenre: %v", err)
		}
	})

	t.Run("Obtener FilmGenre", func(t *testing.T) {
		fg_get := db.GetFilmGenreParams{
			GenreGenreid: id_genre,
			FilmFilmid:   id_film,
		}

		_, err := repository.GetFilmGenre(ctx, fg_get)
		if err != nil {
			t.Fatalf("Error al obtener FilmGenre: %v", err)
		}
	})

	t.Run("Borrado general (FilmGenre y Dependencias)", func(t *testing.T) {
		err = repository.DeleteFilmGenre(ctx, db.DeleteFilmGenreParams{
			GenreGenreid: id_genre, 
			FilmFilmid: id_film,
		})
		if err != nil { 
			t.Fatalf("Error al borrar FilmGenre: %v", err) 
		}

		_ = repository.DeleteFilm(ctx, id_film)
		_ = repository.DeleteGenre(ctx, id_genre)
	})
}