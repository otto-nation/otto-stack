//go:build unit

package project

import (
	"context"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/ui"
	"github.com/otto-nation/otto-stack/test/testhelpers"
)

func TestNewServicesHandler(t *testing.T) {
	handler := NewServicesHandler()

	assert.NotNil(t, handler)
	assert.IsType(t, &ServicesHandler{}, handler)
}

func TestServicesHandler_ValidateArgs_NoArgs(t *testing.T) {
	handler := NewServicesHandler()
	err := handler.ValidateArgs([]string{})
	assert.NoError(t, err, "Services command should accept no arguments")
}

func TestServicesHandler_ValidateArgs_WithServiceNames(t *testing.T) {
	handler := NewServicesHandler()
	err := handler.ValidateArgs([]string{testhelpers.TestServiceName})
	assert.NoError(t, err, "Services command should accept service names")
}

func TestServicesHandler_GetRequiredFlags(t *testing.T) {
	handler := NewServicesHandler()
	flags := handler.GetRequiredFlags()
	assert.Empty(t, flags, "Services command should have no required flags")
}

func TestServicesHandler_Handle(t *testing.T) {
	handler := NewServicesHandler()
	cmd := &cobra.Command{
		Use: core.CommandServices,
	}

	cmd.Flags().String("format", "table", "Output format")
	cmd.Flags().Bool("quiet", false, "Suppress output")

	base := &base.BaseCommand{
		Output: ui.NewOutput(),
	}

	ctx := context.Background()
	args := []string{}

	err := handler.Handle(ctx, cmd, args, base)
	assert.NoError(t, err, "Handler should succeed in listing available services")
}

// TODO: Add unit tests for service validation logic
// TODO: Add tests for different output formats (table, json, yaml)
// TODO: Add tests for service category filtering
// TODO: Extract common test utilities to reduce duplication
// TODO: Add E2E tests for full service listing workflow
// TODO: Add tests for error handling scenarios
