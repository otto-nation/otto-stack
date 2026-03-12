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
	"github.com/otto-nation/otto-stack/internal/pkg/env"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/filesystem"
	"github.com/otto-nation/otto-stack/internal/pkg/messages"
	svc "github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
)

// ProjectManager handles project creation logic
type ProjectManager struct {
	serviceUtils     *svc.ServiceUtils
	configManager    *ConfigManager
	directoryManager *DirectoryManager
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
		configManager:    NewConfigManager(),
		directoryManager: NewDirectoryManager(),
	}
}

// CreateProjectStructure creates the complete project structure
func (pm *ProjectManager) CreateProjectStructure(projectCtx clicontext.Context, base *base.BaseCommand) error {
	if err := pm.directoryManager.CreateDirectoryStructure(); err != nil {
		return pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentProject, messages.ErrorsDirectoryCreateFailed, err)
	}

	if err := pm.configManager.CreateConfigFile(projectCtx, base); err != nil {
		return pkgerrors.NewConfigError(pkgerrors.ErrCodeOperationFail, "", messages.ErrorsConfigWriteFailed, err)
	}

	pm.configManager.GenerateServiceConfigs(projectCtx.Services.Configs, projectCtx.Sharing.Enabled, base)

	// Generate env file with ALL services (shared and non-shared)
	if err := pm.generateEnvFile(projectCtx.Services.Configs, projectCtx.Project.Name, base); err != nil {
		return pkgerrors.NewSystemError(pkgerrors.ErrCodeOperationFail, messages.ValidationFailedGenerateEnv, err)
	}

	// Filter services for project compose file (exclude shared services)
	projectServices := pm.filterProjectServices(projectCtx.Services.Configs, projectCtx.Sharing)
	hasSharingEnabled := projectCtx.Sharing != nil && projectCtx.Sharing.Enabled
	if err := pm.generateDockerComposeWithSharing(projectServices, projectCtx.Project.Name, hasSharingEnabled, base); err != nil {
		return pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, "compose", messages.ErrorsComposeGenerateFailed, err)
	}

	// Generate shared docker-compose.yml if sharing is enabled
	if projectCtx.Sharing != nil && projectCtx.Sharing.Enabled {
		homeDir, _ := os.UserHomeDir()
		sharedRoot := filepath.Join(homeDir, core.OttoStackDir, core.SharedDir)
		if err := pm.generateSharedCompose(projectCtx.Services.Configs, sharedRoot, base); err != nil {
			base.Output.Warning(messages.WarningsComposeGenerateSharedFailed, err)
		}
	}

	if err := pm.createComposeOverrideFile(base); err != nil {
		base.Output.Warning(messages.WarningsComposeGenerateSharedFailed, err)
	}

	if err := pm.createGitignoreEntries(base); err != nil {
		base.Output.Warning(messages.WarningsGitignoreCreateFailed, err)
	}

	if err := pm.createReadme(projectCtx.Project.Name, projectCtx.Services.Configs, projectCtx.Sharing.Enabled, base); err != nil {
		base.Output.Warning(messages.WarningsFailedReadme, err)
	}

	return nil
}

// generateEnvFile generates the .env file
func (pm *ProjectManager) generateEnvFile(serviceConfigs []types.ServiceConfig, projectName string, base *base.BaseCommand) error {
	if err := env.GenerateFile(projectName, serviceConfigs, core.EnvGeneratedFilePath); err != nil {
		return pkgerrors.NewSystemError(pkgerrors.ErrCodeOperationFail, messages.ValidationFailedGenerateEnv, err)
	}

	base.Output.Success(messages.SuccessCreatedEnvFile, core.EnvGeneratedFilePath)
	return nil
}

// generateDockerComposeWithSharing generates the docker-compose.yml file with sharing info
func (pm *ProjectManager) generateDockerComposeWithSharing(serviceConfigs []types.ServiceConfig, projectName string, hasSharingEnabled bool, base *base.BaseCommand) error {
	generator, err := compose.NewGenerator(projectName)
	if err != nil {
		return pkgerrors.NewSystemError(pkgerrors.ErrCodeOperationFail, messages.ErrorsComposeGeneratorCreateFailed, err)
	}

	var header string
	if hasSharingEnabled {
		homeDir, _ := os.UserHomeDir()
		sharedPath := filepath.Join(homeDir, core.OttoStackDir, core.SharedDir)
		header = fmt.Sprintf(core.ComposeHeaderProjectShared, sharedPath)
	} else {
		header = core.ComposeHeaderProject
	}

	content, err := generator.BuildComposeDataWithHeader(serviceConfigs, header)
	if err != nil {
		return pkgerrors.NewSystemError(pkgerrors.ErrCodeOperationFail, messages.ErrorsComposeGenerateFailed, err)
	}

	if err := filesystem.EnsureDir(core.GeneratedDir); err != nil {
		return pkgerrors.NewSystemError(pkgerrors.ErrCodeOperationFail, messages.ErrorsComposeDirCreateFailed, err)
	}

	if err := filesystem.WriteFile(docker.DockerComposeFilePath, content, core.PermReadWrite); err != nil {
		return pkgerrors.NewSystemError(pkgerrors.ErrCodeOperationFail, messages.ErrorsComposeWriteFailed, err)
	}

	base.Output.Success(messages.SuccessCreatedComposeFile, docker.DockerComposeFilePath)
	return nil
}

// createComposeOverrideFile creates the user-owned docker-compose.override.yml stub.
// It is never overwritten — if it already exists, this is a no-op.
func (pm *ProjectManager) createComposeOverrideFile(base *base.BaseCommand) error {
	if _, err := os.Stat(docker.DockerComposeOverrideFilePath); err == nil {
		return nil
	}

	const content = "# docker-compose.override.yml\n" +
		"# This file is yours — otto-stack never overwrites it.\n" +
		"#\n" +
		"# Add service customizations here. Docker Compose merges this file\n" +
		"# with docker-compose.yml automatically on every `otto-stack up`.\n" +
		"#\n" +
		"# Example:\n" +
		"#   services:\n" +
		"#     postgres:\n" +
		"#       volumes:\n" +
		"#         - /my/local/data:/var/lib/postgresql/data\n" +
		"#       environment:\n" +
		"#         POSTGRES_PASSWORD: localonly\n"

	if err := os.WriteFile(docker.DockerComposeOverrideFilePath, []byte(content), core.PermReadWrite); err != nil {
		return err
	}

	base.Output.Success(messages.SuccessCreatedFile, docker.DockerComposeOverrideFilePath)
	return nil
}

// createGitignoreEntries adds entries to .gitignore
func (pm *ProjectManager) createGitignoreEntries(base *base.BaseCommand) error {
	entries := []string{
		"# " + core.AppNameTitle,
		core.OttoStackDir + "/logs/",
		core.ExtENV + core.LocalFileExtension,
		core.LocalConfigFileName,
		docker.DockerComposeOverrideFileName,
		"*.log",
	}

	content := strings.Join(entries, "\n") + "\n"
	gitignorePath := filepath.Join(core.OttoStackDir, core.GitIgnoreFileName)

	if err := os.WriteFile(gitignorePath, []byte(content), core.PermReadWrite); err != nil {
		return pkgerrors.NewSystemError(pkgerrors.ErrCodeOperationFail, messages.ErrorsFileWriteFailed, err)
	}

	base.Output.Success(messages.SuccessUpdatedFile, gitignorePath)
	return nil
}

// createReadme creates README file
func (pm *ProjectManager) createReadme(projectName string, serviceConfigs []types.ServiceConfig, sharingEnabled bool, base *base.BaseCommand) error {
	const readmeTemplate = `# %s

This project was initialized with %s.

## Services
%s%s

## Commands
- ` + "`%s up`" + ` - Start all services
- ` + "`%s down`" + ` - Stop all services
- ` + "`%s status`" + ` - Show service status
- ` + "`%s logs`" + ` - View service logs
- ` + "`%s validate`" + ` - Validate configuration

## Configuration
- Main config: ` + "`%s/%s`" + `
- Service configs: ` + "`%s/%s/`" + `
%s
## File Management
- **Generated files** (` + "`docker-compose.yml`" + `, ` + "`env.generated`" + `): Automatically regenerated on each ` + "`up`" + ` command
- **Service configs** (` + "`service-configs/`" + `): Created during ` + "`init`" + `, preserved across ` + "`up`" + ` commands for user customization
- **User configs** (` + "`otto-stack-config.yml`" + `): Never overwritten, safe to edit

## Documentation
%s
`

	sharedInfo, sharedSection := pm.buildSharedServicesInfo(serviceConfigs, sharingEnabled)

	readmeContent := fmt.Sprintf(readmeTemplate,
		projectName,
		core.AppNameTitle,
		pm.formatServicesList(svc.ExtractServiceNames(serviceConfigs)),
		sharedInfo,
		core.AppName, core.AppName, core.AppName, core.AppName, core.AppName,
		core.OttoStackDir, core.ConfigFileName,
		core.OttoStackDir, core.ServiceConfigsDir,
		sharedSection,
		core.DocsURL,
	)

	readmePath := filepath.Join(core.OttoStackDir, core.ReadmeFileName)
	if err := os.WriteFile(readmePath, []byte(readmeContent), core.PermReadWrite); err != nil {
		return pkgerrors.NewSystemError(pkgerrors.ErrCodeOperationFail, messages.ErrorsFileWriteFailed, err)
	}

	base.Output.Success(messages.SuccessCreatedReadme, readmePath)
	return nil
}

// buildSharedServicesInfo builds shared services information for README
func (pm *ProjectManager) buildSharedServicesInfo(serviceConfigs []types.ServiceConfig, sharingEnabled bool) (string, string) {
	if !sharingEnabled {
		return "", ""
	}

	var sharedServices []string

	for _, config := range serviceConfigs {
		if config.Shareable {
			sharedServices = append(sharedServices, config.Name)
		}
	}

	if len(sharedServices) == 0 {
		return "", ""
	}

	homeDir, _ := os.UserHomeDir()
	sharedPath := filepath.Join(homeDir, core.OttoStackDir, core.SharedDir)

	info := fmt.Sprintf("\n### Shared Services\nThe following services are shared across projects:\n%s", pm.formatServicesList(sharedServices))

	section := fmt.Sprintf("\n## Shared Services\nShared services are managed globally and located at:\n- `%s/`\n- Registry: `%s/containers.yaml`\n- Compose: `%s/generated/docker-compose.yml`\n",
		sharedPath,
		sharedPath,
		sharedPath,
	)

	return info, section
}

// generateSharedCompose generates shared service files by delegating to GenerateSharedFiles.
func (pm *ProjectManager) generateSharedCompose(serviceConfigs []types.ServiceConfig, sharedRoot string, base *base.BaseCommand) error {
	return GenerateSharedFiles(serviceConfigs, sharedRoot, base)
}

// GenerateSharedFiles creates or refreshes all shared service files under sharedRoot:
//   - generated/docker-compose.yml — the compose file for shared containers
//   - generated/.env.generated    — generated environment variables
//   - services/<name>.yml         — per-service configuration stubs
//
// Only shareable services are processed; non-shareable entries are ignored.
// Idempotent: safe to call on every `up`.
func GenerateSharedFiles(serviceConfigs []types.ServiceConfig, sharedRoot string, base *base.BaseCommand) error {
	var sharedConfigs []types.ServiceConfig
	for _, cfg := range serviceConfigs {
		if cfg.Shareable {
			sharedConfigs = append(sharedConfigs, cfg)
		}
	}

	if len(sharedConfigs) == 0 {
		return nil
	}

	if err := os.MkdirAll(sharedRoot, core.PermReadWriteExec); err != nil {
		return err
	}

	generatedDir := filepath.Join(sharedRoot, core.GeneratedDir)
	if err := os.MkdirAll(generatedDir, core.PermReadWriteExec); err != nil {
		return err
	}

	generator, err := compose.NewGenerator("shared")
	if err != nil {
		return err
	}

	content, err := generator.BuildComposeDataWithHeader(sharedConfigs, core.ComposeHeaderShared)
	if err != nil {
		return err
	}

	composePath := filepath.Join(generatedDir, docker.DockerComposeFileName)
	if err := os.WriteFile(composePath, content, core.PermReadWrite); err != nil {
		return err
	}
	base.Output.Success(messages.SuccessCreatedSharedComposeFile, composePath)

	envPath := filepath.Join(generatedDir, core.EnvGeneratedFileName)
	if err := env.GenerateFile("shared", sharedConfigs, envPath); err != nil {
		return err
	}
	base.Output.Success(messages.SuccessCreatedEnvFile, envPath)

	servicesDir := filepath.Join(sharedRoot, core.ServiceConfigsDir)
	if err := os.MkdirAll(servicesDir, core.PermReadWriteExec); err != nil {
		return err
	}
	cm := NewConfigManager()
	for _, svc := range sharedConfigs {
		stubPath := filepath.Join(servicesDir, svc.Name+core.YMLFileExtension)
		stubContent := cm.generateServiceConfigContent(svc.Name)
		if err := os.WriteFile(stubPath, []byte(stubContent), core.PermReadWrite); err != nil {
			return err
		}
	}

	return nil
}

// FilterProjectServices returns services to include in the project compose file.
// When sharing is enabled, shareable services are excluded — unless the per-service
// sharingServices map explicitly disables sharing for a service (value == false).
// This is called both during init and on every `otto-stack up` to ensure the
// compose file never includes shared services.
func FilterProjectServices(serviceConfigs []types.ServiceConfig, sharingEnabled bool, sharingServices map[string]bool) []types.ServiceConfig {
	if !sharingEnabled {
		return serviceConfigs
	}

	var projectServices []types.ServiceConfig
	for _, config := range serviceConfigs {
		// Non-shareable services always belong in the project compose.
		if !config.Shareable {
			projectServices = append(projectServices, config)
			continue
		}
		// Shareable services are excluded (moved to shared compose) unless the
		// per-service map explicitly opts them out of sharing.
		if len(sharingServices) > 0 && !sharingServices[config.Name] {
			projectServices = append(projectServices, config)
		}
	}
	return projectServices
}

// filterProjectServices is an internal wrapper that calls FilterProjectServices
// using the SharingSpec from the init context.
func (pm *ProjectManager) filterProjectServices(serviceConfigs []types.ServiceConfig, sharing *clicontext.SharingSpec) []types.ServiceConfig {
	if sharing == nil {
		return FilterProjectServices(serviceConfigs, false, nil)
	}
	return FilterProjectServices(serviceConfigs, sharing.Enabled, sharing.Services)
}

// formatServicesList formats services for README
func (pm *ProjectManager) formatServicesList(services []string) string {
	var result strings.Builder
	for _, service := range services {
		fmt.Fprintf(&result, "- %s\n", service)
	}
	return result.String()
}
