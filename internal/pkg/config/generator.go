package config

// GenerateConfig generates a otto-stack configuration file using schema-driven approach
func GenerateConfig(projectName string, services []string, validation, advanced map[string]bool) string {
	config, err := GenerateConfigFromSchema(projectName, services, nil)
	if err != nil {
		// This should not happen with a valid schema, but provide a fallback
		return "# Error generating config from schema\n"
	}
	return config
}
