package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/synclabs-io/GateForge/db/generated"
	"github.com/synclabs-io/GateForge/internal/domain/entities"
	domain_errors "github.com/synclabs-io/GateForge/internal/domain/errors"
)

type UsersRepository struct {
	q *db.Queries
}

func NewUsersRepository(pool *pgxpool.Pool) *UsersRepository {
	return &UsersRepository{q: db.New(pool)}
}

func (r *UsersRepository) CreateUser(ctx context.Context, entity entities.User) (entities.User, error) {
	_, err := r.q.CreateUser(ctx, db.CreateUserParams{
		ID:           entity.ID.String(),
		Username:     entity.Username.String(),
		PasswordHash: entity.Password.String(),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return entities.User{}, domain_errors.ErrUserAlreadyExists
		}
		return entities.User{}, fmt.Errorf("UsersRepository.CreateUser: %w", err)
	}
	return entity, nil
}
