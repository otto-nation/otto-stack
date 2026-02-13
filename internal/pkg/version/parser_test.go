//go:build unit

package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseVersion(t *testing.T) {
	t.Run("parses valid semver", func(t *testing.T) {
		v, err := ParseVersion("1.2.3")
		require.NoError(t, err)
		assert.Equal(t, 1, v.Major)
		assert.Equal(t, 2, v.Minor)
		assert.Equal(t, 3, v.Patch)
	})

	t.Run("parses version with v prefix", func(t *testing.T) {
		v, err := ParseVersion("v1.2.3")
		require.NoError(t, err)
		assert.Equal(t, 1, v.Major)
	})

	t.Run("parses version with prerelease", func(t *testing.T) {
		v, err := ParseVersion("1.2.3-alpha")
		require.NoError(t, err)
		assert.Equal(t, "alpha", v.PreRelease)
	})

	t.Run("handles latest", func(t *testing.T) {
		v, err := ParseVersion("latest")
		require.NoError(t, err)
		assert.NotNil(t, v)
	})

	t.Run("returns error for empty string", func(t *testing.T) {
		_, err := ParseVersion("")
		assert.Error(t, err)
	})

	t.Run("returns error for invalid format", func(t *testing.T) {
		_, err := ParseVersion("invalid")
		assert.Error(t, err)
	})
}

func TestVersion_String(t *testing.T) {
	t.Run("formats basic version", func(t *testing.T) {
		v := Version{Major: 1, Minor: 2, Patch: 3}
		assert.Equal(t, "1.2.3", v.String())
	})

	t.Run("includes prerelease", func(t *testing.T) {
		v := Version{Major: 1, Minor: 2, Patch: 3, PreRelease: "alpha"}
		assert.Contains(t, v.String(), "alpha")
	})
}

func TestVersion_Compare(t *testing.T) {
	t.Run("compares equal versions", func(t *testing.T) {
		v1 := Version{Major: 1, Minor: 2, Patch: 3}
		v2 := Version{Major: 1, Minor: 2, Patch: 3}
		assert.Equal(t, VersionEqual, v1.Compare(v2))
	})

	t.Run("compares newer version", func(t *testing.T) {
		v1 := Version{Major: 2, Minor: 0, Patch: 0}
		v2 := Version{Major: 1, Minor: 0, Patch: 0}
		assert.Equal(t, VersionNewer, v1.Compare(v2))
	})

	t.Run("compares older version", func(t *testing.T) {
		v1 := Version{Major: 1, Minor: 0, Patch: 0}
		v2 := Version{Major: 2, Minor: 0, Patch: 0}
		assert.Equal(t, VersionOlder, v1.Compare(v2))
	})

	t.Run("release is newer than prerelease", func(t *testing.T) {
		v1 := Version{Major: 1, Minor: 0, Patch: 0}
		v2 := Version{Major: 1, Minor: 0, Patch: 0, PreRelease: "alpha"}
		assert.Equal(t, VersionNewer, v1.Compare(v2))
	})
}

func TestVersionConstraint_Satisfies(t *testing.T) {
	t.Run("equals operator", func(t *testing.T) {
		constraint := VersionConstraint{
			Operator: "=",
			Version:  Version{Major: 1, Minor: 2, Patch: 3},
		}
		v := Version{Major: 1, Minor: 2, Patch: 3}
		assert.True(t, constraint.Satisfies(v))
	})

	t.Run("greater than operator", func(t *testing.T) {
		constraint := VersionConstraint{
			Operator: ">",
			Version:  Version{Major: 1, Minor: 0, Patch: 0},
		}
		v := Version{Major: 2, Minor: 0, Patch: 0}
		assert.True(t, constraint.Satisfies(v))
	})

	t.Run("less than operator", func(t *testing.T) {
		constraint := VersionConstraint{
			Operator: "<",
			Version:  Version{Major: 2, Minor: 0, Patch: 0},
		}
		v := Version{Major: 1, Minor: 0, Patch: 0}
		assert.True(t, constraint.Satisfies(v))
	})

	t.Run("wildcard operator", func(t *testing.T) {
		constraint := VersionConstraint{
			Operator: "*",
			Version:  Version{Major: 1, Minor: 0, Patch: 0},
		}
		v := Version{Major: 99, Minor: 99, Patch: 99}
		assert.True(t, constraint.Satisfies(v))
	})
}
