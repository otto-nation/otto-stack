package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	pkgConfig "github.com/otto-nation/otto-stack/internal/pkg/config"
)

type templateData struct {
	Commands []commandData
	Flags    []constantData
	Messages []constantData
	Icons    []constantData
}

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

func generateWithTemplates(config *pkgConfig.CommandConfig) error {
	if err := generateConstantsTemplate(config); err != nil {
		return fmt.Errorf("failed to generate constants: %w", err)
	}

	if err := generateFlagsTemplate(config); err != nil {
		return fmt.Errorf("failed to generate flags: %w", err)
	}

	return nil
}

func generateConstantsTemplate(config *pkgConfig.CommandConfig) error {
	tmpl, err := template.ParseFiles("cmd/generate-cli/templates/constants.tmpl")
	if err != nil {
		return fmt.Errorf("failed to parse constants template: %w", err)
	}

	file, err := os.Create("internal/pkg/constants/cli_generated_template.go")
	if err != nil {
		return fmt.Errorf("failed to create constants file: %w", err)
	}
	defer file.Close()

	data := templateData{
		Flags:    collectFlagsData(config),
		Messages: collectMessagesData(config),
		Icons:    collectIconsData(config),
	}

	return tmpl.Execute(file, data)
}

func generateFlagsTemplate(config *pkgConfig.CommandConfig) error {
	tmpl, err := template.ParseFiles("cmd/generate-cli/templates/flags.tmpl")
	if err != nil {
		return fmt.Errorf("failed to parse flags template: %w", err)
	}

	file, err := os.Create("internal/pkg/cli/flags_generated_template.go")
	if err != nil {
		return fmt.Errorf("failed to create flags file: %w", err)
	}
	defer file.Close()

	var commands []commandData
	for cmdName, cmd := range config.Commands {
		commands = append(commands, commandData{
			CommandName: cmdName,
			StructName:  toPascalCase(cmdName) + "Flags",
			FuncName:    "Parse" + toPascalCase(cmdName) + "Flags",
			Fields:      collectCommandFlags(cmd, config),
		})
	}

	return tmpl.Execute(file, templateData{Commands: commands})
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
	var messages []constantData
	if config.Messages == nil {
		return messages
	}

	for category, categoryData := range config.Messages {
		categoryMap := categoryData.(map[string]any)
		for key, value := range categoryMap {
			messages = append(messages, constantData{
				Name:  "Message" + toPascalCase(strings.ReplaceAll(category+"."+key, ".", "_")),
				Value: value.(string),
			})
		}
	}
	return messages
}

func collectIconsData(config *pkgConfig.CommandConfig) []constantData {
	var icons []constantData
	if config.Icons == nil {
		return icons
	}

	for category, categoryData := range config.Icons {
		categoryMap := categoryData.(map[string]any)
		for key, value := range categoryMap {
			icons = append(icons, constantData{
				Name:  "Icon" + toPascalCase(strings.ReplaceAll(category+"."+key, ".", "_")),
				Value: value.(string),
			})
		}
	}
	return icons
}
