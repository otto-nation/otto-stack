//go:build unit

package project

import (
	"testing"

	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestProjectManager_filterProjectServices(t *testing.T) {
	pm := &ProjectManager{}

	configs := []types.ServiceConfig{
		{Name: services.ServicePostgres, Shareable: true},
		{Name: services.ServiceRedis, Shareable: true},
		{Name: "app", Shareable: false},
	}

	result := pm.filterProjectServices(configs, nil)
	assert.Len(t, result, 3)

	sharing := &clicontext.SharingSpec{Enabled: false}
	result = pm.filterProjectServices(configs, sharing)
	assert.Len(t, result, 3)

	sharing = &clicontext.SharingSpec{Enabled: true}
	result = pm.filterProjectServices(configs, sharing)
	assert.Len(t, result, 1)
	assert.Equal(t, "app", result[0].Name)

	configs = []types.ServiceConfig{
		{Name: services.ServicePostgres, Shareable: true},
		{Name: services.ServiceRedis, Shareable: true},
	}

	result = pm.filterProjectServices(configs, sharing)
	assert.Len(t, result, 0)

	configs = []types.ServiceConfig{
		{Name: "app1"},
		{Name: "app2"},
	}
	result = pm.filterProjectServices(configs, sharing)
	assert.Len(t, result, 2)
}

func TestProjectManager_formatServicesList(t *testing.T) {
	pm := &ProjectManager{}

	result := pm.formatServicesList([]string{services.ServicePostgres})
	assert.Equal(t, "- postgres\n", result)

	result = pm.formatServicesList([]string{services.ServicePostgres, services.ServiceRedis, services.ServiceMysql})
	assert.Contains(t, result, "- postgres\n")
	assert.Contains(t, result, "- redis\n")
	assert.Contains(t, result, "- mysql\n")

	result = pm.formatServicesList([]string{})
	assert.Equal(t, "", result)
}
