package services

import (
	"context"
	"log/slog"
	"time"

	"github.com/docker/compose/v5/pkg/api"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/compose"
	"github.com/otto-nation/otto-stack/internal/pkg/config"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/logger"
	"github.com/otto-nation/otto-stack/internal/pkg/messages"
	servicetypes "github.com/otto-nation/otto-stack/internal/pkg/types"
)

// Service provides high-level stack operations with automatic characteristics resolution
type Service struct {
	compose         api.Compose
	characteristics CharacteristicsResolver
	project         ProjectLoader
	DockerClient    *docker.Client // Exposed for direct access
	logger          *slog.Logger
}

// ServiceInterface defines the interface for service operations
type ServiceInterface interface {
	Start(ctx context.Context, req StartRequest) error
	Stop(ctx context.Context, req StopRequest) error
	Logs(ctx context.Context, req LogRequest) error
	Exec(ctx context.Context, req ExecRequest) error
}

// NewServiceWithDependencies creates a service with injected dependencies (for testing)
// deadcode: used for dependency injection in unit tests
func NewServiceWithDependencies(compose api.Compose, characteristics CharacteristicsResolver, project ProjectLoader, dockerClient *docker.Client) *Service {
	return &Service{
		compose:         compose,
		characteristics: characteristics,
		project:         project,
		DockerClient:    dockerClient,
		logger:          logger.GetLogger(),
	}
}

// ResolveUpServices resolves service names and returns their configs with dependencies
func ResolveUpServices(args []string, cfg *config.Config) ([]servicetypes.ServiceConfig, error) {
	serviceNames := args
	if len(serviceNames) == 0 {
		serviceNames = cfg.Stack.Enabled
	}

	// Load the service manager directly
	manager, err := New()
	if err != nil {
		return nil, err
	}

	// Validate services exist
	if err := manager.ValidateServices(serviceNames); err != nil {
		return nil, pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentStack, messages.ErrorsStackResolveServicesFailed, err)
	}

	// Resolve all dependencies recursively
	resolvedNames := make(map[string]bool)
	var allServiceNames []string

	var resolveDependencies func(string) error
	resolveDependencies = func(serviceName string) error {
		if resolvedNames[serviceName] {
			return nil
		}

		service, err := manager.GetService(serviceName)
		if err != nil {
			return err
		}

		// First resolve dependencies
		for _, dep := range service.Service.Dependencies.Required {
			if err := resolveDependencies(dep); err != nil {
				// Skip missing dependencies (they might be virtual or init containers)
				continue
			}
		}

		// Add this service to output
		if !resolvedNames[serviceName] {
			resolvedNames[serviceName] = true
			allServiceNames = append(allServiceNames, serviceName)
		}

		return nil
	}

	// Resolve dependencies for all requested services
	for _, serviceName := range serviceNames {
		if err := resolveDependencies(serviceName); err != nil {
			return nil, err
		}
	}

	// Load ServiceConfigs for resolved services
	var serviceConfigs []servicetypes.ServiceConfig
	for _, serviceName := range allServiceNames {
		service, err := manager.GetService(serviceName)
		if err != nil {
			continue // Skip services that can't be loaded
		}
		serviceConfigs = append(serviceConfigs, *service)
	}

	return serviceConfigs, nil
}

// StartRequest defines parameters for starting a stack
type StartRequest struct {
	Project         string
	ServiceConfigs  []servicetypes.ServiceConfig
	Build           bool
	ForceRecreate   bool
	Detach          bool
	NoDeps          bool
	Timeout         time.Duration
	Characteristics []string
}

// StopRequest defines parameters for stopping a stack
type StopRequest struct {
	Project         string
	ServiceConfigs  []servicetypes.ServiceConfig
	Remove          bool // true = down, false = stop
	RemoveVolumes   bool
	RemoveOrphans   bool
	Timeout         time.Duration
	Characteristics []string
}

// StatusRequest defines parameters for getting service status
type StatusRequest struct {
	Project  string
	Services []string
}

// CleanupRequest defines parameters for cleanup operations
type CleanupRequest struct {
	Project       string
	Force         bool
	RemoveVolumes bool
	RemoveImages  bool
}

// ExecRequest defines parameters for executing commands in containers
type ExecRequest struct {
	Project     string
	Service     string
	Command     []string
	User        string
	WorkingDir  string
	Interactive bool
	TTY         bool
}
type LogRequest struct {
	Project        string
	ServiceConfigs []servicetypes.ServiceConfig
	Follow         bool
	Timestamps     bool
	Tail           string
}

// NewService creates a new stack service
func NewService(compose api.Compose, characteristics CharacteristicsResolver, project ProjectLoader) (*Service, error) {
	dockerClient, err := docker.NewClient(nil)
	if err != nil {
		return nil, pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentDocker, messages.ErrorsStackDockerClientFailed, err)
	}

	return &Service{
		compose:         compose,
		characteristics: characteristics,
		project:         project,
		DockerClient:    dockerClient,
		logger:          logger.GetLogger(),
	}, nil
}

// Start starts services with automatic characteristics resolution
func (s *Service) Start(ctx context.Context, req StartRequest) error {
	s.logger.Debug("Starting services",
		"project", req.Project,
		"serviceCount", len(req.ServiceConfigs),
		"build", req.Build,
		"forceRecreate", req.ForceRecreate)

	// Load and validate service configs from .otto-stack/service-configs/
	req.ServiceConfigs = s.loadAndValidateServiceConfigs(req.ServiceConfigs)

	// Generate docker-compose.yml from service configs
	if err := s.GenerateComposeFile(req.Project, req.ServiceConfigs); err != nil {
		return pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentProject, messages.ErrorsStackComposeGenerateFailed, err)
	}

	// Load project
	project, err := s.project.Load(req.Project)
	if err != nil {
		return pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentProject, messages.ErrorsStackProjectLoadFailed, err)
	}

	s.logger.Debug("Project loaded successfully", pkgerrors.ComponentProject, req.Project)

	// Resolve characteristics to options and convert to SDK format
	options := s.characteristics.ResolveUpOptions(req.Characteristics, req.ServiceConfigs, docker.UpOptions{
		Build:         req.Build,
		ForceRecreate: req.ForceRecreate,
		Detach:        req.Detach,
		NoDeps:        req.NoDeps,
		Timeout:       &req.Timeout,
	})

	err = s.compose.Up(ctx, project, options.ToSDK())
	if err != nil {
		if len(req.ServiceConfigs) > 0 {
			return pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentProject, messages.ErrorsStackServicesStartFailed, err)
		}
		return pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentProject, messages.ErrorsStackProjectStartFailed, err)
	}

	s.logger.Debug("Services started successfully")

	// Execute local init scripts for services that have them
	if err := s.executeLocalInitScripts(ctx, req.ServiceConfigs, req.Project); err != nil {
		return pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentProject, messages.ErrorsStackInitScriptsFailed, err)
	}

	return nil
}

// executeLocalInitScripts executes local init scripts for all services that have them
// Stop stops services with automatic characteristics resolution
func (s *Service) Stop(ctx context.Context, req StopRequest) error {
	s.logger.Debug("Stopping services",
		"project", req.Project,
		messages.ErrorsStackProjectRemoveFailed, req.Remove,
		"serviceCount", len(req.ServiceConfigs))

	// Load project
	project, err := s.project.Load(req.Project)
	if err != nil {
		return pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentProject, messages.ErrorsStackProjectLoadFailed, err)
	}

	if req.Remove {
		// Use down operation
		options := s.characteristics.ResolveDownOptions(req.Characteristics, req.ServiceConfigs, docker.DownOptions{
			RemoveVolumes: req.RemoveVolumes,
			Timeout:       &req.Timeout,
		})
		err = s.compose.Down(ctx, project.Name, options.ToSDK())
		if err != nil {
			return pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentProject, messages.ErrorsStackProjectRemoveFailed, err)
		}
		return nil
	}

	// Use stop operation
	stopOptions := s.characteristics.ResolveStopOptions(req.Characteristics, req.ServiceConfigs, docker.StopOptions{
		Timeout: &req.Timeout,
	})
	err = s.compose.Stop(ctx, project.Name, stopOptions.ToSDK())
	if err != nil {
		serviceNames := ExtractServiceNames(req.ServiceConfigs)
		if len(serviceNames) > 0 {
			return pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentProject, messages.ErrorsStackStopFailed, err)
		}
		return pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentProject, messages.ErrorsStackProjectStopFailed, err)
	}
	return nil
}

// Logs retrieves logs from services
func (s *Service) Logs(ctx context.Context, req LogRequest) error {
	serviceNames := ExtractServiceNames(req.ServiceConfigs)

	options := docker.LogOptions{
		Services:   serviceNames,
		Follow:     req.Follow,
		Timestamps: req.Timestamps,
		Tail:       req.Tail,
	}
	consumer := &docker.SimpleLogConsumer{}
	err := s.compose.Logs(ctx, req.Project, consumer, options.ToSDK())
	if err != nil {
		return pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentProject, messages.ErrorsStackGetLogsFailed, err)
	}
	return nil
}

// Exec executes commands in service containers
func (s *Service) Exec(ctx context.Context, req ExecRequest) error {
	// Load project
	project, err := s.project.Load(req.Project)
	if err != nil {
		return pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentProject, messages.ErrorsStackProjectLoadFailed, err)
	}

	// Use the compose SDK's exec functionality
	options := api.RunOptions{
		Project:     project,
		Service:     req.Service,
		Command:     req.Command,
		User:        req.User,
		WorkingDir:  req.WorkingDir,
		Interactive: req.Interactive,
		Tty:         req.TTY,
		Index:       1, // Default to first container instance
	}

	_, err = s.compose.Exec(ctx, req.Project, options)
	if err != nil {
		return pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentProject, messages.ErrorsStackExecCommandFailed, err)
	}
	return nil
}

// Status retrieves status of services
func (s *Service) Status(ctx context.Context, req StatusRequest) ([]docker.ContainerStatus, error) {
	return s.DockerClient.GetServiceStatus(ctx, req.Project, req.Services)
}

// Cleanup removes containers and resources for a project
func (s *Service) Cleanup(ctx context.Context, req CleanupRequest) error {
	// List containers
	containers, err := s.DockerClient.ListContainers(ctx, req.Project)
	if err != nil {
		return pkgerrors.New(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentDocker, "list containers", err)
	}

	// Remove containers
	for _, container := range containers {
		if err := s.DockerClient.RemoveContainer(ctx, container.ID, req.Force); err != nil {
			return pkgerrors.New(pkgerrors.ErrCodeOperationFail, container.Name, "remove container", err)
		}
	}

	// Remove volumes if requested
	if req.RemoveVolumes {
		if err := s.DockerClient.RemoveResources(ctx, docker.ResourceVolume, req.Project); err != nil {
			return pkgerrors.New(pkgerrors.ErrCodeOperationFail, req.Project, "remove volumes", err)
		}
	}

	// Remove networks
	if err := s.DockerClient.RemoveResources(ctx, docker.ResourceNetwork, req.Project); err != nil {
		return pkgerrors.New(pkgerrors.ErrCodeOperationFail, req.Project, "remove networks", err)
	}

	// Remove images if requested
	if req.RemoveImages {
		if err := s.DockerClient.RemoveResources(ctx, docker.ResourceImage, req.Project); err != nil {
			return pkgerrors.New(pkgerrors.ErrCodeOperationFail, req.Project, "remove images", err)
		}
	}

	return nil
}

// CheckDockerHealth verifies Docker daemon is running
func (s *Service) CheckDockerHealth(ctx context.Context) error {
	_, err := s.DockerClient.GetCli().Info(ctx)
	if err != nil {
		return pkgerrors.New(pkgerrors.ErrCodeUnavailable, pkgerrors.ComponentDocker, "check docker health", err)
	}
	return nil
}

// GenerateComposeFile generates docker-compose.yml from service configs
func (s *Service) GenerateComposeFile(projectName string, serviceConfigs []servicetypes.ServiceConfig) error {
	generator, err := compose.NewGenerator(projectName)
	if err != nil {
		return err
	}

	return generator.GenerateFromServiceConfigs(serviceConfigs, projectName)
}
