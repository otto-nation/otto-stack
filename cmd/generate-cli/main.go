package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

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

	// Extract service constants
	serviceConstants, err := extractServiceConstants()
	if err != nil {
		return fmt.Errorf("failed to extract service constants: %w", err)
	}

	data := struct {
		Commands       []constantData
		Flags          []constantData
		Messages       []constantData
		Icons          []constantData
		CommandsData   []commandData
		ServiceClients []constantData
		ServicePorts   []constantData
		ServiceNames   []constantData
	}{
		Commands:       collectCommandConstants(rawConfig),
		Flags:          collectFlagsData(rawConfig),
		Messages:       collectMessagesData(rawConfig),
		Icons:          collectIconsData(rawConfig),
		CommandsData:   collectCommandsData(rawConfig),
		ServiceClients: serviceConstants.Clients,
		ServicePorts:   serviceConstants.Ports,
		ServiceNames:   serviceConstants.Names,
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

// ServiceConstants holds extracted service constants
type ServiceConstants struct {
	Clients []constantData
	Ports   []constantData
	Names   []constantData
}

// ServiceYAML represents service YAML structure for constant extraction
type ServiceYAML struct {
	Name       string `yaml:"name"`
	Connection struct {
		Client      string `yaml:"client"`
		DefaultPort int    `yaml:"default_port"`
	} `yaml:"connection"`
}

// extractServiceConstants extracts constants from service YAML files
func extractServiceConstants() (*ServiceConstants, error) {
	const servicesDir = "internal/config/services"

	clients := make(map[string]bool)
	ports := make(map[string]int)
	names := make(map[string]bool)

	err := filepath.Walk(servicesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(path, ".yaml") {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		var service ServiceYAML
		if err := yaml.Unmarshal(data, &service); err != nil {
			return err
		}

		// Extract service name
		if service.Name != "" {
			names[service.Name] = true
		}

		// Extract client
		if service.Connection.Client != "" {
			clients[service.Connection.Client] = true
		}

		// Extract port
		if service.Connection.DefaultPort > 0 {
			portKey := fmt.Sprintf("%s_port", strings.ToUpper(service.Name))
			ports[portKey] = service.Connection.DefaultPort
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	result := &ServiceConstants{}

	// Convert to constantData
	for client := range clients {
		constName := "Client" + toPascalCase(strings.ReplaceAll(client, "-", "_"))
		result.Clients = append(result.Clients, constantData{
			Name:  constName,
			Value: client,
		})
	}

	for portKey, port := range ports {
		constName := "DefaultPort" + toPascalCase(portKey)
		result.Ports = append(result.Ports, constantData{
			Name:  constName,
			Value: fmt.Sprintf("%d", port),
		})
	}

	for name := range names {
		constName := "Service" + toPascalCase(name)
		result.Names = append(result.Names, constantData{
			Name:  constName,
			Value: name,
		})
	}

	return result, nil
}
