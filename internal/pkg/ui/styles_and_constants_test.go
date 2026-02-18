//go:build unit

package ui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUIConstants(t *testing.T) {
	assert.NotEmpty(t, IconSuccess)
	assert.NotEmpty(t, IconError)
	assert.NotEmpty(t, IconWarning)
	assert.NotEmpty(t, IconInfo)
	assert.NotEmpty(t, IconBox)
	assert.NotEmpty(t, IconHeader)

	assert.NotEmpty(t, ColorGreen)
	assert.NotEmpty(t, ColorRed)
	assert.NotEmpty(t, ColorYellow)
	assert.NotEmpty(t, ColorBlue)
	assert.NotEmpty(t, ColorGray)
	assert.NotEmpty(t, ColorBold)
	assert.NotEmpty(t, ColorReset)

	assert.Greater(t, SeparatorLength, 0)
	assert.Greater(t, StatusSeparatorLength, 0)
	assert.Greater(t, TableWidth42, 0)
	assert.Greater(t, TableWidth75, 0)
	assert.Greater(t, TableWidth80, 0)
	assert.Greater(t, TableWidth85, 0)
	assert.Greater(t, TableWidth90, 0)
	assert.Greater(t, UIPadding, 0)
	assert.Greater(t, SpinnerIntervalMilliseconds, 0)
}

func TestStyleFunctions(t *testing.T) {
	result := formatColored("test", ColorGreen, false)
	assert.Contains(t, result, "test")

	resultNoColor := formatColored("test", ColorGreen, true)
	assert.Equal(t, "test", resultNoColor)

	result = formatSuccess("success", false)
	assert.Contains(t, result, "success")
	assert.Contains(t, result, IconSuccess)

	resultNoColor = formatSuccess("success", true)
	assert.Contains(t, resultNoColor, "success")
	assert.Contains(t, resultNoColor, IconSuccess)

	result = formatError("error", false)
	assert.Contains(t, result, "error")
	assert.Contains(t, result, IconError)

	resultNoColor = formatError("error", true)
	assert.Contains(t, resultNoColor, "error")
	assert.Contains(t, resultNoColor, IconError)

	result = formatWarning("warning", false)
	assert.Contains(t, result, "warning")
	assert.Contains(t, result, IconWarning)

	resultNoColor = formatWarning("warning", true)
	assert.Contains(t, resultNoColor, "warning")
	assert.Contains(t, resultNoColor, IconWarning)

	result = formatInfo("info", false)
	assert.Contains(t, result, "info")
	assert.Contains(t, result, IconInfo)

	resultNoColor = formatInfo("info", true)
	assert.Contains(t, resultNoColor, "info")
	assert.Contains(t, resultNoColor, IconInfo)

	result = formatHeader("header", false)
	assert.Contains(t, result, "header")
	assert.Contains(t, result, IconHeader)

	resultNoColor = formatHeader("header", true)
	assert.Contains(t, resultNoColor, "header")
	assert.Contains(t, resultNoColor, "===")

	result = formatMuted("muted", false)
	assert.Contains(t, result, "muted")

	resultNoColor = formatMuted("muted", true)
	assert.Equal(t, "muted", resultNoColor)
}

func TestOutput_AdditionalMethods(t *testing.T) {
	output := NewOutput()
	output.Box("Title", "Content")
	assert.NotNil(t, output)

	output.SubHeader("Test SubHeader")
	assert.NotNil(t, output)

	items := []string{"item1", "item2", "item3"}
	output.List(items)
	assert.NotNil(t, output)

	err := output.Progress("Testing...", func() error {
		return nil
	})
	assert.NoError(t, err)

	assert.NotNil(t, DefaultOutput)
	assert.IsType(t, &Output{}, DefaultOutput)
}
