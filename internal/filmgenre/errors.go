package FilmGenre

import "errors"

var (
	ErrInvalidFilmGenreIDs = errors.New("El ID del género y el ID del film deben ser mayores a cero.")
)
