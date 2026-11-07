package Film

import "errors"

var (
	ErrInvalidFilmTitleEmpty = errors.New("El título del film no puede estar vacío.")
	ErrInvalidFilmDuration   = errors.New("La duración del film debe ser mayor a cero.")
)
