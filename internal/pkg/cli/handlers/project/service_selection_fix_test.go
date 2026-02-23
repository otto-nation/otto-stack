//go:build unit

package project

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceSelectionManager_ReturnsOriginalSelection(t *testing.T) {
	manager, err := services.New()
	require.NoError(t, err)

	t.Run("resolved includes dependencies but selection does not", func(t *testing.T) {
		// localstack-sns depends on localstack (hidden)
		localstackSns, err := manager.GetService("localstack-sns")
		require.NoError(t, err)
		assert.Contains(t, localstackSns.Service.Dependencies.Required, "localstack")

		// Resolver includes dependencies
		resolver := services.NewServiceResolver(manager)
		resolved, err := resolver.ResolveServices([]string{"localstack-sns"})
		require.NoError(t, err)

		resolvedNames := services.ExtractServiceNames(resolved)
		assert.Contains(t, resolvedNames, "localstack")
		assert.Contains(t, resolvedNames, "localstack-sns")
		assert.Len(t, resolvedNames, 2)

		// SelectionResult should only contain original user selection
		// This ensures config.Stack.Enabled doesn't include hidden dependencies
	})
}
