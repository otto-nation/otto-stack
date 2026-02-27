package main

import (
	"bytes"
	"strings"

	"gopkg.in/yaml.v3"
)

// YAML key constants used when accessing yaml.Node fields.
const (
	keyAliases         = "aliases"
	keyCategories      = "categories"
	keyCommand         = "command"
	keyCommands        = "commands"
	keyConfigSchema    = "configuration_schema"
	keyDefault         = "default"
	keyDescription     = "description"
	keyExamples        = "examples"
	keyFlags           = "flags"
	keyGlobalFlags     = "global_flags"
	keyIcon            = "icon"
	keyItems           = "items"
	keyLongDescription = "long_description"
	keyMetadata        = "metadata"
	keyName            = "name"
	keyOptions         = "options"
	keyProperties      = "properties"
	keyRelatedCommands = "related_commands"
	keyRequired        = "required"
	keySchema          = "schema"
	keyShort           = "short"
	keyTips            = "tips"
	keyType            = "type"
	keyUsage           = "usage"
	keyValidation      = "validation"
)

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
	var keys []string
	eachEntry(n, func(key string, _ *yaml.Node) {
		keys = append(keys, key)
	})
	return keys
}

// nodeGetStr returns the string value of key in a mapping node.
func nodeGetStr(n *yaml.Node, key string) string {
	return nodeStr(nodeGet(n, key))
}

// eachEntry calls f for each key-value pair in a mapping node, in document order.
func eachEntry(n *yaml.Node, f func(key string, val *yaml.Node)) {
	n = nodeDoc(n)
	if n == nil || n.Kind != yaml.MappingNode {
		return
	}
	for i := 0; i+1 < len(n.Content); i += 2 {
		f(n.Content[i].Value, n.Content[i+1])
	}
}

// appendMappingEntry appends a key-value pair to a YAML mapping node.
func appendMappingEntry(m *yaml.Node, key string, val *yaml.Node) {
	m.Content = append(m.Content, scalarNode(key), val)
}

// scalarNode returns a YAML scalar node with the given value.
func scalarNode(value string) *yaml.Node {
	return &yaml.Node{Kind: yaml.ScalarNode, Value: value}
}

// taggedScalarNode returns a YAML scalar node with the given tag and value.
func taggedScalarNode(tag, value string) *yaml.Node {
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: tag, Value: value}
}

// capitalizeFirst uppercases the first character of s.
func capitalizeFirst(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
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

// quoteEach wraps each string in backticks.
func quoteEach(ss []string) []string {
	out := make([]string, len(ss))
	for i, s := range ss {
		out[i] = "`" + s + "`"
	}
	return out
}
