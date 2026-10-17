/*

Este repository.go permite brindar una capa intermedia entre el codigo generado por sqlc y el resto de la logica del negocio,
permitiendo cambiar en un futuro la implementacion de la base de datos sin afectar el resto del codigo.

*/

package User

// Imports
import (
	db "WebGrupal/db/sqlc"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// Interfaz
type UserRepository interface {
	GetUser(ctx context.Context, id int32) (db.User, error)
	GetUsers(ctx context.Context) ([]db.User, error)
	Create(ctx context.Context, user db.User) (db.User, error)
	Delete(ctx context.Context, id int32) error
	Update(ctx context.Context, user db.User) (db.User, error)
}

type SQLCRepository struct {
	queries *db.Queries
}

// "Constructor"
func NewSQLCRepository(q *db.Queries) *SQLCRepository {
	return &SQLCRepository{queries: q}
}

// A partir de aca, implementaciones:

func (r *SQLCRepository) GetUser(ctx context.Context, id int32) (db.User, error) {
	u, err := r.queries.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.User{}, ErrUserNotFound
		}
		return db.User{}, fmt.Errorf("consultar user %d: %w", id, err)
	}
	return u, nil
}

func (r *SQLCRepository) GetUsers(ctx context.Context) ([]db.User, error) {
	rows, err := r.queries.GetUsers(ctx)
	if err != nil {
		return nil, err
	}
	users := make([]db.User, 0, len(rows))
	for _, u := range rows {
		users = append(users, u)
	}
	return users, nil
}

func (r *SQLCRepository) Create(ctx context.Context, user db.User) (db.User, error) {
	u, err := r.queries.CreateUser(ctx, db.CreateUserParams{
		Username:           user.Username,
		Userpasswordhashed: user.Userpasswordhashed,
		Userbirthday:       user.Userbirthday,
		Usergenre:          user.Usergenre,
	})
	if err != nil {
		return db.User{}, err
	}
	return u, nil
}

func (r *SQLCRepository) Delete(ctx context.Context, id int32) error {
	return r.queries.DeleteUser(ctx, id)
}

func (r *SQLCRepository) Update(ctx context.Context, user db.User) (db.User, error) {
	err := r.queries.UpdateUser(ctx, db.UpdateUserParams{
		Userid:             user.Userid,
		Username:           user.Username,
		Userpasswordhashed: user.Userpasswordhashed,
		Userbirthday:       user.Userbirthday,
		Usergenre:          user.Usergenre,
	})
	if err != nil {
		return db.User{}, err
	}
	return r.GetUser(ctx, user.Userid)
}

