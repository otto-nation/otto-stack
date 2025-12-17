package stack

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/ci"
	"github.com/otto-nation/otto-stack/internal/pkg/compose"
	"github.com/otto-nation/otto-stack/internal/pkg/config"
	"github.com/otto-nation/otto-stack/internal/pkg/env"
	"github.com/otto-nation/otto-stack/internal/pkg/logger"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/internal/scripts"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const (
	// DefaultTimeoutSeconds is the default timeout for operations
	DefaultTimeoutSeconds = 30
)

// StackState tracks the current state of the stack
type StackState struct {
	ConfigHash string   `json:"config_hash"`
	Services   []string `json:"services"`
}

// UpHandler handles the up command
type UpHandler struct{}

// NewUpHandler creates a new up handler
func NewUpHandler() *UpHandler {
	return &UpHandler{}
}

// Handle executes the up command
func (h *UpHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
	// Check initialization first

	// Check initialization first, before any output
	setup, cleanup, err := SetupCoreCommand(ctx, base)
	if err != nil {
		return err // Return error directly without logging or headers
	}
	defer cleanup()

	// Start operation logging only after initialization check passes
	logger.Info(logger.LogMsgStartingOperation, logger.LogFieldOperation, logger.OperationStackUp, logger.LogFieldServices, args)

	ciFlags := ci.GetFlags(cmd)

	if ciFlags.DryRun {
		base.Output.Info("%s", core.MsgDry_run_showing_what_would_happen)

		// Determine services that would be started
		serviceNames := args
		if len(serviceNames) == 0 {
			serviceNames = setup.Config.Stack.Enabled
		}

		base.Output.Info(core.MsgDry_run_would_start_services, fmt.Sprintf("%v", serviceNames))
		base.Output.Info(core.MsgDry_run_would_use_config, filepath.Join(core.OttoStackDir, core.ConfigFileName))
		return nil
	}
	defer func() {
		if r := recover(); r != nil {
			logger.Error(logger.LogMsgOperationFailed, logger.LogFieldOperation, logger.OperationStackUp, logger.LogFieldError, fmt.Errorf("panic: %v", r))
			panic(r)
		}
	}()

	base.Output.Header("%s", core.MsgStarting)
	logger.Info(logger.LogMsgServiceAction, logger.LogFieldAction, logger.ActionStart, logger.LogFieldService, "stack", logger.LogFieldServices, args)

	// Parse all flags with validation
	flags, err := core.ParseUpFlags(cmd)
	if err != nil {
		logger.Error(logger.LogMsgOperationFailed, logger.LogFieldOperation, logger.OperationStackUp, logger.LogFieldError, err)
		return err
	}

	// Parse timeout from string to duration
	timeoutSecs := h.parseTimeoutSeconds(flags.Timeout)

	// Clean usage with no repetitive error handling
	options := docker.StartOptions{
		Build:          flags.Build,
		ForceRecreate:  flags.ForceRecreate,
		Detach:         flags.Detach,
		Timeout:        time.Duration(timeoutSecs) * time.Second,
		NoDeps:         flags.NoDeps,
		ResolveDeps:    flags.ResolveDeps,
		CheckConflicts: flags.CheckConflicts,
		RemoveOrphans:  flags.ForceRecreate, // Auto-remove orphans when force recreating
	}

	// Determine services to start
	serviceNames := args
	if len(serviceNames) == 0 {
		serviceNames = setup.Config.Stack.Enabled
	}

	// Filter services to only include container services
	serviceUtils := services.NewServiceUtils()
	filteredServices, err := serviceUtils.ResolveServices(serviceNames)
	if err != nil {
		return fmt.Errorf(core.MsgStack_failed_resolve_services, err)
	}

	// Check for config changes
	configHash, err := h.getConfigHash(setup.Config)
	if err != nil {
		return fmt.Errorf(core.MsgStack_failed_calculate_hash, err)
	}

	previousState, err := h.loadState()
	if err != nil {
		// Restart operation
		previousState = &StackState{}
	}

	configChanged := previousState.ConfigHash != configHash
	h.handleConfigChange(ctx, setup, previousState.Services, filteredServices, base, configChanged)

	// Generate compose file
	generator, err := compose.NewGenerator(setup.Config.Project.Name, services.ServicesDir)
	if err != nil {
		return fmt.Errorf(core.MsgStack_failed_create_generator, err)
	}

	composeData, err := generator.GenerateYAML(serviceNames)
	if err != nil {
		return fmt.Errorf(core.MsgStack_failed_generate_compose, err)
	}

	// Ensure otto-stack directory exists
	if err := os.MkdirAll(core.OttoStackDir, core.PermReadWriteExec); err != nil {
		return fmt.Errorf(core.MsgStack_failed_create_directory, err)
	}

	composePath := docker.DockerComposeFilePath
	if err := os.MkdirAll(filepath.Dir(composePath), core.PermReadWriteExec); err != nil {
		return fmt.Errorf("failed to create generated directory: %w", err)
	}
	if err := os.WriteFile(composePath, composeData, core.PermReadWrite); err != nil {
		return fmt.Errorf(core.MsgStack_failed_write_compose, err)
	}

	// Generate .env.generated file
	if err := h.generateEnvFile(filteredServices, setup.Config.Project.Name); err != nil {
		base.Output.Warning("Failed to generate .env file: %v", err)
	}

	// Start services first
	if err := setup.DockerClient.ComposeUp(ctx, setup.Config.Project.Name, filteredServices, options); err != nil {
		return fmt.Errorf(core.MsgStack_failed_start_services, err)
	}

	// Run init containers after main services are started
	if err := h.runInitContainers(ctx, setup, filteredServices, base); err != nil {
		base.Output.Warning("Failed to run init containers: %v", err)
	}

	// Save new state
	newState := &StackState{
		ConfigHash: configHash,
		Services:   filteredServices,
	}
	if err := h.saveState(newState); err != nil {
		base.Output.Warning("Failed to save state: %v", err)
	}

	base.Output.Success(core.MsgStartSuccess)
	logger.Info(logger.LogMsgOperationCompleted, logger.LogFieldOperation, logger.OperationStackUp)
	return nil
}

// ValidateArgs validates the command arguments
func (h *UpHandler) ValidateArgs(args []string) error {
	// Service names are optional - if none provided, all enabled services are used
	return nil
}

// GetRequiredFlags returns required flags for this command
func (h *UpHandler) GetRequiredFlags() []string {
	// No flags are strictly required for the up command
	return []string{}
}

// getConfigHash calculates hash of current config
func (h *UpHandler) getConfigHash(config *config.Config) (string, error) {
	data, err := yaml.Marshal(config)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash), nil
}

// loadState loads previous stack state
func (h *UpHandler) loadState() (*StackState, error) {
	// Try new location first
	data, err := os.ReadFile(core.StateFilePath)
	if err != nil {
		// Fall back to old location for backward compatibility
		oldStatePath := filepath.Join(core.OttoStackDir, core.StateFileName)
		data, err = os.ReadFile(oldStatePath)
		if err != nil {
			return nil, err
		}
	}

	var state StackState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}
	return &state, nil
}

// saveState saves current stack state
func (h *UpHandler) saveState(state *StackState) error {
	if err := os.MkdirAll(filepath.Dir(core.StateFilePath), core.PermReadWriteExec); err != nil {
		return fmt.Errorf("failed to create generated directory: %w", err)
	}
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	return os.WriteFile(core.StateFilePath, data, core.PermReadWrite)
}

// cleanupRemovedServices removes services no longer in config
func (h *UpHandler) cleanupRemovedServices(ctx context.Context, setup *CoreSetup, oldServices, newServices []string, base *base.BaseCommand) error {
	removedServices := h.findRemovedServices(oldServices, newServices)
	if len(removedServices) == 0 {
		return nil
	}

	base.Output.Info(core.MsgStack_removing_services, fmt.Sprintf("%v", removedServices))
	return setup.DockerClient.ComposeDown(ctx, setup.Config.Project.Name, docker.StopOptions{
		Remove:        true,
		RemoveVolumes: true,
	})
}

// findRemovedServices compares old vs new service lists
func (h *UpHandler) findRemovedServices(oldServices, newServices []string) []string {
	newServiceMap := make(map[string]bool)
	for _, service := range newServices {
		newServiceMap[service] = true
	}

	var removed []string
	for _, service := range oldServices {
		if !newServiceMap[service] {
			removed = append(removed, service)
		}
	}
	return removed
}

// generateEnvFile generates .env.generated file with resolved services
func (h *UpHandler) generateEnvFile(services []string, projectName string) error {
	envContent, err := env.Generate(projectName, services)
	if err != nil {
		return fmt.Errorf("failed to generate env content: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(core.EnvGeneratedFilePath), core.PermReadWriteExec); err != nil {
		return fmt.Errorf("failed to create generated directory: %w", err)
	}
	return os.WriteFile(core.EnvGeneratedFilePath, envContent, core.PermReadWrite)
}

// runInitContainers discovers and runs initialization containers
func (h *UpHandler) runInitContainers(ctx context.Context, setup *CoreSetup, resolvedServices []string, base *base.BaseCommand) error {
	initServices := h.discoverInitServices(resolvedServices)
	if len(initServices) == 0 {
		return nil
	}

	base.Output.Info("Running initialization containers: %v", initServices)

	for _, initService := range initServices {
		if err := h.runSingleInitContainer(ctx, setup, initService, base); err != nil {
			return fmt.Errorf("failed to run init container %s: %w", initService, err)
		}
	}

	return nil
}

// discoverInitServices auto-discovers init containers based on config file patterns
func (h *UpHandler) discoverInitServices(resolvedServices []string) []string {
	var initServices []string

	configDir := filepath.Join(core.OttoStackDir, core.ServiceConfigsDir)

	// Check if directory exists
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		return initServices
	}

	// Walk through config files
	err := filepath.Walk(configDir, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil || info.IsDir() || !core.IsYAMLFile(info.Name()) {
			return nil
		}

		initService := h.processConfigFileForInit(path, info, resolvedServices)
		h.addUniqueInitService(&initServices, initService)

		return nil
	})

	if err != nil {
		return initServices
	}

	return initServices
}

// processConfigFileForInit processes a single config file and returns init service name if found
func (h *UpHandler) processConfigFileForInit(path string, info os.FileInfo, resolvedServices []string) string {
	// Extract service name from filename (remove extension)
	serviceName := core.TrimYAMLExt(info.Name())
	parts := strings.Split(serviceName, "-")

	const minPartsOfServiceName = 2
	if len(parts) < minPartsOfServiceName {
		return ""
	}

	return h.findMatchingInitServiceForFile(parts, path, resolvedServices)
}

// findMatchingInitServiceForFile finds init service that matches the config file pattern
func (h *UpHandler) findMatchingInitServiceForFile(parts []string, path string, resolvedServices []string) string {
	// Try different combinations of service name parts
	for i := len(parts) - 1; i >= 1; i-- {
		targetService := strings.Join(parts[:i], "-")

		initService := h.checkServiceMatchForFile(targetService, path, resolvedServices)
		if initService != "" {
			return initService
		}
	}
	return ""
}

// checkServiceMatchForFile checks if target service matches any resolved service
func (h *UpHandler) checkServiceMatchForFile(targetService, path string, resolvedServices []string) string {
	for _, resolved := range resolvedServices {
		if !strings.HasPrefix(resolved, targetService) && resolved != targetService {
			continue
		}

		if !h.hasValidConfiguration(path, targetService) {
			continue
		}

		return targetService + "-init"
	}
	return ""
}

// hasValidConfiguration checks if a config file has actual resources to create
func (h *UpHandler) hasValidConfiguration(configPath, targetService string) bool {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return false
	}

	// Parse YAML to check for actual configuration
	var config map[string]any
	if err := yaml.Unmarshal(data, &config); err != nil {
		return false
	}

	// Check based on service type
	switch targetService {
	case services.ServiceLocalstack:
		return h.hasValidLocalStackConfig(config)
	case services.ServicePostgres:
		return h.hasValidPostgresConfig(config)
	case services.ServiceKafka:
		return h.hasValidKafkaConfig(config)
	}

	return false
}

// isValidArray checks if an any is a non-empty array
func (h *UpHandler) isValidArray(value any) bool {
	if arr, ok := value.([]any); ok {
		return len(arr) > 0
	}
	return false
}

// hasValidLocalStackConfig checks for valid LocalStack configuration
func (h *UpHandler) hasValidLocalStackConfig(config map[string]any) bool {
	if queues, exists := config["queues"]; exists && h.isValidArray(queues) {
		return true
	}
	if topics, exists := config["topics"]; exists && h.isValidArray(topics) {
		return true
	}
	if buckets, exists := config["buckets"]; exists && h.isValidArray(buckets) {
		return true
	}
	return false
}

// hasValidPostgresConfig checks for valid PostgreSQL configuration
func (h *UpHandler) hasValidPostgresConfig(config map[string]any) bool {
	if schemas, exists := config["schemas"]; exists && h.isValidArray(schemas) {
		return true
	}
	if databases, exists := config["databases"]; exists && h.isValidArray(databases) {
		return true
	}
	return false
}

// hasValidKafkaConfig checks for valid Kafka configuration
func (h *UpHandler) hasValidKafkaConfig(config map[string]any) bool {
	if topics, exists := config["topics"]; exists && h.isValidArray(topics) {
		return true
	}
	return false
}

// parseTimeoutSeconds parses timeout string or returns default
func (h *UpHandler) parseTimeoutSeconds(timeoutStr string) int {
	if timeoutStr == "" {
		return DefaultTimeoutSeconds
	}

	parsed, err := strconv.Atoi(timeoutStr)
	if err != nil {
		return DefaultTimeoutSeconds
	}
	return parsed
}

// handleConfigChange handles configuration changes and cleanup
func (h *UpHandler) handleConfigChange(ctx context.Context, setup *CoreSetup, oldServices, newServices []string, base *base.BaseCommand, configChanged bool) {
	if !configChanged {
		return
	}

	// Clean up removed services
	if err := h.cleanupRemovedServices(ctx, setup, oldServices, newServices, base); err != nil {
		base.Output.Warning("Failed to clean up removed services: %v", err)
	}
}

// addUniqueInitService adds init service to list if not empty and not already present
func (h *UpHandler) addUniqueInitService(initServices *[]string, initService string) {
	if initService != "" && !slices.Contains(*initServices, initService) {
		*initServices = append(*initServices, initService)
	}
}

// runSingleInitContainer runs a single initialization container
func (h *UpHandler) runSingleInitContainer(ctx context.Context, setup *CoreSetup, initServiceName string, base *base.BaseCommand) error {
	targetService := strings.TrimSuffix(initServiceName, "-init")

	// Create init container configuration
	initConfig := h.createInitContainerConfig(targetService, setup)

	// Run the init container
	containerName := fmt.Sprintf("%s-%s", setup.Config.Project.Name, initServiceName)

	base.Output.Info("Starting init container: %s", containerName)

	// Run container and wait for completion
	return setup.DockerClient.RunInitContainer(ctx, containerName, initConfig)
}

// createInitContainerConfig creates configuration for init containers
func (h *UpHandler) createInitContainerConfig(targetService string, setup *CoreSetup) docker.InitContainerConfig {
	// Get current working directory for absolute path
	cwd, _ := os.Getwd()
	configPath := filepath.Join(cwd, core.OttoStackDir, core.ServiceConfigsDir)

	// Use the shell script from embedded.go
	processedScript := strings.ReplaceAll(scripts.GenericInitScript, "$$", "$")

	// Load service configuration from YAML
	manager, err := services.New()
	if err != nil {
		panic(fmt.Sprintf("Failed to create service manager: %v", err))
	}

	service, err := manager.GetService(targetService)
	if err != nil {
		panic(fmt.Sprintf("Failed to load service configuration for %s: %v", targetService, err))
	}

	// Build configuration from service YAML
	config := docker.InitContainerConfig{
		Image:   services.ImageLocalstack, // Default, will be overridden based on service
		Command: []string{"sh", "-c", processedScript},
		Environment: map[string]string{
			services.InitServiceName: targetService,
			services.InitConfigDir:   "/config",
		},
		Volumes: []string{
			fmt.Sprintf("%s:/config", configPath),
		},
		WorkingDir: "/",
		Networks:   []string{setup.Config.Project.Name + services.NetworkNameSuffix},
	}

	// Extract environment variables from service YAML
	maps.Copy(config.Environment, service.Environment)

	// Set service endpoint URL based on service configuration
	if service.Service.Connection != nil && service.Service.Connection.DefaultPort > 0 {
		config.Environment[services.InitServiceEndpointURL] = fmt.Sprintf("http://%s:%d", targetService, service.Service.Connection.DefaultPort)
	}

	// Customize based on service type using service YAML data
	switch targetService {
	case services.ServiceLocalstack:
		config.Environment[services.InitServiceEndpointURL] = fmt.Sprintf("http://localhost:%d", services.PortLocalstack)
	case services.ServicePostgres:
		// Use postgres image from service YAML
		if service.Container.Image != "" {
			config.Image = service.Container.Image
		}

		// Use PGHOST key from postgres service definition
		config.Environment[services.EnvPostgresPGHOST] = targetService

		// Build connection URL from service data using postgres service format
		user := config.Environment[services.EnvPostgresPOSTGRES_USER]
		password := config.Environment[services.EnvPostgresPOSTGRES_PASSWORD]
		database := config.Environment[services.EnvPostgresPOSTGRES_DB]
		port := fmt.Sprintf("%d", services.PortPostgres)
		if service.Service.Connection != nil && service.Service.Connection.DefaultPort > 0 {
			port = fmt.Sprintf("%d", service.Service.Connection.DefaultPort)
		}
		config.Environment[services.InitServiceEndpointURL] = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", user, password, targetService, port, database)
	case services.ServiceKafka:
		// Use kafka image from service YAML
		if service.Container.Image != "" {
			config.Image = service.Container.Image
		}

		// Set kafka endpoint
		port := fmt.Sprintf("%d", services.PortKafkaBroker)
		if service.Service.Connection != nil && service.Service.Connection.DefaultPort > 0 {
			port = fmt.Sprintf("%d", service.Service.Connection.DefaultPort)
		}
		config.Environment[services.InitServiceEndpointURL] = fmt.Sprintf("%s:%s", targetService, port)
	}

	return config
}
