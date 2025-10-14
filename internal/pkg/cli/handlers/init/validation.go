package init

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/utils"
	"github.com/otto-nation/otto-stack/internal/pkg/ui"
)

// validateInitEnvironment validates the environment before initialization
func (h *InitHandler) validateInitEnvironment() error {
	// Check if already initialized
	if _, err := os.Stat("otto-stack/otto-stack-config.yml"); err == nil {
		return fmt.Errorf("otto-stack is already initialized in this directory")
	}
	if _, err := os.Stat("otto-stack/otto-stack-config.yaml"); err == nil {
		return fmt.Errorf("otto-stack is already initialized in this directory")
	}

	// Check for required tools
	requiredTools := []string{"docker"}
	for _, tool := range requiredTools {
		if !h.isCommandAvailable(tool) {
			return fmt.Errorf("required tool '%s' is not available", tool)
		}
	}

	return nil
}

// validateProjectName validates the project name
func (h *InitHandler) validateProjectName(name string) error {
	if name == "" {
		return fmt.Errorf("project name cannot be empty")
	}

	if len(name) < 2 {
		return fmt.Errorf("project name must be at least 2 characters long")
	}

	if len(name) > 50 {
		return fmt.Errorf("project name must be less than 50 characters")
	}

	// Check for valid characters (alphanumeric, hyphens, underscores)
	for _, char := range name {
		if (char < 'a' || char > 'z') && (char < 'A' || char > 'Z') &&
			(char < '0' || char > '9') && char != '-' && char != '_' {
			return fmt.Errorf("project name can only contain letters, numbers, hyphens, and underscores")
		}
	}

	// Cannot start with hyphen or underscore
	if name[0] == '-' || name[0] == '_' {
		return fmt.Errorf("project name cannot start with a hyphen or underscore")
	}

	return nil
}

// validateServices validates the selected services
func (h *InitHandler) validateServices(services []string) error {
	if len(services) == 0 {
		return fmt.Errorf("at least one service must be selected")
	}

	// Check for duplicates
	seen := make(map[string]bool)
	for _, serviceName := range services {
		if seen[serviceName] {
			return fmt.Errorf("duplicate service '%s' found", serviceName)
		}
		seen[serviceName] = true
	}

	// Validate each service exists
	serviceUtils := utils.NewServiceUtils()
	for _, serviceName := range services {
		if _, err := serviceUtils.LoadServiceConfig(serviceName); err != nil {
			return fmt.Errorf("invalid service '%s': %w", serviceName, err)
		}
	}

	return nil
}

// validateDirectoryStructure ensures the directory structure is valid for initialization
func (h *InitHandler) validateDirectoryStructure() error {
	// Check if we're in a git repository (optional but recommended)
	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		ui.Warning("Not in a git repository. Consider running 'git init' first.")
	}

	// Check for conflicting files
	conflictingFiles := []string{
		"docker-compose.yml",
		"docker-compose.yaml",
	}

	for _, file := range conflictingFiles {
		if _, err := os.Stat(file); err == nil {
			return fmt.Errorf("conflicting file '%s' already exists. Please remove or rename it before initializing", file)
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
