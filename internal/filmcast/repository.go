/*

Este repository.go permite brindar una capa intermedia entre el codigo generado por sqlc y el resto de la logica del negocio,
permitiendo cambiar en un futuro la implementacion de la base de datos sin afectar el resto del codigo.

*/

package FilmCast

// Imports
import (
	db "WebGrupal/db/sqlc"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// Interfaz
// FilmCast es una tabla de union sin campos propios, por lo que no tiene operacion Update.
type FilmCastRepository interface {
	GetFilmCast(ctx context.Context, filmID int32, castMemberID int32) (db.Filmcast, error)
	GetFilmCasts(ctx context.Context) ([]db.Filmcast, error)
	Create(ctx context.Context, filmCast db.Filmcast) (db.Filmcast, error)
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

func (r *SQLCRepository) GetFilmCast(ctx context.Context, filmID int32, castMemberID int32) (db.Filmcast, error) {
	fc, err := r.queries.GetFilmCast(ctx, db.GetFilmCastParams{
		Filmid:       filmID,
		Castmemberid: castMemberID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.Filmcast{}, ErrFilmCastNotFound
		}
		return db.Filmcast{}, fmt.Errorf("consultar filmcast film=%d cast=%d: %w", filmID, castMemberID, err)
	}
	return fc, nil
}

func (r *SQLCRepository) GetFilmCasts(ctx context.Context) ([]db.Filmcast, error) {
	rows, err := r.queries.GetFilmCasts(ctx)
	if err != nil {
		return nil, err
	}
	filmCasts := make([]db.Filmcast, 0, len(rows))
	for _, fc := range rows {
		filmCasts = append(filmCasts, fc)
	}
	return filmCasts, nil
}

func (r *SQLCRepository) Create(ctx context.Context, filmCast db.Filmcast) (db.Filmcast, error) {
	fc, err := r.queries.CreateFilmCast(ctx, db.CreateFilmCastParams{
		Filmid:       filmCast.Filmid,
		Castmemberid: filmCast.Castmemberid,
	})
	if err != nil {
		return db.Filmcast{}, err
	}
	return fc, nil
}

func (r *SQLCRepository) Delete(ctx context.Context, filmID int32, castMemberID int32) error {
	return r.queries.DeleteFilmCast(ctx, db.DeleteFilmCastParams{
		Filmid:       filmID,
		Castmemberid: castMemberID,
	})
}

