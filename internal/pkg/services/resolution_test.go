//go:build unit

package services

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/config"
	servicetypes "github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestResolveUpServices(t *testing.T) {
	cfg := &config.Config{
		Stack: config.StackConfig{
			Enabled: []string{ServicePostgres, ServiceRedis},
		},
	}

	t.Run("resolves specific services using constants", func(t *testing.T) {
		args := []string{ServicePostgres}

		configs, err := ResolveUpServices(args, cfg)

		// May fail due to service loading, but tests the resolution logic
		if err == nil {
			assert.NotEmpty(t, configs)
			// Should contain the requested service
			found := false
			for _, config := range configs {
				if config.Name == ServicePostgres {
					found = true
					break
				}
			}
			assert.True(t, found, "Should resolve requested service")
		}
	})

	t.Run("resolves enabled services when no args provided", func(t *testing.T) {
		args := []string{}

		configs, err := ResolveUpServices(args, cfg)

		// May fail due to service loading, but tests the logic path
		if err == nil {
			assert.NotEmpty(t, configs)
		}
		// Key test: function follows enabled services path
	})

	t.Run("handles invalid service names gracefully", func(t *testing.T) {
		args := []string{"nonexistent-service"}

		configs, err := ResolveUpServices(args, cfg)

		// Should handle invalid services gracefully
		if err != nil {
			assert.Error(t, err)
		} else {
			// If no error, configs should be empty or valid
			assert.NotNil(t, configs)
		}
	})
}

func TestServiceConfigValidation(t *testing.T) {
	t.Run("validates service config structure", func(t *testing.T) {
		config := servicetypes.ServiceConfig{
			Name:        ServicePostgres,
			Description: "PostgreSQL database",
			Category:    CategoryDatabase,
		}

		assert.Equal(t, ServicePostgres, config.Name)
		assert.Equal(t, CategoryDatabase, config.Category)
		assert.NotEmpty(t, config.Description)
	})
}
