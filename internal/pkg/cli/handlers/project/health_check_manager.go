package project

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/common"
	"github.com/otto-nation/otto-stack/internal/pkg/messages"
)

// HealthCheckManager handles system health checks
type HealthCheckManager struct{}

// NewHealthCheckManager creates a new health check manager
func NewHealthCheckManager() *HealthCheckManager {
	return &HealthCheckManager{}
}

// RunAllChecks executes all health checks and returns overall status
func (hcm *HealthCheckManager) RunAllChecks(ctx context.Context, base *base.BaseCommand) bool {
	return hcm.CheckDocker(ctx, base) &&
		hcm.CheckDockerCompose(base) &&
		hcm.CheckProjectInit(base) &&
		hcm.CheckConfiguration(base)
}

// CheckDocker verifies Docker is available and running
func (hcm *HealthCheckManager) CheckDocker(ctx context.Context, base *base.BaseCommand) bool {
	base.Output.Info(messages.DoctorCheckingDocker)

	if !hcm.isCommandAvailable(docker.DockerCmd) {
		base.Output.Error(messages.DoctorDockerNotFound)
		base.Output.Info(messages.DoctorDockerInstallHelp, "https://docs.docker.com/get-docker/")
		return false
	}

	// Use StackService to run docker info command
	stackService, err := common.NewServiceManager(false)
	if err != nil {
		base.Output.Error("Failed to create stack service: %v", err)
		return false
	}

	err = stackService.CheckDockerHealth(ctx)
	if err != nil {
		base.Output.Error(messages.DoctorDockerDaemonNotRunning)
		base.Output.Info(messages.DoctorDockerStartHelp)
		return false
	}

	base.Output.Success(messages.DoctorDockerAvailable)
	return true
}

// CheckDockerCompose verifies Docker Compose is available
func (hcm *HealthCheckManager) CheckDockerCompose(base *base.BaseCommand) bool {
	base.Output.Info(messages.DoctorCheckingDockerCompose)

	if hcm.hasDockerComposePlugin() {
		base.Output.Success(messages.DoctorDockerComposeAvailable)
		return true
	}

	// Check if docker compose command is available
	composeCommand := fmt.Sprintf("%s %s", docker.DockerCmd, docker.DockerComposeCmd)
	if hcm.isCommandAvailable(composeCommand) {
		base.Output.Success(messages.DoctorDockerComposeAvailable)
		return true
	}

	base.Output.Error(messages.DoctorDockerComposeNotFound)
	base.Output.Info(messages.DoctorDockerComposeUpdate)
	return false
}

// CheckProjectInit verifies project is initialized
func (hcm *HealthCheckManager) CheckProjectInit(base *base.BaseCommand) bool {
	base.Output.Info(messages.DoctorCheckingProjectInit)

	if _, err := os.Stat(core.OttoStackDir); os.IsNotExist(err) {
		base.Output.Error(messages.DoctorProjectNotInitialized)
		base.Output.Info(messages.DoctorRunInitHelp, "otto init")
		return false
	}

	base.Output.Success(messages.DoctorProjectInitialized)
	return true
}

// CheckConfiguration verifies configuration files exist and are valid
func (hcm *HealthCheckManager) CheckConfiguration(base *base.BaseCommand) bool {
	base.Output.Info(messages.DoctorCheckingConfig)

	configPath := filepath.Join(core.OttoStackDir, core.ConfigFileName)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		base.Output.Error(messages.DoctorConfigDirMissing)
		base.Output.Info("   Expected: %s", configPath)
		return false
	}

	base.Output.Success(messages.DoctorConfigValid)
	return true
}

// isCommandAvailable checks if a command is available in PATH
func (hcm *HealthCheckManager) isCommandAvailable(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

// hasDockerComposePlugin checks if Docker Compose plugin is available
func (hcm *HealthCheckManager) hasDockerComposePlugin() bool {
	stackService, err := common.NewServiceManager(false)
	if err != nil {
		return false
	}

	// Check if compose is available by checking Docker health
	err = stackService.CheckDockerHealth(context.Background())
	return err == nil
}
