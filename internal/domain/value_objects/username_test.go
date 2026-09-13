package value_objects_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/synclabs-io/GateForge/internal/domain/value_objects"
)

func TestNewUsername(t *testing.T) {
	t.Run("Normal name", func(t *testing.T) {
		_, err := value_objects.NewUsername("alex")
		assert.NoError(t, err)
	})
	t.Run("Short name", func(t *testing.T) {
		_, err := value_objects.NewUsername("jo")
		assert.Error(t, err)
	})
	t.Run("Long name", func(t *testing.T) {
		_, err := value_objects.NewUsername("very_long_username_that_should_fail")
		assert.Error(t, err)
	})
}
