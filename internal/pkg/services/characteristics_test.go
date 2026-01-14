//go:build unit

package services

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/core/docker"
	servicetypes "github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/otto-nation/otto-stack/test/testhelpers"
	"github.com/stretchr/testify/assert"
)

func TestNewDefaultCharacteristicsResolver(t *testing.T) {
	resolver, err := NewDefaultCharacteristicsResolver()
	testhelpers.AssertValidConstructor(t, resolver, err, "DefaultCharacteristicsResolver")
	assert.IsType(t, &DefaultCharacteristicsResolver{}, resolver)
}

func TestDefaultCharacteristicsResolver_ResolveUpOptions(t *testing.T) {
	resolver := &DefaultCharacteristicsResolver{}

	serviceConfigs := []servicetypes.ServiceConfig{
		{Name: ServicePostgres},
		{Name: ServiceRedis},
	}

	baseOptions := docker.UpOptions{}
	characteristics := []string{} // Empty characteristics to avoid nil pointer

	// Test that the function doesn't panic and extracts service names
	result := resolver.ResolveUpOptions(characteristics, serviceConfigs, baseOptions)

	// Should extract service names using constants
	expected := []string{ServicePostgres, ServiceRedis}
	assert.Equal(t, expected, result.Services)
}

func TestDefaultCharacteristicsResolver_ResolveDownOptions(t *testing.T) {
	resolver := &DefaultCharacteristicsResolver{}

	serviceConfigs := []servicetypes.ServiceConfig{
		{Name: ServicePostgres},
	}

	baseOptions := docker.DownOptions{}
	characteristics := []string{} // Empty characteristics to avoid nil pointer

	result := resolver.ResolveDownOptions(characteristics, serviceConfigs, baseOptions)

	// Should extract service names using constants
	expected := []string{ServicePostgres}
	assert.Equal(t, expected, result.Services)
}

func TestDefaultCharacteristicsResolver_ResolveStopOptions(t *testing.T) {
	resolver := &DefaultCharacteristicsResolver{}

	serviceConfigs := []servicetypes.ServiceConfig{
		{Name: ServiceRedis},
	}

	baseOptions := docker.StopOptions{}
	characteristics := []string{} // Empty characteristics to avoid nil pointer

	result := resolver.ResolveStopOptions(characteristics, serviceConfigs, baseOptions)

	// Should extract service names using constants
	expected := []string{ServiceRedis}
	assert.Equal(t, expected, result.Services)
}

func TestServiceConstantsValidation(t *testing.T) {
	t.Run("extracts names from service configs using constants", func(t *testing.T) {
		serviceConfigs := []servicetypes.ServiceConfig{
			{Name: ServicePostgres},
			{Name: ServiceRedis},
			{Name: ServiceMysql},
		}

		names := ExtractServiceNames(serviceConfigs)

		expected := []string{ServicePostgres, ServiceRedis, ServiceMysql}
		assert.Equal(t, expected, names)
	})

	t.Run("handles empty service configs", func(t *testing.T) {
		names := ExtractServiceNames([]servicetypes.ServiceConfig{})
		assert.Empty(t, names)
	})

	t.Run("validates service constants are not empty", func(t *testing.T) {
		constants := []string{
			ServicePostgres,
			ServiceRedis,
			ServiceMysql,
			ServiceLocalstack,
		}

		for _, constant := range constants {
			assert.NotEmpty(t, constant, "Service constant should not be empty")
		}
	})
}

func TestServiceConstants(t *testing.T) {
	t.Run("validates all service constants", func(t *testing.T) {
		services := map[string]string{
			"postgres":   ServicePostgres,
			"redis":      ServiceRedis,
			"mysql":      ServiceMysql,
			"localstack": ServiceLocalstack,
		}

		for expected, actual := range services {
			assert.Equal(t, expected, actual, "Service constant mismatch")
		}
	})

	t.Run("validates category constants", func(t *testing.T) {
		categories := []string{
			CategoryDatabase,
			CategoryCache,
			CategoryCloud,
		}

		for _, category := range categories {
			assert.NotEmpty(t, category, "Category constant should not be empty")
		}
	})
}

func TestServiceCatalogConstants(t *testing.T) {
	t.Run("validates catalog format constants", func(t *testing.T) {
		formats := map[string]string{
			"json":  ServiceCatalogJSONFormat,
			"yaml":  ServiceCatalogYAMLFormat,
			"table": ServiceCatalogTableFormat,
		}

		for expected, actual := range formats {
			assert.Equal(t, expected, actual, "Catalog format constant mismatch")
		}
	})

	t.Run("validates catalog message constants", func(t *testing.T) {
		messages := []string{
			MsgServiceCatalogHeader,
			MsgServiceCount,
			MsgCategoryServiceCount,
			SummaryTotal,
			SummaryRunning,
			SummaryHealthy,
		}

		for _, msg := range messages {
			assert.NotEmpty(t, msg, "Message constant should not be empty")
		}
	})
}
