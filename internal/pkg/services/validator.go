package services

import (
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/messages"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
)

// Validator provides service validation functionality
type Validator struct {
	utils *ServiceUtils
}

// NewValidator creates a new service validator
func NewValidator() *Validator {
	return &Validator{
		utils: NewServiceUtils(),
	}
}

// ValidateServiceNames validates a list of service names
// Checks: existence, not hidden
func (v *Validator) ValidateServiceNames(serviceNames []string) error {
	if len(serviceNames) == 0 {
		return pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldServiceName, messages.ValidationNoServicesSelected, nil)
	}

	for _, name := range serviceNames {
		if err := v.validateServiceName(name); err != nil {
			return err
		}
	}
	return nil
}

// ValidateServiceConfigs validates a list of service configs
// Checks: not empty, no duplicates, each service exists and is loadable
func (v *Validator) ValidateServiceConfigs(serviceConfigs []types.ServiceConfig) error {
	if len(serviceConfigs) == 0 {
		return pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldServiceName, messages.ValidationNoServicesSelected, nil)
	}

	// Check for duplicates
	seen := make(map[string]bool)
	for _, cfg := range serviceConfigs {
		if seen[cfg.Name] {
			return pkgerrors.NewValidationErrorf(pkgerrors.ErrCodeInvalid, pkgerrors.FieldServiceName, messages.ValidationDuplicateService, cfg.Name)
		}
		seen[cfg.Name] = true
	}

	// Validate each service exists and is loadable
	for _, cfg := range serviceConfigs {
		if _, err := v.utils.LoadServiceConfig(cfg.Name); err != nil {
			return pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldServiceName, messages.ValidationInvalidService, err)
		}
	}

	return nil
}

// validateServiceName validates a single service name
func (v *Validator) validateServiceName(name string) error {
	cfg, err := v.utils.LoadServiceConfig(name)
	if err != nil {
		return pkgerrors.NewValidationErrorf(pkgerrors.ErrCodeInvalid, pkgerrors.FieldServiceName, messages.ErrorsServiceUnknown, name)
	}

	if cfg.Hidden {
		return pkgerrors.NewValidationErrorf(pkgerrors.ErrCodeInvalid, pkgerrors.FieldServiceName, messages.ErrorsServiceNotAccessible, name)
	}

	return nil
}
