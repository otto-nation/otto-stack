//go:build unit

package operations

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogsHandler_GetRequiredFlags(t *testing.T) {
	handler := NewLogsHandler()
	flags := handler.GetRequiredFlags()
	assert.Empty(t, flags)
}

func TestLogsHandler_ValidateArgs_WithArgs(t *testing.T) {
	handler := NewLogsHandler()
	assert.NoError(t, handler.ValidateArgs(nil))
	assert.NoError(t, handler.ValidateArgs([]string{"redis", "postgres"}))
}
