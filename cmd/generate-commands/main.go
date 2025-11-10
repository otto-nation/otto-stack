package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	pkgConfig "github.com/otto-nation/otto-stack/internal/pkg/config"
)

const (
	TemplateDir          = "cmd/generate-commands/templates/"
	ApatersTemplateFile  = "adapters.tmpl"
	CommandsTemplateFile = "commands.tmpl"
	GeneratedFilePath    = "internal/pkg/cli/commands_generated.go"
)

type templateData struct {
	Commands    []commandData
	GlobalFlags []flagData
}

type commandData struct {
	Name            string
	FuncName        string
	HandlerName     string
	ConstantName    string
	Handler         string
	Description     string
	LongDescription string
	Flags           []flagData
}

type flagData struct {
	Name        string
	Type        string
	Short       string
	Description string
	Default     string
}

func main() {
	commandConfig, err := pkgConfig.LoadCommandConfigStruct()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load commands config: %v\n", err)
		os.Exit(1)
	}

	if err := generateCommands(*commandConfig); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to generate commands: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated %d commands\n", len(commandConfig.Commands))
}

func generateCommands(commandConfig pkgConfig.CommandConfig) error {
	// Parse templates
	tmpl, err := template.ParseFiles(
		TemplateDir+CommandsTemplateFile,
		TemplateDir+ApatersTemplateFile,
	)
	if err != nil {
		return fmt.Errorf("failed to parse templates: %w", err)
	}

	// Create output file
	file, err := os.Create(GeneratedFilePath)
	if err != nil {
		// Try creating the directory and retry
		if err := os.MkdirAll(filepath.Dir(GeneratedFilePath), 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
		file, err = os.Create(GeneratedFilePath)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
	}
	defer func() { _ = file.Close() }()

	// Prepare template data
	var commands []commandData
	for name, cmd := range commandConfig.Commands {
		commands = append(commands, commandData{
			Name:            name,
			FuncName:        toPascalCase(name),
			HandlerName:     toPascalCase(name),
			ConstantName:    toPascalCase(name),
			Handler:         cmd.Handler,
			Description:     escape(cmd.Description),
			LongDescription: escape(cmd.LongDescription),
			Flags:           extractFlags(cmd.Flags),
		})
	}

	data := templateData{
		Commands:    commands,
		GlobalFlags: extractFlags(commandConfig.Global.Flags),
	}

	// Execute main template
	if err := tmpl.ExecuteTemplate(file, CommandsTemplateFile, data); err != nil {
		return fmt.Errorf("failed to execute commands template: %w", err)
	}

	// Execute adapters template
	if err := tmpl.ExecuteTemplate(file, ApatersTemplateFile, nil); err != nil {
		return fmt.Errorf("failed to execute adapters template: %w", err)
	}

	return nil
}

func toPascalCase(s string) string {
	parts := strings.Split(s, "-")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return strings.Join(parts, "")
}

func extractFlags(flags map[string]pkgConfig.FlagConfig) []flagData {
	var result []flagData
	for name, flag := range flags {
		result = append(result, flagData{
			Name:        name,
			Type:        flag.Type,
			Short:       flag.Short,
			Description: escape(flag.Description),
			Default:     fmt.Sprintf("%v", flag.Default),
		})
	}
	return result
}

func escape(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	return s
}
