/*

Este service.go permite brindar una capa intermedia entre el codigo generado por sqlc y el resto de la logica del negocio,
permitiendo aplicar reglas de negocio y validaciones independiente de la implementacion de la base de datos.

*/

package Review

// Imports
import (
	dto "WebGrupal/internal"
	"context"
)

type ReviewService struct {
	repo ReviewRepository
}

func NewReviewService(repo ReviewRepository) *ReviewService {
	return &ReviewService{repo: repo}
}

func (s *ReviewService) GetReviews(ctx context.Context) ([]dto.Review, error) {
	return s.repo.GetReviews(ctx)
}

func (s *ReviewService) Create(ctx context.Context, review dto.Review) (dto.Review, error) {
	if err := validateReview(review); err != nil {
		return dto.Review{}, err
	}
	return s.repo.Create(ctx, review) // @Nachito aca para crear no se si te pinta buscar que los dos ids existan en la base de datos, sino no tiene sentido crear un review para un film que no existe o un usuario que no existe, pero capaz es mucho quilombo.
}

// validateReview aplica las reglas de negocio de Create.
func validateReview(review dto.Review) error {
	if review.Reviewdescription == "" {
		return ErrInvalidReviewDescriptionEmpty
	}
	if review.Reviewstars < 1 || review.Reviewstars > 5 {
		return ErrInvalidReviewStars
	}
	return nil
}

// Delete borra la reseña.
func (s *ReviewService) Delete(ctx context.Context, filmID, userID int32) error {
	return s.repo.Delete(ctx, filmID, userID)
}

// Get devuelve la reseña.
func (s *ReviewService) Get(ctx context.Context, filmID, userID int32) (dto.Review, error) {
	return s.repo.Get(ctx, filmID, userID)
}

// Update valida la reseña con las mismas reglas que Create y guarda los cambios.
func (s *ReviewService) Update(ctx context.Context, review dto.Review) error {
	if err := validateReview(review); err != nil {
		return err
	}
	return s.repo.Update(ctx, review)
}
