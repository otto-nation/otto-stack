package version

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/stretchr/testify/assert"
)

func TestGetBuildInfo(t *testing.T) {
	info := GetBuildInfo()
	assert.NotNil(t, info)
	assert.NotEmpty(t, info.Version)
	assert.NotEmpty(t, info.GoVersion)
	assert.NotEmpty(t, info.Platform)
	assert.NotEmpty(t, info.Arch)
}

func TestGetAppVersion(t *testing.T) {
	version := GetAppVersion()
	assert.NotEmpty(t, version)
}

func TestGetShortVersion(t *testing.T) {
	version := GetShortVersion()
	assert.NotEmpty(t, version)
}

func TestGetFullVersion(t *testing.T) {
	version := GetFullVersion()
	assert.NotEmpty(t, version)
	assert.Contains(t, version, core.AppName)
}

func TestGetFormattedBuildInfo(t *testing.T) {
	info := GetFormattedBuildInfo()
	assert.NotEmpty(t, info)
	assert.Contains(t, info, "Version:")
	assert.Contains(t, info, "Git Commit:")
	assert.Contains(t, info, "Build Date:")
}

func TestIsDevBuild(t *testing.T) {
	// Should return true for default/dev builds
	isDev := IsDevBuild()
	assert.True(t, isDev) // Since we're running tests with default values
}

func TestGetUserAgent(t *testing.T) {
	ua := GetUserAgent()
	assert.NotEmpty(t, ua)
	assert.Contains(t, ua, core.AppName)
}

func TestParseVersion(t *testing.T) {
	tests := []struct {
		input    string
		expected Version
		hasError bool
	}{
		{"1.2.3", Version{Major: 1, Minor: 2, Patch: 3, Original: "1.2.3"}, false},
		{"v1.2.3", Version{Major: 1, Minor: 2, Patch: 3, Original: "v1.2.3"}, false},
		{"1.2.3-alpha.1", Version{Major: 1, Minor: 2, Patch: 3, PreRelease: "alpha.1", Original: "1.2.3-alpha.1"}, false},
		{"1.2.3+build.123", Version{Major: 1, Minor: 2, Patch: 3, Build: "build.123", Original: "1.2.3+build.123"}, false},
		{"latest", Version{Major: MaxVersionNumber, Minor: MaxVersionNumber, Patch: MaxVersionNumber, Original: "latest"}, false},
		{"", Version{}, true},
		{"invalid", Version{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := ParseVersion(tt.input)
			if tt.hasError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expected.Major, result.Major)
			assert.Equal(t, tt.expected.Minor, result.Minor)
			assert.Equal(t, tt.expected.Patch, result.Patch)
			assert.Equal(t, tt.expected.PreRelease, result.PreRelease)
			assert.Equal(t, tt.expected.Build, result.Build)
			assert.Equal(t, tt.expected.Original, result.Original)
		})
	}
}

func TestVersionCompare(t *testing.T) {
	tests := []struct {
		v1       string
		v2       string
		expected int
	}{
		{"1.2.3", "1.2.3", VersionEqual},
		{"1.2.4", "1.2.3", VersionNewer},
		{"1.2.2", "1.2.3", VersionOlder},
		{"1.3.0", "1.2.3", VersionNewer},
		{"1.1.0", "1.2.3", VersionOlder},
		{"2.0.0", "1.2.3", VersionNewer},
		{"0.9.0", "1.2.3", VersionOlder},
		{"1.2.3", "1.2.3-alpha", VersionNewer},
		{"1.2.3-alpha", "1.2.3", VersionOlder},
	}

	for _, tt := range tests {
		t.Run(tt.v1+"_vs_"+tt.v2, func(t *testing.T) {
			v1, err := ParseVersion(tt.v1)
			assert.NoError(t, err)
			v2, err := ParseVersion(tt.v2)
			assert.NoError(t, err)

			result := v1.Compare(*v2)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestVersionString(t *testing.T) {
	tests := []struct {
		version  Version
		expected string
	}{
		{Version{Major: 1, Minor: 2, Patch: 3}, "1.2.3"},
		{Version{Major: 1, Minor: 2, Patch: 3, PreRelease: "alpha"}, "1.2.3-alpha"},
		{Version{Major: 1, Minor: 2, Patch: 3, Build: "build"}, "1.2.3+build"},
		{Version{Major: 1, Minor: 2, Patch: 3, PreRelease: "alpha", Build: "build"}, "1.2.3-alpha+build"},
		{Version{Original: "custom"}, "custom"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.version.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseVersionConstraint(t *testing.T) {
	tests := []struct {
		input    string
		operator string
		version  string
		hasError bool
	}{
		{">=1.2.3", ">=", "1.2.3", false},
		{"1.2.3", "=", "1.2.3", false},
		{">1.0.0", ">", "1.0.0", false},
		{"*", "*", "", false},
		{"", "*", "", false},
		{"invalid", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := ParseVersionConstraint(tt.input)
			if tt.hasError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.operator, result.Operator)
			if tt.version != "" {
				assert.Equal(t, tt.version, result.Version.String())
			}
		})
	}
}

func TestVersionConstraintSatisfies(t *testing.T) {
	tests := []struct {
		constraint string
		version    string
		satisfies  bool
	}{
		{">=1.2.0", "1.2.3", true},
		{">=1.2.0", "1.1.0", false},
		{">1.2.0", "1.2.1", true},
		{">1.2.0", "1.2.0", false},
		{"<=1.2.0", "1.1.0", true},
		{"<=1.2.0", "1.3.0", false},
		{"<1.2.0", "1.1.0", true},
		{"<1.2.0", "1.2.0", false},
		{"=1.2.0", "1.2.0", true},
		{"=1.2.0", "1.2.1", false},
		{"*", "1.2.3", true},
	}

	for _, tt := range tests {
		t.Run(tt.constraint+"_"+tt.version, func(t *testing.T) {
			constraint, err := ParseVersionConstraint(tt.constraint)
			assert.NoError(t, err)

			version, err := ParseVersion(tt.version)
			assert.NoError(t, err)

			result := constraint.Satisfies(*version)
			assert.Equal(t, tt.satisfies, result)
		})
	}
}

func TestNewUpdateChecker(t *testing.T) {
	checker := NewUpdateChecker("1.0.0")
	assert.NotNil(t, checker)
	assert.Equal(t, "1.0.0", checker.currentVersion)
	assert.NotNil(t, checker.client)
}

func TestDetectProjectVersion(t *testing.T) {
	constraint, err := DetectProjectVersion("/tmp")
	assert.NoError(t, err)
	assert.NotNil(t, constraint)
	assert.Equal(t, "*", constraint.Operator)
}

func TestValidateProjectVersion(t *testing.T) {
	err := ValidateProjectVersion("/tmp")
	assert.NoError(t, err) // Should pass with wildcard constraint
}
