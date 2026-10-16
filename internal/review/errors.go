package Review

import "errors"

var (
	ErrReviewNotFound              = errors.New("Review no encontrada.")
	ErrInvalidReview               = errors.New("Datos de la review inválidos.")
	ErrMetodoNoPermitido           = errors.New("Método no permitido.")
	ErrInvalidReviewStars          = errors.New("La cantidad de estrellas debe estar entre 1 y 5.")
	ErrInvalidReviewDescriptionEmpty = errors.New("La descripción de la review no puede estar vacía.")
)

