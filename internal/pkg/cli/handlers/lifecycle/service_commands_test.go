//go:build unit

package lifecycle

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/common"
	"github.com/stretchr/testify/assert"
)

func TestServiceCommand(t *testing.T) {
	t.Run("creates new service command", func(t *testing.T) {
		stateManager := &common.StateManager{}
		cmd := NewServiceCommand("up", stateManager)
		assert.NotNil(t, cmd)
	})
}

func TestCleanupOperations(t *testing.T) {
	t.Run("validates cleanup handler creation", func(t *testing.T) {
		handler := NewCleanupHandler()
		assert.NotNil(t, handler)
	})
}

func TestRestartOperations(t *testing.T) {
	t.Run("validates restart handler creation", func(t *testing.T) {
		handler := NewRestartHandler()
		assert.NotNil(t, handler)
	})
}

func TestUpHandler(t *testing.T) {
	t.Run("validates up handler creation", func(t *testing.T) {
		handler := NewUpHandler()
		assert.NotNil(t, handler)
	})
}

func TestDownHandler(t *testing.T) {
	t.Run("validates down handler creation", func(t *testing.T) {
		handler := NewDownHandler()
		assert.NotNil(t, handler)
	})
}
