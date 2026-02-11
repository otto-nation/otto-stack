package project

import (
	"os"
	"os/exec"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/messages"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
)

type ValidationFunc func(*InitHandler, []types.ServiceConfig, *base.BaseCommand) error
type CheckFunc func(*InitHandler, []types.ServiceConfig, *base.BaseCommand)

var ValidationRegistry = map[string]ValidationFunc{
	core.ValidationDocker:             validateDocker,
	core.ValidationConfigSyntax:       validateNoFileConflicts,
	core.ValidationServiceDefinitions: validateServices,
}

var CheckRegistry = map[string]CheckFunc{
	core.ValidationFilePermissions: checkGitRepository,
}

func validateDocker(h *InitHandler, serviceConfigs []types.ServiceConfig, base *base.BaseCommand) error {
	if !isCommandAvailable(docker.DockerCmd) {
		return pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, "", messages.ValidationRequiredToolUnavailable, nil)
	}
	return nil
}

func validateNoFileConflicts(h *InitHandler, serviceConfigs []types.ServiceConfig, base *base.BaseCommand) error {
	// Skip validation if force flag is set
	if h.forceOverwrite {
		return nil
	}

	conflictingFiles := []string{docker.DockerComposeFileName, docker.DockerComposeFileNameYaml}
	for _, file := range conflictingFiles {
		if _, err := os.Stat(file); err == nil {
			return pkgerrors.NewValidationError(pkgerrors.ErrCodeAlreadyExists, "", messages.ValidationConflictingFileExists, nil)
		}
	}
	return nil
}

func validateServices(h *InitHandler, serviceConfigs []types.ServiceConfig, base *base.BaseCommand) error {
	if len(serviceConfigs) == 0 {
		return pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldServiceName, messages.ValidationNoServicesSelected, nil)
	}

	// Check for duplicates
	seen := make(map[string]bool)
	for _, cfg := range serviceConfigs {
		if seen[cfg.Name] {
			return pkgerrors.NewValidationErrorf(pkgerrors.ErrCodeInvalid, pkgerrors.FieldServiceName, messages.ValidationDuplicateService, cfg.Name)
		}
		seen[cfg.Name] = true
	}

	// Validate each service exists and is loadable
	serviceUtils := services.NewServiceUtils()
	for _, serviceConfig := range serviceConfigs {
		if _, err := serviceUtils.LoadServiceConfig(serviceConfig.Name); err != nil {
			return pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldServiceName, messages.ValidationInvalidService, err)
		}
	}
	return nil
}

func checkGitRepository(h *InitHandler, serviceConfigs []types.ServiceConfig, base *base.BaseCommand) {
	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		base.Output.Warning("%s", messages.WarningsNotGitRepository)
	}
}

func isCommandAvailable(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}
