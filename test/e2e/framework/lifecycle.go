package framework

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/test/testutil"
	"github.com/stretchr/testify/require"
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

	// Allocate ports
	ports, err := env.PortManager.AllocateServicePorts(serviceList)
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

	return &TestLifecycle{
		Environment: env,
		CLI:         cli,
		ProjectName: projectName,
		Services:    serviceList,
	}
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
	return result.Error
}

func (tl *TestLifecycle) StartServices() error {
	result := tl.CLI.RunWithBuilder(func() []string {
		return NewCLIBuilder(core.CommandUp).
			BoolFlag(core.FlagForceRecreate).
			BuildArgs()
	})
	return result.Error
}

func (tl *TestLifecycle) StopServices() error {
	result := tl.CLI.RunWithBuilder(func() []string {
		return NewCLIBuilder(core.CommandDown).BuildArgs()
	})
	return result.Error
}

func (tl *TestLifecycle) WaitForServices() error {
	result := tl.CLI.RunWithBuilder(func() []string {
		return NewCLIBuilder(core.CommandStatus).BuildArgs()
	})
	return result.Error
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
}

func joinServices(services []string) string {
	if len(services) == 0 {
		return ""
	}
	result := services[0]
	for i := 1; i < len(services); i++ {
		result += "," + services[i]
	}
	return result
}
