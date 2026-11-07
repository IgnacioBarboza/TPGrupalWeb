/*

Este repository.go permite brindar una capa intermedia entre el codigo generado por sqlc y el resto de la logica del negocio,
permitiendo cambiar en un futuro la implementacion de la base de datos sin afectar el resto del codigo.

*/

package FilmGenre

// Imports
import (
	db "WebGrupal/db/sqlc"
	dto "WebGrupal/internal"
	"context"
)

// Interfaz
// FilmGenre es una tabla de union sin campos propios, por lo que no tiene operacion Update.
type FilmGenreRepository interface {
	GetFilmGenres(ctx context.Context) ([]dto.Filmgenre, error)
	Create(ctx context.Context, filmGenre dto.Filmgenre) (dto.Filmgenre, error)
}

type SQLCRepository struct {
	queries *db.Queries
}

// "Constructor"
func NewSQLCRepository(q *db.Queries) *SQLCRepository {
	return &SQLCRepository{queries: q}
}

// A partir de aca, implementaciones:

func (r *SQLCRepository) GetFilmGenres(ctx context.Context) ([]dto.Filmgenre, error) {
	rows, err := r.queries.GetFilmGenres(ctx)
	if err != nil {
		return nil, err
	}
	filmGenres := make([]dto.Filmgenre, 0, len(rows))
	for _, fg := range rows {
		filmGenres = append(filmGenres, dto.Filmgenre(fg))
	}
	return filmGenres, nil
}

func (r *SQLCRepository) Create(ctx context.Context, filmGenre dto.Filmgenre) (dto.Filmgenre, error) {
	fg, err := r.queries.CreateFilmGenre(ctx, db.CreateFilmGenreParams{
		GenreGenreid: filmGenre.GenreGenreid,
		FilmFilmid:   filmGenre.FilmFilmid,
	})
	if err != nil {
		return dto.Filmgenre{}, err
	}
	return dto.Filmgenre(fg), nil
}
