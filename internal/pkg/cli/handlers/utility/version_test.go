//go:build unit

package utility

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionDisplayManager_Methods(t *testing.T) {
	t.Run("tests NewVersionDisplayManager", func(t *testing.T) {
		manager := NewVersionDisplayManager()
		assert.NotNil(t, manager)
	})

	t.Run("tests DisplayBasic", func(t *testing.T) {
		manager := NewVersionDisplayManager()

		manager.DisplayBasic("1.0.0", "text")
		// Should not panic
		assert.True(t, true)
	})

	t.Run("tests DisplayFull", func(t *testing.T) {
		manager := NewVersionDisplayManager()

		manager.DisplayFull("1.0.0", "text")
		// Should not panic
		assert.True(t, true)
	})

	t.Run("tests GetCurrentVersion", func(t *testing.T) {
		manager := NewVersionDisplayManager()

		version := manager.GetCurrentVersion()
		assert.IsType(t, "", version)
	})
}
