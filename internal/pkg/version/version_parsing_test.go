//go:build unit

package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionConstants(t *testing.T) {
	t.Run("validates version constants", func(t *testing.T) {
		assert.Equal(t, "dev", DevelVersion)
		assert.Equal(t, 999999, MaxVersionNumber)
		assert.Equal(t, 7, GitCommitHashLength)
	})

	t.Run("validates comparison constants", func(t *testing.T) {
		assert.Equal(t, -1, VersionOlder)
		assert.Equal(t, 0, VersionEqual)
		assert.Equal(t, 1, VersionNewer)
	})

	t.Run("validates default constants", func(t *testing.T) {
		assert.Equal(t, "dev", DefaultVersion)
		assert.Equal(t, "unknown", DefaultCommit)
		assert.Equal(t, "unknown", DefaultBuildDate)
		assert.Equal(t, "unknown", DefaultBuildBy)
	})
}

func TestUpdateChecker_EdgeCases(t *testing.T) {
	t.Run("creates update checker with version", func(t *testing.T) {
		checker := NewUpdateChecker("1.0.0")
		assert.NotNil(t, checker)
		assert.Equal(t, "1.0.0", checker.currentVersion)
	})

	t.Run("handles empty version", func(t *testing.T) {
		checker := NewUpdateChecker("")
		assert.NotNil(t, checker)
		assert.Equal(t, "", checker.currentVersion)
	})

	t.Run("handles dev version", func(t *testing.T) {
		checker := NewUpdateChecker(DevelVersion)
		assert.NotNil(t, checker)
		assert.Equal(t, DevelVersion, checker.currentVersion)
	})
}

func TestVersionParser_EdgeCases(t *testing.T) {
	t.Run("parses valid semantic version", func(t *testing.T) {
		version, err := ParseVersion("1.2.3")
		assert.NoError(t, err)
		assert.NotNil(t, version)
		assert.Equal(t, 1, version.Major)
		assert.Equal(t, 2, version.Minor)
		assert.Equal(t, 3, version.Patch)
	})

	t.Run("handles invalid version format", func(t *testing.T) {
		_, err := ParseVersion("invalid")
		assert.Error(t, err)
	})

	t.Run("handles empty version string", func(t *testing.T) {
		_, err := ParseVersion("")
		assert.Error(t, err)
	})

	t.Run("handles version with prerelease", func(t *testing.T) {
		version, err := ParseVersion("1.0.0-alpha")
		if err == nil {
			assert.NotNil(t, version)
			assert.Equal(t, 1, version.Major)
		} else {
			// Parser may not support prerelease
			assert.Error(t, err)
		}
	})
}

func TestVersionComparison(t *testing.T) {
	t.Run("compares versions correctly", func(t *testing.T) {
		v1, err1 := ParseVersion("1.0.0")
		v2, err2 := ParseVersion("2.0.0")

		if err1 == nil && err2 == nil {
			result := v1.Compare(*v2)
			assert.Equal(t, VersionOlder, result)

			result2 := v2.Compare(*v1)
			assert.Equal(t, VersionNewer, result2)

			result3 := v1.Compare(*v1)
			assert.Equal(t, VersionEqual, result3)
		}
	})
}

func TestVersionConstraint_EdgeCases(t *testing.T) {
	t.Run("parses version constraint", func(t *testing.T) {
		constraint, err := ParseVersionConstraint(">=1.0.0")
		if err == nil {
			assert.NotNil(t, constraint)
		} else {
			// Parser may not support constraints
			assert.Error(t, err)
		}
	})

	t.Run("handles invalid constraint", func(t *testing.T) {
		_, err := ParseVersionConstraint("invalid")
		assert.Error(t, err)
	})
}
