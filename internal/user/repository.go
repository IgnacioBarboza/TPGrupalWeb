/*

Este repository.go permite brindar una capa intermedia entre el codigo generado por sqlc y el resto de la logica del negocio,
permitiendo cambiar en un futuro la implementacion de la base de datos sin afectar el resto del codigo.

*/

package User

// Imports
import (
	db "WebGrupal/db/sqlc"
	dto "WebGrupal/internal"
	"context"
)

// Interfaz
type UserRepository interface {
	GetUsers(ctx context.Context) ([]dto.User, error)
	Create(ctx context.Context, user dto.User) (dto.User, error)
}

type SQLCRepository struct {
	queries *db.Queries
}

// "Constructor"
func NewSQLCRepository(q *db.Queries) *SQLCRepository {
	return &SQLCRepository{queries: q}
}

// A partir de aca, implementaciones:

func (r *SQLCRepository) GetUsers(ctx context.Context) ([]dto.User, error) {
	rows, err := r.queries.GetUsers(ctx)
	if err != nil {
		return nil, err
	}
	users := make([]dto.User, 0, len(rows))
	for _, u := range rows {
		users = append(users, dto.User(u))
	}
	return users, nil
}

func (r *SQLCRepository) Create(ctx context.Context, user dto.User) (dto.User, error) {
	u, err := r.queries.CreateUser(ctx, db.CreateUserParams{
		Username:           user.Username,
		Userpasswordhashed: user.Userpasswordhashed,
		Userbirthday:       user.Userbirthday,
		Usergenre:          user.Usergenre,
	})
	if err != nil {
		return dto.User{}, err
	}
	return dto.User(u), nil
}
