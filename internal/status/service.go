/*

Este service.go permite brindar una capa intermedia entre el codigo generado por sqlc y el resto de la logica del negocio,
permitiendo aplicar reglas de negocio y validaciones independiente de la implementacion de la base de datos.

*/

package Status

// Imports
import (
	dto "WebGrupal/internal"
	"context"
)

type StatusService struct {
	repo StatusRepository
}

func NewStatusService(repo StatusRepository) *StatusService {
	return &StatusService{repo: repo}
}

func (s *StatusService) GetStatuses(ctx context.Context) ([]dto.Status, error) {
	return s.repo.GetStatuses(ctx)
}

func (s *StatusService) Create(ctx context.Context, status dto.Status) (dto.Status, error) {
	if err := validateStatus(status); err != nil {
		return dto.Status{}, err
	}
	return s.repo.Create(ctx, status)
}

// validateStatus aplica las reglas de negocio de Create.
func validateStatus(status dto.Status) error {
	if status.Statusname == "" {
		return ErrInvalidStatusNameEmpty
	}
	return nil
}

// Delete borra el estado.
func (s *StatusService) Delete(ctx context.Context, id int32) error {
	return s.repo.Delete(ctx, id)
}

// Get devuelve el estado.
func (s *StatusService) Get(ctx context.Context, id int32) (dto.Status, error) {
	return s.repo.Get(ctx, id)
}

// Update valida el estado con las mismas reglas que Create y guarda los cambios.
func (s *StatusService) Update(ctx context.Context, status dto.Status) error {
	if err := validateStatus(status); err != nil {
		return err
	}
	return s.repo.Update(ctx, status)
}
