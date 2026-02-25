//go:build integration

package e2e

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/test/e2e/framework"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestE2E_ErrorMessages_NoDoubleWrapping(t *testing.T) {
	projectName := fmt.Sprintf("error-e2e-%d", time.Now().UnixNano())
	lifecycle := framework.NewTestLifecycle(t, projectName, []string{services.ServicePostgres})
	defer lifecycle.Cleanup()

	err := lifecycle.InitializeStack()
	require.NoError(t, err)

	t.Run("status with invalid service shows clean error", func(t *testing.T) {
		result := lifecycle.CLI.RunExpectFailure(core.CommandStatus, "invalid-service-xyz")

		// Should contain the service name
		assert.Contains(t, result.Stderr, "invalid-service-xyz")

		// Should NOT contain format placeholders or nested wrapping
		assert.NotContains(t, result.Stderr, "%v", "Error should not contain format placeholders")
		assert.NotContains(t, result.Stderr, "%s", "Error should not contain format placeholders")

		// Check for nested error patterns
		errorLines := strings.Split(result.Stderr, "\n")
		for _, line := range errorLines {
			if strings.Contains(line, "Failed to resolve") {
				// Count occurrences of "Failed to resolve"
				count := strings.Count(line, "Failed to resolve")
				assert.Equal(t, 1, count, "Error message should not have nested 'Failed to resolve': %s", line)
			}
		}
	})

	t.Run("up with invalid service shows clean error", func(t *testing.T) {
		result := lifecycle.CLI.RunExpectFailure(core.CommandUp, "nonexistent-service")

		assert.Contains(t, result.Stderr, "nonexistent-service")
		assert.NotContains(t, result.Stderr, "%v")
		assert.NotContains(t, result.Stderr, "%s")
	})

	t.Run("down with invalid service shows clean error", func(t *testing.T) {
		result := lifecycle.CLI.RunExpectFailure(core.CommandDown, "invalid-xyz")

		assert.Contains(t, result.Stderr, "invalid-xyz")
		assert.NotContains(t, result.Stderr, "%v")
	})
}

func TestE2E_ErrorMessages_MatchMessagesYAML(t *testing.T) {
	projectName := fmt.Sprintf("msg-e2e-%d", time.Now().UnixNano())
	lifecycle := framework.NewTestLifecycle(t, projectName, []string{services.ServicePostgres})
	defer lifecycle.Cleanup()

	err := lifecycle.InitializeStack()
	require.NoError(t, err)

	t.Run("unknown service error is user-friendly", func(t *testing.T) {
		result := lifecycle.CLI.RunExpectFailure(core.CommandUp, "unknown-service")

		// Should have a clear, single error message
		errorOutput := result.Stderr

		// Should mention the service name
		assert.Contains(t, errorOutput, "unknown-service")

		// Should not have technical stack traces or multiple error layers
		lines := strings.Split(errorOutput, "\n")
		errorCount := 0
		for _, line := range lines {
			if strings.Contains(line, "✗") || strings.Contains(line, "Error:") {
				errorCount++
			}
		}
		assert.LessOrEqual(t, errorCount, 2, "Should have at most 2 error indicators, got %d", errorCount)
	})
}
