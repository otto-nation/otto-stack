package stack

import (
	"context"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/internal/scripts"
	"gopkg.in/yaml.v3"
)

// InitContainerManager handles initialization container discovery and execution
type InitContainerManager struct{}

// NewInitContainerManager creates a new init container manager
func NewInitContainerManager() *InitContainerManager {
	return &InitContainerManager{}
}

// DiscoverAndRun discovers and runs initialization containers for the given services
func (m *InitContainerManager) DiscoverAndRun(ctx context.Context, setup *CoreSetup, resolvedServices []string, base *base.BaseCommand) error {
	initServices := m.discoverInitServices(resolvedServices)
	if len(initServices) == 0 {
		return nil
	}

	base.Output.Info("Running initialization containers: %v", initServices)

	for _, initService := range initServices {
		if err := m.runSingleInitContainer(ctx, setup, initService, base); err != nil {
			return pkgerrors.NewServiceError(ComponentStack, ActionRunInitContainer, err)
		}
	}

	return nil
}

// discoverInitServices auto-discovers init containers based on config file patterns
func (m *InitContainerManager) discoverInitServices(resolvedServices []string) []string {
	var initServices []string

	configDir := filepath.Join(core.OttoStackDir, core.ServiceConfigsDir)

	// Check if directory exists
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		return initServices
	}

	// Walk through config files
	err := filepath.Walk(configDir, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil || info.IsDir() || !core.IsYAMLFile(info.Name()) {
			return nil
		}

		initService := m.processConfigFileForInit(path, info, resolvedServices)
		m.addUniqueInitService(&initServices, initService)

		return nil
	})

	if err != nil {
		return initServices
	}

	return initServices
}

// processConfigFileForInit processes a single config file and returns init service name if found
func (m *InitContainerManager) processConfigFileForInit(path string, info os.FileInfo, resolvedServices []string) string {
	// Extract service name from filename (remove extension)
	serviceName := core.TrimYAMLExt(info.Name())
	parts := strings.Split(serviceName, "-")

	const minPartsOfServiceName = 2
	if len(parts) < minPartsOfServiceName {
		return ""
	}

	return m.findMatchingInitServiceForFile(parts, path, resolvedServices)
}

// findMatchingInitServiceForFile finds init service that matches the config file pattern
func (m *InitContainerManager) findMatchingInitServiceForFile(parts []string, path string, resolvedServices []string) string {
	// Try different combinations of service name parts
	for i := len(parts) - 1; i >= 1; i-- {
		targetService := strings.Join(parts[:i], "-")

		initService := m.checkServiceMatchForFile(targetService, path, resolvedServices)
		if initService != "" {
			return initService
		}
	}
	return ""
}

// checkServiceMatchForFile checks if target service matches any resolved service
func (m *InitContainerManager) checkServiceMatchForFile(targetService, path string, resolvedServices []string) string {
	for _, resolved := range resolvedServices {
		if !strings.HasPrefix(resolved, targetService) && resolved != targetService {
			continue
		}

		if !m.hasValidConfiguration(path, targetService) {
			continue
		}

		return targetService + "-init"
	}
	return ""
}

// hasValidConfiguration checks if a config file has actual resources to create
func (m *InitContainerManager) hasValidConfiguration(configPath, targetService string) bool {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return false
	}

	// Parse YAML to check for actual configuration
	var config map[string]any
	if err := yaml.Unmarshal(data, &config); err != nil {
		return false
	}

	// Check based on service type
	switch targetService {
	case services.ServiceLocalstack:
		return m.hasValidLocalStackConfig(config)
	case services.ServicePostgres:
		return m.hasValidPostgresConfig(config)
	case services.ServiceKafka:
		return m.hasValidKafkaConfig(config)
	}

	return false
}

// isValidArray checks if an any is a non-empty array
func (m *InitContainerManager) isValidArray(value any) bool {
	if arr, ok := value.([]any); ok {
		return len(arr) > 0
	}
	return false
}

// hasValidLocalStackConfig checks for valid LocalStack configuration
func (m *InitContainerManager) hasValidLocalStackConfig(config map[string]any) bool {
	if queues, exists := config["queues"]; exists && m.isValidArray(queues) {
		return true
	}
	if topics, exists := config["topics"]; exists && m.isValidArray(topics) {
		return true
	}
	if buckets, exists := config["buckets"]; exists && m.isValidArray(buckets) {
		return true
	}
	return false
}

// hasValidPostgresConfig checks for valid PostgreSQL configuration
func (m *InitContainerManager) hasValidPostgresConfig(config map[string]any) bool {
	if schemas, exists := config["schemas"]; exists && m.isValidArray(schemas) {
		return true
	}
	if databases, exists := config["databases"]; exists && m.isValidArray(databases) {
		return true
	}
	return false
}

// hasValidKafkaConfig checks for valid Kafka configuration
func (m *InitContainerManager) hasValidKafkaConfig(config map[string]any) bool {
	if topics, exists := config["topics"]; exists && m.isValidArray(topics) {
		return true
	}
	return false
}

// addUniqueInitService adds init service to list if not empty and not already present
func (m *InitContainerManager) addUniqueInitService(initServices *[]string, initService string) {
	if initService != "" && !slices.Contains(*initServices, initService) {
		*initServices = append(*initServices, initService)
	}
}

// runSingleInitContainer runs a single initialization container
func (m *InitContainerManager) runSingleInitContainer(ctx context.Context, setup *CoreSetup, initServiceName string, base *base.BaseCommand) error {
	targetService := strings.TrimSuffix(initServiceName, "-init")

	// Create init container configuration
	initConfig, err := m.createInitContainerConfig(targetService, setup)
	if err != nil {
		return pkgerrors.NewServiceError(ComponentStack, ActionCreateInitConfig, err)
	}

	// Run the init container
	containerName := fmt.Sprintf("%s-%s", setup.Config.Project.Name, initServiceName)

	base.Output.Info("Starting init container: %s", containerName)

	// Run container and wait for completion
	return setup.DockerClient.RunInitContainer(ctx, containerName, initConfig)
}

// CreateInitContainerConfig creates configuration for init containers (public for testing)
func (m *InitContainerManager) CreateInitContainerConfig(targetService string, setup *CoreSetup) (docker.InitContainerConfig, error) {
	return m.createInitContainerConfig(targetService, setup)
}

// createInitContainerConfig creates configuration for init containers
func (m *InitContainerManager) createInitContainerConfig(targetService string, setup *CoreSetup) (docker.InitContainerConfig, error) {
	service, err := m.loadServiceConfig(targetService)
	if err != nil {
		return docker.InitContainerConfig{}, err
	}

	config := m.buildBaseInitConfig(targetService, setup)
	m.applyServiceEnvironment(config, service)
	m.customizeForServiceType(config, targetService, service)

	return *config, nil
}

func (m *InitContainerManager) loadServiceConfig(targetService string) (*services.ServiceConfig, error) {
	manager, err := GetServicesManager()
	if err != nil {
		return nil, pkgerrors.NewServiceError(ComponentServiceManager, ActionCreateManager, err)
	}

	service, err := manager.GetService(targetService)
	if err != nil {
		return nil, pkgerrors.NewServiceError(ComponentServices, ActionLoadServiceConfig, err)
	}
	return service, nil
}

func (m *InitContainerManager) buildBaseInitConfig(targetService string, setup *CoreSetup) *docker.InitContainerConfig {
	cwd, _ := os.Getwd()
	configPath := filepath.Join(cwd, core.OttoStackDir, core.ServiceConfigsDir)
	processedScript := strings.ReplaceAll(scripts.GenericInitScript, "$$", "$")

	return &docker.InitContainerConfig{
		Image:   services.ImageLocalstack,
		Command: []string{"sh", "-c", processedScript},
		Environment: map[string]string{
			services.InitServiceName: targetService,
			services.InitConfigDir:   "/config",
		},
		Volumes: []string{
			fmt.Sprintf("%s:/config", configPath),
		},
		WorkingDir: "/",
		Networks:   []string{setup.Config.Project.Name + services.NetworkNameSuffix},
	}
}

func (m *InitContainerManager) applyServiceEnvironment(config *docker.InitContainerConfig, service *services.ServiceConfig) {
	maps.Copy(config.Environment, service.Environment)

	if service.Service.Connection != nil && service.Service.Connection.DefaultPort > 0 {
		config.Environment[services.InitServiceEndpointURL] = fmt.Sprintf("http://%s:%d",
			config.Environment[services.InitServiceName], service.Service.Connection.DefaultPort)
	}
}

func (m *InitContainerManager) customizeForServiceType(config *docker.InitContainerConfig, targetService string, service *services.ServiceConfig) {
	switch targetService {
	case services.ServiceLocalstack:
		m.configureLocalstack(config)
	case services.ServicePostgres:
		m.configurePostgres(config, targetService, service)
	case services.ServiceKafka:
		m.configureKafka(config, targetService, service)
	}
}

func (m *InitContainerManager) configureLocalstack(config *docker.InitContainerConfig) {
	config.Environment[services.InitServiceEndpointURL] = fmt.Sprintf("http://localhost:%d", services.PortLocalstack)
}

func (m *InitContainerManager) configurePostgres(config *docker.InitContainerConfig, targetService string, service *services.ServiceConfig) {
	if service.Container.Image != "" {
		config.Image = service.Container.Image
	}

	config.Environment[services.EnvPostgresPGHOST] = targetService

	user := config.Environment[services.EnvPostgresPOSTGRES_USER]
	password := config.Environment[services.EnvPostgresPOSTGRES_PASSWORD]
	database := config.Environment[services.EnvPostgresPOSTGRES_DB]

	port := fmt.Sprintf("%d", services.PortPostgres)
	if service.Service.Connection != nil && service.Service.Connection.DefaultPort > 0 {
		port = fmt.Sprintf("%d", service.Service.Connection.DefaultPort)
	}

	config.Environment[services.InitServiceEndpointURL] = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		user, password, targetService, port, database)
}

func (m *InitContainerManager) configureKafka(config *docker.InitContainerConfig, targetService string, service *services.ServiceConfig) {
	if service.Container.Image != "" {
		config.Image = service.Container.Image
	}

	port := fmt.Sprintf("%d", services.PortKafkaBroker)
	if service.Service.Connection != nil && service.Service.Connection.DefaultPort > 0 {
		port = fmt.Sprintf("%d", service.Service.Connection.DefaultPort)
	}

	config.Environment[services.InitServiceEndpointURL] = fmt.Sprintf("%s:%s", targetService, port)
}
