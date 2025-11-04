package types

import (
	"time"

	"github.com/otto-nation/otto-stack/internal/pkg/constants"
)

// StartOptions defines options for starting services
type StartOptions struct {
	Build          bool
	ForceRecreate  bool
	NoDeps         bool
	Detach         bool
	Timeout        time.Duration
	ResolveDeps    bool
	CheckConflicts bool
}

// StopOptions defines options for stopping services
type StopOptions struct {
	Timeout       int
	Remove        bool
	RemoveVolumes bool
	RemoveOrphans bool
	RemoveImages  string
}

// RestartOptions defines options for restarting services
type RestartOptions struct {
	Timeout time.Duration
	Build   bool
}

// ExecOptions defines options for executing commands in containers
type ExecOptions struct {
	User        string
	WorkingDir  string
	Env         []string
	Interactive bool
	TTY         bool
	Detach      bool
}

// LogOptions defines options for retrieving container logs
type LogOptions struct {
	Follow     bool
	Timestamps bool
	Tail       string
	Since      string
}

// ConnectOptions defines options for connecting to services
type ConnectOptions struct {
	User     string
	Database string
	Host     string
	Port     string
	ReadOnly bool
}

// BackupOptions defines options for backing up service data
type BackupOptions struct {
	OutputDir string
	Compress  bool
	Format    string
	Database  string
	User      string
}

// RestoreOptions defines options for restoring service data
type RestoreOptions struct {
	Database string
	User     string
	Clean    bool
	CreateDB bool
}

// CleanupOptions defines options for cleaning up resources
type CleanupOptions struct {
	RemoveVolumes  bool
	RemoveImages   bool
	RemoveNetworks bool
	All            bool
	DryRun         bool
}

// NewStartOptions returns StartOptions with default values
func NewStartOptions() *StartOptions {
	return &StartOptions{
		Timeout: time.Duration(constants.DefaultStartTimeoutSeconds) * time.Second,
	}
}

// NewStopOptions returns StopOptions with default values
func NewStopOptions() *StopOptions {
	return &StopOptions{
		Timeout: constants.DefaultStopTimeoutSeconds,
	}
}

// NewExecOptions returns ExecOptions with default values
func NewExecOptions() *ExecOptions {
	return &ExecOptions{
		Interactive: true,
		TTY:         true,
	}
}

// NewBackupOptions returns BackupOptions with default values
func NewBackupOptions() *BackupOptions {
	return &BackupOptions{
		Format:   constants.FormatJSON,
		Compress: true,
	}
}
