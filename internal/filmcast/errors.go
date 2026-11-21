package FilmCast

import "errors"

var (
	ErrInvalidFilmCastIDs = errors.New("El ID del film y el ID del miembro de reparto deben ser mayores a cero.")
	ErrFilmCastNotFound   = errors.New("Esa película no tiene a ese miembro en el elenco.")
)
