package services

import (
	"bytes"
	"context"
	"fmt"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"text/template"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/logger"
	servicetypes "github.com/otto-nation/otto-stack/internal/pkg/types"
)

func (s *Service) executeLocalInitScripts(ctx context.Context, serviceConfigs []servicetypes.ServiceConfig, projectName string) error {
	s.logger.Debug("Executing local init scripts for services")
	for _, config := range serviceConfigs {
		if s.hasLocalInitScripts(config) {
			if err := s.executeServiceInitScripts(ctx, config, serviceConfigs, projectName); err != nil {
				return err
			}
		}
	}
	return nil
}

// hasInitScripts checks if a service has init scripts enabled
func (s *Service) hasInitScripts(config servicetypes.ServiceConfig) bool {
	hasInit := config.InitService != nil && config.InitService.Enabled

	// Log debug information about init service configuration
	if config.InitService == nil {
		s.logger.Debug("Service has no InitService configuration", "service", config.Name)
	} else {
		s.logger.Debug("Service InitService configuration",
			"service", config.Name,
			"enabled", config.InitService.Enabled,
			"mode", config.InitService.Mode)
	}

	return hasInit
}

// hasLocalInitScripts checks if a service has local init scripts enabled (for backward compatibility)
func (s *Service) hasLocalInitScripts(config servicetypes.ServiceConfig) bool {
	return s.hasInitScripts(config) && config.InitService.Mode == docker.InitServiceModeLocal
}

// executeServiceInitScripts executes all init scripts for a single service
func (s *Service) executeServiceInitScripts(ctx context.Context, config servicetypes.ServiceConfig, allConfigs []servicetypes.ServiceConfig, projectName string) error {
	for _, script := range config.InitService.Scripts {
		// Process template variables in script content
		processedScript, err := s.processScriptTemplate(script.Content, config, allConfigs)
		if err != nil {
			return fmt.Errorf("failed to process template for service %s: %w", config.Name, err)
		}

		// Execute based on mode
		if config.InitService.Mode == docker.InitServiceModeContainer {
			if err := s.executeScriptInContainer(ctx, processedScript, config, projectName); err != nil {
				return err
			}
		} else {
			// local mode
			env := make(map[string]string)
			if config.InitService.Environment != nil {
				maps.Copy(env, config.InitService.Environment)
			}
			env["DOCKER_IMAGE"] = config.InitService.Image
			env["DOCKER_NETWORK"] = projectName + docker.NetworkNameSuffix

			if err := s.executeScript(ctx, processedScript, env, config.Name); err != nil {
				return err
			}
		}
	}
	return nil
}

// processScriptTemplate processes Go template variables in script content
func (s *Service) processScriptTemplate(scriptContent string, config servicetypes.ServiceConfig, allConfigs []servicetypes.ServiceConfig) (string, error) {
	// Create template data by collecting from dependent services
	templateData := s.collectTemplateData(config, allConfigs)

	// Parse and execute template
	tmpl, err := template.New("script").Parse(scriptContent)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, templateData); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// collectTemplateData collects template data from enabled services that depend on this service
func (s *Service) collectTemplateData(config servicetypes.ServiceConfig, allConfigs []servicetypes.ServiceConfig) map[string]any {
	templateData := make(map[string]any)

	// Collect data from services that depend on the current service (config.Name)
	for _, serviceConfig := range allConfigs {
		if s.serviceDependsOn(serviceConfig, config.Name) {
			s.addConfigData(templateData, serviceConfig)
		}
	}

	s.logger.Debug("Collected template data", "data", templateData)

	return templateData
}

// serviceDependsOn checks if serviceConfig depends on the given service name
func (s *Service) serviceDependsOn(serviceConfig servicetypes.ServiceConfig, serviceName string) bool {
	if serviceConfig.Service.Dependencies.Required != nil {
		return slices.Contains(serviceConfig.Service.Dependencies.Required, serviceName)
	}
	return false
}

// addConfigData adds configuration data from a service to template data
func (s *Service) addConfigData(templateData map[string]any, serviceConfig servicetypes.ServiceConfig) {
	v := reflect.ValueOf(serviceConfig)

	for i := 0; i < v.NumField(); i++ {
		s.processServiceField(templateData, v.Field(i))
	}
}

func (s *Service) processServiceField(templateData map[string]any, field reflect.Value) {
	if field.Kind() != reflect.Pointer || field.IsNil() {
		return
	}

	structValue := field.Elem()
	structType := structValue.Type()

	if !strings.HasSuffix(structType.Name(), "Config") {
		return
	}

	s.extractFieldsFromStruct(templateData, structValue, structType)
}

func (s *Service) extractFieldsFromStruct(templateData map[string]any, structValue reflect.Value, structType reflect.Type) {
	for j := 0; j < structValue.NumField(); j++ {
		s.processStructField(templateData, structValue.Field(j), structType.Field(j))
	}
}

func (s *Service) processStructField(templateData map[string]any, structField reflect.Value, structFieldType reflect.StructField) {
	if !s.isPopulatedSlice(structField) {
		return
	}

	fieldName := s.getYAMLFieldName(structFieldType)
	if fieldName != "" {
		templateData[fieldName] = structField.Interface()
	}
}

func (s *Service) isPopulatedSlice(field reflect.Value) bool {
	return field.Kind() == reflect.Slice && field.Len() > 0
}

func (s *Service) getYAMLFieldName(structField reflect.StructField) string {
	yamlTag := structField.Tag.Get(ServiceCatalogYAMLFormat)
	fieldName := strings.Split(yamlTag, ",")[0]
	if fieldName != "-" {
		return fieldName
	}
	return ""
}

// loadAndValidateServiceConfigs loads user service config files and validates them
func (s *Service) loadAndValidateServiceConfigs(serviceConfigs []servicetypes.ServiceConfig) []servicetypes.ServiceConfig {
	enrichedConfigs := make([]servicetypes.ServiceConfig, 0, len(serviceConfigs))

	for _, config := range serviceConfigs {
		configData, err := loadServiceConfigFile(config.Name)
		if err != nil {
			// File doesn't exist - that's OK, not all services need config files
			s.logger.Debug("No config file for service", "service", config.Name)
			enrichedConfigs = append(enrichedConfigs, config)
			continue
		}

		s.logger.Debug("Loaded config file", "service", config.Name, "data", configData)

		// Merge config data into ServiceConfig struct
		enrichedConfig := mergeConfigIntoStruct(config, configData)
		s.logger.Debug("Merged config", "service", config.Name, "enrichedConfig", enrichedConfig)
		enrichedConfigs = append(enrichedConfigs, enrichedConfig)
	}

	return enrichedConfigs
}

// loadServiceConfigFile loads a service config file from .otto-stack/service-configs/
func loadServiceConfigFile(serviceName string) (map[string]any, error) {
	configPath := filepath.Join(core.OttoStackDir, core.ServiceConfigsDir, serviceName+core.YMLFileExtension)

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var configData map[string]any
	if err := yaml.Unmarshal(data, &configData); err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", configPath, err)
	}

	return configData, nil
}

// mergeConfigIntoStruct merges config file data into the ServiceConfig struct using reflection
func mergeConfigIntoStruct(config servicetypes.ServiceConfig, configData map[string]any) servicetypes.ServiceConfig {
	configValue := reflect.ValueOf(&config).Elem()
	configType := configValue.Type()

	for i := 0; i < configValue.NumField(); i++ {
		field := configValue.Field(i)
		fieldType := configType.Field(i)

		if !isConfigStructField(field, fieldType) {
			continue
		}

		mergeFieldFromConfigData(field, configData)
	}

	return config
}

func isConfigStructField(field reflect.Value, fieldType reflect.StructField) bool {
	return field.Kind() == reflect.Ptr &&
		field.Type().Elem().Kind() == reflect.Struct &&
		strings.Contains(fieldType.Type.String(), "Config")
}

func mergeFieldFromConfigData(field reflect.Value, configData map[string]any) {
	structType := field.Type().Elem()

	for j := 0; j < structType.NumField(); j++ {
		structField := structType.Field(j)
		yamlTag := getYAMLFieldNameFromTag(structField)

		if yamlTag == "" {
			continue
		}

		if data, exists := configData[yamlTag]; exists {
			logger.GetLogger().Debug("Assigning field data", "field", yamlTag, "data", data)
			assignConfigFieldData(field, structType, j, data)
			break
		}
	}
}

func getYAMLFieldNameFromTag(structField reflect.StructField) string {
	yamlTag := structField.Tag.Get(ServiceCatalogYAMLFormat)
	if yamlTag == "" || yamlTag == "-" {
		return ""
	}
	fieldName := strings.Split(yamlTag, ",")[0]
	return fieldName
}

func assignConfigFieldData(field reflect.Value, structType reflect.Type, fieldIndex int, data any) {
	if field.IsNil() {
		newStruct := reflect.New(structType)
		field.Set(newStruct)
	}

	structField := field.Elem().Field(fieldIndex)
	if structField.Kind() == reflect.Slice {
		if sliceData, ok := data.([]any); ok {
			structField.Set(reflect.ValueOf(convertToMapSlice(sliceData)))
		}
	}
}

// convertToMapSlice converts []any to []map[string]any
func convertToMapSlice(data []any) []map[string]any {
	result := make([]map[string]any, 0, len(data))
	for _, item := range data {
		if m, ok := item.(map[string]any); ok {
			result = append(result, m)
		}
	}
	return result
}

// executeScript executes a single script with environment variables on the host
func (s *Service) executeScript(ctx context.Context, scriptContent string, env map[string]string, serviceName string) error {
	cmd := exec.CommandContext(ctx, docker.ShellSh, docker.ShellC, scriptContent)

	// Start with parent environment
	cmd.Env = os.Environ()

	// Add/override with service environment variables
	for key, value := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to execute init script for service %s: %w", serviceName, err)
	}
	return nil
}

// executeScriptInContainer executes a script inside a Docker container
func (s *Service) executeScriptInContainer(ctx context.Context, scriptContent string, config servicetypes.ServiceConfig, projectName string) error {
	// Build init container config
	initConfig := docker.InitContainerConfig{
		Image:       config.InitService.Image,
		Command:     []string{docker.ShellSh, docker.ShellC, scriptContent},
		Environment: config.InitService.Environment,
		Networks:    []string{projectName + docker.NetworkNameSuffix},
	}

	// Use docker client to run init container
	containerName := fmt.Sprintf("%s-init-%d", config.Name, time.Now().Unix())
	if err := s.DockerClient.RunInitContainer(ctx, containerName, initConfig); err != nil {
		return fmt.Errorf("failed to execute init container for service %s: %w", config.Name, err)
	}

	return nil
}
