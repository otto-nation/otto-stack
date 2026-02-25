package lifecycle

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRestartHandler_ValidateArgs(t *testing.T) {
	handler := &RestartHandler{}
	assert.NoError(t, handler.ValidateArgs([]string{}))
}

func TestRestartHandler_GetRequiredFlags(t *testing.T) {
	handler := &RestartHandler{}
	flags := handler.GetRequiredFlags()
	assert.Empty(t, flags)
}
