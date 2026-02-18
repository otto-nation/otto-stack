//go:build unit

package utility

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/otto-nation/otto-stack/test/fixtures"
	"github.com/stretchr/testify/assert"
)

func TestWebInterfacesHandler_extractWebInterfaces(t *testing.T) {
	handler := NewWebInterfacesHandler()

	configs := []types.ServiceConfig{
		fixtures.NewServiceConfig(services.ServicePostgres).Build(),
		fixtures.NewServiceConfig(services.ServiceRedis).Build(),
	}
	configs[0].Documentation.WebInterfaces = []types.WebInterface{
		{Name: "Admin", URL: "http://localhost:5432", Description: "DB Admin"},
	}
	configs[1].Documentation.WebInterfaces = []types.WebInterface{
		{Name: "Commander", URL: "http://localhost:6379", Description: "Redis UI"},
	}

	result := handler.extractWebInterfaces(configs, nil, true)
	assert.Len(t, result, 2)
	assert.Equal(t, services.ServicePostgres, result[0].Service)
	assert.Equal(t, "Admin", result[0].Name)
	assert.Equal(t, services.ServiceRedis, result[1].Service)

	running := map[string]bool{services.ServicePostgres: true, services.ServiceRedis: false}
	result = handler.extractWebInterfaces(configs, running, false)
	assert.Len(t, result, 1)
	assert.Equal(t, services.ServicePostgres, result[0].Service)

	configs = []types.ServiceConfig{
		fixtures.NewServiceConfig(services.ServicePostgres).Build(),
	}
	result = handler.extractWebInterfaces(configs, nil, true)
	assert.Len(t, result, 0)

	configs = []types.ServiceConfig{
		fixtures.NewServiceConfig("grafana").Build(),
	}
	configs[0].Documentation.WebInterfaces = []types.WebInterface{
		{Name: "Dashboard", URL: "http://localhost:3000"},
		{Name: "API", URL: "http://localhost:3000/api"},
	}
	result = handler.extractWebInterfaces(configs, nil, true)
	assert.Len(t, result, 2)
	assert.Equal(t, "Dashboard", result[0].Name)
	assert.Equal(t, "API", result[1].Name)
}

func TestWebInterfacesHandler_createWebInterfaces(t *testing.T) {
	handler := NewWebInterfacesHandler()

	webIfaces := []types.WebInterface{
		{Name: "Admin", URL: "http://localhost:8080", Description: "Admin Panel"},
		{Name: "API", URL: "http://localhost:8080/api", Description: "REST API"},
	}

	result := handler.createWebInterfaces("myservice", webIfaces)
	assert.Len(t, result, 2)
	assert.Equal(t, "myservice", result[0].Service)
	assert.Equal(t, "Admin", result[0].Name)
	assert.Equal(t, "http://localhost:8080", result[0].URL)
	assert.Equal(t, "Admin Panel", result[0].Description)
	assert.Equal(t, "API", result[1].Name)

	result = handler.createWebInterfaces("myservice", []types.WebInterface{})
	assert.Len(t, result, 0)
}

func TestWebInterfacesHandler_shouldIncludeService(t *testing.T) {
	handler := NewWebInterfacesHandler()

	assert.True(t, handler.shouldIncludeService(services.ServicePostgres, nil, true))
	assert.True(t, handler.shouldIncludeService(services.ServiceRedis, map[string]bool{services.ServiceRedis: false}, true))

	assert.True(t, handler.shouldIncludeService(services.ServicePostgres, nil, false))

	running := map[string]bool{services.ServicePostgres: true, services.ServiceRedis: false}
	assert.True(t, handler.shouldIncludeService(services.ServicePostgres, running, false))
	assert.False(t, handler.shouldIncludeService(services.ServiceRedis, running, false))

	running = map[string]bool{services.ServicePostgres: true}
	assert.False(t, handler.shouldIncludeService(services.ServiceRedis, running, false))
}

func TestWebInterfacesHandler_formatStatus(t *testing.T) {
	handler := NewWebInterfacesHandler()

	result := handler.formatStatus(true)
	assert.Contains(t, result, "Available")

	result = handler.formatStatus(false)
	assert.Contains(t, result, "Not Available")
}

func TestWebInterfacesHandler_formatStatusFromResponse(t *testing.T) {
	handler := NewWebInterfacesHandler()

	result := handler.formatStatusFromResponse(200)
	assert.Contains(t, result, "Available")

	result = handler.formatStatusFromResponse(404)
	assert.Contains(t, result, "Not Available")

	result = handler.formatStatusFromResponse(500)
	assert.Contains(t, result, "Not Available")
}

func TestWebInterfacesHandler_ValidateArgs(t *testing.T) {
	handler := NewWebInterfacesHandler()

	err := handler.ValidateArgs([]string{})
	assert.NoError(t, err)

	err = handler.ValidateArgs([]string{services.ServicePostgres, services.ServiceRedis})
	assert.NoError(t, err)
}
