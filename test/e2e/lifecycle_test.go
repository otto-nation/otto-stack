//go:build integration

package e2e

import (
	"fmt"
	"testing"
	"time"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/test/e2e/framework"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestE2E_ServiceLifecycle(t *testing.T) {
	projectName := fmt.Sprintf("lifecycle-e2e-%d", time.Now().UnixNano())
	lifecycle := framework.NewTestLifecycle(t, projectName, []string{services.ServicePostgres})
	defer lifecycle.Cleanup()

	err := lifecycle.InitializeStack()
	require.NoError(t, err)

	err = lifecycle.StartServices()
	require.NoError(t, err)

	result := lifecycle.CLI.RunExpectSuccess(core.CommandStatus)
	assert.Contains(t, result.Stdout, services.ServicePostgres)
	assert.Contains(t, result.Stdout, "running")

	err = lifecycle.StopServices()
	require.NoError(t, err)

	err = lifecycle.StartServices()
	require.NoError(t, err)

	result = lifecycle.CLI.RunExpectSuccess(core.CommandStatus)
	assert.Contains(t, result.Stdout, "running")
}
