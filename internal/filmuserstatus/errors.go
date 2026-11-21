package FilmUserStatus

import "errors"

var (
	ErrInvalidFilmUserStatusIDs = errors.New("El ID del film, del usuario y del estado deben ser mayores a cero.")
	ErrFilmUserStatusNotFound   = errors.New("Ese usuario no tiene estado para esa película.")
)
