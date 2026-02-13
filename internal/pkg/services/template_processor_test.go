//go:build unit

package services

import (
	"testing"

	servicetypes "github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTemplateProcessor(t *testing.T) {
	processor := NewTemplateProcessor()
	assert.NotNil(t, processor)
}

func TestTemplateProcessor_Process(t *testing.T) {
	processor := NewTemplateProcessor()

	t.Run("processes simple template", func(t *testing.T) {
		config := servicetypes.ServiceConfig{
			Name: "test-service",
		}
		script := "echo 'Hello World'"

		result, err := processor.Process(script, config, []servicetypes.ServiceConfig{})
		require.NoError(t, err)
		assert.Equal(t, "echo 'Hello World'", result)
	})

	t.Run("handles invalid template", func(t *testing.T) {
		config := servicetypes.ServiceConfig{Name: "test"}
		script := "{{.Invalid"

		_, err := processor.Process(script, config, []servicetypes.ServiceConfig{})
		assert.Error(t, err)
	})
}
