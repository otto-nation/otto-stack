//go:build unit

package ui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServiceOption_Fields(t *testing.T) {
	t.Run("validates ServiceOption structure", func(t *testing.T) {
		option := ServiceOption{
			Name:         "postgres",
			Description:  "PostgreSQL database",
			Dependencies: []string{"network"},
			Category:     "database",
		}

		assert.Equal(t, "postgres", option.Name)
		assert.Equal(t, "PostgreSQL database", option.Description)
		assert.Len(t, option.Dependencies, 1)
		assert.Equal(t, "database", option.Category)
	})
}

func TestUIConstants(t *testing.T) {
	t.Run("validates icon constants", func(t *testing.T) {
		assert.NotEmpty(t, IconSuccess)
		assert.NotEmpty(t, IconError)
		assert.NotEmpty(t, IconWarning)
		assert.NotEmpty(t, IconInfo)
		assert.NotEmpty(t, IconBox)
		assert.NotEmpty(t, IconHeader)
	})

	t.Run("validates color constants", func(t *testing.T) {
		assert.NotEmpty(t, ColorGreen)
		assert.NotEmpty(t, ColorRed)
		assert.NotEmpty(t, ColorYellow)
		assert.NotEmpty(t, ColorBlue)
		assert.NotEmpty(t, ColorGray)
		assert.NotEmpty(t, ColorBold)
		assert.NotEmpty(t, ColorReset)
	})

	t.Run("validates display constants", func(t *testing.T) {
		assert.Greater(t, SeparatorLength, 0)
		assert.Greater(t, StatusSeparatorLength, 0)
		assert.Greater(t, TableWidth42, 0)
		assert.Greater(t, TableWidth75, 0)
		assert.Greater(t, TableWidth80, 0)
		assert.Greater(t, TableWidth85, 0)
		assert.Greater(t, TableWidth90, 0)
		assert.Greater(t, UIPadding, 0)
		assert.Greater(t, SpinnerIntervalMilliseconds, 0)
	})

	t.Run("validates status constants", func(t *testing.T) {
		assert.NotEmpty(t, StatusHeaderService)
		assert.NotEmpty(t, StatusHeaderProvidedBy)
		assert.NotEmpty(t, StatusHeaderState)
		assert.NotEmpty(t, StatusHeaderHealth)
		assert.NotEmpty(t, StatusSeparator)
		assert.NotEmpty(t, StatusSeparatorWithProvider)
	})
}

func TestStyleFunctions(t *testing.T) {
	t.Run("formats colored text", func(t *testing.T) {
		result := formatColored("test", ColorGreen, false)
		assert.Contains(t, result, "test")

		// With no color
		resultNoColor := formatColored("test", ColorGreen, true)
		assert.Equal(t, "test", resultNoColor)
	})

	t.Run("formats success text", func(t *testing.T) {
		result := formatSuccess("success", false)
		assert.Contains(t, result, "success")
		assert.Contains(t, result, IconSuccess)

		resultNoColor := formatSuccess("success", true)
		assert.Contains(t, resultNoColor, "success")
		assert.Contains(t, resultNoColor, IconSuccess)
	})

	t.Run("formats error text", func(t *testing.T) {
		result := formatError("error", false)
		assert.Contains(t, result, "error")
		assert.Contains(t, result, IconError)

		resultNoColor := formatError("error", true)
		assert.Contains(t, resultNoColor, "error")
		assert.Contains(t, resultNoColor, IconError)
	})

	t.Run("formats warning text", func(t *testing.T) {
		result := formatWarning("warning", false)
		assert.Contains(t, result, "warning")
		assert.Contains(t, result, IconWarning)

		resultNoColor := formatWarning("warning", true)
		assert.Contains(t, resultNoColor, "warning")
		assert.Contains(t, resultNoColor, IconWarning)
	})

	t.Run("formats info text", func(t *testing.T) {
		result := formatInfo("info", false)
		assert.Contains(t, result, "info")
		assert.Contains(t, result, IconInfo)

		resultNoColor := formatInfo("info", true)
		assert.Contains(t, resultNoColor, "info")
		assert.Contains(t, resultNoColor, IconInfo)
	})

	t.Run("formats header text", func(t *testing.T) {
		result := formatHeader("header", false)
		assert.Contains(t, result, "header")
		assert.Contains(t, result, IconHeader)

		resultNoColor := formatHeader("header", true)
		assert.Contains(t, resultNoColor, "header")
		assert.Contains(t, resultNoColor, "===")
	})

	t.Run("formats muted text", func(t *testing.T) {
		result := formatMuted("muted", false)
		assert.Contains(t, result, "muted")

		resultNoColor := formatMuted("muted", true)
		assert.Equal(t, "muted", resultNoColor)
	})
}

func TestOutput_AdditionalMethods(t *testing.T) {
	t.Run("tests box method", func(t *testing.T) {
		output := NewOutput()
		// Box method should not panic
		output.Box("Title", "Content")
		assert.NotNil(t, output)
	})

	t.Run("tests subheader method", func(t *testing.T) {
		output := NewOutput()
		// SubHeader method should not panic
		output.SubHeader("Test SubHeader")
		assert.NotNil(t, output)
	})

	t.Run("tests list method", func(t *testing.T) {
		output := NewOutput()
		items := []string{"item1", "item2", "item3"}
		// List method should not panic
		output.List(items)
		assert.NotNil(t, output)
	})

	t.Run("tests progress method", func(t *testing.T) {
		output := NewOutput()
		// Progress method should not panic
		err := output.Progress("Testing...", func() error {
			return nil
		})
		assert.NoError(t, err)
	})

	t.Run("tests default output instance", func(t *testing.T) {
		assert.NotNil(t, DefaultOutput)
		assert.IsType(t, &Output{}, DefaultOutput)
	})
}
