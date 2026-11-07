/*

Este service.go Permite brindar una capa intermedia entre el codigo generado por sqlc y el resto de la logica del negocio,
permitiendo aplicar reglas de negocio y validaciones independiente de la implementacion de la base de datos.

*/

package Film

// Imports
import (
	dto "WebGrupal/internal"
	"context"
)

type FilmService struct {
	repo FilmRepository
}

func NewFilmService(repo FilmRepository) *FilmService {
	return &FilmService{repo: repo}
}

func (s *FilmService) GetFilms(ctx context.Context) ([]dto.Film, error) {
	return s.repo.GetFilms(ctx)
}

func (s *FilmService) Create(ctx context.Context, film dto.Film) (dto.Film, error) {
	validationErr := validateFilm(film)
	if validationErr != nil {
		return dto.Film{}, validationErr
	}
	return s.repo.Create(ctx, film)
}

// validateFilm aplica las reglas de negocio de Create.
func validateFilm(film dto.Film) error {
	if film.Filmtitle == "" {
		return ErrInvalidFilmTitleEmpty
	}
	if film.Filmduration <= 0 {
		return ErrInvalidFilmDuration
	}
	return nil
}
