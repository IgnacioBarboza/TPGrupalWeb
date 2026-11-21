/*

Este repository.go permite brindar una capa intermedia entre el codigo generado por sqlc y el resto de la logica del negocio,
permitiendo cambiar en un futuro la implementacion de la base de datos sin afectar el resto del codigo.

*/

package Status

// Imports
import (
	db "WebGrupal/db/sqlc"
	dto "WebGrupal/internal"
	"context"
	"database/sql"
	"errors"
)

// Interfaz
type StatusRepository interface {
	GetStatuses(ctx context.Context) ([]dto.Status, error)
	Create(ctx context.Context, status dto.Status) (dto.Status, error)
	Delete(ctx context.Context, id int32) error
	Get(ctx context.Context, id int32) (dto.Status, error)
	Update(ctx context.Context, status dto.Status) error
}

type SQLCRepository struct {
	queries *db.Queries
}

// "Constructor"
func NewSQLCRepository(q *db.Queries) *SQLCRepository {
	return &SQLCRepository{queries: q}
}

// A partir de aca, implementaciones:

func (r *SQLCRepository) GetStatuses(ctx context.Context) ([]dto.Status, error) {
	rows, err := r.queries.GetStatuses(ctx)
	if err != nil {
		return nil, err
	}
	statuses := make([]dto.Status, 0, len(rows))
	for _, s := range rows {
		statuses = append(statuses, dto.Status(s))
	}
	return statuses, nil
}

func (r *SQLCRepository) Create(ctx context.Context, status dto.Status) (dto.Status, error) {
	s, err := r.queries.CreateStatus(ctx, status.Statusname)
	if err != nil {
		return dto.Status{}, err
	}
	return dto.Status(s), nil
}

// Delete borra el estado con la query DeleteStatus de sqlc.
func (r *SQLCRepository) Delete(ctx context.Context, id int32) error {
	return r.queries.DeleteStatus(ctx, id)
}

// Get devuelve el estado, o ErrStatusNotFound si no existe.
func (r *SQLCRepository) Get(ctx context.Context, id int32) (dto.Status, error) {
	s, err := r.queries.GetStatus(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return dto.Status{}, ErrStatusNotFound
	}
	if err != nil {
		return dto.Status{}, err
	}
	return dto.Status(s), nil
}

// Update guarda los cambios de el estado con la query UpdateStatus de sqlc.
func (r *SQLCRepository) Update(ctx context.Context, status dto.Status) error {
	return r.queries.UpdateStatus(ctx, db.UpdateStatusParams{
		Statusid:   status.Statusid,
		Statusname: status.Statusname,
	})
}
