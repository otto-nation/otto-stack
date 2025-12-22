package docker

import (
	"fmt"
	"log/slog"
	"os/exec"
	"strings"

	"github.com/otto-nation/otto-stack/internal/core"
)

// NewCommand creates a new Docker command builder
func NewCommand(command string) *core.CommandBuilder {
	return core.NewCommandBuilder(command)
}

// ComposeBuilder provides a fluent interface for building docker-compose commands
type ComposeBuilder struct {
	builder       *core.CommandBuilder
	project       string
	services      []string
	detach        bool
	verbose       bool
	timeout       int
	removeVolumes bool
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

// Timeout sets the timeout for operations
func (cb *ComposeBuilder) Timeout(timeout int) *ComposeBuilder {
	cb.timeout = timeout
	return cb
}

// RemoveVolumes sets whether to remove volumes
func (cb *ComposeBuilder) RemoveVolumes(remove bool) *ComposeBuilder {
	cb.removeVolumes = remove
	return cb
}

// File sets the compose file path
func (cb *ComposeBuilder) File(path string) *ComposeBuilder {
	cb.builder.Flag(FlagFile, path)
	return cb
}

// Detach sets whether to run in detached mode
func (cb *ComposeBuilder) Detach(detach bool) *ComposeBuilder {
	cb.detach = detach
	return cb
}

// Verbose sets whether to enable verbose logging
func (cb *ComposeBuilder) Verbose(verbose bool) *ComposeBuilder {
	cb.verbose = verbose
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
	args := []string{DockerComposeCmd}

	// Global flags
	if cb.project != "" {
		args = append(args, FlagPrefix+core.FlagProjectName, cb.project)
	}
	args = append(args, FlagPrefix+FlagFile, DockerComposeFilePath)

	// Subcommand
	args = append(args, ComposeUpCmd)

	// Subcommand flags
	if cb.detach {
		args = append(args, FlagPrefix+core.FlagDetach)
	}

	// Services
	args = append(args, cb.services...)

	if cb.verbose {
		slog.Info("Executing "+DockerCmd+" "+DockerComposeCmd+" "+ComposeUpCmd, "command", strings.Join(append([]string{DockerCmd}, args...), " "))
	}
	return exec.Command(DockerCmd, args...).Run()
}

// Down executes compose down
func (cb *ComposeBuilder) Down() error {
	if cb.project != "" {
		cb.builder.Flag(core.FlagProjectName, cb.project)
	}

	cb.builder.Args(ComposeDownCmd)

	if cb.removeVolumes {
		cb.builder.SubcommandBoolFlag(core.FlagVolumes)
	}

	if len(cb.services) > 0 {
		cb.builder.Args(cb.services...)
	}

	// Add timeout if specified
	if cb.timeout > 0 {
		cb.builder.SubcommandFlag(FlagTimeout, fmt.Sprintf("%d", cb.timeout))
	}

	return cb.builder.Build().Run()
}

// Stop executes compose stop
func (cb *ComposeBuilder) Stop() error {
	if cb.project != "" {
		cb.builder.Flag(core.FlagProjectName, cb.project)
	}

	cb.builder.Args(ComposeStopCmd)

	if len(cb.services) > 0 {
		cb.builder.Args(cb.services...)
	}

	// Add timeout if specified
	if cb.timeout > 0 {
		cb.builder.SubcommandFlag(FlagTimeout, fmt.Sprintf("%d", cb.timeout))
	}

	err := cb.builder.Build().Run()
	// Make stop idempotent - ignore errors when no containers are running
	if err != nil && isNoContainersError(err) {
		return nil
	}
	return err
}

// isNoContainersError checks if error is due to no containers running
func isNoContainersError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	// Docker compose stop returns exit status 1 but with no stderr when no containers
	// This is a common pattern for idempotent operations
	return errStr == "exit status 1"
}

// Logs executes compose logs
func (cb *ComposeBuilder) Logs() error {
	if cb.project != "" {
		cb.builder.Flag(core.FlagProjectName, cb.project)
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
		cb.builder.Flag(core.FlagProjectName, cb.project)
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
