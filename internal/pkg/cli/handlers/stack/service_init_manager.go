package stack

import (
	"bytes"
	"context"
	"fmt"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"gopkg.in/yaml.v3"
)

// ServiceInitManager handles service-specific init container configuration
type ServiceInitManager struct {
	stackService *services.Service
}

// NewServiceInitManager creates a new service init manager
func NewServiceInitManager() (*ServiceInitManager, error) {
	stackService, err := NewStackService(false)
	if err != nil {
		return nil, err
	}

	return &ServiceInitManager{
		stackService: stackService,
	}, nil
}

// RunInitContainers runs init containers for the specified services
func (m *ServiceInitManager) RunInitContainers(ctx context.Context, serviceConfigs map[string]*services.ServiceConfig, projectName string) error {
	fmt.Printf("DEBUG: Processing %d services for init containers\n", len(serviceConfigs))
	for serviceName, service := range serviceConfigs {
		fmt.Printf("DEBUG: Checking service %s, has init service: %v\n", serviceName, service.InitService != nil)

		// Skip if no init service or not enabled
		if service.InitService == nil || !service.InitService.Enabled {
			continue
		}

		// Handle local mode
		if service.InitService.Mode == "local" {
			if err := m.executeLocalInit(serviceName, service, projectName); err != nil {
				return pkgerrors.NewServiceError("stack", "execute local init", err)
			}
			continue
		}

		// Handle container mode
		config := m.buildInitContainerConfig(serviceName, service, projectName)
		for i, script := range service.InitService.Scripts {
			containerName := fmt.Sprintf("%s-%s-init-%d", projectName, serviceName, i)
			scriptConfig := config

			// Process template with service config data
			processedScript, err := m.processTemplate(script.Content, serviceName)
			if err != nil {
				return pkgerrors.NewServiceError("stack", "process template", err)
			}

			// Substitute environment variables in the script
			processedScript = m.substituteEnvVars(processedScript, config.Environment)

			scriptConfig.Command = []string{docker.ShellSh, docker.ShellC, processedScript}

			if err := m.stackService.DockerClient.RunInitContainer(ctx, containerName, scriptConfig); err != nil {
				return pkgerrors.NewServiceError("stack", "run init container", err)
			}
		}
	}
	return nil
}

// executeLocalInit executes init scripts locally using Docker
func (m *ServiceInitManager) executeLocalInit(serviceName string, service *services.ServiceConfig, projectName string) error {
	// Wait for service to be ready
	if err := m.waitForServiceReady(serviceName, projectName); err != nil {
		return fmt.Errorf("service %s not ready: %w", serviceName, err)
	}

	for _, script := range service.InitService.Scripts {
		// Process template with service config data
		processedScript, err := m.processTemplate(script.Content, serviceName)
		if err != nil {
			return pkgerrors.NewServiceError("stack", "process template", err)
		}

		// Build complete environment from service config
		env := make(map[string]string)

		// Add all service environment variables
		maps.Copy(env, service.Environment)

		// Add init service specific environment variables
		maps.Copy(env, service.InitService.Environment)

		// Substitute environment variables
		processedScript = m.substituteEnvVars(processedScript, env)

		fmt.Printf("DEBUG: Executing local init script for %s:\n%s\n", serviceName, processedScript)

		// Execute script locally using shell
		if err := m.executeLocalScript(service.InitService.Image, service.InitService.Environment, processedScript, projectName); err != nil {
			return pkgerrors.NewServiceError("stack", "execute docker command", err)
		}
	}
	return nil
}

// buildInitContainerConfig converts service config to init container config
func (m *ServiceInitManager) buildInitContainerConfig(serviceName string, service *services.ServiceConfig, projectName string) docker.InitContainerConfig {
	image := "alpine:latest" // Default image
	if service.InitService.Image != "" {
		image = service.InitService.Image
	}

	env := make(map[string]string)
	env["SERVICE_NAME"] = serviceName

	// Copy service environment
	maps.Copy(env, service.Environment)

	// Copy init service environment
	if service.InitService != nil {
		maps.Copy(env, service.InitService.Environment)
	}

	return docker.InitContainerConfig{
		Image:       image,
		Environment: env,
		Networks:    []string{projectName + docker.NetworkNameSuffix},
		WorkingDir:  "/",
	}
}

// processTemplate processes Go templates in init container scripts with service config data
func (m *ServiceInitManager) processTemplate(content string, serviceName string) (string, error) {
	// Load service config files and process template
	configData, err := m.loadServiceConfigData(serviceName)
	if err != nil {
		return content, nil // Return original content if no config data
	}

	tmpl, err := template.New("init").Parse(content)
	if err != nil {
		return "", pkgerrors.NewServiceError("stack", "parse template", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, configData); err != nil {
		return "", pkgerrors.NewServiceError("stack", "execute template", err)
	}

	return buf.String(), nil
}

// loadServiceConfigData loads service config data for template processing
func (m *ServiceInitManager) loadServiceConfigData(serviceName string) (map[string]any, error) {
	// Load config files from service-configs directory
	configDir := filepath.Join(core.OttoStackDir, core.ServiceConfigsDir)

	// Find config files that match the service pattern
	pattern := fmt.Sprintf("%s-*.yml", serviceName)
	matches, err := filepath.Glob(filepath.Join(configDir, pattern))
	if err != nil {
		return nil, err
	}

	configData := make(map[string]any)

	for _, match := range matches {
		data, err := os.ReadFile(match)
		if err != nil {
			continue
		}

		var config map[string]any
		if err := yaml.Unmarshal(data, &config); err != nil {
			continue
		}

		// Merge config data
		maps.Copy(configData, config)
	}

	return configData, nil
}

// substituteEnvVars replaces ${VAR} patterns with environment variable values
func (m *ServiceInitManager) substituteEnvVars(script string, env map[string]string) string {
	for key, value := range env {
		script = strings.ReplaceAll(script, "${"+key+"}", value)
	}
	return script
}

// executeLocalScript executes a script using local shell
func (m *ServiceInitManager) executeLocalScript(image string, env map[string]string, script string, projectName string) error {
	// Set environment variables for the shell
	cmd := exec.Command("sh", "-c", script)

	// Add environment variables
	for key, value := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}

	// Add standard environment variables
	cmd.Env = append(cmd.Env, fmt.Sprintf("DOCKER_IMAGE=%s", image))
	cmd.Env = append(cmd.Env, fmt.Sprintf("DOCKER_NETWORK=%s", projectName+docker.NetworkNameSuffix))

	fmt.Printf("DEBUG: Executing shell script:\n%s\n", script)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("shell command failed: %w, output: %s", err, string(output))
	}

	fmt.Printf("Shell output: %s\n", string(output))
	return nil
}

// waitForServiceReady polls until service is healthy
func (m *ServiceInitManager) waitForServiceReady(serviceName string, projectName string) error {
	const maxRetries = 30
	const retryDelay = 2 * time.Second

	for i := range maxRetries {
		// Check if container is healthy
		cmd := exec.Command("docker", "ps", "--filter", fmt.Sprintf("name=%s-%s-1", projectName, serviceName), "--filter", "health=healthy", "--format", "{{.Names}}")
		output, err := cmd.Output()
		if err == nil && strings.TrimSpace(string(output)) != "" {
			fmt.Printf("DEBUG: Service %s is ready\n", serviceName)
			return nil
		}

		fmt.Printf("DEBUG: Waiting for %s to be ready... (%d/%d)\n", serviceName, i+1, maxRetries)
		time.Sleep(retryDelay)
	}
	return fmt.Errorf("service %s did not become ready within timeout", serviceName)
}
