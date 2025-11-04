package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/otto-nation/otto-stack/internal/pkg/constants"
)

var (
	// Colors
	primaryColor = lipgloss.Color("#00D4AA")
	successColor = lipgloss.Color("#00C851")
	warningColor = lipgloss.Color("#FF8800")
	errorColor   = lipgloss.Color("#FF4444")
	infoColor    = lipgloss.Color("#33B5E5")
	mutedColor   = lipgloss.Color("#666666")

	// Base styles
	baseStyle = lipgloss.NewStyle().
			Padding(0, 1)

	// Message styles
	SuccessStyle = baseStyle.
			Foreground(successColor).
			Bold(true)

	ErrorStyle = baseStyle.
			Foreground(errorColor).
			Bold(true)

	WarningStyle = baseStyle.
			Foreground(warningColor).
			Bold(true)

	InfoStyle = baseStyle.
			Foreground(infoColor)

	MutedStyle = baseStyle.
			Foreground(mutedColor)

	// Header styles
	HeaderStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Padding(1, constants.UIPadding).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor)

	SubHeaderStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Padding(0, 1)

	// List styles
	ListItemStyle = lipgloss.NewStyle().
			Padding(0, constants.UIPadding)

	SelectedItemStyle = ListItemStyle.
				Foreground(primaryColor).
				Bold(true)

	// Box styles
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(mutedColor).
			Padding(1, constants.UIPadding).
			Margin(1, 0)

	HighlightBoxStyle = BoxStyle.
				BorderForeground(primaryColor)
)
