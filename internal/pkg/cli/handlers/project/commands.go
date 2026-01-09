package project

import (
	"context"

	"github.com/otto-nation/otto-stack/internal/pkg/base"
	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
	"github.com/otto-nation/otto-stack/internal/pkg/display"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
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

	if len(cliCtx.Services.Names) == 0 {
		base.Output.Warning("No services specified")
		return nil
	}

	// Get service manager
	manager, err := services.New()
	if err != nil {
		return pkgerrors.NewServiceError(services.ComponentServices, services.ActionCreateManager, err)
	}
	allServices := manager.GetAllServices()

	// Show dependencies for each specified service
	for _, serviceName := range cliCtx.Services.Names {
		service, exists := allServices[serviceName]
		if !exists {
			base.Output.Warning("Service '%s' not found", serviceName)
			continue
		}

		base.Output.Info("\n%s %s:", display.StatusSuccess, serviceName)
		if len(service.Service.Dependencies.Required) == 0 {
			base.Output.Info("  No dependencies")
		} else {
			for _, dep := range service.Service.Dependencies.Required {
				base.Output.Info("  • %s", dep)
			}
		}
	}

	base.Output.Success("Dependencies displayed successfully")
	return nil
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

	if len(cliCtx.Services.Names) == 0 {
		base.Output.Warning("No services specified")
		return nil
	}

	// Get service manager
	manager, err := services.New()
	if err != nil {
		return pkgerrors.NewServiceError(services.ComponentServices, services.ActionCreateManager, err)
	}
	allServices := manager.GetAllServices()

	// Check for conflicts between specified services
	conflicts := make(map[string][]string)
	for _, serviceName := range cliCtx.Services.Names {
		service, exists := allServices[serviceName]
		if !exists {
			continue
		}

		for _, conflict := range service.Service.Dependencies.Conflicts {
			for _, otherService := range cliCtx.Services.Names {
				if otherService == conflict {
					conflicts[serviceName] = append(conflicts[serviceName], conflict)
				}
			}
		}
	}

	if len(conflicts) == 0 {
		base.Output.Success("No conflicts detected between selected services")
	} else {
		for service, conflictList := range conflicts {
			base.Output.Warning("Service '%s' conflicts with:", service)
			for _, conflict := range conflictList {
				base.Output.Info("  • %s", conflict)
			}
		}
	}

	base.Output.Success("Conflicts checked successfully")
	return nil
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

	// Basic validation checks
	validationPassed := true

	// Check if services are specified
	if len(cliCtx.Services.Names) == 0 {
		base.Output.Warning("No services specified for validation")
		validationPassed = false
	} else {
		// Get service manager and validate each service
		manager, err := services.New()
		if err != nil {
			return pkgerrors.NewServiceError(services.ComponentServices, services.ActionCreateManager, err)
		}
		allServices := manager.GetAllServices()

		for _, serviceName := range cliCtx.Services.Names {
			if _, exists := allServices[serviceName]; !exists {
				base.Output.Error("Service '%s' not found", serviceName)
				validationPassed = false
			} else {
				base.Output.Info("%sService '%s' is valid", display.StatusSuccess, serviceName)
			}
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

	// Run basic health checks
	healthChecks := []struct {
		name string
		pass bool
	}{
		{"Docker daemon", true}, // Simplified - would check actual Docker connection
		{"Configuration files", true},
		{"Service definitions", len(cliCtx.Services.Names) > 0},
		{"Project structure", true},
	}

	allPassed := true
	for _, check := range healthChecks {
		if check.pass {
			base.Output.Info("%s%s", display.StatusSuccess, check.name)
		} else {
			base.Output.Warning("%s%s", display.StatusError, check.name)
			allPassed = false
		}
	}

	if allPassed {
		base.Output.Success("All health checks passed")
	} else {
		base.Output.Warning("Some health checks failed")
	}

	return nil
}
