package integration_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/synclabs-io/GateForge/internal/cache"
	"github.com/synclabs-io/GateForge/internal/config"
	"github.com/synclabs-io/GateForge/internal/domain/value_objects"
	redis_repository "github.com/synclabs-io/GateForge/internal/repositories/redis"
	"github.com/synclabs-io/GateForge/internal/repositories/redis/dto"
)

func TestSessionRepository_Integration(t *testing.T) {
	ctx := context.Background()
	cfg, err := config.Load()
	require.NoError(t, err, "failed to load config")

	redisClient, err := cache.NewRedisClient(ctx, cfg)
	require.NoError(t, err, "failed to initialize redis client")

	repo := redis_repository.NewSessionRepository(redisClient)

	mockUserID := value_objects.NewID()
	input := dto.SaveRequestDTO{
		UserID: mockUserID,
		TTL:    10 * time.Second,
	}

	response, err := repo.Save(ctx, input)

	assert.NoError(t, err)
	assert.NotEmpty(t, response.Session.Refresh, "repository must return raw refresh token string")

	expectedRedisKey := fmt.Sprintf("session:%s", response.Session.Refresh.Hashed())
	exists, err := redisClient.Client.Exists(ctx, expectedRedisKey).Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), exists, "session key should exist in redis database")

	storedHash, err := redisClient.Client.HGetAll(ctx, expectedRedisKey).Result()
	assert.NoError(t, err)
	assert.Equal(t, mockUserID.String(), storedHash["user_id"], "stored json hash must contain user id string")
}
