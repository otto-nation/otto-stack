package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

const (
	CommandsYAMLPath  = "internal/config/commands.yaml"
	TemplateFilePath  = "cmd/generate-cli/templates/core.tmpl"
	GeneratedFilePath = "internal/core/cli_generated.go"
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
	Key   string
}

type flagField struct {
	Name     string
	Type     string
	FlagName string
	YAMLType string
}

func main() {
	rawConfig, err := loadCommandConfig()
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

func loadCommandConfig() (map[string]any, error) {
	data, err := os.ReadFile(CommandsYAMLPath)
	if err != nil {
		return nil, err
	}

	var config map[string]any
	err = yaml.Unmarshal(data, &config)
	return config, err
}

func generateConstants() error {
	rawConfig, err := loadCommandConfig()
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
			Value: strconv.Quote(cmdName),
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
			Value: strconv.Quote(flagName),
		})
	}
	return flags
}

func collectMessagesData(rawConfig map[string]any) []constantData {
	var messages []constantData
	msgs := getMessages(rawConfig)
	for category, categoryData := range msgs {
		messages = append(messages, processMessageCategory(category, categoryData)...)
	}
	return messages
}

func processMessageCategory(category string, categoryData any) []constantData {
	var messages []constantData
	categoryMap := getCategoryMap(categoryData)
	for key, value := range categoryMap {
		if msg := createMessageConstant(category, key, value); msg.Name != "" {
			messages = append(messages, msg)
		}
	}
	return messages
}

func createMessageConstant(category, key string, value any) constantData {
	valueStr := getStringValue(value)
	if valueStr == "" {
		return constantData{}
	}

	constName := "Msg" + toPascalCase(category) + "_" + strings.ReplaceAll(key, "-", "_")
	return constantData{
		Name:  constName,
		Value: strconv.Quote(valueStr),
	}
}

func collectIconsData(rawConfig map[string]any) []constantData {
	var icons []constantData
	icns := getIcons(rawConfig)
	for category, categoryData := range icns {
		icons = append(icons, processIconCategory(category, categoryData)...)
	}
	return icons
}

func processIconCategory(category string, categoryData any) []constantData {
	var icons []constantData
	categoryMap := getCategoryMap(categoryData)
	for key, value := range categoryMap {
		if icon := createIconConstant(category, key, value); icon.Name != "" {
			icons = append(icons, icon)
		}
	}
	return icons
}

func createIconConstant(category, key string, value any) constantData {
	valueStr := getStringValue(value)
	if valueStr == "" {
		return constantData{}
	}

	// Generate snake_case names to match existing code expectations
	constName := "Icon" + toPascalCase(category) + "_" + strings.ToLower(strings.ReplaceAll(key, "-", "_"))
	return constantData{
		Name:  constName,
		Value: strconv.Quote(valueStr),
		Key:   category + "_" + key,
	}
}

func collectCommandsData(rawConfig map[string]any) []commandData {
	var commands []commandData
	cmds := getCommands(rawConfig)
	for cmdName, cmdData := range cmds {
		if cmd := createCommandData(cmdName, cmdData); cmd.CommandName != "" {
			commands = append(commands, cmd)
		}
	}
	return commands
}

func createCommandData(cmdName string, cmdData any) commandData {
	cmd := getCommandMap(cmdData)
	if cmd == nil {
		return commandData{}
	}

	return commandData{
		CommandName: cmdName,
		StructName:  toPascalCase(cmdName) + "Flags",
		FuncName:    "Parse" + toPascalCase(cmdName) + "Flags",
		Fields:      extractCommandFlags(cmd),
	}
}

func extractCommandFlags(cmd map[string]any) []flagField {
	var fields []flagField
	flags := getFlags(cmd)
	for flagName, flagData := range flags {
		if field := createFlagField(flagName, flagData); field.Name != "" {
			fields = append(fields, field)
		}
	}
	return fields
}

const (
	defaultFlagType = "string"
)

func createFlagField(flagName string, flagData any) flagField {
	flag := getFlagMap(flagData)
	if flag == nil {
		return flagField{}
	}

	flagType := getStringValue(flag["type"])
	if flagType == "" {
		flagType = defaultFlagType
	}

	return flagField{
		Name:     toPascalCase(flagName),
		Type:     goTypeFromYAML(flagType),
		FlagName: flagName,
		YAMLType: flagType,
	}
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
	case defaultFlagType:
		return defaultFlagType
	case "int":
		return "int"
	case "bool":
		return "bool"
	case "stringArray":
		return "stringArray"
	default:
		return defaultFlagType
	}
}
