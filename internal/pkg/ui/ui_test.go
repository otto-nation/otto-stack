//go:build unit

package ui

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewOutput(t *testing.T) {
	output := NewOutput()
	assert.NotNil(t, output)
	assert.False(t, output.Quiet)
	assert.False(t, output.NoColor)
}

func TestOutput_MessageFormatting(t *testing.T) {
	tests := []struct {
		name     string
		method   func(*Output, string, ...any)
		message  string
		args     []any
		expected string
	}{
		{"success", (*Output).Success, "Task completed on port %d", []any{5432}, "✅ Task completed on port 5432"},
		{"warning", (*Output).Warning, "Port %d in use", []any{5432}, "⚠️  Port 5432 in use"},
		{"info", (*Output).Info, "Connecting to %s", []any{"database"}, "ℹ️  Connecting to database"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := &Output{NoColor: true}

			// Capture stdout
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			tt.method(output, tt.message, tt.args...)

			w.Close()
			os.Stdout = old

			var buf bytes.Buffer
			io.Copy(&buf, r)

			assert.Contains(t, buf.String(), tt.expected)
		})
	}
}

func TestOutput_QuietMode(t *testing.T) {
	output := &Output{Quiet: true}

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	output.Success("test message")
	output.Warning("test warning")
	output.Info("test info")

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)

	assert.Empty(t, buf.String(), "Quiet mode should suppress output")
}

func TestOutput_NoColorMode(t *testing.T) {
	output := &Output{NoColor: true}

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	output.Success("Success message")
	output.Header("Header message")

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	result := buf.String()

	assert.Contains(t, result, "✅ Success message")
	assert.Contains(t, result, "=== Header message ===")
	assert.NotContains(t, result, "\033[", "No color mode should not contain ANSI codes")
}

func TestOutput_List(t *testing.T) {
	output := NewOutput()
	items := []string{"item1", "item2", "item3"}

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	output.List(items)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	result := buf.String()

	for _, item := range items {
		assert.Contains(t, result, "• "+item)
	}
}

func TestOutput_Box(t *testing.T) {
	output := &Output{NoColor: true}

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	output.Box("Title", "Content")

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	result := buf.String()

	assert.Contains(t, result, "Title")
	assert.Contains(t, result, "Content")
}

func TestFormatFunctions(t *testing.T) {
	tests := []struct {
		name     string
		function func(string, bool) string
		input    string
		noColor  bool
		contains string
	}{
		{"success with color", formatSuccess, "test", false, "✅"},
		{"success no color", formatSuccess, "test", true, "✅ test"},
		{"error with color", formatError, "test", false, "❌"},
		{"error no color", formatError, "test", true, "❌ test"},
		{"warning with color", formatWarning, "test", false, "⚠️"},
		{"info with color", formatInfo, "test", false, "ℹ️"},
		{"header no color", formatHeader, "test", true, "=== test ==="},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.function(tt.input, tt.noColor)
			assert.Contains(t, result, tt.contains)
		})
	}
}

func TestDefaultOutput(t *testing.T) {
	assert.NotNil(t, DefaultOutput)
	assert.IsType(t, &Output{}, DefaultOutput)
}
