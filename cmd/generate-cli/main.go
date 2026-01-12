package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/otto-nation/otto-stack/cmd/codegen"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
)

const (
	CommandsYAMLPath     = "internal/config/commands.yaml"
	CoreTemplateFilePath = "cmd/generate-cli/templates/core.tmpl"
	CLITemplateFilePath  = "cmd/generate-cli/templates/ci.tmpl"
	CLICommandsTemplate  = "cmd/generate-cli/templates/cli.tmpl"
	CoreGeneratedPath    = "internal/core/constants_generated.go"
	CIGeneratedPath      = "internal/pkg/ci/generated.go"
	CLICommandsPath      = "internal/pkg/cli/cli_generated.go"
	CoreTemplateName     = "core"
)

const (
	// YAML keys
	KeyCommands    = "commands"
	KeyMessages    = "messages"
	KeyIcons       = "icons"
	KeyGlobal      = "global"
	KeyFlags       = "flags"
	KeyValidation  = "validation"
	KeyType        = "type"
	KeyDescription = "description"
)

const (
	// Type strings
	TypeString      = "string"
	TypeInt         = "int"
	TypeBool        = "bool"
	TypeStringArray = "stringArray"
)

const (
	// Prefix strings
	PrefixCommand = "Command"
	PrefixFlag    = "Flag"
	PrefixMsg     = "Msg"
	PrefixIcon    = "Icon"
	PrefixParse   = "Parse"
	SuffixFlags   = "Flags"
)

const (
	// Special case strings for toPascalCase
	CaseJSON  = "json"
	CaseYAML  = "yaml"
	CaseXML   = "xml"
	CaseHTTP  = "http"
	CaseHTTPS = "https"
	CaseTTY   = "tty"
	CaseURL   = "url"
	CaseAPI   = "api"
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

type validationOption struct {
	Key         string
	Description string
	Required    bool
}

func main() {
	rawConfig, err := loadCommandConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	if err := generateCoreConstants(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to generate core constants: %v\n", err)
		os.Exit(1)
	}

	if err := generateCIUtilities(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to generate CI utilities: %v\n", err)
		os.Exit(1)
	}

	if err := generateCLICommands(rawConfig); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to generate CLI commands: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated CLI code for %d flags and %d commands\n",
		countFlags(rawConfig), countCommands(rawConfig))
}

func loadCommandConfig() (map[string]any, error) {
	return codegen.LoadYAMLConfig(CommandsYAMLPath)
}

func generateCoreConstants() error {
	rawConfig, err := loadCommandConfig()
	if err != nil {
		return pkgerrors.NewServiceError("generator", "load raw config", err)
	}

	tmpl, err := codegen.ParseTemplate(CoreTemplateFilePath, CoreTemplateName)
	if err != nil {
		return err
	}

	file, err := os.Create(CoreGeneratedPath)
	if err != nil {
		// Try creating the directory and retry
		if err := codegen.EnsureDir(filepath.Dir(CoreGeneratedPath)); err != nil {
			return pkgerrors.NewServiceError("generator", "create directory", err)
		}
		file, err = os.Create(CoreGeneratedPath)
		if err != nil {
			return pkgerrors.NewServiceError("generator", "create constants file", err)
		}
	}
	defer func() { _ = file.Close() }()

	data := struct {
		Commands          []constantData
		Flags             []constantData
		Messages          []constantData
		Icons             []constantData
		CommandsData      []commandData
		ValidationOptions []validationOption
	}{
		Commands:          collectCommandConstants(rawConfig),
		Flags:             collectFlagsData(rawConfig),
		Messages:          collectMessagesData(rawConfig),
		Icons:             collectIconsData(rawConfig),
		CommandsData:      collectCommandsData(rawConfig),
		ValidationOptions: collectValidationOptions(rawConfig),
	}

	return tmpl.Execute(file, data)
}

func generateCIUtilities() error {
	tmpl, err := template.ParseFiles(CLITemplateFilePath)
	if err != nil {
		return pkgerrors.NewServiceError("generator", "parse CI template", err)
	}

	file, err := os.Create(CIGeneratedPath)
	if err != nil {
		return pkgerrors.NewServiceError("generator", "create CI file", err)
	}
	defer func() { _ = file.Close() }()

	// CI utilities don't need data from YAML, they're static
	return tmpl.Execute(file, nil)
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
			Name:  PrefixCommand + toPascalCase(cmdName),
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
			Name:  PrefixFlag + toPascalCase(flagName),
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
	defaultFlagType = TypeString
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
	if commands, ok := rawConfig[KeyCommands].(map[string]any); ok {
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
	if flags, ok := global["flags"].(map[string]any); ok {
		// Add persistent flags
		if persistent, ok := flags["persistent"].(map[string]any); ok {
			for flagName := range persistent {
				flagNames[flagName] = true
			}
		}
		// Add conditional flags
		if conditional, ok := flags["conditional"].(map[string]any); ok {
			for flagName := range conditional {
				flagNames[flagName] = true
			}
		}
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

func collectValidationOptions(rawConfig map[string]any) []validationOption {
	var options []validationOption

	validation := getValidation(rawConfig)
	for key, data := range validation {
		if option := createValidationOption(key, data); option.Key != "" {
			options = append(options, option)
		}
	}

	return options
}

func createValidationOption(key string, data any) validationOption {
	dataMap := getCategoryMap(data)
	desc := getStringValue(dataMap["description"])
	if desc == "" {
		return validationOption{}
	}

	required := false
	if reqVal, ok := dataMap["required"].(bool); ok {
		required = reqVal
	}

	return validationOption{Key: key, Description: desc, Required: required}
}

func getValidation(rawConfig map[string]any) map[string]any {
	if validation, ok := rawConfig["validation"].(map[string]any); ok {
		return validation
	}
	return make(map[string]any)
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

func generateCLICommands(rawConfig map[string]any) error {
	tmpl, err := template.New(filepath.Base(CLICommandsTemplate)).Funcs(template.FuncMap{
		"toPascalCase":       toPascalCase,
		"getCommandCategory": func(cmdName string) string { return getCommandCategory(rawConfig, cmdName) },
	}).ParseFiles(CLICommandsTemplate)
	if err != nil {
		return pkgerrors.NewServiceError("generator", "parse CLI template", err)
	}

	file, err := os.Create(CLICommandsPath)
	if err != nil {
		if err := codegen.EnsureDir(filepath.Dir(CLICommandsPath)); err != nil {
			return pkgerrors.NewServiceError("generator", "create directory", err)
		}
		file, err = os.Create(CLICommandsPath)
		if err != nil {
			return pkgerrors.NewServiceError("generator", "create CLI file", err)
		}
	}
	defer func() { _ = file.Close() }()

	return tmpl.Execute(file, rawConfig)
}

func getCommandCategory(rawConfig map[string]any, cmdName string) string {
	categories, ok := rawConfig["categories"].(map[string]any)
	if !ok {
		return ""
	}

	for catName, catData := range categories {
		catMap, ok := catData.(map[string]any)
		if !ok {
			continue
		}

		commands, ok := catMap["commands"].([]any)
		if !ok {
			continue
		}

		for _, cmd := range commands {
			if cmdStr, ok := cmd.(string); ok && cmdStr == cmdName {
				return catName
			}
		}
	}
	return ""
}
