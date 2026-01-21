package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/otto-nation/otto-stack/cmd/codegen"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
)

const (
	CommandsYAMLPath               = "internal/config/commands.yaml"
	MessagesYAMLPath               = "internal/config/messages.yaml"
	CoreTemplateFilePath           = "cmd/generate-cli/templates/core.tmpl"
	CLICommandsTemplate            = "cmd/generate-cli/templates/cli.tmpl"
	RegisterTemplateFilePath       = "cmd/generate-cli/templates/register.tmpl"
	ValidationTestTemplateFilePath = "cmd/generate-cli/templates/validation_test.tmpl"
	CoreGeneratedPath              = "internal/core/constants_generated.go"
	CLICommandsPath                = "internal/pkg/cli/cli_generated.go"
	ValidationTestGeneratedPath    = "internal/pkg/cli/flags_validation_generated_test.go"
	HandlersBasePath               = "internal/pkg/cli/handlers"
	CoreTemplateName               = "core"
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
	KeyCategories  = "categories"
	KeyPersistent  = "persistent"
	KeyConditional = "conditional"
	KeyRequired    = "required"
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

const (
	// Built-in commands
	CommandHelp = "help"
)

const (
	// Handler registration paths
	HandlerBasePath     = "internal/pkg/cli/handlers"
	RegisterGenFileName = "register_generated.go"
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

	if err := generateCLICommands(rawConfig); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to generate CLI commands: %v\n", err)
		os.Exit(1)
	}

	if err := generateHandlerRegistrations(rawConfig); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to generate handler registrations: %v\n", err)
		os.Exit(1)
	}

	if err := generateFlagValidationTest(rawConfig); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to generate flag validation test: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated CLI code for %d flags and %d commands\n",
		countFlags(rawConfig), countCommands(rawConfig))
}

func loadCommandConfig() (map[string]any, error) {
	return codegen.LoadYAMLConfig(CommandsYAMLPath)
}

func loadMessagesConfig() (map[string]any, error) {
	return codegen.LoadYAMLConfig(MessagesYAMLPath)
}

func generateCoreConstants() error {
	rawConfig, err := loadCommandConfig()
	if err != nil {
		return pkgerrors.NewServiceError("generator", "load commands config", err)
	}

	messagesConfig, err := loadMessagesConfig()
	if err != nil {
		return pkgerrors.NewServiceError("generator", "load messages config", err)
	}

	// Merge messages into rawConfig for template
	rawConfig[KeyMessages] = messagesConfig

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

	flagType := getStringValue(flag[KeyType])
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
	if messages, ok := rawConfig[KeyMessages].(map[string]any); ok {
		return messages
	}
	return make(map[string]any)
}

func getIcons(rawConfig map[string]any) map[string]any {
	if icons, ok := rawConfig[KeyIcons].(map[string]any); ok {
		return icons
	}
	return make(map[string]any)
}

func getGlobal(rawConfig map[string]any) map[string]any {
	if global, ok := rawConfig[KeyGlobal].(map[string]any); ok {
		return global
	}
	return make(map[string]any)
}

func getFlags(cmd map[string]any) map[string]any {
	if flags, ok := cmd[KeyFlags].(map[string]any); ok {
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
	if flags, ok := global[KeyFlags].(map[string]any); ok {
		// Add persistent flags
		if persistent, ok := flags[KeyPersistent].(map[string]any); ok {
			for flagName := range persistent {
				flagNames[flagName] = true
			}
		}
		// Add conditional flags
		if conditional, ok := flags[KeyConditional].(map[string]any); ok {
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
	desc := getStringValue(dataMap[KeyDescription])
	if desc == "" {
		return validationOption{}
	}

	required := false
	if reqVal, ok := dataMap[KeyRequired].(bool); ok {
		required = reqVal
	}

	return validationOption{Key: key, Description: desc, Required: required}
}

func getValidation(rawConfig map[string]any) map[string]any {
	if validation, ok := rawConfig[KeyValidation].(map[string]any); ok {
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
	categories := getCategories(rawConfig)

	for catName, catData := range categories {
		catMap := getCategoryMap(catData)
		commands := getCategoryCommandList(catMap)

		for _, cmd := range commands {
			if cmdStr := getStringValue(cmd); cmdStr == cmdName {
				return catName
			}
		}
	}
	return ""
}

type handlerRegistrationData struct {
	Package  string
	Commands []handlerCommand
}

type handlerCommand struct {
	ConstName   string
	HandlerName string
}

func generateHandlerRegistrations(rawConfig map[string]any) error {
	categories := getCategories(rawConfig)
	if len(categories) == 0 {
		return pkgerrors.NewServiceError("generator", "load categories", fmt.Errorf("no categories found in config"))
	}

	tmpl, err := template.ParseFiles(RegisterTemplateFilePath)
	if err != nil {
		return pkgerrors.NewServiceError("generator", "parse register template", err)
	}

	for categoryName, catData := range categories {
		if err := generateCategoryRegistration(tmpl, categoryName, catData); err != nil {
			return err
		}
	}

	return nil
}

func generateCategoryRegistration(tmpl *template.Template, categoryName string, catData any) error {
	commands := extractCategoryCommands(catData)
	if len(commands) == 0 {
		return nil
	}

	data := handlerRegistrationData{
		Package:  categoryName,
		Commands: commands,
	}

	outputPath := filepath.Join(HandlerBasePath, categoryName, RegisterGenFileName)
	file, err := os.Create(outputPath)
	if err != nil {
		if err := codegen.EnsureDir(filepath.Dir(outputPath)); err != nil {
			return pkgerrors.NewServiceError("generator", "create directory", err)
		}
		file, err = os.Create(outputPath)
		if err != nil {
			return pkgerrors.NewServiceError("generator", "create register file", err)
		}
	}
	defer func() { _ = file.Close() }()

	return tmpl.Execute(file, data)
}

func extractCategoryCommands(catData any) []handlerCommand {
	var handlerCommands []handlerCommand
	catMap := getCategoryMap(catData)
	commands := getCategoryCommandList(catMap)

	for _, cmd := range commands {
		if cmdStr := getStringValue(cmd); cmdStr != "" && !isBuiltInCommand(cmdStr) {
			handlerCommands = append(handlerCommands, handlerCommand{
				ConstName:   toPascalCase(cmdStr),
				HandlerName: toPascalCase(cmdStr),
			})
		}
	}

	return handlerCommands
}

func getCategories(rawConfig map[string]any) map[string]any {
	if categories, ok := rawConfig[KeyCategories].(map[string]any); ok {
		return categories
	}
	return make(map[string]any)
}

func getCategoryCommandList(catMap map[string]any) []any {
	if commands, ok := catMap[KeyCommands].([]any); ok {
		return commands
	}
	return []any{}
}

func isBuiltInCommand(cmd string) bool {
	return cmd == CommandHelp
}

func generateFlagValidationTest(rawConfig map[string]any) error {
	globalFlags := extractGlobalFlags(rawConfig)
	commands := buildCommandFlagInfo(rawConfig, globalFlags)

	data := struct {
		Commands []CommandInfo
	}{
		Commands: commands,
	}

	tmpl, err := template.ParseFiles(ValidationTestTemplateFilePath)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	outputFile, err := os.Create(ValidationTestGeneratedPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	if err := tmpl.Execute(outputFile, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

type CommandInfo struct {
	Name        string
	HandlerPath string
	Flags       []string
}

func extractGlobalFlags(rawConfig map[string]any) map[string]bool {
	globalFlags := make(map[string]bool)

	globalData, ok := rawConfig[KeyGlobal].(map[string]any)
	if !ok {
		return globalFlags
	}

	flagsData, ok := globalData[KeyFlags].(map[string]any)
	if !ok {
		return globalFlags
	}

	// Get persistent flags
	if persistent, ok := flagsData[KeyPersistent].(map[string]any); ok {
		for flagName := range persistent {
			globalFlags[flagName] = true
		}
	}

	// Get conditional flags
	if conditional, ok := flagsData[KeyConditional].(map[string]any); ok {
		for flagName := range conditional {
			globalFlags[flagName] = true
		}
	}

	return globalFlags
}

func buildCommandFlagInfo(rawConfig map[string]any, globalFlags map[string]bool) []CommandInfo {
	var commands []CommandInfo

	categoriesData, ok := rawConfig[KeyCategories].(map[string]any)
	if !ok {
		return commands
	}

	commandsData, ok := rawConfig[KeyCommands].(map[string]any)
	if !ok {
		return commands
	}

	// Build map of command -> category
	cmdToCategory := buildCommandCategoryMap(categoriesData)

	// Process each command (sorted for deterministic output)
	var cmdNames []string
	for cmdName := range commandsData {
		cmdNames = append(cmdNames, cmdName)
	}
	sort.Strings(cmdNames)

	for _, cmdName := range cmdNames {
		cmdData := commandsData[cmdName]
		cmdMap, ok := cmdData.(map[string]any)
		if !ok {
			continue
		}

		category, hasCategory := cmdToCategory[cmdName]
		if !hasCategory {
			continue
		}

		handlerPath := fmt.Sprintf("%s/%s/%s.go", HandlersBasePath, category, cmdName)

		flagsData, ok := cmdMap[KeyFlags].(map[string]any)
		if !ok || len(flagsData) == 0 {
			continue
		}

		flags := extractNonGlobalFlags(flagsData, globalFlags)
		if len(flags) > 0 {
			commands = append(commands, CommandInfo{
				Name:        cmdName,
				HandlerPath: handlerPath,
				Flags:       flags,
			})
		}
	}

	return commands
}

func buildCommandCategoryMap(categoriesData map[string]any) map[string]string {
	cmdToCategory := make(map[string]string)

	for catName, catData := range categoriesData {
		catMap, ok := catData.(map[string]any)
		if !ok {
			continue
		}

		cmdList := getCategoryCommandList(catMap)
		for _, cmd := range cmdList {
			if cmdStr, ok := cmd.(string); ok {
				cmdToCategory[cmdStr] = catName
			}
		}
	}

	return cmdToCategory
}

func extractNonGlobalFlags(flagsData map[string]any, globalFlags map[string]bool) []string {
	var flags []string

	for flagName := range flagsData {
		if !globalFlags[flagName] {
			flags = append(flags, flagName)
		}
	}

	sort.Strings(flags) // Sort for deterministic output
	return flags
}
