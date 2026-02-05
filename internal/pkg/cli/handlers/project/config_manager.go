package project

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
	"github.com/otto-nation/otto-stack/internal/pkg/config"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/messages"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
)

// ConfigManager handles configuration file operations
type ConfigManager struct{}

// NewConfigManager creates a new configuration manager
func NewConfigManager() *ConfigManager {
	return &ConfigManager{}
}

// CreateConfigFile creates the main otto-stack configuration file
func (cm *ConfigManager) CreateConfigFile(ctx clicontext.Context, base *base.BaseCommand) error {
	configBytes, err := config.GenerateConfig(ctx)
	if err != nil {
		return err
	}

	configPath := core.OttoStackDir + "/" + core.ConfigFileName
	if err := os.WriteFile(configPath, configBytes, core.PermReadWrite); err != nil {
		return pkgerrors.NewConfigError(pkgerrors.ErrCodeOperationFail, configPath, messages.ErrorsConfigWriteFailed, err)
	}

	base.Output.Success("Created configuration file: %s", configPath)
	return nil
}

// GenerateServiceConfigs generates individual service configuration files
func (cm *ConfigManager) GenerateServiceConfigs(serviceConfigs []types.ServiceConfig, base *base.BaseCommand) {
	for _, config := range serviceConfigs {
		if config.Hidden {
			continue
		}
		if err := cm.generateServiceConfig(config.Name); err != nil {
			base.Output.Warning("Failed to generate config for service %s: %v", config.Name, err)
		}
	}
}

// generateServiceConfig creates a configuration file for a specific service
func (cm *ConfigManager) generateServiceConfig(serviceName string) error {
	content := cm.generateServiceConfigContent(serviceName)
	configPath := filepath.Join(core.OttoStackDir, core.ServiceConfigsDir, serviceName+core.YMLFileExtension)
	return os.WriteFile(configPath, []byte(content), core.PermReadWrite)
}

// generateServiceConfigContent generates the YAML content for a service configuration
func (cm *ConfigManager) generateServiceConfigContent(serviceName string) string {
	config := ServiceConfig{
		Name:        serviceName,
		Description: fmt.Sprintf("Configuration for %s service", serviceName),
	}

	data, _ := yaml.Marshal(&config) // Simple struct marshal cannot fail

	// Add comment header
	header := fmt.Sprintf("# Documentation: %s/services/%s\n\n", core.DocsURL, serviceName)
	return header + string(data)
}
