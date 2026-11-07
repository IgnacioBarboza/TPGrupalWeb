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
)

// Interfaz
type GenreRepository interface {
	GetGenres(ctx context.Context) ([]dto.Genre, error)
	Create(ctx context.Context, genre dto.Genre) (dto.Genre, error)
}

type SQLCRepository struct {
	queries *db.Queries
}

// "Constructor"
func NewSQLCRepository(q *db.Queries) *SQLCRepository {
	return &SQLCRepository{queries: q}
}

// A partir de aca, implementaciones:

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
