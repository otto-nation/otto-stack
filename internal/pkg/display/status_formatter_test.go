//go:build unit

package display

import (
	"bytes"
	"testing"
	"time"

	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStatusFormatter(t *testing.T) {
	buf := &bytes.Buffer{}
	formatter := NewStatusFormatter(buf)
	assert.NotNil(t, formatter)
}

func TestStatusFormatter_FormatTable_Compact(t *testing.T) {
	buf := &bytes.Buffer{}
	formatter := NewStatusFormatter(buf)

	serviceStatuses := []ServiceStatus{
		{
			Name:      services.ServicePostgres,
			State:     "running",
			Health:    "healthy",
			Ports:     []string{"5432:5432"},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Uptime:    time.Hour,
		},
	}

	err := formatter.FormatTable(serviceStatuses, Options{Compact: true})
	require.NoError(t, err)
	output := buf.String()
	assert.Contains(t, output, services.ServicePostgres)
	assert.Contains(t, output, "running")
}

func TestStatusFormatter_FormatTable_Full(t *testing.T) {
	buf := &bytes.Buffer{}
	formatter := NewStatusFormatter(buf)

	serviceStatuses := []ServiceStatus{
		{
			Name:   services.ServiceRedis,
			State:  "running",
			Health: "healthy",
		},
	}

	err := formatter.FormatTable(serviceStatuses, Options{Compact: false})
	require.NoError(t, err)
	output := buf.String()
	assert.Contains(t, output, services.ServiceRedis)
}

func TestStatusFormatter_FormatResourceSummary(t *testing.T) {
	buf := &bytes.Buffer{}
	formatter := NewStatusFormatter(buf)

	services := []ServiceStatus{
		{Name: "service1", State: "running"},
		{Name: "service2", State: "running"},
		{Name: "service3", State: "stopped"},
	}

	formatter.formatResourceSummary(services)
	output := buf.String()
	assert.Contains(t, output, "Summary: 3 total")
	assert.Contains(t, output, "running")
}

func TestStatusFormatter_CreateSummary(t *testing.T) {
	formatter := NewStatusFormatter(&bytes.Buffer{})

	services := []ServiceStatus{
		{State: "running"},
		{State: "running"},
		{State: "stopped"},
	}

	summary := formatter.createSummary(services)
	assert.Equal(t, 2, summary["running"])
	assert.Equal(t, 1, summary["stopped"])
}

func TestStatusFormatter_getIcon(t *testing.T) {
	sf := NewStatusFormatter(&bytes.Buffer{})

	tests := []struct {
		state    string
		expected string
	}{
		{"running", "✓ "},
		{"healthy", "✓ "},
		{"stopped", "✗ "},
		{"unhealthy", "✗ "},
		{"starting", "! "},
		{"unknown", "— "},
	}

	for _, tt := range tests {
		t.Run(tt.state, func(t *testing.T) {
			result := sf.getIcon(tt.state)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestStatusFormatter_formatDuration(t *testing.T) {
	sf := NewStatusFormatter(&bytes.Buffer{})

	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{"seconds", 30 * time.Second, "30s"},
		{"minutes", 5 * time.Minute, "5m"},
		{"hours", 3 * time.Hour, "3h"},
		{"days", 48 * time.Hour, "2d"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sf.formatDuration(tt.duration)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestStatusFormatter_formatPorts(t *testing.T) {
	sf := NewStatusFormatter(&bytes.Buffer{})

	tests := []struct {
		name     string
		ports    []string
		expected string
	}{
		{"empty", []string{}, "-"},
		{"single", []string{"8080"}, "8080"},
		{"multiple", []string{"8080", "8081"}, "8080,8081"},
		{"many", []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13"}, "1,2,3,4,5,6,7,8,9,10,11,12..."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sf.formatPorts(tt.ports)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestStatusFormatter_colorizeState(t *testing.T) {
	tests := []struct {
		name     string
		state    string
		wantANSI bool
	}{
		{"running applies green", "running", true},
		{"healthy applies green", "healthy", true},
		{"stopped applies red", "stopped", true},
		{"unhealthy applies red", "unhealthy", true},
		{"starting applies yellow", "starting", true},
		{"unknown applies gray", "unknown", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sfColor := &StatusFormatter{writer: &bytes.Buffer{}, noColor: false}
			colored := sfColor.colorizeState("text", tt.state)
			assert.Contains(t, colored, "\033[", "should contain ANSI escape when color enabled")
			assert.Contains(t, colored, "text")

			sfNoColor := &StatusFormatter{writer: &bytes.Buffer{}, noColor: true}
			plain := sfNoColor.colorizeState("text", tt.state)
			assert.Equal(t, "text", plain, "noColor should return text unchanged")
		})
	}
}

func TestBuildServiceContainerMap(t *testing.T) {
	configs := []types.ServiceConfig{
		{Name: "postgres"},
		{Name: "redis", Service: types.ServiceSpec{
			Dependencies: types.ServiceDependencies{Required: []string{"redis-container"}},
		}},
		{Name: "hidden-init", Hidden: true},
	}

	m := buildServiceContainerMap(configs)

	assert.Equal(t, "postgres", m["postgres"])
	assert.Equal(t, "redis-container", m["redis"])
	assert.Equal(t, "hidden-init", m["hidden-init"])
}

func TestConvertToServiceStatuses(t *testing.T) {
	configs := []types.ServiceConfig{
		{Name: "postgres", Shareable: true},
		{Name: "init-job", Container: types.ContainerSpec{Restart: types.RestartPolicyNo}},
		{Name: "hidden-svc", Hidden: true},
	}
	serviceToContainer := map[string]string{
		"postgres":   "postgres",
		"init-job":   "init-job",
		"hidden-svc": "hidden-svc",
	}
	containerStatuses := []docker.ContainerStatus{
		{Name: "postgres", State: "running", Health: "healthy"},
	}

	result := convertToServiceStatuses(containerStatuses, configs, serviceToContainer)

	// init-job (RestartPolicyNo) and hidden-svc (Hidden) should be excluded
	require.Len(t, result, 1)
	assert.Equal(t, "postgres", result[0].Name)
	assert.Equal(t, "running", result[0].State)
	assert.Equal(t, "healthy", result[0].Health)
	assert.Equal(t, ScopeShared, result[0].Scope)
}

func TestConvertToServiceStatuses_NotFound(t *testing.T) {
	configs := []types.ServiceConfig{
		{Name: "postgres"},
	}
	serviceToContainer := map[string]string{"postgres": "postgres"}
	// No matching container status — service should appear as not found
	result := convertToServiceStatuses([]docker.ContainerStatus{}, configs, serviceToContainer)

	require.Len(t, result, 1)
	assert.Equal(t, StateNotFound, result[0].State)
	assert.Equal(t, StateUnknown, result[0].Health)
}

func TestRenderStatusTable(t *testing.T) {
	configs := []types.ServiceConfig{
		{Name: "postgres", Shareable: true},
		{Name: "redis"},
	}
	containerStatuses := []docker.ContainerStatus{
		{Name: "postgres", State: "running", Health: "healthy"},
		{Name: "redis", State: "running", Health: "healthy"},
	}

	buf := &bytes.Buffer{}
	err := RenderStatusTable(buf, containerStatuses, configs, true, true)
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "postgres")
	assert.Contains(t, output, "redis")
	assert.Contains(t, output, "running")
	// noColor=true — no ANSI escape codes
	assert.NotContains(t, output, "\033[")
}

func TestRenderStatusTable_Empty(t *testing.T) {
	buf := &bytes.Buffer{}
	// Empty service configs — should produce a table with headers only (no panic)
	err := RenderStatusTable(buf, []docker.ContainerStatus{}, []types.ServiceConfig{}, true, true)
	require.NoError(t, err)
}
