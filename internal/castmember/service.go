/*

Este service.go permite brindar una capa intermedia entre el codigo generado por sqlc y el resto de la logica del negocio,
permitiendo aplicar reglas de negocio y validaciones independiente de la implementacion de la base de datos.

*/

package CastMember

// Imports
import (
	dto "WebGrupal/internal"
	"context"
)

type CastMemberService struct {
	repo CastMemberRepository
}

func NewCastMemberService(repo CastMemberRepository) *CastMemberService {
	return &CastMemberService{repo: repo}
}

func (s *CastMemberService) GetCastMembers(ctx context.Context) ([]dto.Castmember, error) {
	return s.repo.GetCastMembers(ctx)
}

func (s *CastMemberService) Create(ctx context.Context, castMember dto.Castmember) (dto.Castmember, error) {
	if err := validateCastMember(castMember); err != nil {
		return dto.Castmember{}, err
	}
	return s.repo.Create(ctx, castMember)
}

// validateCastMember aplica las reglas de negocio de Create.
func validateCastMember(castMember dto.Castmember) error {
	if castMember.Castmembername == "" {
		return ErrInvalidCastMemberNameEmpty
	}
	if len(castMember.Castmemberrole) != 1 {
		return ErrInvalidCastMemberRole
	}
	if len(castMember.Castmembergender) != 1 {
		return ErrInvalidCastMemberGender
	}
	return nil
}
