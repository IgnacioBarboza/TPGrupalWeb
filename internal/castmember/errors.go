package CastMember

import "errors"

var (
	ErrCastMemberNotFound        = errors.New("CastMember no encontrado.")
	ErrInvalidCastMember         = errors.New("Datos del castMember inválidos.")
	ErrMetodoNoPermitido         = errors.New("Método no permitido.")
	ErrInvalidCastMemberID       = errors.New("ID del castMember inválido.")
	ErrInvalidCastMemberNameEmpty = errors.New("El nombre del castMember no puede estar vacío.")
	ErrInvalidCastMemberRole     = errors.New("El rol del castMember debe ser un único carácter.")
	ErrInvalidCastMemberGender   = errors.New("El género del castMember debe ser un único carácter.")
)

