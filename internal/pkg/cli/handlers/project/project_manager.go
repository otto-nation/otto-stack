package project

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
)

// ProjectManager handles project creation logic
type ProjectManager struct {
	serviceUtils  *services.ServiceUtils
	configService services.ConfigService
}

// NewProjectManager creates a new project manager
func NewProjectManager() *ProjectManager {
	return &ProjectManager{
		serviceUtils:  services.NewServiceUtils(),
		configService: services.NewConfigService(),
	}
}

// CreateProjectStructure creates the complete project structure
func (pm *ProjectManager) CreateProjectStructure(projectName string, services []string, validation, advanced map[string]bool, base *base.BaseCommand) error {
	if err := pm.createDirectoryStructure(); err != nil {
		return pkgerrors.NewServiceError(ComponentProject, ActionCreateDirectories, err)
	}

	if err := pm.createConfigFile(projectName, services, validation, base); err != nil {
		return pkgerrors.NewConfigError("", ActionCreateConfigFile, err)
	}

	pm.generateServiceConfigs(services, base)

	if err := pm.generateInitialComposeFiles(services, projectName, validation, advanced, base); err != nil {
		return pkgerrors.NewServiceError(ComponentCompose, ActionGenerateFiles, err)
	}

	if err := pm.createGitignoreEntries(base); err != nil {
		base.Output.Warning("Failed to create .gitignore entries: %v", err)
	}

	if err := pm.createReadme(projectName, services, base); err != nil {
		base.Output.Warning("Failed to create README: %v", err)
	}

	return nil
}

// createDirectoryStructure creates the necessary directory structure
func (pm *ProjectManager) createDirectoryStructure() error {
	directories := []string{
		core.OttoStackDir,
		filepath.Join(core.OttoStackDir, core.ServiceConfigsDir),
	}

	for _, dir := range directories {
		if err := os.MkdirAll(dir, core.PermReadWriteExec); err != nil {
			return pkgerrors.NewConfigError(dir, MsgFailedToCreateDirectory, err)
		}
	}

	return nil
}

// createConfigFile creates the main configuration file
func (pm *ProjectManager) createConfigFile(projectName string, services []string, validationOptions map[string]bool, base *base.BaseCommand) error {
	configContent := pm.generateConfig(projectName, services, validationOptions)

	configPath := core.OttoStackDir + "/" + core.ConfigFileName
	if err := os.WriteFile(configPath, []byte(configContent), core.PermReadWrite); err != nil {
		return pkgerrors.NewConfigError(configPath, MsgFailedToWriteConfigFile, err)
	}

	base.Output.Success("Created configuration file: %s", configPath)
	return nil
}

// generateConfig generates the configuration content
func (pm *ProjectManager) generateConfig(name string, services []string, validationOptions map[string]bool) string {
	return fmt.Sprintf(`project:
  name: %s
  services: [%s]
validation:
  enabled: %t
`, name, strings.Join(services, ", "), len(validationOptions) > 0)
}

// generateServiceConfigs generates service configurations
func (pm *ProjectManager) generateServiceConfigs(services []string, base *base.BaseCommand) {
	for _, serviceName := range services {
		if err := pm.generateServiceConfig(serviceName); err != nil {
			base.Output.Warning("Failed to generate config for service %s: %v", serviceName, err)
		}
	}
}

// generateServiceConfig generates configuration for a single service
func (pm *ProjectManager) generateServiceConfig(serviceName string) error {
	content := pm.generateServiceConfigContent(serviceName)
	configPath := filepath.Join(core.OttoStackDir, core.ServiceConfigsDir, serviceName+".yml")
	return os.WriteFile(configPath, []byte(content), core.PermReadWrite)
}

// generateServiceConfigContent generates the content for a service configuration
func (pm *ProjectManager) generateServiceConfigContent(serviceName string) string {
	return fmt.Sprintf(`name: %s
description: Configuration for %s service
`, serviceName, serviceName)
}

// generateInitialComposeFiles generates Docker Compose files
func (pm *ProjectManager) generateInitialComposeFiles(services []string, projectName string, _, _ map[string]bool, base *base.BaseCommand) error {
	if err := pm.generateEnvFile(services, projectName, base); err != nil {
		return err
	}

	if err := pm.generateDockerCompose(services, projectName, base); err != nil {
		return err
	}

	return nil
}

// generateEnvFile generates the .env file
func (pm *ProjectManager) generateEnvFile(_ []string, projectName string, base *base.BaseCommand) error {
	envContent := fmt.Sprintf("PROJECT_NAME=%s\n", projectName)
	envPath := ".env"

	if err := os.WriteFile(envPath, []byte(envContent), core.PermReadWrite); err != nil {
		return err
	}

	base.Output.Success("Created environment file: %s", envPath)
	return nil
}

// generateDockerCompose generates the docker-compose.yml file
func (pm *ProjectManager) generateDockerCompose(services []string, _ string, base *base.BaseCommand) error {
	composeContent := fmt.Sprintf(`version: '3.8'
services:
%s
`, pm.formatServicesForCompose(services))

	composePath := "docker-compose.yml"
	if err := os.WriteFile(composePath, []byte(composeContent), core.PermReadWrite); err != nil {
		return err
	}

	base.Output.Success("Created Docker Compose file: %s", composePath)
	return nil
}

// formatServicesForCompose formats services for docker-compose.yml
func (pm *ProjectManager) formatServicesForCompose(services []string) string {
	var result strings.Builder
	for _, service := range services {
		result.WriteString(fmt.Sprintf("  %s:\n    image: %s:latest\n", service, service))
	}
	return result.String()
}

// createGitignoreEntries adds entries to .gitignore
func (pm *ProjectManager) createGitignoreEntries(base *base.BaseCommand) error {
	gitignorePath := ".gitignore"

	entries := []string{
		"# Otto Stack",
		".otto-stack/logs/",
		".env.local",
		"*.log",
	}

	content := strings.Join(entries, "\n") + "\n"

	if err := os.WriteFile(gitignorePath, []byte(content), core.PermReadWrite); err != nil {
		return err
	}

	base.Output.Success("Updated .gitignore file")
	return nil
}

// createReadme creates README file
func (pm *ProjectManager) createReadme(projectName string, services []string, base *base.BaseCommand) error {
	readmeContent := fmt.Sprintf(`# %s

This project was initialized with Otto Stack.

## Services
%s

## Getting Started
1. Review configuration in .otto-stack/config.yml
2. Start the stack: otto up
3. Check status: otto status
`, projectName, pm.formatServicesList(services))

	readmePath := filepath.Join(core.OttoStackDir, "README.md")
	if err := os.WriteFile(readmePath, []byte(readmeContent), core.PermReadWrite); err != nil {
		return err
	}

	base.Output.Success("Created README file: %s", readmePath)
	return nil
}

// formatServicesList formats services for README
func (pm *ProjectManager) formatServicesList(services []string) string {
	var result strings.Builder
	for _, service := range services {
		result.WriteString(fmt.Sprintf("- %s\n", service))
	}
	return result.String()
}
