package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	commandsYAMLPath = "internal/config/commands.yaml"
	schemaYAMLPath   = "internal/config/schema.yaml"
	servicesDirPath  = "internal/config/services"
	contributingPath = "CONTRIBUTING.md"
	readmePath       = "README.md"
	outputDirPath    = "docs-site/content"
	staticDate       = "2025-10-01"
)

var (
	exampleServices         = []string{"postgres", "redis"}
	completeExampleServices = []string{"postgres", "redis", "kafka"}
)

// ---- Frontmatter ----

type frontmatter struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Lead        string `yaml:"lead"`
	Date        string `yaml:"date"`
	Lastmod     string `yaml:"lastmod"`
	Draft       bool   `yaml:"draft"`
	Weight      int    `yaml:"weight"`
	Toc         bool   `yaml:"toc"`
}

func today() string {
	return time.Now().Format("2006-01-02")
}

func newFrontmatter(title, description, lead string, weight int) frontmatter {
	return frontmatter{
		Title:       title,
		Description: description,
		Lead:        lead,
		Date:        staticDate,
		Lastmod:     today(),
		Draft:       false,
		Weight:      weight,
		Toc:         true,
	}
}

func formatDocument(fm frontmatter, content string) (string, error) {
	fmBytes, err := yaml.Marshal(fm)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("---\n%s---\n\n%s", fmBytes, content), nil
}

// ---- yaml.Node helpers ----

// nodeDoc unwraps a DocumentNode to its first child.
func nodeDoc(n *yaml.Node) *yaml.Node {
	if n != nil && n.Kind == yaml.DocumentNode && len(n.Content) > 0 {
		return n.Content[0]
	}
	return n
}

// nodeGet returns the value node for the given key in a mapping node.
func nodeGet(n *yaml.Node, key string) *yaml.Node {
	n = nodeDoc(n)
	if n == nil || n.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i+1 < len(n.Content); i += 2 {
		if n.Content[i].Value == key {
			return n.Content[i+1]
		}
	}
	return nil
}

// nodeKeys returns the keys of a mapping node in document order.
func nodeKeys(n *yaml.Node) []string {
	n = nodeDoc(n)
	if n == nil || n.Kind != yaml.MappingNode {
		return nil
	}
	keys := make([]string, 0, len(n.Content)/2)
	for i := 0; i < len(n.Content); i += 2 {
		keys = append(keys, n.Content[i].Value)
	}
	return keys
}

// nodeStr returns the string value of a scalar node.
func nodeStr(n *yaml.Node) string {
	if n == nil {
		return ""
	}
	return n.Value
}

// nodeBool returns the boolean value of a scalar node.
func nodeBool(n *yaml.Node) bool {
	if n == nil {
		return false
	}
	return n.Value == "true"
}

// nodeStringSlice returns the string values of a sequence node.
func nodeStringSlice(n *yaml.Node) []string {
	if n == nil || n.Kind != yaml.SequenceNode {
		return nil
	}
	result := make([]string, 0, len(n.Content))
	for _, item := range n.Content {
		result = append(result, item.Value)
	}
	return result
}

// marshalYAML marshals a yaml.Node to string with 2-space indent.
func marshalYAML(n *yaml.Node) (string, error) {
	if n == nil {
		return "", nil
	}
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(n); err != nil {
		return "", err
	}
	_ = enc.Close()
	return buf.String(), nil
}

// ---- Ordered schema property types ----

type schemaField struct {
	Name        string
	Type        string
	Description string
	Default     *yaml.Node
	Required    bool
	Items       *schemaItems
	Properties  []*schemaField
}

type schemaItems struct {
	Type       string
	Properties []*schemaField
}

// extractSchemaFields extracts ordered fields from a configuration_schema node.
// The node should be the value of the "configuration_schema" key.
func extractSchemaFields(schemaNode *yaml.Node) []*schemaField {
	if schemaNode == nil {
		return nil
	}
	propsNode := nodeGet(schemaNode, "properties")
	if propsNode == nil || propsNode.Kind != yaml.MappingNode {
		return nil
	}
	return extractPropertiesNode(propsNode)
}

func extractPropertiesNode(propsNode *yaml.Node) []*schemaField {
	if propsNode == nil || propsNode.Kind != yaml.MappingNode {
		return nil
	}
	var fields []*schemaField
	for i := 0; i+1 < len(propsNode.Content); i += 2 {
		keyNode := propsNode.Content[i]
		valNode := propsNode.Content[i+1]
		field := &schemaField{Name: keyNode.Value}
		field.Type = nodeStr(nodeGet(valNode, "type"))
		field.Description = nodeStr(nodeGet(valNode, "description"))
		field.Default = nodeGet(valNode, "default")
		if reqNode := nodeGet(valNode, "required"); reqNode != nil {
			field.Required = nodeBool(reqNode)
		}
		if itemsNode := nodeGet(valNode, "items"); itemsNode != nil {
			field.Items = &schemaItems{
				Type:       nodeStr(nodeGet(itemsNode, "type")),
				Properties: extractPropertiesNode(nodeGet(itemsNode, "properties")),
			}
		}
		if subPropsNode := nodeGet(valNode, "properties"); subPropsNode != nil {
			field.Properties = extractPropertiesNode(subPropsNode)
		}
		fields = append(fields, field)
	}
	return fields
}

// buildExamplesNode builds an ordered yaml.Node representing the example configuration.
// This mirrors the JS SchemaParser.generateSchemaExamples logic.
func buildExamplesNode(fields []*schemaField) *yaml.Node {
	mapping := &yaml.Node{Kind: yaml.MappingNode}
	for _, f := range fields {
		var valNode *yaml.Node
		switch f.Type {
		case "string":
			if f.Default != nil && f.Default.Value != "" {
				valNode = &yaml.Node{Kind: yaml.ScalarNode, Value: f.Default.Value}
			}
		case "integer":
			if f.Default != nil {
				valNode = &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!int", Value: f.Default.Value}
			}
		case "boolean":
			if f.Default != nil {
				valNode = &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!bool", Value: f.Default.Value}
			}
		case "array":
			if f.Items != nil {
				itemNode := buildItemExampleNode(f.Items)
				if itemNode != nil {
					valNode = &yaml.Node{Kind: yaml.SequenceNode, Content: []*yaml.Node{itemNode}}
				}
			}
		case "object":
			if len(f.Properties) > 0 {
				valNode = buildObjectExampleNode(f.Properties)
			}
		}
		if valNode == nil {
			continue
		}
		mapping.Content = append(mapping.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Value: f.Name},
			valNode,
		)
	}
	if len(mapping.Content) == 0 {
		return nil
	}
	return mapping
}

func buildItemExampleNode(items *schemaItems) *yaml.Node {
	if len(items.Properties) == 0 {
		return &yaml.Node{Kind: yaml.MappingNode}
	}
	return buildObjectExampleNode(items.Properties)
}

func buildObjectExampleNode(props []*schemaField) *yaml.Node {
	mapping := &yaml.Node{Kind: yaml.MappingNode}
	for _, p := range props {
		var valNode *yaml.Node
		if p.Default != nil {
			valNode = &yaml.Node{Kind: yaml.ScalarNode, Tag: p.Default.Tag, Value: p.Default.Value}
		} else {
			switch p.Type {
			case "string":
				valNode = &yaml.Node{Kind: yaml.ScalarNode, Value: "example-" + p.Name}
			case "integer":
				valNode = &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!int", Value: "1"}
			case "boolean":
				valNode = &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!bool", Value: "true"}
			}
		}
		if valNode == nil {
			continue
		}
		mapping.Content = append(mapping.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Value: p.Name},
			valNode,
		)
	}
	return mapping
}

// ---- Service types ----

type serviceConfig struct {
	Name          string                `yaml:"name"`
	Description   string                `yaml:"description"`
	Hidden        bool                  `yaml:"hidden"`
	Environment   map[string]string     `yaml:"environment"`
	Documentation *serviceDocumentation `yaml:"documentation"`
	// configSchema is parsed separately via yaml.Node
	configSchemaFields []*schemaField
}

type serviceDocumentation struct {
	UseCases []string `yaml:"use_cases"`
	Examples []string `yaml:"examples"`
}

type loadedService struct {
	name     string
	config   serviceConfig
	category string
}

func loadAllServices() ([]loadedService, error) {
	var services []loadedService
	err := filepath.Walk(servicesDirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		ext := filepath.Ext(path)
		if ext != ".yaml" && ext != ".yml" {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}

		var svc serviceConfig
		if err := yaml.Unmarshal(data, &svc); err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}
		if svc.Hidden {
			return nil
		}

		// Extract schema fields from yaml.Node for order preservation
		var rootNode yaml.Node
		if err := yaml.Unmarshal(data, &rootNode); err != nil {
			return fmt.Errorf("parse node %s: %w", path, err)
		}
		schemaNode := nodeGet(&rootNode, "configuration_schema")
		svc.configSchemaFields = extractSchemaFields(schemaNode)

		// Determine category from directory structure
		relPath, _ := filepath.Rel(servicesDirPath, path)
		parts := strings.Split(filepath.ToSlash(relPath), "/")
		category := "other"
		if len(parts) >= 2 {
			category = parts[0]
		}

		name := strings.TrimSuffix(filepath.Base(path), ext)
		services = append(services, loadedService{name: name, config: svc, category: category})
		return nil
	})
	return services, err
}

type categoryConfig struct {
	icon  string
	order int
}

var categoryConfigs = map[string]categoryConfig{
	"database":      {icon: "🗄️", order: 1},
	"cache":         {icon: "⚡", order: 2},
	"messaging":     {icon: "📨", order: 3},
	"cloud":         {icon: "☁️", order: 4},
	"observability": {icon: "🔍", order: 5},
	"other":         {icon: "🔧", order: 99},
}

func getCategoryConfig(name string) categoryConfig {
	if c, ok := categoryConfigs[name]; ok {
		return c
	}
	return categoryConfigs["other"]
}

// ---- Generator: services-guide ----

func generateServicesGuide() error {
	services, err := loadAllServices()
	if err != nil {
		return fmt.Errorf("load services: %w", err)
	}

	// Group by category
	byCategory := make(map[string][]loadedService)
	for _, svc := range services {
		byCategory[svc.category] = append(byCategory[svc.category], svc)
	}

	// Sort categories by order
	var categories []string
	for cat := range byCategory {
		categories = append(categories, cat)
	}
	sort.Slice(categories, func(i, j int) bool {
		return getCategoryConfig(categories[i]).order < getCategoryConfig(categories[j]).order
	})

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# Available Services\n\n%d services available for your development stack.\n\n", len(services)))
	sb.WriteString("Each service can be configured through the `service_configuration` section in your `otto-stack-config.yaml` file. For detailed configuration instructions, see the [Configuration Guide](/otto-stack/configuration/).\n\n")

	for _, cat := range categories {
		catCfg := getCategoryConfig(cat)
		catTitle := strings.ToUpper(cat[:1]) + cat[1:]
		sb.WriteString(fmt.Sprintf("## %s %s\n\n", catCfg.icon, catTitle))

		// Sort services alphabetically within category
		svcs := byCategory[cat]
		sort.Slice(svcs, func(i, j int) bool {
			return svcs[i].name < svcs[j].name
		})

		for _, svc := range svcs {
			sb.WriteString(renderServiceSection(svc))
		}
	}

	fm := newFrontmatter(
		"Services",
		"Available services and configuration options",
		"Explore all the services you can use with otto-stack",
		30,
	)
	out, err := formatDocument(fm, sb.String())
	if err != nil {
		return err
	}
	return writeOutput("services.md", out)
}

func renderServiceSection(svc loadedService) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("### %s\n\n%s\n\n", svc.name, svc.config.Description))

	if len(svc.config.configSchemaFields) > 0 {
		sb.WriteString("#### Configuration Options\n\n")
		for _, field := range svc.config.configSchemaFields {
			sb.WriteString(renderSchemaField(field, "####"))
		}

		// Build examples
		examplesNode := buildExamplesNode(svc.config.configSchemaFields)
		if examplesNode != nil {
			exYAML, err := marshalYAML(examplesNode)
			if err == nil && strings.TrimSpace(exYAML) != "" {
				sb.WriteString("\n##### Example Configuration\n\n")
				sb.WriteString("```yaml\n")
				sb.WriteString(strings.TrimRight(exYAML, "\n"))
				sb.WriteString("\n```\n\n")
			}
		}
	}

	useCases := svc.config.Documentation
	if useCases != nil && len(useCases.UseCases) > 0 {
		sb.WriteString("#### Use Cases\n\n")
		for _, uc := range useCases.UseCases {
			sb.WriteString(fmt.Sprintf("- %s\n\n", uc))
		}
	}

	if useCases != nil && len(useCases.Examples) > 0 {
		sb.WriteString("#### Examples\n\n")
		for _, ex := range useCases.Examples {
			sb.WriteString("```bash\n")
			sb.WriteString(ex)
			sb.WriteString("\n```\n\n")
		}
	}

	sb.WriteString("---\n\n")
	return sb.String()
}

func renderSchemaField(field *schemaField, headingLevel string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s %s\n\n", headingLevel, field.Name))
	if field.Description != "" {
		sb.WriteString(field.Description + "\n\n")
	}
	sb.WriteString(fmt.Sprintf("- Type: `%s`\n", field.Type))
	if field.Default != nil && field.Default.Value != "" {
		sb.WriteString(fmt.Sprintf("- Default: `%s`\n", field.Default.Value))
	}
	if field.Required {
		sb.WriteString("- Required: Yes\n")
	}
	sb.WriteString("\n")

	if field.Items != nil && len(field.Items.Properties) > 0 {
		sb.WriteString("**Items:**\n\n")
		for _, itemProp := range field.Items.Properties {
			sb.WriteString(renderItemProperty(itemProp))
		}
		sb.WriteString("\n")
	}

	if len(field.Properties) > 0 {
		sb.WriteString("**Properties:**\n\n")
		for _, subProp := range field.Properties {
			sb.WriteString(renderItemProperty(subProp))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func renderItemProperty(p *schemaField) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("- **%s** (`%s`)", p.Name, p.Type))
	if p.Required {
		sb.WriteString(" _required_")
	}
	if p.Default != nil && p.Default.Value != "" {
		sb.WriteString(fmt.Sprintf(" = `%s`", p.Default.Value))
	}
	if p.Description != "" {
		sb.WriteString(": " + p.Description)
	}
	sb.WriteString("\n\n")
	return sb.String()
}

// ---- Generator: cli-reference ----

func generateCLIReference() error {
	data, err := os.ReadFile(commandsYAMLPath)
	if err != nil {
		return fmt.Errorf("read commands.yaml: %w", err)
	}

	var rootNode yaml.Node
	if err := yaml.Unmarshal(data, &rootNode); err != nil {
		return fmt.Errorf("parse commands.yaml node: %w", err)
	}

	// Get description from metadata
	metadataNode := nodeGet(&rootNode, "metadata")
	description := nodeStr(nodeGet(metadataNode, "description"))
	if description == "" {
		description = "A powerful development stack management tool for streamlined local development automation"
	}

	var sb strings.Builder
	sb.WriteString(`<!--
  ⚠️  AUTO-GENERATED FILE - DO NOT EDIT DIRECTLY
  This file is generated from internal/config/commands.yaml
  To make changes, edit the source file and run: task generate:docs
-->

# otto-stack CLI Reference

`)
	sb.WriteString(description + "\n\n")
	sb.WriteString("## Command Categories\n\n")

	// Categories in YAML file order
	categoriesNode := nodeGet(&rootNode, "categories")
	for _, catKey := range nodeKeys(categoriesNode) {
		catNode := nodeGet(categoriesNode, catKey)
		icon := nodeStr(nodeGet(catNode, "icon"))
		name := nodeStr(nodeGet(catNode, "name"))
		desc := nodeStr(nodeGet(catNode, "description"))
		cmds := nodeStringSlice(nodeGet(catNode, "commands"))

		quotedCmds := make([]string, len(cmds))
		for i, c := range cmds {
			quotedCmds[i] = "`" + c + "`"
		}

		sb.WriteString(fmt.Sprintf("### %s %s\n\n%s\n\n**Commands:** %s\n\n", icon, name, desc, strings.Join(quotedCmds, ", ")))
	}

	sb.WriteString("## Commands\n\n")

	// Commands in YAML file order
	commandsNode := nodeGet(&rootNode, "commands")
	for _, cmdKey := range nodeKeys(commandsNode) {
		cmdNode := nodeGet(commandsNode, cmdKey)
		sb.WriteString(renderCommandSection(cmdKey, cmdNode))
	}

	// Global flags (if present at top level as global_flags)
	globalFlagsNode := nodeGet(&rootNode, "global_flags")
	if globalFlagsNode != nil {
		sb.WriteString("## Global Flags\n\nThese flags are available for all commands:\n\n")
		for _, flagKey := range nodeKeys(globalFlagsNode) {
			flagNode := nodeGet(globalFlagsNode, flagKey)
			short := nodeStr(nodeGet(flagNode, "short"))
			flagDesc := nodeStr(nodeGet(flagNode, "description"))
			defaultNode := nodeGet(flagNode, "default")

			line := fmt.Sprintf("- `--%s`", flagKey)
			if short != "" {
				line += fmt.Sprintf(", `-%s`", short)
			}
			line += ": " + flagDesc
			if defaultNode != nil {
				line += fmt.Sprintf(" (default: `%s`)", defaultNode.Value)
			}
			sb.WriteString(line + "\n")
		}
		sb.WriteString("\n")
	}

	fm := newFrontmatter(
		"CLI Reference",
		"Complete command reference for otto-stack CLI",
		"Comprehensive reference for all otto-stack CLI commands and their usage",
		50,
	)
	out, err := formatDocument(fm, sb.String())
	if err != nil {
		return err
	}
	return writeOutput("cli-reference.md", out)
}

func renderCommandSection(name string, cmdNode *yaml.Node) string {
	var sb strings.Builder

	desc := nodeStr(nodeGet(cmdNode, "description"))
	longDesc := nodeStr(nodeGet(cmdNode, "long_description"))
	usage := nodeStr(nodeGet(cmdNode, "usage"))
	aliases := nodeStringSlice(nodeGet(cmdNode, "aliases"))

	sb.WriteString(fmt.Sprintf("### `%s`\n\n%s\n\n", name, desc))

	if longDesc != "" {
		sb.WriteString(strings.TrimSpace(longDesc) + "\n\n")
	}

	if usage != "" {
		sb.WriteString(fmt.Sprintf("**Usage:** `otto-stack %s`\n\n", usage))
	}

	if len(aliases) > 0 {
		quoted := make([]string, len(aliases))
		for i, a := range aliases {
			quoted[i] = "`" + a + "`"
		}
		sb.WriteString("**Aliases:** " + strings.Join(quoted, ", ") + "\n\n")
	}

	// Examples
	examplesNode := nodeGet(cmdNode, "examples")
	if examplesNode != nil && examplesNode.Kind == yaml.SequenceNode {
		sb.WriteString("**Examples:**\n\n")
		for _, exNode := range examplesNode.Content {
			cmd := nodeStr(nodeGet(exNode, "command"))
			exDesc := nodeStr(nodeGet(exNode, "description"))
			sb.WriteString(fmt.Sprintf("```bash\n%s\n```\n\n", cmd))
			if exDesc != "" {
				sb.WriteString(exDesc + "\n\n")
			}
		}
	}

	// Flags
	flagsNode := nodeGet(cmdNode, "flags")
	if flagsNode != nil && len(nodeKeys(flagsNode)) > 0 {
		sb.WriteString("**Flags:**\n\n")
		for _, flagKey := range nodeKeys(flagsNode) {
			flagNode := nodeGet(flagsNode, flagKey)
			short := nodeStr(nodeGet(flagNode, "short"))
			flagType := nodeStr(nodeGet(flagNode, "type"))
			flagDesc := nodeStr(nodeGet(flagNode, "description"))
			defaultNode := nodeGet(flagNode, "default")
			optionsNode := nodeGet(flagNode, "options")

			line := fmt.Sprintf("- `--%s`", flagKey)
			if short != "" {
				line += fmt.Sprintf(", `-%s`", short)
			}
			if flagType != "" {
				line += fmt.Sprintf(" (`%s`)", flagType)
			}
			line += ": " + flagDesc
			if defaultNode != nil {
				line += fmt.Sprintf(" (default: `%s`)", defaultNode.Value)
			}
			if optionsNode != nil {
				opts := nodeStringSlice(optionsNode)
				quoted := make([]string, len(opts))
				for i, o := range opts {
					quoted[i] = "`" + o + "`"
				}
				line += " (options: " + strings.Join(quoted, ", ") + ")"
			}
			sb.WriteString(line + "\n")
		}
		sb.WriteString("\n")
	}

	// Related commands
	relatedNode := nodeGet(cmdNode, "related_commands")
	if relatedNode != nil {
		related := nodeStringSlice(relatedNode)
		if len(related) > 0 {
			links := make([]string, len(related))
			for i, r := range related {
				links[i] = fmt.Sprintf("[`%s`](#%s)", r, r)
			}
			sb.WriteString("**Related Commands:** " + strings.Join(links, ", ") + "\n\n")
		}
	}

	// Tips
	tipsNode := nodeGet(cmdNode, "tips")
	if tipsNode != nil {
		tips := nodeStringSlice(tipsNode)
		if len(tips) > 0 {
			sb.WriteString("**Tips:**\n\n")
			for _, tip := range tips {
				sb.WriteString("- " + tip + "\n")
			}
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

// ---- Generator: configuration-guide ----

// schemaSection represents a parsed section from schema.yaml with ordered properties.
type schemaSection struct {
	name        string
	description string
	properties  []*schemaSectionProp
}

type schemaSectionProp struct {
	key         string
	propType    string
	description string
	defaultVal  string
	isTemplate  bool // default starts with {{
}

func generateConfigurationGuide() error {
	// Load schema
	schemaData, err := os.ReadFile(schemaYAMLPath)
	if err != nil {
		return fmt.Errorf("read schema.yaml: %w", err)
	}
	var schemaRoot yaml.Node
	if err := yaml.Unmarshal(schemaData, &schemaRoot); err != nil {
		return fmt.Errorf("parse schema.yaml: %w", err)
	}

	// Load services for env var examples
	services, err := loadAllServices()
	if err != nil {
		return fmt.Errorf("load services: %w", err)
	}
	svcMap := make(map[string]loadedService)
	for _, svc := range services {
		svcMap[svc.name] = svc
	}

	schemaNode := nodeGet(&schemaRoot, "schema")
	sections := extractSchemaSections(schemaNode)

	templateData := struct {
		FileStructure        string
		ConfigStructure      string
		ConfigSections       string
		ServiceConfigExample string
		CustomEnvExample     string
		CompleteExample      string
		CompleteEnvExample   string
	}{
		FileStructure:        generateFileStructure(),
		ConfigStructure:      generateConfigStructure(sections),
		ConfigSections:       generateConfigSections(sections),
		ServiceConfigExample: generateServiceConfigExample(svcMap),
		CustomEnvExample:     generateCustomEnvExample(svcMap),
		CompleteExample:      generateCompleteExample(schemaNode),
		CompleteEnvExample:   generateCompleteEnvExample(svcMap),
	}

	fm := newFrontmatter(
		"Configuration Guide",
		"Configure your otto-stack development environment",
		"Learn how to configure your development stack",
		25,
	)
	fmBytes, err := yaml.Marshal(fm)
	if err != nil {
		return err
	}

	// Use strings.Builder to avoid backtick-in-backtick issues
	const fence = "```"
	var sb strings.Builder

	sb.WriteString("---\n")
	sb.WriteString(string(fmBytes))
	sb.WriteString("---\n\n")

	sb.WriteString("<!--\n")
	sb.WriteString("  \u26a0\ufe0f  PARTIALLY GENERATED FILE\n")
	sb.WriteString("  - Sections marked with triple braces are auto-generated from internal/config/schema.yaml\n")
	sb.WriteString("  - Custom content (like \"Sharing Configuration Details\") is maintained in docs-site/templates/configuration-guide.md\n")
	sb.WriteString("  - To regenerate, run: task generate:docs\n")
	sb.WriteString("-->\n\n")

	sb.WriteString("# Configuration Guide\n\n")
	sb.WriteString("Otto-stack uses `.otto-stack/config.yaml` to define your development stack.\n\n")

	sb.WriteString("## File Structure\n\n")
	sb.WriteString("After running `otto-stack init`, you'll have:\n\n")
	sb.WriteString(fence + "\n")
	sb.WriteString(templateData.FileStructure + "\n")
	sb.WriteString(fence + "\n\n")

	sb.WriteString("## Main Configuration\n\n")
	sb.WriteString("**`.otto-stack/config.yaml`:**\n\n")
	sb.WriteString(fence + "yaml\n")
	sb.WriteString(strings.TrimRight(templateData.ConfigStructure, "\n") + "\n")
	sb.WriteString(fence + "\n\n")

	sb.WriteString(templateData.ConfigSections + "\n\n")

	sb.WriteString("### Sharing Configuration Details\n\n")
	sb.WriteString("When sharing is enabled:\n")
	sb.WriteString("1. Containers are prefixed with `otto-stack-` (e.g., `otto-stack-redis`)\n")
	sb.WriteString("2. A registry at `~/.otto-stack/shared/containers.yaml` tracks which projects use each shared container\n")
	sb.WriteString("3. The `down` command prompts before stopping shared containers used by other projects\n")
	sb.WriteString("4. Shared containers persist across project switches\n\n")

	sb.WriteString("**Example configurations:**\n\n")
	sb.WriteString(fence + "yaml\n")
	sb.WriteString("# Share all services (default)\n")
	sb.WriteString("sharing:\n  enabled: true\n\n")
	sb.WriteString("# Share specific services only\n")
	sb.WriteString("sharing:\n  enabled: true\n  services:\n    postgres: true\n    redis: true\n    kafka: false  # Not shared\n\n")
	sb.WriteString("# Disable sharing completely\n")
	sb.WriteString("sharing:\n  enabled: false\n")
	sb.WriteString(fence + "\n\n")

	sb.WriteString("**Registry location:** `~/.otto-stack/shared/containers.yaml`\n\n")

	sb.WriteString("## Service Configuration\n\n")
	sb.WriteString("Services are configured through environment variables. Otto-stack generates `.otto-stack/generated/.env.generated` showing all available variables with defaults:\n\n")
	sb.WriteString("**Example `.env.generated`:**\n\n")
	sb.WriteString(fence + "bash\n")
	sb.WriteString("# " + templateData.ServiceConfigExample + "\n")
	sb.WriteString(fence + "\n\n")

	sb.WriteString("### Customizing Services\n\n")
	sb.WriteString("Create a `.env` file in your project root to override defaults:\n\n")
	sb.WriteString(fence + "bash\n")
	sb.WriteString(templateData.CustomEnvExample + "\n")
	sb.WriteString(fence + "\n\n")

	sb.WriteString("These values will be used by Docker Compose when starting services.\n\n")

	sb.WriteString("## Service Metadata Files\n\n")
	sb.WriteString("Files in `.otto-stack/services/` contain service metadata:\n\n")
	sb.WriteString("**`.otto-stack/services/postgres.yml`:**\n\n")
	sb.WriteString(fence + "yaml\n")
	sb.WriteString("name: postgres\ndescription: Configuration for postgres service\n")
	sb.WriteString(fence + "\n\n")
	sb.WriteString("These are informational and don't affect service behavior. Configuration happens via environment variables.\n\n")

	sb.WriteString("## Complete Example\n\n")
	sb.WriteString("**`.otto-stack/config.yaml`:**\n\n")
	sb.WriteString(fence + "yaml\n")
	sb.WriteString(strings.TrimRight(templateData.CompleteExample, "\n") + "\n")
	sb.WriteString(fence + "\n\n")

	sb.WriteString("**`.env` (your customizations):**\n\n")
	sb.WriteString(fence + "bash\n")
	sb.WriteString(templateData.CompleteEnvExample + "\n")
	sb.WriteString(fence + "\n\n")

	sb.WriteString("## Next Steps\n\n")
	sb.WriteString("- **[Services Guide](/otto-stack/services/)** - Available services and environment variables\n")
	sb.WriteString("- **[CLI Reference](/otto-stack/cli-reference/)** - Command usage\n")
	sb.WriteString("- **[Troubleshooting](/otto-stack/troubleshooting/)** - Common issues\n")

	return writeOutput("configuration.md", sb.String())
}

func extractSchemaSections(schemaNode *yaml.Node) []schemaSection {
	if schemaNode == nil {
		return nil
	}
	var sections []schemaSection
	for i := 0; i+1 < len(schemaNode.Content); i += 2 {
		sectionName := schemaNode.Content[i].Value
		sectionNode := schemaNode.Content[i+1]
		propsNode := nodeGet(sectionNode, "properties")
		if propsNode == nil || propsNode.Kind != yaml.MappingNode {
			continue
		}
		section := schemaSection{
			name:        sectionName,
			description: nodeStr(nodeGet(sectionNode, "description")),
		}
		for j := 0; j+1 < len(propsNode.Content); j += 2 {
			propKey := propsNode.Content[j].Value
			propNode := propsNode.Content[j+1]
			defaultNode := nodeGet(propNode, "default")
			defaultVal := ""
			isTemplate := false
			if defaultNode != nil {
				defaultVal = defaultNode.Value
				isTemplate = strings.HasPrefix(defaultVal, "{{")
			}
			section.properties = append(section.properties, &schemaSectionProp{
				key:         propKey,
				propType:    nodeStr(nodeGet(propNode, "type")),
				description: nodeStr(nodeGet(propNode, "description")),
				defaultVal:  defaultVal,
				isTemplate:  isTemplate,
			})
		}
		sections = append(sections, section)
	}
	return sections
}

func generateFileStructure() string {
	return `.otto-stack/
├── config.yaml              # Main configuration
├── generated/
│   ├── .env.generated       # Available environment variables
│   └── docker-compose.yml   # Generated Docker Compose
├── services/                # Service metadata
│   ├── postgres.yml
│   └── redis.yml
├── .gitignore
└── README.md`
}

func generateConfigStructure(sections []schemaSection) string {
	// Build an ordered yaml.Node mapping
	mapping := &yaml.Node{Kind: yaml.MappingNode}

	for _, section := range sections {
		sectionMapping := &yaml.Node{Kind: yaml.MappingNode}
		for _, prop := range section.properties {
			var valNode *yaml.Node
			switch {
			case prop.defaultVal != "" && !prop.isTemplate && prop.propType == "string":
				// Use the default, with special case for project.name
				val := prop.defaultVal
				if section.name == "project" && prop.key == "name" {
					val = "my-app"
				}
				valNode = &yaml.Node{Kind: yaml.ScalarNode, Value: val}
			case prop.isTemplate && prop.propType == "string":
				// Template default: use special handling
				if section.name == "project" && prop.key == "name" {
					valNode = &yaml.Node{Kind: yaml.ScalarNode, Value: "my-app"}
				} else {
					valNode = &yaml.Node{Kind: yaml.ScalarNode, Value: ""}
				}
			case prop.propType == "array":
				seq := &yaml.Node{Kind: yaml.SequenceNode}
				if section.name == "stack" {
					for _, s := range exampleServices {
						seq.Content = append(seq.Content, &yaml.Node{Kind: yaml.ScalarNode, Value: s})
					}
				}
				valNode = seq
			case prop.propType == "boolean":
				valNode = &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!bool", Value: "false"}
			default:
				// object, other - skip
				continue
			}
			sectionMapping.Content = append(sectionMapping.Content,
				&yaml.Node{Kind: yaml.ScalarNode, Value: prop.key},
				valNode,
			)
		}
		if len(sectionMapping.Content) == 0 {
			continue
		}
		mapping.Content = append(mapping.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Value: section.name},
			sectionMapping,
		)
	}

	result, err := marshalYAML(mapping)
	if err != nil {
		return ""
	}
	return result
}

func generateConfigSections(sections []schemaSection) string {
	var parts []string
	for _, section := range sections {
		title := titleCase(section.name)
		var propLines []string
		for _, prop := range section.properties {
			propLines = append(propLines, fmt.Sprintf("- **%s**: %s", prop.key, prop.description))
		}
		parts = append(parts, fmt.Sprintf("### %s\n\n%s\n\n%s", title, section.description, strings.Join(propLines, "\n")))
	}
	return strings.Join(parts, "\n\n")
}

func titleCase(s string) string {
	parts := strings.Split(s, "_")
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return strings.Join(parts, " ")
}

// sortedEnvKeys returns sorted env var keys for a service.
func sortedEnvKeys(env map[string]string) []string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func generateServiceConfigExample(svcMap map[string]loadedService) string {
	var lines []string
	for _, name := range completeExampleServices {
		svc, ok := svcMap[name]
		if !ok || len(svc.config.Environment) == 0 {
			continue
		}
		lines = append(lines, "# "+strings.ToUpper(name))
		keys := sortedEnvKeys(svc.config.Environment)
		limit := 4
		if len(keys) < limit {
			limit = len(keys)
		}
		for _, k := range keys[:limit] {
			lines = append(lines, k+"="+svc.config.Environment[k])
		}
	}
	return strings.Join(lines, "\n")
}

func generateCustomEnvExample(svcMap map[string]loadedService) string {
	var lines []string
	for _, name := range exampleServices {
		svc, ok := svcMap[name]
		if !ok || len(svc.config.Environment) == 0 {
			continue
		}
		capName := strings.ToUpper(name[:1]) + name[1:]
		lines = append(lines, "# "+capName)
		keys := sortedEnvKeys(svc.config.Environment)
		limit := 2
		if len(keys) < limit {
			limit = len(keys)
		}
		for _, k := range keys[:limit] {
			lines = append(lines, k+"=my_custom_value")
		}
	}
	return strings.Join(lines, "\n")
}

func generateCompleteExample(schemaNode *yaml.Node) string {
	// Build ordered YAML with project, stack, and optionally validation
	mapping := &yaml.Node{Kind: yaml.MappingNode}

	// project section
	projectMapping := &yaml.Node{Kind: yaml.MappingNode, Content: []*yaml.Node{
		{Kind: yaml.ScalarNode, Value: "name"},
		{Kind: yaml.ScalarNode, Value: "my-fullstack-app"},
		{Kind: yaml.ScalarNode, Value: "type"},
		{Kind: yaml.ScalarNode, Value: "docker"},
	}}
	mapping.Content = append(mapping.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Value: "project"},
		projectMapping,
	)

	// stack section
	stackSeq := &yaml.Node{Kind: yaml.SequenceNode}
	for _, s := range completeExampleServices {
		stackSeq.Content = append(stackSeq.Content, &yaml.Node{Kind: yaml.ScalarNode, Value: s})
	}
	stackMapping := &yaml.Node{Kind: yaml.MappingNode, Content: []*yaml.Node{
		{Kind: yaml.ScalarNode, Value: "enabled"},
		stackSeq,
	}}
	mapping.Content = append(mapping.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Value: "stack"},
		stackMapping,
	)

	// validation section (if schema has it)
	if nodeGet(schemaNode, "validation") != nil {
		optionsMapping := &yaml.Node{Kind: yaml.MappingNode, Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Value: "config-syntax"},
			{Kind: yaml.ScalarNode, Tag: "!!bool", Value: "true"},
			{Kind: yaml.ScalarNode, Value: "docker"},
			{Kind: yaml.ScalarNode, Tag: "!!bool", Value: "true"},
		}}
		validationMapping := &yaml.Node{Kind: yaml.MappingNode, Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Value: "options"},
			optionsMapping,
		}}
		mapping.Content = append(mapping.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Value: "validation"},
			validationMapping,
		)
	}

	result, err := marshalYAML(mapping)
	if err != nil {
		return ""
	}
	return result
}

func generateCompleteEnvExample(svcMap map[string]loadedService) string {
	var lines []string
	for _, name := range exampleServices {
		svc, ok := svcMap[name]
		if !ok || len(svc.config.Environment) == 0 {
			continue
		}
		capName := strings.ToUpper(name[:1]) + name[1:]
		lines = append(lines, "# "+capName)
		keys := sortedEnvKeys(svc.config.Environment)
		limit := 2
		if len(keys) < limit {
			limit = len(keys)
		}
		for _, k := range keys[:limit] {
			lines = append(lines, k+"=production_value")
		}
	}
	return strings.Join(lines, "\n")
}

// ---- Generator: homepage ----

func generateHomepage() error {
	data, err := os.ReadFile(readmePath)
	if err != nil {
		return fmt.Errorf("read README.md: %w", err)
	}

	content := string(data)

	// Find the first heading and take everything from there
	lines := strings.Split(content, "\n")
	start := 0
	for i, line := range lines {
		if strings.HasPrefix(line, "# ") {
			start = i
			break
		}
	}
	content = strings.Join(lines[start:], "\n")

	// Fix links for Hugo with baseURL subdirectory
	// Convert docs-site/content/file.md -> file/
	reDocs := regexp.MustCompile(`docs-site/content/([^)]+)\.md`)
	content = reDocs.ReplaceAllString(content, "$1/")

	// Convert remaining .md links to Hugo format
	reMd := regexp.MustCompile(`\]\(([^)]+)\.md\)`)
	content = reMd.ReplaceAllString(content, "]($1/)")

	// Fix docs-site root links
	reDocsRoot := regexp.MustCompile(`\[([^\]]+)\]\(docs-site/\)`)
	content = reDocsRoot.ReplaceAllStringFunc(content, func(match string) string {
		inner := reDocsRoot.FindStringSubmatch(match)
		if len(inner) > 1 {
			return fmt.Sprintf(`[%s]({{< ref "/" >}})`, inner[1])
		}
		return match
	})

	// Fix LICENSE link
	reLicense := regexp.MustCompile(`\[([^\]]+)\]\(LICENSE\)`)
	content = reLicense.ReplaceAllString(content, "[$1](https://github.com/otto-nation/otto-stack/blob/main/LICENSE)")

	fm := frontmatter{
		Title:       "otto-stack",
		Description: "A powerful development stack management tool built in Go for streamlined local development automation",
		Lead:        "Streamline your local development with powerful CLI tools and automated service management",
		Date:        staticDate,
		Lastmod:     today(),
		Draft:       false,
		Weight:      50,
		Toc:         true,
	}
	fmBytes, err := yaml.Marshal(fm)
	if err != nil {
		return err
	}
	out := fmt.Sprintf("---\n%s---\n\n%s", fmBytes, content)
	return writeOutput("_index.md", out)
}

// ---- Generator: contributing-guide ----

func generateContributingGuide() error {
	data, err := os.ReadFile(contributingPath)
	if err != nil {
		return fmt.Errorf("read CONTRIBUTING.md: %w", err)
	}

	fm := frontmatter{
		Title:       "Contributing",
		Description: "Guide for contributing to otto-stack development",
		Lead:        "Learn how to contribute to otto-stack development",
		Date:        staticDate,
		Lastmod:     today(),
		Draft:       false,
		Weight:      60,
		Toc:         true,
	}
	fmBytes, err := yaml.Marshal(fm)
	if err != nil {
		return err
	}
	out := fmt.Sprintf("---\n%s---\n\n%s", fmBytes, string(data))
	return writeOutput("contributing.md", out)
}

// ---- Output helpers ----

func writeOutput(filename, content string) error {
	outPath := filepath.Join(outputDirPath, filename)
	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}
	if err := os.WriteFile(outPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("write %s: %w", outPath, err)
	}
	fmt.Printf("generated %s\n", outPath)
	return nil
}

// ---- Main ----

func main() {
	generatorFlag := flag.String("generator", "", "Run a specific generator by name")
	flag.Parse()

	type generatorFn struct {
		name string
		run  func() error
	}

	allGenerators := []generatorFn{
		{"cli-reference", generateCLIReference},
		{"services-guide", generateServicesGuide},
		{"configuration-guide", generateConfigurationGuide},
		{"homepage", generateHomepage},
		{"contributing-guide", generateContributingGuide},
	}

	var toRun []generatorFn
	if *generatorFlag != "" {
		for _, g := range allGenerators {
			if g.name == *generatorFlag {
				toRun = append(toRun, g)
				break
			}
		}
		if len(toRun) == 0 {
			fmt.Fprintf(os.Stderr, "unknown generator: %s\n", *generatorFlag)
			os.Exit(1)
		}
	} else {
		toRun = allGenerators
	}

	failed := false
	for _, g := range toRun {
		if err := g.run(); err != nil {
			fmt.Fprintf(os.Stderr, "generator %s failed: %v\n", g.name, err)
			failed = true
		}
	}
	if failed {
		os.Exit(1)
	}
}
