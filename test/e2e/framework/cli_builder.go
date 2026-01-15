package framework

const (
	FlagPrefix = "--"
)

// CLIBuilder provides a fluent interface for building CLI commands in tests
type CLIBuilder struct {
	command   string
	flags     map[string]string
	boolFlags []string
	args      []string
}

// NewCLIBuilder creates a new CLI command builder for tests
func NewCLIBuilder(command string) *CLIBuilder {
	return &CLIBuilder{
		command:   command,
		flags:     make(map[string]string),
		boolFlags: make([]string, 0),
		args:      make([]string, 0),
	}
}

// Flag adds a flag with value
func (cb *CLIBuilder) Flag(flag, value string) *CLIBuilder {
	cb.flags[flag] = value
	return cb
}

// BoolFlag adds a boolean flag
func (cb *CLIBuilder) BoolFlag(flag string) *CLIBuilder {
	cb.boolFlags = append(cb.boolFlags, flag)
	return cb
}

// Args adds arguments
func (cb *CLIBuilder) Args(args ...string) *CLIBuilder {
	cb.args = append(cb.args, args...)
	return cb
}

// BuildArgs builds the command arguments array
func (cb *CLIBuilder) BuildArgs() []string {
	result := []string{cb.command}

	// Add boolean flags
	for _, flag := range cb.boolFlags {
		result = append(result, FlagPrefix+flag)
	}

	// Add flags with values
	for flag, value := range cb.flags {
		result = append(result, FlagPrefix+flag, value)
	}

	// Add arguments
	result = append(result, cb.args...)

	return result
}
