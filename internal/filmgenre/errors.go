package FilmGenre

import "errors"

var (
	ErrFilmGenreNotFound   = errors.New("FilmGenre no encontrado.")
	ErrInvalidFilmGenre    = errors.New("Datos del FilmGenre inválidos.")
	ErrMetodoNoPermitido   = errors.New("Método no permitido.")
	ErrInvalidFilmGenreIDs = errors.New("El ID del género y el ID del film deben ser mayores a cero.")
)

