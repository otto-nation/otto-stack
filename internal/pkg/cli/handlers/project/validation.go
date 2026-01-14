package project

import (
	"github.com/otto-nation/otto-stack/internal/core"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
)

// validateProjectName validates the project name
func (h *InitHandler) validateProjectName(name string) error {
	if name == "" {
		return pkgerrors.NewValidationError(pkgerrors.FieldProjectName, MsgProjectNameEmpty, nil)
	}

	if len(name) < core.MinProjectNameLength {
		return pkgerrors.NewValidationError(pkgerrors.FieldProjectName, MsgProjectNameTooShort, nil)
	}

	if len(name) > core.MaxProjectNameLength {
		return pkgerrors.NewValidationError(pkgerrors.FieldProjectName, MsgProjectNameTooLong, nil)
	}

	// Cannot start with hyphen or underscore
	if name[0] == '-' || name[0] == '_' {
		return pkgerrors.NewValidationError(pkgerrors.FieldProjectName, core.MsgValidation_project_name_invalid_start, nil)
	}

	// Check for valid characters (alphanumeric, hyphens, underscores)
	for _, char := range name {
		if (char < 'a' || char > 'z') && (char < 'A' || char > 'Z') &&
			(char < '0' || char > '9') && char != '-' && char != '_' {
			return pkgerrors.NewValidationError(pkgerrors.FieldProjectName, core.MsgValidation_project_name_invalid_chars, nil)
		}
	}

	return nil
}

// validateServiceConfigs validates ServiceConfig objects
func (h *InitHandler) validateServiceConfigs(serviceConfigs []types.ServiceConfig) error {
	if len(serviceConfigs) == 0 {
		return pkgerrors.NewValidationError(pkgerrors.FieldServiceName, core.MsgValidation_no_services_selected, nil)
	}

	// Check for duplicates by name
	seen := make(map[string]bool)
	for _, serviceConfig := range serviceConfigs {
		if seen[serviceConfig.Name] {
			return pkgerrors.NewValidationErrorf(pkgerrors.FieldServiceName, core.MsgValidation_duplicate_service, serviceConfig.Name)
		}
		seen[serviceConfig.Name] = true
	}

	// ServiceConfigs are already validated when loaded, so no need to re-validate existence
	return nil
}
