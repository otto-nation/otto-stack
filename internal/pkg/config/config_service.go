package config

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/messages"
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
		return nil, pkgerrors.NewConfigError(pkgerrors.ErrCodeOperationFail, "", messages.ErrorsConfigLoadFailed, err)
	}
	return cfg, nil
}

// ValidateConfig validates the configuration
func (s *configService) ValidateConfig(cfg *Config) error {
	if cfg == nil {
		return pkgerrors.NewSystemError(pkgerrors.ErrCodeInvalid, messages.ErrorsConfigNil, nil)
	}

	if cfg.Project.Name == "" {
		return pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldProjectName, messages.ValidationProjectNameEmpty, nil)
	}

	if len(cfg.Stack.Enabled) == 0 {
		return pkgerrors.NewSystemError(pkgerrors.ErrCodeInvalid, messages.ValidationNoServicesSelected, nil)
	}

	return nil
}

// GetConfigHash returns a hash of the current configuration
func (s *configService) GetConfigHash(cfg *Config) (string, error) {
	if cfg == nil {
		return "", pkgerrors.NewSystemError(pkgerrors.ErrCodeInvalid, messages.ErrorsConfigNil, nil)
	}

	data, err := json.Marshal(cfg)
	if err != nil {
		return "", pkgerrors.NewConfigError(pkgerrors.ErrCodeOperationFail, "", messages.ErrorsConfigMarshalFailed, err)
	}

	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash), nil
}
