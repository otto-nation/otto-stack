//go:build unit

package project

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseServices(t *testing.T) {
	t.Run("parses single service", func(t *testing.T) {
		result := parseServices("postgres")
		assert.Equal(t, []string{"postgres"}, result)
	})

	t.Run("parses multiple services", func(t *testing.T) {
		result := parseServices("postgres,redis,mysql")
		assert.Equal(t, []string{"postgres", "redis", "mysql"}, result)
	})

	t.Run("trims whitespace", func(t *testing.T) {
		result := parseServices("postgres , redis , mysql")
		assert.Equal(t, []string{"postgres", "redis", "mysql"}, result)
	})

	t.Run("handles extra spaces", func(t *testing.T) {
		result := parseServices("  postgres  ,  redis  ")
		assert.Equal(t, []string{"postgres", "redis"}, result)
	})

	t.Run("handles empty string", func(t *testing.T) {
		result := parseServices("")
		assert.Equal(t, []string{""}, result)
	})
}

func TestGetDefaultValidation(t *testing.T) {
	t.Run("returns validation map with all keys enabled", func(t *testing.T) {
		result := getDefaultValidation()
		assert.NotNil(t, result)

		// All validation keys should be true
		for key, value := range result {
			assert.True(t, value, "validation key %s should be true", key)
		}
	})

	t.Run("includes all validation registry keys", func(t *testing.T) {
		result := getDefaultValidation()

		// Should have same number of keys as ValidationRegistry
		assert.Equal(t, len(ValidationRegistry), len(result))
	})
}

func TestInitHandler_ValidateProjectName(t *testing.T) {
	handler := NewInitHandler()

	t.Run("validates valid project name", func(t *testing.T) {
		err := handler.ValidateProjectName("my-project")
		assert.NoError(t, err)
	})

	t.Run("validates project name with numbers", func(t *testing.T) {
		err := handler.ValidateProjectName("project123")
		assert.NoError(t, err)
	})

	t.Run("validates project name with underscores", func(t *testing.T) {
		err := handler.ValidateProjectName("my_project")
		assert.NoError(t, err)
	})
}
