package project

import (
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/messages"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
)

// validateServiceConfigs validates ServiceConfig objects
// This is called before the validation registry runs, so it only does basic checks
func (h *InitHandler) validateServiceConfigs(serviceConfigs []types.ServiceConfig) error {
	// Detailed validation happens in validateServices() in the registry
	// This is just a quick sanity check
	if len(serviceConfigs) == 0 {
		return pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldServiceName, messages.ValidationNoServicesSelected, nil)
	}
	return nil
}
