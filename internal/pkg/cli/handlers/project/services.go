package project

import (
	"context"
	"fmt"

	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/utils"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/types"
	"github.com/otto-nation/otto-stack/internal/pkg/constants"
	"github.com/otto-nation/otto-stack/internal/pkg/display"
	"github.com/spf13/cobra"
)

// ServicesHandler handles the services command
type ServicesHandler struct{}

// NewServicesHandler creates a new services handler
func NewServicesHandler() *ServicesHandler {
	return &ServicesHandler{}
}

// Handle executes the services command
func (h *ServicesHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *types.BaseCommand) error {
	// Get flags
	format, _ := cmd.Flags().GetString(constants.FlagFormat)
	category, _ := cmd.Flags().GetString(constants.FlagCategory)

	// Load services by category
	serviceUtils := utils.NewServiceUtils()
	categorizedServices, err := serviceUtils.GetServicesByCategory()
	if err != nil {
		return fmt.Errorf(constants.MsgFailedLoadServices.Content, err)
	}

	// Build service catalog
	catalog := h.buildServiceCatalog(categorizedServices)

	// Create formatter
	formatter, err := display.CreateFormatter(format, cmd.OutOrStdout())
	if err != nil {
		return fmt.Errorf("failed to create formatter: %w", err)
	}

	// Format and display
	options := display.ServiceCatalogOptions{
		Category: category,
		Format:   format,
	}

	return formatter.FormatServiceCatalog(catalog, options)
}

// buildServiceCatalog converts service data to catalog format
func (h *ServicesHandler) buildServiceCatalog(categorizedServices map[string][]types.ServiceInfo) display.ServiceCatalog {
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
