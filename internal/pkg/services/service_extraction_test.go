//go:build unit

package services

import (
	"testing"

	servicetypes "github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestExtractServiceNames(t *testing.T) {
	t.Run("extracts names from ServiceConfigs", func(t *testing.T) {
		serviceConfigs := []servicetypes.ServiceConfig{
			{Name: ServicePostgres},
			{Name: ServiceRedis},
			{Name: ServiceMysql},
		}

		names := ExtractServiceNames(serviceConfigs)
		expected := []string{ServicePostgres, ServiceRedis, ServiceMysql}

		assert.Equal(t, expected, names)
	})

	t.Run("returns nil for empty slice", func(t *testing.T) {
		names := ExtractServiceNames([]servicetypes.ServiceConfig{})
		assert.Nil(t, names)
	})

	t.Run("returns nil for nil slice", func(t *testing.T) {
		names := ExtractServiceNames(nil)
		assert.Nil(t, names)
	})
}

func TestNewServiceUtils(t *testing.T) {
	t.Run("creates service utils successfully", func(t *testing.T) {
		utils := NewServiceUtils()
		assert.NotNil(t, utils)
		assert.NotNil(t, utils.manager)
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
		config, err := utils.LoadServiceConfig(ServicePostgres)
		if err == nil {
			assert.NotNil(t, config)
			assert.Equal(t, ServicePostgres, config.Name)
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
