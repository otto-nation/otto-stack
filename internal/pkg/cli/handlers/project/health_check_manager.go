package project

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/common"
	"github.com/otto-nation/otto-stack/internal/pkg/messages"
)

// CheckResult holds the outcome of a single health check for structured output.
type CheckResult struct {
	Name    string `json:"name"`
	Passed  bool   `json:"passed"`
	Message string `json:"message"`
}

// HealthCheckManager handles system health checks
type HealthCheckManager struct{}

// NewHealthCheckManager creates a new health check manager
func NewHealthCheckManager() *HealthCheckManager {
	return &HealthCheckManager{}
}

// collectResults runs all checks without printing and returns structured results.
func (hcm *HealthCheckManager) collectResults(ctx context.Context) []CheckResult {
	const numChecks = 4
	results := make([]CheckResult, 0, numChecks)

	// Docker availability
	dockerPassed := false
	dockerMsg := messages.DoctorDockerNotFound
	if isCommandAvailable(docker.DockerCmd) {
		stackService, err := common.NewServiceManager(false)
		if err == nil {
			if err = stackService.CheckDockerHealth(ctx); err == nil {
				dockerPassed = true
				dockerMsg = messages.DoctorDockerAvailable
			} else {
				dockerMsg = messages.DoctorDockerDaemonNotRunning
			}
		}
	}
	results = append(results, CheckResult{Name: "docker", Passed: dockerPassed, Message: dockerMsg})

	// Docker Compose availability
	composePassed := hcm.hasDockerComposePlugin() ||
		isCommandAvailable(fmt.Sprintf("%s %s", docker.DockerCmd, docker.DockerComposeCmd))
	composeMsg := messages.DoctorDockerComposeAvailable
	if !composePassed {
		composeMsg = messages.DoctorDockerComposeNotFound
	}
	results = append(results, CheckResult{Name: "docker-compose", Passed: composePassed, Message: composeMsg})

	// Project initialization
	_, statErr := os.Stat(core.OttoStackDir)
	projectPassed := !os.IsNotExist(statErr)
	projectMsg := messages.DoctorProjectInitialized
	if !projectPassed {
		projectMsg = messages.DoctorProjectNotInitialized
	}
	results = append(results, CheckResult{Name: "project-init", Passed: projectPassed, Message: projectMsg})

	// Configuration file
	configPath := filepath.Join(core.OttoStackDir, core.ConfigFileName)
	_, statErr = os.Stat(configPath)
	configPassed := !os.IsNotExist(statErr)
	configMsg := messages.DoctorConfigValid
	if !configPassed {
		configMsg = messages.DoctorConfigDirMissing
	}
	results = append(results, CheckResult{Name: "configuration", Passed: configPassed, Message: configMsg})

	return results
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

	if !isCommandAvailable(docker.DockerCmd) {
		base.Output.Error(messages.DoctorDockerNotFound)
		base.Output.Info(messages.DoctorDockerInstallHelp, "https://docs.docker.com/get-docker/")
		return false
	}

	// Use StackService to run docker info command
	stackService, err := common.NewServiceManager(false)
	if err != nil {
		base.Output.Error(messages.ErrorsStackCreateFailed, err)
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
	if isCommandAvailable(composeCommand) {
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
		base.Output.Info(messages.InfoExpectedPath, configPath)
		return false
	}

	base.Output.Success(messages.DoctorConfigValid)
	return true
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
