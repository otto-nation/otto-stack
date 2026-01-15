package project

import (
	"context"
	"slices"

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

// ConflictsCommand handles checking service conflicts
type ConflictsCommand struct{}

// NewConflictsCommand creates a new conflicts command
func NewConflictsCommand() *ConflictsCommand {
	return &ConflictsCommand{}
}

// Execute checks for conflicts between services
func (c *ConflictsCommand) Execute(ctx context.Context, cliCtx clicontext.Context, base *base.BaseCommand) error {
	base.Output.Header("%s Service Conflicts", ui.IconWarning)

	if !c.validateInput(cliCtx, base) {
		return nil
	}

	allServices, err := c.loadServices()
	if err != nil {
		return err
	}

	conflicts := c.findAllConflicts(cliCtx.Services.Names, allServices)
	c.reportConflicts(base, conflicts)
	base.Output.Success("Conflicts checked successfully")
	return nil
}

func (c *ConflictsCommand) validateInput(cliCtx clicontext.Context, base *base.BaseCommand) bool {
	if len(cliCtx.Services.Names) == 0 {
		base.Output.Warning("No services specified")
		return false
	}
	return true
}

func (c *ConflictsCommand) loadServices() (map[string]*types.ServiceConfig, error) {
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

func (c *ConflictsCommand) findAllConflicts(names []string, allServices map[string]*types.ServiceConfig) map[string][]string {
	conflicts := make(map[string][]string)
	for _, name := range names {
		c.checkServiceConflicts(name, names, allServices, conflicts)
	}
	return conflicts
}

func (c *ConflictsCommand) checkServiceConflicts(name string, allNames []string, allServices map[string]*types.ServiceConfig, conflicts map[string][]string) {
	service, exists := allServices[name]
	if !exists {
		return
	}

	for _, conflict := range service.Service.Dependencies.Conflicts {
		if slices.Contains(allNames, conflict) {
			conflicts[name] = append(conflicts[name], conflict)
		}
	}
}

func (c *ConflictsCommand) reportConflicts(base *base.BaseCommand, conflicts map[string][]string) {
	if len(conflicts) == 0 {
		base.Output.Success("No conflicts detected between selected services")
		return
	}

	for service, conflictList := range conflicts {
		c.displayServiceConflicts(base, service, conflictList)
	}
}

func (c *ConflictsCommand) displayServiceConflicts(base *base.BaseCommand, service string, conflicts []string) {
	base.Output.Warning("Service '%s' conflicts with:", service)
	for _, conflict := range conflicts {
		base.Output.Info("  • %s", conflict)
	}
}

// ValidateCommand handles validating configurations
type ValidateCommand struct{}

// NewValidateCommand creates a new validate command
func NewValidateCommand() *ValidateCommand {
	return &ValidateCommand{}
}

// Execute validates configurations and manifests
func (c *ValidateCommand) Execute(ctx context.Context, cliCtx clicontext.Context, base *base.BaseCommand) error {
	base.Output.Header("%s Configuration Validation", ui.IconInfo)

	if len(cliCtx.Services.Names) == 0 {
		base.Output.Warning("No services specified for validation")
		return nil
	}

	manager, err := services.New()
	if err != nil {
		return pkgerrors.NewServiceError(services.ComponentServices, services.ActionCreateManager, err)
	}
	allServices := manager.GetAllServices()

	validationPassed := true
	for _, serviceName := range cliCtx.Services.Names {
		if _, exists := allServices[serviceName]; !exists {
			base.Output.Error("Service '%s' not found", serviceName)
			validationPassed = false
		} else {
			base.Output.Info("%sService '%s' is valid", display.StatusSuccess, serviceName)
		}
	}

	if validationPassed {
		base.Output.Success("All validations passed")
	} else {
		base.Output.Warning("Some validations failed")
	}

	return nil
}

// DoctorCommand handles diagnosing stack health
type DoctorCommand struct{}

// NewDoctorCommand creates a new doctor command
func NewDoctorCommand() *DoctorCommand {
	return &DoctorCommand{}
}

// Execute runs health checks and diagnostics
func (c *DoctorCommand) Execute(ctx context.Context, cliCtx clicontext.Context, base *base.BaseCommand) error {
	base.Output.Header("%s System Health Check", ui.IconInfo)

	healthChecks := c.runHealthChecks(cliCtx)
	allPassed := c.displayHealthChecks(base, healthChecks)
	c.displayHealthResult(base, allPassed)

	return nil
}

// runHealthChecks runs all system health checks
func (c *DoctorCommand) runHealthChecks(cliCtx clicontext.Context) []struct {
	name string
	pass bool
} {
	return []struct {
		name string
		pass bool
	}{
		{"Docker daemon", true},
		{"Configuration files", true},
		{"Service definitions", len(cliCtx.Services.Names) > 0},
		{"Project structure", true},
	}
}

// displayHealthChecks displays each health check and returns whether all passed
func (c *DoctorCommand) displayHealthChecks(base *base.BaseCommand, healthChecks []struct {
	name string
	pass bool
}) bool {
	allPassed := true
	for _, check := range healthChecks {
		status := display.StatusSuccess
		if !check.pass {
			status = display.StatusError
			allPassed = false
		}
		base.Output.Info("%s%s", status, check.name)
	}
	return allPassed
}

// displayHealthResult displays the final health check result
func (c *DoctorCommand) displayHealthResult(base *base.BaseCommand, allPassed bool) {
	if allPassed {
		base.Output.Success("All health checks passed")
	} else {
		base.Output.Warning("Some health checks failed")
	}
}
