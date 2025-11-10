package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServiceUtils(t *testing.T) {
	t.Run("creates service utils successfully", func(t *testing.T) {
		utils := NewServiceUtils()
		assert.NotNil(t, utils)
		assert.NotNil(t, utils.manager)
	})
}

func TestServiceUtils_ResolveServices(t *testing.T) {
	utils := NewServiceUtils()

	t.Run("resolves services through manager", func(t *testing.T) {
		resolved, err := utils.ResolveServices([]string{"postgres"})
		if err == nil {
			assert.NotEmpty(t, resolved)
		}
	})

	t.Run("handles empty input", func(t *testing.T) {
		resolved, err := utils.ResolveServices([]string{})
		assert.NoError(t, err)
		assert.Empty(t, resolved)
	})
}

func TestServiceUtils_LoadServicesByCategory(t *testing.T) {
	utils := NewServiceUtils()

	t.Run("loads services organized by category", func(t *testing.T) {
		categories, err := utils.LoadServicesByCategory()
		assert.NoError(t, err)
		assert.NotEmpty(t, categories)

		// Verify structure
		for categoryName, services := range categories {
			assert.NotEmpty(t, categoryName, "Category should have a name")
			assert.NotEmpty(t, services, "Category should have services")

			for _, service := range services {
				assert.False(t, service.Hidden, "Hidden services should be filtered out")
			}
		}
	})
}

func TestServiceUtils_LoadServiceConfig(t *testing.T) {
	utils := NewServiceUtils()

	t.Run("loads specific service config", func(t *testing.T) {
		config, err := utils.LoadServiceConfig("postgres")
		if err == nil {
			assert.NotNil(t, config)
			assert.Equal(t, "postgres", config.Name)
		}
	})

	t.Run("returns error for nonexistent service", func(t *testing.T) {
		config, err := utils.LoadServiceConfig("nonexistent")
		assert.Error(t, err)
		assert.Nil(t, config)
	})
}

func TestServiceUtils_GetServicesByCategory(t *testing.T) {
	utils := NewServiceUtils()

	t.Run("alias for LoadServicesByCategory", func(t *testing.T) {
		categories1, err1 := utils.LoadServicesByCategory()
		categories2, err2 := utils.GetServicesByCategory()

		assert.Equal(t, err1, err2)
		if err1 == nil && err2 == nil {
			assert.Equal(t, len(categories1), len(categories2))
		}
	})
}

func TestServiceUtils_LoadAllServiceDependencies(t *testing.T) {
	utils := NewServiceUtils()

	t.Run("returns empty map (deprecated)", func(t *testing.T) {
		deps, err := utils.LoadAllServiceDependencies()
		assert.NoError(t, err)
		assert.Empty(t, deps)
	})
}

func TestIsYAMLFile(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		expected bool
	}{
		{"yaml extension", "config.yaml", true},
		{"yml extension", "config.yml", true},
		{"no extension", "config", false},
		{"other extension", "config.json", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsYAMLFile(tt.filename)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTrimYAMLExt(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		expected string
	}{
		{"yaml extension", "config.yaml", "config"},
		{"yml extension", "config.yml", "config"},
		{"no extension", "config", "config"},
		{"other extension", "config.json", "config.json"},
		{"empty string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TrimYAMLExt(tt.filename)
			assert.Equal(t, tt.expected, result)
		})
	}
}
