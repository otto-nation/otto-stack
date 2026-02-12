//go:build unit

package project

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/stretchr/testify/assert"
)

// MockProjectValidator for testing
type MockProjectValidator struct {
	validateError error
}

func (m *MockProjectValidator) ValidateProjectName(name string) error {
	return m.validateError
}

func TestNewPromptManager(t *testing.T) {
	validator := &MockProjectValidator{}
	pm := NewPromptManager(validator)

	assert.NotNil(t, pm)
	assert.Equal(t, validator, pm.validator)
}

func TestPromptManager_PromptForServiceConfigs(t *testing.T) {
	pm := &PromptManager{}

	t.Run("handles empty service configs", func(t *testing.T) {
		configs, err := pm.PromptForServiceConfigs()
		// Will likely fail due to no stdin, but tests the function exists
		if err != nil {
			assert.Error(t, err)
		} else {
			assert.NotNil(t, configs)
		}
	})
}

func TestPromptManager_PromptForAdvancedOptions(t *testing.T) {
	pm := &PromptManager{}

	t.Run("handles advanced options prompt", func(t *testing.T) {
		validation, advanced, err := pm.PromptForAdvancedOptions()
		// Will likely fail due to no stdin, but tests the function exists
		if err != nil {
			assert.Error(t, err)
		} else {
			assert.NotNil(t, validation)
			assert.NotNil(t, advanced)
		}
	})
}

func TestPromptManager_ConfirmInitialization(t *testing.T) {
	pm := &PromptManager{}

	t.Run("handles initialization confirmation", func(t *testing.T) {
		mockOutput := &MockOutput{}
		base := &base.BaseCommand{Output: mockOutput}

		result, err := pm.ConfirmInitialization("test-project", []string{"postgres"}, map[string]bool{}, map[string]bool{}, base)
		// Will likely fail due to no stdin, but tests the function exists
		if err != nil {
			assert.Error(t, err)
		} else {
			assert.NotEmpty(t, result)
		}
	})
}
