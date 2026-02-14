//go:build unit

package utility

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWebInterfacesHandler_ShouldIncludeService(t *testing.T) {
	handler := NewWebInterfacesHandler()

	t.Run("includes all when showAll is true", func(t *testing.T) {
		runningServices := map[string]bool{"postgres": true}
		result := handler.shouldIncludeService("redis", runningServices, true)
		assert.True(t, result)
	})

	t.Run("includes all when runningServices is nil", func(t *testing.T) {
		result := handler.shouldIncludeService("redis", nil, false)
		assert.True(t, result)
	})

	t.Run("includes only running services when showAll is false", func(t *testing.T) {
		runningServices := map[string]bool{"postgres": true, "redis": false}
		assert.True(t, handler.shouldIncludeService("postgres", runningServices, false))
		assert.False(t, handler.shouldIncludeService("redis", runningServices, false))
	})
}

func TestWebInterfacesHandler_FormatStatus(t *testing.T) {
	handler := NewWebInterfacesHandler()

	t.Run("formats available status", func(t *testing.T) {
		result := handler.formatStatus(true)
		assert.Contains(t, result, "Available")
	})

	t.Run("formats unavailable status", func(t *testing.T) {
		result := handler.formatStatus(false)
		assert.Contains(t, result, "Not Available")
	})
}
