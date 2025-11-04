package ui

import (
	"fmt"
	"os"
	"time"

	"github.com/otto-nation/otto-stack/internal/pkg/constants"

	"github.com/briandowns/spinner"
	"github.com/otto-nation/otto-stack/internal/pkg/logger"
)

// Output handles all user-facing output with consistent styling
type Output struct {
	Quiet   bool
	NoColor bool
}

// NewOutput creates a new output handler
func NewOutput() *Output {
	return &Output{}
}

// Success prints a success message
func (o *Output) Success(msg string, args ...any) {
	if o.Quiet {
		return
	}
	formatted := fmt.Sprintf(msg, args...)
	logger.Info("Success", "message", formatted)
	if o.NoColor {
		fmt.Printf("✅ %s\n", formatted)
	} else {
		fmt.Println(SuccessStyle.Render("✅ " + formatted))
	}
}

// Error prints an error message
func (o *Output) Error(msg string, args ...any) {
	formatted := fmt.Sprintf(msg, args...)
	logger.Error("UI Error", "message", formatted)
	if o.NoColor {
		fmt.Fprintf(os.Stderr, "❌ %s\n", formatted)
	} else {
		fmt.Fprintln(os.Stderr, ErrorStyle.Render("❌ "+formatted))
	}
}

// Warning prints a warning message
func (o *Output) Warning(msg string, args ...any) {
	if o.Quiet {
		return
	}
	formatted := fmt.Sprintf(msg, args...)
	logger.Warn("UI Warning", "message", formatted)
	if o.NoColor {
		fmt.Printf("⚠️  %s\n", formatted)
	} else {
		fmt.Println(WarningStyle.Render("⚠️  " + formatted))
	}
}

// Info prints an info message
func (o *Output) Info(msg string, args ...any) {
	if o.Quiet {
		return
	}
	formatted := fmt.Sprintf(msg, args...)
	logger.Info("UI Info", "message", formatted)
	if o.NoColor {
		fmt.Printf("ℹ️  %s\n", formatted)
	} else {
		fmt.Println(InfoStyle.Render("ℹ️  " + formatted))
	}
}

// Header prints a styled header
func (o *Output) Header(msg string, args ...any) {
	if o.Quiet {
		return
	}
	formatted := fmt.Sprintf(msg, args...)
	logger.Info("UI Header", "message", formatted)
	if o.NoColor {
		fmt.Printf("\n=== %s ===\n\n", formatted)
	} else {
		fmt.Println(HeaderStyle.Render("🚀 " + formatted))
	}
}

// SubHeader prints a styled sub-header
func (o *Output) SubHeader(msg string, args ...any) {
	if o.Quiet {
		return
	}
	formatted := fmt.Sprintf(msg, args...)
	logger.Info("UI SubHeader", "message", formatted)
	if o.NoColor {
		fmt.Printf("\n--- %s ---\n", formatted)
	} else {
		fmt.Println(SubHeaderStyle.Render("📦 " + formatted))
	}
}

// List prints a styled list
func (o *Output) List(items []string) {
	if o.Quiet {
		return
	}
	logger.Info("UI List", "items", items)
	for _, item := range items {
		if o.NoColor {
			fmt.Printf("  • %s\n", item)
		} else {
			fmt.Println(ListItemStyle.Render("• " + item))
		}
	}
}

// Progress shows a spinner with message
func (o *Output) Progress(msg string, fn func() error) error {
	if o.Quiet || o.NoColor {
		return fn()
	}

	s := spinner.New(spinner.CharSets[14], constants.SpinnerIntervalMilliseconds*time.Millisecond)
	s.Suffix = " " + msg
	s.Start()
	defer s.Stop()

	return fn()
}

// Muted prints muted text
func (o *Output) Muted(msg string, args ...any) {
	if o.Quiet {
		return
	}
	formatted := fmt.Sprintf(msg, args...)
	logger.Debug("UI Muted", "message", formatted)
	if o.NoColor {
		fmt.Printf("  %s\n", formatted)
	} else {
		fmt.Println(MutedStyle.Render(formatted))
	}
}

// Box prints content in a styled box
func (o *Output) Box(title, content string) {
	if o.Quiet {
		return
	}
	logger.Info("UI Box", "title", title, "content", content)
	if o.NoColor {
		fmt.Printf("\n┌─ %s ─\n│ %s\n└─\n", title, content)
	} else {
		boxContent := fmt.Sprintf("%s\n\n%s", SubHeaderStyle.Render(title), content)
		fmt.Println(BoxStyle.Render(boxContent))
	}
}

// Global output instance
var DefaultOutput = NewOutput()

// Convenience functions for global use
func Success(msg string, args ...any)            { DefaultOutput.Success(msg, args...) }
func Error(msg string, args ...any)              { DefaultOutput.Error(msg, args...) }
func Warning(msg string, args ...any)            { DefaultOutput.Warning(msg, args...) }
func Info(msg string, args ...any)               { DefaultOutput.Info(msg, args...) }
func Header(msg string, args ...any)             { DefaultOutput.Header(msg, args...) }
func SubHeader(msg string, args ...any)          { DefaultOutput.SubHeader(msg, args...) }
func List(items []string)                        { DefaultOutput.List(items) }
func Progress(msg string, fn func() error) error { return DefaultOutput.Progress(msg, fn) }
func Muted(msg string, args ...any)              { DefaultOutput.Muted(msg, args...) }
func Box(title, content string)                  { DefaultOutput.Box(title, content) }
