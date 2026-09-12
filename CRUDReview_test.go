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