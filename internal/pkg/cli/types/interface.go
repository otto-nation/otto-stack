package types

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// CommandHandler defines the interface for all command handlers
type CommandHandler interface {
	// Handle executes the command with the given context, command, arguments, and base command
	Handle(ctx context.Context, cmd *cobra.Command, args []string, base *BaseCommand) error

	// ValidateArgs validates the command arguments before execution
	ValidateArgs(args []string) error

	// GetRequiredFlags returns a list of required flags for this command
	GetRequiredFlags() []string
}

// BaseCommand provides common functionality for all commands
type BaseCommand struct {
	ProjectDir string
	Manager    ServiceManager
	Logger     Logger
}

// Close cleans up resources
func (b *BaseCommand) Close() error {
	if b.Manager != nil {
		return b.Manager.Close()
	}
	return nil
}

// ServiceManager interface for service operations
type ServiceManager interface {
	StartServices(ctx context.Context, serviceNames []string, options StartOptions) error
	StopServices(ctx context.Context, serviceNames []string, options StopOptions) error
	GetServiceStatus(ctx context.Context, serviceNames []string) ([]ServiceStatus, error)
	Close() error
}

// Logger interface for logging operations
type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
}

// ValidateServices validates service names against available services
func (b *BaseCommand) ValidateServices(serviceNames []string) error {
	// Try to load services from embedded config first
	servicesFile := "internal/config/services/services.yaml"

	// Check if we're in the otto-stack project directory
	if _, err := os.Stat(filepath.Join(b.ProjectDir, servicesFile)); os.IsNotExist(err) {
		// We're not in the otto-stack project directory, skip validation for now
		// TODO: Use embedded services configuration
		b.Logger.Debug("services.yaml not found, skipping service validation")
		return nil
	}

	// We're in the otto-stack project directory, use local services.yaml
	fullPath := filepath.Join(b.ProjectDir, servicesFile)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return fmt.Errorf("failed to read services.yaml: %w", err)
	}

	var servicesConfig map[string]interface{}
	if err := yaml.Unmarshal(data, &servicesConfig); err != nil {
		return fmt.Errorf("failed to parse services.yaml: %w", err)
	}

	// Check each service name
	for _, serviceName := range serviceNames {
		if _, exists := servicesConfig[serviceName]; !exists {
			availableServices := make([]string, 0, len(servicesConfig))
			for name := range servicesConfig {
				availableServices = append(availableServices, name)
			}
			return fmt.Errorf("unknown service '%s'. Available services: %v", serviceName, availableServices)
		}
	}

	return nil
}
