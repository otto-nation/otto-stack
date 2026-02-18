//go:build unit

package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAppVersion(t *testing.T) {
	version := GetAppVersion()
	assert.NotEmpty(t, version)
}

func TestGetShortVersion(t *testing.T) {
	version := GetShortVersion()
	assert.NotEmpty(t, version)
}

func TestIsDevBuild(t *testing.T) {
	isDev := IsDevBuild()
	assert.True(t, isDev) // In test environment, should be dev build
}

func TestGetUserAgent(t *testing.T) {
	userAgent := GetUserAgent()
	assert.NotEmpty(t, userAgent)
	assert.Contains(t, userAgent, "otto-stack")
}
