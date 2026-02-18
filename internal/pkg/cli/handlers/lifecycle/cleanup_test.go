package lifecycle

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCleanupHandler_ValidateArgs(t *testing.T) {
	handler := &CleanupHandler{}
	assert.NoError(t, handler.ValidateArgs([]string{}))
}

func TestCleanupHandler_GetRequiredFlags(t *testing.T) {
	handler := &CleanupHandler{}
	flags := handler.GetRequiredFlags()
	assert.Empty(t, flags)
}
