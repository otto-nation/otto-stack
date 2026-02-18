//go:build unit

package project

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
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
	configs, err := pm.PromptForServiceConfigs()
	if err != nil {
		assert.Error(t, err)
	} else {
		assert.NotNil(t, configs)
	}
}

func TestPromptManager_PromptForAdvancedOptions(t *testing.T) {
	pm := &PromptManager{}
	validation, advanced, err := pm.PromptForAdvancedOptions()
	if err != nil {
		assert.Error(t, err)
	} else {
		assert.NotNil(t, validation)
		assert.NotNil(t, advanced)
	}
}

func TestPromptManager_ConfirmInitialization(t *testing.T) {
	pm := &PromptManager{}
	mockOutput := &MockOutput{}
	base := &base.BaseCommand{Output: mockOutput}

	result, err := pm.ConfirmInitialization("test-project", []string{services.ServicePostgres}, map[string]bool{}, map[string]bool{}, base)
	if err != nil {
		assert.Error(t, err)
	} else {
		assert.NotEmpty(t, result)
	}
}
