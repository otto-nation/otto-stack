package core

import (
	"context"
	"os/exec"
	"strings"
)

// CommandBuilder provides a fluent interface for building shell commands
type CommandBuilder struct {
	command     string
	subcommands []string
	args        []string
	flags       map[string]string
	boolFlags   []string
	ctx         context.Context
}

// NewCommandBuilder creates a new command builder
func NewCommandBuilder(command string) *CommandBuilder {
	return &CommandBuilder{
		command:   command,
		flags:     make(map[string]string),
		boolFlags: make([]string, 0),
	}
}

// Subcommand adds a subcommand
func (cb *CommandBuilder) Subcommand(subcommand string) *CommandBuilder {
	cb.subcommands = append(cb.subcommands, subcommand)
	return cb
}

// Args adds arguments
func (cb *CommandBuilder) Args(args ...string) *CommandBuilder {
	cb.args = append(cb.args, args...)
	return cb
}

// Flag adds a flag with value
func (cb *CommandBuilder) Flag(flag, value string) *CommandBuilder {
	cb.flags[flag] = value
	return cb
}

// BoolFlag adds a boolean flag
func (cb *CommandBuilder) BoolFlag(flag string) *CommandBuilder {
	cb.boolFlags = append(cb.boolFlags, flag)
	return cb
}

// Context sets the context for the command
func (cb *CommandBuilder) Context(ctx context.Context) *CommandBuilder {
	cb.ctx = ctx
	return cb
}

// Build constructs the final command
func (cb *CommandBuilder) Build() *exec.Cmd {
	args := make([]string, 0)

	// Add subcommands
	args = append(args, cb.subcommands...)

	// Add flags with values
	for flag, value := range cb.flags {
		args = append(args, "--"+flag, value)
	}

	// Add boolean flags
	for _, flag := range cb.boolFlags {
		args = append(args, "--"+flag)
	}

	// Add arguments
	args = append(args, cb.args...)

	var cmd *exec.Cmd
	if cb.ctx != nil {
		cmd = exec.CommandContext(cb.ctx, cb.command, args...)
	} else {
		cmd = exec.Command(cb.command, args...)
	}

	return cmd
}

// BuildArgs returns just the arguments as a string slice
func (cb *CommandBuilder) BuildArgs() []string {
	args := make([]string, 0)

	// Add command
	args = append(args, cb.command)

	// Add subcommands
	args = append(args, cb.subcommands...)

	// Add flags with values
	for flag, value := range cb.flags {
		args = append(args, "--"+flag, value)
	}

	// Add boolean flags
	for _, flag := range cb.boolFlags {
		args = append(args, "--"+flag)
	}

	// Add arguments
	args = append(args, cb.args...)

	return args
}

// String returns the command as a string
func (cb *CommandBuilder) String() string {
	return strings.Join(cb.BuildArgs(), " ")
}
