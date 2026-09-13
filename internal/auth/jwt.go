package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/synclabs-io/GateForge/internal/domain/value_objects"
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	secret        []byte
	tokenDuration time.Duration
}

func NewJWTManager(secret string, tokenDuration time.Duration) *JWTManager {
	return &JWTManager{
		secret:        []byte(secret),
		tokenDuration: tokenDuration,
	}
}

func (m *JWTManager) GenerateToken(userID *value_objects.ID) (string, error) {
	claims := Claims{
		UserID: userID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}
