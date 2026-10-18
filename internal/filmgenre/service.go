/*

Este service.go permite brindar una capa intermedia entre el codigo generado por sqlc y el resto de la logica del negocio,
permitiendo aplicar reglas de negocio y validaciones independiente de la implementacion de la base de datos.

*/

package FilmGenre

// Imports
import (
	db "WebGrupal/db/sqlc"
	"context"
)

type FilmGenreService struct {
	repo FilmGenreRepository
}

func NewFilmGenreService(repo FilmGenreRepository) *FilmGenreService {
	return &FilmGenreService{repo: repo}
}

func (s *FilmGenreService) GetFilmGenre(ctx context.Context, genreID int32, filmID int32) (db.Filmgenre, error) {
	return s.repo.GetFilmGenre(ctx, genreID, filmID)
}

func (s *FilmGenreService) GetFilmGenres(ctx context.Context) ([]db.Filmgenre, error) {
	return s.repo.GetFilmGenres(ctx)
}

func (s *FilmGenreService) Create(ctx context.Context, filmGenre db.Filmgenre) (db.Filmgenre, error) {
	if err := validateFilmGenre(filmGenre); err != nil {
		return db.Filmgenre{}, err
	}
	return s.repo.Create(ctx, filmGenre)
}

// validateFilmGenre aplica las reglas de negocio comunes a Create.
func validateFilmGenre(filmGenre db.Filmgenre) error {
	if filmGenre.GenreGenreid <= 0 || filmGenre.FilmFilmid <= 0 {
		return ErrInvalidFilmGenreIDs
	}
	return nil
}

func (s *FilmGenreService) CreateMultipleFilmGenres(ctx context.Context, filmGenres []db.Filmgenre) ([]db.Filmgenre, error) {
	if len(filmGenres) == 0 {
		return nil, ErrInvalidFilmGenre
	}
	// Valida todos
	for _, fg := range filmGenres {
		if err := validateFilmGenre(fg); err != nil {
			return nil, err
		}
	}

	created := make([]db.Filmgenre, 0, len(filmGenres))
	for _, fg := range filmGenres {
		c, err := s.repo.Create(ctx, fg)
		if err != nil {
			return nil, err
		}
		created = append(created, c)
	}
	return created, nil
}

func (s *FilmGenreService) Delete(ctx context.Context, genreID int32, filmID int32) error {
	return s.repo.Delete(ctx, genreID, filmID)
}

