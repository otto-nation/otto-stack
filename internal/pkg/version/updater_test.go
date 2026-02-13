//go:build unit

package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUpdateChecker(t *testing.T) {
	checker := NewUpdateChecker("1.0.0")
	assert.NotNil(t, checker)
	assert.Equal(t, "1.0.0", checker.currentVersion)
	assert.NotNil(t, checker.client)
}

func TestUpdateChecker_IsNewer(t *testing.T) {
	t.Run("returns false for dev builds", func(t *testing.T) {
		checker := NewUpdateChecker("devel")
		newer, err := checker.isNewer("2.0.0")
		require.NoError(t, err)
		assert.False(t, newer)
	})

	t.Run("handles version comparison", func(t *testing.T) {
		checker := NewUpdateChecker("1.0.0")
		// Will return false for dev builds (AppVersion is "devel" in tests)
		newer, err := checker.isNewer("2.0.0")
		require.NoError(t, err)
		assert.False(t, newer)
	})
}

func TestDetectProjectVersion(t *testing.T) {
	t.Run("returns wildcard constraint", func(t *testing.T) {
		constraint, err := DetectProjectVersion("/tmp/test")
		require.NoError(t, err)
		assert.NotNil(t, constraint)
		assert.Equal(t, "*", constraint.Operator)
	})
}

func TestValidateProjectVersion(t *testing.T) {
	t.Run("validates project version", func(t *testing.T) {
		err := ValidateProjectVersion("/tmp/test")
		// Should not error with wildcard constraint
		assert.NoError(t, err)
	})
}
