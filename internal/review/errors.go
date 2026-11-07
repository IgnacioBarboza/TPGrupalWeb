package Review

import "errors"

var (
	ErrInvalidReviewStars            = errors.New("La cantidad de estrellas debe estar entre 1 y 5.")
	ErrInvalidReviewDescriptionEmpty = errors.New("La descripción de la review no puede estar vacía.")
)
