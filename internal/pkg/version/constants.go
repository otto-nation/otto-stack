package version

import "time"

// Version change types (used for both drift and updates)
const (
	ChangeTypeNone       = "none"
	ChangeTypeMajor      = "major"
	ChangeTypeMinor      = "minor"
	ChangeTypePatch      = "patch"
	ChangeTypePrerelease = "prerelease"
)

// Severity levels with their numeric values for comparison
var SeverityLevels = map[string]int{
	"low":    1, // Optional updates, minor drift
	"medium": 2, // Recommended updates, moderate drift
	"high":   3, // Critical updates, major drift
}

// Enforcement actions
const (
	EnforcementActionNone    = "none"
	EnforcementActionWarn    = "warn"
	EnforcementActionSwitch  = "switch"
	EnforcementActionInstall = "install"
)

// Default durations
const (
	DefaultMaxDriftDuration = 7 * 24 * time.Hour // 1 week
	DefaultCheckInterval    = 24 * time.Hour     // 1 day
)

// Default enforcement policy
const (
	DefaultStrictMode     = false
	DefaultAllowDrift     = true
	DefaultAutoSync       = false
	DefaultNotifyUpdates  = true
	DefaultShowPrerelease = false
	DefaultAutoCheck      = true
)

// Version source types
const (
	VersionSourceGitHub = "github"
	VersionSourceLocal  = "local"
	VersionSourceBinary = "binary"
)
