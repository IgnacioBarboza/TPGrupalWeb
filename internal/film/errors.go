package Film

import "errors"

var (
	ErrFilmNotFound          = errors.New("Film no encontrado.")
	ErrInvalidFilm           = errors.New("Datos del film inválidos.")
	ErrMetodoNoPermitido     = errors.New("Método no permitido.")
	ErrInvalidFilmID         = errors.New("ID del film inválido.")
	ErrInvalidFilmTitleEmpty = errors.New("El título del film no puede estar vacío.")
	ErrInvalidFilmDuration   = errors.New("La duración del film debe ser mayor a cero.")
)
