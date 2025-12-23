package services

import (
	"context"
	"fmt"
	"time"

	"github.com/compose-spec/compose-go/v2/types"
	"github.com/docker/compose/v5/pkg/api"
	"github.com/otto-nation/otto-stack/internal/core/docker"
)

// Service provides high-level stack operations with automatic characteristics resolution
type Service struct {
	compose         api.Compose
	characteristics CharacteristicsResolver
	project         ProjectLoader
	DockerClient    *docker.Client // Exposed for direct access
}

// StartRequest defines parameters for starting a stack
type StartRequest struct {
	Project         string
	Services        []string
	Build           bool
	ForceRecreate   bool
	Characteristics []string
}

// StopRequest defines parameters for stopping a stack
type StopRequest struct {
	Project         string
	Services        []string
	Remove          bool // true = down, false = stop
	RemoveVolumes   bool
	Timeout         time.Duration
	Characteristics []string
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
	Project    string
	Services   []string
	Follow     bool
	Timestamps bool
	Tail       string
}

// NewService creates a new stack service
func NewService(compose api.Compose, characteristics CharacteristicsResolver, project ProjectLoader) (*Service, error) {
	dockerClient, err := docker.NewClient(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %w", err)
	}

	return &Service{
		compose:         compose,
		characteristics: characteristics,
		project:         project,
		DockerClient:    dockerClient,
	}, nil
}

// Start starts services with automatic characteristics resolution
func (s *Service) Start(ctx context.Context, req StartRequest) error {
	// Load project
	project, err := s.project.Load(req.Project)
	if err != nil {
		return fmt.Errorf("failed to load project %s: %w", req.Project, err)
	}

	// Filter services if specified
	if len(req.Services) > 0 {
		project = s.filterServices(project, req.Services)
	}

	// Resolve characteristics to options and convert to SDK format
	options := s.characteristics.ResolveUpOptions(req.Characteristics, UpOptions{
		Build:         req.Build,
		ForceRecreate: true, // Force recreate for diagnosis
	})

	err = s.compose.Up(ctx, project, options.ToSDK())
	if err != nil {
		if len(req.Services) > 0 {
			return fmt.Errorf("failed to start services %v in project %s: %w", req.Services, req.Project, err)
		}
		return fmt.Errorf("failed to start project %s: %w", req.Project, err)
	}
	return nil
}

// Stop stops services with automatic characteristics resolution
func (s *Service) Stop(ctx context.Context, req StopRequest) error {
	// Load project
	project, err := s.project.Load(req.Project)
	if err != nil {
		return fmt.Errorf("failed to load project %s: %w", req.Project, err)
	}

	// Filter services if specified
	if len(req.Services) > 0 {
		project = s.filterServices(project, req.Services)
	}

	if req.Remove {
		// Use down operation
		options := s.characteristics.ResolveDownOptions(req.Characteristics, DownOptions{
			RemoveVolumes: req.RemoveVolumes,
			Timeout:       req.Timeout,
		})
		err = s.compose.Down(ctx, project.Name, options.ToSDK())
		if err != nil {
			return fmt.Errorf("failed to remove project %s: %w", req.Project, err)
		}
		return nil
	}

	// Use stop operation
	options := s.characteristics.ResolveStopOptions(req.Characteristics, StopOptions{
		Services: req.Services,
		Timeout:  req.Timeout,
	})
	err = s.compose.Stop(ctx, project.Name, options.ToSDK())
	if err != nil {
		if len(req.Services) > 0 {
			return fmt.Errorf("failed to stop services %v in project %s: %w", req.Services, req.Project, err)
		}
		return fmt.Errorf("failed to stop project %s: %w", req.Project, err)
	}
	return nil
}

// Logs retrieves logs from services
func (s *Service) Logs(ctx context.Context, req LogRequest) error {
	options := LogOptions{
		Services:   req.Services,
		Follow:     req.Follow,
		Timestamps: req.Timestamps,
		Tail:       req.Tail,
	}
	consumer := &docker.SimpleLogConsumer{}
	err := s.compose.Logs(ctx, req.Project, consumer, options.ToSDK())
	if err != nil {
		if len(req.Services) > 0 {
			return fmt.Errorf("failed to get logs for services %v in project %s: %w", req.Services, req.Project, err)
		}
		return fmt.Errorf("failed to get logs for project %s: %w", req.Project, err)
	}
	return nil
}

// Exec executes commands in service containers
func (s *Service) Exec(ctx context.Context, req ExecRequest) error {
	// Load project
	project, err := s.project.Load(req.Project)
	if err != nil {
		return fmt.Errorf("failed to load project %s: %w", req.Project, err)
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
		return fmt.Errorf("failed to exec command %v in service %s (project %s): %w", req.Command, req.Service, req.Project, err)
	}
	return nil
}

func (s *Service) filterServices(project *types.Project, services []string) *types.Project {
	filtered := make(types.Services)
	for _, service := range services {
		if serviceConfig, exists := project.Services[service]; exists {
			filtered[service] = serviceConfig
		}
	}
	project.Services = filtered
	return project
}
