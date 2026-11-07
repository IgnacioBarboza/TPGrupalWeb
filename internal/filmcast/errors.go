package FilmCast

import "errors"

var (
	ErrInvalidFilmCastIDs = errors.New("El ID del film y el ID del miembro de reparto deben ser mayores a cero.")
)
