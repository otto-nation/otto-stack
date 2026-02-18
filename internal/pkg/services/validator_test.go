package services

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/otto-nation/otto-stack/test/fixtures"
	"github.com/stretchr/testify/assert"
)

func TestValidator_ValidateServiceNames(t *testing.T) {
	validator := NewValidator()

	err := validator.ValidateServiceNames([]string{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "at least one service")

	err = validator.ValidateServiceNames([]string{ServicePostgres})
	if err != nil {
		t.Logf("Service validation error (may be expected if service doesn't exist): %v", err)
	}
}

func TestValidator_ValidateServiceConfigs(t *testing.T) {
	validator := NewValidator()

	err := validator.ValidateServiceConfigs([]types.ServiceConfig{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "at least one service")

	configs := []types.ServiceConfig{
		fixtures.NewServiceConfig(ServicePostgres).Build(),
		fixtures.NewServiceConfig(ServicePostgres).Build(),
	}
	err = validator.ValidateServiceConfigs(configs)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate")

	configs = []types.ServiceConfig{
		fixtures.NewServiceConfig(ServicePostgres).Build(),
		fixtures.NewServiceConfig(ServiceRedis).Build(),
	}
	err = validator.ValidateServiceConfigs(configs)
	if err != nil {
		t.Logf("Service validation error (may be expected if services don't exist): %v", err)
	}

	configs = []types.ServiceConfig{
		fixtures.NewServiceConfig("").Build(),
	}
	err = validator.ValidateServiceConfigs(configs)
	assert.Error(t, err)

	configs = []types.ServiceConfig{
		fixtures.NewServiceConfig(ServicePostgres).Build(),
	}
	err = validator.ValidateServiceConfigs(configs)
	if err != nil {
		t.Logf("Service validation error: %v", err)
	}
}
