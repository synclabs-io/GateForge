package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/synclabs-io/GateForge/internal/domain/value_objects"
)

var ErrFailedToCreateSession = errors.New("failed to create session due to invalid components")

type Session struct {
	ID        *value_objects.ID           `json:"id"`
	UserID    *value_objects.ID           `json:"user_id"`
	Refresh   *value_objects.RefreshToken `json:"_"`
	ExpiresAt time.Time                   `json:"expires_at"`
	CreatedAt time.Time                   `json:"created_at"`
}

func NewSession(userID *value_objects.ID, ttl time.Duration) (*Session, error) {
	refresh, err := value_objects.NewRefreshToken(uuid.NewString())
	if err != nil {
		return nil, ErrFailedToCreateSession
	}

	return &Session{
		ID:        value_objects.NewID(),
		UserID:    userID,
		Refresh:   &refresh,
		ExpiresAt: time.Now().UTC().Add(ttl),
		CreatedAt: time.Now().UTC(),
	}, nil
}
