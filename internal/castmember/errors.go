package CastMember

import "errors"

var (
	ErrInvalidCastMemberNameEmpty = errors.New("El nombre del castMember no puede estar vacío.")
	ErrInvalidCastMemberRole      = errors.New("El rol del castMember debe ser un único carácter.")
	ErrInvalidCastMemberGender    = errors.New("El género del castMember debe ser un único carácter.")
)
