package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	pkgConfig "github.com/otto-nation/otto-stack/internal/pkg/config"
	"github.com/otto-nation/otto-stack/internal/pkg/constants"
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
	rawConfig, err := pkgConfig.LoadCommandConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	if err := generateConstants(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to generate constants: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated CLI code for %d flags and %d commands\n",
		countFlags(rawConfig), countCommands(rawConfig))
}

func generateConstants() error {
	rawConfig, err := pkgConfig.LoadCommandConfig()
	if err != nil {
		return fmt.Errorf("failed to load raw config: %w", err)
	}

	tmpl, err := template.ParseFiles(TemplateFilePath)
	if err != nil {
		return fmt.Errorf("failed to parse constants template: %w", err)
	}

	file, err := os.Create(GeneratedFilePath)
	if err != nil {
		return fmt.Errorf("failed to create constants file: %w", err)
	}
	defer func() { _ = file.Close() }()

	data := struct {
		Commands     []constantData
		Flags        []constantData
		Messages     []constantData
		Icons        []constantData
		CommandsData []commandData
	}{
		Commands:     collectCommandConstants(rawConfig),
		Flags:        collectFlagsData(rawConfig),
		Messages:     collectMessagesData(rawConfig),
		Icons:        collectIconsData(rawConfig),
		CommandsData: collectCommandsData(rawConfig),
	}

	return tmpl.Execute(file, data)
}

func countFlags(rawConfig map[string]any) int {
	flagNames := make(map[string]bool)
	addCommandFlags(rawConfig, flagNames)
	addGlobalFlags(rawConfig, flagNames)
	return len(flagNames)
}

func countCommands(rawConfig map[string]any) int {
	commands := getCommands(rawConfig)
	return len(commands)
}

func collectCommandConstants(rawConfig map[string]any) []constantData {
	var constants []constantData
	commands := getCommands(rawConfig)
	for cmdName := range commands {
		constants = append(constants, constantData{
			Name:  "Command" + toPascalCase(cmdName),
			Value: cmdName,
		})
	}
	return constants
}

func collectFlagsData(rawConfig map[string]any) []constantData {
	flagNames := make(map[string]bool)
	var flags []constantData

	addCommandFlags(rawConfig, flagNames)
	addGlobalFlags(rawConfig, flagNames)

	for flagName := range flagNames {
		flags = append(flags, constantData{
			Name:  "Flag" + toPascalCase(flagName),
			Value: flagName,
		})
	}
	return flags
}

func collectMessagesData(rawConfig map[string]any) []constantData {
	var messages []constantData
	msgs := getMessages(rawConfig)
	for category, categoryData := range msgs {
		categoryMap := getCategoryMap(categoryData)
		for key, value := range categoryMap {
			valueStr := getStringValue(value)
			if valueStr != "" {
				constName := "Msg" + toPascalCase(category) + "_" + strings.ReplaceAll(key, "-", "_")
				messages = append(messages, constantData{
					Name:  constName,
					Value: valueStr,
				})
			}
		}
	}
	return messages
}

func collectIconsData(rawConfig map[string]any) []constantData {
	var icons []constantData
	icns := getIcons(rawConfig)
	for category, categoryData := range icns {
		categoryMap := getCategoryMap(categoryData)
		for key, value := range categoryMap {
			valueStr := getStringValue(value)
			if valueStr != "" {
				constName := "Icon" + toPascalCase(strings.ReplaceAll(category+"."+key, ".", "_"))
				icons = append(icons, constantData{
					Name:  constName,
					Value: valueStr,
				})
			}
		}
	}
	return icons
}

func collectCommandsData(rawConfig map[string]any) []commandData {
	var commands []commandData
	cmds := getCommands(rawConfig)
	for cmdName, cmdData := range cmds {
		cmd := getCommandMap(cmdData)
		if cmd != nil {
			commands = append(commands, commandData{
				CommandName: cmdName,
				StructName:  toPascalCase(cmdName) + "Flags",
				FuncName:    "Parse" + toPascalCase(cmdName) + "Flags",
				Fields:      extractCommandFlags(cmd),
			})
		}
	}
	return commands
}

func extractCommandFlags(cmd map[string]any) []flagField {
	var fields []flagField
	flags := getFlags(cmd)
	for flagName, flagData := range flags {
		flag := getFlagMap(flagData)
		if flag != nil {
			flagType := getStringValue(flag["type"])
			if flagType == "" {
				flagType = "string"
			}
			fields = append(fields, flagField{
				Name:     toPascalCase(flagName),
				Type:     goTypeFromYAML(flagType),
				FlagName: flagName,
				YAMLType: flagType,
			})
		}
	}
	return fields
}

// Helper functions to reduce nesting
func getCommands(rawConfig map[string]any) map[string]any {
	if commands, ok := rawConfig["commands"].(map[string]any); ok {
		return commands
	}
	return make(map[string]any)
}

func getMessages(rawConfig map[string]any) map[string]any {
	if messages, ok := rawConfig["messages"].(map[string]any); ok {
		return messages
	}
	return make(map[string]any)
}

func getIcons(rawConfig map[string]any) map[string]any {
	if icons, ok := rawConfig["icons"].(map[string]any); ok {
		return icons
	}
	return make(map[string]any)
}

func getGlobal(rawConfig map[string]any) map[string]any {
	if global, ok := rawConfig["global"].(map[string]any); ok {
		return global
	}
	return make(map[string]any)
}

func getFlags(cmd map[string]any) map[string]any {
	if flags, ok := cmd["flags"].(map[string]any); ok {
		return flags
	}
	return make(map[string]any)
}

func getCommandMap(cmdData any) map[string]any {
	if cmd, ok := cmdData.(map[string]any); ok {
		return cmd
	}
	return nil
}

func getFlagMap(flagData any) map[string]any {
	if flag, ok := flagData.(map[string]any); ok {
		return flag
	}
	return nil
}

func getCategoryMap(categoryData any) map[string]any {
	if categoryMap, ok := categoryData.(map[string]any); ok {
		return categoryMap
	}
	return make(map[string]any)
}

func getStringValue(value any) string {
	if valueStr, ok := value.(string); ok {
		return valueStr
	}
	return ""
}

func addCommandFlags(rawConfig map[string]any, flagNames map[string]bool) {
	commands := getCommands(rawConfig)
	for _, cmdData := range commands {
		cmd := getCommandMap(cmdData)
		if cmd == nil {
			continue
		}
		cmdFlags := getFlags(cmd)
		for flagName := range cmdFlags {
			flagNames[flagName] = true
		}
	}
}

func addGlobalFlags(rawConfig map[string]any, flagNames map[string]bool) {
	global := getGlobal(rawConfig)
	globalFlags := getFlags(global)
	for flagName := range globalFlags {
		flagNames[flagName] = true
	}
}

func toPascalCase(s string) string {
	parts := strings.Split(s, "-")
	for i, part := range parts {
		if len(part) > 0 {
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
