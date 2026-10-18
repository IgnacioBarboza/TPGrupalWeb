package FilmUserStatus

import "errors"

var (
	ErrFilmUserStatusNotFound   = errors.New("FilmUserStatus no encontrado.")
	ErrInvalidFilmUserStatus    = errors.New("Datos del FilmUserStatus inválidos.")
	ErrMetodoNoPermitido        = errors.New("Método no permitido.")
	ErrInvalidFilmUserStatusIDs = errors.New("El ID del film, del usuario y del estado deben ser mayores a cero.")
)

