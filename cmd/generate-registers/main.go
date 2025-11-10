package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"text/template"

	pkgConfig "github.com/otto-nation/otto-stack/internal/pkg/config"
)

const (
	TemplateFilePath  = "cmd/generate-registers/templates/register.tmpl"
	GeneratedFilePath = "internal/pkg/cli/handlers/%s/register_generated.go"
)

type templateData struct {
	Package  string
	Commands []commandData
}

type commandData struct {
	Constant    string
	HandlerFunc string
}

func main() {
	// Load commands config
	commandConfig, err := pkgConfig.LoadCommandConfigStruct()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load commands config: %v\n", err)
		os.Exit(1)
	}

	// Group commands by handler
	handlerCommands := make(map[string][]string)
	for cmdName, cmd := range commandConfig.Commands {
		if cmd.Handler != "" {
			handlerCommands[cmd.Handler] = append(handlerCommands[cmd.Handler], cmdName)
		}
	}

	// Generate register.go for each handler
	for handler, commands := range handlerCommands {
		sort.Strings(commands)
		if err := generateRegisterFile(handler, commands); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to generate %s: %v\n", handler, err)
			continue
		}
		fmt.Printf("Generated %s\n", fmt.Sprintf(GeneratedFilePath, handler))
	}

	fmt.Printf("Generated register.go files for %d handlers\n", len(handlerCommands))
}

func generateRegisterFile(handler string, commands []string) error {
	tmpl, err := template.ParseFiles(TemplateFilePath)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	outputPath := fmt.Sprintf(GeneratedFilePath, handler)
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer func() { _ = file.Close() }()

	var commandsData []commandData
	for _, cmd := range commands {
		commandsData = append(commandsData, commandData{
			Constant:    "Command" + toPascalCase(cmd),
			HandlerFunc: "New" + toPascalCase(cmd) + "Handler",
		})
	}

	data := templateData{
		Package:  handler,
		Commands: commandsData,
	}

	return tmpl.Execute(file, data)
}

// toPascalCase converts kebab-case to PascalCase
func toPascalCase(s string) string {
	parts := strings.Split(s, "-")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
		}
	}
	return strings.Join(parts, "")
}
