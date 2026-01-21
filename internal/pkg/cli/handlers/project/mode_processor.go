package project

import (
	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
)

// ModeProcessor handles command processing for different modes (interactive/non-interactive)
type ModeProcessor interface {
	Process(flags *core.InitFlags, base *base.BaseCommand) (clicontext.Context, error)
}

// NewModeProcessor creates the appropriate processor based on mode
func NewModeProcessor(nonInteractive bool, handler *InitHandler) ModeProcessor {
	if nonInteractive {
		return &NonInteractiveProcessor{handler: handler}
	}
	return &InteractiveProcessor{handler: handler}
}
