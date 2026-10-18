/*

Este repository.go permite brindar una capa intermedia entre el codigo generado por sqlc y el resto de la logica del negocio,
permitiendo cambiar en un futuro la implementacion de la base de datos sin afectar el resto del codigo.

*/

package FilmGenre

// Imports
import (
	db "WebGrupal/db/sqlc"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// Interfaz
// FilmGenre es una tabla de union sin campos propios, por lo que no tiene operacion Update.
type FilmGenreRepository interface {
	GetFilmGenre(ctx context.Context, genreID int32, filmID int32) (db.Filmgenre, error)
	GetFilmGenres(ctx context.Context) ([]db.Filmgenre, error)
	Create(ctx context.Context, filmGenre db.Filmgenre) (db.Filmgenre, error)
	Delete(ctx context.Context, genreID int32, filmID int32) error
}

type SQLCRepository struct {
	queries *db.Queries
}

// "Constructor"
func NewSQLCRepository(q *db.Queries) *SQLCRepository {
	return &SQLCRepository{queries: q}
}

// A partir de aca, implementaciones:

func (r *SQLCRepository) GetFilmGenre(ctx context.Context, genreID int32, filmID int32) (db.Filmgenre, error) {
	fg, err := r.queries.GetFilmGenre(ctx, db.GetFilmGenreParams{
		GenreGenreid: genreID,
		FilmFilmid:   filmID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.Filmgenre{}, ErrFilmGenreNotFound
		}
		return db.Filmgenre{}, fmt.Errorf("consultar filmgenre genre=%d film=%d: %w", genreID, filmID, err)
	}
	return fg, nil
}

func (r *SQLCRepository) GetFilmGenres(ctx context.Context) ([]db.Filmgenre, error) {
	rows, err := r.queries.GetFilmGenres(ctx)
	if err != nil {
		return nil, err
	}
	filmGenres := make([]db.Filmgenre, 0, len(rows))
	for _, fg := range rows {
		filmGenres = append(filmGenres, fg)
	}
	return filmGenres, nil
}

func (r *SQLCRepository) Create(ctx context.Context, filmGenre db.Filmgenre) (db.Filmgenre, error) {
	fg, err := r.queries.CreateFilmGenre(ctx, db.CreateFilmGenreParams{
		GenreGenreid: filmGenre.GenreGenreid,
		FilmFilmid:   filmGenre.FilmFilmid,
	})
	if err != nil {
		return db.Filmgenre{}, err
	}
	return fg, nil
}

func (r *SQLCRepository) Delete(ctx context.Context, genreID int32, filmID int32) error {
	return r.queries.DeleteFilmGenre(ctx, db.DeleteFilmGenreParams{
		GenreGenreid: genreID,
		FilmFilmid:   filmID,
	})
}

