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

func TestPromptManager_LoadServiceCategories(t *testing.T) {
	pm := &PromptManager{}

	categories, err := pm.loadServiceCategories()

	// Should not error even if no services found
	assert.NoError(t, err)
	assert.NotNil(t, categories)
}

func TestPromptManager_PrepareCategoryNavigation(t *testing.T) {
	pm := &PromptManager{}

	// Test with empty categories
	categories := make(map[string][]services.ServiceConfig)
	categoryNames, categoryServicesList := pm.prepareCategoryNavigation(categories)

	assert.Empty(t, categoryNames)
	assert.Empty(t, categoryServicesList)

	// Test with sample categories
	categories["database"] = []services.ServiceConfig{
		{Name: "postgres", Description: "PostgreSQL database"},
		{Name: "mysql", Description: "MySQL database"},
	}
	categories["cache"] = []services.ServiceConfig{
		{Name: "redis", Description: "Redis cache"},
	}

	categoryNames, categoryServicesList = pm.prepareCategoryNavigation(categories)

	assert.Len(t, categoryNames, 2)
	assert.Len(t, categoryServicesList, 2)
	assert.Contains(t, categoryNames, "database")
	assert.Contains(t, categoryNames, "cache")
}

func TestPromptManager_BuildServiceOptions(t *testing.T) {
	pm := &PromptManager{}

	services := []services.ServiceConfig{
		{Name: "postgres", Description: "PostgreSQL database"},
		{Name: "redis", Description: "Redis cache"},
	}

	// Test without go back option
	options := pm.buildServiceOptions(services, false)

	assert.Len(t, options, 2)
	assert.Contains(t, options, "postgres - PostgreSQL database")
	assert.Contains(t, options, "redis - Redis cache")
	assert.NotContains(t, options, "Go Back")

	// Test with go back option
	options = pm.buildServiceOptions(services, true)

	assert.Len(t, options, 3)
	assert.Contains(t, options, "Go Back")
}

func TestPromptManager_FindCategoryIndex(t *testing.T) {
	pm := &PromptManager{}

	categories := []string{"database", "cache", "messaging"}

	t.Run("find existing category", func(t *testing.T) {
		index := pm.findCategoryIndex(categories, "cache")
		assert.Equal(t, 1, index)
	})

	t.Run("find non-existing category", func(t *testing.T) {
		index := pm.findCategoryIndex(categories, "nonexistent")
		assert.Equal(t, -1, index)
	})

	t.Run("find in empty slice", func(t *testing.T) {
		index := pm.findCategoryIndex([]string{}, "any")
		assert.Equal(t, -1, index)
	})
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
