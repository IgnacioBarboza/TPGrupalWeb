/*

Este repostory.go Permite brindar una capa intermedia entre el codigo generado por sqlc y el resto de la logica del negocio,
permitiendo cambiar en un futuro la implementacion de la base de datos sin afectar el resto del codigo.

*/

package Film

// Imports
import (
	db "WebGrupal/db/sqlc"
	dto "WebGrupal/internal"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// Interfaz
type FilmRepository interface {
	GetFilm(ctx context.Context, id int32) (dto.Film, error)
	GetFilms(ctx context.Context) ([]dto.Film, error)
	Create(ctx context.Context, film dto.Film) (dto.Film, error)
	Delete(ctx context.Context, id int32) error
	Update(ctx context.Context, film dto.Film) (dto.Film, error)
}

type SQLCRepository struct {
	queries *db.Queries
}

// "Constructor"
func NewSQLCRepository(q *db.Queries) *SQLCRepository {
	return &SQLCRepository{queries: q}
}

//A partir de aca, implementaciones:

func (r *SQLCRepository) GetFilm(ctx context.Context, id int32) (dto.Film, error) {
	f, err := r.queries.GetFilm(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.Film{}, ErrFilmNotFound
		}
		return dto.Film{}, fmt.Errorf("consultar film %d: %w", id, err)
	}
	return dto.Film(f), nil
}

func (r *SQLCRepository) GetFilms(ctx context.Context) ([]dto.Film, error) {
	rows, err := r.queries.GetFilms(ctx)
	if err != nil {
		return nil, err
	}
	films := make([]dto.Film, 0, len(rows))
	for _, f := range rows {
		films = append(films, dto.Film(f))
	}
	return films, nil
}

func (r *SQLCRepository) Create(ctx context.Context, film dto.Film) (dto.Film, error) {
	f, err := r.queries.CreateFilm(ctx, db.CreateFilmParams{
		Filmtitle:       film.Filmtitle,
		Filmreleasedate: film.Filmreleasedate,
		Filmduration:    film.Filmduration,
	})
	if err != nil {
		return dto.Film{}, err
	}
	return dto.Film(f), nil
}

func (r *SQLCRepository) Delete(ctx context.Context, id int32) error {
	return r.queries.DeleteFilm(ctx, id)
}

func (r *SQLCRepository) Update(ctx context.Context, film dto.Film) (dto.Film, error) {
	err := r.queries.UpdateFilm(ctx, db.UpdateFilmParams{
		Filmid:          film.Filmid,
		Filmtitle:       film.Filmtitle,
		Filmreleasedate: film.Filmreleasedate,
		Filmduration:    film.Filmduration,
	})
	if err != nil {
		return dto.Film{}, err
	}
	return r.GetFilm(ctx, film.Filmid)
}
