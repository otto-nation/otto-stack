package framework

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/test/testutil"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

type TestLifecycle struct {
	Environment *E2EEnvironment
	CLI         *CLIRunner
	ProjectName string
	Services    []string
}

type E2EEnvironment struct {
	WorkDir      string
	PortManager  *testutil.PortManager
	ServicePorts map[string]int
	t            *testing.T
}

func NewTestLifecycle(t *testing.T, projectName string, serviceList []string) *TestLifecycle {
	tempDir := t.TempDir()
	workDir := filepath.Join(tempDir, "workspace")

	err := os.MkdirAll(workDir, core.PermReadWriteExec)
	require.NoError(t, err)

	env := &E2EEnvironment{
		WorkDir:      workDir,
		PortManager:  testutil.NewPortManager(),
		ServicePorts: make(map[string]int),
		t:            t,
	}

	// Resolve service dependencies and allocate ports for all services
	allServices := resolveServiceDependencies(serviceList)
	ports, err := env.PortManager.AllocateServicePorts(allServices)
	require.NoError(t, err)
	env.ServicePorts = ports

	// Create project directory but NOT .otto-stack (let init create it)
	projectDir := filepath.Join(workDir, projectName)
	err = os.MkdirAll(projectDir, core.PermReadWriteExec)
	require.NoError(t, err)

	// Build binary
	binPath := filepath.Join(tempDir, core.AppName)
	builder := NewBinaryBuilder(t)
	binPath = builder.Build(binPath)

	cli := NewCLIRunner(t, binPath, projectDir)
	cli.SetEnv(core.EnvOttoNonInteractive, "true")

	// Set port environment variables for dynamic port allocation
	portMappings := map[string]string{
		services.ServicePostgres:   services.EnvKeyPOSTGRES_PORT,
		services.ServiceRedis:      services.EnvKeyREDIS_PORT,
		services.ServiceLocalstack: services.EnvKeyLOCALSTACK_PORT,
	}

	for serviceName, port := range env.ServicePorts {
		if envVar, exists := portMappings[serviceName]; exists {
			cli.SetEnv(envVar, fmt.Sprintf("%d", port))
		}
	}

	tl := &TestLifecycle{
		Environment: env,
		CLI:         cli,
		ProjectName: projectName,
		Services:    serviceList,
	}

	// Register emergency cleanup that runs even on panic/timeout
	t.Cleanup(func() {
		tl.emergencyCleanup()
	})

	return tl
}

func (e *E2EEnvironment) GetServicePort(serviceName string) int {
	if port, exists := e.ServicePorts[serviceName]; exists {
		return port
	}
	// Fallback to default ports
	switch serviceName {
	case services.ServicePostgres:
		return services.PortPostgres
	case services.ServiceRedis:
		return services.PortRedis
	case services.ServiceLocalstack:
		return services.PortLocalstack
	default:
		return 0
	}
}

func (tl *TestLifecycle) InitializeStack() error {
	result := tl.CLI.RunExpectSuccessWithBuilder(func() []string {
		return NewCLIBuilder(core.CommandInit).
			BoolFlag(core.FlagNonInteractive).
			Flag(core.FlagProjectName, tl.ProjectName).
			Flag(core.FlagServices, joinServices(tl.Services)).
			BuildArgs()
	})
	if result.Error != nil {
		return result.Error
	}

	// After initialization, change working directory to the project directory
	// since otto-stack init creates a project directory
	projectDir := filepath.Join(tl.Environment.WorkDir, tl.ProjectName)
	if _, err := os.Stat(projectDir); err == nil {
		tl.CLI.SetWorkDir(projectDir)
		tl.Environment.WorkDir = projectDir
	}

	return nil
}

func (tl *TestLifecycle) AddServiceConfig(serviceName string, config map[string]any) error {
	configPath := filepath.Join(tl.Environment.WorkDir, core.OttoStackDir, core.ConfigFileName)

	// Read existing config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return pkgerrors.NewServiceError("test", "read config file", err)
	}

	var configData map[string]any
	if err := yaml.Unmarshal(data, &configData); err != nil {
		return pkgerrors.NewServiceError("test", "unmarshal config", err)
	}

	// Add service configuration
	if configData[docker.ComposeFieldServices] == nil {
		configData[docker.ComposeFieldServices] = make(map[string]any)
	}

	services := configData[docker.ComposeFieldServices].(map[string]any)
	services[serviceName] = config

	// Write back to file
	updatedData, err := yaml.Marshal(configData)
	if err != nil {
		return pkgerrors.NewServiceError("test", "marshal config", err)
	}

	if err := os.WriteFile(configPath, updatedData, core.PermReadWrite); err != nil {
		return pkgerrors.NewServiceError("test", "write config file", err)
	}

	return nil
}

// CreateServiceConfigFile creates a service config file for init containers
func (tl *TestLifecycle) CreateServiceConfigFile(filename string, config map[string]any) error {
	serviceConfigsDir := filepath.Join(tl.Environment.WorkDir, core.OttoStackDir, core.ServiceConfigsDir)

	// Create service-configs directory if it doesn't exist
	if err := os.MkdirAll(serviceConfigsDir, core.PermReadWrite); err != nil {
		return pkgerrors.NewServiceError("test", "create directory", err)
	}

	// Write service config file
	configPath := filepath.Join(serviceConfigsDir, filename)
	configData, err := yaml.Marshal(config)
	if err != nil {
		return pkgerrors.NewServiceError("test", "marshal service config", err)
	}

	if err := os.WriteFile(configPath, configData, core.PermReadWrite); err != nil {
		return pkgerrors.NewServiceError("test", "write service config file", err)
	}

	return nil
}

func (tl *TestLifecycle) StartServices() error {
	result := tl.CLI.RunWithBuilder(func() []string {
		return NewCLIBuilder(core.CommandUp).
			BuildArgs()
	})
	if result.Error != nil {
		tl.Environment.t.Logf("StartServices failed - stdout: %s, stderr: %s", result.Stdout, result.Stderr)
	}
	return result.Error
}

func (tl *TestLifecycle) StopServices() error {
	result := tl.CLI.RunWithBuilder(func() []string {
		return NewCLIBuilder(core.CommandDown).BuildArgs()
	})
	if result.Error != nil {
		tl.Environment.t.Logf("StopServices failed - stdout: %s, stderr: %s", result.Stdout, result.Stderr)
	}
	return result.Error
}

func (tl *TestLifecycle) WaitForServices() error {
	// Quick health check instead of long sleep
	const maxRetries = 30
	for range maxRetries {
		result := tl.CLI.RunWithBuilder(func() []string {
			return NewCLIBuilder(core.CommandStatus).BuildArgs()
		})
		if result.Error == nil {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("services failed to become ready within 30 seconds")
}

func (tl *TestLifecycle) Cleanup() {
	// Ensure cleanup happens regardless of failures
	defer func() {
		if r := recover(); r != nil {
			tl.Environment.t.Logf("Panic during cleanup: %v", r)
		}
		// Always release allocated ports
		for _, port := range tl.Environment.ServicePorts {
			tl.Environment.PortManager.ReleasePort(port)
		}
	}()

	tl.emergencyCleanup()
}

func (tl *TestLifecycle) emergencyCleanup() {
	// Force stop all containers first
	tl.CLI.RunWithBuilder(func() []string {
		return NewCLIBuilder(core.CommandDown).
			BoolFlag(core.FlagForce).
			BuildArgs()
	})

	// Clean up any remaining resources
	tl.CLI.RunWithBuilder(func() []string {
		return NewCLIBuilder(core.CommandCleanup).
			BoolFlag(core.FlagForce).
			BuildArgs()
	})

	// Additional Docker SDK cleanup for any remaining containers/networks
	logger := slog.Default()
	if dockerClient, err := docker.NewClient(logger); err == nil {
		defer func() { _ = dockerClient.Close() }()
		ctx := context.Background()

		// Clean up containers with our project name
		if containers, err := dockerClient.ListProjectContainers(ctx, tl.ProjectName); err == nil {
			for _, container := range containers {
				_ = dockerClient.RemoveContainer(ctx, container.ID, true) // force remove
			}
		}

		// Clean up networks and volumes for this project
		_ = dockerClient.RemoveResources(ctx, docker.ResourceNetwork, tl.ProjectName)
		_ = dockerClient.RemoveResources(ctx, docker.ResourceVolume, tl.ProjectName)
	}
}

func resolveServiceDependencies(serviceList []string) []string {
	// For E2E tests, we need to include dependency services for port allocation
	allServices := make(map[string]bool)

	// Add requested services
	for _, service := range serviceList {
		allServices[service] = true
	}

	// Add known dependencies
	for _, service := range serviceList {
		switch service {
		case services.ServiceLocalstackSqs, services.ServiceLocalstackSns,
			services.ServiceLocalstackS3, services.ServiceLocalstackDynamodb:
			allServices[services.ServiceLocalstack] = true
		}
	}

	// Convert back to slice
	result := make([]string, 0, len(allServices))
	for service := range allServices {
		result = append(result, service)
	}

	return result
}

func joinServices(services []string) string {
	return strings.Join(services, ",")
}
