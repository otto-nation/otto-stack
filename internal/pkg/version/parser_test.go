//go:build unit

package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseVersion(t *testing.T) {
	v, err := ParseVersion("1.2.3")
	require.NoError(t, err)
	assert.Equal(t, 1, v.Major)
	assert.Equal(t, 2, v.Minor)
	assert.Equal(t, 3, v.Patch)

	v, err = ParseVersion("v1.2.3")
	require.NoError(t, err)
	assert.Equal(t, 1, v.Major)

	v, err = ParseVersion("1.2.3-alpha")
	require.NoError(t, err)
	assert.Equal(t, "alpha", v.PreRelease)

	v, err = ParseVersion("latest")
	require.NoError(t, err)
	assert.NotNil(t, v)

	_, err = ParseVersion("")
	assert.Error(t, err)

	_, err = ParseVersion("invalid")
	assert.Error(t, err)
}

func TestVersion_String(t *testing.T) {
	v := Version{Major: 1, Minor: 2, Patch: 3}
	assert.Equal(t, "1.2.3", v.String())

	v = Version{Major: 1, Minor: 2, Patch: 3, PreRelease: "alpha"}
	assert.Contains(t, v.String(), "alpha")

	v = Version{Major: 1, Minor: 2, Patch: 3, Build: "20210101"}
	assert.Contains(t, v.String(), "20210101")

	v = Version{Major: 1, Minor: 2, Patch: 3, Original: "v1.2.3"}
	assert.Equal(t, "v1.2.3", v.String())
}

func TestVersion_Compare(t *testing.T) {
	v1 := Version{Major: 1, Minor: 2, Patch: 3}
	v2 := Version{Major: 1, Minor: 2, Patch: 3}
	assert.Equal(t, VersionEqual, v1.Compare(v2))

	v1 = Version{Major: 2, Minor: 0, Patch: 0}
	v2 = Version{Major: 1, Minor: 0, Patch: 0}
	assert.Equal(t, VersionNewer, v1.Compare(v2))

	v1 = Version{Major: 1, Minor: 0, Patch: 0}
	v2 = Version{Major: 2, Minor: 0, Patch: 0}
	assert.Equal(t, VersionOlder, v1.Compare(v2))

	v1 = Version{Major: 1, Minor: 0, Patch: 0}
	v2 = Version{Major: 1, Minor: 0, Patch: 0, PreRelease: "alpha"}
	assert.Equal(t, VersionNewer, v1.Compare(v2))

	v1 = Version{Major: 1, Minor: 2, Patch: 0}
	v2 = Version{Major: 1, Minor: 1, Patch: 0}
	assert.Equal(t, VersionNewer, v1.Compare(v2))

	v1 = Version{Major: 1, Minor: 0, Patch: 2}
	v2 = Version{Major: 1, Minor: 0, Patch: 1}
	assert.Equal(t, VersionNewer, v1.Compare(v2))

	v1 = Version{Major: 1, Minor: 0, Patch: 0, PreRelease: "beta"}
	v2 = Version{Major: 1, Minor: 0, Patch: 0, PreRelease: "alpha"}
	assert.Equal(t, VersionNewer, v1.Compare(v2))

	v1 = Version{Major: 1, Minor: 0, Patch: 0, PreRelease: "alpha"}
	v2 = Version{Major: 1, Minor: 0, Patch: 0}
	assert.Equal(t, VersionOlder, v1.Compare(v2))

	v1 = Version{Major: 1, Minor: 1, Patch: 0}
	v2 = Version{Major: 1, Minor: 2, Patch: 0}
	assert.Equal(t, VersionOlder, v1.Compare(v2))

	v1 = Version{Major: 1, Minor: 0, Patch: 1}
	v2 = Version{Major: 1, Minor: 0, Patch: 2}
	assert.Equal(t, VersionOlder, v1.Compare(v2))

	v1 = Version{Major: 1, Minor: 0, Patch: 0, PreRelease: "alpha"}
	v2 = Version{Major: 1, Minor: 0, Patch: 0, PreRelease: "beta"}
	assert.Equal(t, VersionOlder, v1.Compare(v2))
}

func TestVersionConstraint_Satisfies(t *testing.T) {
	constraint := VersionConstraint{Operator: "=", Version: Version{Major: 1, Minor: 2, Patch: 3}}
	v := Version{Major: 1, Minor: 2, Patch: 3}
	assert.True(t, constraint.Satisfies(v))

	constraint = VersionConstraint{Operator: ">", Version: Version{Major: 1, Minor: 0, Patch: 0}}
	v = Version{Major: 2, Minor: 0, Patch: 0}
	assert.True(t, constraint.Satisfies(v))

	constraint = VersionConstraint{Operator: "<", Version: Version{Major: 2, Minor: 0, Patch: 0}}
	v = Version{Major: 1, Minor: 0, Patch: 0}
	assert.True(t, constraint.Satisfies(v))

	constraint = VersionConstraint{Operator: "*", Version: Version{Major: 1, Minor: 0, Patch: 0}}
	v = Version{Major: 99, Minor: 99, Patch: 99}
	assert.True(t, constraint.Satisfies(v))

	constraint = VersionConstraint{Operator: ">=", Version: Version{Major: 1, Minor: 0, Patch: 0}}
	v1 := Version{Major: 1, Minor: 0, Patch: 0}
	v2 := Version{Major: 2, Minor: 0, Patch: 0}
	assert.True(t, constraint.Satisfies(v1))
	assert.True(t, constraint.Satisfies(v2))

	constraint = VersionConstraint{Operator: "<=", Version: Version{Major: 2, Minor: 0, Patch: 0}}
	v1 = Version{Major: 2, Minor: 0, Patch: 0}
	v2 = Version{Major: 1, Minor: 0, Patch: 0}
	assert.True(t, constraint.Satisfies(v1))
	assert.True(t, constraint.Satisfies(v2))

	constraint = VersionConstraint{Operator: "!=", Version: Version{Major: 1, Minor: 0, Patch: 0}}
	v = Version{Major: 2, Minor: 0, Patch: 0}
	assert.True(t, constraint.Satisfies(v))

	constraint = VersionConstraint{Operator: "==", Version: Version{Major: 1, Minor: 0, Patch: 0}}
	v = Version{Major: 1, Minor: 0, Patch: 0}
	assert.True(t, constraint.Satisfies(v))

	constraint = VersionConstraint{Operator: "", Version: Version{Major: 1, Minor: 0, Patch: 0}}
	v = Version{Major: 1, Minor: 0, Patch: 0}
	assert.True(t, constraint.Satisfies(v))

	constraint = VersionConstraint{Operator: "~>", Version: Version{Major: 1, Minor: 0, Patch: 0}}
	v = Version{Major: 1, Minor: 0, Patch: 0}
	assert.False(t, constraint.Satisfies(v))
}

func TestParseVersion_EdgeCases(t *testing.T) {
	v, err := ParseVersion("1.2.3+build123")
	require.NoError(t, err)
	assert.Equal(t, "build123", v.Build)

	v, err = ParseVersion("1.2.3-alpha+build123")
	require.NoError(t, err)
	assert.Equal(t, "alpha", v.PreRelease)
	assert.Equal(t, "build123", v.Build)

	v, err = ParseVersion("*")
	require.NoError(t, err)
	assert.NotNil(t, v)

	v, err = ParseVersion("  1.2.3  ")
	require.NoError(t, err)
	assert.Equal(t, 1, v.Major)

	v, err = ParseVersion("latest")
	require.NoError(t, err)
	assert.Equal(t, MaxVersionNumber, v.Major)
	assert.Equal(t, "latest", v.Original)

	v, err = ParseVersion("v1.2.3-beta+build")
	require.NoError(t, err)
	assert.Equal(t, "v1.2.3-beta+build", v.Original)
}

func TestGetAppVersion_EdgeCases(t *testing.T) {
	version := GetAppVersion()
	assert.NotEmpty(t, version)
	assert.IsType(t, "", version)
}

func TestGetShortVersion_EdgeCases(t *testing.T) {
	version := GetShortVersion()
	assert.NotEmpty(t, version)
	assert.IsType(t, "", version)
}
