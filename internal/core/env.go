package core

import (
	"os"
	"strings"
)

// ResolveVar resolves environment variable syntax ${VAR:-default}
func ResolveVar(value string) string {
	if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
		inner := value[2 : len(value)-1]
		const splitLen = 2
		if parts := strings.Split(inner, ":-"); len(parts) == splitLen {
			envVar := parts[0]
			defaultValue := parts[1]

			// Check if environment variable is set
			if envValue := os.Getenv(envVar); envValue != "" {
				return envValue
			}
			return defaultValue // Return default value if env var not set
		}
	}
	return value
}
