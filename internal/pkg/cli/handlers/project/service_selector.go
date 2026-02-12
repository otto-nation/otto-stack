package project

import (
	"fmt"
	"sort"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/messages"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// ServiceSelector handles service selection prompts
type ServiceSelector struct{}

// NewServiceSelector creates a new service selector
func NewServiceSelector() *ServiceSelector {
	return &ServiceSelector{}
}

// SelectServices prompts user to select services and returns ServiceConfigs
func (ss *ServiceSelector) SelectServices() ([]types.ServiceConfig, error) {
	categories, err := ss.loadServiceCategories()
	if err != nil {
		return nil, err
	}

	if len(categories) == 0 {
		return nil, pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldServiceName, messages.ErrorsServiceSelectionFailed, nil)
	}

	return ss.selectServicesFromAllCategories(categories)
}

func (ss *ServiceSelector) loadServiceCategories() (map[string][]types.ServiceConfig, error) {
	utils := services.NewServiceUtils()
	return utils.GetServicesByCategory()
}

func (ss *ServiceSelector) selectServicesFromAllCategories(categories map[string][]types.ServiceConfig) ([]types.ServiceConfig, error) {
	allServices, serviceOptions := ss.buildServiceList(categories)

	selectedNames, err := ss.promptForServiceSelection(serviceOptions)
	if err != nil {
		return nil, err
	}

	if len(selectedNames) == 0 {
		return nil, pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldServiceName, messages.ValidationNoServicesSelected, nil)
	}

	return ss.mapSelectedServicesByName(selectedNames, allServices), nil
}

func (ss *ServiceSelector) buildServiceList(categories map[string][]types.ServiceConfig) ([]types.ServiceConfig, []string) {
	var allServices []types.ServiceConfig
	var serviceOptions []string
	caser := cases.Title(language.English)

	// Sort category names for consistent ordering
	categoryNames := make([]string, 0, len(categories))
	for categoryName := range categories {
		categoryNames = append(categoryNames, categoryName)
	}
	sort.Strings(categoryNames)

	for _, categoryName := range categoryNames {
		categoryServices := categories[categoryName]

		// Sort services within category by name
		sort.Slice(categoryServices, func(i, j int) bool {
			return categoryServices[i].Name < categoryServices[j].Name
		})

		for _, service := range categoryServices {
			displayName := fmt.Sprintf("[%s] %s - %s",
				caser.String(categoryName), service.Name, service.Description)
			serviceOptions = append(serviceOptions, displayName)
			allServices = append(allServices, service)
		}
	}

	return allServices, serviceOptions
}

func (ss *ServiceSelector) promptForServiceSelection(serviceOptions []string) ([]string, error) {
	prompt := &survey.MultiSelect{
		Message: messages.PromptsSelectServices,
		Options: serviceOptions,
		Help:    messages.PromptsSelectServicesHelp,
	}

	var selected []string
	if err := survey.AskOne(prompt, &selected); err != nil {
		return nil, err
	}

	// Extract service names from selected display strings
	serviceNames := make([]string, 0, len(selected))
	for _, displayStr := range selected {
		name := ss.extractServiceName(displayStr)
		if name != "" {
			serviceNames = append(serviceNames, name)
		}
	}

	return serviceNames, nil
}

func (ss *ServiceSelector) mapSelectedServicesByName(serviceNames []string, allServices []types.ServiceConfig) []types.ServiceConfig {
	var selectedConfigs []types.ServiceConfig

	for _, serviceName := range serviceNames {
		if config := ss.findServiceConfig(serviceName, allServices); config != nil {
			selectedConfigs = append(selectedConfigs, *config)
		}
	}

	return selectedConfigs
}

func (ss *ServiceSelector) extractServiceName(selection string) string {
	const minParts = 2
	parts := strings.Split(selection, "] ")
	if len(parts) < minParts {
		return ""
	}

	serviceNamePart := strings.Split(parts[1], " - ")
	if len(serviceNamePart) < 1 {
		return ""
	}

	return serviceNamePart[0]
}

func (ss *ServiceSelector) findServiceConfig(serviceName string, allServices []types.ServiceConfig) *types.ServiceConfig {
	for i := range allServices {
		if allServices[i].Name == serviceName {
			return &allServices[i]
		}
	}
	return nil
}
