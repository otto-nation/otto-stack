package project

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/compose"
	pkgConfig "github.com/otto-nation/otto-stack/internal/pkg/config"
	"github.com/otto-nation/otto-stack/internal/pkg/constants"
	"github.com/otto-nation/otto-stack/internal/pkg/env"
)

// generateConfig generates config using simplified approach
func (h *InitHandler) generateConfig(name string, services []string) string {
	configData, err := pkgConfig.GenerateConfig(name, services)
	if err != nil {
		return "# Error generating config\n"
	}
	return string(configData)
}

// generateInitialComposeFiles generates initial compose files during init
func (h *InitHandler) generateInitialComposeFiles(services []string, projectName string, _, _ map[string]bool, base *base.BaseCommand) error {
	// Generate .env.generated
	if err := h.generateEnvFile(services, projectName, base); err != nil {
		return fmt.Errorf("failed to generate .env file: %w", err)
	}

	// Generate docker-compose.yml
	if err := h.generateDockerCompose(services, projectName, base); err != nil {
		return fmt.Errorf("failed to generate docker-compose.yml: %w", err)
	}

	base.Output.Success(constants.MsgSuccess_generated_files,
		filepath.Join(constants.OttoStackDir, constants.DockerComposeFileName),
		filepath.Join(constants.OttoStackDir, constants.EnvGeneratedFileName))

	return nil
}

// generateEnvFile generates .env.generated using programmatic generation
func (h *InitHandler) generateEnvFile(services []string, projectName string, base *base.BaseCommand) error {
	base.Output.Info("%s", constants.MsgProcess_generating_env)

	envContent, err := env.Generate(projectName, services)
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
func (h *InitHandler) generateDockerCompose(services []string, projectName string, base *base.BaseCommand) error {
	base.Output.Info("%s", constants.MsgProcess_generating_compose)

	generator, err := compose.NewGenerator(projectName, constants.ServicesDir)
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
