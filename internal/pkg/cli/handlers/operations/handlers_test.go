//go:build unit

package operations

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConnectHandler(t *testing.T) {
	handler := NewConnectHandler()

	assert.NotNil(t, handler)
}

func TestConnectHandler_ValidateArgs(t *testing.T) {
	handler := NewConnectHandler()

	t.Run("valid args", func(t *testing.T) {
		args := []string{"postgres"}
		err := handler.ValidateArgs(args)
		assert.NoError(t, err)
	})

	t.Run("no args provided", func(t *testing.T) {
		args := []string{}
		err := handler.ValidateArgs(args)
		assert.Error(t, err)
	})

	t.Run("multiple args allowed", func(t *testing.T) {
		args := []string{"postgres", "redis"}
		err := handler.ValidateArgs(args)
		assert.NoError(t, err) // Connect handler allows multiple args
	})
}

func TestConnectHandler_GetRequiredFlags(t *testing.T) {
	handler := NewConnectHandler()

	flags := handler.GetRequiredFlags()

	assert.NotNil(t, flags)
	assert.IsType(t, []string{}, flags)
}

func TestNewExecHandler(t *testing.T) {
	handler := NewExecHandler()

	assert.NotNil(t, handler)
}

func TestExecHandler_ValidateArgs(t *testing.T) {
	handler := NewExecHandler()

	t.Run("valid args", func(t *testing.T) {
		args := []string{"postgres", "psql", "-c", "SELECT 1"}
		err := handler.ValidateArgs(args)
		assert.NoError(t, err)
	})

	t.Run("insufficient args", func(t *testing.T) {
		args := []string{"postgres"}
		err := handler.ValidateArgs(args)
		assert.Error(t, err)
	})
}

func TestExecHandler_GetRequiredFlags(t *testing.T) {
	handler := NewExecHandler()

	flags := handler.GetRequiredFlags()

	assert.NotNil(t, flags)
	assert.IsType(t, []string{}, flags)
}

func TestLogsHandler_GetRequiredFlags(t *testing.T) {
	handler := NewLogsHandler()

	flags := handler.GetRequiredFlags()

	assert.NotNil(t, flags)
	assert.IsType(t, []string{}, flags)
}
