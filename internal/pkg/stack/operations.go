package stack

import (
	"context"
	"time"

	"github.com/otto-nation/otto-stack/internal/pkg/types/generated"
)

// Operations defines the core stack operations
type Operations interface {
	Up(ctx context.Context, req UpRequest) error
	Down(ctx context.Context, req DownRequest) error
	Restart(ctx context.Context, req RestartRequest) error
	Status(ctx context.Context, req StatusRequest) ([]ServiceStatus, error)
	Logs(ctx context.Context, req LogsRequest) error
}

// UpRequest represents a stack up operation
type UpRequest struct {
	Project        string
	ServiceConfigs []generated.ServiceConfig
	Build          bool
	SkipConflicts  bool
}

// DownRequest represents a stack down operation
type DownRequest struct {
	Project        string
	ServiceConfigs []generated.ServiceConfig
	RemoveVolumes  bool
	Timeout        time.Duration
}

// RestartRequest represents a stack restart operation
type RestartRequest struct {
	Project        string
	ServiceConfigs []generated.ServiceConfig
	Timeout        time.Duration
}

// StatusRequest represents a stack status query
type StatusRequest struct {
	Project        string
	ServiceConfigs []generated.ServiceConfig
}

// LogsRequest represents a stack logs operation
type LogsRequest struct {
	Project        string
	ServiceConfigs []generated.ServiceConfig
	Follow         bool
	Timestamps     bool
	Tail           string
}

// ServiceStatus represents the status of a service
type ServiceStatus struct {
	Name   string
	Status string
	Health string
}
