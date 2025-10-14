package init

import (
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/otto-nation/otto-stack/internal/config"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/utils"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/types"
	"github.com/otto-nation/otto-stack/internal/pkg/compose"
	pkgConfig "github.com/otto-nation/otto-stack/internal/pkg/config"
	"github.com/otto-nation/otto-stack/internal/pkg/ui"
)

// generateConfig generates config using code generation
func (h *InitHandler) generateConfig(name, environment string, services []string, validation, advanced map[string]bool) (string, error) {
	return pkgConfig.GenerateConfig(name, environment, services, validation, advanced), nil
}

// generateInitialComposeFiles generates initial compose files during init
func (h *InitHandler) generateInitialComposeFiles(services []string, projectName, environment string, validation, advanced map[string]bool) error {
	// Create a temporary project config structure
	projectConfig := struct {
		Project struct {
			Name        string
			Environment string
		}
		Stack struct {
			Enabled []string
		}
	}{
		Project: struct {
			Name        string
			Environment string
		}{
			Name:        projectName,
			Environment: environment,
		},
		Stack: struct {
			Enabled []string
		}{
			Enabled: services,
		},
	}

	// Generate .env.generated
	if err := h.generateInitEnvFile(services, &projectConfig); err != nil {
		return fmt.Errorf("failed to generate .env file: %w", err)
	}

	// Generate docker-compose.yml
	if err := h.generateInitDockerCompose(services, &projectConfig); err != nil {
		return fmt.Errorf("failed to generate docker-compose.yml: %w", err)
	}

	ui.Success("Generated otto-stack/docker-compose.yml and otto-stack/.env.generated")
	return nil
}

// generateInitEnvFile generates .env.generated during init using template
func (h *InitHandler) generateInitEnvFile(services []string, projectConfig interface{}) error {
	pc := projectConfig.(*struct {
		Project struct {
			Name        string
			Environment string
		}
		Stack struct {
			Enabled []string
		}
	})

	// Load template
	var templateContent []byte
	candidates := []string{
		"internal/config/env.template",
		"config/env.template",
		"otto-stack/env.template",
	}

	if templatePath, err := h.findTemplateFile(candidates, "env template"); err == nil {
		content, err := os.ReadFile(templatePath)
		if err != nil {
			return fmt.Errorf("failed to read env template: %w", err)
		}
		templateContent = content
	} else {
		templateContent = config.EmbeddedEnvTemplate
		if len(templateContent) == 0 {
			return fmt.Errorf("no env template found and no embedded template available")
		}
	}

	// Parse template with custom functions
	tmpl, err := template.New("env").Funcs(template.FuncMap{
		"ToUpper": strings.ToUpper,
	}).Parse(string(templateContent))
	if err != nil {
		return fmt.Errorf("failed to parse env template: %w", err)
	}

	// Prepare template data
	var templateServices []struct {
		Name   string
		Config *types.ServiceConfig
	}

	for _, serviceName := range services {
		serviceConfig, err := utils.NewServiceUtils().LoadServiceConfig(serviceName)
		if err != nil {
			ui.Warning("Failed to load config for %s: %v", serviceName, err)
			continue
		}
		templateServices = append(templateServices, struct {
			Name   string
			Config *types.ServiceConfig
		}{
			Name:   serviceName,
			Config: serviceConfig,
		})
	}

	data := struct {
		ProjectName string
		Environment string
		GeneratedAt string
		Services    []struct {
			Name   string
			Config *types.ServiceConfig
		}
	}{
		ProjectName: pc.Project.Name,
		Environment: pc.Project.Environment,
		GeneratedAt: time.Now().Format(time.RFC1123),
		Services:    templateServices,
	}

	// Execute template
	var result strings.Builder
	if err := tmpl.Execute(&result, data); err != nil {
		return fmt.Errorf("failed to execute env template: %w", err)
	}

	return os.WriteFile("otto-stack/.env.generated", []byte(result.String()), 0644)
}

// generateInitDockerCompose generates docker-compose.yml during init using programmatic generation
func (h *InitHandler) generateInitDockerCompose(services []string, projectConfig interface{}) error {
	pc := projectConfig.(*struct {
		Project struct {
			Name        string
			Environment string
		}
		Stack struct {
			Enabled []string
		}
	})

	ui.Info("Generating docker-compose files...")

	// Create compose generator
	generator, err := compose.NewGenerator(pc.Project.Name, "internal/config/services")
	if err != nil {
		return fmt.Errorf("failed to create compose generator: %w", err)
	}

	// Generate docker-compose.yml
	composeYAML, err := generator.GenerateYAML(services)
	if err != nil {
		return fmt.Errorf("failed to generate docker-compose YAML: %w", err)
	}

	// Write docker-compose.yml
	if err := os.WriteFile("otto-stack/docker-compose.yml", composeYAML, 0644); err != nil {
		return fmt.Errorf("failed to write docker-compose.yml: %w", err)
	}

	return nil
}
