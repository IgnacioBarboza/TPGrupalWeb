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

func (s *FilmService) GetFilm(ctx context.Context, id int32) (dto.Film, error) {
	ff, err := s.repo.GetFilm(ctx, id)

	if err != nil {
		return dto.Film{}, err
	}

	// EJEMPLO DE REGLA DE NEGOCIO
	if ff.Filmid < 0 {
		return dto.Film{}, ErrInvalidFilmID
	}
	return ff, nil
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

// validateFilm aplica las reglas de negocio comunes a Create y Update.
func validateFilm(film dto.Film) error {
	if film.Filmtitle == "" {
		return ErrInvalidFilmTitleEmpty
	}
	if film.Filmduration <= 0 {
		return ErrInvalidFilmDuration
	}
	return nil
}

func (s *FilmService) CreateMultipleFilms(ctx context.Context, films []dto.Film) ([]dto.Film, error) {
	if len(films) == 0 {
		return nil, ErrInvalidFilmTitleEmpty
	}
	// Valida todos
	for _, film := range films {
		if err := validateFilm(film); err != nil {
			return nil, err
		}
	}

	created := make([]dto.Film, 0, len(films))
	for _, film := range films {
		c, err := s.repo.Create(ctx, film)
		if err != nil {
			return nil, err
		}
		created = append(created, c)
	}
	return created, nil
}

func (s *FilmService) Delete(ctx context.Context, id int32) error {
	return s.repo.Delete(ctx, id)
}

func (s *FilmService) Update(ctx context.Context, film dto.Film) (dto.Film, error) {
	validationErr := validateFilm(film)
	if validationErr != nil {
		return dto.Film{}, validationErr
	}
	return s.repo.Update(ctx, film)
}
