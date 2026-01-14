package project

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
	"github.com/otto-nation/otto-stack/internal/pkg/compose"
	"github.com/otto-nation/otto-stack/internal/pkg/config"
	"github.com/otto-nation/otto-stack/internal/pkg/env"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	svc "github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
)

// ProjectManager handles project creation logic
type ProjectManager struct {
	serviceUtils     *svc.ServiceUtils
	configService    config.ConfigService
	configManager    *ConfigManager
	directoryManager *DirectoryManager
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
		serviceUtils:     svc.NewServiceUtils(),
		configService:    config.NewConfigService(),
		configManager:    NewConfigManager(),
		directoryManager: NewDirectoryManager(),
	}
}

// CreateProjectStructure creates the complete project structure
func (pm *ProjectManager) CreateProjectStructure(projectCtx clicontext.Context, base *base.BaseCommand) error {
	if err := pm.directoryManager.CreateDirectoryStructure(); err != nil {
		return pkgerrors.NewServiceError(ComponentProject, ActionCreateDirectories, err)
	}

	if err := pm.configManager.CreateConfigFile(projectCtx.Project.Name, projectCtx.Services.Names, projectCtx.Options.Validation, base); err != nil {
		return pkgerrors.NewConfigError("", ActionCreateConfigFile, err)
	}

	pm.configManager.GenerateServiceConfigs(projectCtx.Services.Configs, base)

	if err := pm.generateInitialComposeFiles(projectCtx.Services.Configs, projectCtx.Project.Name, projectCtx.Options.Validation, projectCtx.Options.Advanced, base); err != nil {
		return pkgerrors.NewServiceError(ComponentCompose, ActionGenerateFiles, err)
	}

	if err := pm.createGitignoreEntries(base); err != nil {
		base.Output.Warning("Failed to create .gitignore entries: %v", err)
	}

	if err := pm.createReadme(projectCtx.Project.Name, projectCtx.Services.Configs, base); err != nil {
		base.Output.Warning("Failed to create README: %v", err)
	}

	return nil
}

// generateInitialComposeFiles generates Docker Compose files
func (pm *ProjectManager) generateInitialComposeFiles(serviceConfigs []types.ServiceConfig, projectName string, _, _ map[string]bool, base *base.BaseCommand) error {
	if err := pm.generateEnvFile(serviceConfigs, projectName, base); err != nil {
		return err
	}

	if err := pm.generateDockerCompose(serviceConfigs, projectName, base); err != nil {
		return err
	}

	return nil
}

// generateEnvFile generates the .env file
func (pm *ProjectManager) generateEnvFile(serviceConfigs []types.ServiceConfig, projectName string, base *base.BaseCommand) error {
	if err := env.GenerateFile(projectName, serviceConfigs, core.EnvGeneratedFileName); err != nil {
		return err
	}

	base.Output.Success("Created environment file: %s", core.EnvGeneratedFileName)
	return nil
}

// generateDockerCompose generates the docker-compose.yml file
func (pm *ProjectManager) generateDockerCompose(serviceConfigs []types.ServiceConfig, projectName string, base *base.BaseCommand) error {
	manager, err := svc.New()
	if err != nil {
		return err
	}

	generator, err := compose.NewGenerator(projectName, "", manager)
	if err != nil {
		return err
	}
	err = generator.GenerateFromServiceConfigs(serviceConfigs, projectName)
	if err != nil {
		return err
	}

	base.Output.Success("Created Docker Compose file: %s", docker.DockerComposeFilePath)
	return nil
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
func (pm *ProjectManager) createReadme(projectName string, serviceConfigs []types.ServiceConfig, base *base.BaseCommand) error {
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
		pm.formatServicesList(svc.ExtractServiceNames(serviceConfigs)),
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
