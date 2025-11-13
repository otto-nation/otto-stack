package project

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
)

// createDirectoryStructure creates the necessary directory structure
func (h *InitHandler) createDirectoryStructure() error {
	directories := []string{
		core.OttoStackDir,
		filepath.Join(core.OttoStackDir, core.ServiceConfigsDir),
	}

	for _, dir := range directories {
		if err := os.MkdirAll(dir, core.PermReadWriteExec); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// createConfigFile creates the main configuration file
func (h *InitHandler) createConfigFile(projectName string, services []string, validationOptions map[string]bool, base *base.BaseCommand) error {
	configContent := h.generateConfig(projectName, services, validationOptions)

	configPath := core.OttoStackDir + "/" + core.ConfigFileName
	if err := os.WriteFile(configPath, []byte(configContent), core.PermReadWrite); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	base.Output.Success(core.MsgSuccess_created_file, configPath)
	return nil
}

// createGitignoreEntries adds entries to .gitignore
func (h *InitHandler) createGitignoreEntries(base *base.BaseCommand) error {
	gitignorePath := core.GitIgnoreFileName

	// Check if .gitignore exists
	var existingContent []byte
	if content, err := os.ReadFile(gitignorePath); err == nil {
		existingContent = content
	}

	// Check if entries already exist
	existingStr := string(existingContent)
	hasDevStackEntries := false
	for _, entry := range core.GitignoreEntries {
		if entry != "" && contains(existingStr, entry) {
			hasDevStackEntries = true
			break
		}
	}

	if hasDevStackEntries {
		base.Output.Info("%s", core.MsgFiles_gitignore_exists)
		return nil
	}

	// Append entries
	file, err := os.OpenFile(gitignorePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, core.PermReadWrite)
	if err != nil {
		return fmt.Errorf("failed to open .gitignore: %w", err)
	}
	defer func() { _ = file.Close() }()

	for _, entry := range core.GitignoreEntries {
		if _, err := file.WriteString(entry + "\n"); err != nil {
			return fmt.Errorf("failed to write to .gitignore: %w", err)
		}
	}

	base.Output.Success("%s", core.MsgSuccess_updated_gitignore)
	return nil
}

// createReadme creates a README file for the project
func (h *InitHandler) createReadme(projectName string, services []string, base *base.BaseCommand) error {
	readmeContent := fmt.Sprintf(`# %s Otto Stack

This directory contains the development stack configuration for %s.

## Services

The following services are configured:

%s

## Quick Start

1. Start the stack:
   `+"```bash"+`
   %s up
   `+"```"+`

2. Check status:
   `+"```bash"+`
   %s status
   `+"```"+`

3. Stop the stack:
   `+"```bash"+`
   %s down
   `+"```"+`

## Configuration

- Main config: `+"`%s`"+`
- Docker Compose: `+"`%s`"+`
- Environment: `+"`.env.generated`"+`

## Commands

Run `+"`%s --help`"+` for a full list of available commands.
`, projectName, projectName, formatServicesList(services),
		core.AppName, core.AppName, core.AppName,
		core.ConfigFileName, docker.DockerComposeFileName, core.AppName)

	readmePath := core.OttoStackDir + "/" + core.ReadmeFileName
	if err := os.WriteFile(readmePath, []byte(readmeContent), core.PermReadWrite); err != nil {
		return fmt.Errorf("failed to create README: %w", err)
	}

	base.Output.Success(core.MsgSuccess_created_file, readmePath)
	return nil
}

// formatServicesList formats the services list for README
func formatServicesList(services []string) string {
	if len(services) == 0 {
		return "- No services configured"
	}

	var builder strings.Builder
	for _, service := range services {
		builder.WriteString(fmt.Sprintf("- %s\n", service))
	}
	return builder.String()
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) &&
			(s[:len(substr)] == substr ||
				s[len(s)-len(substr):] == substr ||
				containsSubstring(s, substr))))
}

// containsSubstring checks if string contains substring anywhere
func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
