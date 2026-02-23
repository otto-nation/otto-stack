package operations

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/display"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestNewStatusConverter(t *testing.T) {
	converter := NewStatusConverter()
	assert.NotNil(t, converter)
}

func TestStatusConverter_buildContainerMap(t *testing.T) {
	converter := NewStatusConverter()
	statuses := []docker.ContainerStatus{
		{Name: "service1", State: "running"},
		{Name: "service2", State: "exited"},
	}

	containerMap := converter.buildContainerMap(statuses)
	assert.Len(t, containerMap, 2)
	assert.Equal(t, "service1", containerMap["service1"].Name)
	assert.Equal(t, "service2", containerMap["service2"].Name)
}

func TestStatusConverter_shouldSkipService(t *testing.T) {
	converter := NewStatusConverter()

	assert.True(t, converter.shouldSkipService(types.ServiceConfig{Hidden: true}))
	assert.False(t, converter.shouldSkipService(types.ServiceConfig{Hidden: false}))
}

func TestStatusConverter_getProviderName(t *testing.T) {
	converter := NewStatusConverter()

	assert.Equal(t, "", converter.getProviderName("postgres", "postgres"))
	assert.Equal(t, "custom-provider", converter.getProviderName("postgres", "custom-provider"))
}

func TestStatusConverter_buildFoundStatus(t *testing.T) {
	converter := NewStatusConverter()
	containerStatus := docker.ContainerStatus{
		Name:   "postgres",
		State:  "running",
		Health: "healthy",
	}

	status := converter.buildFoundStatus("postgres", "postgres", containerStatus, true)
	assert.Equal(t, "postgres", status.Name)
	assert.Equal(t, display.ScopeShared, status.Scope)
	assert.Equal(t, "running", status.State)
	assert.Equal(t, "healthy", status.Health)
	assert.Equal(t, "postgres", status.Provider)

	statusLocal := converter.buildFoundStatus("postgres", "postgres", containerStatus, false)
	assert.Equal(t, display.ScopeLocal, statusLocal.Scope)
}

func TestStatusConverter_buildNotFoundStatus(t *testing.T) {
	converter := NewStatusConverter()

	status := converter.buildNotFoundStatus("postgres", "postgres", true)
	assert.Equal(t, "postgres", status.Name)
	assert.Equal(t, display.ScopeShared, status.Scope)
	assert.Equal(t, display.StateNotFound, status.State)
	assert.Equal(t, display.StateUnknown, status.Health)
	assert.Equal(t, "postgres", status.Provider)

	statusLocal := converter.buildNotFoundStatus("postgres", "postgres", false)
	assert.Equal(t, display.ScopeLocal, statusLocal.Scope)
}

func TestStatusConverter_createServiceStatus(t *testing.T) {
	converter := NewStatusConverter()
	config := types.ServiceConfig{Name: "postgres"}
	serviceToContainer := map[string]string{"postgres": "postgres"}
	containerMap := map[string]docker.ContainerStatus{
		"postgres": {
			Name:  "postgres",
			State: "running",
		},
	}

	status := converter.createServiceStatus(config, serviceToContainer, containerMap)
	assert.Equal(t, "postgres", status.Name)
	assert.Equal(t, "running", status.State)
}

func TestStatusConverter_createServiceStatus_NotFound(t *testing.T) {
	converter := NewStatusConverter()
	config := types.ServiceConfig{Name: "postgres"}
	serviceToContainer := map[string]string{}
	containerMap := map[string]docker.ContainerStatus{}

	status := converter.createServiceStatus(config, serviceToContainer, containerMap)
	assert.Equal(t, "postgres", status.Name)
	assert.Equal(t, "not found", status.State)
}

func TestStatusConverter_buildDisplayStatuses(t *testing.T) {
	converter := NewStatusConverter()
	configs := []types.ServiceConfig{
		{Name: "postgres", Hidden: false},
		{Name: "redis", Hidden: true},
	}
	serviceToContainer := map[string]string{"postgres": "postgres"}
	containerMap := map[string]docker.ContainerStatus{
		"postgres": {Name: "postgres", State: "running"},
	}

	statuses := converter.buildDisplayStatuses(configs, serviceToContainer, containerMap)
	assert.Len(t, statuses, 1)
	assert.Equal(t, "postgres", statuses[0].Name)
}

func TestStatusConverter_ConvertToDisplayStatuses(t *testing.T) {
	converter := NewStatusConverter()
	containerStatuses := []docker.ContainerStatus{
		{Name: "postgres", State: "running"},
	}
	serviceConfigs := []types.ServiceConfig{
		{Name: "postgres"},
	}
	serviceToContainer := map[string]string{"postgres": "postgres"}

	statuses := converter.ConvertToDisplayStatuses(containerStatuses, serviceConfigs, serviceToContainer)
	assert.Len(t, statuses, 1)
	assert.Equal(t, "postgres", statuses[0].Name)
	assert.Equal(t, "running", statuses[0].State)
}

func TestStatusConverter_buildFoundStatus_ZeroStartTime(t *testing.T) {
	converter := NewStatusConverter()

	status := converter.buildFoundStatus("test", "docker", docker.ContainerStatus{
		State:  "running",
		Health: "healthy",
		// StartedAt is zero
	}, false)

	if status.Name != "test" {
		t.Errorf("expected name 'test', got '%s'", status.Name)
	}
	if status.Uptime != 0 {
		t.Errorf("expected uptime 0, got %v", status.Uptime)
	}
}
