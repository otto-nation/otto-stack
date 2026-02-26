package main

import "gopkg.in/yaml.v3"

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
func extractSchemaFields(schemaNode *yaml.Node) []*schemaField {
	if schemaNode == nil {
		return nil
	}
	propsNode := nodeGet(schemaNode, keyProperties)
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
		fields = append(fields, extractField(propsNode.Content[i].Value, propsNode.Content[i+1]))
	}
	return fields
}

func extractField(name string, valNode *yaml.Node) *schemaField {
	field := &schemaField{
		Name:        name,
		Type:        nodeStr(nodeGet(valNode, keyType)),
		Description: nodeStr(nodeGet(valNode, keyDescription)),
		Default:     nodeGet(valNode, keyDefault),
	}
	if reqNode := nodeGet(valNode, keyRequired); reqNode != nil {
		field.Required = nodeBool(reqNode)
	}
	if itemsNode := nodeGet(valNode, keyItems); itemsNode != nil {
		field.Items = &schemaItems{
			Type:       nodeStr(nodeGet(itemsNode, keyType)),
			Properties: extractPropertiesNode(nodeGet(itemsNode, keyProperties)),
		}
	}
	if subPropsNode := nodeGet(valNode, keyProperties); subPropsNode != nil {
		field.Properties = extractPropertiesNode(subPropsNode)
	}
	return field
}

// buildExamplesNode builds an ordered yaml.Node representing example configuration.
func buildExamplesNode(fields []*schemaField) *yaml.Node {
	mapping := &yaml.Node{Kind: yaml.MappingNode}
	for _, f := range fields {
		valNode := buildFieldExampleNode(f)
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

func buildFieldExampleNode(f *schemaField) *yaml.Node {
	switch f.Type {
	case "string":
		if f.Default != nil && f.Default.Value != "" {
			return &yaml.Node{Kind: yaml.ScalarNode, Value: f.Default.Value}
		}
	case "integer":
		if f.Default != nil {
			return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!int", Value: f.Default.Value}
		}
	case "boolean":
		if f.Default != nil {
			return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!bool", Value: f.Default.Value}
		}
	case "array":
		if f.Items != nil {
			return buildItemSequenceNode(f.Items)
		}
	case "object":
		if len(f.Properties) > 0 {
			return buildObjectExampleNode(f.Properties)
		}
	}
	return nil
}

func buildItemSequenceNode(items *schemaItems) *yaml.Node {
	itemNode := buildItemExampleNode(items)
	if itemNode == nil {
		return nil
	}
	return &yaml.Node{Kind: yaml.SequenceNode, Content: []*yaml.Node{itemNode}}
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
		valNode := buildObjectPropNode(p)
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

func buildObjectPropNode(p *schemaField) *yaml.Node {
	if p.Default != nil {
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: p.Default.Tag, Value: p.Default.Value}
	}
	switch p.Type {
	case "string":
		return &yaml.Node{Kind: yaml.ScalarNode, Value: "example-" + p.Name}
	case "integer":
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!int", Value: "1"}
	case "boolean":
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!bool", Value: "true"}
	}
	return nil
}
