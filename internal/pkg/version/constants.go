package version

// Version constants
const (
	DevelVersion        = "dev"
	MaxVersionNumber    = 999999
	GitCommitHashLength = 7
)

// Version comparison results
const (
	VersionOlder = -1
	VersionEqual = 0
	VersionNewer = 1
)

// Version defaults
const (
	DefaultVersion   = "dev"
	DefaultCommit    = "unknown"
	DefaultBuildDate = "unknown"
	DefaultBuildBy   = "unknown"
)
