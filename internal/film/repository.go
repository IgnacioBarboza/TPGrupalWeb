/*

Este repostory.go Permite brindar una capa intermedia entre el codigo generado por sqlc y el resto de la logica del negocio,
permitiendo cambiar en un futuro la implementacion de la base de datos sin afectar el resto del codigo.

*/

package Film

// Imports
import (
	db "WebGrupal/db/sqlc"
	dto "WebGrupal/internal"
	"context"
)

// Interfaz
type FilmRepository interface {
	GetFilms(ctx context.Context) ([]dto.Film, error)
	Create(ctx context.Context, film dto.Film) (dto.Film, error)
}

type SQLCRepository struct {
	queries *db.Queries
}

// "Constructor"
func NewSQLCRepository(q *db.Queries) *SQLCRepository {
	return &SQLCRepository{queries: q}
}

//A partir de aca, implementaciones:

func (r *SQLCRepository) GetFilms(ctx context.Context) ([]dto.Film, error) {
	rows, err := r.queries.GetFilms(ctx)
	if err != nil {
		return nil, err
	}
	films := make([]dto.Film, 0, len(rows))
	for _, f := range rows {
		films = append(films, dto.Film(f))
	}
	return films, nil
}

func (r *SQLCRepository) Create(ctx context.Context, film dto.Film) (dto.Film, error) {
	f, err := r.queries.CreateFilm(ctx, db.CreateFilmParams{
		Filmtitle:       film.Filmtitle,
		Filmreleasedate: film.Filmreleasedate,
		Filmduration:    film.Filmduration,
	})
	if err != nil {
		return dto.Film{}, err
	}
	return dto.Film(f), nil
}
