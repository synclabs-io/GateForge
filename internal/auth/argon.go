package auth

import (
	"fmt"

	"github.com/alexedwards/argon2id"
)

type PasswordHasher struct {
	params *argon2id.Params
}

func NewPasswordHasher(
	memory uint32,
	iterations uint32,
	parallelism uint8,
	saltLength uint32,
	keyLength uint32,
) *PasswordHasher {
	return &PasswordHasher{params: &argon2id.Params{
		Memory:      memory,
		Iterations:  iterations,
		Parallelism: parallelism,
		SaltLength:  saltLength,
		KeyLength:   keyLength,
	}}
}

func (h *PasswordHasher) Hash(plainPassword string) (string, error) {
	hash, err := argon2id.CreateHash(plainPassword, h.params)
	if err != nil {
		return "", fmt.Errorf("create hash: %w", err)
	}
	return hash, nil
}

func (h *PasswordHasher) Verify(plainPassword string, hashPassword string) (bool, error) {
	verify, err := argon2id.ComparePasswordAndHash(plainPassword, hashPassword)
	if err != nil {
		return false, fmt.Errorf("verify hash: %w", err)
	}
	return verify, nil
}
