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
	checker := NewUpdateChecker(DevelVersion)
	newer, err := checker.isNewer("2.0.0")
	require.NoError(t, err)
	assert.False(t, newer)

	checker = NewUpdateChecker("1.0.0")
	newer, err = checker.isNewer("2.0.0")
	require.NoError(t, err)
	assert.False(t, newer)

	checker = NewUpdateChecker("2.0.0")
	newer, err = checker.isNewer("1.0.0")
	require.NoError(t, err)
	assert.False(t, newer)
}

func TestDetectProjectVersion(t *testing.T) {
	constraint, err := DetectProjectVersion("/tmp/test")
	require.NoError(t, err)
	assert.NotNil(t, constraint)
	assert.Equal(t, "*", constraint.Operator)
}

func TestValidateProjectVersion(t *testing.T) {
	err := ValidateProjectVersion("/tmp/test")
	assert.NoError(t, err)

	err = ValidateProjectVersion("/tmp/test")
	assert.NoError(t, err)
}
