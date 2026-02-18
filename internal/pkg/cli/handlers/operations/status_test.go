//go:build unit

package operations

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/test/testhelpers"
	"github.com/stretchr/testify/assert"
)

const (
	defaultLogTailLines = core.DefaultLogTailLines
)

func TestNewStatusHandler(t *testing.T) {
	handler := NewStatusHandler()

	assert.NotNil(t, handler)
	assert.IsType(t, &StatusHandler{}, handler)
	assert.NotNil(t, handler.logger, "Logger should be initialized")
}

func TestStatusHandler_ValidateArgs(t *testing.T) {
	handler := NewStatusHandler()

	err := handler.ValidateArgs([]string{})
	assert.NoError(t, err, "Status command should accept no arguments")

	err = handler.ValidateArgs([]string{testhelpers.TestServiceName})
	assert.NoError(t, err, "Status command should accept service names")
}

func TestStatusHandler_GetRequiredFlags(t *testing.T) {
	handler := NewStatusHandler()

	flags := handler.GetRequiredFlags()
	assert.Empty(t, flags, "Status command should have no required flags")
}

func TestStatusHandler_Handle(t *testing.T) {
	handler := NewStatusHandler()
	assert.NotNil(t, handler, "Handler should exist")
}

func TestNewLogsHandler(t *testing.T) {
	handler := NewLogsHandler()

	assert.NotNil(t, handler)
	assert.IsType(t, &LogsHandler{}, handler)
}

func TestLogsHandler_ValidateArgs(t *testing.T) {
	handler := NewLogsHandler()

	err := handler.ValidateArgs([]string{})
	assert.NoError(t, err, "Logs command should accept no arguments")

	err = handler.ValidateArgs([]string{testhelpers.TestServiceName})
	assert.NoError(t, err, "Logs command should accept service names")
}

func TestDefaultLogTailConstant(t *testing.T) {
	assert.Equal(t, "100", defaultLogTailLines,
		"Default log tail lines should match the hardcoded value in commands.go")
}
