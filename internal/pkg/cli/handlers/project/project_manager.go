package project

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
)

// ProjectManager handles project creation logic
type ProjectManager struct {
	serviceUtils  *services.ServiceUtils
	configService services.ConfigService
	fileGenerator *services.FileGenerator
}

// OttoStackConfig represents the otto-stack configuration structure
type OttoStackConfig struct {
	Project struct {
		Name string `yaml:"name"`
		Type string `yaml:"type"`
	} `yaml:"project"`
	Stack struct {
		Enabled []string `yaml:"enabled"`
	} `yaml:"stack"`
	Validation struct {
		Options map[string]bool `yaml:"options"`
	} `yaml:"validation"`
}

// ServiceConfig represents a service configuration file
type ServiceConfig struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

// NewProjectManager creates a new project manager
func NewProjectManager() *ProjectManager {
	return &ProjectManager{
		serviceUtils:  services.NewServiceUtils(),
		configService: services.NewConfigService(),
		fileGenerator: services.NewFileGenerator(),
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
	config := OttoStackConfig{}
	config.Project.Name = name
	config.Project.Type = "docker"
	config.Stack.Enabled = services

	config.Validation.Options = make(map[string]bool)
	for _, key := range core.ValidationKeys {
		if len(validationOptions) == 0 || validationOptions[key] {
			config.Validation.Options[key] = true
		}
	}

	data, err := yaml.Marshal(&config)
	if err != nil {
		// Fallback to empty config if marshal fails
		return "# Error generating configuration\n"
	}

	return string(data)
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
	config := ServiceConfig{
		Name:        serviceName,
		Description: fmt.Sprintf("Configuration for %s service", serviceName),
	}

	data, err := yaml.Marshal(&config)
	if err != nil {
		// Fallback if marshal fails
		return fmt.Sprintf("name: %s\ndescription: Configuration for %s service\n", serviceName, serviceName)
	}

	// Add comment header
	header := fmt.Sprintf("# Documentation: %s/services/%s\n\n", core.DocsURL, serviceName)
	return header + string(data)
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
func (pm *ProjectManager) generateEnvFile(services []string, projectName string, base *base.BaseCommand) error {
	resolvedServices := pm.resolveServiceDependencies(services)

	if err := pm.fileGenerator.GenerateEnvFile(resolvedServices, projectName); err != nil {
		return err
	}

	base.Output.Success("Created environment file: %s", core.ExtENV)
	return nil
}

// generateDockerCompose generates the docker-compose.yml file
func (pm *ProjectManager) generateDockerCompose(services []string, projectName string, base *base.BaseCommand) error {
	resolvedServices := pm.resolveServiceDependencies(services)

	if err := pm.fileGenerator.GenerateComposeFileWithOriginal(resolvedServices, services, projectName); err != nil {
		return err
	}

	base.Output.Success("Created Docker Compose file: %s", docker.DockerComposeFilePath)
	return nil
}

// resolveServiceDependencies resolves configuration services to their container dependencies
func (pm *ProjectManager) resolveServiceDependencies(serviceList []string) []string {
	resolved := make(map[string]bool)

	for _, service := range serviceList {
		// For localstack-* services, resolve to localstack
		if strings.HasPrefix(service, services.ServiceLocalstack) {
			resolved[services.ServiceLocalstack] = true
		} else {
			resolved[service] = true
		}
	}

	var result []string
	for service := range resolved {
		result = append(result, service)
	}
	return result
}

// createGitignoreEntries adds entries to .gitignore
func (pm *ProjectManager) createGitignoreEntries(base *base.BaseCommand) error {
	entries := []string{
		"# " + core.AppNameTitle,
		core.OttoStackDir + "/logs/",
		core.ExtENV + core.LocalFileExtension,
		"*.log",
	}

	content := strings.Join(entries, "\n") + "\n"
	gitignorePath := filepath.Join(core.OttoStackDir, core.GitIgnoreFileName)

	if err := os.WriteFile(gitignorePath, []byte(content), core.PermReadWrite); err != nil {
		return err
	}

	base.Output.Success("Updated %s file", gitignorePath)
	return nil
}

// createReadme creates README file
func (pm *ProjectManager) createReadme(projectName string, services []string, base *base.BaseCommand) error {
	const readmeTemplate = `# %s

This project was initialized with %s.

## Services
%s

## Commands
- ` + "`%s up`" + ` - Start all services
- ` + "`%s down`" + ` - Stop all services
- ` + "`%s status`" + ` - Show service status
- ` + "`%s logs`" + ` - View service logs
- ` + "`%s validate`" + ` - Validate configuration

## Configuration
- Main config: ` + "`%s/%s`" + `
- Service configs: ` + "`%s/%s/`" + `

## File Management
- **Generated files** (` + "`docker-compose.yml`" + `, ` + "`env.generated`" + `): Automatically regenerated on each ` + "`up`" + ` command
- **Service configs** (` + "`service-configs/`" + `): Created during ` + "`init`" + `, preserved across ` + "`up`" + ` commands for user customization
- **User configs** (` + "`otto-stack-config.yml`" + `): Never overwritten, safe to edit

## Documentation
%s
`

	readmeContent := fmt.Sprintf(readmeTemplate,
		projectName,
		core.AppNameTitle,
		pm.formatServicesList(services),
		core.AppName, core.AppName, core.AppName, core.AppName, core.AppName,
		core.OttoStackDir, core.ConfigFileName,
		core.OttoStackDir, core.ServiceConfigsDir,
		core.DocsURL,
	)

	readmePath := filepath.Join(core.OttoStackDir, core.ReadmeFileName)
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
		fmt.Fprintf(&result, "- %s\n", service)
	}
	return result.String()
}
