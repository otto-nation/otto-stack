package services

import (
	"fmt"

	"github.com/otto-nation/otto-stack/internal/config"
	"github.com/otto-nation/otto-stack/internal/pkg/constants"
	"gopkg.in/yaml.v3"
)

// ConnectionConfig defines connection configuration for a service
type ConnectionConfig struct {
	Client      string
	DefaultUser string
	DefaultPort int
	UserFlag    string
	HostFlag    string
	PortFlag    string
	DBFlag      string
	ExtraFlags  []string
}

// ServiceYAML represents the actual YAML structure of service files
type ServiceYAML struct {
	Name        string            `yaml:"name"`
	Category    string            `yaml:"category"`
	Environment map[string]string `yaml:"environment"`
	Connection  struct {
		Client       string   `yaml:"client"`
		DefaultUser  string   `yaml:"default_user"`
		DefaultPort  int      `yaml:"default_port"`
		UserFlag     string   `yaml:"user_flag"`
		HostFlag     string   `yaml:"host_flag"`
		PortFlag     string   `yaml:"port_flag"`
		DatabaseFlag string   `yaml:"database_flag"`
		ExtraFlags   []string `yaml:"extra_flags"`
	} `yaml:"connection"`
	Docker struct {
		Environment []string `yaml:"environment"`
	} `yaml:"docker"`
}

// GetServiceConnectionConfig returns connection config for a service by reading YAML directly
func GetServiceConnectionConfig(serviceName string) (*ConnectionConfig, error) {
	// Load service YAML directly
	serviceYAML, err := loadServiceYAML(serviceName)
	if err != nil {
		return nil, fmt.Errorf("failed to load service %s: %w", serviceName, err)
	}

	// Check if service has connection configuration
	if serviceYAML.Connection.Client == "" {
		return nil, fmt.Errorf("no connection client configured for service: %s", serviceName)
	}

	config := &ConnectionConfig{
		Client:      serviceYAML.Connection.Client,
		DefaultUser: serviceYAML.Connection.DefaultUser,
		DefaultPort: serviceYAML.Connection.DefaultPort,
		UserFlag:    serviceYAML.Connection.UserFlag,
		HostFlag:    serviceYAML.Connection.HostFlag,
		PortFlag:    serviceYAML.Connection.PortFlag,
		DBFlag:      serviceYAML.Connection.DatabaseFlag,
		ExtraFlags:  serviceYAML.Connection.ExtraFlags,
	}

	return config, nil
}

// loadServiceYAML loads a service YAML file directly from embedded FS
func loadServiceYAML(serviceName string) (*ServiceYAML, error) {
	// Dynamically discover categories from services directory structure
	entries, err := config.EmbeddedServicesFS.ReadDir(constants.EmbeddedServicesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read services directory: %w", err)
	}

	// Try each category directory to find the service
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		category := entry.Name()
		servicePath := fmt.Sprintf("%s/%s/%s.yaml", constants.EmbeddedServicesDir, category, serviceName)
		data, err := config.EmbeddedServicesFS.ReadFile(servicePath)
		if err != nil {
			continue // Try next category
		}

		var serviceYAML ServiceYAML
		if err := yaml.Unmarshal(data, &serviceYAML); err != nil {
			return nil, fmt.Errorf("failed to parse service YAML: %w", err)
		}

		return &serviceYAML, nil
	}

	return nil, fmt.Errorf("service not found: %s", serviceName)
}
