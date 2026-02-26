package main

import (
	"fmt"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Schema section name constants used for conditional example generation.
const (
	schemaSectionProject = "project"
	schemaSectionStack   = "stack"
	schemaPropName       = "name"
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

	const fence = "```"
	var sb strings.Builder

	sb.WriteString(htmlComment(
		"\u26a0\ufe0f  PARTIALLY GENERATED FILE",
		"- Sections marked with triple braces are auto-generated from "+schemaYAMLPath,
		`- Custom content (like "Sharing Configuration Details") is maintained in docs-site/templates/configuration-guide.md`,
		"- To regenerate, run: task generate:docs",
	))

	sb.WriteString("# " + docs.Pages["configuration"].Heading + "\n\n")
	sb.WriteString(docs.Pages["configuration"].Intro + "\n\n")

	writeFileStructureSection(&sb, fence)
	writeMainConfigSection(&sb, fence, generateConfigStructure(sections))
	writeConfigSections(&sb, sections)
	writeSharingSection(&sb, fence)
	writeServiceConfigSection(&sb, fence, generateServiceConfigExample(svcMap), generateCustomEnvExample(svcMap))
	writeServiceMetadataSection(&sb, fence)
	writeCompleteExampleSection(&sb, fence, generateCompleteExample(schemaNode), generateCompleteEnvExample(svcMap))
	writeNextStepsSection(&sb)

	out, err := formatDocument(pageFM("configuration"), sb.String())
	if err != nil {
		return err
	}
	return writeOutput(pageOutput("configuration"), out)
}

func writeFileStructureSection(sb *strings.Builder, fence string) {
	s := docs.Pages["configuration"].Sections.FileStructure
	sb.WriteString(s.Heading + "\n\n")
	sb.WriteString(s.Intro + "\n\n")
	sb.WriteString(fence + "\n")
	sb.WriteString(docs.Pages["configuration"].FileStructure + "\n")
	sb.WriteString(fence + "\n\n")
}

func writeMainConfigSection(sb *strings.Builder, fence, configStructure string) {
	s := docs.Pages["configuration"].Sections.MainConfig
	sb.WriteString(s.Heading + "\n\n")
	sb.WriteString(s.FileLabel + "\n\n")
	sb.WriteString(fence + "yaml\n")
	sb.WriteString(strings.TrimRight(configStructure, "\n") + "\n")
	sb.WriteString(fence + "\n\n")
}

func writeConfigSections(sb *strings.Builder, sections []schemaSection) {
	sb.WriteString(generateConfigSections(sections))
	sb.WriteString("\n\n")
}

func writeSharingSection(sb *strings.Builder, fence string) {
	s := docs.Pages["configuration"].Sections.Sharing
	sb.WriteString(s.Heading + "\n\n")
	sb.WriteString(s.Intro + "\n")
	for i, behavior := range s.Behaviors {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, behavior))
	}
	sb.WriteString("\n")
	sb.WriteString(s.ExampleLabel + "\n\n")
	sb.WriteString(fence + "yaml\n")
	sb.WriteString(strings.TrimRight(s.Examples, "\n") + "\n")
	sb.WriteString(fence + "\n\n")
	sb.WriteString(s.RegistryNote + "\n\n")
}

func writeServiceConfigSection(sb *strings.Builder, fence, serviceConfigExample, customEnvExample string) {
	s := docs.Pages["configuration"].Sections.ServiceConfig
	sb.WriteString(s.Heading + "\n\n")
	sb.WriteString(s.Intro + "\n\n")
	sb.WriteString(s.EnvGeneratedLabel + "\n\n")
	sb.WriteString(fence + "bash\n")
	sb.WriteString("# " + serviceConfigExample + "\n")
	sb.WriteString(fence + "\n\n")
	sb.WriteString(s.CustomizingHeading + "\n\n")
	sb.WriteString(s.CustomizingIntro + "\n\n")
	sb.WriteString(fence + "bash\n")
	sb.WriteString(customEnvExample + "\n")
	sb.WriteString(fence + "\n\n")
	sb.WriteString(s.CustomizingNote + "\n\n")
}

func writeServiceMetadataSection(sb *strings.Builder, fence string) {
	s := docs.Pages["configuration"].Sections.ServiceMetadata
	sb.WriteString(s.Heading + "\n\n")
	sb.WriteString(s.Intro + "\n\n")
	sb.WriteString(s.ExampleLabel + "\n\n")
	sb.WriteString(fence + "yaml\n")
	sb.WriteString(strings.TrimRight(s.ExampleContent, "\n") + "\n")
	sb.WriteString(fence + "\n\n")
	sb.WriteString(s.Note + "\n\n")
}

func writeCompleteExampleSection(sb *strings.Builder, fence, completeExample, completeEnvExample string) {
	s := docs.Pages["configuration"].Sections.CompleteExample
	sb.WriteString(s.Heading + "\n\n")
	sb.WriteString(s.ConfigLabel + "\n\n")
	sb.WriteString(fence + "yaml\n")
	sb.WriteString(strings.TrimRight(completeExample, "\n") + "\n")
	sb.WriteString(fence + "\n\n")
	sb.WriteString(s.EnvLabel + "\n\n")
	sb.WriteString(fence + "bash\n")
	sb.WriteString(completeEnvExample + "\n")
	sb.WriteString(fence + "\n\n")
}

func writeNextStepsSection(sb *strings.Builder) {
	sb.WriteString(docs.Pages["configuration"].Sections.NextStepsSection + "\n\n")
	for _, link := range docs.Pages["configuration"].NextSteps {
		sb.WriteString(fmt.Sprintf("- **[%s](%s)** - %s\n", link.Label, link.URL, link.Description))
	}
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
		return &yaml.Node{Kind: yaml.ScalarNode, Value: docs.Examples.ProjectName}
	}
	if prop.defaultVal != "" && !prop.isTemplate {
		return &yaml.Node{Kind: yaml.ScalarNode, Value: prop.defaultVal}
	}
	return &yaml.Node{Kind: yaml.ScalarNode, Value: ""}
}

func buildArrayValueNode(sectionName string) *yaml.Node {
	seq := &yaml.Node{Kind: yaml.SequenceNode}
	if sectionName == schemaSectionStack {
		for _, s := range docs.Examples.Services {
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
	return envExample(svcMap, docs.Examples.CompleteServices, docs.Examples.EnvVarDisplayLimit,
		strings.ToUpper,
		func(svc loadedService, k string) string { return svc.config.Environment[k] },
	)
}

func generateCustomEnvExample(svcMap map[string]loadedService) string {
	return envExample(svcMap, docs.Examples.Services, docs.Examples.CustomEnvDisplayLimit,
		func(s string) string { return strings.ToUpper(s[:1]) + s[1:] },
		func(_ loadedService, _ string) string { return "my_custom_value" },
	)
}

func generateCompleteEnvExample(svcMap map[string]loadedService) string {
	return envExample(svcMap, docs.Examples.Services, docs.Examples.CustomEnvDisplayLimit,
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
		{Kind: yaml.ScalarNode, Value: docs.Examples.FullstackProjectName},
		{Kind: yaml.ScalarNode, Value: keyType},
		{Kind: yaml.ScalarNode, Value: docs.Examples.ProjectType},
	}}
	return []*yaml.Node{
		{Kind: yaml.ScalarNode, Value: schemaSectionProject},
		projectMapping,
	}
}

func stackExampleNodes() []*yaml.Node {
	stackSeq := &yaml.Node{Kind: yaml.SequenceNode}
	for _, s := range docs.Examples.CompleteServices {
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
