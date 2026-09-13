package entities

import (
	"time"

	"github.com/synclabs-io/GateForge/internal/domain/value_objects"
)

type User struct {
	ID        *value_objects.ID       `json:"id"`
	Username  *value_objects.Username `json:"username"`
	Password  *value_objects.Password `json:"password"`
	CreatedAt time.Time               `json:"created_at"`
	UpdatedAt time.Time               `json:"updated_at"`
}

func NewUser(plainUsername string, plainPassword string, hasher value_objects.PasswordHash) (*User, error) {
	username, err := value_objects.NewUsername(plainUsername)
	if err != nil {
		return nil, err
	}

	password, err := value_objects.NewPassword(plainPassword, hasher)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:        value_objects.NewID(),
		Username:  username,
		Password:  password,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}, nil
}
