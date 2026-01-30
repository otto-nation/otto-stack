package project

import (
	"context"
	"sort"

	"github.com/spf13/cobra"

	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/display"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	servicetypes "github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/otto-nation/otto-stack/internal/pkg/ui"
)

// ServicesHandler handles the services command
type ServicesHandler struct{}

// NewServicesHandler creates a new services handler
func NewServicesHandler() *ServicesHandler {
	return &ServicesHandler{}
}

// Handle executes the services command
func (h *ServicesHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
	base.Output.Header("%s Available Services", ui.IconBox)

	servicesByCategory, totalCount, err := h.loadServices()
	if err != nil {
		return err
	}

	groups := h.buildTableGroups(servicesByCategory)

	headers := []string{display.HeaderService, display.HeaderCategory, display.HeaderDescription}
	display.RenderTableWithSeparators(base.Output.Writer(), headers, groups)

	base.Output.Info("\nTotal services available: %d", totalCount)
	base.Output.Success("Services listed successfully")
	return nil
}

func (h *ServicesHandler) loadServices() (map[string][]servicetypes.ServiceConfig, int, error) {
	manager, err := services.New()
	if err != nil {
		return nil, 0, pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentServices, "operation failed", err)
	}

	utils := services.NewServiceUtils()
	servicesByCategory, err := utils.GetServicesByCategory()
	if err != nil {
		return nil, 0, pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentServices, "operation failed", err)
	}

	return servicesByCategory, len(manager.GetAllServices()), nil
}

func (h *ServicesHandler) buildTableGroups(servicesByCategory map[string][]servicetypes.ServiceConfig) [][][]string {
	categories := h.sortedCategories(servicesByCategory)
	groups := make([][][]string, len(categories))

	for i, category := range categories {
		groups[i] = h.buildCategoryGroup(servicesByCategory[category], category)
	}

	return groups
}

func (h *ServicesHandler) sortedCategories(servicesByCategory map[string][]servicetypes.ServiceConfig) []string {
	categories := make([]string, 0, len(servicesByCategory))
	for category := range servicesByCategory {
		categories = append(categories, category)
	}
	sort.Strings(categories)
	return categories
}

func (h *ServicesHandler) buildCategoryGroup(categoryServices []servicetypes.ServiceConfig, category string) [][]string {
	sort.Slice(categoryServices, func(a, b int) bool {
		return categoryServices[a].Name < categoryServices[b].Name
	})

	group := make([][]string, len(categoryServices))
	for i, svc := range categoryServices {
		group[i] = []string{svc.Name, category, svc.Description}
	}
	return group
}

// ValidateArgs validates the command arguments
func (h *ServicesHandler) ValidateArgs(args []string) error {
	return nil
}

// GetRequiredFlags returns required flags for this command
func (h *ServicesHandler) GetRequiredFlags() []string {
	return []string{}
}
