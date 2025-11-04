package types

import (
	"testing"
	"time"

	"github.com/otto-nation/otto-stack/internal/pkg/constants"
	"github.com/stretchr/testify/assert"
)

func TestServiceState(t *testing.T) {
	tests := []struct {
		name      string
		state     ServiceState
		isRunning bool
		isStopped bool
	}{
		{"running", ServiceStateRunning, true, false},
		{"stopped", ServiceStateStopped, false, true},
		{"created", ServiceStateCreated, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.isRunning, tt.state.IsRunning())
			assert.Equal(t, tt.isStopped, tt.state.IsStopped())
			assert.Equal(t, string(tt.state), tt.state.String())
		})
	}
}

func TestHealthStatus(t *testing.T) {
	tests := []struct {
		name        string
		health      HealthStatus
		isHealthy   bool
		isUnhealthy bool
		isStarting  bool
	}{
		{"healthy", HealthStatusHealthy, true, false, false},
		{"unhealthy", HealthStatusUnhealthy, false, true, false},
		{"starting", HealthStatusStarting, false, false, true},
		{"none", HealthStatusNone, false, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.isHealthy, tt.health.IsHealthy())
			assert.Equal(t, tt.isUnhealthy, tt.health.IsUnhealthy())
			assert.Equal(t, tt.isStarting, tt.health.IsStarting())
			assert.Equal(t, string(tt.health), tt.health.String())
		})
	}
}

func TestShellType(t *testing.T) {
	tests := []struct {
		name    string
		shell   ShellType
		isValid bool
	}{
		{"bash", ShellTypeBash, true},
		{"zsh", ShellTypeZsh, true},
		{"fish", ShellTypeFish, true},
		{"powershell", ShellTypePowerShell, true},
		{"invalid", ShellType("invalid"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.isValid, tt.shell.IsValid())
			assert.Equal(t, string(tt.shell), tt.shell.String())
		})
	}
}

func TestProject(t *testing.T) {
	project := Project{
		Name:        "test-project",
		Type:        constants.ProjectTypeDocker,
		Environment: constants.DefaultEnvironment,
		Services:    []string{"postgres", "redis"},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	assert.Equal(t, "test-project", project.Name)
	assert.Equal(t, constants.ProjectTypeDocker, project.Type)
	assert.Equal(t, constants.DefaultEnvironment, project.Environment)
	assert.Len(t, project.Services, 2)
	assert.Contains(t, project.Services, "postgres")
	assert.Contains(t, project.Services, "redis")
}

func TestService(t *testing.T) {
	service := Service{
		Name:  "postgres",
		Type:  "database",
		Image: "postgres:13",
		Ports: []PortMapping{
			{Host: "5432", Container: "5432", Protocol: "tcp"},
		},
		Environment: map[string]string{
			"POSTGRES_DB": "testdb",
		},
		DependsOn: []string{"network"},
	}

	assert.Equal(t, "postgres", service.Name)
	assert.Equal(t, "database", service.Type)
	assert.Equal(t, "postgres:13", service.Image)
	assert.Len(t, service.Ports, 1)
	assert.Equal(t, "5432", service.Ports[0].Host)
	assert.Equal(t, "testdb", service.Environment["POSTGRES_DB"])
	assert.Contains(t, service.DependsOn, "network")
}

func TestHealthCheck(t *testing.T) {
	healthCheck := HealthCheck{
		Test:        []string{"CMD", "pg_isready"},
		Interval:    constants.DefaultStartTimeoutSeconds * time.Second,
		Timeout:     5 * time.Second,
		Retries:     3,
		StartPeriod: 10 * time.Second,
	}

	assert.Equal(t, []string{"CMD", "pg_isready"}, healthCheck.Test)
	assert.Equal(t, constants.DefaultStartTimeoutSeconds*time.Second, healthCheck.Interval)
	assert.Equal(t, 5*time.Second, healthCheck.Timeout)
	assert.Equal(t, 3, healthCheck.Retries)
	assert.Equal(t, 10*time.Second, healthCheck.StartPeriod)
}

func TestServiceStatus(t *testing.T) {
	status := ServiceStatus{
		Name:     "postgres",
		State:    ServiceStateRunning,
		Health:   HealthStatusHealthy,
		Uptime:   time.Hour,
		CPUUsage: 15.5,
		Memory: MemoryUsage{
			Used:  1024 * 1024 * 100, // 100MB
			Limit: 1024 * 1024 * 512, // 512MB
		},
		CreatedAt: time.Now(),
	}

	assert.Equal(t, "postgres", status.Name)
	assert.True(t, status.State.IsRunning())
	assert.True(t, status.Health.IsHealthy())
	assert.Equal(t, time.Hour, status.Uptime)
	assert.Equal(t, 15.5, status.CPUUsage)
	assert.Equal(t, uint64(1024*1024*100), status.Memory.Used)
}
