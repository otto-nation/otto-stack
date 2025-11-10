package config

import "time"

// Config defines the single source of truth for otto-stack configuration
type Config struct {
	Project              ProjectConfig  `yaml:"project" json:"project"`
	Stack                StackConfig    `yaml:"stack" json:"stack"`
	ServiceConfiguration map[string]any `yaml:"service_configuration,omitempty" json:"service_configuration,omitempty"`
}

// ProjectConfig defines project-level configuration
type ProjectConfig struct {
	Name      string    `yaml:"name" json:"name" validate:"required"`
	Type      string    `yaml:"type" json:"type" validate:"required"`
	Services  []string  `yaml:"services,omitempty" json:"services,omitempty"`
	CreatedAt time.Time `yaml:"created_at,omitempty" json:"created_at"`
	UpdatedAt time.Time `yaml:"updated_at,omitempty" json:"updated_at"`
}

// StackConfig defines stack-level configuration
type StackConfig struct {
	Enabled []string `yaml:"enabled" json:"enabled"`
}
