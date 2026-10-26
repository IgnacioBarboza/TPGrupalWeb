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

func (s *ReviewService) GetReview(ctx context.Context, filmID int32, userID int32) (dto.Review, error) {
	return s.repo.GetReview(ctx, filmID, userID)
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

func (s *ReviewService) CreateMultipleReviews(ctx context.Context, reviews []dto.Review) ([]dto.Review, error) {
	if len(reviews) == 0 {
		return nil, ErrInvalidReview
	}
	// Valida todas
	for _, review := range reviews {
		if err := validateReview(review); err != nil {
			return nil, err
		}
	}

	created := make([]dto.Review, 0, len(reviews))
	for _, review := range reviews {
		c, err := s.repo.Create(ctx, review)
		if err != nil {
			return nil, err
		}
		created = append(created, c)
	}
	return created, nil
}

func (s *ReviewService) Delete(ctx context.Context, filmID int32, userID int32) error {
	return s.repo.Delete(ctx, filmID, userID)
}

func (s *ReviewService) Update(ctx context.Context, review dto.Review) (dto.Review, error) {
	if err := validateReview(review); err != nil {
		return dto.Review{}, err
	}
	return s.repo.Update(ctx, review)
}

// validateReview aplica las reglas de negocio comunes a Create y Update.
func validateReview(review dto.Review) error {
	if review.Reviewdescription == "" {
		return ErrInvalidReviewDescriptionEmpty
	}
	if review.Reviewstars < 1 || review.Reviewstars > 5 {
		return ErrInvalidReviewStars
	}
	return nil
}
