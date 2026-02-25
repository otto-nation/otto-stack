package project

import (
	"context"
	"encoding/json"
	"sort"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/ci"
	"github.com/otto-nation/otto-stack/internal/pkg/display"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/messages"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	servicetypes "github.com/otto-nation/otto-stack/internal/pkg/types"
)

// ServicesHandler handles the services command
type ServicesHandler struct{}

// NewServicesHandler creates a new services handler
func NewServicesHandler() *ServicesHandler {
	return &ServicesHandler{}
}

// Handle executes the services command
func (h *ServicesHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
	flags, err := core.ParseServicesFlags(cmd)
	if err != nil {
		return pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldFlags, messages.ValidationFailedParseFlags, err)
	}

	base.Output.Header(messages.ServicesHeader)

	servicesByCategory, err := h.loadServices()
	if err != nil {
		return err
	}

	if flags.Category != "" {
		servicesByCategory = h.filterByCategory(servicesByCategory, flags.Category)
		if len(servicesByCategory) == 0 {
			base.Output.Info(messages.ServicesNoneInCategory, flags.Category)
			return nil
		}
	}

	totalCount := h.countServices(servicesByCategory)

	switch flags.Format {
	case "json":
		return h.outputJSON(servicesByCategory, totalCount, base)
	case "yaml":
		return h.outputYAML(servicesByCategory, base)
	case "table":
		rows := h.buildFlatRows(servicesByCategory)
		headers := []string{display.HeaderService, display.HeaderCategory, display.HeaderDescription}
		display.RenderTable(base.Output.Writer(), headers, rows)
	default:
		// "group" and any unrecognised value: grouped table with category separators
		groups := h.buildTableGroups(servicesByCategory)
		headers := []string{display.HeaderService, display.HeaderCategory, display.HeaderDescription}
		display.RenderTableWithSeparators(base.Output.Writer(), headers, groups)
	}

	base.Output.Info(messages.ServicesTotalAvailable, totalCount)
	return nil
}

func (h *ServicesHandler) loadServices() (map[string][]servicetypes.ServiceConfig, error) {
	utils := services.NewServiceUtils()
	servicesByCategory, err := utils.GetServicesByCategory()
	if err != nil {
		return nil, pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentServices, messages.ErrorsServiceOperationFailed, err)
	}
	return servicesByCategory, nil
}

func (h *ServicesHandler) filterByCategory(servicesByCategory map[string][]servicetypes.ServiceConfig, category string) map[string][]servicetypes.ServiceConfig {
	filtered := make(map[string][]servicetypes.ServiceConfig)
	for cat, svcs := range servicesByCategory {
		if cat == category {
			filtered[cat] = svcs
		}
	}
	return filtered
}

func (h *ServicesHandler) countServices(servicesByCategory map[string][]servicetypes.ServiceConfig) int {
	total := 0
	for _, svcs := range servicesByCategory {
		total += len(svcs)
	}
	return total
}

func (h *ServicesHandler) outputJSON(servicesByCategory map[string][]servicetypes.ServiceConfig, totalCount int, base *base.BaseCommand) error {
	type serviceEntry struct {
		Name        string `json:"name"`
		Category    string `json:"category"`
		Description string `json:"description"`
	}

	output := ci.ServicesOutput{
		Services: make([]any, 0, totalCount),
		Count:    totalCount,
	}

	for _, cat := range h.sortedCategories(servicesByCategory) {
		for _, svc := range servicesByCategory[cat] {
			output.Services = append(output.Services, serviceEntry{
				Name:        svc.Name,
				Category:    cat,
				Description: svc.Description,
			})
		}
	}

	return json.NewEncoder(base.Output.Writer()).Encode(output)
}

func (h *ServicesHandler) outputYAML(servicesByCategory map[string][]servicetypes.ServiceConfig, base *base.BaseCommand) error {
	type serviceEntry struct {
		Name        string `yaml:"name"`
		Category    string `yaml:"category"`
		Description string `yaml:"description"`
	}

	var entries []serviceEntry
	for _, cat := range h.sortedCategories(servicesByCategory) {
		for _, svc := range servicesByCategory[cat] {
			entries = append(entries, serviceEntry{
				Name:        svc.Name,
				Category:    cat,
				Description: svc.Description,
			})
		}
	}

	return yaml.NewEncoder(base.Output.Writer()).Encode(entries)
}

func (h *ServicesHandler) buildFlatRows(servicesByCategory map[string][]servicetypes.ServiceConfig) [][]string {
	var rows [][]string
	for _, cat := range h.sortedCategories(servicesByCategory) {
		for _, svc := range servicesByCategory[cat] {
			rows = append(rows, []string{svc.Name, cat, svc.Description})
		}
	}
	return rows
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
