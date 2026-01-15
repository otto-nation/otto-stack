//go:build unit

package utils

import (
	"context"
	"testing"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/common"
	"github.com/spf13/cobra"
)

func TestExecuteLifecycleCommand(t *testing.T) {
	ctx := context.Background()
	cmd := &cobra.Command{}
	args := []string{}
	base := &base.BaseCommand{}
	stateManager := common.NewStateManager()

	// Test that the function doesn't panic and handles missing config gracefully
	err := ExecuteLifecycleCommand(ctx, cmd, args, base, core.CommandUp, stateManager)

	// We expect an error due to missing config in test environment
	if err == nil {
		t.Error("Expected error due to missing config in test environment")
	}

	// Test with down command as well
	err = ExecuteLifecycleCommand(ctx, cmd, args, base, core.CommandDown, stateManager)
	if err == nil {
		t.Error("Expected error due to missing config in test environment")
	}
}
