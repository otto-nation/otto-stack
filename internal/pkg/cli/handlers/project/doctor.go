package project

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/utils"
	"github.com/otto-nation/otto-stack/internal/pkg/constants"
	"github.com/otto-nation/otto-stack/internal/pkg/logger"
	"github.com/otto-nation/otto-stack/internal/pkg/ui"
	"github.com/spf13/cobra"
)

type DoctorHandler struct {
	output *ui.Output
}

func NewDoctorHandler() *DoctorHandler {
	return &DoctorHandler{
		output: ui.NewOutput(),
	}
}

func (h *DoctorHandler) ValidateArgs(args []string) error {
	return nil
}

func (h *DoctorHandler) GetRequiredFlags() []string {
	return []string{}
}

func (h *DoctorHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
	// Check initialization first
	if err := utils.CheckInitialization(); err != nil {
		return err
	}

	logger.Info(constants.LogMsgProjectAction, constants.LogFieldAction, constants.CommandDoctor, constants.LogFieldProject, "health_check")

	base.Output.Header("🩺 Otto Stack Health Check")
	logger.Info("Starting health checks")

	allGood := true &&
		h.checkDocker(base) &&
		h.checkDockerCompose(base) &&
		h.checkProjectInit(base) &&
		h.checkConfiguration(base)

	if allGood {
		base.Output.Success("All checks passed! Your otto-stack is healthy.")
		logger.Info("All health checks passed")
		return nil
	} else {
		base.Output.Error("Some issues found")
		logger.Error("Health checks failed")
		return fmt.Errorf("health check failed")
	}
}

func (h *DoctorHandler) checkDocker(base *base.BaseCommand) bool {
	base.Output.Info("%s", constants.MsgDoctor_checking_docker)

	if !h.isCommandAvailable(constants.DockerCmd) {
		base.Output.Error("%s", constants.MsgDoctor_docker_not_found)
		base.Output.Muted(constants.MsgDoctor_docker_install_help, constants.DockerInstallURL)
		return false
	}

	// Check if Docker daemon is running
	cmd := exec.Command(constants.DockerCmd, constants.DockerInfoCmd)
	if err := cmd.Run(); err != nil {
		base.Output.Error("%s", constants.MsgDoctor_docker_daemon_not_running)
		base.Output.Muted("%s", constants.MsgDoctor_docker_start_help)
		return false
	}

	base.Output.Success("%s", constants.MsgDoctor_docker_available)
	return true
}

func (h *DoctorHandler) checkDockerCompose(base *base.BaseCommand) bool {
	base.Output.Info("%s", constants.MsgDoctor_checking_docker_compose)

	if !h.hasDockerComposePlugin() {
		base.Output.Error("%s", constants.MsgDoctor_docker_compose_not_found)
		base.Output.Muted("%s", constants.MsgDoctor_docker_compose_integrated)
		base.Output.Muted("%s", constants.MsgDoctor_docker_compose_update)
		return false
	}

	base.Output.Success("%s", constants.MsgDoctor_docker_compose_available)
	return true
}

func (h *DoctorHandler) checkProjectInit(base *base.BaseCommand) bool {
	base.Output.Info("%s", constants.MsgDoctor_checking_project_init)

	configPath := filepath.Join(constants.OttoStackDir, constants.ConfigFileName)

	if _, err := os.Stat(configPath); err != nil {
		base.Output.Error("%s", constants.MsgDoctor_project_not_initialized)
		base.Output.Muted(constants.MsgDoctor_run_init_help, constants.AppName+" init")
		return false
	}

	base.Output.Success("%s", constants.MsgDoctor_project_initialized)
	return true
}

func (h *DoctorHandler) checkConfiguration(base *base.BaseCommand) bool {
	base.Output.Info("%s", constants.MsgDoctor_checking_config)

	// Check if otto-stack directory exists
	if _, err := os.Stat(constants.OttoStackDir); os.IsNotExist(err) {
		base.Output.Error("%s", constants.MsgDoctor_config_dir_missing)
		base.Output.Muted(constants.MsgDoctor_run_init_help, constants.AppName+" init")
		return false
	}

	// Check if docker-compose file exists
	composePath := filepath.Join(constants.OttoStackDir, constants.DockerComposeFileName)
	if _, err := os.Stat(composePath); os.IsNotExist(err) {
		base.Output.Error("%s", constants.MsgDoctor_docker_compose_missing)
		base.Output.Muted("%s", constants.MsgDoctor_config_incomplete)
		return false
	}

	base.Output.Success("%s", constants.MsgDoctor_config_valid)
	return true
}

func (h *DoctorHandler) isCommandAvailable(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

func (h *DoctorHandler) hasDockerComposePlugin() bool {
	cmd := exec.Command(constants.DockerCmd, constants.DockerComposeCmd, constants.DockerVersionCmd)
	return cmd.Run() == nil
}
