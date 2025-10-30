package version

import (
	"os"
	"path/filepath"

	"github.com/otto-nation/otto-stack/internal/pkg/constants"
	"gopkg.in/yaml.v3"
)

// DetectProjectVersion detects the required version for a project
func DetectProjectVersion(projectPath string) (*VersionConstraint, error) {
	// Try to detect from main configuration file
	constraint, err := detectFromConfigFiles(projectPath)
	if err != nil {
		return nil, err
	}

	if constraint != nil {
		return constraint, nil
	}

	// Default to any version if no specific requirement found
	return &VersionConstraint{
		Operator: "*",
		Version:  Version{Major: 0, Minor: 0, Patch: 0},
		Original: "*",
	}, nil
}

// detectFromConfigFiles tries to detect version requirements from other config files
func detectFromConfigFiles(projectPath string) (*VersionConstraint, error) {
	// Check for otto-stack configuration file
	configPath := filepath.Join(projectPath, constants.ConfigFileName)
	if _, err := os.Stat(configPath); err == nil {
		constraint, err := parseVersionFromConfig(configPath)
		if err == nil && constraint != nil {
			return constraint, nil
		}
	}

	return nil, nil
}

// parseVersionFromConfig extracts version requirements from a config file
func parseVersionFromConfig(configPath string) (*VersionConstraint, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config map[string]any
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	versionStr := getNestedString(config, "version_config", "required_version")
	if versionStr == "" {
		return nil, nil
	}

	return ParseVersionConstraint(versionStr)
}

// getNestedString safely extracts a nested string value from a map
func getNestedString(data map[string]any, keys ...string) string {
	current := data
	for _, key := range keys[:len(keys)-1] {
		next, exists := current[key]
		if !exists {
			return ""
		}
		if nextMap, ok := next.(map[string]any); ok {
			current = nextMap
		} else {
			return ""
		}
	}

	finalKey := keys[len(keys)-1]
	if value, exists := current[finalKey]; exists {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return ""
}

// ValidateProjectVersion validates the version constraint in a project
func ValidateProjectVersion(projectPath string) error {
	constraint, err := DetectProjectVersion(projectPath)
	if err != nil {
		return err
	}

	if constraint == nil {
		return NewVersionError(ErrProjectConfig,
			"no version constraint found", nil)
	}

	// Validate the constraint format
	return ValidateConstraint(constraint.Original)
}
