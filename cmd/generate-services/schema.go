package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// SchemaField represents a field in the schema
type SchemaField struct {
	Name     string
	Type     string
	YamlTag  string
	Required bool
	Children map[string]*SchemaField
}

// Schema represents the extracted schema
type Schema struct {
	Fields map[string]*SchemaField
}

// extractSchema analyzes all YAML files to build schema
func extractSchema() (*Schema, error) {
	schema := &Schema{
		Fields: make(map[string]*SchemaField),
	}

	err := filepath.Walk(ServicesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || !strings.HasSuffix(path, ".yaml") {
			return err
		}
		return analyzeYAMLFile(path, schema)
	})

	return schema, err
}

// analyzeYAMLFile processes a single YAML file
func analyzeYAMLFile(path string, schema *Schema) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var yamlData map[string]any
	if err := yaml.Unmarshal(data, &yamlData); err != nil {
		return err
	}

	// Analyze top-level structure
	for key, value := range yamlData {
		field := analyzeField(key, value)
		if existing, exists := schema.Fields[key]; exists {
			mergeFields(existing, field)
		} else {
			schema.Fields[key] = field
		}
	}

	return nil
}

// analyzeField analyzes a single field
const (
	structType = "struct"
)

func analyzeField(name string, value any) *SchemaField {
	field := &SchemaField{
		Name:     name,
		YamlTag:  fmt.Sprintf(`yaml:"%s,omitempty"`, name),
		Children: make(map[string]*SchemaField),
	}

	switch v := value.(type) {
	case string:
		field.Type = "string"
	case int, int64:
		field.Type = "int"
	case bool:
		field.Type = "bool"
	case []any:
		if len(v) > 0 {
			if _, ok := v[0].(string); ok {
				field.Type = "[]string"
			} else {
				field.Type = "[]map[string]any"
			}
		} else {
			field.Type = "[]any"
		}
	case map[string]any:
		field.Type = structType
		for k, val := range v {
			field.Children[k] = analyzeField(k, val)
		}
	default:
		field.Type = "any"
	}

	return field
}

// mergeFields merges two field definitions
func mergeFields(existing, new *SchemaField) {
	// Merge children for struct types
	if existing.Type == "struct" && new.Type == "struct" {
		for k, v := range new.Children {
			if existingChild, exists := existing.Children[k]; exists {
				mergeFields(existingChild, v)
			} else {
				existing.Children[k] = v
			}
		}
	}
}
