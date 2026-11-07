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
)

// Interfaz
type CastMemberRepository interface {
	GetCastMembers(ctx context.Context) ([]dto.Castmember, error)
	Create(ctx context.Context, castMember dto.Castmember) (dto.Castmember, error)
}

type SQLCRepository struct {
	queries *db.Queries
}

// "Constructor"
func NewSQLCRepository(q *db.Queries) *SQLCRepository {
	return &SQLCRepository{queries: q}
}

// A partir de aca, implementaciones:

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
