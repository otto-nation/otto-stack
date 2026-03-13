package main

import "gopkg.in/yaml.v3"

const (
	fieldTypeString  = "string"
	fieldTypeBoolean = "boolean"
	fieldTypeInteger = "integer"
	fieldTypeArray   = "array"
	fieldTypeObject  = "object"
)

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
	return extractPropertiesNode(nodeGet(schemaNode, keyProperties))
}

func extractPropertiesNode(propsNode *yaml.Node) []*schemaField {
	var fields []*schemaField
	eachEntry(propsNode, func(key string, val *yaml.Node) {
		fields = append(fields, extractField(key, val))
	})
	return fields
}

func extractField(name string, valNode *yaml.Node) *schemaField {
	field := &schemaField{
		Name:        name,
		Type:        nodeGetStr(valNode, keyType),
		Description: nodeGetStr(valNode, keyDescription),
		Default:     nodeGet(valNode, keyDefault),
	}
	if reqNode := nodeGet(valNode, keyRequired); reqNode != nil {
		field.Required = nodeBool(reqNode)
	}
	if itemsNode := nodeGet(valNode, keyItems); itemsNode != nil {
		field.Items = &schemaItems{
			Type:       nodeGetStr(itemsNode, keyType),
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
		appendMappingEntry(mapping, f.Name, valNode)
	}
	if len(mapping.Content) == 0 {
		return nil
	}
	return mapping
}

func buildFieldExampleNode(f *schemaField) *yaml.Node {
	switch f.Type {
	case fieldTypeString:
		if f.Default != nil && f.Default.Value != "" {
			return scalarNode(f.Default.Value)
		}
	case fieldTypeInteger:
		if f.Default != nil {
			return taggedScalarNode("!!int", f.Default.Value)
		}
	case fieldTypeBoolean:
		if f.Default != nil {
			return taggedScalarNode("!!bool", f.Default.Value)
		}
	case fieldTypeArray:
		if f.Items != nil {
			return buildItemSequenceNode(f.Items)
		}
	case fieldTypeObject:
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
		appendMappingEntry(mapping, p.Name, valNode)
	}
	return mapping
}

func buildObjectPropNode(p *schemaField) *yaml.Node {
	if p.Default != nil {
		return taggedScalarNode(p.Default.Tag, p.Default.Value)
	}
	switch p.Type {
	case fieldTypeString:
		return scalarNode("example-" + p.Name)
	case fieldTypeInteger:
		return taggedScalarNode("!!int", "1")
	case fieldTypeBoolean:
		return taggedScalarNode("!!bool", "true")
	}
	return nil
}
