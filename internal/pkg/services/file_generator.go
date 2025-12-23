package services

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"strings"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"gopkg.in/yaml.v3"
)

// FileGenerator handles generation of Docker Compose and environment files
type FileGenerator struct{}

// NewFileGenerator creates a new file generator
func NewFileGenerator() *FileGenerator {
	return &FileGenerator{}
}

// GenerateComposeFile creates and writes the docker-compose file
func (fg *FileGenerator) GenerateComposeFile(serviceNames []string, projectName string) error {
	return fg.GenerateComposeFileWithOriginal(serviceNames, serviceNames, projectName)
}

// GenerateComposeFileWithOriginal creates and writes the docker-compose file with original service names for optimization
func (fg *FileGenerator) GenerateComposeFileWithOriginal(serviceNames []string, originalServices []string, projectName string) error {
	if len(serviceNames) == 0 {
		return pkgerrors.NewValidationError(core.FlagServices, "no services provided for compose generation", nil)
	}

	// Create services manager
	manager, err := New()
	if err != nil {
		return pkgerrors.NewServiceError("compose", "create_manager", err)
	}

	// Generate compose structure
	composeData := fg.buildComposeStructure(serviceNames, originalServices, projectName, manager)

	// Marshal to YAML
	composeContent, err := yaml.Marshal(composeData)
	if err != nil {
		return pkgerrors.NewServiceError("compose", "marshal_yaml", err)
	}

	// Ensure directory exists
	const dirPerm = 0755
	composePath := docker.DockerComposeFilePath
	if err := os.MkdirAll(filepath.Dir(composePath), dirPerm); err != nil {
		return pkgerrors.NewServiceError("compose", "create_directory", err)
	}

	// Write compose file
	const filePerm = 0644
	if err := os.WriteFile(composePath, composeContent, filePerm); err != nil {
		return pkgerrors.NewServiceError("compose", "write_file", err)
	}

	return nil
}

// buildComposeStructure creates the compose structure
func (fg *FileGenerator) buildComposeStructure(serviceNames []string, originalServices []string, projectName string, manager *Manager) map[string]any {
	// Resolve service dependencies
	resolvedServices, _ := manager.ResolveServices(serviceNames)

	return map[string]any{
		docker.ComposeFieldServices: fg.buildServices(resolvedServices, originalServices, manager),
		docker.ComposeFieldNetworks: map[string]any{
			"default": map[string]any{
				docker.ComposeFieldName: projectName + "-network",
			},
		},
	}
}

// buildServices creates the services section
func (fg *FileGenerator) buildServices(resolvedServices []string, originalServices []string, manager *Manager) map[string]any {
	services := make(map[string]any)
	localstackServices := fg.detectLocalStackServices(originalServices)

	for _, serviceName := range resolvedServices {
		serviceConfig, err := manager.GetService(serviceName)
		if err != nil {
			continue
		}

		if fg.shouldSkipService(serviceConfig) {
			continue
		}

		serviceMap := fg.buildServiceConfig(serviceConfig, serviceName)
		fg.applyServiceOptimizations(serviceMap, serviceName, localstackServices)
		services[serviceName] = serviceMap
	}

	return services
}

// shouldSkipService determines if a service should be skipped during compose generation
func (fg *FileGenerator) shouldSkipService(config *ServiceConfig) bool {
	return config.ServiceType == ServiceTypeConfiguration || config.Container.Image == ""
}

// applyServiceOptimizations applies service-specific optimizations
func (fg *FileGenerator) applyServiceOptimizations(serviceMap map[string]any, serviceName string, localstackServices []string) {
	if serviceName == ServiceLocalstack && len(localstackServices) > 0 {
		fg.addLocalStackServicesEnv(serviceMap, localstackServices)
	}
}

// buildServiceConfig builds configuration for a single service
func (fg *FileGenerator) buildServiceConfig(config *ServiceConfig, serviceName string) map[string]any {
	service := map[string]any{
		"image": config.Container.Image,
	}

	// Add ports
	if len(config.Container.Ports) > 0 {
		var ports []string
		for _, port := range config.Container.Ports {
			portStr := fmt.Sprintf("%s:%s", fg.resolveEnvVar(port.External), port.Internal)
			if port.Protocol != "" && port.Protocol != "tcp" {
				portStr += "/" + port.Protocol
			}
			ports = append(ports, portStr)
		}
		service["ports"] = ports
	}

	// Add environment variables
	envVars := make(map[string]string)
	// Add top-level environment variables first
	maps.Copy(envVars, config.Environment)
	// Add container-level environment variables (these take precedence)
	maps.Copy(envVars, config.Container.Environment)
	if len(envVars) > 0 {
		service[docker.ComposeFieldEnvironment] = envVars
	}

	// Add restart policy
	if config.Container.Restart != "" {
		service["restart"] = string(config.Container.Restart)
	}

	// Add command
	if len(config.Container.Command) > 0 {
		service["command"] = config.Container.Command
	}

	// Add memory limit
	if config.Container.MemoryLimit != "" {
		service["mem_limit"] = config.Container.MemoryLimit
	}

	// Add labels
	service["labels"] = map[string]string{
		"otto.managed":      "true",
		"otto.project":      "dev", // This would come from project name
		"otto.service":      serviceName,
		"otto.version":      "dev",
		"otto.sharing-mode": "isolated",
	}

	return service
}

// detectLocalStackServices finds LocalStack services in the service list
func (fg *FileGenerator) detectLocalStackServices(serviceNames []string) []string {
	var localstackServices []string

	for _, serviceName := range serviceNames {
		// Match pattern: localstack-<service-name>
		if awsService, found := strings.CutPrefix(serviceName, ServiceLocalstack+"-"); found {
			localstackServices = append(localstackServices, awsService)
		}
	}

	return localstackServices
}

// addLocalStackServicesEnv adds SERVICES environment variable to LocalStack
func (fg *FileGenerator) addLocalStackServicesEnv(serviceMap map[string]any, services []string) {
	var envMap map[string]string

	if env, exists := serviceMap[docker.ComposeFieldEnvironment]; exists {
		if existingEnv, ok := env.(map[string]string); ok {
			envMap = existingEnv
		}
	}

	if envMap == nil {
		envMap = make(map[string]string)
		serviceMap[docker.ComposeFieldEnvironment] = envMap
	}

	envMap["SERVICES"] = strings.Join(services, ",")
}

// resolveEnvVar resolves environment variables in the format ${VAR:-default}
func (fg *FileGenerator) resolveEnvVar(value string) string {
	if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
		inner := value[2 : len(value)-1]
		if parts := strings.Split(inner, ":-"); len(parts) == expectedEnvParts {
			return parts[1] // Return default value
		}
	}
	return value
}

// GenerateEnvFile generates .env.generated file with resolved services
func (fg *FileGenerator) GenerateEnvFile(services []string, projectName string) error {
	if projectName == "" {
		return pkgerrors.NewValidationError(pkgerrors.FieldProjectName, "project name cannot be empty", nil)
	}

	// Create environment content
	envContent := fg.buildEnvContent(services, projectName)

	// Write environment file
	const filePerm = 0644
	envPath := ".env.generated"
	if err := os.WriteFile(envPath, []byte(envContent), filePerm); err != nil {
		return pkgerrors.NewServiceError("env", "write_file", err)
	}

	return nil
}

// buildEnvContent creates the environment file content
func (fg *FileGenerator) buildEnvContent(services []string, projectName string) string {
	var content strings.Builder

	// Add project configuration
	content.WriteString(fmt.Sprintf("# Generated environment file for %s\n", projectName))
	content.WriteString(fmt.Sprintf("PROJECT_NAME=%s\n", projectName))
	content.WriteString(fmt.Sprintf("SERVICES=%s\n", strings.Join(services, ",")))

	// Add service-specific environment variables
	manager, err := New()
	if err != nil {
		// If we can't load services, just return basic content
		return content.String()
	}

	// Add environment variables for each service
	for _, serviceName := range services {
		serviceConfig, err := manager.GetService(serviceName)
		if err != nil {
			continue // Skip services we can't find
		}

		// Add service-specific environment variables
		content.WriteString(fmt.Sprintf("\n# %s service environment\n", strings.ToUpper(serviceName)))
		for key, value := range serviceConfig.Environment {
			content.WriteString(fmt.Sprintf("%s=%s\n", key, value))
		}
	}

	return content.String()
}
