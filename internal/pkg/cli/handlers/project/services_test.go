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

// TestNewServicesHandler tests the services handler constructor
func TestNewServicesHandler(t *testing.T) {
	t.Run("creates handler successfully", func(t *testing.T) {
		handler := NewServicesHandler()

		assert.NotNil(t, handler)
		assert.IsType(t, &ServicesHandler{}, handler)
	})
}

// TestServicesHandler_ValidateArgs tests argument validation
func TestServicesHandler_ValidateArgs(t *testing.T) {
	handler := NewServicesHandler()

	t.Run("accepts no arguments", func(t *testing.T) {
		err := handler.ValidateArgs([]string{})
		assert.NoError(t, err, "Services command should accept no arguments")
	})

	t.Run("accepts service names as arguments", func(t *testing.T) {
		err := handler.ValidateArgs([]string{testhelpers.TestServiceName})
		assert.NoError(t, err, "Services command should accept service names")
	})
}

// TestServicesHandler_GetRequiredFlags tests required flags
func TestServicesHandler_GetRequiredFlags(t *testing.T) {
	handler := NewServicesHandler()

	t.Run("returns empty slice for required flags", func(t *testing.T) {
		flags := handler.GetRequiredFlags()
		assert.Empty(t, flags, "Services command should have no required flags")
	})
}

// TestServicesHandler_Handle tests the main handler execution
func TestServicesHandler_Handle(t *testing.T) {
	handler := NewServicesHandler()

	t.Run("handles basic execution flow", func(t *testing.T) {
		cmd := &cobra.Command{
			Use: core.CommandServices,
		}

		// Add flags that services command might expect
		cmd.Flags().String("format", "table", "Output format")
		cmd.Flags().Bool("quiet", false, "Suppress output")

		base := &base.BaseCommand{
			Output: ui.NewOutput(),
		}

		ctx := context.Background()
		args := []string{}

		// This handler doesn't use ci.HandleError, so it should be testable
		err := handler.Handle(ctx, cmd, args, base)

		// The handler should succeed in listing services
		assert.NoError(t, err, "Handler should succeed in listing available services")
	})
}

// TestNewServicesCommand tests the services command constructor
func TestNewServicesCommand(t *testing.T) {
	t.Run("creates command successfully", func(t *testing.T) {
		command := NewServicesCommand()

		assert.NotNil(t, command)
		assert.IsType(t, &ServicesCommand{}, command)
	})
}

// TODO: CRITICAL ARCHITECTURE INCONSISTENCY IDENTIFIED:
// Project handlers use individual command structs (ServicesCommand, DepsCommand, etc.)
// while lifecycle/operations use consolidated ServiceCommand pattern.
// This creates maintenance burden and inconsistent patterns.
// RECOMMENDATION: Consolidate project commands to use same pattern as lifecycle/operations

// TODO: Add unit tests for service validation logic
// TODO: Add tests for different output formats (table, json, yaml)
// TODO: Add tests for service category filtering
// TODO: Extract common test utilities to reduce duplication
// TODO: Add E2E tests for full service listing workflow
// TODO: Add tests for error handling scenarios
