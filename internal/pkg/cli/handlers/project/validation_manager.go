package project

import (
	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
)

// ValidationManager handles validation logic
type ValidationManager struct{}

// NewValidationManager creates a new validation manager
func NewValidationManager() *ValidationManager {
	return &ValidationManager{}
}

// RunValidations executes selected validation functions
func (vm *ValidationManager) RunValidations(selectedValidations map[string]bool, handler *InitHandler, serviceConfigs []services.ServiceConfig, base *base.BaseCommand) error {
	// Always run required validations
	for validationKey := range core.ValidationOptions {
		validationFunc, exists := ValidationRegistry[validationKey]
		if !exists {
			continue
		}

		// Run if it's required OR if user selected it
		isRequired := vm.isRequiredValidation(validationKey)
		isSelected := selectedValidations[validationKey]

		if isRequired || isSelected {
			if err := validationFunc(handler, serviceConfigs, base); err != nil {
				return pkgerrors.NewValidationError(validationKey, "validation failed", err)
			}
		}
	}
	return nil
}

// isRequiredValidation checks if a validation is required based on YAML config
func (vm *ValidationManager) isRequiredValidation(key string) bool {
	return core.ValidationRequired[key]
}
