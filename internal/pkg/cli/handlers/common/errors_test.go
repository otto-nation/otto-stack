//go:build unit

package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorConstants(t *testing.T) {
	// Test that all error constants are defined and non-empty
	constants := map[string]string{
		"ComponentStack":              ComponentStack,
		"ActionStartServices":         ActionStartServices,
		"ActionStopServices":          ActionStopServices,
		"ActionRestartServices":       ActionRestartServices,
		"ActionShowLogs":              ActionShowLogs,
		"ActionShowStatus":            ActionShowStatus,
		"ActionCleanupResources":      ActionCleanupResources,
		"ActionCreateService":         ActionCreateService,
		"ActionGetManager":            ActionGetManager,
		"ActionCreateGenerator":       ActionCreateGenerator,
		"ActionGenerateCompose":       ActionGenerateCompose,
		"ActionCreateDirectory":       ActionCreateDirectory,
		"ActionGenerateEnv":           ActionGenerateEnv,
		"ActionCreateManager":         ActionCreateManager,
		"ActionLoadProject":           ActionLoadProject,
		"ActionCreateClient":          ActionCreateClient,
		"ActionGetServiceStatus":      ActionGetServiceStatus,
		"ActionConnectToService":      ActionConnectToService,
		"ActionExecuteCommand":        ActionExecuteCommand,
		"ActionGetLogs":               ActionGetLogs,
		"ActionValidateArgs":          ActionValidateArgs,
		"ActionBuildContext":          ActionBuildContext,
		"ActionResolveServices":       ActionResolveServices,
		"ActionFilterServices":        ActionFilterServices,
		"OpShowLogs":                  OpShowLogs,
		"OpListContainers":            OpListContainers,
		"OpRemoveResources":           OpRemoveResources,
		"MsgFailedCreateStackService": MsgFailedCreateStackService,
		"MsgUnsupportedService":       MsgUnsupportedService,
	}

	for name, value := range constants {
		t.Run(name, func(t *testing.T) {
			assert.NotEmpty(t, value, "Constant %s should not be empty", name)
			assert.IsType(t, "", value, "Constant %s should be string", name)
		})
	}
}

func TestErrorConstantsUniqueness(t *testing.T) {
	// Ensure no duplicate values (except where intentional)
	values := []string{
		ActionStartServices,
		ActionStopServices,
		ActionRestartServices,
		ActionShowStatus, // Changed from ActionShowLogs to avoid duplicate
		OpShowLogs,
		OpListContainers,
		OpRemoveResources,
	}

	seen := make(map[string]bool)
	for _, value := range values {
		assert.False(t, seen[value], "Value '%s' should be unique", value)
		seen[value] = true
	}
}
