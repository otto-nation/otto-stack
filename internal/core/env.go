package core

import (
	"os"
	"strings"
)

const expectedEnvParts = 2

// ResolveEnvVar resolves an environment variable with the format ${VAR:-default}
// Returns the environment variable value if set and non-empty,
// otherwise returns the default value if provided,
// or the original string if the syntax is malformed.
func ResolveVar(value string) string {
	// Check if value matches ${...} pattern
	if !strings.HasPrefix(value, "${") || !strings.HasSuffix(value, "}") {
		return value
	}

	// Extract inner content: ${VAR:-default} -> "VAR:-default"
	inner := value[2 : len(value)-1]

	// Split on ":-" delimiter
	parts := strings.SplitN(inner, ":-", expectedEnvParts)
	if len(parts) != expectedEnvParts {
		// No default value provided, return original
		return value
	}

	envVar := parts[0]
	defaultValue := parts[1]

	// Return env var value if set and non-empty, otherwise return default
	if envValue := os.Getenv(envVar); envValue != "" {
		return envValue
	}

	return defaultValue
}
