/*

Este service.go permite brindar una capa intermedia entre el codigo generado por sqlc y el resto de la logica del negocio,
permitiendo aplicar reglas de negocio y validaciones independiente de la implementacion de la base de datos.

*/

package CastMember

// Imports
import (
	db "WebGrupal/db/sqlc"
	"context"
)

type CastMemberService struct {
	repo CastMemberRepository
}

func NewCastMemberService(repo CastMemberRepository) *CastMemberService {
	return &CastMemberService{repo: repo}
}

func (s *CastMemberService) GetCastMember(ctx context.Context, id int32) (db.Castmember, error) {
	cm, err := s.repo.GetCastMember(ctx, id)

	if err != nil {
		return db.Castmember{}, err
	}

	// EJEMPLO DE REGLA DE NEGOCIO
	if cm.Castmemberid < 0 {
		return db.Castmember{}, ErrInvalidCastMemberID
	}
	return cm, nil
}

func (s *CastMemberService) GetCastMembers(ctx context.Context) ([]db.Castmember, error) {
	return s.repo.GetCastMembers(ctx)
}

func (s *CastMemberService) Create(ctx context.Context, castMember db.Castmember) (db.Castmember, error) {
	if err := validateCastMember(castMember); err != nil {
		return db.Castmember{}, err
	}
	return s.repo.Create(ctx, castMember)
}

// validateCastMember aplica las reglas de negocio comunes a Create y Update.
func validateCastMember(castMember db.Castmember) error {
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

func (s *CastMemberService) CreateMultipleCastMembers(ctx context.Context, castMembers []db.Castmember) ([]db.Castmember, error) {
	if len(castMembers) == 0 {
		return nil, ErrInvalidCastMember
	}
	// Valida todos
	for _, castMember := range castMembers {
		if err := validateCastMember(castMember); err != nil {
			return nil, err
		}
	}

	created := make([]db.Castmember, 0, len(castMembers))
	for _, castMember := range castMembers {
		c, err := s.repo.Create(ctx, castMember)
		if err != nil {
			return nil, err
		}
		created = append(created, c)
	}
	return created, nil
}

func (s *CastMemberService) Delete(ctx context.Context, id int32) error {
	return s.repo.Delete(ctx, id)
}

func (s *CastMemberService) Update(ctx context.Context, castMember db.Castmember) (db.Castmember, error) {
	if err := validateCastMember(castMember); err != nil {
		return db.Castmember{}, err
	}
	return s.repo.Update(ctx, castMember)
}

