package ui

// formatColored applies color if color is enabled
func formatColored(text, color string, noColor bool) string {
	if noColor {
		return text
	}
	return color + text + ColorReset
}

// formatSuccess formats success messages
func formatSuccess(text string, noColor bool) string {
	return formatColored(IconOK+" "+text, ColorGreen+ColorBold, noColor)
}

// formatError formats error messages
func formatError(text string, noColor bool) string {
	return formatColored(IconFail+" "+text, ColorRed+ColorBold, noColor)
}

// formatWarning formats warning messages
func formatWarning(text string, noColor bool) string {
	return formatColored(IconWarn+" "+text, ColorYellow+ColorBold, noColor)
}

// formatInfo formats info messages — no prefix symbol, color is sufficient
func formatInfo(text string, noColor bool) string {
	return formatColored(text, ColorBlue, noColor)
}

// formatHeader formats header messages
func formatHeader(text string, noColor bool) string {
	if noColor {
		return "\n" + text + "\n"
	}
	return formatColored("\n"+text+"\n", ColorGreen+ColorBold, noColor)
}

// formatMuted formats muted text
func formatMuted(text string, noColor bool) string {
	return formatColored(text, ColorGray, noColor)
}
