package value_objects

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"github.com/google/uuid"
)

var ErrInvalidRefreshToken = errors.New("refresh token must be a valid uuid")

type RefreshToken struct {
	raw    string
	hashed string
}

func NewRefreshToken(rawStr string) (RefreshToken, error) {
	parsed, err := uuid.Parse(rawStr)
	if err != nil {
		return RefreshToken{}, ErrInvalidRefreshToken
	}

	hash := sha256.Sum256([]byte(parsed.String()))

	return RefreshToken{
		raw:    parsed.String(),
		hashed: hex.EncodeToString(hash[:]),
	}, nil
}

func GenerateRefreshToken() RefreshToken {
	rawUUID := uuid.NewString()
	hash := sha256.Sum256([]byte(rawUUID))

	return RefreshToken{
		raw:    rawUUID,
		hashed: hex.EncodeToString(hash[:]),
	}
}

func (r *RefreshToken) Raw() string {
	return r.raw
}

func (r *RefreshToken) Hashed() string {
	return r.hashed
}
