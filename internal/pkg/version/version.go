package version

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/otto-nation/otto-stack/internal/pkg/constants"
)

// Build variables set by ldflags
var (
	AppVersion = constants.DefaultVersion
	GitCommit  = constants.DefaultCommit
	BuildDate  = constants.DefaultBuildDate
	BuildBy    = constants.DefaultBuildBy
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
func GetBuildInfo() *BuildInfo {
	buildTime, _ := time.Parse(time.RFC3339, BuildDate)
	if buildTime.IsZero() {
		buildTime = time.Now()
	}

	return &BuildInfo{
		Version:   AppVersion,
		GitCommit: GitCommit,
		BuildDate: BuildDate,
		BuildBy:   BuildBy,
		GoVersion: runtime.Version(),
		Platform:  runtime.GOOS,
		Arch:      runtime.GOARCH,
		BuildTime: buildTime,
	}
}

// GetAppVersion returns the application version string
func GetAppVersion() string {
	if AppVersion == constants.DefaultVersion {
		if buildInfo, ok := debug.ReadBuildInfo(); ok {
			return buildInfo.Main.Version
		}
	}
	return AppVersion
}

// GetShortVersion returns a short version string
func GetShortVersion() string {
	version := GetAppVersion()
	if version == constants.DevelVersion || version == constants.DefaultVersion {
		return constants.DefaultVersion
	}
	return version
}

// GetFullVersion returns a detailed version string
func GetFullVersion() string {
	info := GetBuildInfo()
	version := info.Version
	if version == constants.DefaultVersion || version == constants.DevelVersion {
		version = constants.DefaultVersion
	}

	result := fmt.Sprintf("%s %s", constants.AppName, version)
	if info.GitCommit != constants.DefaultCommit && info.GitCommit != "" {
		if len(info.GitCommit) > constants.GitCommitHashLength {
			result += fmt.Sprintf(" (%s)", info.GitCommit[:constants.GitCommitHashLength])
		} else {
			result += fmt.Sprintf(" (%s)", info.GitCommit)
		}
	}
	return result
}

// GetFormattedBuildInfo returns formatted build information
func GetFormattedBuildInfo() string {
	info := GetBuildInfo()
	return fmt.Sprintf(`Version:    %s
Git Commit: %s
Build Date: %s
Built By:   %s
Go Version: %s
Platform:   %s/%s`,
		info.Version, info.GitCommit, info.BuildDate,
		info.BuildBy, info.GoVersion, info.Platform, info.Arch)
}

// IsDevBuild returns true if this is a development build
func IsDevBuild() bool {
	return AppVersion == constants.DefaultVersion || AppVersion == constants.DevelVersion
}

// GetUserAgent returns a user agent string for HTTP requests
func GetUserAgent() string {
	return fmt.Sprintf("%s/%s (%s; %s)", constants.AppName, GetShortVersion(), runtime.GOOS, runtime.GOARCH)
}
