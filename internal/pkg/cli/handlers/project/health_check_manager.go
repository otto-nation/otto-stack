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
)

// HealthCheckManager handles system health checks
type HealthCheckManager struct{}

// NewHealthCheckManager creates a new health check manager
func NewHealthCheckManager() *HealthCheckManager {
	return &HealthCheckManager{}
}

// RunAllChecks executes all health checks and returns overall status
func (hcm *HealthCheckManager) RunAllChecks(base *base.BaseCommand) bool {
	return hcm.CheckDocker(base) &&
		hcm.CheckDockerCompose(base) &&
		hcm.CheckProjectInit(base) &&
		hcm.CheckConfiguration(base)
}

// CheckDocker verifies Docker is available and running
func (hcm *HealthCheckManager) CheckDocker(base *base.BaseCommand) bool {
	base.Output.Info(core.MsgDoctor_checking_docker)

	if !hcm.isCommandAvailable(docker.DockerCmd) {
		base.Output.Error(core.MsgDoctor_docker_not_found)
		base.Output.Info(core.MsgDoctor_docker_install_help, "https://docs.docker.com/get-docker/")
		return false
	}

	// Use StackService to run docker info command
	stackService, err := common.NewServiceManager(false)
	if err != nil {
		base.Output.Error("Failed to create stack service: %v", err)
		return false
	}

	_, err = stackService.DockerClient.GetCli().Info(context.Background())
	if err != nil {
		base.Output.Error(core.MsgDoctor_docker_daemon_not_running)
		base.Output.Info(core.MsgDoctor_docker_start_help)
		return false
	}

	base.Output.Success(core.MsgDoctor_docker_available)
	return true
}

// CheckDockerCompose verifies Docker Compose is available
func (hcm *HealthCheckManager) CheckDockerCompose(base *base.BaseCommand) bool {
	base.Output.Info(core.MsgDoctor_checking_docker_compose)

	if hcm.hasDockerComposePlugin() {
		base.Output.Success(core.MsgDoctor_docker_compose_available)
		return true
	}

	// Check if docker compose command is available
	composeCommand := fmt.Sprintf("%s %s", docker.DockerCmd, docker.DockerComposeCmd)
	if hcm.isCommandAvailable(composeCommand) {
		base.Output.Success(core.MsgDoctor_docker_compose_available)
		return true
	}

	base.Output.Error(core.MsgDoctor_docker_compose_not_found)
	base.Output.Info(core.MsgDoctor_docker_compose_update)
	return false
}

// CheckProjectInit verifies project is initialized
func (hcm *HealthCheckManager) CheckProjectInit(base *base.BaseCommand) bool {
	base.Output.Info(core.MsgDoctor_checking_project_init)

	if _, err := os.Stat(core.OttoStackDir); os.IsNotExist(err) {
		base.Output.Error(core.MsgDoctor_project_not_initialized)
		base.Output.Info(core.MsgDoctor_run_init_help, "otto init")
		return false
	}

	base.Output.Success(core.MsgDoctor_project_initialized)
	return true
}

// CheckConfiguration verifies configuration files exist and are valid
func (hcm *HealthCheckManager) CheckConfiguration(base *base.BaseCommand) bool {
	base.Output.Info(core.MsgDoctor_checking_config)

	configPath := filepath.Join(core.OttoStackDir, core.ConfigFileName)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		base.Output.Error(core.MsgDoctor_config_dir_missing)
		base.Output.Info("   Expected: %s", configPath)
		return false
	}

	base.Output.Success(core.MsgDoctor_config_valid)
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

	// Check if compose is available by trying to create a compose service
	_, err = stackService.DockerClient.GetCli().Info(context.Background())
	return err == nil
}
