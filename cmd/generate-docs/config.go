package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Example values used in the generated configuration guide.
const (
	exampleProjectName          = "my-app"
	exampleFullstackProjectName = "my-fullstack-app"
	exampleProjectType          = "docker"

	// envVarDisplayLimit is the number of env vars shown in the full service config example.
	envVarDisplayLimit = 4
	// customEnvDisplayLimit is the number of env vars shown in customization examples.
	customEnvDisplayLimit = 2
)

// Schema section name constants used for conditional example generation.
const (
	schemaSectionProject = "project"
	schemaSectionStack   = "stack"
	schemaPropName       = "name"
)

// Example services used in config examples (order matters for readability).
var (
	exampleServices         = []string{"postgres", "redis"}
	completeExampleServices = []string{"postgres", "redis", "kafka"}
)

// schemaSection holds an ordered set of properties parsed from schema.yaml.
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
	isTemplate  bool // true when default starts with "{{" (Go template expression)
}

func generateConfigurationGuide() error {
	var schemaRoot yaml.Node
	if err := loadYAML(schemaYAMLPath, &schemaRoot); err != nil {
		return fmt.Errorf("load schema.yaml: %w", err)
	}

	services, err := loadAllServices()
	if err != nil {
		return fmt.Errorf("load services: %w", err)
	}
	svcMap := indexServices(services)

	schemaNode := nodeGet(&schemaRoot, keySchema)
	sections := extractSchemaSections(schemaNode)

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

	writeFileStructureSection(&sb, fence)
	writeMainConfigSection(&sb, fence, generateConfigStructure(sections))
	writeConfigSections(&sb, sections)
	writeSharingSection(&sb, fence)
	writeServiceConfigSection(&sb, fence, generateServiceConfigExample(svcMap), generateCustomEnvExample(svcMap))
	writeServiceMetadataSection(&sb, fence)
	writeCompleteExampleSection(&sb, fence, generateCompleteExample(schemaNode), generateCompleteEnvExample(svcMap))
	writeNextStepsSection(&sb)

	return writeOutput("configuration.md", sb.String())
}

func writeFileStructureSection(sb *strings.Builder, fence string) {
	sb.WriteString("## File Structure\n\n")
	sb.WriteString("After running `otto-stack init`, you'll have:\n\n")
	sb.WriteString(fence + "\n")
	sb.WriteString(generateFileStructure() + "\n")
	sb.WriteString(fence + "\n\n")
}

func writeMainConfigSection(sb *strings.Builder, fence, configStructure string) {
	sb.WriteString("## Main Configuration\n\n")
	sb.WriteString("**`.otto-stack/config.yaml`:**\n\n")
	sb.WriteString(fence + "yaml\n")
	sb.WriteString(strings.TrimRight(configStructure, "\n") + "\n")
	sb.WriteString(fence + "\n\n")
}

func writeConfigSections(sb *strings.Builder, sections []schemaSection) {
	sb.WriteString(generateConfigSections(sections))
	sb.WriteString("\n\n")
}

func writeSharingSection(sb *strings.Builder, fence string) {
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
}

func writeServiceConfigSection(sb *strings.Builder, fence, serviceConfigExample, customEnvExample string) {
	sb.WriteString("## Service Configuration\n\n")
	sb.WriteString("Services are configured through environment variables. Otto-stack generates `.otto-stack/generated/.env.generated` showing all available variables with defaults:\n\n")
	sb.WriteString("**Example `.env.generated`:**\n\n")
	sb.WriteString(fence + "bash\n")
	sb.WriteString("# " + serviceConfigExample + "\n")
	sb.WriteString(fence + "\n\n")

	sb.WriteString("### Customizing Services\n\n")
	sb.WriteString("Create a `.env` file in your project root to override defaults:\n\n")
	sb.WriteString(fence + "bash\n")
	sb.WriteString(customEnvExample + "\n")
	sb.WriteString(fence + "\n\n")
	sb.WriteString("These values will be used by Docker Compose when starting services.\n\n")
}

func writeServiceMetadataSection(sb *strings.Builder, fence string) {
	sb.WriteString("## Service Metadata Files\n\n")
	sb.WriteString("Files in `.otto-stack/services/` contain service metadata:\n\n")
	sb.WriteString("**`.otto-stack/services/postgres.yml`:**\n\n")
	sb.WriteString(fence + "yaml\n")
	sb.WriteString("name: postgres\ndescription: Configuration for postgres service\n")
	sb.WriteString(fence + "\n\n")
	sb.WriteString("These are informational and don't affect service behavior. Configuration happens via environment variables.\n\n")
}

func writeCompleteExampleSection(sb *strings.Builder, fence, completeExample, completeEnvExample string) {
	sb.WriteString("## Complete Example\n\n")
	sb.WriteString("**`.otto-stack/config.yaml`:**\n\n")
	sb.WriteString(fence + "yaml\n")
	sb.WriteString(strings.TrimRight(completeExample, "\n") + "\n")
	sb.WriteString(fence + "\n\n")

	sb.WriteString("**`.env` (your customizations):**\n\n")
	sb.WriteString(fence + "bash\n")
	sb.WriteString(completeEnvExample + "\n")
	sb.WriteString(fence + "\n\n")
}

func writeNextStepsSection(sb *strings.Builder) {
	sb.WriteString("## Next Steps\n\n")
	sb.WriteString("- **[Services Guide](/otto-stack/services/)** - Available services and environment variables\n")
	sb.WriteString("- **[CLI Reference](/otto-stack/cli-reference/)** - Command usage\n")
	sb.WriteString("- **[Troubleshooting](/otto-stack/troubleshooting/)** - Common issues\n")
}

func extractSchemaSections(schemaNode *yaml.Node) []schemaSection {
	if schemaNode == nil {
		return nil
	}
	var sections []schemaSection
	for i := 0; i+1 < len(schemaNode.Content); i += 2 {
		section, ok := extractSchemaSection(schemaNode.Content[i].Value, schemaNode.Content[i+1])
		if ok {
			sections = append(sections, section)
		}
	}
	return sections
}

func extractSchemaSection(name string, sectionNode *yaml.Node) (schemaSection, bool) {
	propsNode := nodeGet(sectionNode, keyProperties)
	if propsNode == nil || propsNode.Kind != yaml.MappingNode {
		return schemaSection{}, false
	}
	section := schemaSection{
		name:        name,
		description: nodeStr(nodeGet(sectionNode, keyDescription)),
	}
	for j := 0; j+1 < len(propsNode.Content); j += 2 {
		section.properties = append(section.properties,
			extractSchemaProp(propsNode.Content[j].Value, propsNode.Content[j+1]))
	}
	return section, true
}

func extractSchemaProp(key string, propNode *yaml.Node) *schemaSectionProp {
	defaultVal := nodeStr(nodeGet(propNode, keyDefault))
	return &schemaSectionProp{
		key:         key,
		propType:    nodeStr(nodeGet(propNode, keyType)),
		description: nodeStr(nodeGet(propNode, keyDescription)),
		defaultVal:  defaultVal,
		isTemplate:  strings.HasPrefix(defaultVal, "{{"),
	}
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
	mapping := &yaml.Node{Kind: yaml.MappingNode}
	for _, section := range sections {
		sectionMapping := buildSectionMapping(section)
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

func buildSectionMapping(section schemaSection) *yaml.Node {
	mapping := &yaml.Node{Kind: yaml.MappingNode}
	for _, prop := range section.properties {
		valNode := buildSectionPropValueNode(section.name, prop)
		if valNode == nil {
			continue
		}
		mapping.Content = append(mapping.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Value: prop.key},
			valNode,
		)
	}
	return mapping
}

func buildSectionPropValueNode(sectionName string, prop *schemaSectionProp) *yaml.Node {
	switch prop.propType {
	case "string":
		return buildStringValueNode(sectionName, prop)
	case "array":
		return buildArrayValueNode(sectionName)
	case "boolean":
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!bool", Value: "false"}
	}
	return nil
}

func buildStringValueNode(sectionName string, prop *schemaSectionProp) *yaml.Node {
	if sectionName == schemaSectionProject && prop.key == schemaPropName {
		return &yaml.Node{Kind: yaml.ScalarNode, Value: exampleProjectName}
	}
	if prop.defaultVal != "" && !prop.isTemplate {
		return &yaml.Node{Kind: yaml.ScalarNode, Value: prop.defaultVal}
	}
	return &yaml.Node{Kind: yaml.ScalarNode, Value: ""}
}

func buildArrayValueNode(sectionName string) *yaml.Node {
	seq := &yaml.Node{Kind: yaml.SequenceNode}
	if sectionName == schemaSectionStack {
		for _, s := range exampleServices {
			seq.Content = append(seq.Content, &yaml.Node{Kind: yaml.ScalarNode, Value: s})
		}
	}
	return seq
}

func generateConfigSections(sections []schemaSection) string {
	var parts []string
	for _, section := range sections {
		var propLines []string
		for _, prop := range section.properties {
			propLines = append(propLines, fmt.Sprintf("- **%s**: %s", prop.key, prop.description))
		}
		parts = append(parts, fmt.Sprintf("### %s\n\n%s\n\n%s", titleCase(section.name), section.description, strings.Join(propLines, "\n")))
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

func sortedEnvKeys(env map[string]string) []string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// envExample builds a block of env var lines from the named services.
// limit caps the number of vars shown per service.
// label formats the service name for the comment header.
// value returns the value string for each key.
func envExample(
	svcMap map[string]loadedService,
	names []string,
	limit int,
	label func(string) string,
	value func(loadedService, string) string,
) string {
	var lines []string
	for _, name := range names {
		svc, ok := svcMap[name]
		if !ok || len(svc.config.Environment) == 0 {
			continue
		}
		lines = append(lines, "# "+label(name))
		keys := sortedEnvKeys(svc.config.Environment)
		n := limit
		if len(keys) < n {
			n = len(keys)
		}
		for _, k := range keys[:n] {
			lines = append(lines, k+"="+value(svc, k))
		}
	}
	return strings.Join(lines, "\n")
}

func generateServiceConfigExample(svcMap map[string]loadedService) string {
	return envExample(svcMap, completeExampleServices, envVarDisplayLimit,
		strings.ToUpper,
		func(svc loadedService, k string) string { return svc.config.Environment[k] },
	)
}

func generateCustomEnvExample(svcMap map[string]loadedService) string {
	return envExample(svcMap, exampleServices, customEnvDisplayLimit,
		func(s string) string { return strings.ToUpper(s[:1]) + s[1:] },
		func(_ loadedService, _ string) string { return "my_custom_value" },
	)
}

func generateCompleteEnvExample(svcMap map[string]loadedService) string {
	return envExample(svcMap, exampleServices, customEnvDisplayLimit,
		func(s string) string { return strings.ToUpper(s[:1]) + s[1:] },
		func(_ loadedService, _ string) string { return "production_value" },
	)
}

func generateCompleteExample(schemaNode *yaml.Node) string {
	mapping := buildCompleteExampleMapping(schemaNode)
	result, err := marshalYAML(mapping)
	if err != nil {
		return ""
	}
	return result
}

func buildCompleteExampleMapping(schemaNode *yaml.Node) *yaml.Node {
	mapping := &yaml.Node{Kind: yaml.MappingNode}
	mapping.Content = append(mapping.Content, projectExampleNodes()...)
	mapping.Content = append(mapping.Content, stackExampleNodes()...)
	if nodeGet(schemaNode, keyValidation) != nil {
		mapping.Content = append(mapping.Content, validationExampleNodes()...)
	}
	return mapping
}

func projectExampleNodes() []*yaml.Node {
	projectMapping := &yaml.Node{Kind: yaml.MappingNode, Content: []*yaml.Node{
		{Kind: yaml.ScalarNode, Value: keyName},
		{Kind: yaml.ScalarNode, Value: exampleFullstackProjectName},
		{Kind: yaml.ScalarNode, Value: keyType},
		{Kind: yaml.ScalarNode, Value: exampleProjectType},
	}}
	return []*yaml.Node{
		{Kind: yaml.ScalarNode, Value: schemaSectionProject},
		projectMapping,
	}
}

func stackExampleNodes() []*yaml.Node {
	stackSeq := &yaml.Node{Kind: yaml.SequenceNode}
	for _, s := range completeExampleServices {
		stackSeq.Content = append(stackSeq.Content, &yaml.Node{Kind: yaml.ScalarNode, Value: s})
	}
	stackMapping := &yaml.Node{Kind: yaml.MappingNode, Content: []*yaml.Node{
		{Kind: yaml.ScalarNode, Value: "enabled"},
		stackSeq,
	}}
	return []*yaml.Node{
		{Kind: yaml.ScalarNode, Value: schemaSectionStack},
		stackMapping,
	}
}

func validationExampleNodes() []*yaml.Node {
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
	return []*yaml.Node{
		{Kind: yaml.ScalarNode, Value: keyValidation},
		validationMapping,
	}
}

// loadYAML reads a YAML file and unmarshals it into the given node.
func loadYAML(path string, out *yaml.Node) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	if err := yaml.Unmarshal(data, out); err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}
	return nil
}
