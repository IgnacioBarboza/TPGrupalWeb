/*

Este service.go permite brindar una capa intermedia entre el codigo generado por sqlc y el resto de la logica del negocio,
permitiendo aplicar reglas de negocio y validaciones independiente de la implementacion de la base de datos.

*/

package Status

// Imports
import (
	db "WebGrupal/db/sqlc"
	"context"
)

type StatusService struct {
	repo StatusRepository
}

func NewStatusService(repo StatusRepository) *StatusService {
	return &StatusService{repo: repo}
}

func (s *StatusService) GetStatus(ctx context.Context, id int32) (db.Status, error) {
	st, err := s.repo.GetStatus(ctx, id)

	if err != nil {
		return db.Status{}, err
	}

	// EJEMPLO DE REGLA DE NEGOCIO
	if st.Statusid < 0 {
		return db.Status{}, ErrInvalidStatusID
	}
	return st, nil
}

func (s *StatusService) GetStatuses(ctx context.Context) ([]db.Status, error) {
	return s.repo.GetStatuses(ctx)
}

func (s *StatusService) Create(ctx context.Context, status db.Status) (db.Status, error) {
	if err := validateStatus(status); err != nil {
		return db.Status{}, err
	}
	return s.repo.Create(ctx, status)
}

// validateStatus aplica las reglas de negocio comunes a Create y Update.
func validateStatus(status db.Status) error {
	if status.Statusname == "" {
		return ErrInvalidStatusNameEmpty
	}
	return nil
}

func (s *StatusService) CreateMultipleStatuses(ctx context.Context, statuses []db.Status) ([]db.Status, error) {
	if len(statuses) == 0 {
		return nil, ErrInvalidStatus
	}
	// Valida todos
	for _, status := range statuses {
		if err := validateStatus(status); err != nil {
			return nil, err
		}
	}

	created := make([]db.Status, 0, len(statuses))
	for _, status := range statuses {
		c, err := s.repo.Create(ctx, status)
		if err != nil {
			return nil, err
		}
		created = append(created, c)
	}
	return created, nil
}

func (s *StatusService) Delete(ctx context.Context, id int32) error {
	return s.repo.Delete(ctx, id)
}

func (s *StatusService) Update(ctx context.Context, status db.Status) (db.Status, error) {
	if err := validateStatus(status); err != nil {
		return db.Status{}, err
	}
	return s.repo.Update(ctx, status)
}

