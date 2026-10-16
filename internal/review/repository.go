/*

Este repository.go permite brindar una capa intermedia entre el codigo generado por sqlc y el resto de la logica del negocio,
permitiendo cambiar en un futuro la implementacion de la base de datos sin afectar el resto del codigo.

*/

package Review

// Imports
import (
	db "WebGrupal/db/sqlc"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// Interfaz
type ReviewRepository interface {
	GetReview(ctx context.Context, filmID int32, userID int32) (db.Review, error)
	GetReviews(ctx context.Context) ([]db.Review, error)
	Create(ctx context.Context, review db.Review) (db.Review, error)
	Delete(ctx context.Context, filmID int32, userID int32) error
	Update(ctx context.Context, review db.Review) (db.Review, error)
}

type SQLCRepository struct {
	queries *db.Queries
}

// "Constructor"
func NewSQLCRepository(q *db.Queries) *SQLCRepository {
	return &SQLCRepository{queries: q}
}

// A partir de aca, implementaciones:

func (r *SQLCRepository) GetReview(ctx context.Context, filmID int32, userID int32) (db.Review, error) {
	rev, err := r.queries.GetReview(ctx, db.GetReviewParams{
		Reviewfilmid: filmID,
		Reviewuserid: userID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.Review{}, ErrReviewNotFound
		}
		return db.Review{}, fmt.Errorf("consultar review film=%d user=%d: %w", filmID, userID, err)
	}
	return rev, nil
}

func (r *SQLCRepository) GetReviews(ctx context.Context) ([]db.Review, error) {
	rows, err := r.queries.GetReviews(ctx)
	if err != nil {
		return nil, err
	}
	reviews := make([]db.Review, 0, len(rows))
	for _, rev := range rows {
		reviews = append(reviews, rev)
	}
	return reviews, nil
}

func (r *SQLCRepository) Create(ctx context.Context, review db.Review) (db.Review, error) {
	rev, err := r.queries.CreateReview(ctx, db.CreateReviewParams{
		Reviewstars:       review.Reviewstars,
		Reviewdescription: review.Reviewdescription,
		Reviewuserid:      review.Reviewuserid,
		Reviewfilmid:      review.Reviewfilmid,
	})
	if err != nil {
		return db.Review{}, err
	}
	return rev, nil
}

func (r *SQLCRepository) Delete(ctx context.Context, filmID int32, userID int32) error {
	return r.queries.DeleteReview(ctx, db.DeleteReviewParams{
		Reviewfilmid: filmID,
		Reviewuserid: userID,
	})
}

func (r *SQLCRepository) Update(ctx context.Context, review db.Review) (db.Review, error) {
	err := r.queries.UpdateReview(ctx, db.UpdateReviewParams{
		Reviewfilmid:      review.Reviewfilmid,
		Reviewuserid:      review.Reviewuserid,
		Reviewstars:       review.Reviewstars,
		Reviewdescription: review.Reviewdescription,
	})
	if err != nil {
		return db.Review{}, err
	}
	return r.GetReview(ctx, review.Reviewfilmid, review.Reviewuserid)
}

