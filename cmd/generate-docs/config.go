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
		htmlComment(
			"\u26a0\ufe0f  AUTO-GENERATED FILE - DO NOT EDIT DIRECTLY",
			"This file is generated from "+schemaYAMLPath+" and "+docsConfigPath,
			"To make changes, edit source files and run: task generate:docs",
		),
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

	out, err := formatDocument(pageFM(pageConfiguration), content)
	if err != nil {
		return err
	}
	return writeOutput(pageOutput(pageConfiguration), out)
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
	sb.WriteString(s.Heading + "\n\n" + s.Intro + "\n")
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

// titleLabel capitalises the first letter of s, used as an env var comment header.
func titleLabel(s string) string {
	return strings.ToUpper(s[:1]) + s[1:]
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
		n := limit
		if len(keys) < n {
			n = len(keys)
		}
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
		titleLabel,
		func(_ loadedService, _ string) string { return docs.Examples.CustomEnvValue },
	)
}

func generateCompleteEnvExample(svcMap map[string]loadedService) string {
	return envExample(svcMap, docs.Examples.Services, docs.Examples.CustomEnvDisplayLimit,
		titleLabel,
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
