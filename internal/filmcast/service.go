/*

Este service.go permite brindar una capa intermedia entre el codigo generado por sqlc y el resto de la logica del negocio,
permitiendo aplicar reglas de negocio y validaciones independiente de la implementacion de la base de datos.

*/

package FilmCast

// Imports
import (
	dto "WebGrupal/internal"
	"context"
)

type FilmCastService struct {
	repo FilmCastRepository
}

func NewFilmCastService(repo FilmCastRepository) *FilmCastService {
	return &FilmCastService{repo: repo}
}

func (s *FilmCastService) GetFilmCast(ctx context.Context, filmID int32, castMemberID int32) (dto.Filmcast, error) {
	return s.repo.GetFilmCast(ctx, filmID, castMemberID)
}

func (s *FilmCastService) GetFilmCasts(ctx context.Context) ([]dto.Filmcast, error) {
	return s.repo.GetFilmCasts(ctx)
}

func (s *FilmCastService) Create(ctx context.Context, filmCast dto.Filmcast) (dto.Filmcast, error) {
	if err := validateFilmCast(filmCast); err != nil {
		return dto.Filmcast{}, err
	}
	return s.repo.Create(ctx, filmCast)
}

// validateFilmCast aplica las reglas de negocio comunes a Create.
func validateFilmCast(filmCast dto.Filmcast) error {
	if filmCast.Filmid <= 0 || filmCast.Castmemberid <= 0 {
		return ErrInvalidFilmCastIDs
	}
	return nil
}

func (s *FilmCastService) CreateMultipleFilmCasts(ctx context.Context, filmCasts []dto.Filmcast) ([]dto.Filmcast, error) {
	if len(filmCasts) == 0 {
		return nil, ErrInvalidFilmCast
	}
	// Valida todos
	for _, fc := range filmCasts {
		if err := validateFilmCast(fc); err != nil {
			return nil, err
		}
	}

	created := make([]dto.Filmcast, 0, len(filmCasts))
	for _, fc := range filmCasts {
		c, err := s.repo.Create(ctx, fc)
		if err != nil {
			return nil, err
		}
		created = append(created, c)
	}
	return created, nil
}

func (s *FilmCastService) Delete(ctx context.Context, filmID int32, castMemberID int32) error {
	return s.repo.Delete(ctx, filmID, castMemberID)
}
