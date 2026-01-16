package project

import (
	"context"

	"github.com/otto-nation/otto-stack/internal/pkg/base"
	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
	"github.com/otto-nation/otto-stack/internal/pkg/display"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/otto-nation/otto-stack/internal/pkg/ui"
)

// ServicesCommand handles listing available services
type ServicesCommand struct{}

// NewServicesCommand creates a new services command
func NewServicesCommand() *ServicesCommand {
	return &ServicesCommand{}
}

// Execute lists available services by category
func (c *ServicesCommand) Execute(ctx context.Context, cliCtx clicontext.Context, base *base.BaseCommand) error {
	base.Output.Header("%s Available Services", ui.IconBox)

	// Get service manager and list services
	manager, err := services.New()
	if err != nil {
		return pkgerrors.NewServiceError(services.ComponentServices, services.ActionCreateManager, err)
	}
	allServices := manager.GetAllServices()

	// Get services by category
	utils := services.NewServiceUtils()
	servicesByCategory, err := utils.GetServicesByCategory()
	if err != nil {
		return pkgerrors.NewServiceError(services.ComponentServices, services.ActionLoadServices, err)
	}

	// Display services by category
	for category, categoryServices := range servicesByCategory {
		base.Output.Info("\n%s %s:", display.StatusSuccess, category)
		for _, service := range categoryServices {
			base.Output.Info("  • %s - %s", service.Name, service.Description)
		}
	}

	base.Output.Info("\nTotal services available: %d", len(allServices))
	base.Output.Success("Services listed successfully")
	return nil
}

// DepsCommand handles showing service dependencies
type DepsCommand struct{}

// NewDepsCommand creates a new deps command
func NewDepsCommand() *DepsCommand {
	return &DepsCommand{}
}

// Execute shows dependencies for a service
func (c *DepsCommand) Execute(ctx context.Context, cliCtx clicontext.Context, base *base.BaseCommand) error {
	base.Output.Header("%s Service Dependencies", ui.IconInfo)

	if !c.validateInput(cliCtx, base) {
		return nil
	}

	allServices, err := c.loadServices()
	if err != nil {
		return err
	}

	c.showAllDependencies(base, cliCtx.Services.Names, allServices)
	base.Output.Success("Dependencies displayed successfully")
	return nil
}

func (c *DepsCommand) validateInput(cliCtx clicontext.Context, base *base.BaseCommand) bool {
	if len(cliCtx.Services.Names) == 0 {
		base.Output.Warning("No services specified")
		return false
	}
	return true
}

func (c *DepsCommand) loadServices() (map[string]*types.ServiceConfig, error) {
	manager, err := services.New()
	if err != nil {
		return nil, pkgerrors.NewServiceError(services.ComponentServices, services.ActionCreateManager, err)
	}

	allServices := manager.GetAllServices()
	result := make(map[string]*types.ServiceConfig, len(allServices))
	for name, service := range allServices {
		svc := service
		result[name] = &svc
	}
	return result, nil
}

func (c *DepsCommand) showAllDependencies(base *base.BaseCommand, names []string, allServices map[string]*types.ServiceConfig) {
	for _, name := range names {
		c.showServiceDeps(base, name, allServices)
	}
}

func (c *DepsCommand) showServiceDeps(base *base.BaseCommand, name string, allServices map[string]*types.ServiceConfig) {
	service, exists := allServices[name]
	if !exists {
		base.Output.Warning("Service '%s' not found", name)
		return
	}

	base.Output.Info("\n%s %s:", display.StatusSuccess, name)
	c.displayDeps(base, service.Service.Dependencies.Required)
}

func (c *DepsCommand) displayDeps(base *base.BaseCommand, deps []string) {
	if len(deps) == 0 {
		base.Output.Info("  No dependencies")
		return
	}
	for _, dep := range deps {
		base.Output.Info("  • %s", dep)
	}
}

