/*

Este service.go permite brindar una capa intermedia entre el codigo generado por sqlc y el resto de la logica del negocio,
permitiendo aplicar reglas de negocio y validaciones independiente de la implementacion de la base de datos.

*/

package Genre

// Imports
import (
	dto "WebGrupal/internal"
	"context"
)

type GenreService struct {
	repo GenreRepository
}

func NewGenreService(repo GenreRepository) *GenreService {
	return &GenreService{repo: repo}
}

func (s *GenreService) GetGenres(ctx context.Context) ([]dto.Genre, error) {
	return s.repo.GetGenres(ctx)
}

func (s *GenreService) Create(ctx context.Context, genre dto.Genre) (dto.Genre, error) {
	if err := validateGenre(genre); err != nil {
		return dto.Genre{}, err
	}
	return s.repo.Create(ctx, genre)
}

// validateGenre aplica las reglas de negocio de Create.
func validateGenre(genre dto.Genre) error {
	if genre.Genrename == "" {
		return ErrInvalidGenreNameEmpty
	}
	return nil
}
