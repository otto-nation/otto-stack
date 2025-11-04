package constants

// File and directory permissions
const (
	FilePermReadWrite    = 0644 // Standard file read/write permissions
	DirPermReadWriteExec = 0755 // Standard directory permissions
	FilePermReadWriteAll = 0666 // File permissions for all users
)

// Validation limits
const (
	MinProjectNameLength = 2
	MaxProjectNameLength = 50
	MinArgumentCount     = 2 // For command parsing
	MinFieldCount        = 2 // For field parsing
	MaxCategoryCommands  = 10
	GitCommitHashLength  = 7
)

// Timeouts and intervals (in seconds)
const (
	DefaultStopTimeoutSeconds   = 10
	DefaultStartTimeoutSeconds  = 30
	HealthCheckIntervalSeconds  = 2
	SpinnerIntervalMilliseconds = 100
)

// Display formatting
const (
	SeparatorLength       = 50
	StatusSeparatorLength = 45
	TableWidth42          = 42
	TableWidth75          = 75
	TableWidth80          = 80
	TableWidth85          = 85
	TableWidth90          = 90
	HoursPerDay           = 24
)

// Validation thresholds (percentages)
const (
	PercentageMultiplier   = 100
	MinExampleCoverage     = 80
	MinTipsCoverage        = 50
	MinDescriptionCoverage = 60
	MinConfigurationScore  = 80
	MaxValidationErrors    = 5
	BaseValidationScore    = 100.0
	ErrorWeight            = 10
	WarningWeight          = 2
)

// Version and parsing
const (
	MaxVersionNumber = 999
	KeyValueParts    = 2 // For splitting "key=value" strings
	PortSearchRange  = 100
	HexDivisor       = 2 // For hex string conversion
)

// Version defaults
const (
	DefaultVersion   = "dev"
	DefaultCommit    = "unknown"
	DefaultBuildDate = "unknown"
	DefaultBuildBy   = "unknown"
	DevelVersion     = "(devel)"
)

// Version comparison results
const (
	VersionEqual   = 0
	VersionNewer   = 1
	VersionOlder   = -1
	VersionInvalid = -999
)

// UI padding and spacing
const (
	UIPadding = 2
)

// UI ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorGreen  = "\033[32m"
	ColorRed    = "\033[31m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorGray   = "\033[90m"
	ColorBold   = "\033[1m"
)

// UI message prefixes
const (
	IconSuccess = "✅"
	IconError   = "❌"
	IconWarning = "⚠️ "
	IconInfo    = "ℹ️ "
	IconHeader  = "🚀"
	IconBox     = "📦"
)
