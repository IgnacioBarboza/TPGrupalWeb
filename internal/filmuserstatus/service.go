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

func (s *FilmUserStatusService) GetFilmUserStatus(ctx context.Context, filmID int32, userID int32) (dto.Filmuserstatus, error) {
	return s.repo.GetFilmUserStatus(ctx, filmID, userID)
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

// validateFilmUserStatus aplica las reglas de negocio comunes a Create y Update.
func validateFilmUserStatus(fus dto.Filmuserstatus) error {
	if fus.FilmFilmid <= 0 || fus.UserUserid <= 0 || fus.StatusStatusid <= 0 {
		return ErrInvalidFilmUserStatusIDs
	}
	return nil
}

func (s *FilmUserStatusService) CreateMultipleFilmUserStatuses(ctx context.Context, statuses []dto.Filmuserstatus) ([]dto.Filmuserstatus, error) {
	if len(statuses) == 0 {
		return nil, ErrInvalidFilmUserStatus
	}
	// Valida todos
	for _, fus := range statuses {
		if err := validateFilmUserStatus(fus); err != nil {
			return nil, err
		}
	}

	created := make([]dto.Filmuserstatus, 0, len(statuses))
	for _, fus := range statuses {
		c, err := s.repo.Create(ctx, fus)
		if err != nil {
			return nil, err
		}
		created = append(created, c)
	}
	return created, nil
}

func (s *FilmUserStatusService) Delete(ctx context.Context, filmID int32, userID int32) error {
	return s.repo.Delete(ctx, filmID, userID)
}

func (s *FilmUserStatusService) Update(ctx context.Context, fus dto.Filmuserstatus) (dto.Filmuserstatus, error) {
	if err := validateFilmUserStatus(fus); err != nil {
		return dto.Filmuserstatus{}, err
	}
	return s.repo.Update(ctx, fus)
}
