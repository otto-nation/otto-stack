package project

import (
	"fmt"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
)

// validateProjectName validates the project name
func (h *InitHandler) validateProjectName(name string) error {
	if name == "" {
		return fmt.Errorf("%s", core.MsgValidation_project_name_empty)
	}

	if len(name) < core.MinProjectNameLength {
		return fmt.Errorf("%s", core.MsgValidation_project_name_too_short)
	}

	if len(name) > core.MaxProjectNameLength {
		return fmt.Errorf("%s", core.MsgValidation_project_name_too_long)
	}

	// Cannot start with hyphen or underscore
	if name[0] == '-' || name[0] == '_' {
		return fmt.Errorf("%s", core.MsgValidation_project_name_invalid_start)
	}

	// Check for valid characters (alphanumeric, hyphens, underscores)
	for _, char := range name {
		if (char < 'a' || char > 'z') && (char < 'A' || char > 'Z') &&
			(char < '0' || char > '9') && char != '-' && char != '_' {
			return fmt.Errorf("%s", core.MsgValidation_project_name_invalid_chars)
		}
	}

	return nil
}

// validateServices validates the selected services
func (h *InitHandler) validateServices(serviceNames []string) error {
	if len(serviceNames) == 0 {
		return fmt.Errorf("%s", core.MsgValidation_no_services_selected)
	}

	// Check for duplicates
	seen := make(map[string]bool)
	for _, serviceName := range serviceNames {
		if seen[serviceName] {
			return fmt.Errorf(core.MsgValidation_duplicate_service, serviceName)
		}
		seen[serviceName] = true
	}

	// Validate each service exists
	serviceUtils := services.NewServiceUtils()
	for _, serviceName := range serviceNames {
		if _, err := serviceUtils.LoadServiceConfig(serviceName); err != nil {
			return fmt.Errorf(core.MsgValidation_invalid_service, serviceName, err)
		}
	}

	return nil
}
