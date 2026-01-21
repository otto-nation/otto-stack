package config

import "time"

// Config defines the single source of truth for otto-stack configuration
type Config struct {
	Project    ProjectConfig     `yaml:"project" json:"project"`
	Stack      StackConfig       `yaml:"stack" json:"stack"`
	Sharing    *SharingConfig    `yaml:"sharing,omitempty" json:"sharing,omitempty"`
	Validation *ValidationConfig `yaml:"validation,omitempty" json:"validation,omitempty"`
}

// ProjectConfig defines project-level configuration
type ProjectConfig struct {
	Name      string    `yaml:"name" json:"name" validate:"required"`
	Type      string    `yaml:"type" json:"type" validate:"required"`
	CreatedAt time.Time `yaml:"created_at,omitempty" json:"created_at"`
	UpdatedAt time.Time `yaml:"updated_at,omitempty" json:"updated_at"`
}

// StackConfig defines stack-level configuration
type StackConfig struct {
	Enabled []string `yaml:"enabled" json:"enabled"`
}

// SharingConfig defines container sharing configuration
type SharingConfig struct {
	Enabled  bool            `yaml:"enabled" json:"enabled"`
	Services map[string]bool `yaml:"services,omitempty" json:"services,omitempty"`
}

// ValidationConfig defines validation settings
type ValidationConfig struct {
	Options map[string]bool `yaml:"options,omitempty" json:"options,omitempty"`
}
