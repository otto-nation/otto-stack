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
	schemaSections := extractSchemaSections(schemaNode)

	page := docs.Pages[pageConfiguration]
	content := strings.Join([]string{
		"# " + page.Heading + "\n\n" + page.Intro + "\n\n",
		fileStructureSection(),
		mainConfigSection(generateConfigStructure(schemaSections)),
		configSchemaSections(schemaSections),
		sharingSection(),
		serviceConfigSection(generateServiceConfigExample(svcMap), generateCustomEnvExample(svcMap)),
		serviceMetadataSection(),
		completeExampleSection(generateCompleteExample(schemaNode), generateCompleteEnvExample(svcMap)),
		nextStepsSection(),
	}, "")

	return writePage(pageConfiguration, content)
}

func fileStructureSection() string {
	s := docs.ConfigSections.FileStructure
	return s.Heading + "\n\n" + s.Intro + "\n\n" +
		codeBlock("", docs.Pages[pageConfiguration].FileStructure)
}

func mainConfigSection(configStructure string) string {
	s := docs.ConfigSections.MainConfig
	return s.Heading + "\n\n" + s.FileLabel + "\n\n" + codeBlock("yaml", configStructure)
}

func sharingSection() string {
	s := docs.ConfigSections.Sharing
	var sb strings.Builder
	sb.WriteString(s.Heading + "\n\n" + s.Intro + "\n\n")
	for i, behavior := range s.Behaviors {
		fmt.Fprintf(&sb, "%d. %s\n", i+1, behavior)
	}
	sb.WriteString("\n" + s.ExampleLabel + "\n\n")
	sb.WriteString(codeBlock("yaml", s.Examples))
	sb.WriteString(s.RegistryNote + "\n\n")
	return sb.String()
}

func serviceConfigSection(serviceConfigExample, customEnvExample string) string {
	s := docs.ConfigSections.ServiceConfig
	return s.Heading + "\n\n" + s.Intro + "\n\n" + s.EnvGeneratedLabel + "\n\n" +
		codeBlock("bash", "# "+serviceConfigExample) +
		s.CustomizingHeading + "\n\n" + s.CustomizingIntro + "\n\n" +
		codeBlock("bash", customEnvExample) +
		s.CustomizingNote + "\n\n"
}

func serviceMetadataSection() string {
	s := docs.ConfigSections.ServiceMetadata
	return s.Heading + "\n\n" + s.Intro + "\n\n" + s.ExampleLabel + "\n\n" +
		codeBlock("yaml", s.ExampleContent) +
		s.Note + "\n\n"
}

func completeExampleSection(completeExample, completeEnvExample string) string {
	s := docs.ConfigSections.CompleteExample
	return s.Heading + "\n\n" + s.ConfigLabel + "\n\n" +
		codeBlock("yaml", completeExample) +
		s.EnvLabel + "\n\n" +
		codeBlock("bash", completeEnvExample)
}

func nextStepsSection() string {
	var sb strings.Builder
	sb.WriteString(docs.ConfigSections.NextStepsSection + "\n\n")
	for _, link := range docs.Pages[pageConfiguration].NextSteps {
		fmt.Fprintf(&sb, "- **[%s](%s)** - %s\n", link.Label, link.URL, link.Description)
	}
	return sb.String()
}

func extractSchemaSections(schemaNode *yaml.Node) []schemaSection {
	var sections []schemaSection
	eachEntry(schemaNode, func(key string, val *yaml.Node) {
		if section, ok := extractSchemaSection(key, val); ok {
			sections = append(sections, section)
		}
	})
	return sections
}

func extractSchemaSection(name string, sectionNode *yaml.Node) (schemaSection, bool) {
	propsNode := nodeGet(sectionNode, keyProperties)
	if propsNode == nil || propsNode.Kind != yaml.MappingNode {
		return schemaSection{}, false
	}
	section := schemaSection{
		name:        name,
		description: nodeGetStr(sectionNode, keyDescription),
	}
	eachEntry(propsNode, func(key string, val *yaml.Node) {
		section.properties = append(section.properties, extractSchemaProp(key, val))
	})
	return section, true
}

func extractSchemaProp(key string, propNode *yaml.Node) *schemaSectionProp {
	defaultVal := nodeGetStr(propNode, keyDefault)
	return &schemaSectionProp{
		key:         key,
		propType:    nodeGetStr(propNode, keyType),
		description: nodeGetStr(propNode, keyDescription),
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
		appendMappingEntry(mapping, section.name, sectionMapping)
	}
	result, err := marshalYAML(mapping)
	if err != nil {
		// marshalYAML on an in-memory node tree cannot fail in practice.
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
		appendMappingEntry(mapping, prop.key, valNode)
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
		return taggedScalarNode("!!bool", "false")
	}
	return nil
}

func buildStringValueNode(sectionName string, prop *schemaSectionProp) *yaml.Node {
	if sectionName == schemaSectionProject && prop.key == schemaPropName {
		return scalarNode(docs.Examples.ProjectName)
	}
	if prop.defaultVal != "" && !prop.isTemplate {
		return scalarNode(prop.defaultVal)
	}
	return scalarNode("")
}

func buildArrayValueNode(sectionName string) *yaml.Node {
	seq := &yaml.Node{Kind: yaml.SequenceNode}
	if sectionName == schemaSectionStack {
		for _, s := range docs.Examples.Services {
			seq.Content = append(seq.Content, scalarNode(s))
		}
	}
	return seq
}

// configSchemaSections renders schema.yaml sections as ### headings with property lists.
func configSchemaSections(sections []schemaSection) string {
	var sb strings.Builder
	for i, section := range sections {
		if i > 0 {
			sb.WriteString("\n\n")
		}
		fmt.Fprintf(&sb, "### %s\n\n%s\n\n", titleCase(section.name), section.description)
		for j, prop := range section.properties {
			if j > 0 {
				sb.WriteString("\n")
			}
			fmt.Fprintf(&sb, "- **%s**: %s", prop.key, prop.description)
		}
	}
	sb.WriteString("\n\n")
	return sb.String()
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
	var sb strings.Builder
	for _, name := range names {
		svc, ok := svcMap[name]
		if !ok || len(svc.config.Environment) == 0 {
			continue
		}
		sb.WriteString("# " + label(name) + "\n")
		keys := sortedEnvKeys(svc.config.Environment)
		n := min(limit, len(keys))
		for _, k := range keys[:n] {
			sb.WriteString(k + "=" + value(svc, k) + "\n")
		}
	}
	// Trim trailing newline — callers wrap output in a code block which adds its own.
	return strings.TrimRight(sb.String(), "\n")
}

func generateServiceConfigExample(svcMap map[string]loadedService) string {
	return envExample(svcMap, docs.Examples.CompleteServices, docs.Examples.EnvVarDisplayLimit,
		strings.ToUpper,
		func(svc loadedService, k string) string { return svc.config.Environment[k] },
	)
}

func generateCustomEnvExample(svcMap map[string]loadedService) string {
	return envExample(svcMap, docs.Examples.Services, docs.Examples.CustomEnvDisplayLimit,
		capitalizeFirst,
		func(_ loadedService, _ string) string { return docs.Examples.CustomEnvValue },
	)
}

func generateCompleteEnvExample(svcMap map[string]loadedService) string {
	return envExample(svcMap, docs.Examples.Services, docs.Examples.CustomEnvDisplayLimit,
		capitalizeFirst,
		func(_ loadedService, _ string) string { return docs.Examples.CompleteEnvValue },
	)
}

func generateCompleteExample(schemaNode *yaml.Node) string {
	mapping := buildCompleteExampleMapping(schemaNode)
	result, err := marshalYAML(mapping)
	if err != nil {
		// marshalYAML on an in-memory node tree cannot fail in practice.
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
	projectMapping := &yaml.Node{Kind: yaml.MappingNode}
	appendMappingEntry(projectMapping, keyName, scalarNode(docs.Examples.FullstackProjectName))
	appendMappingEntry(projectMapping, keyType, scalarNode(docs.Examples.ProjectType))
	return []*yaml.Node{scalarNode(schemaSectionProject), projectMapping}
}

func stackExampleNodes() []*yaml.Node {
	stackSeq := &yaml.Node{Kind: yaml.SequenceNode}
	for _, s := range docs.Examples.CompleteServices {
		stackSeq.Content = append(stackSeq.Content, scalarNode(s))
	}
	stackMapping := &yaml.Node{Kind: yaml.MappingNode}
	appendMappingEntry(stackMapping, "enabled", stackSeq)
	return []*yaml.Node{scalarNode(schemaSectionStack), stackMapping}
}

func validationExampleNodes() []*yaml.Node {
	optionsMapping := &yaml.Node{Kind: yaml.MappingNode}
	appendMappingEntry(optionsMapping, "config-syntax", taggedScalarNode("!!bool", "true"))
	appendMappingEntry(optionsMapping, "docker", taggedScalarNode("!!bool", "true"))
	validationMapping := &yaml.Node{Kind: yaml.MappingNode}
	appendMappingEntry(validationMapping, "options", optionsMapping)
	return []*yaml.Node{scalarNode(keyValidation), validationMapping}
}
