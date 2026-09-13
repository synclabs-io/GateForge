package value_objects_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/synclabs-io/GateForge/internal/domain/value_objects"
)

type MockHasher struct{}

func (h MockHasher) Hash(password string) (string, error) {
	return password + "_hashed_secret", nil
}

func (h MockHasher) Verify(plainPassword string, hashPassword string) (bool, error) {
	return hashPassword == plainPassword+"_hashed_secret", nil
}

func TestNewPassword(t *testing.T) {
	hasher := MockHasher{}

	t.Run("Normal password", func(t *testing.T) {
		_, err := value_objects.NewPassword("test_password_123", hasher)
		assert.NoError(t, err)
	})
	t.Run("Short password", func(t *testing.T) {
		_, err := value_objects.NewPassword("pass", hasher)
		assert.Error(t, err)
	})
	t.Run("Long password", func(t *testing.T) {
		_, err := value_objects.NewPassword(strings.Repeat("a", 256), hasher)
		assert.Error(t, err)
	})
}

func TestNewPassword_Regex(t *testing.T) {
	hasher := MockHasher{}

	t.Run("Valid characters in password", func(t *testing.T) {
		_, err := value_objects.NewPassword("Admin@123!_-", hasher)
		assert.NoError(t, err)
	})
	t.Run("Invalid characters (emoji)", func(t *testing.T) {
		_, err := value_objects.NewPassword("Secure🚀123", hasher)
		assert.Error(t, err)
	})
	t.Run("Invalid characters (cyrillic)", func(t *testing.T) {
		_, err := value_objects.NewPassword("МойПароль123!", hasher)
		assert.Error(t, err)
	})
}
