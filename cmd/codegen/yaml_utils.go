package codegen

import (
	"os"
	"strings"

	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"gopkg.in/yaml.v3"
)

// LoadYAMLConfig loads and parses a YAML file into a map
func LoadYAMLConfig(path string) (map[string]any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config map[string]any
	err = yaml.Unmarshal(data, &config)
	return config, err
}

// ToPascalCase converts a string to PascalCase
func ToPascalCase(s string) string {
	if s == "" {
		return ""
	}

	parts := strings.FieldsFunc(s, func(c rune) bool {
		return c == '-' || c == '_' || c == ' '
	})

	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
		}
	}

	return strings.Join(parts, "")
}

// YAMLLoader handles YAML file loading operations
type YAMLLoader struct{}

// NewYAMLLoader creates a new YAML loader
func NewYAMLLoader() *YAMLLoader {
	return &YAMLLoader{}
}

// LoadYAMLFile loads and parses a YAML file into the provided structure
func (yl *YAMLLoader) LoadYAMLFile(filePath string, target any) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return pkgerrors.NewConfigError("file", "failed to read YAML file", err)
	}

	if err := yaml.Unmarshal(data, target); err != nil {
		return pkgerrors.NewConfigError("yaml", "failed to parse YAML content", err)
	}

	return nil
}

// StringUtils provides common string manipulation utilities for code generation
type StringUtils struct{}

// NewStringUtils creates a new string utilities instance
func NewStringUtils() *StringUtils {
	return &StringUtils{}
}

// ToPascalCase converts a string to PascalCase
func (su *StringUtils) ToPascalCase(s string) string {
	return ToPascalCase(s)
}

// ToSnakeCase converts a string to snake_case
func (su *StringUtils) ToSnakeCase(s string) string {
	if s == "" {
		return ""
	}

	var result strings.Builder
	for i, r := range s {
		if i > 0 && (r >= 'A' && r <= 'Z') {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}

	return strings.ToLower(result.String())
}

// ToConstantCase converts a string to CONSTANT_CASE
func (su *StringUtils) ToConstantCase(s string) string {
	return strings.ToUpper(su.ToSnakeCase(s))
}
