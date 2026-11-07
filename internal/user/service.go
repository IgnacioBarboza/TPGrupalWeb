/*

Este service.go permite brindar una capa intermedia entre el codigo generado por sqlc y el resto de la logica del negocio,
permitiendo aplicar reglas de negocio y validaciones independiente de la implementacion de la base de datos.

*/

package User

// Imports
import (
	dto "WebGrupal/internal"
	"context"
)

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUsers(ctx context.Context) ([]dto.User, error) {
	return s.repo.GetUsers(ctx)
}

func (s *UserService) Create(ctx context.Context, user dto.User) (dto.User, error) {
	if err := validateUser(user); err != nil {
		return dto.User{}, err
	}
	return s.repo.Create(ctx, user)
}

// validateUser aplica las reglas de negocio de Create.
func validateUser(user dto.User) error {
	if user.Username == "" {
		return ErrInvalidUserNameEmpty
	}
	if user.Userpasswordhashed == "" {
		return ErrInvalidUserPasswordEmpty
	}
	if len(user.Usergenre) != 1 {
		return ErrInvalidUserGenre
	}
	return nil
}
