package User

import "errors"

var (
	ErrInvalidUserNameEmpty     = errors.New("El nombre del usuario no puede estar vacío.")
	ErrInvalidUserPasswordEmpty = errors.New("La contraseña del usuario no puede estar vacía.")
	ErrInvalidUserGenre         = errors.New("El género del usuario debe ser un único carácter.")
)
