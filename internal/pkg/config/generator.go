package config

import (
	"fmt"
	"strings"

	"github.com/otto-nation/otto-stack/internal/pkg/constants"
)

// GenerateConfig generates a otto-stack configuration file
func GenerateConfig(projectName, environment string, services []string, validation, advanced map[string]bool) string {
	var builder strings.Builder

	// Header comment
	builder.WriteString(fmt.Sprintf("# %s Configuration\n", constants.AppNameTitle))
	builder.WriteString(fmt.Sprintf("# Generated for project: %s\n", projectName))
	builder.WriteString(fmt.Sprintf("# Documentation: %s\n\n", constants.ConfigDocsURL))

	// Project section
	builder.WriteString("project:\n")
	builder.WriteString(fmt.Sprintf("  name: %s\n", projectName))
	builder.WriteString(fmt.Sprintf("  environment: %s\n\n", environment))

	// Stack section
	builder.WriteString("stack:\n")
	builder.WriteString("  enabled:\n")
	for _, service := range services {
		builder.WriteString(fmt.Sprintf("    - %s\n", service))
	}
	builder.WriteString("\n")

	// Service configuration section
	builder.WriteString("# Service-specific configuration\n")
	builder.WriteString(fmt.Sprintf("# Service configuration options: %s\n", constants.ServiceConfigURL))
	builder.WriteString("service_configuration: {}\n\n")

	// Validation section
	builder.WriteString("validation:\n")

	validationDefaults := map[string]bool{
		"skip_warnings":            constants.DefaultSkipWarnings,
		"allow_multiple_databases": constants.DefaultAllowMultipleDBs,
	}

	for key, defaultValue := range validationDefaults {
		builder.WriteString(fmt.Sprintf("  %s: %t\n", key, getBoolValue(validation, key, defaultValue)))
	}
	builder.WriteString("\n")

	// Advanced section
	builder.WriteString("advanced:\n")

	advancedDefaults := map[string]bool{
		"auto_start":          constants.DefaultAutoStart,
		"pull_latest_images":  constants.DefaultPullLatestImages,
		"cleanup_on_recreate": constants.DefaultCleanupOnRecreate,
	}

	for key, defaultValue := range advancedDefaults {
		builder.WriteString(fmt.Sprintf("  %s: %t\n", key, getBoolValue(advanced, key, defaultValue)))
	}

	return builder.String()
}

// getBoolValue gets a boolean value from map with fallback to default
func getBoolValue(m map[string]bool, key string, defaultValue bool) bool {
	if value, exists := m[key]; exists {
		return value
	}
	return defaultValue
}
