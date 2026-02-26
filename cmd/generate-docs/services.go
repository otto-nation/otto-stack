package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

type serviceConfig struct {
	Name               string                `yaml:"name"`
	Description        string                `yaml:"description"`
	Hidden             bool                  `yaml:"hidden"`
	Environment        map[string]string     `yaml:"environment"`
	Documentation      *serviceDocumentation `yaml:"documentation"`
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

type categoryConfig struct {
	Icon  string `yaml:"icon"`
	Order int    `yaml:"order"`
}

func getCategoryConfig(name string) categoryConfig {
	if c, ok := docs.Categories[name]; ok {
		return c
	}
	return docs.Categories["other"]
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
		svc, err := loadService(path)
		if err != nil {
			return err
		}
		if svc != nil {
			services = append(services, *svc)
		}
		return nil
	})
	return services, err
}

func loadService(path string) (*loadedService, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}

	var svc serviceConfig
	if err := yaml.Unmarshal(data, &svc); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	if svc.Hidden {
		return nil, nil
	}

	var rootNode yaml.Node
	if err := yaml.Unmarshal(data, &rootNode); err != nil {
		return nil, fmt.Errorf("parse node %s: %w", path, err)
	}
	svc.configSchemaFields = extractSchemaFields(nodeGet(&rootNode, keyConfigSchema))

	ext := filepath.Ext(path)
	name := strings.TrimSuffix(filepath.Base(path), ext)
	return &loadedService{name: name, config: svc, category: inferCategory(path)}, nil
}

func inferCategory(path string) string {
	relPath, _ := filepath.Rel(servicesDirPath, path)
	parts := strings.Split(filepath.ToSlash(relPath), "/")
	if len(parts) >= 2 {
		return parts[0]
	}
	return "other"
}

func indexServices(services []loadedService) map[string]loadedService {
	svcMap := make(map[string]loadedService, len(services))
	for _, svc := range services {
		svcMap[svc.name] = svc
	}
	return svcMap
}

// ---- Generator: services-guide ----

func generateServicesGuide() error {
	services, err := loadAllServices()
	if err != nil {
		return fmt.Errorf("load services: %w", err)
	}

	byCategory := groupByCategory(services)
	categories := sortedCategories(byCategory)

	page := docs.Pages["services"]
	var sb strings.Builder
	sb.WriteString("# " + page.Heading + "\n\n")
	sb.WriteString(fmt.Sprintf(page.ServiceCount, len(services)) + "\n\n")
	sb.WriteString(page.Intro + "\n\n")

	for _, cat := range categories {
		catCfg := getCategoryConfig(cat)
		catTitle := strings.ToUpper(cat[:1]) + cat[1:]
		sb.WriteString(fmt.Sprintf("## %s %s\n\n", catCfg.Icon, catTitle))

		svcs := byCategory[cat]
		sort.Slice(svcs, func(i, j int) bool { return svcs[i].name < svcs[j].name })
		for _, svc := range svcs {
			sb.WriteString(renderServiceSection(svc))
		}
	}

	out, err := formatDocument(pageFM("services"), sb.String())
	if err != nil {
		return err
	}
	return writeOutput(pageOutput("services"), out)
}

func groupByCategory(services []loadedService) map[string][]loadedService {
	byCategory := make(map[string][]loadedService)
	for _, svc := range services {
		byCategory[svc.category] = append(byCategory[svc.category], svc)
	}
	return byCategory
}

func sortedCategories(byCategory map[string][]loadedService) []string {
	categories := make([]string, 0, len(byCategory))
	for cat := range byCategory {
		categories = append(categories, cat)
	}
	sort.Slice(categories, func(i, j int) bool {
		return getCategoryConfig(categories[i]).Order < getCategoryConfig(categories[j]).Order
	})
	return categories
}

func renderServiceSection(svc loadedService) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("### %s\n\n%s\n\n", svc.name, svc.config.Description))

	if len(svc.config.configSchemaFields) > 0 {
		sb.WriteString(renderServiceSchemaSection(svc.config.configSchemaFields))
	}
	if svc.config.Documentation != nil {
		sb.WriteString(renderServiceDocumentation(svc.config.Documentation))
	}

	sb.WriteString("---\n\n")
	return sb.String()
}

func renderServiceSchemaSection(fields []*schemaField) string {
	sections := docs.ServicesSections
	var sb strings.Builder
	sb.WriteString(sections.ConfigOptions + "\n\n")
	for _, field := range fields {
		sb.WriteString(renderSchemaField(field, "####"))
	}

	examplesNode := buildExamplesNode(fields)
	if examplesNode == nil {
		return sb.String()
	}
	exYAML, err := marshalYAML(examplesNode)
	if err != nil || strings.TrimSpace(exYAML) == "" {
		return sb.String()
	}
	sb.WriteString("\n" + sections.ExampleConfig + "\n\n")
	codeBlock(&sb, "yaml", exYAML)
	return sb.String()
}

func renderServiceDocumentation(svcDocs *serviceDocumentation) string {
	sections := docs.ServicesSections
	var sb strings.Builder
	if len(svcDocs.UseCases) > 0 {
		sb.WriteString(sections.UseCases + "\n\n")
		for _, uc := range svcDocs.UseCases {
			sb.WriteString(fmt.Sprintf("- %s\n\n", uc))
		}
	}
	if len(svcDocs.Examples) > 0 {
		sb.WriteString(sections.ExamplesHeading + "\n\n")
		for _, ex := range svcDocs.Examples {
			codeBlock(&sb, "bash", ex)
		}
	}
	return sb.String()
}

func renderSchemaField(field *schemaField, headingLevel string) string {
	labels := docs.Labels
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s %s\n\n", headingLevel, field.Name))
	if field.Description != "" {
		sb.WriteString(field.Description + "\n\n")
	}
	sb.WriteString(fmt.Sprintf("- %s: `%s`\n", labels.FieldType, field.Type))
	if field.Default != nil && field.Default.Value != "" {
		sb.WriteString(fmt.Sprintf("- %s: `%s`\n", labels.FieldDefault, field.Default.Value))
	}
	if field.Required {
		sb.WriteString("- " + labels.FieldRequiredYes + "\n")
	}
	sb.WriteString("\n")

	if field.Items != nil && len(field.Items.Properties) > 0 {
		sb.WriteString(labels.Items + "\n\n")
		for _, itemProp := range field.Items.Properties {
			sb.WriteString(renderItemProperty(itemProp))
		}
		sb.WriteString("\n")
	}
	if len(field.Properties) > 0 {
		sb.WriteString(labels.Properties + "\n\n")
		for _, subProp := range field.Properties {
			sb.WriteString(renderItemProperty(subProp))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func renderItemProperty(p *schemaField) string {
	labels := docs.Labels
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("- **%s** (`%s`)", p.Name, p.Type))
	if p.Required {
		sb.WriteString(" " + labels.FieldRequiredIndicator)
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
