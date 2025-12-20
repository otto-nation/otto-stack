package docker

import (
	"context"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/scripts"
)

// InitContainerManager manages init container operations
type InitContainerManager struct {
	dockerClient *Client
}

// NewInitContainerManager creates a new init container manager
func NewInitContainerManager(dockerClient *Client) *InitContainerManager {
	return &InitContainerManager{
		dockerClient: dockerClient,
	}
}

// RunInitContainer runs an init container for the specified service
func (m *InitContainerManager) RunInitContainer(ctx context.Context, serviceName string, serviceConfig ServiceConfigInterface, projectName string) error {
	initSpec := serviceConfig.GetInitContainerSpec()
	if initSpec == nil || !initSpec.Enabled {
		return nil // No init container needed
	}

	config := m.buildBaseInitConfig(serviceName, serviceConfig, projectName)
	m.applyServiceEnvironment(config, serviceConfig)

	// Execute scripts in order
	for i, script := range initSpec.Scripts {
		if err := m.executeScript(ctx, config, script, i, serviceName, projectName); err != nil {
			return fmt.Errorf("failed to execute script %d (%s): %w", i, script.Type, err)
		}
	}

	return nil
}

// ServiceConfigInterface defines the interface for service configuration
type ServiceConfigInterface interface {
	GetInitContainerImage() string
	GetInitContainerSpec() *InitContainerSpec
	GetEnvironment() map[string]string
	GetConnectionPort() int
	HasConnection() bool
}

// buildBaseInitConfig creates the base configuration for init containers
func (m *InitContainerManager) buildBaseInitConfig(serviceName string, serviceConfig ServiceConfigInterface, projectName string) *InitContainerConfig {
	cwd, _ := os.Getwd()
	configPath := filepath.Join(cwd, core.OttoStackDir, core.ServiceConfigsDir)
	processedScript := strings.ReplaceAll(scripts.GenericInitScript, "$$", "$")

	return &InitContainerConfig{
		Image:   m.getInitContainerImage(serviceConfig),
		Command: []string{ShellSh, ShellC, processedScript},
		Environment: map[string]string{
			InitServiceName: serviceName,
			InitConfigDir:   ContainerConfigPath,
		},
		Volumes: []string{
			fmt.Sprintf("%s:%s", configPath, ContainerConfigPath),
		},
		WorkingDir: ContainerRootPath,
		Networks:   []string{projectName + NetworkNameSuffix},
	}
}

// applyServiceEnvironment applies service environment variables with template resolution
func (m *InitContainerManager) applyServiceEnvironment(config *InitContainerConfig, serviceConfig ServiceConfigInterface) {
	maps.Copy(config.Environment, serviceConfig.GetEnvironment())

	// Resolve template variables to their defaults
	for key, value := range config.Environment {
		config.Environment[key] = m.resolveTemplate(value)
	}

	if serviceConfig.HasConnection() && serviceConfig.GetConnectionPort() > 0 {
		config.Environment[InitServiceEndpointURL] = fmt.Sprintf("http://%s:%d",
			config.Environment[InitServiceName], serviceConfig.GetConnectionPort())
	}
}

// resolveTemplate resolves ${VAR:-default} template patterns
func (m *InitContainerManager) resolveTemplate(value string) string {
	if strings.Contains(value, "${") && strings.Contains(value, ":-") {
		re := regexp.MustCompile(`\$\{[^}]*:-([^}]*)\}`)
		return re.ReplaceAllString(value, "$1")
	}
	return value
}

// getInitContainerImage returns the appropriate init container image
func (m *InitContainerManager) getInitContainerImage(serviceConfig ServiceConfigInterface) string {
	image := serviceConfig.GetInitContainerImage()
	if image != "" {
		return image
	}
	return AlpineLatestImage
}

// executeScript executes a single init script
func (m *InitContainerManager) executeScript(ctx context.Context, config *InitContainerConfig, script InitScript, scriptIndex int, serviceName, projectName string) error {
	containerName := fmt.Sprintf("%s-%s-init-%d", projectName, serviceName, scriptIndex)

	resolvedContent := m.resolveScriptContent(script.Content, serviceName)

	scriptConfig := *config
	scriptConfig.Command = []string{ShellSh, ShellC, resolvedContent}

	return m.dockerClient.RunInitContainer(ctx, containerName, scriptConfig)
}

// resolveScriptContent resolves templates and environment variables in script content
func (m *InitContainerManager) resolveScriptContent(content, serviceName string) string {
	resolved := m.resolveTemplate(content)
	resolved = strings.ReplaceAll(resolved, "${SERVICE_NAME}", serviceName)
	return resolved
}
