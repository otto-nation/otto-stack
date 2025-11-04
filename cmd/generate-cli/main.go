package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/otto-nation/otto-stack/internal/config"
	pkgConfig "github.com/otto-nation/otto-stack/internal/pkg/config"
	"github.com/otto-nation/otto-stack/internal/pkg/constants"
	"gopkg.in/yaml.v3"
)

const (
	TemplateFilePath  = "cmd/generate-cli/templates/constants.tmpl"
	GeneratedFilePath = "internal/pkg/constants/cli_generated.go"
)

type commandData struct {
	CommandName string
	StructName  string
	FuncName    string
	Fields      []flagField
}

type constantData struct {
	Name  string
	Value string
}

type flagField struct {
	Name     string
	Type     string
	FlagName string
	YAMLType string
}

func main() {
	commandConfig, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	if err := generateConstants(commandConfig); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to generate constants: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated CLI code for %d flags and %d commands\n",
		countFlags(commandConfig), len(commandConfig.Commands))
}

func generateConstants(config *pkgConfig.CommandConfig) error {
	tmpl, err := template.ParseFiles(TemplateFilePath)
	if err != nil {
		return fmt.Errorf("failed to parse constants template: %w", err)
	}

	file, err := os.Create(GeneratedFilePath)
	if err != nil {
		return fmt.Errorf("failed to create constants file: %w", err)
	}
	defer func() { _ = file.Close() }()

	// Collect command data for flag parsing
	var commands []commandData
	for cmdName, cmd := range config.Commands {
		commands = append(commands, commandData{
			CommandName: cmdName,
			StructName:  toPascalCase(cmdName) + "Flags",
			FuncName:    "Parse" + toPascalCase(cmdName) + "Flags",
			Fields:      collectCommandFlags(cmd),
		})
	}

	data := struct {
		Commands     []constantData
		Flags        []constantData
		Messages     []constantData
		Icons        []constantData
		CommandsData []commandData
	}{
		Commands:     collectCommandsData(config),
		Flags:        collectFlagsData(config),
		Messages:     collectMessagesData(config),
		Icons:        collectIconsData(config),
		CommandsData: commands,
	}

	return tmpl.Execute(file, data)
}

func loadConfig() (*pkgConfig.CommandConfig, error) {
	var commandConfig pkgConfig.CommandConfig
	if err := yaml.Unmarshal(config.EmbeddedCommandsYAML, &commandConfig); err != nil {
		return nil, fmt.Errorf("failed to parse commands.yaml: %w", err)
	}
	return &commandConfig, nil
}

func countFlags(config *pkgConfig.CommandConfig) int {
	flagNames := make(map[string]bool)
	for _, cmd := range config.Commands {
		for flagName := range cmd.Flags {
			flagNames[flagName] = true
		}
	}
	for flagName := range config.Global.Flags {
		flagNames[flagName] = true
	}
	return len(flagNames)
}

func collectCommandsData(config *pkgConfig.CommandConfig) []constantData {
	var commands []constantData
	for cmdName := range config.Commands {
		commands = append(commands, constantData{
			Name:  "Command" + toPascalCase(cmdName),
			Value: cmdName,
		})
	}
	return commands
}

func collectFlagsData(config *pkgConfig.CommandConfig) []constantData {
	flagNames := make(map[string]bool)
	var flags []constantData

	for _, cmd := range config.Commands {
		for flagName := range cmd.Flags {
			if !flagNames[flagName] {
				flagNames[flagName] = true
				flags = append(flags, constantData{
					Name:  "Flag" + toPascalCase(flagName),
					Value: flagName,
				})
			}
		}
	}

	for flagName := range config.Global.Flags {
		if !flagNames[flagName] {
			flagNames[flagName] = true
			flags = append(flags, constantData{
				Name:  "Flag" + toPascalCase(flagName),
				Value: flagName,
			})
		}
	}

	return flags
}

func collectMessagesData(config *pkgConfig.CommandConfig) []constantData {
	messageMap := make(map[string]string)
	if config.Messages != nil {
		for category, categoryData := range config.Messages {
			categoryMap := categoryData.(map[string]any)
			for key, value := range categoryMap {
				constName := "Msg" + toPascalCase(category) + "_" + strings.ReplaceAll(key, "-", "_")
				messageMap[constName] = value.(string)
			}
		}
	}

	var messages []constantData
	for name, value := range messageMap {
		messages = append(messages, constantData{
			Name:  name,
			Value: value,
		})
	}
	return messages
}

func collectIconsData(config *pkgConfig.CommandConfig) []constantData {
	iconMap := make(map[string]string)
	if config.Icons != nil {
		for category, categoryData := range config.Icons {
			categoryMap := categoryData.(map[string]any)
			for key, value := range categoryMap {
				constName := "Icon" + toPascalCase(strings.ReplaceAll(category+"."+key, ".", "_"))
				iconMap[constName] = value.(string)
			}
		}
	}

	var icons []constantData
	for name, value := range iconMap {
		icons = append(icons, constantData{
			Name:  name,
			Value: value,
		})
	}
	return icons
}

func collectCommandFlags(cmd pkgConfig.Command) []flagField {
	var flagFields []flagField

	for flagName, flag := range cmd.Flags {
		flagFields = append(flagFields, flagField{
			Name:     toPascalCase(flagName),
			Type:     goTypeFromYAML(flag.Type),
			FlagName: flagName,
			YAMLType: flag.Type,
		})
	}

	return flagFields
}

func toPascalCase(s string) string {
	parts := strings.Split(s, "-")
	for i, part := range parts {
		if len(part) > 0 {
			// Handle special cases
			switch strings.ToLower(part) {
			case "json":
				parts[i] = "JSON"
			case "yaml":
				parts[i] = "YAML"
			case "xml":
				parts[i] = "XML"
			case "http":
				parts[i] = "HTTP"
			case "https":
				parts[i] = "HTTPS"
			case "tty":
				parts[i] = "TTY"
			case "url":
				parts[i] = "URL"
			case "api":
				parts[i] = "API"
			default:
				parts[i] = strings.ToUpper(part[:1]) + part[1:]
			}
		}
	}
	return strings.Join(parts, "")
}

func goTypeFromYAML(yamlType string) string {
	switch yamlType {
	case constants.TypeString:
		return constants.TypeString
	case constants.TypeInt:
		return constants.TypeInt
	case constants.TypeBool:
		return constants.TypeBool
	case constants.TypeStringArray:
		return constants.TypeStringArray
	default:
		return constants.TypeString
	}
}
