package value_objects

import (
	"regexp"

	domain_errors "github.com/synclabs-io/GateForge/internal/domain/errors"
)

var passwordRegex = regexp.MustCompile(`^[a-zA-Z0-9!@#$%^&*()_+\-=\[\]{};':"\\|,.<>/?]+$`)

type PasswordHash interface {
	Hash(plainPassword string) (string, error)
	Verify(plainPassword string, hashPassword string) (bool, error)
}

type Password struct {
	value  string
	hasher PasswordHash
}

func NewPassword(value string, hasher PasswordHash) (*Password, error) {
	if err := validatePassword(value); err != nil {
		return nil, err
	}

	passwordHash, err := hasher.Hash(value)
	if err != nil {
		return nil, err
	}

	return &Password{value: passwordHash, hasher: hasher}, nil
}

func NewPasswordFromHash(hashedPassword string, hasher PasswordHash) *Password {
	return &Password{value: hashedPassword, hasher: hasher}
}

func validatePassword(value string) error {
	if value == "" {
		return domain_errors.ErrPasswordEmpty
	}
	if len(value) < domain_errors.MinPasswordLength {
		return domain_errors.ErrPasswordTooShort
	}
	if len(value) > domain_errors.MaxPasswordLength {
		return domain_errors.ErrPasswordTooLong
	}
	if !passwordRegex.MatchString(value) {
		return domain_errors.ErrPasswordInvalid
	}

	return nil
}

func (p *Password) Validate(plainPassword string) (bool, error) {
	validate, err := p.hasher.Verify(plainPassword, p.value)
	if err != nil {
		return false, err
	}

	return validate, nil
}

func (p *Password) String() string {
	return p.value
}
