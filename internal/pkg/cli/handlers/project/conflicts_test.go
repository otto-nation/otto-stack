//go:build unit

package project

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConflictsHandler_ValidateArgs(t *testing.T) {
	handler := &ConflictsHandler{}

	t.Run("accepts no arguments", func(t *testing.T) {
		err := handler.ValidateArgs([]string{})
		assert.NoError(t, err)
	})

	t.Run("accepts service names", func(t *testing.T) {
		err := handler.ValidateArgs([]string{"postgres", "redis"})
		assert.NoError(t, err)
	})
}

func TestConflictsHandler_GetRequiredFlags(t *testing.T) {
	handler := &ConflictsHandler{}
	flags := handler.GetRequiredFlags()
	assert.Empty(t, flags)
}

func TestDepsHandler_ValidateArgs(t *testing.T) {
	handler := &DepsHandler{}

	t.Run("accepts no arguments", func(t *testing.T) {
		err := handler.ValidateArgs([]string{})
		assert.NoError(t, err)
	})

	t.Run("accepts service names", func(t *testing.T) {
		err := handler.ValidateArgs([]string{"postgres"})
		assert.NoError(t, err)
	})
}

func TestDepsHandler_GetRequiredFlags(t *testing.T) {
	handler := &DepsHandler{}
	flags := handler.GetRequiredFlags()
	assert.Empty(t, flags)
}
