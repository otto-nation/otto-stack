package constants

import (
	"fmt"

	"github.com/otto-nation/otto-stack/internal/pkg/ui"
)

// Message levels
const (
	LevelInfo    = "info"
	LevelSuccess = "success"
	LevelWarning = "warning"
	LevelError   = "error"
	LevelHeader  = "header"
)

// Message represents a structured message with level and content
type Message struct {
	Level   string
	Content string
}

// Init process messages
var (
	MsgGoingBack     = Message{LevelInfo, "Going back to service selection..."}
	MsgInitCancelled = Message{LevelInfo, "Initialization cancelled"}
	MsgInitSummary   = Message{LevelInfo, "Initialization Summary:"}
	MsgNextSteps     = Message{LevelInfo, "Next steps:"}

	// Summary templates (use with fmt.Sprintf)
	MsgProject     = Message{LevelInfo, "  Project: %s"}
	MsgEnvironment = Message{LevelInfo, "  Environment: %s"}
	MsgServices    = Message{LevelInfo, "  Services: %s"}
	MsgValidation  = Message{LevelInfo, "  Validation: %s"}
	MsgAdvanced    = Message{LevelInfo, "  Advanced: %s"}

	// Next steps templates
	MsgStep1 = Message{LevelInfo, "  1. Review the configuration in %s/%s"}
	MsgStep2 = Message{LevelInfo, "  2. Start your stack with: %s"}
	MsgStep3 = Message{LevelInfo, "  3. Check status with: %s"}

	// Service selection
	MsgSelectServices = Message{LevelInfo, "Select %s services:"}
)

// File operation messages
var (
	MsgCreatedFile       = Message{LevelSuccess, "Created %s"}
	MsgUpdatedGitignore  = Message{LevelSuccess, "Updated .gitignore with otto-stack entries"}
	MsgGitignoreExists   = Message{LevelInfo, ".gitignore already contains otto-stack entries"}
	MsgFailedGitignore   = Message{LevelWarning, "Failed to update .gitignore: %v"}
	MsgFailedReadme      = Message{LevelWarning, "Failed to create README: %v"}
	MsgGeneratingEnv     = Message{LevelInfo, "Generating environment file..."}
	MsgGeneratingCompose = Message{LevelInfo, "Generating docker-compose files..."}
)

// Service messages
var (
	MsgNoServicesConfigured = Message{LevelInfo, "- No services configured"}
	MsgServiceListItem      = Message{LevelInfo, "- %s"}
	MsgNoDescription        = Message{LevelInfo, "No description available"}
)

// Command error messages
var (
	MsgRequiresServiceName       = Message{LevelError, "requires service name"}
	MsgRequiresServiceAndCommand = Message{LevelError, "requires service name and command"}
	MsgFailedResolveServices     = Message{LevelError, "failed to resolve services: %w"}
	MsgFailedGetConfirmation     = Message{LevelError, "failed to get confirmation: %w"}
	MsgCleanupFailed             = Message{LevelError, "cleanup failed: %w"}
	MsgUnsupportedServiceType    = Message{LevelError, "unsupported service type: %s. Supported: postgres, mysql, redis, mongodb"}
)

// Cleanup messages
var (
	MsgCleanupDryRun              = Message{LevelInfo, "Dry run - showing what would be cleaned up:"}
	MsgCleanupUnusedVolumes       = Message{LevelInfo, "  - Unused volumes"}
	MsgCleanupUnusedImages        = Message{LevelInfo, "  - Unused images"}
	MsgCleanupUnusedNetworks      = Message{LevelInfo, "  - Unused networks"}
	MsgCleanupStoppedContainers   = Message{LevelInfo, "  - Stopped containers"}
	MsgCleanupWarning             = Message{LevelWarning, "This will remove Docker resources"}
	MsgCleanupConfirm             = Message{LevelInfo, "Continue with cleanup?"}
	MsgCleanupCancelled           = Message{LevelInfo, "Cleanup cancelled"}
	MsgCleanupSuccess             = Message{LevelSuccess, "Cleanup completed successfully"}
	MsgRemovingContainers         = Message{LevelInfo, "Removing stopped containers..."}
	MsgCleanupOperationsCompleted = Message{LevelSuccess, "Cleanup operations completed"}
	MsgFailedRemoveContainers     = Message{LevelWarning, "Failed to remove some containers: %v"}
)

// Doctor messages
var (
	MsgHealthCheckHeader      = Message{LevelHeader, "🩺 %s Health Check"}
	MsgAllChecksPassedHealthy = Message{LevelSuccess, "All checks passed! Your %s is healthy."}
	MsgSomeIssuesFound        = Message{LevelError, "Some issues found. Please address them above."}
	MsgHealthCheckFailed      = Message{LevelError, "health check failed"}
	MsgCheckingDocker         = Message{LevelInfo, "Checking Docker installation..."}
	MsgDockerNotFound         = Message{LevelError, "Docker not found"}
	MsgInstallDocker          = Message{LevelInfo, "Install Docker: %s"}
)

// Services messages
var (
	MsgAvailableServices     = Message{LevelHeader, "Available Services"}
	MsgNoServicesAvailable   = Message{LevelInfo, "No services available"}
	MsgFailedLoadServices    = Message{LevelError, "failed to load services: %w"}
	MsgFailedCreateFormatter = Message{LevelError, "failed to create formatter: %w"}
)

// Completion messages
var (
	MsgCompletionRequiresOneArg = Message{LevelError, "completion requires exactly one argument (%s)"}
	MsgUnsupportedShell         = Message{LevelError, "unsupported shell: %s (supported: %v)"}
	MsgUnsupportedShellSimple   = Message{LevelError, "unsupported shell: %s"}
)

// Validation messages
var (
	MsgValidationFailed          = Message{LevelError, "validation failed: %w"}
	MsgDirectoryValidationFailed = Message{LevelError, "directory validation failed: %w"}
)

// Helper function to send message with appropriate UI method
func SendMessage(msg Message, args ...interface{}) {
	content := msg.Content
	if len(args) > 0 {
		content = fmt.Sprintf(msg.Content, args...)
	}

	switch msg.Level {
	case LevelInfo:
		ui.Info("%s", content)
	case LevelSuccess:
		ui.Success("%s", content)
	case LevelWarning:
		ui.Warning("%s", content)
	case LevelError:
		ui.Error("%s", content)
	case LevelHeader:
		ui.Header("%s", content)
	default:
		ui.Info("%s", content)
	}
}
