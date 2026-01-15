package project

import (
	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
)

// ValidationManager handles validation logic
type ValidationManager struct{}

// NewValidationManager creates a new validation manager
func NewValidationManager() *ValidationManager {
	return &ValidationManager{}
}

// RunValidations executes selected validation functions
func (vm *ValidationManager) RunValidations(selectedValidations map[string]bool, handler *InitHandler, serviceConfigs []types.ServiceConfig, base *base.BaseCommand) error {
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

// ValidateProjectName validates a project name
func (vm *ValidationManager) ValidateProjectName(name string) error {
	if name == "" {
		return pkgerrors.NewValidationError(pkgerrors.FieldProjectName, MsgProjectNameEmpty, nil)
	}

	if len(name) < core.MinProjectNameLength {
		return pkgerrors.NewValidationError(pkgerrors.FieldProjectName, MsgProjectNameTooShort, nil)
	}

	if len(name) > core.MaxProjectNameLength {
		return pkgerrors.NewValidationError(pkgerrors.FieldProjectName, MsgProjectNameTooLong, nil)
	}

	if name[0] == '-' || name[0] == '_' {
		return pkgerrors.NewValidationError(pkgerrors.FieldProjectName, core.MsgValidation_project_name_invalid_start, nil)
	}

	for _, char := range name {
		if (char < 'a' || char > 'z') && (char < 'A' || char > 'Z') &&
			(char < '0' || char > '9') && char != '-' && char != '_' {
			return pkgerrors.NewValidationError(pkgerrors.FieldProjectName, core.MsgValidation_project_name_invalid_chars, nil)
		}
	}

	return nil
}
