package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/otto-nation/otto-stack/cmd/codegen"
	pkgConfig "github.com/otto-nation/otto-stack/internal/pkg/config"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
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

	// Define command categories for target folder structure
	commandCategories := map[string]string{
		// Lifecycle commands (from stack)
		"up":      "lifecycle",
		"down":    "lifecycle",
		"restart": "lifecycle",
		"cleanup": "lifecycle",

		// Operations commands (from stack)
		"status":  "operations",
		"logs":    "operations",
		"exec":    "operations",
		"connect": "operations",

		// Utility commands
		"web-interfaces": "utility",

		// Project commands
		"init":      "project",
		"services":  "project",
		"deps":      "project",
		"conflicts": "project",
		"validate":  "project",
		"doctor":    "project",
		"version":   "project",
	}

	// Group commands by category instead of handler
	categoryCommands := make(map[string][]string)
	for cmdName, cmd := range commandConfig.Commands {
		if cmd.Handler != "" {
			category := commandCategories[cmdName]
			if category == "" {
				// Fallback to handler if no category defined
				category = cmd.Handler
			}
			categoryCommands[category] = append(categoryCommands[category], cmdName)
		}
	}

	// Generate register.go for each category
	for category, commands := range categoryCommands {
		sort.Strings(commands)
		if err := generateRegisterFile(category, commands); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to generate %s: %v\n", category, err)
			continue
		}
		fmt.Printf("Generated %s\n", fmt.Sprintf(GeneratedFilePath, category))
	}

	fmt.Printf("Generated register.go files for %d categories\n", len(categoryCommands))
}

func generateRegisterFile(handler string, commands []string) error {
	tmpl, err := codegen.ParseTemplate(TemplateFilePath, "register")
	if err != nil {
		return err
	}

	outputPath := fmt.Sprintf(GeneratedFilePath, handler)
	file, err := os.Create(outputPath)
	if err != nil {
		return pkgerrors.NewServiceError("generator", "create file", err)
	}
	defer func() { _ = file.Close() }()

	var commandsData []commandData
	for _, cmd := range commands {
		commandsData = append(commandsData, commandData{
			Constant:    "Command" + codegen.ToPascalCase(cmd),
			HandlerFunc: "New" + codegen.ToPascalCase(cmd) + "Handler",
		})
	}

	data := templateData{
		Package:  handler,
		Commands: commandsData,
	}

	return tmpl.Execute(file, data)
}
