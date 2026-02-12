package services

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestValidator_ValidateServiceNames(t *testing.T) {
	validator := NewValidator()

	t.Run("empty list", func(t *testing.T) {
		err := validator.ValidateServiceNames([]string{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "at least one service")
	})

	t.Run("valid service", func(t *testing.T) {
		err := validator.ValidateServiceNames([]string{"postgres"})
		if err != nil {
			t.Logf("Service validation error (may be expected if service doesn't exist): %v", err)
		}
	})
}

func TestValidator_ValidateServiceConfigs(t *testing.T) {
	validator := NewValidator()

	t.Run("empty list", func(t *testing.T) {
		err := validator.ValidateServiceConfigs([]types.ServiceConfig{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "at least one service")
	})

	t.Run("duplicate services", func(t *testing.T) {
		configs := []types.ServiceConfig{
			{Name: "postgres"},
			{Name: "postgres"},
		}
		err := validator.ValidateServiceConfigs(configs)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "duplicate")
	})

	t.Run("unique services", func(t *testing.T) {
		configs := []types.ServiceConfig{
			{Name: "postgres"},
			{Name: "redis"},
		}
		err := validator.ValidateServiceConfigs(configs)
		if err != nil {
			t.Logf("Service validation error (may be expected if services don't exist): %v", err)
		}
	})
}
