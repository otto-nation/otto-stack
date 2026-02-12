package ui

import (
	"fmt"
	"io"
	"os"
	"time"

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

// Writer returns the underlying writer (stdout)
func (o *Output) Writer() io.Writer {
	return os.Stdout
}

// Success prints a success message
func (o *Output) Success(msg string, args ...any) {
	o.print(formatSuccess(fmt.Sprintf(msg, args...), o.NoColor), "Success")
}

// Error prints an error message
func (o *Output) Error(msg string, args ...any) {
	if o.Quiet {
		return
	}
	formatted := formatError(fmt.Sprintf(msg, args...), o.NoColor)
	fmt.Fprintln(os.Stderr, formatted)
}

// Warning prints a warning message
func (o *Output) Warning(msg string, args ...any) {
	o.print(formatWarning(fmt.Sprintf(msg, args...), o.NoColor), "Warning")
}

// Info prints an info message
func (o *Output) Info(msg string, args ...any) {
	o.print(formatInfo(fmt.Sprintf(msg, args...), o.NoColor), "Info")
}

// Header prints a styled header
func (o *Output) Header(msg string, args ...any) {
	o.print(formatHeader(fmt.Sprintf(msg, args...), o.NoColor), "Header")
}

// SubHeader prints a styled sub-header
func (o *Output) SubHeader(msg string, args ...any) {
	formatted := fmt.Sprintf(msg, args...)
	if o.NoColor {
		o.print("\n--- "+formatted+" ---\n", "SubHeader")
	} else {
		o.print(formatColored("\n"+IconBox+" "+formatted+"\n", ColorGreen+ColorBold, o.NoColor), "SubHeader")
	}
}

// List prints a styled list
func (o *Output) List(items []string) {
	if o.Quiet {
		return
	}
	logger.Info("UI List", "items", items)
	for _, item := range items {
		_, _ = fmt.Fprintf(os.Stdout, "  • %s\n", item)
	}
}

// Progress shows a spinner with message
func (o *Output) Progress(msg string, fn func() error) error {
	if o.Quiet || o.NoColor {
		return fn()
	}

	s := spinner.New(spinner.CharSets[14], SpinnerIntervalMilliseconds*time.Millisecond)
	s.Suffix = " " + msg
	s.Start()
	defer s.Stop()

	return fn()
}

// Muted prints muted text
func (o *Output) Muted(msg string, args ...any) {
	o.print(formatMuted(fmt.Sprintf(msg, args...), o.NoColor), "Muted")
}

// Box prints content in a styled box
func (o *Output) Box(title, content string) {
	if o.Quiet {
		return
	}
	logger.Info("UI Box", "title", title, "content", content)
	if o.NoColor {
		_, _ = fmt.Fprintf(os.Stdout, "\n┌─ %s ─\n│ %s\n└─\n", title, content)
	} else {
		_, _ = fmt.Fprintf(os.Stdout, "\n%s\n\n%s\n", formatColored(title, ColorGreen+ColorBold, false), content)
	}
}

// print is a helper method to reduce repetition
func (o *Output) print(formatted, logType string) {
	if o.Quiet {
		return
	}
	logger.Info("UI "+logType, "message", formatted)
	_, _ = fmt.Fprintln(os.Stdout, formatted)
}

// Global output instance
var DefaultOutput = NewOutput()
