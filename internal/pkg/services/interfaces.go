package services

import (
	"context"

	"github.com/otto-nation/otto-stack/internal/core/docker"
)

// ManagerAdapter interface for service management operations
type ManagerAdapter interface {
	StartServices(ctx context.Context, serviceNames []string, options docker.StartOptions) error
	StopServices(ctx context.Context, serviceNames []string, options docker.StopOptions) error
	GetServiceStatus(ctx context.Context, serviceNames []string) ([]docker.DockerServiceStatus, error)
}
