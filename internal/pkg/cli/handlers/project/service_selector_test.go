//go:build unit

package project

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/otto-nation/otto-stack/test/fixtures"
	"github.com/stretchr/testify/assert"
)

func TestServiceSelector_buildServiceList(t *testing.T) {
	selector := NewServiceSelector()

	categories := map[string][]types.ServiceConfig{
		"database": {
			fixtures.NewServiceConfig(services.ServicePostgres).Build(),
			fixtures.NewServiceConfig(services.ServiceMysql).Build(),
		},
		"cache": {
			fixtures.NewServiceConfig(services.ServiceRedis).Build(),
		},
	}
	categories["database"][0].Description = "PostgreSQL"
	categories["database"][1].Description = "MySQL"
	categories["cache"][0].Description = "Redis"

	allServices, options := selector.buildServiceList(categories)

	assert.Len(t, allServices, 3)
	assert.Len(t, options, 3)

	for _, opt := range options {
		assert.Contains(t, opt, "[")
		assert.Contains(t, opt, "]")
	}

	categories = map[string][]types.ServiceConfig{
		"zeta": {
			fixtures.NewServiceConfig("zebra").Build(),
			fixtures.NewServiceConfig("alpha").Build(),
		},
		"beta": {
			fixtures.NewServiceConfig("bravo").Build(),
		},
	}
	categories["zeta"][0].Description = "Z service"
	categories["zeta"][1].Description = "A service"
	categories["beta"][0].Description = "B service"

	_, options = selector.buildServiceList(categories)
	assert.Contains(t, options[0], "Beta")

	categories = map[string][]types.ServiceConfig{}
	allServices, options = selector.buildServiceList(categories)
	assert.Empty(t, allServices)
	assert.Empty(t, options)
}

func TestServiceSelector_mapSelectedServicesByName(t *testing.T) {
	selector := NewServiceSelector()

	allServices := []types.ServiceConfig{
		fixtures.NewServiceConfig(services.ServicePostgres).Build(),
		fixtures.NewServiceConfig(services.ServiceRedis).Build(),
		fixtures.NewServiceConfig(services.ServiceMysql).Build(),
	}
	allServices[0].Description = "PostgreSQL"
	allServices[1].Description = "Redis"
	allServices[2].Description = "MySQL"

	selectedNames := []string{services.ServicePostgres, services.ServiceMysql}

	result := selector.mapSelectedServicesByName(selectedNames, allServices)

	assert.Len(t, result, 2)
	assert.Equal(t, services.ServicePostgres, result[0].Name)
	assert.Equal(t, services.ServiceMysql, result[1].Name)

	allServices = []types.ServiceConfig{
		fixtures.NewServiceConfig(services.ServicePostgres).Build(),
	}
	allServices[0].Description = "PostgreSQL"

	selectedNames = []string{"nonexistent"}
	result = selector.mapSelectedServicesByName(selectedNames, allServices)
	assert.Empty(t, result)

	selectedNames = []string{}
	result = selector.mapSelectedServicesByName(selectedNames, allServices)
	assert.Empty(t, result)
}

func TestServiceSelector_loadServiceCategories(t *testing.T) {
	selector := NewServiceSelector()

	categories, err := selector.loadServiceCategories()

	assert.NoError(t, err)
	assert.NotNil(t, categories)
	assert.NotEmpty(t, categories)
}

func TestNewServiceSelector(t *testing.T) {
	selector := NewServiceSelector()
	assert.NotNil(t, selector)
}
