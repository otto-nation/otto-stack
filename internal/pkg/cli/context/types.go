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
