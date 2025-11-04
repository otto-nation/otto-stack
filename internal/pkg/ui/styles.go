package ui

import "github.com/otto-nation/otto-stack/internal/pkg/constants"

// formatColored applies color if color is enabled
func formatColored(text, color string, noColor bool) string {
	if noColor {
		return text
	}
	return color + text + constants.ColorReset
}

// formatSuccess formats success messages
func formatSuccess(text string, noColor bool) string {
	return formatColored(constants.IconSuccess+" "+text, constants.ColorGreen+constants.ColorBold, noColor)
}

// formatError formats error messages
func formatError(text string, noColor bool) string {
	return formatColored(constants.IconError+" "+text, constants.ColorRed+constants.ColorBold, noColor)
}

// formatWarning formats warning messages
func formatWarning(text string, noColor bool) string {
	return formatColored(constants.IconWarning+" "+text, constants.ColorYellow+constants.ColorBold, noColor)
}

// formatInfo formats info messages
func formatInfo(text string, noColor bool) string {
	return formatColored(constants.IconInfo+" "+text, constants.ColorBlue, noColor)
}

// formatHeader formats header messages
func formatHeader(text string, noColor bool) string {
	if noColor {
		return "\n=== " + text + " ===\n"
	}
	return formatColored("\n"+constants.IconHeader+" "+text+"\n", constants.ColorGreen+constants.ColorBold, noColor)
}

// formatMuted formats muted text
func formatMuted(text string, noColor bool) string {
	return formatColored(text, constants.ColorGray, noColor)
}
