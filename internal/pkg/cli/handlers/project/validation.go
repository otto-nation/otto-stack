package project

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/utils"
	"github.com/otto-nation/otto-stack/internal/pkg/constants"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
)

// validateInitEnvironment validates the environment before initialization
func (h *InitHandler) validateInitEnvironment(_ *types.BaseCommand) error {
	// Check if already initialized
	configPath := filepath.Join(constants.OttoStackDir, constants.ConfigFileName)
	if _, err := os.Stat(configPath); err == nil {
		return fmt.Errorf("%s is already initialized in this directory", constants.AppNameLower)
	}
	// Also check for .yaml extension
	yamlConfigPath := filepath.Join(constants.OttoStackDir, constants.ConfigFileName)
	if _, err := os.Stat(yamlConfigPath); err == nil {
		return fmt.Errorf("%s is already initialized in this directory", constants.AppNameLower)
	}

	// Check for required tools
	requiredTools := []string{"docker"}
	for _, tool := range requiredTools {
		if !h.isCommandAvailable(tool) {
			return fmt.Errorf(constants.Messages[constants.MsgValidation_required_tool_unavailable], tool)
		}
	}

	return nil
}

// validateProjectName validates the project name
func (h *InitHandler) validateProjectName(name string) error {
	if name == "" {
		return fmt.Errorf("%s", constants.Messages[constants.MsgValidation_project_name_empty])
	}

	if len(name) < constants.MinProjectNameLength {
		return fmt.Errorf("%s", constants.Messages[constants.MsgValidation_project_name_too_short])
	}

	if len(name) > constants.MaxProjectNameLength {
		return fmt.Errorf("%s", constants.Messages[constants.MsgValidation_project_name_too_long])
	}

	// Check for valid characters (alphanumeric, hyphens, underscores)
	for _, char := range name {
		if (char < 'a' || char > 'z') && (char < 'A' || char > 'Z') &&
			(char < '0' || char > '9') && char != '-' && char != '_' {
			return fmt.Errorf("%s", constants.Messages[constants.MsgValidation_project_name_invalid_chars])
		}
	}

	// Cannot start with hyphen or underscore
	if name[0] == '-' || name[0] == '_' {
		return fmt.Errorf("%s", constants.Messages[constants.MsgValidation_project_name_invalid_start])
	}

	return nil
}

// validateServices validates the selected services
func (h *InitHandler) validateServices(services []string) error {
	if len(services) == 0 {
		return fmt.Errorf("%s", constants.Messages[constants.MsgValidation_no_services_selected])
	}

	// Check for duplicates
	seen := make(map[string]bool)
	for _, serviceName := range services {
		if seen[serviceName] {
			return fmt.Errorf(constants.Messages[constants.MsgValidation_duplicate_service], serviceName)
		}
		seen[serviceName] = true
	}

	// Validate each service exists
	serviceUtils := utils.NewServiceUtils()
	for _, serviceName := range services {
		if _, err := serviceUtils.LoadServiceConfig(serviceName); err != nil {
			return fmt.Errorf(constants.Messages[constants.MsgValidation_invalid_service], serviceName, err)
		}
	}

	return nil
}

// validateDirectoryStructure ensures the directory structure is valid for initialization
func (h *InitHandler) validateDirectoryStructure(base *types.BaseCommand) error {
	// Check if we're in a git repository (optional but recommended)
	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		base.Output.Warning("%s", constants.Messages[constants.MsgWarnings_not_git_repository])
	}

	// Check for conflicting files
	conflictingFiles := []string{
		"docker-compose.yml",
		"docker-compose.yaml",
	}

	for _, file := range conflictingFiles {
		if _, err := os.Stat(file); err == nil {
			return fmt.Errorf(constants.Messages[constants.MsgValidation_conflicting_file_exists], file)
		}
	}

	return nil
}

// isCommandAvailable checks if a command is available in PATH
func (h *InitHandler) isCommandAvailable(command string) bool {
	if command == "" {
		return false
	}

	// Use exec.LookPath which is cross-platform
	_, err := exec.LookPath(command)
	return err == nil
}
