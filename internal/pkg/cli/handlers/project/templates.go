package project

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/otto-nation/otto-stack/internal/pkg/compose"
	pkgConfig "github.com/otto-nation/otto-stack/internal/pkg/config"
	"github.com/otto-nation/otto-stack/internal/pkg/constants"
	"github.com/otto-nation/otto-stack/internal/pkg/env"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
)

// ProjectConfig represents the project configuration structure
type ProjectConfig struct {
	Project struct {
		Name        string
		Environment string
	}
	Stack struct {
		Enabled []string
	}
}

// generateConfig generates config using code generation
func (h *InitHandler) generateConfig(name string, services []string, validation, advanced map[string]bool) string {
	return pkgConfig.GenerateConfig(name, services, validation, advanced)
}

// generateInitialComposeFiles generates initial compose files during init
func (h *InitHandler) generateInitialComposeFiles(services []string, projectName string, _, _ map[string]bool, base *types.BaseCommand) error {
	projectConfig := &ProjectConfig{
		Project: struct {
			Name        string
			Environment string
		}{
			Name:        projectName,
			Environment: constants.DefaultEnvironment,
		},
		Stack: struct {
			Enabled []string
		}{
			Enabled: services,
		},
	}

	// Generate .env.generated
	if err := h.generateEnvFile(services, projectConfig, base); err != nil {
		return fmt.Errorf("failed to generate .env file: %w", err)
	}

	// Generate docker-compose.yml
	if err := h.generateDockerCompose(services, projectConfig, base); err != nil {
		return fmt.Errorf("failed to generate docker-compose.yml: %w", err)
	}

	base.Output.Success(constants.Messages[constants.MsgSuccess_generated_files],
		filepath.Join(constants.OttoStackDir, constants.DockerComposeFileName),
		filepath.Join(constants.OttoStackDir, constants.EnvGeneratedFileName))

	return nil
}

// generateEnvFile generates .env.generated using programmatic generation
func (h *InitHandler) generateEnvFile(services []string, config *ProjectConfig, base *types.BaseCommand) error {
	base.Output.Info("%s", constants.Messages[constants.MsgProcess_generating_env])

	generator := env.NewGenerator(config.Project.Name, config.Project.Environment)

	envContent, err := generator.Generate(services)
	if err != nil {
		return fmt.Errorf("failed to generate env content: %w", err)
	}

	envPath := filepath.Join(constants.OttoStackDir, constants.EnvGeneratedFileName)
	if err := os.WriteFile(envPath, envContent, constants.FilePermReadWrite); err != nil {
		return fmt.Errorf("failed to write %s: %w", envPath, err)
	}

	return nil
}

// generateDockerCompose generates docker-compose.yml using programmatic generation
func (h *InitHandler) generateDockerCompose(services []string, config *ProjectConfig, base *types.BaseCommand) error {
	base.Output.Info("%s", constants.Messages[constants.MsgProcess_generating_compose])

	generator, err := compose.NewGenerator(config.Project.Name, constants.ServicesDir)
	if err != nil {
		return fmt.Errorf("failed to create compose generator: %w", err)
	}

	composeYAML, err := generator.GenerateYAML(services)
	if err != nil {
		return fmt.Errorf("failed to generate docker-compose YAML: %w", err)
	}

	composePath := filepath.Join(constants.OttoStackDir, constants.DockerComposeFileName)
	if err := os.WriteFile(composePath, composeYAML, constants.FilePermReadWrite); err != nil {
		return fmt.Errorf("failed to write %s: %w", composePath, err)
	}

	return nil
}
