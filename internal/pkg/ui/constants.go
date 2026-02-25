package ui

// UI constants
const (
	IconOK      = "✓" // success, running, healthy
	IconFail    = "✗" // error, stopped, unhealthy
	IconWarn    = "!" // warning, starting, caution
	IconUnknown = "—" // not found, unknown, indeterminate

	ColorGreen   = "\033[32m"
	ColorRed     = "\033[31m"
	ColorYellow  = "\033[33m"
	ColorBlue    = "\033[34m"
	ColorCyan    = "\033[36m"
	ColorMagenta = "\033[35m"
	ColorGray    = "\033[90m"
	ColorBold    = "\033[1m"
	ColorReset   = "\033[0m"

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
	UIPadding             = 2
)
