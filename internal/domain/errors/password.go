package domain_errors

import (
	"errors"
	"fmt"
)

const (
	MinPasswordLength = 8
	MaxPasswordLength = 255
)

var (
	ErrPasswordEmpty    = errors.New("password cannot be empty")
	ErrPasswordTooLong  = fmt.Errorf("the password must not exceed %d characters", MaxUsernameLength)
	ErrPasswordTooShort = fmt.Errorf("the password must not be shorter than %d characters", MinUsernameLength)
	ErrPasswordInvalid  = errors.New("password contains invalid characters")
)
