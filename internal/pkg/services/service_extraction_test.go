//go:build unit

package services

import (
	"testing"

	servicetypes "github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/otto-nation/otto-stack/test/fixtures"
	"github.com/stretchr/testify/assert"
)

func TestExtractServiceNames(t *testing.T) {
	serviceConfigs := []servicetypes.ServiceConfig{
		fixtures.NewServiceConfig(ServicePostgres).Build(),
		fixtures.NewServiceConfig(ServiceRedis).Build(),
		fixtures.NewServiceConfig(ServiceMysql).Build(),
	}

	names := ExtractServiceNames(serviceConfigs)
	expected := []string{ServicePostgres, ServiceRedis, ServiceMysql}
	assert.Equal(t, expected, names)

	names = ExtractServiceNames([]servicetypes.ServiceConfig{})
	assert.Nil(t, names)

	names = ExtractServiceNames(nil)
	assert.Nil(t, names)
}

func TestNewServiceUtils(t *testing.T) {
	utils := NewServiceUtils()
	assert.NotNil(t, utils)
	assert.NotNil(t, utils.manager)
}

func TestServiceUtils_LoadServicesByCategory(t *testing.T) {
	utils := NewServiceUtils()

	categories, err := utils.LoadServicesByCategory()
	assert.NoError(t, err)
	assert.NotEmpty(t, categories)

	for categoryName, services := range categories {
		assert.NotEmpty(t, categoryName, "Category should have a name")
		assert.NotEmpty(t, services, "Category should have services")

		for _, service := range services {
			assert.False(t, service.Hidden, "Hidden services should be filtered out")
		}
	}
}

func TestServiceUtils_LoadServiceConfig(t *testing.T) {
	utils := NewServiceUtils()

	config, err := utils.LoadServiceConfig(ServicePostgres)
	if err == nil {
		assert.NotNil(t, config)
		assert.Equal(t, ServicePostgres, config.Name)
	}

	config, err = utils.LoadServiceConfig("nonexistent")
	assert.Error(t, err)
	assert.Nil(t, config)
}

func TestServiceUtils_GetServicesByCategory(t *testing.T) {
	utils := NewServiceUtils()

	categories1, err1 := utils.LoadServicesByCategory()
	categories2, err2 := utils.GetServicesByCategory()

	assert.Equal(t, err1, err2)
	if err1 == nil && err2 == nil {
		assert.Equal(t, len(categories1), len(categories2))
	}
}
