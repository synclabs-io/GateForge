package entities_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/synclabs-io/GateForge/internal/domain/entities"
	"github.com/synclabs-io/GateForge/internal/domain/value_objects"
)

func TestNewSession(t *testing.T) {
	mockUserID := value_objects.NewID()
	ttl := 1 * time.Hour
	now := time.Now().UTC()

	t.Run("Create session successfully", func(t *testing.T) {
		session, err := entities.NewSession(mockUserID, ttl)

		assert.NoError(t, err)
		assert.NotNil(t, session)
		assert.Equal(t, mockUserID, session.UserID)
		assert.NotEmpty(t, session.Refresh.Raw())

		assert.WithinDuration(t, now, session.CreatedAt, 1*time.Second)
		assert.WithinDuration(t, now.Add(ttl), session.ExpiresAt, 1*time.Second)
	})
}
