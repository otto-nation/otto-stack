//go:build unit

package lifecycle

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/common"
	"github.com/otto-nation/otto-stack/test/testhelpers"
)

func TestNewServiceCommand(t *testing.T) {
	stateManager := &common.StateManager{}
	cmd := NewServiceCommand(core.CommandUp, stateManager)
	testhelpers.AssertValidConstructor(t, cmd, nil, "ServiceCommand")
}

func TestServiceCommand_Operations(t *testing.T) {
	tests := []struct {
		name      string
		operation string
	}{
		{"up operation", core.CommandUp},
		{"down operation", core.CommandDown},
		{"restart operation", core.CommandRestart},
		{"cleanup operation", core.CommandCleanup},
	}

	stateManager := &common.StateManager{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewServiceCommand(tt.operation, stateManager)
			testhelpers.AssertValidConstructor(t, cmd, nil, "ServiceCommand")
		})
	}
}

func TestServiceCommand_ExecuteOperations(t *testing.T) {
	stateManager := &common.StateManager{}

	t.Run("service command with different operations", func(t *testing.T) {
		operations := []string{core.CommandUp, core.CommandDown, core.CommandLogs, core.CommandStatus, core.CommandRestart, core.CommandConnect, core.CommandExec, core.CommandCleanup}
		for _, op := range operations {
			cmd := NewServiceCommand(op, stateManager)
			if cmd.operation != op {
				t.Errorf("Expected operation %s, got %s", op, cmd.operation)
			}
		}
	})
}
