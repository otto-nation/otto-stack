//go:build unit

package project

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/test/fixtures"
	"github.com/stretchr/testify/assert"
)

func TestConflictsHandler_ValidateArgs(t *testing.T) {
	handler := &ConflictsHandler{}

	err := handler.ValidateArgs([]string{})
	assert.NoError(t, err)

	err = handler.ValidateArgs([]string{services.ServicePostgres, services.ServiceRedis})
	assert.NoError(t, err)
}

func TestConflictsHandler_GetRequiredFlags(t *testing.T) {
	handler := &ConflictsHandler{}
	flags := handler.GetRequiredFlags()
	assert.Empty(t, flags)
}

func TestConflictsHandler_ParsePort(t *testing.T) {
	handler := &ConflictsHandler{}

	port := handler.parsePort("8080")
	assert.Equal(t, 8080, port)

	port = handler.parsePort("invalid")
	assert.Equal(t, 0, port)

	port = handler.parsePort("")
	assert.Equal(t, 0, port)
}

func TestConflictsHandler_ExtractPortsFromService(t *testing.T) {
	handler := &ConflictsHandler{}

	service := fixtures.NewServiceConfig("test").
		WithPort("8080", "8080").
		WithPort("9090", "9090").
		Build()

	ports := handler.extractPortsFromService(&service)
	assert.Len(t, ports, 2)
	assert.Contains(t, ports, 8080)
	assert.Contains(t, ports, 9090)

	service = fixtures.NewServiceConfig("test").
		WithPort("8080", "8080").
		WithPort("invalid", "invalid").
		Build()

	ports = handler.extractPortsFromService(&service)
	assert.Len(t, ports, 1)
	assert.Contains(t, ports, 8080)

	service = fixtures.NewServiceConfig("test").Build()
	ports = handler.extractPortsFromService(&service)
	assert.Empty(t, ports)
}

func TestDepsHandler_ValidateArgs(t *testing.T) {
	handler := &DepsHandler{}

	err := handler.ValidateArgs([]string{})
	assert.NoError(t, err)

	err = handler.ValidateArgs([]string{services.ServicePostgres})
	assert.NoError(t, err)
}

func TestDepsHandler_GetRequiredFlags(t *testing.T) {
	handler := &DepsHandler{}
	flags := handler.GetRequiredFlags()
	assert.Empty(t, flags)
}
