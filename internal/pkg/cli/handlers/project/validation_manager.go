package project

import (
	"regexp"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/messages"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
)

var projectNameRegex = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]*$`)

// ValidationManager handles validation logic
type ValidationManager struct{}

// NewValidationManager creates a new validation manager
func NewValidationManager() *ValidationManager {
	return &ValidationManager{}
}

// RunValidations executes selected validation functions
func (vm *ValidationManager) RunValidations(selectedValidations map[string]bool, handler *InitHandler, serviceConfigs []types.ServiceConfig, base *base.BaseCommand) error {
	// Run validations
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
				return err
			}
		}
	}

	// Run checks (advisory only, never fail)
	for checkKey := range core.ValidationOptions {
		checkFunc, exists := CheckRegistry[checkKey]
		if !exists {
			continue
		}

		isSelected := selectedValidations[checkKey]
		if isSelected {
			checkFunc(handler, serviceConfigs, base)
		}
	}

	return nil
}

// isRequiredValidation checks if a validation is required based on YAML config
func (vm *ValidationManager) isRequiredValidation(key string) bool {
	return core.ValidationRequired[key]
}

// validationError creates a validation error for project name
func (vm *ValidationManager) validationError(message string) error {
	return pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldProjectName, message, nil)
}

// ValidateProjectName validates a project name
func (vm *ValidationManager) ValidateProjectName(name string) error {
	if name == "" {
		return vm.validationError(messages.ValidationProjectNameEmpty)
	}

	if len(name) < core.MinProjectNameLength {
		return vm.validationError(messages.ValidationProjectNameTooShort)
	}

	if len(name) > core.MaxProjectNameLength {
		return vm.validationError(messages.ValidationProjectNameTooLong)
	}

	if name[0] == '-' || name[0] == '_' {
		return vm.validationError(messages.ValidationProjectNameInvalidStart)
	}

	if !projectNameRegex.MatchString(name) {
		return vm.validationError(messages.ValidationProjectNameInvalidChars)
	}

	return nil
}
