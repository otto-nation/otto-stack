package context

import (
	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
)

// Context encapsulates all data needed for CLI operations
type Context struct {
	Project  ProjectSpec
	Services ServiceSpec
	Options  OptionSet
	Runtime  RuntimeSpec
	Sharing  *SharingSpec
	Advanced *AdvancedSpec
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

// SharingSpec contains container sharing information.
//
// Services semantics: nil/empty = share all catalog-shareable services;
// non-empty = whitelist where true entries are shared and false/absent entries run locally.
type SharingSpec struct {
	Enabled bool
	// Services is the optional whitelist built during init. true = shared globally,
	// false/absent = project-local. Empty means share all shareable catalog services.
	Services map[string]bool
}

// AdvancedSpec contains advanced init-time options
type AdvancedSpec struct {
	AutoStart bool
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

// WithRuntimeFlags sets runtime flags from InitFlags
func (b *Builder) WithRuntimeFlags(flags *core.InitFlags, interactive bool) *Builder {
	b.ctx.Runtime = RuntimeSpec{
		Force:       flags.Force,
		Interactive: interactive,
		DryRun:      false,
	}
	return b
}

// WithSharing sets sharing configuration
func (b *Builder) WithSharing(sharing *SharingSpec) *Builder {
	b.ctx.Sharing = sharing
	return b
}

// WithAdvancedSpec sets advanced options
func (b *Builder) WithAdvancedSpec(advanced *AdvancedSpec) *Builder {
	b.ctx.Advanced = advanced
	return b
}

// Build returns the constructed Context
func (b *Builder) Build() Context {
	return b.ctx
}
