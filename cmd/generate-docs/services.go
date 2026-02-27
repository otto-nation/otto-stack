package main

import (
	"fmt"
	"sort"
	"strings"
)

func generateServicesGuide() error {
	services, err := loadAllServices()
	if err != nil {
		return fmt.Errorf("load services: %w", err)
	}

	byCategory := groupByCategory(services)
	categories := sortedCategories(byCategory)

	page := docs.Pages[pageServices]
	content := strings.Join([]string{
		"# " + page.Heading + "\n\n",
		fmt.Sprintf(page.ServiceCount, len(services)) + "\n\n",
		page.Intro + "\n\n",
		categoriesSection(byCategory, categories),
	}, "")

	out, err := formatDocument(pageFM(pageServices), content)
	if err != nil {
		return err
	}
	return writeOutput(pageOutput(pageServices), out)
}

func categoriesSection(byCategory map[string][]loadedService, categories []string) string {
	var sb strings.Builder
	for _, cat := range categories {
		catCfg := getCategoryConfig(cat)
		catTitle := strings.ToUpper(cat[:1]) + cat[1:]
		fmt.Fprintf(&sb, "## %s %s\n\n", catCfg.Icon, catTitle)

		svcs := byCategory[cat]
		sort.Slice(svcs, func(i, j int) bool { return svcs[i].name < svcs[j].name })
		for _, svc := range svcs {
			sb.WriteString(renderServiceSection(svc))
		}
	}
	return sb.String()
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
	sb.WriteString(codeBlock("yaml", exYAML))
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
			sb.WriteString(codeBlock("bash", ex))
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
