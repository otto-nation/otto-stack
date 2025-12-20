package docker

import (
	"os"

	"gopkg.in/yaml.v3"
)

// ServiceCharacteristicsConfig defines Docker behaviors for service characteristics
type ServiceCharacteristicsConfig struct {
	ServiceCharacteristics map[string]ServiceCharacteristic `yaml:"service_characteristics"`
}

// ServiceCharacteristic defines flags for different Docker operations
type ServiceCharacteristic struct {
	ComposeUpFlags   []string `yaml:"compose_up_flags"`
	ComposeDownFlags []string `yaml:"compose_down_flags"`
	RunFlags         []string `yaml:"run_flags"`
}

// ServiceCharacteristicsResolver resolves Docker flags based on service characteristics
type ServiceCharacteristicsResolver struct {
	config *ServiceCharacteristicsConfig
}

// NewServiceCharacteristicsResolver creates a new service characteristics resolver
func NewServiceCharacteristicsResolver() (*ServiceCharacteristicsResolver, error) {
	config, err := loadServiceCharacteristicsConfig()
	if err != nil {
		return nil, err
	}

	return &ServiceCharacteristicsResolver{
		config: config,
	}, nil
}

// ResolveComposeUpFlags resolves flags for compose up based on service characteristics
func (scr *ServiceCharacteristicsResolver) ResolveComposeUpFlags(characteristics []string) []string {
	flags := []string{}

	// Add characteristic-based flags
	for _, characteristic := range characteristics {
		if serviceChar, exists := scr.config.ServiceCharacteristics[characteristic]; exists {
			flags = append(flags, serviceChar.ComposeUpFlags...)
		}
	}

	return flags
}

// ResolveComposeDownFlags resolves flags for compose down based on service characteristics
func (scr *ServiceCharacteristicsResolver) ResolveComposeDownFlags(characteristics []string) []string {
	flags := []string{}

	// Add characteristic-based flags
	for _, characteristic := range characteristics {
		if serviceChar, exists := scr.config.ServiceCharacteristics[characteristic]; exists {
			flags = append(flags, serviceChar.ComposeDownFlags...)
		}
	}

	return flags
}

// loadServiceCharacteristicsConfig loads the service characteristics configuration
func loadServiceCharacteristicsConfig() (*ServiceCharacteristicsConfig, error) {
	data, err := os.ReadFile(ServiceCharacteristicsConfigPath)
	if err != nil {
		return nil, err
	}

	var config ServiceCharacteristicsConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
