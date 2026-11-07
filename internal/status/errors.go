package Status

import "errors"

var (
	ErrInvalidStatusNameEmpty = errors.New("El nombre del status no puede estar vacío.")
)
