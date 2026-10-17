package Status

import "errors"

var (
	ErrStatusNotFound        = errors.New("Status no encontrado.")
	ErrInvalidStatus         = errors.New("Datos del status inválidos.")
	ErrMetodoNoPermitido     = errors.New("Método no permitido.")
	ErrInvalidStatusID       = errors.New("ID del status inválido.")
	ErrInvalidStatusNameEmpty = errors.New("El nombre del status no puede estar vacío.")
)

