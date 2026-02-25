//go:build unit

package services

import (
	"testing"

	servicetypes "github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidationService_ValidateUserServices(t *testing.T) {
	manager, err := New()
	require.NoError(t, err)
	validator := NewValidationService(manager)

	t.Run("empty list returns error", func(t *testing.T) {
		err := validator.ValidateUserServices([]string{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "at least one service must be selected")
	})

	t.Run("valid service passes", func(t *testing.T) {
		err := validator.ValidateUserServices([]string{ServicePostgres})
		assert.NoError(t, err)
	})

	t.Run("invalid service fails", func(t *testing.T) {
		err := validator.ValidateUserServices([]string{"nonexistent"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "nonexistent")
	})

	t.Run("hidden service fails for user", func(t *testing.T) {
		// localstack is hidden
		err := validator.ValidateUserServices([]string{"localstack"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "internal service")
	})
}

func TestValidationService_ValidateWithContext(t *testing.T) {
	manager, err := New()
	require.NoError(t, err)
	validator := NewValidationService(manager)

	t.Run("user context rejects hidden service", func(t *testing.T) {
		ctx := NewUserValidationContext()
		err := validator.ValidateWithContext("localstack", ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "internal service")
	})

	t.Run("dependency context allows hidden service", func(t *testing.T) {
		ctx := NewDependencyValidationContext()
		err := validator.ValidateWithContext("localstack", ctx)
		assert.NoError(t, err)
	})

	t.Run("internal context allows hidden service", func(t *testing.T) {
		ctx := NewInternalValidationContext()
		err := validator.ValidateWithContext("localstack", ctx)
		assert.NoError(t, err)
	})

	t.Run("all contexts reject nonexistent service", func(t *testing.T) {
		contexts := []ValidationContext{
			NewUserValidationContext(),
			NewDependencyValidationContext(),
			NewInternalValidationContext(),
		}

		for _, ctx := range contexts {
			err := validator.ValidateWithContext("nonexistent", ctx)
			assert.Error(t, err)
		}
	})
}

func TestValidationService_ValidateResolvedServices(t *testing.T) {
	manager, err := New()
	require.NoError(t, err)
	validator := NewValidationService(manager)

	t.Run("empty list returns error", func(t *testing.T) {
		err := validator.ValidateResolvedServices([]servicetypes.ServiceConfig{})
		assert.Error(t, err)
	})

	t.Run("valid services pass", func(t *testing.T) {
		postgres, _ := manager.GetService(ServicePostgres)
		redis, _ := manager.GetService(ServiceRedis)

		err := validator.ValidateResolvedServices([]servicetypes.ServiceConfig{*postgres, *redis})
		assert.NoError(t, err)
	})

	t.Run("allows hidden services in resolved list", func(t *testing.T) {
		localstack, _ := manager.GetService("localstack")
		postgres, _ := manager.GetService(ServicePostgres)

		err := validator.ValidateResolvedServices([]servicetypes.ServiceConfig{*localstack, *postgres})
		assert.NoError(t, err)
	})

	t.Run("detects duplicates", func(t *testing.T) {
		postgres, _ := manager.GetService(ServicePostgres)

		err := validator.ValidateResolvedServices([]servicetypes.ServiceConfig{*postgres, *postgres})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "duplicate")
	})
}
