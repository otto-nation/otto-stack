package project

import (
	"github.com/otto-nation/otto-stack/internal/core"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
)

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
