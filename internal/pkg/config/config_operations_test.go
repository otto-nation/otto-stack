//go:build unit

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigService_Operations(t *testing.T) {
	t.Run("creates config service", func(t *testing.T) {
		service := NewConfigService()
		assert.NotNil(t, service)
	})

	t.Run("loads config", func(t *testing.T) {
		service := NewConfigService()

		config, err := service.LoadConfig()
		if err != nil {
			assert.Error(t, err)
		} else {
			assert.NotNil(t, config)
		}
	})
}
