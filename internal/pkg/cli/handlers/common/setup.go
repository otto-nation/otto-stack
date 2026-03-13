package common

import (
	"context"
	"fmt"

	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
	"github.com/otto-nation/otto-stack/internal/pkg/config"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/messages"
	"github.com/otto-nation/otto-stack/internal/pkg/registry"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/otto-nation/otto-stack/internal/pkg/validation"
)

// CoreSetup contains common setup data for handlers
type CoreSetup struct {
	Config       *config.Config
	DockerClient *docker.Client
}

// SetupCoreCommand provides common setup for handlers that need Docker and config
func SetupCoreCommand(ctx context.Context, base *base.BaseCommand) (*CoreSetup, func(), error) {
	// Check if otto-stack is initialized
	if err := validation.CheckInitialization(); err != nil {
		return nil, nil, err
	}

	// Load project configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, nil, err
	}

	// Create Docker client
	dockerClient, err := docker.NewClient(nil)
	if err != nil {
		return nil, nil, pkgerrors.NewDockerError(pkgerrors.ErrCodeOperationFail, messages.ErrorsDockerClientCreateFailed, err)
	}

	cleanup := func() {
		_ = dockerClient.Close()
	}

	return &CoreSetup{
		Config:       cfg,
		DockerClient: dockerClient,
	}, cleanup, nil
}

// ResolveServiceConfigs resolves service configurations from args or enabled services
func ResolveServiceConfigs(args []string, setup *CoreSetup) ([]types.ServiceConfig, error) {
	if len(args) > 0 {
		return services.ResolveUpServices(args, setup.Config)
	}
	return services.ResolveUpServices(setup.Config.Stack.Enabled, setup.Config)
}

// DetectExecutionContext creates a detector and returns the current execution context.
func DetectExecutionContext() (clicontext.ExecutionMode, error) {
	detector, err := clicontext.NewDetector()
	if err != nil {
		return nil, pkgerrors.NewSystemError(pkgerrors.ErrCodeInternal, messages.ErrorsContextDetectorCreateFailed, err)
	}
	execCtx, err := detector.DetectContext()
	if err != nil {
		return nil, pkgerrors.NewSystemError(pkgerrors.ErrCodeInternal, messages.ErrorsContextDetectFailed, err)
	}
	return execCtx, nil
}

// VerifyServicesInRegistry checks that all named services exist in the shared registry.
// Returns an error naming the first service not found.
func VerifyServicesInRegistry(serviceNames []string, reg *registry.Registry) error {
	for _, name := range serviceNames {
		if _, exists := reg.Containers[name]; !exists {
			return pkgerrors.NewValidationError(
				pkgerrors.ErrCodeInvalid,
				pkgerrors.FieldServiceName,
				fmt.Sprintf(messages.SharedServiceNotInRegistry, name),
				nil,
			)
		}
	}
	return nil
}
