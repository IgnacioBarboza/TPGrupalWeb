/*

Este repository.go permite brindar una capa intermedia entre el codigo generado por sqlc y el resto de la logica del negocio,
permitiendo cambiar en un futuro la implementacion de la base de datos sin afectar el resto del codigo.

*/

package CastMember

// Imports
import (
	db "WebGrupal/db/sqlc"
	dto "WebGrupal/internal"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// Interfaz
type CastMemberRepository interface {
	GetCastMember(ctx context.Context, id int32) (dto.Castmember, error)
	GetCastMembers(ctx context.Context) ([]dto.Castmember, error)
	Create(ctx context.Context, castMember dto.Castmember) (dto.Castmember, error)
	Delete(ctx context.Context, id int32) error
	Update(ctx context.Context, castMember dto.Castmember) (dto.Castmember, error)
}

type SQLCRepository struct {
	queries *db.Queries
}

// "Constructor"
func NewSQLCRepository(q *db.Queries) *SQLCRepository {
	return &SQLCRepository{queries: q}
}

// A partir de aca, implementaciones:

func (r *SQLCRepository) GetCastMember(ctx context.Context, id int32) (dto.Castmember, error) {
	c, err := r.queries.GetCastMember(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.Castmember{}, ErrCastMemberNotFound
		}
		return dto.Castmember{}, fmt.Errorf("consultar castmember %d: %w", id, err)
	}
	return dto.Castmember(c), nil
}

func (r *SQLCRepository) GetCastMembers(ctx context.Context) ([]dto.Castmember, error) {
	rows, err := r.queries.GetCastMembers(ctx)
	if err != nil {
		return nil, err
	}
	castMembers := make([]dto.Castmember, 0, len(rows))
	for _, c := range rows {
		castMembers = append(castMembers, dto.Castmember(c))
	}
	return castMembers, nil
}

func (r *SQLCRepository) Create(ctx context.Context, castMember dto.Castmember) (dto.Castmember, error) {
	c, err := r.queries.CreateCastMember(ctx, db.CreateCastMemberParams{
		Castmembername:   castMember.Castmembername,
		Castmemberrole:   castMember.Castmemberrole,
		Castmembergender: castMember.Castmembergender,
	})
	if err != nil {
		return dto.Castmember{}, err
	}
	return dto.Castmember(c), nil
}

func (r *SQLCRepository) Delete(ctx context.Context, id int32) error {
	return r.queries.DeleteCastMember(ctx, id)
}

func (r *SQLCRepository) Update(ctx context.Context, castMember dto.Castmember) (dto.Castmember, error) {
	err := r.queries.UpdateCastMember(ctx, db.UpdateCastMemberParams{
		Castmemberid:     castMember.Castmemberid,
		Castmembername:   castMember.Castmembername,
		Castmemberrole:   castMember.Castmemberrole,
		Castmembergender: castMember.Castmembergender,
	})
	if err != nil {
		return dto.Castmember{}, err
	}
	return r.GetCastMember(ctx, castMember.Castmemberid)
}
