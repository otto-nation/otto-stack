package version

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/otto-nation/otto-stack/internal/pkg/constants"
)

// Default values
const (
	DefaultVersion   = "dev"
	DefaultCommit    = "unknown"
	DefaultBuildDate = "unknown"
	DefaultBuildBy   = "unknown"
	DevelVersion     = "(devel)"
)

// Use constants from brand.go
var (
	AppNameTemplate   = constants.AppNameTemplate
	UserAgentTemplate = constants.UserAgentTemplate
)

// Version-specific templates (not branding)
const (
	VersionInfoTemplate = `Version:    %s
Git Commit: %s
Build Date: %s
Built By:   %s
Go Version: %s
Platform:   %s/%s
`
)

// Version comparison results
const (
	VersionEqual   = 0
	VersionNewer   = 1
	VersionOlder   = -1
	VersionInvalid = -999
)

// Special version values
var SpecialVersions = map[string]struct{}{
	"latest": {},
	"*":      {},
}

var (
	// These variables are set by the build process using ldflags
	AppVersion = DefaultVersion
	GitCommit  = DefaultCommit
	BuildDate  = DefaultBuildDate
	BuildBy    = DefaultBuildBy
)

// BuildInfo contains comprehensive build information
type BuildInfo struct {
	Version    string           `json:"version"`
	GitCommit  string           `json:"git_commit"`
	BuildDate  string           `json:"build_date"`
	BuildBy    string           `json:"build_by"`
	GoVersion  string           `json:"go_version"`
	Platform   string           `json:"platform"`
	Arch       string           `json:"arch"`
	BuildTime  time.Time        `json:"build_time"`
	ModuleInfo *debug.BuildInfo `json:"module_info,omitempty"`
}

// GetBuildInfo returns comprehensive build information
func GetBuildInfo() *BuildInfo {
	buildTime, _ := time.Parse(time.RFC3339, BuildDate)
	if buildTime.IsZero() {
		buildTime = time.Now()
	}

	info := &BuildInfo{
		Version:   AppVersion,
		GitCommit: GitCommit,
		BuildDate: BuildDate,
		BuildBy:   BuildBy,
		GoVersion: runtime.Version(),
		Platform:  runtime.GOOS,
		Arch:      runtime.GOARCH,
		BuildTime: buildTime,
	}

	// Try to get module information
	if buildInfo, ok := debug.ReadBuildInfo(); ok {
		info.ModuleInfo = buildInfo
	}

	return info
}

// GetAppVersion returns the application version string
func GetAppVersion() string {
	if AppVersion == DefaultVersion {
		// Try to get version from module info
		if buildInfo, ok := debug.ReadBuildInfo(); ok {
			return buildInfo.Main.Version
		}
	}
	return AppVersion
}

// GetShortVersion returns a short version string
func GetShortVersion() string {
	version := GetAppVersion()
	if version == DevelVersion || version == DefaultVersion {
		return DefaultVersion
	}
	return version
}

// GetFullVersion returns a detailed version string
func GetFullVersion() string {
	info := GetBuildInfo()

	version := info.Version
	if version == DefaultVersion || version == DevelVersion {
		version = DefaultVersion
	}

	result := fmt.Sprintf(AppNameTemplate, version)

	if info.GitCommit != DefaultCommit && info.GitCommit != "" {
		if len(info.GitCommit) > 7 {
			result += fmt.Sprintf(" (%s)", info.GitCommit[:7])
		} else {
			result += fmt.Sprintf(" (%s)", info.GitCommit)
		}
	}

	return result
}

// GetFormattedBuildInfo returns formatted build information
func GetFormattedBuildInfo() string {
	info := GetBuildInfo()
	return fmt.Sprintf(VersionInfoTemplate,
		info.Version, info.GitCommit, info.BuildDate,
		info.BuildBy, info.GoVersion, info.Platform, info.Arch)
}

// IsDevBuild returns true if this is a development build
func IsDevBuild() bool {
	return AppVersion == DefaultVersion || AppVersion == DevelVersion
}

// GetUserAgent returns a user agent string for HTTP requests
func GetUserAgent() string {
	return fmt.Sprintf(UserAgentTemplate, GetShortVersion(), runtime.GOOS, runtime.GOARCH)
}

// IsAppVersionCompatible checks if the current version is compatible with required version
func IsAppVersionCompatible(requiredVersion string) bool {
	currentVersion := GetShortVersion()

	// Development builds are always compatible
	if IsDevBuild() {
		return true
	}

	// For now, use simple string comparison
	// In production, you'd implement proper semantic versioning
	return currentVersion >= requiredVersion
}
