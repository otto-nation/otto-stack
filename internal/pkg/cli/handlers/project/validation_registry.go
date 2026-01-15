package project

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
)

type ValidationFunc func(*InitHandler, []types.ServiceConfig, *base.BaseCommand) error

var ValidationRegistry = map[string]ValidationFunc{
	core.ValidationDocker:             validateDocker,
	core.ValidationConfigSyntax:       validateConfigSyntax,
	core.ValidationServiceDefinitions: validateServiceDefinitions,
	core.ValidationFilePermissions:    validateFilePermissions,
}

func validateDocker(h *InitHandler, serviceConfigs []types.ServiceConfig, base *base.BaseCommand) error {
	if !isCommandAvailable(docker.DockerCmd) {
		return fmt.Errorf(core.MsgValidation_required_tool_unavailable, docker.DockerCmd)
	}
	return nil
}

func validateConfigSyntax(h *InitHandler, serviceConfigs []types.ServiceConfig, base *base.BaseCommand) error {
	// Skip validation if force flag is set
	if h.forceOverwrite {
		return nil
	}

	conflictingFiles := []string{docker.DockerComposeFileName, docker.DockerComposeFileNameYaml}
	for _, file := range conflictingFiles {
		if _, err := os.Stat(file); err == nil {
			return fmt.Errorf(core.MsgValidation_conflicting_file_exists, file)
		}
	}
	return nil
}

func validateServiceDefinitions(h *InitHandler, serviceConfigs []types.ServiceConfig, base *base.BaseCommand) error {
	if len(serviceConfigs) == 0 {
		return fmt.Errorf("%s", core.MsgValidation_no_services_selected)
	}

	serviceUtils := services.NewServiceUtils()
	for _, serviceConfig := range serviceConfigs {
		if _, err := serviceUtils.LoadServiceConfig(serviceConfig.Name); err != nil {
			return fmt.Errorf(core.MsgValidation_invalid_service, serviceConfig.Name, err)
		}
	}
	return nil
}

func validateFilePermissions(h *InitHandler, serviceConfigs []types.ServiceConfig, base *base.BaseCommand) error {
	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		base.Output.Warning("%s", core.MsgWarnings_not_git_repository)
	}
	return nil
}

func isCommandAvailable(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}
