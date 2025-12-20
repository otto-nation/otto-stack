package core

// NewCLICommandBuilder creates a command builder for CLI commands
func NewCLICommandBuilder(command string) *CommandBuilder {
	return NewCommandBuilder(command)
}

// CLIBuilder provides a fluent interface for building CLI test commands
type CLIBuilder struct {
	builder *CommandBuilder
}

// NewCLIBuilder creates a new CLI command builder
func NewCLIBuilder(command string) *CLIBuilder {
	return &CLIBuilder{
		builder: NewCommandBuilder(command),
	}
}

// Flag adds a flag with value
func (cb *CLIBuilder) Flag(flag, value string) *CLIBuilder {
	cb.builder.Flag(flag, value)
	return cb
}

// BoolFlag adds a boolean flag
func (cb *CLIBuilder) BoolFlag(flag string) *CLIBuilder {
	cb.builder.BoolFlag(flag)
	return cb
}

// Args adds arguments
func (cb *CLIBuilder) Args(args ...string) *CLIBuilder {
	cb.builder.Args(args...)
	return cb
}

// BuildArgs returns the command arguments as a string slice
func (cb *CLIBuilder) BuildArgs() []string {
	return cb.builder.BuildArgs()
}
