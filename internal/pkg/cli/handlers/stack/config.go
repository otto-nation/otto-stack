package stack

import (
	"fmt"
	"log/slog"

	"github.com/otto-nation/otto-stack/internal/pkg/utils"
	"gopkg.in/yaml.v3"
)

// loggerAdapter interface for accessing underlying slog.Logger
type loggerAdapter interface {
	SlogLogger() *slog.Logger
}

// ProjectConfig represents the otto-stack project configuration
type ProjectConfig struct {
	Project struct {
		Name        string `yaml:"name"`
		Environment string `yaml:"environment"`
	} `yaml:"project"`
	Stack struct {
		Enabled []string `yaml:"enabled"`
	} `yaml:"stack"`
}

// LoadProjectConfig loads the otto-stack project configuration
func LoadProjectConfig(configPath string) (*ProjectConfig, error) {
	data, err := utils.ReadFileLines(configPath)
	if err != nil {
		return nil, err
	}

	content := ""
	for _, line := range data {
		content += line + "\n"
	}

	var cfg ProjectConfig
	if err := yaml.Unmarshal([]byte(content), &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}
