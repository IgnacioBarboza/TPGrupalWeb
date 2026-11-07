/*

Este repository.go permite brindar una capa intermedia entre el codigo generado por sqlc y el resto de la logica del negocio,
permitiendo cambiar en un futuro la implementacion de la base de datos sin afectar el resto del codigo.

*/

package FilmUserStatus

// Imports
import (
	db "WebGrupal/db/sqlc"
	dto "WebGrupal/internal"
	"context"
)

// Interfaz
type FilmUserStatusRepository interface {
	GetFilmUserStatuses(ctx context.Context) ([]dto.Filmuserstatus, error)
	Create(ctx context.Context, fus dto.Filmuserstatus) (dto.Filmuserstatus, error)
}

type SQLCRepository struct {
	queries *db.Queries
}

// "Constructor"
func NewSQLCRepository(q *db.Queries) *SQLCRepository {
	return &SQLCRepository{queries: q}
}

// A partir de aca, implementaciones:

func (r *SQLCRepository) GetFilmUserStatuses(ctx context.Context) ([]dto.Filmuserstatus, error) {
	rows, err := r.queries.GetFilmUserStatuses(ctx)
	if err != nil {
		return nil, err
	}
	statuses := make([]dto.Filmuserstatus, 0, len(rows))
	for _, fus := range rows {
		statuses = append(statuses, dto.Filmuserstatus(fus))
	}
	return statuses, nil
}

func (r *SQLCRepository) Create(ctx context.Context, fus dto.Filmuserstatus) (dto.Filmuserstatus, error) {
	created, err := r.queries.CreateFilmUserStatus(ctx, db.CreateFilmUserStatusParams{
		StatusStatusid: fus.StatusStatusid,
		FilmFilmid:     fus.FilmFilmid,
		UserUserid:     fus.UserUserid,
	})
	if err != nil {
		return dto.Filmuserstatus{}, err
	}
	return dto.Filmuserstatus(created), nil
}
