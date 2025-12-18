package project

import (
	"github.com/otto-nation/otto-stack/internal/core"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
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

// validateServices validates the selected services
func (h *InitHandler) validateServices(serviceNames []string) error {
	if len(serviceNames) == 0 {
		return pkgerrors.NewValidationError(pkgerrors.FieldServiceName, core.MsgValidation_no_services_selected, nil)
	}

	// Check for duplicates
	seen := make(map[string]bool)
	for _, serviceName := range serviceNames {
		if seen[serviceName] {
			return pkgerrors.NewValidationErrorf(pkgerrors.FieldServiceName, core.MsgValidation_duplicate_service, serviceName)
		}
		seen[serviceName] = true
	}

	// Validate each service exists
	serviceUtils := services.NewServiceUtils()
	for _, serviceName := range serviceNames {
		if _, err := serviceUtils.LoadServiceConfig(serviceName); err != nil {
			return pkgerrors.NewValidationError(pkgerrors.FieldServiceName, "invalid service", err)
		}
	}

	return nil
}
