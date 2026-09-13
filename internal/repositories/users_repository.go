package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/synclabs-io/GateForge/db/generated"
	domain_errors "github.com/synclabs-io/GateForge/internal/domain/errors"
)

type UsersRepository struct {
	q *db.Queries
}

func NewUsersRepository(pool *pgxpool.Pool) *UsersRepository {
	return &UsersRepository{q: db.New(pool)}
}

func (r *UsersRepository) CreateUser(ctx context.Context, params db.CreateUserParams) (db.User, error) {
	user, err := r.q.CreateUser(ctx, params)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return db.User{}, domain_errors.ErrUserAlreadyExists
		}
		return db.User{}, fmt.Errorf("UsersRepository.CreateUser: %w", err)
	}
	return user, nil
}
