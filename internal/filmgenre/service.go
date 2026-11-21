/*

Este service.go permite brindar una capa intermedia entre el codigo generado por sqlc y el resto de la logica del negocio,
permitiendo aplicar reglas de negocio y validaciones independiente de la implementacion de la base de datos.

*/

package FilmGenre

// Imports
import (
	dto "WebGrupal/internal"
	"context"
)

type FilmGenreService struct {
	repo FilmGenreRepository
}

func NewFilmGenreService(repo FilmGenreRepository) *FilmGenreService {
	return &FilmGenreService{repo: repo}
}

func (s *FilmGenreService) GetFilmGenres(ctx context.Context) ([]dto.Filmgenre, error) {
	return s.repo.GetFilmGenres(ctx)
}

func (s *FilmGenreService) Create(ctx context.Context, filmGenre dto.Filmgenre) (dto.Filmgenre, error) {
	if err := validateFilmGenre(filmGenre); err != nil {
		return dto.Filmgenre{}, err
	}
	return s.repo.Create(ctx, filmGenre)
}

// validateFilmGenre aplica las reglas de negocio comunes a Create.
func validateFilmGenre(filmGenre dto.Filmgenre) error {
	if filmGenre.GenreGenreid <= 0 || filmGenre.FilmFilmid <= 0 {
		return ErrInvalidFilmGenreIDs
	}
	return nil
}

// Delete borra la relación película-género.
func (s *FilmGenreService) Delete(ctx context.Context, genreID, filmID int32) error {
	return s.repo.Delete(ctx, genreID, filmID)
}

// Get devuelve la relación película-género.
func (s *FilmGenreService) Get(ctx context.Context, genreID, filmID int32) (dto.Filmgenre, error) {
	return s.repo.Get(ctx, genreID, filmID)
}

// Update valida el par nuevo con las mismas reglas que Create y reemplaza al par (genreID, filmID).
func (s *FilmGenreService) Update(ctx context.Context, genreID, filmID int32, filmGenre dto.Filmgenre) error {
	if err := validateFilmGenre(filmGenre); err != nil {
		return err
	}
	return s.repo.Update(ctx, genreID, filmID, filmGenre)
}
