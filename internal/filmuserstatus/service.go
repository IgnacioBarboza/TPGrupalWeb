/*

Este service.go permite brindar una capa intermedia entre el codigo generado por sqlc y el resto de la logica del negocio,
permitiendo aplicar reglas de negocio y validaciones independiente de la implementacion de la base de datos.

*/

package FilmUserStatus

// Imports
import (
	dto "WebGrupal/internal"
	"context"
)

type FilmUserStatusService struct {
	repo FilmUserStatusRepository
}

func NewFilmUserStatusService(repo FilmUserStatusRepository) *FilmUserStatusService {
	return &FilmUserStatusService{repo: repo}
}

func (s *FilmUserStatusService) GetFilmUserStatuses(ctx context.Context) ([]dto.Filmuserstatus, error) {
	return s.repo.GetFilmUserStatuses(ctx)
}

func (s *FilmUserStatusService) Create(ctx context.Context, fus dto.Filmuserstatus) (dto.Filmuserstatus, error) {
	if err := validateFilmUserStatus(fus); err != nil {
		return dto.Filmuserstatus{}, err
	}
	return s.repo.Create(ctx, fus)
}

// validateFilmUserStatus aplica las reglas de negocio de Create.
func validateFilmUserStatus(fus dto.Filmuserstatus) error {
	if fus.FilmFilmid <= 0 || fus.UserUserid <= 0 || fus.StatusStatusid <= 0 {
		return ErrInvalidFilmUserStatusIDs
	}
	return nil
}
