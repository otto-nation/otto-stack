package config

import "time"

// Config defines the single source of truth for otto-stack configuration
type Config struct {
	Project    ProjectConfig     `yaml:"project" json:"project"`
	Stack      StackConfig       `yaml:"stack" json:"stack"`
	Sharing    *SharingConfig    `yaml:"sharing,omitempty" json:"sharing,omitempty"`
	Validation *ValidationConfig `yaml:"validation,omitempty" json:"validation,omitempty"`
	Advanced   *AdvancedConfig   `yaml:"advanced,omitempty" json:"advanced,omitempty"`
	Version    *VersionConfig    `yaml:"version_config,omitempty" json:"version_config,omitempty"`
}

// ProjectConfig defines project-level configuration
type ProjectConfig struct {
	Name      string    `yaml:"name" json:"name" validate:"required"`
	Type      string    `yaml:"type" json:"type" validate:"required"`
	CreatedAt time.Time `yaml:"created_at,omitempty" json:"created_at"`
	UpdatedAt time.Time `yaml:"updated_at,omitempty" json:"updated_at"`
}

// StackConfig defines stack-level configuration
// StackConfig defines stack service configuration
type StackConfig struct {
	// Enabled contains user-selected services only
	// Dependencies are resolved automatically at runtime
	Enabled []string `yaml:"enabled" json:"enabled"`
}

// SharingConfig defines container sharing configuration.
//
// Semantics for Services:
//   - nil or empty map: all services marked shareable in the catalog are shared globally.
//   - non-empty map: only entries with value true are shared; all other shareable services
//     run as project-local containers.
type SharingConfig struct {
	Enabled bool `yaml:"enabled" json:"enabled"`
	// Services is an optional whitelist of service names to share globally.
	// true = run as a shared container; false/absent = run project-local.
	// An empty map shares every service the catalog marks as shareable.
	Services map[string]bool `yaml:"services,omitempty" json:"services,omitempty"`
}

// ValidationConfig defines validation settings
type ValidationConfig struct {
	Options map[string]bool `yaml:"options,omitempty" json:"options,omitempty"`
}

// AdvancedConfig defines advanced operational settings
type AdvancedConfig struct {
	AutoStart         bool `yaml:"auto_start" json:"auto_start"`
	PullLatestImages  bool `yaml:"pull_latest_images" json:"pull_latest_images"`
	CleanupOnRecreate bool `yaml:"cleanup_on_recreate" json:"cleanup_on_recreate"`
}

// VersionConfig defines version constraint settings
type VersionConfig struct {
	RequiredVersion string `yaml:"required_version,omitempty" json:"required_version,omitempty"`
}
