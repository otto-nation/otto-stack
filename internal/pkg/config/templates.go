package config

import (
	"fmt"
	"strings"

	"github.com/otto-nation/otto-stack/internal/pkg/constants"
)

// DefaultConfigTemplate generates a basic configuration template
func DefaultConfigTemplate(projectName string, services []string) string {
	servicesYAML := ""
	if len(services) > 0 {
		servicesYAML = fmt.Sprintf("  - %s", strings.Join(services, "\n  - "))
	}

	return fmt.Sprintf(`# Otto Stack Configuration
# Generated for project: %s

project:
  name: %s
  environment: %s

stack:
  enabled:
%s

# Service-specific configuration (optional)
# service_configuration:
#   postgres:
#     POSTGRES_DB: myapp
#   redis:
#     REDIS_MAXMEMORY: 256mb
`, projectName, projectName, constants.DefaultEnvironment, servicesYAML)
}

// LocalConfigTemplate generates a local override template
func LocalConfigTemplate() string {
	return fmt.Sprintf(`# Otto Stack Local Configuration
# This file overrides settings in %s
# Add it to .gitignore to keep local settings private

project:
  environment: local

# Override enabled services for local development
# stack:
#   enabled:
#     - postgres
#     - redis

# Local service configuration
# service_configuration:
#   postgres:
#     POSTGRES_PASSWORD: localdev
`, constants.ConfigFileName)
}
