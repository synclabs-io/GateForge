package dto

import (
	"time"

	"github.com/synclabs-io/GateForge/internal/domain/entities"
	"github.com/synclabs-io/GateForge/internal/domain/value_objects"
)

type SaveRequestDTO struct {
	UserID *value_objects.ID
	TTL    time.Duration
}

type SaveResponseDTO struct {
	Session *entities.Session
}
