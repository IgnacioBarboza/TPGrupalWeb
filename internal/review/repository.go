/*

Este repository.go permite brindar una capa intermedia entre el codigo generado por sqlc y el resto de la logica del negocio,
permitiendo cambiar en un futuro la implementacion de la base de datos sin afectar el resto del codigo.

*/

package Review

// Imports
import (
	db "WebGrupal/db/sqlc"
	dto "WebGrupal/internal"
	"context"
)

// Interfaz
type ReviewRepository interface {
	GetReviews(ctx context.Context) ([]dto.Review, error)
	Create(ctx context.Context, review dto.Review) (dto.Review, error)
}

type SQLCRepository struct {
	queries *db.Queries
}

// "Constructor"
func NewSQLCRepository(q *db.Queries) *SQLCRepository {
	return &SQLCRepository{queries: q}
}

// A partir de aca, implementaciones:

func (r *SQLCRepository) GetReviews(ctx context.Context) ([]dto.Review, error) {
	rows, err := r.queries.GetReviews(ctx)
	if err != nil {
		return nil, err
	}
	reviews := make([]dto.Review, 0, len(rows))
	for _, rev := range rows {
		reviews = append(reviews, dto.Review(rev))
	}
	return reviews, nil
}

func (r *SQLCRepository) Create(ctx context.Context, review dto.Review) (dto.Review, error) {
	rev, err := r.queries.CreateReview(ctx, db.CreateReviewParams{
		Reviewstars:       review.Reviewstars,
		Reviewdescription: review.Reviewdescription,
		Reviewuserid:      review.Reviewuserid,
		Reviewfilmid:      review.Reviewfilmid,
	})
	if err != nil {
		return dto.Review{}, err
	}
	return dto.Review(rev), nil
}
