package utility

import (
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/project"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/stack"
)

// NewWebInterfacesHandler creates a new web interfaces handler
func NewWebInterfacesHandler() base.CommandHandler {
	return stack.NewWebInterfacesHandler()
}

// NewDoctorHandler creates a new doctor handler
func NewDoctorHandler() base.CommandHandler {
	return project.NewDoctorHandler()
}

// NewVersionHandler creates a new version handler
func NewVersionHandler() base.CommandHandler {
	return project.NewVersionHandler()
}
