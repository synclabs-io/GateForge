package redis_repository

import (
	"context"
	"fmt"
	"time"

	"github.com/synclabs-io/GateForge/internal/cache"
	"github.com/synclabs-io/GateForge/internal/domain/entities"
	"github.com/synclabs-io/GateForge/internal/repositories/redis/dto"
)

type SessionRepository struct {
	client *cache.RedisClient
}

func NewSessionRepository(client *cache.RedisClient) *SessionRepository {
	return &SessionRepository{client: client}
}

func (r *SessionRepository) Save(ctx context.Context, input dto.SaveRequestDTO) (dto.SaveResponseDTO, error) {
	session, err := entities.NewSession(input.UserID, input.TTL)
	if err != nil {
		return dto.SaveResponseDTO{}, err
	}
	key := fmt.Sprintf("session:%s", session.Refresh.Hashed())

	pipe := r.client.Client.Pipeline()
	pipe.HSet(ctx, key, map[string]any{
		"id":         session.ID.String(),
		"user_id":    session.UserID.String(),
		"created_at": session.CreatedAt.Format(time.RFC3339),
	})
	pipe.Expire(ctx, key, input.TTL)
	if _, err := pipe.Exec(ctx); err != nil {
		return dto.SaveResponseDTO{}, err
	}

	return dto.SaveResponseDTO{
		Session: session,
	}, nil
}
