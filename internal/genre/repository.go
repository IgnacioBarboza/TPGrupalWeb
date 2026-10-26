/*

Este repository.go permite brindar una capa intermedia entre el codigo generado por sqlc y el resto de la logica del negocio,
permitiendo cambiar en un futuro la implementacion de la base de datos sin afectar el resto del codigo.

*/

package Genre

// Imports
import (
	db "WebGrupal/db/sqlc"
	dto "WebGrupal/internal"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// Interfaz
type GenreRepository interface {
	GetGenre(ctx context.Context, id int32) (dto.Genre, error)
	GetGenres(ctx context.Context) ([]dto.Genre, error)
	Create(ctx context.Context, genre dto.Genre) (dto.Genre, error)
	Delete(ctx context.Context, id int32) error
	Update(ctx context.Context, genre dto.Genre) (dto.Genre, error)
}

type SQLCRepository struct {
	queries *db.Queries
}

// "Constructor"
func NewSQLCRepository(q *db.Queries) *SQLCRepository {
	return &SQLCRepository{queries: q}
}

// A partir de aca, implementaciones:

func (r *SQLCRepository) GetGenre(ctx context.Context, id int32) (dto.Genre, error) {
	row, err := r.queries.GetGenre(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.Genre{}, ErrGenreNotFound
		}
		return dto.Genre{}, fmt.Errorf("consultar género %d: %w", id, err)
	}
	return dto.Genre{Genreid: row.Genreid, Genrename: row.Genrename}, nil
}

func (r *SQLCRepository) GetGenres(ctx context.Context) ([]dto.Genre, error) {
	rows, err := r.queries.GetGenres(ctx)
	if err != nil {
		return nil, err
	}
	genres := make([]dto.Genre, 0, len(rows))
	for _, row := range rows {
		genres = append(genres, dto.Genre{Genreid: row.Genreid, Genrename: row.Genrename})
	}
	return genres, nil
}

func (r *SQLCRepository) Create(ctx context.Context, genre dto.Genre) (dto.Genre, error) {
	row, err := r.queries.CreateGenre(ctx, genre.Genrename)
	if err != nil {
		return dto.Genre{}, err
	}
	return dto.Genre{Genreid: row.Genreid, Genrename: row.Genrename}, nil
}

func (r *SQLCRepository) Delete(ctx context.Context, id int32) error {
	return r.queries.DeleteGenre(ctx, id)
}

func (r *SQLCRepository) Update(ctx context.Context, genre dto.Genre) (dto.Genre, error) {
	err := r.queries.UpdateGenre(ctx, db.UpdateGenreParams{
		Genreid:   genre.Genreid,
		Genrename: genre.Genrename,
	})
	if err != nil {
		return dto.Genre{}, err
	}
	return r.GetGenre(ctx, genre.Genreid)
}
