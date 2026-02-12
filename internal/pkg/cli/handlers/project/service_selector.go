package project

import (
	"fmt"
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

	selected, err := ss.promptForServiceSelection(serviceOptions)
	if err != nil {
		return nil, err
	}

	if len(selected) == 0 {
		return nil, pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldServiceName, messages.ValidationNoServicesSelected, nil)
	}

	return ss.mapSelectedServices(selected, allServices), nil
}

func (ss *ServiceSelector) buildServiceList(categories map[string][]types.ServiceConfig) ([]types.ServiceConfig, []string) {
	var allServices []types.ServiceConfig
	var serviceOptions []string
	caser := cases.Title(language.English)

	for categoryName, categoryServices := range categories {
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
		Message: "Select services for your project:",
		Options: serviceOptions,
		Help:    "Use space to select, enter to confirm. Services are grouped by category.",
	}

	var selected []string
	if err := survey.AskOne(prompt, &selected); err != nil {
		return nil, err
	}

	return selected, nil
}

func (ss *ServiceSelector) mapSelectedServices(selected []string, allServices []types.ServiceConfig) []types.ServiceConfig {
	var selectedConfigs []types.ServiceConfig
	selectedMap := make(map[string]bool)

	for _, selection := range selected {
		serviceName := ss.extractServiceName(selection)
		if serviceName == "" || selectedMap[serviceName] {
			continue
		}
		selectedMap[serviceName] = true

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
