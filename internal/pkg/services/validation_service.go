package services

import (
	"log/slog"

	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/logger"
	"github.com/otto-nation/otto-stack/internal/pkg/messages"
	servicetypes "github.com/otto-nation/otto-stack/internal/pkg/types"
)

// ValidationService provides centralized service validation
type ValidationService struct {
	manager *Manager
	logger  *slog.Logger
}

// NewValidationService creates a new validation service
func NewValidationService(manager *Manager) *ValidationService {
	return &ValidationService{
		manager: manager,
		logger:  logger.GetLogger(),
	}
}

// ValidateUserServices validates user-requested service names (rejects hidden)
func (v *ValidationService) ValidateUserServices(serviceNames []string) error {
	if len(serviceNames) == 0 {
		return pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldServiceName, messages.ValidationNoServicesSelected, nil)
	}

	for _, name := range serviceNames {
		if err := v.ValidateWithContext(name, NewUserValidationContext()); err != nil {
			return err
		}
	}
	return nil
}

// ValidateResolvedServices validates a list of resolved service configs (allows hidden)
func (v *ValidationService) ValidateResolvedServices(serviceConfigs []servicetypes.ServiceConfig) error {
	if len(serviceConfigs) == 0 {
		return pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldServiceName, messages.ValidationNoServicesSelected, nil)
	}

	seen := make(map[string]bool)
	for _, cfg := range serviceConfigs {
		if seen[cfg.Name] {
			return pkgerrors.NewValidationErrorf(pkgerrors.ErrCodeInvalid, pkgerrors.FieldServiceName, messages.ValidationDuplicateService, cfg.Name)
		}
		seen[cfg.Name] = true

		// Use internal context for resolved services (allows hidden)
		if err := v.ValidateWithContext(cfg.Name, NewInternalValidationContext()); err != nil {
			return err
		}
	}

	return nil
}

// ValidateWithContext validates a service name with the given context
func (v *ValidationService) ValidateWithContext(name string, ctx ValidationContext) error {
	service, err := v.manager.GetService(name)
	if err != nil {
		return pkgerrors.NewValidationErrorf(pkgerrors.ErrCodeInvalid, pkgerrors.FieldServiceName, messages.ErrorsServiceUnknown, name)
	}

	// Check if hidden services are allowed in this context
	if service.Hidden && !ctx.AllowHidden {
		v.logger.Debug("Rejecting hidden service in user context", "service", name, "context", ctx)
		return pkgerrors.NewValidationErrorf(pkgerrors.ErrCodeInvalid, pkgerrors.FieldServiceName, messages.ErrorsServiceNotAccessible, name)
	}

	return nil
}
