/*

Este repository.go permite brindar una capa intermedia entre el codigo generado por sqlc y el resto de la logica del negocio,
permitiendo cambiar en un futuro la implementacion de la base de datos sin afectar el resto del codigo.

*/

package FilmCast

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
// FilmCast es una tabla de union sin campos propios, por lo que no tiene operacion Update.
type FilmCastRepository interface {
	GetFilmCast(ctx context.Context, filmID int32, castMemberID int32) (dto.Filmcast, error)
	GetFilmCasts(ctx context.Context) ([]dto.Filmcast, error)
	Create(ctx context.Context, filmCast dto.Filmcast) (dto.Filmcast, error)
	Delete(ctx context.Context, filmID int32, castMemberID int32) error
}

type SQLCRepository struct {
	queries *db.Queries
}

// "Constructor"
func NewSQLCRepository(q *db.Queries) *SQLCRepository {
	return &SQLCRepository{queries: q}
}

// A partir de aca, implementaciones:

func (r *SQLCRepository) GetFilmCast(ctx context.Context, filmID int32, castMemberID int32) (dto.Filmcast, error) {
	fc, err := r.queries.GetFilmCast(ctx, db.GetFilmCastParams{
		Filmid:       filmID,
		Castmemberid: castMemberID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.Filmcast{}, ErrFilmCastNotFound
		}
		return dto.Filmcast{}, fmt.Errorf("consultar filmcast film=%d cast=%d: %w", filmID, castMemberID, err)
	}
	return dto.Filmcast(fc), nil
}

func (r *SQLCRepository) GetFilmCasts(ctx context.Context) ([]dto.Filmcast, error) {
	rows, err := r.queries.GetFilmCasts(ctx)
	if err != nil {
		return nil, err
	}
	filmCasts := make([]dto.Filmcast, 0, len(rows))
	for _, fc := range rows {
		filmCasts = append(filmCasts, dto.Filmcast(fc))
	}
	return filmCasts, nil
}

func (r *SQLCRepository) Create(ctx context.Context, filmCast dto.Filmcast) (dto.Filmcast, error) {
	fc, err := r.queries.CreateFilmCast(ctx, db.CreateFilmCastParams{
		Filmid:       filmCast.Filmid,
		Castmemberid: filmCast.Castmemberid,
	})
	if err != nil {
		return dto.Filmcast{}, err
	}
	return dto.Filmcast(fc), nil
}

func (r *SQLCRepository) Delete(ctx context.Context, filmID int32, castMemberID int32) error {
	return r.queries.DeleteFilmCast(ctx, db.DeleteFilmCastParams{
		Filmid:       filmID,
		Castmemberid: castMemberID,
	})
}
