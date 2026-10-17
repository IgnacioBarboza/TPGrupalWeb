package User

import "errors"

var (
	ErrUserNotFound             = errors.New("Usuario no encontrado.")
	ErrInvalidUser              = errors.New("Datos del usuario inválidos.")
	ErrMetodoNoPermitido        = errors.New("Método no permitido.")
	ErrInvalidUserID            = errors.New("ID del usuario inválido.")
	ErrInvalidUserNameEmpty     = errors.New("El nombre del usuario no puede estar vacío.")
	ErrInvalidUserPasswordEmpty = errors.New("La contraseña del usuario no puede estar vacía.")
	ErrInvalidUserGenre         = errors.New("El género del usuario debe ser un único carácter.")
)

