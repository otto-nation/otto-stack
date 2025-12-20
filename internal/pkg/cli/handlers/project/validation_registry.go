package project

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
)

type ValidationFunc func(*InitHandler, *base.BaseCommand) error

var ValidationRegistry = map[string]ValidationFunc{
	core.ValidationDocker:             validateDocker,
	core.ValidationConfigSyntax:       validateConfigSyntax,
	core.ValidationServiceDefinitions: validateServiceDefinitions,
	core.ValidationFilePermissions:    validateFilePermissions,
}

func validateDocker(h *InitHandler, base *base.BaseCommand) error {
	if !isCommandAvailable(docker.DockerCmd) {
		return fmt.Errorf(core.MsgValidation_required_tool_unavailable, docker.DockerCmd)
	}
	return nil
}

func validateConfigSyntax(h *InitHandler, base *base.BaseCommand) error {
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

func validateServiceDefinitions(h *InitHandler, base *base.BaseCommand) error {
	if len(h.selectedServices) == 0 {
		return fmt.Errorf("%s", core.MsgValidation_no_services_selected)
	}

	serviceUtils := services.NewServiceUtils()
	for _, serviceName := range h.selectedServices {
		if _, err := serviceUtils.LoadServiceConfig(serviceName); err != nil {
			return fmt.Errorf(core.MsgValidation_invalid_service, serviceName, err)
		}
	}
	return nil
}

func validateFilePermissions(h *InitHandler, base *base.BaseCommand) error {
	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		base.Output.Warning("%s", core.MsgWarnings_not_git_repository)
	}
	return nil
}

func isCommandAvailable(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}
