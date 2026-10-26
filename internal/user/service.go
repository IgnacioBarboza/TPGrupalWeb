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

func (s *UserService) GetUser(ctx context.Context, id int32) (dto.User, error) {
	u, err := s.repo.GetUser(ctx, id)

	if err != nil {
		return dto.User{}, err
	}

	// EJEMPLO DE REGLA DE NEGOCIO
	if u.Userid < 0 {
		return dto.User{}, ErrInvalidUserID
	}
	return u, nil
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

// validateUser aplica las reglas de negocio comunes a Create y Update.
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

func (s *UserService) CreateMultipleUsers(ctx context.Context, users []dto.User) ([]dto.User, error) {
	if len(users) == 0 {
		return nil, ErrInvalidUser
	}
	// Valida todos
	for _, user := range users {
		if err := validateUser(user); err != nil {
			return nil, err
		}
	}

	created := make([]dto.User, 0, len(users))
	for _, user := range users {
		c, err := s.repo.Create(ctx, user)
		if err != nil {
			return nil, err
		}
		created = append(created, c)
	}
	return created, nil
}

func (s *UserService) Delete(ctx context.Context, id int32) error {
	return s.repo.Delete(ctx, id)
}

func (s *UserService) Update(ctx context.Context, user dto.User) (dto.User, error) {
	if err := validateUser(user); err != nil {
		return dto.User{}, err
	}
	return s.repo.Update(ctx, user)
}
