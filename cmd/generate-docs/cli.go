package main

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func generateCLIReference() error {
	data, err := os.ReadFile(commandsYAMLPath)
	if err != nil {
		return fmt.Errorf("read commands.yaml: %w", err)
	}

	var rootNode yaml.Node
	if err := yaml.Unmarshal(data, &rootNode); err != nil {
		return fmt.Errorf("parse commands.yaml node: %w", err)
	}

	metadataNode := nodeGet(&rootNode, keyMetadata)
	description := nodeStr(nodeGet(metadataNode, keyDescription))
	if description == "" {
		description = "A powerful development stack management tool for streamlined local development automation"
	}

	var sb strings.Builder
	sb.WriteString("<!--\n")
	sb.WriteString("  \u26a0\ufe0f  AUTO-GENERATED FILE - DO NOT EDIT DIRECTLY\n")
	sb.WriteString("  This file is generated from internal/config/commands.yaml\n")
	sb.WriteString("  To make changes, edit the source file and run: task generate:docs\n")
	sb.WriteString("-->\n\n")
	sb.WriteString("# otto-stack CLI Reference\n\n")
	sb.WriteString(description + "\n\n")

	sb.WriteString("## Command Categories\n\n")
	categoriesNode := nodeGet(&rootNode, keyCategories)
	for _, catKey := range nodeKeys(categoriesNode) {
		catNode := nodeGet(categoriesNode, catKey)
		icon := nodeStr(nodeGet(catNode, keyIcon))
		name := nodeStr(nodeGet(catNode, keyName))
		desc := nodeStr(nodeGet(catNode, keyDescription))
		cmds := nodeStringSlice(nodeGet(catNode, keyCommands))
		sb.WriteString(fmt.Sprintf("### %s %s\n\n%s\n\n**Commands:** %s\n\n",
			icon, name, desc, strings.Join(quoteEach(cmds), ", ")))
	}

	sb.WriteString("## Commands\n\n")
	commandsNode := nodeGet(&rootNode, keyCommands)
	for _, cmdKey := range nodeKeys(commandsNode) {
		sb.WriteString(renderCommandSection(cmdKey, nodeGet(commandsNode, cmdKey)))
	}

	globalFlagsNode := nodeGet(&rootNode, keyGlobalFlags)
	if globalFlagsNode != nil {
		sb.WriteString("## Global Flags\n\nThese flags are available for all commands:\n\n")
		sb.WriteString(renderFlagLines(globalFlagsNode))
	}

	out, err := formatDocument(pageFM("cli-reference"), sb.String())
	if err != nil {
		return err
	}
	return writeOutput(pageOutput("cli-reference"), out)
}

func renderCommandSection(name string, cmdNode *yaml.Node) string {
	var sb strings.Builder

	desc := nodeStr(nodeGet(cmdNode, keyDescription))
	longDesc := nodeStr(nodeGet(cmdNode, keyLongDescription))
	usage := nodeStr(nodeGet(cmdNode, keyUsage))
	aliases := nodeStringSlice(nodeGet(cmdNode, keyAliases))

	sb.WriteString(fmt.Sprintf("### `%s`\n\n%s\n\n", name, desc))
	if longDesc != "" {
		sb.WriteString(strings.TrimSpace(longDesc) + "\n\n")
	}
	if usage != "" {
		sb.WriteString(fmt.Sprintf("**Usage:** `otto-stack %s`\n\n", usage))
	}
	if len(aliases) > 0 {
		sb.WriteString("**Aliases:** " + strings.Join(quoteEach(aliases), ", ") + "\n\n")
	}

	sb.WriteString(renderCommandExamples(nodeGet(cmdNode, keyExamples)))
	sb.WriteString(renderCommandFlags(nodeGet(cmdNode, keyFlags)))
	sb.WriteString(renderCommandRelated(nodeGet(cmdNode, keyRelatedCommands)))
	sb.WriteString(renderCommandTips(nodeGet(cmdNode, keyTips)))

	return sb.String()
}

func renderCommandExamples(examplesNode *yaml.Node) string {
	if examplesNode == nil || examplesNode.Kind != yaml.SequenceNode {
		return ""
	}
	var sb strings.Builder
	sb.WriteString("**Examples:**\n\n")
	for _, exNode := range examplesNode.Content {
		cmd := nodeStr(nodeGet(exNode, keyCommand))
		desc := nodeStr(nodeGet(exNode, keyDescription))
		sb.WriteString(fmt.Sprintf("```bash\n%s\n```\n\n", cmd))
		if desc != "" {
			sb.WriteString(desc + "\n\n")
		}
	}
	return sb.String()
}

func renderCommandFlags(flagsNode *yaml.Node) string {
	if flagsNode == nil || len(nodeKeys(flagsNode)) == 0 {
		return ""
	}
	return "**Flags:**\n\n" + renderFlagLines(flagsNode)
}

func renderFlagLines(flagsNode *yaml.Node) string {
	var sb strings.Builder
	for _, flagKey := range nodeKeys(flagsNode) {
		sb.WriteString(renderFlagLine(flagKey, nodeGet(flagsNode, flagKey)))
	}
	sb.WriteString("\n")
	return sb.String()
}

func renderFlagLine(flagKey string, flagNode *yaml.Node) string {
	short := nodeStr(nodeGet(flagNode, keyShort))
	flagType := nodeStr(nodeGet(flagNode, keyType))
	desc := nodeStr(nodeGet(flagNode, keyDescription))
	defaultNode := nodeGet(flagNode, keyDefault)
	optionsNode := nodeGet(flagNode, keyOptions)

	line := fmt.Sprintf("- `--%s`", flagKey)
	if short != "" {
		line += fmt.Sprintf(", `-%s`", short)
	}
	if flagType != "" {
		line += fmt.Sprintf(" (`%s`)", flagType)
	}
	line += ": " + desc
	if defaultNode != nil {
		line += fmt.Sprintf(" (default: `%s`)", defaultNode.Value)
	}
	if optionsNode != nil {
		line += " (options: " + strings.Join(quoteEach(nodeStringSlice(optionsNode)), ", ") + ")"
	}
	return line + "\n"
}

func renderCommandRelated(relatedNode *yaml.Node) string {
	related := nodeStringSlice(relatedNode)
	if len(related) == 0 {
		return ""
	}
	links := make([]string, len(related))
	for i, r := range related {
		links[i] = fmt.Sprintf("[`%s`](#%s)", r, r)
	}
	return "**Related Commands:** " + strings.Join(links, ", ") + "\n\n"
}

func renderCommandTips(tipsNode *yaml.Node) string {
	tips := nodeStringSlice(tipsNode)
	if len(tips) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteString("**Tips:**\n\n")
	for _, tip := range tips {
		sb.WriteString("- " + tip + "\n")
	}
	sb.WriteString("\n")
	return sb.String()
}
