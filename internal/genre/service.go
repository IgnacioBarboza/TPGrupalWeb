/*

Este service.go permite brindar una capa intermedia entre el codigo generado por sqlc y el resto de la logica del negocio,
permitiendo aplicar reglas de negocio y validaciones independiente de la implementacion de la base de datos.

*/

package Genre

// Imports
import (
	db "WebGrupal/db/sqlc"
	"context"
)

type GenreService struct {
	repo GenreRepository
}

func NewGenreService(repo GenreRepository) *GenreService {
	return &GenreService{repo: repo}
}

func (s *GenreService) GetGenre(ctx context.Context, id int32) (db.Genre, error) {
	g, err := s.repo.GetGenre(ctx, id)

	if err != nil {
		return db.Genre{}, err
	}

	// EJEMPLO DE REGLA DE NEGOCIO
	if g.Genreid < 0 {
		return db.Genre{}, ErrInvalidGenre
	}
	return g, nil
}

func (s *GenreService) GetGenres(ctx context.Context) ([]db.Genre, error) {
	return s.repo.GetGenres(ctx)
}

func (s *GenreService) Create(ctx context.Context, genre db.Genre) (db.Genre, error) {
	if err := validateGenre(genre); err != nil {
		return db.Genre{}, err
	}
	return s.repo.Create(ctx, genre)
}

// validateGenre aplica las reglas de negocio comunes a Create y Update.
func validateGenre(genre db.Genre) error {
	if genre.Genrename == "" {
		return ErrInvalidGenreNameEmpty
	}
	return nil
}

func (s *GenreService) CreateMultipleGenres(ctx context.Context, genres []db.Genre) ([]db.Genre, error) {
	if len(genres) == 0 {
		return nil, ErrInvalidGenre
	}
	// Valida todos
	for _, genre := range genres {
		if err := validateGenre(genre); err != nil {
			return nil, err
		}
	}

	created := make([]db.Genre, 0, len(genres))
	for _, genre := range genres {
		c, err := s.repo.Create(ctx, genre)
		if err != nil {
			return nil, err
		}
		created = append(created, c)
	}
	return created, nil
}

func (s *GenreService) Delete(ctx context.Context, id int32) error {
	return s.repo.Delete(ctx, id)
}

func (s *GenreService) Update(ctx context.Context, genre db.Genre) (db.Genre, error) {
	if err := validateGenre(genre); err != nil {
		return db.Genre{}, err
	}
	return s.repo.Update(ctx, genre)
}

