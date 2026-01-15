package context

import (
	"github.com/otto-nation/otto-stack/internal/pkg/types"
)

// Context encapsulates all data needed for CLI operations
type Context struct {
	Project  ProjectSpec
	Services ServiceSpec
	Options  OptionSet
	Runtime  RuntimeSpec
}

// ProjectSpec contains project-related information
type ProjectSpec struct {
	Name string
	Path string
}

// ServiceSpec contains service-related information
type ServiceSpec struct {
	Names   []string
	Configs []types.ServiceConfig
}

// OptionSet contains configuration options
type OptionSet struct {
	Validation map[string]bool
	Advanced   map[string]bool
}

// RuntimeSpec contains runtime flags and settings
type RuntimeSpec struct {
	Force       bool
	Interactive bool
	DryRun      bool
}

// Builder provides fluent API for building Context
type Builder struct {
	ctx Context
}

// NewBuilder creates a new context builder
func NewBuilder() *Builder {
	return &Builder{
		ctx: Context{
			Options: OptionSet{
				Validation: make(map[string]bool),
				Advanced:   make(map[string]bool),
			},
		},
	}
}

// WithProject sets project information
func (b *Builder) WithProject(name, path string) *Builder {
	b.ctx.Project = ProjectSpec{
		Name: name,
		Path: path,
	}
	return b
}

// WithServices sets service information
func (b *Builder) WithServices(names []string, configs []types.ServiceConfig) *Builder {
	b.ctx.Services = ServiceSpec{
		Names:   names,
		Configs: configs,
	}
	return b
}

// WithValidation sets validation options
func (b *Builder) WithValidation(validation map[string]bool) *Builder {
	b.ctx.Options.Validation = validation
	return b
}

// WithAdvanced sets advanced options
func (b *Builder) WithAdvanced(advanced map[string]bool) *Builder {
	b.ctx.Options.Advanced = advanced
	return b
}

// WithRuntime sets runtime flags
func (b *Builder) WithRuntime(force, interactive, dryRun bool) *Builder {
	b.ctx.Runtime = RuntimeSpec{
		Force:       force,
		Interactive: interactive,
		DryRun:      dryRun,
	}
	return b
}

// Build returns the constructed Context
func (b *Builder) Build() Context {
	return b.ctx
}
