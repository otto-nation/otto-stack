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
	core.ValidationDocker:       validateDocker,
	core.ValidationConfigSyntax: validateNoFileConflicts,
	core.ValidationServiceDefinitions: func(h *InitHandler, serviceConfigs []types.ServiceConfig, base *base.BaseCommand) error {
		return services.NewValidator().ValidateServiceConfigs(serviceConfigs)
	},
}

var CheckRegistry = map[string]CheckFunc{
	core.ValidationFilePermissions: checkGitRepository,
}

func validateDocker(h *InitHandler, serviceConfigs []types.ServiceConfig, base *base.BaseCommand) error {
	if !isCommandAvailable(docker.DockerCmd) {
		return pkgerrors.NewSystemError(pkgerrors.ErrCodeInvalid, messages.ValidationRequiredToolUnavailable, nil)
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
			return pkgerrors.NewSystemError(pkgerrors.ErrCodeAlreadyExists, messages.ValidationConflictingFileExists, nil)
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
