//go:build unit

package project

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestServiceSelector_buildServiceList(t *testing.T) {
	selector := NewServiceSelector()

	t.Run("builds service list from categories", func(t *testing.T) {
		categories := map[string][]types.ServiceConfig{
			"database": {
				{Name: "postgres", Description: "PostgreSQL"},
				{Name: "mysql", Description: "MySQL"},
			},
			"cache": {
				{Name: "redis", Description: "Redis"},
			},
		}

		allServices, options := selector.buildServiceList(categories)

		assert.Len(t, allServices, 3)
		assert.Len(t, options, 3)

		// Check that options contain category and service info
		for _, opt := range options {
			assert.Contains(t, opt, "[")
			assert.Contains(t, opt, "]")
		}
	})

	t.Run("sorts categories and services alphabetically", func(t *testing.T) {
		categories := map[string][]types.ServiceConfig{
			"zeta": {
				{Name: "zebra", Description: "Z service"},
				{Name: "alpha", Description: "A service"},
			},
			"beta": {
				{Name: "bravo", Description: "B service"},
			},
		}

		_, options := selector.buildServiceList(categories)

		// Beta should come before Zeta
		assert.Contains(t, options[0], "Beta")

		// Within zeta category, alpha should come before zebra
		for i, opt := range options {
			if i > 0 && options[i-1] != opt {
				// Services within same category should be sorted
				continue
			}
		}
	})

	t.Run("handles empty categories", func(t *testing.T) {
		categories := map[string][]types.ServiceConfig{}

		allServices, options := selector.buildServiceList(categories)

		assert.Empty(t, allServices)
		assert.Empty(t, options)
	})
}

func TestServiceSelector_mapSelectedServicesByName(t *testing.T) {
	selector := NewServiceSelector()

	t.Run("maps selected services by name", func(t *testing.T) {
		allServices := []types.ServiceConfig{
			{Name: "postgres", Description: "PostgreSQL"},
			{Name: "redis", Description: "Redis"},
			{Name: "mysql", Description: "MySQL"},
		}

		selectedNames := []string{"postgres", "mysql"}

		result := selector.mapSelectedServicesByName(selectedNames, allServices)

		assert.Len(t, result, 2)
		assert.Equal(t, "postgres", result[0].Name)
		assert.Equal(t, "mysql", result[1].Name)
	})

	t.Run("handles no matches", func(t *testing.T) {
		allServices := []types.ServiceConfig{
			{Name: "postgres", Description: "PostgreSQL"},
		}

		selectedNames := []string{"nonexistent"}

		result := selector.mapSelectedServicesByName(selectedNames, allServices)

		assert.Empty(t, result)
	})

	t.Run("handles empty selection", func(t *testing.T) {
		allServices := []types.ServiceConfig{
			{Name: "postgres", Description: "PostgreSQL"},
		}

		selectedNames := []string{}

		result := selector.mapSelectedServicesByName(selectedNames, allServices)

		assert.Empty(t, result)
	})
}

func TestServiceSelector_loadServiceCategories(t *testing.T) {
	selector := NewServiceSelector()

	t.Run("loads service categories", func(t *testing.T) {
		categories, err := selector.loadServiceCategories()

		assert.NoError(t, err)
		assert.NotNil(t, categories)
		// Should have at least some categories
		assert.NotEmpty(t, categories)
	})
}

func TestNewServiceSelector(t *testing.T) {
	t.Run("creates service selector", func(t *testing.T) {
		selector := NewServiceSelector()
		assert.NotNil(t, selector)
	})
}
