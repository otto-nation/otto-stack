package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/otto-nation/otto-stack/internal/config"
	pkgConfig "github.com/otto-nation/otto-stack/internal/pkg/config"
	"gopkg.in/yaml.v3"
)

const (
	TemplateDir          = "cmd/generate-commands/templates/"
	ApatersTemplateFile  = "adapters.tmpl"
	CommandsTemplateFile = "commands.tmpl"
	GeneratedFilePath    = "internal/pkg/cli/commands_generated.go"
)

type templateData struct {
	Commands []commandData
}

type commandData struct {
	Name            string
	FuncName        string
	HandlerName     string
	ConstantName    string
	Handler         string
	Description     string
	LongDescription string
}

func main() {
	var commandConfig pkgConfig.CommandConfig
	if err := yaml.Unmarshal(config.EmbeddedCommandsYAML, &commandConfig); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse commands.yaml: %v\n", err)
		os.Exit(1)
	}

	if err := generateCommands(commandConfig); err != nil {
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
		return fmt.Errorf("failed to create file: %w", err)
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
		})
	}

	data := templateData{Commands: commands}

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

func escape(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	return s
}
