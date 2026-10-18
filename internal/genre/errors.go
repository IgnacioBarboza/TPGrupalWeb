package Genre

import "errors"

var (
	ErrGenreNotFound         = errors.New("Género no encontrado.")
	ErrInvalidGenre          = errors.New("Datos del género inválidos.")
	ErrMetodoNoPermitido     = errors.New("Método no permitido.")
	ErrInvalidGenreNameEmpty = errors.New("El nombre del género no puede estar vacío.")
)

