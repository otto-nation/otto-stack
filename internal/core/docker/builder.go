package docker

import (
	"context"
	"os/exec"
	"strings"

	"github.com/otto-nation/otto-stack/internal/core"
)

// CommandBuilder provides a fluent interface for building docker commands
type CommandBuilder struct {
	cmd    string
	subcmd string
	flags  []string
	args   []string
	env    map[string]string
	ctx    context.Context
}

// NewCommand creates a new docker command builder
func NewCommand(cmd string) *CommandBuilder {
	return &CommandBuilder{
		cmd: cmd,
		env: make(map[string]string),
	}
}

// Subcommand sets the docker subcommand
func (cb *CommandBuilder) Subcommand(subcmd string) *CommandBuilder {
	cb.subcmd = subcmd
	return cb
}

// Flag adds a flag with optional values
func (cb *CommandBuilder) Flag(name string, values ...string) *CommandBuilder {
	flag := "--" + name
	cb.flags = append(cb.flags, flag)
	cb.flags = append(cb.flags, values...)
	return cb
}

// BoolFlag adds a boolean flag (no value)
func (cb *CommandBuilder) BoolFlag(name string) *CommandBuilder {
	cb.flags = append(cb.flags, "--"+name)
	return cb
}

// Args adds arguments to the command
func (cb *CommandBuilder) Args(args ...string) *CommandBuilder {
	cb.args = append(cb.args, args...)
	return cb
}

// Env sets environment variables
func (cb *CommandBuilder) Env(key, value string) *CommandBuilder {
	cb.env[key] = value
	return cb
}

// Context sets the context for the command
func (cb *CommandBuilder) Context(ctx context.Context) *CommandBuilder {
	cb.ctx = ctx
	return cb
}

// Build creates the exec.Cmd
func (cb *CommandBuilder) Build() *exec.Cmd {
	cmdArgs := []string{}

	if cb.subcmd != "" {
		cmdArgs = append(cmdArgs, cb.subcmd)
	}

	cmdArgs = append(cmdArgs, cb.flags...)
	cmdArgs = append(cmdArgs, cb.args...)

	var cmd *exec.Cmd
	if cb.ctx != nil {
		cmd = exec.CommandContext(cb.ctx, cb.cmd, cmdArgs...)
	} else {
		cmd = exec.Command(cb.cmd, cmdArgs...)
	}

	// Set environment variables
	if len(cb.env) > 0 {
		env := cmd.Environ()
		for key, value := range cb.env {
			env = append(env, key+"="+value)
		}
		cmd.Env = env
	}

	return cmd
}

// Run executes the command
func (cb *CommandBuilder) Run() error {
	return cb.Build().Run()
}

// ComposeBuilder provides compose-specific functionality
type ComposeBuilder struct {
	*CommandBuilder
	project  string
	services []string
}

// NewComposeBuilder creates a new compose command builder
func NewComposeBuilder() *ComposeBuilder {
	return &ComposeBuilder{
		CommandBuilder: NewCommand(DockerCmd).Subcommand(DockerComposeCmd),
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

// WithFlags adds multiple flags
func (cb *ComposeBuilder) WithFlags(flags ...string) *ComposeBuilder {
	for _, flag := range flags {
		if strings.HasPrefix(flag, "--") {
			cb.flags = append(cb.flags, flag)
		} else {
			cb.BoolFlag(flag)
		}
	}
	return cb
}

// Up executes compose up
func (cb *ComposeBuilder) Up() error {
	if cb.project != "" {
		cb.Flag("project-name", cb.project)
	}

	cb.Args(core.CommandUp)
	cb.BoolFlag(core.FlagDetach)

	if len(cb.services) > 0 {
		cb.Args(cb.services...)
	}

	return cb.Run()
}

// Down executes compose down
func (cb *ComposeBuilder) Down() error {
	if cb.project != "" {
		cb.Flag("project-name", cb.project)
	}

	cb.Args(core.CommandDown)
	cb.BoolFlag(core.FlagVolumes)

	return cb.Run()
}

// Logs executes compose logs
func (cb *ComposeBuilder) Logs() error {
	if cb.project != "" {
		cb.Flag("project-name", cb.project)
	}

	cb.Args(core.CommandLogs)

	if len(cb.services) > 0 {
		cb.Args(cb.services...)
	}

	return cb.Run()
}
