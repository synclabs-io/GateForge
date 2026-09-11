package domain_errors

import (
	"errors"
	"fmt"
)

const (
	MinUsernameLength = 3
	MaxUsernameLength = 25
)

var (
	ErrUsernameEmpty    = errors.New("username cannot be empty")
	ErrUsernameTooLong  = fmt.Errorf("the username must not exceed %d characters", MaxUsernameLength)
	ErrUsernameTooShort = fmt.Errorf("the username must not be shorter than %d characters", MinUsernameLength)
	ErrUsernameInvalid  = errors.New("username can only contain letters, numbers, underscores, and hyphens")
)
