package context

import (
	"path/filepath"

	"github.com/otto-nation/otto-stack/internal/core"
)

// ContextType represents the execution context
type ContextType string

const (
	// Project context - running inside a project directory
	Project ContextType = "project"
	// Global context - running outside any project
	Global ContextType = "global"
)

// ProjectInfo contains project-specific information
type ProjectInfo struct {
	Root       string // Absolute path to project root (.otto-stack parent)
	ConfigDir  string // Absolute path to .otto-stack directory
	ConfigFile string // Absolute path to config.yaml
}

// SharedInfo contains shared container information
type SharedInfo struct {
	Root string // Absolute path to ~/.otto-stack/shared
}

// ExecutionContext represents the current execution context
type ExecutionContext struct {
	Type    ContextType
	Project *ProjectInfo // nil if Global context
	Shared  *SharedInfo  // Always populated
}

// IsProject returns true if running in project context
func (c *ExecutionContext) IsProject() bool {
	return c.Type == Project
}

// IsGlobal returns true if running in global context
func (c *ExecutionContext) IsGlobal() bool {
	return c.Type == Global
}

// GetProjectRoot returns the project root or empty string if global
func (c *ExecutionContext) GetProjectRoot() string {
	if c.Project != nil {
		return c.Project.Root
	}
	return ""
}

// GetConfigFile returns the config file path or empty string if global
func (c *ExecutionContext) GetConfigFile() string {
	if c.Project != nil {
		return c.Project.ConfigFile
	}
	return ""
}

// GetSharedRoot returns the shared containers root directory
func (c *ExecutionContext) GetSharedRoot() string {
	return c.Shared.Root
}

// SharingConfig defines container sharing configuration
type SharingConfig struct {
	Enabled  bool            `yaml:"enabled" json:"enabled"`
	Services map[string]bool `yaml:"services,omitempty" json:"services,omitempty"`
}

// NewProjectInfo creates a ProjectInfo from a config directory path
func NewProjectInfo(configDir string) *ProjectInfo {
	return &ProjectInfo{
		Root:       filepath.Dir(configDir),
		ConfigDir:  configDir,
		ConfigFile: filepath.Join(configDir, core.ConfigFileName),
	}
}
