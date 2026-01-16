package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// JSONSchema represents a JSON Schema document
type JSONSchema struct {
	Schema      string              `json:"$schema"`
	ID          string              `json:"$id"`
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Type        string              `json:"type"`
	Required    []string            `json:"required"`
	Properties  map[string]Property `json:"properties"`
	XEnums      map[string]EnumDef  `json:"x-enums"`
}

// Property represents a JSON Schema property
type Property struct {
	Type        string              `json:"type"`
	Description string              `json:"description"`
	Enum        []string            `json:"enum"`
	Properties  map[string]Property `json:"properties"`
	Items       *Property           `json:"items"`
}

// EnumDef represents an enum definition in x-enums
type EnumDef struct {
	Description string   `json:"description"`
	Values      []string `json:"values"`
}

// ServiceSchema represents the parsed schema for code generation
type ServiceSchema struct {
	Enums          map[string]EnumDefinition
	YAMLStructure  map[string]StructureSection
	TopLevelFields []string
}

// StructureSection defines a section of the YAML structure
type StructureSection struct {
	Description string
	Fields      []string
}

// EnumDefinition defines an enum type
type EnumDefinition struct {
	Description string
	Values      []string
}

// loadServiceSchema loads the JSON schema and converts it to ServiceSchema
func loadServiceSchema(path string) (*ServiceSchema, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema: %w", err)
	}

	var jsonSchema JSONSchema
	if err := json.Unmarshal(data, &jsonSchema); err != nil {
		return nil, fmt.Errorf("failed to parse JSON schema: %w", err)
	}

	// Convert JSON Schema to ServiceSchema
	schema := &ServiceSchema{
		Enums:          make(map[string]EnumDefinition),
		YAMLStructure:  make(map[string]StructureSection),
		TopLevelFields: make([]string, 0),
	}

	// Extract enums from x-enums extension
	for enumName, enumDef := range jsonSchema.XEnums {
		schema.Enums[enumName] = EnumDefinition{
			Description: enumDef.Description,
			Values:      enumDef.Values,
		}
	}

	// Extract structure from properties
	for propName, prop := range jsonSchema.Properties {
		if prop.Type == "object" && len(prop.Properties) > 0 {
			// This is a section with nested fields
			fields := make([]string, 0, len(prop.Properties))
			for fieldName := range prop.Properties {
				fields = append(fields, fieldName)
			}
			schema.YAMLStructure[propName] = StructureSection{
				Description: prop.Description,
				Fields:      fields,
			}
		} else {
			// This is a top-level field
			schema.TopLevelFields = append(schema.TopLevelFields, propName)
		}
	}

	return schema, nil
}

// generateEnumConstants generates Go constants from enum definitions
func (s *ServiceSchema) generateEnumConstants() []EnumConstantGroup {
	var groups []EnumConstantGroup

	for enumName, enumDef := range s.Enums {
		group := EnumConstantGroup{
			Name:        enumName,
			Description: enumDef.Description,
			Constants:   make([]EnumConstant, 0, len(enumDef.Values)),
		}

		for _, value := range enumDef.Values {
			group.Constants = append(group.Constants, EnumConstant{
				Name:  toConstantName(enumName, value),
				Value: value,
			})
		}

		groups = append(groups, group)
	}

	return groups
}

// generateValidationTag generates validator tag for an enum
func (s *ServiceSchema) generateValidationTag(enumName string) string {
	if enumDef, exists := s.Enums[enumName]; exists {
		return "oneof=" + strings.Join(enumDef.Values, " ")
	}
	return ""
}

// generateYAMLKeys generates YAML structure constants
func (s *ServiceSchema) generateYAMLKeys() YAMLKeysData {
	data := YAMLKeysData{
		Sections: make([]YAMLSection, 0),
		TopLevel: make([]YAMLKey, 0),
	}

	// Track generated constants to avoid duplicates
	generated := make(map[string]bool)

	// Generate constants for each section
	for sectionName, section := range s.YAMLStructure {
		yamlSection := YAMLSection{
			Name:        sectionName,
			Description: section.Description,
			Fields:      make([]YAMLKey, 0),
		}

		// Add root key first
		rootConstName := "YAMLKey" + toPascalCaseWithAcronyms(sectionName)
		if !generated[rootConstName] {
			yamlSection.Fields = append(yamlSection.Fields, YAMLKey{
				Name:  rootConstName,
				Value: sectionName,
			})
			generated[rootConstName] = true
		}

		// Add field keys
		for _, field := range section.Fields {
			constName := "YAMLKey" + toPascalCaseWithAcronyms(field)
			// Skip if already generated
			if generated[constName] {
				continue
			}
			yamlSection.Fields = append(yamlSection.Fields, YAMLKey{
				Name:  constName,
				Value: field,
			})
			generated[constName] = true
		}

		data.Sections = append(data.Sections, yamlSection)
	}

	// Generate top-level field constants
	for _, field := range s.TopLevelFields {
		constName := "YAMLKey" + toPascalCaseWithAcronyms(field)
		// Skip if already generated
		if generated[constName] {
			continue
		}
		data.TopLevel = append(data.TopLevel, YAMLKey{
			Name:  constName,
			Value: field,
		})
		generated[constName] = true
	}

	return data
}

// YAMLKeysData represents generated YAML key constants
type YAMLKeysData struct {
	Sections []YAMLSection
	TopLevel []YAMLKey
}

// YAMLSection represents a section with fields
type YAMLSection struct {
	Name        string
	Description string
	Fields      []YAMLKey
}

// YAMLKey represents a single YAML key constant
type YAMLKey struct {
	Name  string
	Value string
}

// EnumConstantGroup represents a group of related constants
type EnumConstantGroup struct {
	Name        string
	Description string
	Constants   []EnumConstant
}

// EnumConstant represents a single constant
type EnumConstant struct {
	Name  string
	Value string
}

// toConstantName converts enum name and value to Go constant name
func toConstantName(enumName, value string) string {
	// Convert: "service_type" + "container" -> "ServiceTypeContainer"
	parts := strings.Split(enumName, "_")
	for i, part := range parts {
		parts[i] = toPascalCaseWithAcronyms(part)
	}
	prefix := strings.Join(parts, "")

	valueParts := strings.Split(value, "_")
	for i, part := range valueParts {
		valueParts[i] = toPascalCaseWithAcronyms(part)
	}
	suffix := strings.Join(valueParts, "")

	return prefix + suffix
}

// toPascalCaseWithAcronyms converts a string to PascalCase with acronym handling
func toPascalCaseWithAcronyms(s string) string {
	// Handle common acronyms
	acronyms := map[string]string{
		"sql":   "SQL",
		"aws":   "AWS",
		"http":  "HTTP",
		"https": "HTTPS",
		"api":   "API",
		"url":   "URL",
		"id":    "ID",
	}

	// Split on underscores and hyphens
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-'
	})

	result := ""
	for _, part := range parts {
		lower := strings.ToLower(part)
		if acronym, exists := acronyms[lower]; exists {
			result += acronym
		} else {
			result += strings.Title(part)
		}
	}

	return result
}
