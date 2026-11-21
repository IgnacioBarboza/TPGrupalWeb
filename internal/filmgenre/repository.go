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
	"database/sql"
	"errors"
)

// Interfaz
// FilmGenre es una tabla de union sin campos propios, por lo que no tiene operacion Update.
type FilmGenreRepository interface {
	GetFilmGenres(ctx context.Context) ([]dto.Filmgenre, error)
	Create(ctx context.Context, filmGenre dto.Filmgenre) (dto.Filmgenre, error)
	Delete(ctx context.Context, genreID, filmID int32) error
	Get(ctx context.Context, genreID, filmID int32) (dto.Filmgenre, error)
	Update(ctx context.Context, genreID, filmID int32, filmGenre dto.Filmgenre) error
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

// Delete borra la relación película-género con la query DeleteFilmGenre de sqlc.
func (r *SQLCRepository) Delete(ctx context.Context, genreID, filmID int32) error {
	return r.queries.DeleteFilmGenre(ctx, db.DeleteFilmGenreParams{GenreGenreid: genreID, FilmFilmid: filmID})
}

// Get devuelve la relación película-género, o ErrFilmGenreNotFound si no existe.
func (r *SQLCRepository) Get(ctx context.Context, genreID, filmID int32) (dto.Filmgenre, error) {
	row, err := r.queries.GetFilmGenre(ctx, db.GetFilmGenreParams{GenreGenreid: genreID, FilmFilmid: filmID})
	if errors.Is(err, sql.ErrNoRows) {
		return dto.Filmgenre{}, ErrFilmGenreNotFound
	}
	if err != nil {
		return dto.Filmgenre{}, err
	}
	return dto.Filmgenre(row), nil
}

// Update reemplaza el par (genreID, filmID) por el de filmGenre, con la query UpdateFilmGenre de sqlc.
func (r *SQLCRepository) Update(ctx context.Context, genreID, filmID int32, filmGenre dto.Filmgenre) error {
	return r.queries.UpdateFilmGenre(ctx, db.UpdateFilmGenreParams{
		NewGenreID: filmGenre.GenreGenreid,
		NewFilmID:  filmGenre.FilmFilmid,
		GenreID:    genreID,
		FilmID:     filmID,
	})
}
