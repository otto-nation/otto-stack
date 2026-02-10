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
	// Shared context - managing shared containers
	Shared ContextType = "shared"
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
	Type             ContextType
	Project          *ProjectInfo // nil if Shared context
	SharedContainers *SharedInfo  // Always populated
}

// SharingConfig defines container sharing configuration
type SharingConfig struct {
	Enabled  bool            `yaml:"enabled" json:"enabled"`
	Services map[string]bool `yaml:"services,omitempty" json:"services,omitempty"`
}

// ExecutionMode interface for type-safe context discrimination
type ExecutionMode interface {
	SharedRoot() string
	isExecutionMode() // unexported marker prevents external implementation
}

// ProjectMode for project-scoped operations
type ProjectMode struct {
	Project *ProjectInfo
	Shared  *SharedInfo
}

func (p *ProjectMode) SharedRoot() string { return p.Shared.Root }
func (p *ProjectMode) isExecutionMode()   {}

// SharedMode for global shared container operations
type SharedMode struct {
	Shared *SharedInfo
}

func (s *SharedMode) SharedRoot() string { return s.Shared.Root }
func (s *SharedMode) isExecutionMode()   {}

// NewProjectInfo creates a ProjectInfo from a config directory path
func NewProjectInfo(configDir string) *ProjectInfo {
	return &ProjectInfo{
		Root:       filepath.Dir(configDir),
		ConfigDir:  configDir,
		ConfigFile: filepath.Join(configDir, core.ConfigFileName),
	}
}
