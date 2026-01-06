package project

import (
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	svc "github.com/otto-nation/otto-stack/internal/pkg/services"
)

// ModeProcessor handles command processing for different modes (interactive/non-interactive)
type ModeProcessor interface {
	Process(flags any, base *base.BaseCommand) (projectName string, originalServiceNames []string, serviceConfigs []svc.ServiceConfig, validation map[string]bool, advanced map[string]bool, error error)
}

// NewModeProcessor creates the appropriate processor based on mode
func NewModeProcessor(nonInteractive bool, handler *InitHandler) ModeProcessor {
	if nonInteractive {
		return &NonInteractiveProcessor{handler: handler}
	}
	return &InteractiveProcessor{handler: handler}
}
