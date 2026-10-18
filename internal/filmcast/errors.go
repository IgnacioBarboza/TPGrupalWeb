package FilmCast

import "errors"

var (
	ErrFilmCastNotFound   = errors.New("FilmCast no encontrado.")
	ErrInvalidFilmCast    = errors.New("Datos del FilmCast inválidos.")
	ErrMetodoNoPermitido  = errors.New("Método no permitido.")
	ErrInvalidFilmCastIDs = errors.New("El ID del film y el ID del miembro de reparto deben ser mayores a cero.")
)

