//go:build unit

package utility

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewVersionHandler(t *testing.T) {
	handler := NewVersionHandler()
	assert.NotNil(t, handler)
}

func TestVersionHandler_ValidateArgs(t *testing.T) {
	handler := NewVersionHandler()
	err := handler.ValidateArgs([]string{})
	assert.NoError(t, err)
}

func TestVersionHandler_GetRequiredFlags(t *testing.T) {
	handler := NewVersionHandler()
	flags := handler.GetRequiredFlags()
	assert.Empty(t, flags)
}

func TestNewEnforcementHandler(t *testing.T) {
	handler := NewEnforcementHandler(nil)
	assert.NotNil(t, handler)
}
