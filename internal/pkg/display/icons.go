package display

import "github.com/otto-nation/otto-stack/internal/pkg/constants"

// StateIcons maps container states to display icons
var StateIcons = map[string]string{
	constants.StateRunning: "🟢",
	constants.StateStopped: "🔴",
	constants.StateCreated: "🟡",
	"starting":             "🟡",
	"paused":               "⏸️",
}

// HealthIcons maps health statuses to display icons
var HealthIcons = map[string]string{
	constants.HealthHealthy:   "✅",
	constants.HealthUnhealthy: "❌",
	constants.HealthStarting:  "🟡",
	constants.HealthNone:      "❓",
}

// GetStateIcon returns the icon for a given state
func GetStateIcon(state string) string {
	if icon, exists := StateIcons[state]; exists {
		return icon
	}
	return "⚪" // default
}

// GetHealthIcon returns the icon for a given health status
func GetHealthIcon(health string) string {
	if icon, exists := HealthIcons[health]; exists {
		return icon
	}
	return "❓" // default
}
