package config

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
)

// configService implements ConfigService interface
type configService struct{}

// NewConfigService creates a new config service
func NewConfigService() ConfigService {
	return &configService{}
}

// LoadConfig loads the project configuration
func (s *configService) LoadConfig() (*Config, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return nil, pkgerrors.NewConfigError(pkgerrors.ErrCodeOperationFail, "", "failed to load configuration", err)
	}
	return cfg, nil
}

// SaveConfig saves the project configuration
func (s *configService) SaveConfig(cfg *Config) error {
	// For now, return not implemented since there's no SaveConfig in the config package
	return pkgerrors.NewConfigError(pkgerrors.ErrCodeOperationFail, "", "save configuration not implemented", nil)
}

// ValidateConfig validates the configuration
func (s *configService) ValidateConfig(cfg *Config) error {
	if cfg == nil {
		return pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, "config", "configuration is nil", nil)
	}

	if cfg.Project.Name == "" {
		return pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldProjectName, "project name is required", nil)
	}

	if len(cfg.Stack.Enabled) == 0 {
		return pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, "stack.enabled", "at least one service must be enabled", nil)
	}

	return nil
}

// GetConfigHash returns a hash of the current configuration
func (s *configService) GetConfigHash(cfg *Config) (string, error) {
	if cfg == nil {
		return "", pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, "config", "configuration is nil", nil)
	}

	data, err := json.Marshal(cfg)
	if err != nil {
		return "", pkgerrors.NewConfigError(pkgerrors.ErrCodeOperationFail, "", "failed to marshal configuration for hashing", err)
	}

	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash), nil
}
