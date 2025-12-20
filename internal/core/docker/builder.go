package docker

import (
	"github.com/otto-nation/otto-stack/internal/core"
)

// NewCommand creates a new Docker command builder
func NewCommand(command string) *core.CommandBuilder {
	return core.NewCommandBuilder(command)
}

// ComposeBuilder provides a fluent interface for building docker-compose commands
type ComposeBuilder struct {
	builder  *core.CommandBuilder
	project  string
	services []string
}

// NewComposeBuilder creates a new compose builder
func NewComposeBuilder() *ComposeBuilder {
	return &ComposeBuilder{
		builder:  core.NewCommandBuilder(DockerCmd).Subcommand(DockerComposeCmd),
		services: make([]string, 0),
	}
}

// Project sets the compose project name
func (cb *ComposeBuilder) Project(name string) *ComposeBuilder {
	cb.project = name
	return cb
}

// Services sets the services to operate on
func (cb *ComposeBuilder) Services(names ...string) *ComposeBuilder {
	cb.services = names
	return cb
}

// File sets the compose file path
func (cb *ComposeBuilder) File(path string) *ComposeBuilder {
	cb.builder.Flag(FlagFile, path)
	return cb
}

// User sets the user for exec command
func (cb *ComposeBuilder) User(user string) *ComposeBuilder {
	if user != "" {
		cb.builder.Flag(FlagUser, user)
	}
	return cb
}

// Workdir sets the working directory for exec command
func (cb *ComposeBuilder) Workdir(workdir string) *ComposeBuilder {
	if workdir != "" {
		cb.builder.Flag(FlagWorkdir, workdir)
	}
	return cb
}

// WithFlags adds multiple flags
func (cb *ComposeBuilder) WithFlags(flags ...string) *ComposeBuilder {
	for _, flag := range flags {
		cb.builder.BoolFlag(flag)
	}
	return cb
}

// Up executes compose up
func (cb *ComposeBuilder) Up() error {
	if cb.project != "" {
		cb.builder.Flag(FlagProjectName, cb.project)
	}

	cb.builder.Args(ComposeUpCmd)
	cb.builder.BoolFlag(FlagDetach)

	if len(cb.services) > 0 {
		cb.builder.Args(cb.services...)
	}

	return cb.builder.Build().Run()
}

// Down executes compose down
func (cb *ComposeBuilder) Down() error {
	if cb.project != "" {
		cb.builder.Flag(FlagProjectName, cb.project)
	}

	cb.builder.Args(ComposeDownCmd)
	cb.builder.BoolFlag(FlagVolumes)

	if len(cb.services) > 0 {
		cb.builder.Args(cb.services...)
	}

	return cb.builder.Build().Run()
}

// Logs executes compose logs
func (cb *ComposeBuilder) Logs() error {
	if cb.project != "" {
		cb.builder.Flag(FlagProjectName, cb.project)
	}

	cb.builder.Args(ComposeLogsCmd)

	if len(cb.services) > 0 {
		cb.builder.Args(cb.services...)
	}

	return cb.builder.Build().Run()
}

// Exec executes compose exec
func (cb *ComposeBuilder) Exec(serviceName string, command ...string) *ComposeBuilder {
	if cb.project != "" {
		cb.builder.Flag(FlagProjectName, cb.project)
	}

	cb.builder.Args(ComposeExecCmd)
	cb.builder.Args(serviceName)
	cb.builder.Args(command...)

	return cb
}

// Run executes the built command
func (cb *ComposeBuilder) Run() error {
	return cb.builder.Build().Run()
}
