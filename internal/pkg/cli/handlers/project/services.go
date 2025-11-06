package project

import (
	"context"
	"fmt"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/display"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/spf13/cobra"
)

// ServicesHandler handles the services command
type ServicesHandler struct{}

// NewServicesHandler creates a new services handler
func NewServicesHandler() *ServicesHandler {
	return &ServicesHandler{}
}

// Handle executes the services command
func (h *ServicesHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
	// Parse all flags with validation - single line!
	flags, err := core.ParseServicesFlags(cmd)
	if err != nil {
		return err
	}

	// Load services by category
	serviceUtils := services.NewServiceUtils()
	categorizedServices, err := serviceUtils.GetServicesByCategory()
	if err != nil {
		return fmt.Errorf(core.MsgErrors_failed_load_services, err)
	}

	// Build service catalog
	catalog := h.buildServiceCatalog(categorizedServices)

	// Create formatter and display
	formatter := display.New(cmd.OutOrStdout(), base.Output)
	options := display.Options{
		Format: flags.Format,
	}

	// Filter by category if specified
	if flags.Category != "" {
		catalog = display.FilterCatalogByCategory(catalog, flags.Category)
	}

	return formatter.FormatServiceCatalog(catalog, options)
}

// buildServiceCatalog converts service data to catalog format
func (h *ServicesHandler) buildServiceCatalog(categorizedServices map[string][]services.ServiceConfigV2) display.ServiceCatalog {
	catalog := display.ServiceCatalog{
		Categories: make(map[string][]display.ServiceInfo),
		Total:      0,
	}

	for categoryName, services := range categorizedServices {
		var catalogServices []display.ServiceInfo
		for _, service := range services {
			catalogServices = append(catalogServices, display.ServiceInfo{
				Name:        service.Name,
				Description: service.Description,
				Category:    categoryName,
			})
		}
		catalog.Categories[categoryName] = catalogServices
		catalog.Total += len(catalogServices)
	}

	return catalog
}

// ValidateArgs validates the command arguments
func (h *ServicesHandler) ValidateArgs(args []string) error {
	return nil
}

// GetRequiredFlags returns required flags for this command
func (h *ServicesHandler) GetRequiredFlags() []string {
	return []string{}
}
