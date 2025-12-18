package project

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/compose"
	pkgConfig "github.com/otto-nation/otto-stack/internal/pkg/config"
	"github.com/otto-nation/otto-stack/internal/pkg/env"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	pkgServices "github.com/otto-nation/otto-stack/internal/pkg/services"
)

// generateConfig generates config using simplified approach
func (h *InitHandler) generateConfig(name string, services []string, validationOptions map[string]bool) string {
	configData, err := pkgConfig.GenerateConfigWithValidation(name, services, validationOptions)
	if err != nil {
		return "# Error generating config\n"
	}
	return string(configData)
}

// generateInitialComposeFiles generates initial compose files during init
func (h *InitHandler) generateInitialComposeFiles(services []string, projectName string, _, _ map[string]bool, base *base.BaseCommand) error {
	// Generate .env.generated
	if err := h.generateEnvFile(services, projectName, base); err != nil {
		return pkgerrors.NewServiceError(ComponentEnv, ActionGenerateEnvFile, err)
	}

	// Generate docker-compose.yml
	if err := h.generateDockerCompose(services, projectName, base); err != nil {
		return pkgerrors.NewServiceError(ComponentCompose, ActionGenerateDockerCompose, err)
	}

	base.Output.Success(core.MsgSuccess_generated_files,
		docker.DockerComposeFilePath,
		core.EnvGeneratedFilePath)

	return nil
}

// generateEnvFile generates .env.generated using programmatic generation
func (h *InitHandler) generateEnvFile(services []string, projectName string, base *base.BaseCommand) error {
	base.Output.Info("%s", core.MsgProcess_generating_env)

	// Resolve services to get actual container services
	manager, err := pkgServices.New()
	if err != nil {
		return pkgerrors.NewServiceError(ComponentServiceManager, ActionCreateManager, err)
	}

	resolvedServices, err := manager.ResolveServices(services)
	if err != nil {
		return pkgerrors.NewServiceError(ComponentServices, ActionResolveServices, err)
	}

	envContent, err := env.Generate(projectName, resolvedServices, manager)
	if err != nil {
		return pkgerrors.NewConfigError(core.EnvGeneratedFilePath, ActionGenerateEnv, err)
	}

	if err := os.MkdirAll(filepath.Dir(core.EnvGeneratedFilePath), core.PermReadWriteExec); err != nil {
		return pkgerrors.NewServiceError(ComponentDirectory, ActionCreateDirectory, err)
	}
	if err := os.WriteFile(core.EnvGeneratedFilePath, envContent, core.PermReadWrite); err != nil {
		return pkgerrors.NewServiceError(ComponentFile, ActionWriteFile, err)
	}

	return nil
}

// generateDockerCompose generates docker-compose.yml using programmatic generation
func (h *InitHandler) generateDockerCompose(services []string, projectName string, base *base.BaseCommand) error {
	base.Output.Info("%s", core.MsgProcess_generating_compose)

	manager, err := pkgServices.New()
	if err != nil {
		return pkgerrors.NewServiceError(ComponentServiceManager, ActionCreateManager, err)
	}

	generator, err := compose.NewGenerator(projectName, pkgServices.ServicesDir, manager)
	if err != nil {
		return pkgerrors.NewServiceError(ComponentCompose, ActionCreateComposeGenerator, err)
	}

	composeYAML, err := generator.GenerateYAML(services)
	if err != nil {
		return pkgerrors.NewServiceError(ComponentCompose, ActionGenerateDockerComposeYAML, err)
	}

	if err := os.MkdirAll(filepath.Dir(docker.DockerComposeFilePath), core.PermReadWriteExec); err != nil {
		return pkgerrors.NewServiceError(ComponentDirectory, ActionCreateDirectory, err)
	}
	if err := os.WriteFile(docker.DockerComposeFilePath, composeYAML, core.PermReadWrite); err != nil {
		return pkgerrors.NewServiceError(ComponentFile, ActionWriteFile, err)
	}

	return nil
}

// generateServiceConfigs creates service configuration files with documentation links
func (h *InitHandler) generateServiceConfigs(services []string, base *base.BaseCommand) {
	base.Output.Info("Generating service configuration files...")

	for _, serviceName := range services {
		if err := h.generateServiceConfig(serviceName); err != nil {
			base.Output.Warning("Failed to generate config for %s: %v", serviceName, err)
		}
	}
}

// generateServiceConfig creates a single service configuration file
func (h *InitHandler) generateServiceConfig(serviceName string) error {
	configContent := h.generateServiceConfigContent(serviceName)

	configPath := filepath.Join(core.OttoStackDir, core.ServiceConfigsDir, serviceName+core.YMLFileExtension)
	return os.WriteFile(configPath, []byte(configContent), core.PermReadWrite)
}

// generateServiceConfigContent generates the content for a service config file
func (h *InitHandler) generateServiceConfigContent(serviceName string) string {
	docURL := fmt.Sprintf(core.DocsURL+"/services/%s/", serviceName)
	title := strings.ReplaceAll(serviceName, "-", " ")

	return fmt.Sprintf(`# %s Configuration
# For all available options, see: %s
#
# This file contains service-specific configuration.
# Uncomment and modify the options you need.

# Add your %s configuration here
`, title, docURL, serviceName)
}
