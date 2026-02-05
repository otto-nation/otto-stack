package compose

import (
	"fmt"
	"log/slog"
	"maps"
	"os"
	"path/filepath"
	"strings"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/filesystem"
	"github.com/otto-nation/otto-stack/internal/pkg/logger"
	"github.com/otto-nation/otto-stack/internal/pkg/messages"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"gopkg.in/yaml.v3"
)

const expectedEnvParts = 2

// Generator handles docker-compose file generation
type Generator struct {
	projectName string
	logger      *slog.Logger
}

// NewGenerator creates a new compose generator
func NewGenerator(projectName string) (*Generator, error) {
	return &Generator{
		projectName: projectName,
		logger:      logger.GetLogger(),
	}, nil
}

// buildComposeStructure creates the compose structure from ServiceConfigs
func (g *Generator) buildComposeStructure(serviceConfigs []types.ServiceConfig) (map[string]any, error) {
	if g.projectName == "" {
		return nil, pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, "input", messages.ValidationProjectNameEmpty, nil)
	}

	services, err := g.buildServicesFromConfigs(serviceConfigs)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		docker.ComposeFieldServices: services,
		docker.ComposeFieldNetworks: map[string]any{
			docker.DefaultNetworkName: map[string]any{
				docker.ComposeFieldName:   g.projectName + docker.NetworkNameSuffix,
				docker.ComposeFieldLabels: g.buildOttoLabels("network"),
			},
		},
	}, nil
}

// buildServicesFromConfigs creates the services section from ServiceConfigs
func (g *Generator) buildServicesFromConfigs(serviceConfigs []types.ServiceConfig) (map[string]any, error) {
	serviceList := make(map[string]any)
	processedServices := make(map[string]bool)

	// Create a map for quick lookup of ServiceConfigs by name
	configMap := make(map[string]*types.ServiceConfig)
	for i := range serviceConfigs {
		configMap[serviceConfigs[i].Name] = &serviceConfigs[i]
	}

	for _, config := range serviceConfigs {
		if err := g.processServiceConfigAndDependencies(&config, configMap, serviceList, processedServices); err != nil {
			return nil, err
		}
	}

	return serviceList, nil
}

// processServiceConfigAndDependencies processes a service config and its dependencies
func (g *Generator) processServiceConfigAndDependencies(config *types.ServiceConfig, configMap map[string]*types.ServiceConfig, serviceList map[string]any, processed map[string]bool) error {
	if processed[config.Name] {
		return nil
	}
	processed[config.Name] = true

	// Process dependencies first (only if they exist in our configMap)
	for _, dep := range config.Service.Dependencies.Required {
		if depConfig, exists := configMap[dep]; exists {
			if err := g.processServiceConfigAndDependencies(depConfig, configMap, serviceList, processed); err != nil {
				return err
			}
		}
	}

	// Skip configuration services (they don't generate containers)
	if config.ServiceType == types.ServiceTypeConfiguration {
		return nil
	}

	// Build the service configuration using existing logic
	serviceConfig := g.buildService(config)
	if serviceConfig != nil {
		serviceList[config.Name] = serviceConfig
	}

	return nil
}

func (g *Generator) buildService(config *types.ServiceConfig) map[string]any {
	if config.Container.Image == "" {
		return nil
	}

	service := g.createBaseService(config)
	g.addServicePorts(service, config)
	g.addServiceEnvironment(service, config)
	g.addServiceVolumes(service, config)
	g.addServiceConfiguration(service, config)
	g.addServiceHealthCheck(service, config)
	g.addServiceLabels(service, config)

	return service
}

// createBaseService creates the base service configuration
func (g *Generator) createBaseService(config *types.ServiceConfig) map[string]any {
	service := map[string]any{
		docker.ComposeFieldImage: config.Container.Image,
	}

	if len(config.Container.Entrypoint) > 0 {
		service[docker.ComposeFieldEntrypoint] = config.Container.Entrypoint
	}

	return service
}

// addServicePorts adds port configuration to the service
func (g *Generator) addServicePorts(service map[string]any, config *types.ServiceConfig) {
	if len(config.Container.Ports) == 0 {
		return
	}

	var ports []string
	for _, port := range config.Container.Ports {
		portStr := fmt.Sprintf("%s:%s", g.resolveEnvVar(port.External), port.Internal)
		if port.Protocol != "" && port.Protocol != "tcp" {
			portStr += docker.ProtocolSeparator + port.Protocol
		}
		ports = append(ports, portStr)
	}
	service[docker.ComposeFieldPorts] = ports
}

// resolveEnvVar resolves environment variables in the format ${VAR:-default}
func (g *Generator) resolveEnvVar(value string) string {
	if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
		inner := value[2 : len(value)-1]
		if parts := strings.Split(inner, ":-"); len(parts) == expectedEnvParts {
			envVar := parts[0]
			defaultValue := parts[1]

			// Check if environment variable is set
			if envValue := os.Getenv(envVar); envValue != "" {
				return envValue
			}
			return defaultValue // Return default value if env var not set
		}
	}
	return value
}

// addServiceEnvironment adds environment variables to the service
func (g *Generator) addServiceEnvironment(service map[string]any, config *types.ServiceConfig) {
	envVars := g.mergeEnvironmentVariables(config)
	if len(envVars) > 0 {
		service[docker.ComposeFieldEnvironment] = envVars
	}
}

// mergeEnvironmentVariables merges top-level and container-level environment variables
func (g *Generator) mergeEnvironmentVariables(config *types.ServiceConfig) map[string]string {
	envVars := make(map[string]string)

	// Add top-level environment variables first
	maps.Copy(envVars, config.Environment)

	// Add container-level environment variables (these take precedence)
	maps.Copy(envVars, config.Container.Environment)

	return envVars
}

// addServiceVolumes adds volume configuration to the service
func (g *Generator) addServiceVolumes(service map[string]any, config *types.ServiceConfig) {
	if len(config.Container.Volumes) == 0 {
		return
	}

	var volumes []string
	for _, vol := range config.Container.Volumes {
		volStr := fmt.Sprintf("%s:%s", vol.Name, vol.Mount)
		if vol.ReadOnly {
			volStr += docker.VolumeReadOnlySuffix
		}
		volumes = append(volumes, volStr)
	}
	service[docker.ComposeFieldVolumes] = volumes
}

// addServiceConfiguration adds basic service configuration options
func (g *Generator) addServiceConfiguration(service map[string]any, config *types.ServiceConfig) {
	if config.Container.Restart != "" {
		service[docker.ComposeFieldRestart] = string(config.Container.Restart)
	}

	if len(config.Container.Command) > 0 {
		service[docker.ComposeFieldCommand] = config.Container.Command
	}

	if config.Container.MemoryLimit != "" {
		service[docker.ComposeFieldMemLimit] = config.Container.MemoryLimit
	}
}

// addServiceHealthCheck adds health check configuration to the service
func (g *Generator) addServiceHealthCheck(service map[string]any, config *types.ServiceConfig) {
	if config.Container.HealthCheck == nil {
		g.logger.Debug("No healthcheck configured", "service", config.Name)
		return
	}

	g.logger.Debug("Adding healthcheck", "service", config.Name, "test", config.Container.HealthCheck.Test)

	healthCheck := map[string]any{
		docker.HealthCheckFieldTest: config.Container.HealthCheck.Test,
	}

	g.addHealthCheckTiming(healthCheck, config.Container.HealthCheck)
	service[docker.ComposeFieldHealthCheck] = healthCheck
}

// addHealthCheckTiming adds timing configuration to health check
func (g *Generator) addHealthCheckTiming(healthCheck map[string]any, hc *types.HealthCheckSpec) {
	if hc.Interval > 0 {
		healthCheck[docker.HealthCheckFieldInterval] = hc.Interval.String()
	}
	if hc.Timeout > 0 {
		healthCheck[docker.HealthCheckFieldTimeout] = hc.Timeout.String()
	}
	if hc.Retries > 0 {
		healthCheck[docker.HealthCheckFieldRetries] = hc.Retries
	}
	if hc.StartPeriod > 0 {
		healthCheck[docker.HealthCheckFieldStartPeriod] = hc.StartPeriod.String()
	}
}

// addServiceLabels adds Otto Stack labels to the service
func (g *Generator) addServiceLabels(service map[string]any, config *types.ServiceConfig) {
	service[docker.ComposeFieldLabels] = g.buildOttoLabels(config.Name)
}

// buildOttoLabels creates Otto Stack management labels
func (g *Generator) buildOttoLabels(serviceName string) map[string]string {
	return map[string]string{
		docker.LabelOttoManaged:     "true",
		docker.LabelOttoProject:     g.projectName,
		docker.LabelOttoService:     serviceName,
		docker.LabelOttoVersion:     "dev",
		docker.LabelOttoSharingMode: "isolated",
	}
}

// BuildComposeData generates docker-compose YAML content from ServiceConfigs without writing to disk
func (g *Generator) BuildComposeData(serviceConfigs []types.ServiceConfig) ([]byte, error) {
	composeData, err := g.buildComposeStructure(serviceConfigs)
	if err != nil {
		return nil, err
	}

	composeContent, err := yaml.Marshal(composeData)
	if err != nil {
		return nil, pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, "compose", messages.ErrorsComposeMarshalFailed, err)
	}
	return composeContent, nil
}

// WriteComposeFile writes docker-compose content to the specified directory
func (g *Generator) WriteComposeFile(content []byte, outputDir string) error {
	if err := filesystem.EnsureDir(outputDir); err != nil {
		return pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, "compose", messages.ErrorsComposeDirCreateFailed, err)
	}

	filePath := outputDir + "/docker-compose.yml"
	if err := filesystem.WriteFile(filePath, content, core.PermReadWrite); err != nil {
		return pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, "compose", messages.ErrorsComposeWriteFailed, err)
	}

	return nil
}

// GenerateFromServiceConfigs creates a docker-compose file from ServiceConfigs
func (g *Generator) GenerateFromServiceConfigs(serviceConfigs []types.ServiceConfig, projectName string) error {
	content, err := g.BuildComposeData(serviceConfigs)
	if err != nil {
		return err
	}

	return g.WriteComposeFile(content, filepath.Dir(docker.DockerComposeFilePath))
}
