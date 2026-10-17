/*

Este repository.go permite brindar una capa intermedia entre el codigo generado por sqlc y el resto de la logica del negocio,
permitiendo cambiar en un futuro la implementacion de la base de datos sin afectar el resto del codigo.

*/

package Status

// Imports
import (
	db "WebGrupal/db/sqlc"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// Interfaz
type StatusRepository interface {
	GetStatus(ctx context.Context, id int32) (db.Status, error)
	GetStatuses(ctx context.Context) ([]db.Status, error)
	Create(ctx context.Context, status db.Status) (db.Status, error)
	Delete(ctx context.Context, id int32) error
	Update(ctx context.Context, status db.Status) (db.Status, error)
}

type SQLCRepository struct {
	queries *db.Queries
}

// "Constructor"
func NewSQLCRepository(q *db.Queries) *SQLCRepository {
	return &SQLCRepository{queries: q}
}

// A partir de aca, implementaciones:

func (r *SQLCRepository) GetStatus(ctx context.Context, id int32) (db.Status, error) {
	s, err := r.queries.GetStatus(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.Status{}, ErrStatusNotFound
		}
		return db.Status{}, fmt.Errorf("consultar status %d: %w", id, err)
	}
	return s, nil
}

func (r *SQLCRepository) GetStatuses(ctx context.Context) ([]db.Status, error) {
	rows, err := r.queries.GetStatuses(ctx)
	if err != nil {
		return nil, err
	}
	statuses := make([]db.Status, 0, len(rows))
	for _, s := range rows {
		statuses = append(statuses, s)
	}
	return statuses, nil
}

func (r *SQLCRepository) Create(ctx context.Context, status db.Status) (db.Status, error) {
	s, err := r.queries.CreateStatus(ctx, status.Statusname)
	if err != nil {
		return db.Status{}, err
	}
	return s, nil
}

func (r *SQLCRepository) Delete(ctx context.Context, id int32) error {
	return r.queries.DeleteStatus(ctx, id)
}

func (r *SQLCRepository) Update(ctx context.Context, status db.Status) (db.Status, error) {
	err := r.queries.UpdateStatus(ctx, db.UpdateStatusParams{
		Statusid:   status.Statusid,
		Statusname: status.Statusname,
	})
	if err != nil {
		return db.Status{}, err
	}
	return r.GetStatus(ctx, status.Statusid)
}

