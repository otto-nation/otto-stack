package version

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/otto-nation/otto-stack/internal/core"
)

// Build variables set by ldflags
var (
	AppVersion = DefaultVersion
	GitCommit  = DefaultCommit
	BuildDate  = DefaultBuildDate
	BuildBy    = DefaultBuildBy
)

// BuildInfo contains comprehensive build information
type BuildInfo struct {
	Version   string    `json:"version"`
	GitCommit string    `json:"git_commit"`
	BuildDate string    `json:"build_date"`
	BuildBy   string    `json:"build_by"`
	GoVersion string    `json:"go_version"`
	Platform  string    `json:"platform"`
	Arch      string    `json:"arch"`
	BuildTime time.Time `json:"build_time"`
}

// GetBuildInfo returns comprehensive build information
// GetAppVersion returns the application version string
func GetAppVersion() string {
	if AppVersion != DefaultVersion {
		return AppVersion
	}

	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return AppVersion
	}

	return buildInfo.Main.Version
}

// GetShortVersion returns a short version string
func GetShortVersion() string {
	version := GetAppVersion()
	if version == DevelVersion || version == DefaultVersion {
		return DefaultVersion
	}
	return version
}

// IsDevBuild returns true if this is a development build
func IsDevBuild() bool {
	return AppVersion == DefaultVersion || AppVersion == DevelVersion
}

// GetUserAgent returns a user agent string for HTTP requests
func GetUserAgent() string {
	return fmt.Sprintf("%s/%s (%s; %s)", core.AppName, GetShortVersion(), runtime.GOOS, runtime.GOARCH)
}
