package project

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/ci"
	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/common"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/logger"
	"github.com/otto-nation/otto-stack/internal/pkg/messages"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
)

// InitHandler handles the init command
type InitHandler struct {
	forceOverwrite    bool
	base              *base.BaseCommand
	promptManager     *PromptManager
	validationManager *ValidationManager
	projectManager    *ProjectManager
}

// NewInitHandler creates a new InitHandler
func NewInitHandler() *InitHandler {
	handler := &InitHandler{
		validationManager: NewValidationManager(),
		projectManager:    NewProjectManager(),
	}
	handler.promptManager = NewPromptManager(handler)
	return handler
}

// ValidateProjectName implements ProjectValidator interface
func (h *InitHandler) ValidateProjectName(name string) error {
	return h.validateProjectName(name)
}

// Handle executes the init command
func (h *InitHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
	logger.Info(logger.LogMsgStartingOperation, logger.LogFieldOperation, logger.OperationInit)
	defer h.handlePanic()

	if err := h.validateInitFlags(cmd); err != nil {
		return pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldFlags, messages.ValidationFailed, err)
	}

	ciFlags := ci.GetFlags(cmd)
	initFlags, _ := core.ParseInitFlags(cmd)

	// Check if already initialized before any other processing
	if err := h.checkAlreadyInitialized(initFlags.Force); err != nil {
		return err
	}

	h.logDebugInfo(base, cmd, initFlags)

	if err := h.setDefaultProjectName(initFlags, base, cmd); err != nil {
		return pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldProjectName, messages.ValidationFailedSetProjectName, err)
	}

	h.forceOverwrite = initFlags.Force
	h.base = base

	projectCtx, err := h.gather(ciFlags, initFlags, base)
	if err != nil {
		return pkgerrors.NewSystemError(pkgerrors.ErrCodeOperationFail, messages.ValidationFailedProcessInitMode, err)
	}

	// Empty project name means the user cancelled the interactive session.
	if projectCtx.Project.Name == "" {
		return nil
	}

	if err := h.executeInit(projectCtx, base); err != nil {
		return pkgerrors.NewSystemError(pkgerrors.ErrCodeOperationFail, messages.ValidationFailedExecuteInit, err)
	}

	h.displaySuccessMessage(projectCtx.Project.Name, base)

	if projectCtx.Advanced != nil && projectCtx.Advanced.AutoStart {
		base.Output.Info("%s", messages.InfoAutoStarting)
		// Non-fatal: init succeeded; auto-start is best-effort.
		if err := h.autoStartServices(ctx, projectCtx); err != nil {
			base.Output.Warning(messages.WarningsAutoStartFailed, err)
		}
	}

	return nil
}

// gather collects all user input — either from flags (non-interactive) or prompts (interactive).
func (h *InitHandler) gather(ciFlags ci.Flags, initFlags *core.InitFlags, base *base.BaseCommand) (clicontext.Context, error) {
	if ciFlags.NonInteractive {
		return h.gatherNonInteractive(initFlags)
	}
	return h.gatherInteractive(initFlags, base)
}

// gatherInteractive runs the interactive prompt flow: project name → service selection → advanced
// options → confirmation. Returns an empty context (Project.Name == "") if the user cancels.
func (h *InitHandler) gatherInteractive(initFlags *core.InitFlags, base *base.BaseCommand) (clicontext.Context, error) {
	base.Output.Header("%s", messages.ProcessInitializing)
	logger.Info(logger.LogMsgProjectAction, logger.LogFieldAction, core.CommandInit, logger.LogFieldProject, "current_directory")

	projectName, err := h.promptManager.PromptForProjectDetails()
	if err != nil {
		return clicontext.Context{}, err
	}

	for {
		result, err := h.runSelectionCycle(projectName, base)
		if err != nil {
			return clicontext.Context{}, err
		}

		switch result.action {
		case core.ActionProceed:
			serviceNames := services.ExtractServiceNames(result.serviceConfigs)
			return clicontext.NewBuilder().
				WithProject(projectName, "").
				WithServices(serviceNames, result.serviceConfigs).
				WithValidation(result.validation).
				WithAdvanced(map[string]bool{}).
				WithAdvancedSpec(result.advanced).
				WithRuntimeFlags(initFlags, true).
				WithSharing(result.sharing).
				Build(), nil
		case core.ActionBack:
			base.Output.Info("%s", messages.InitGoingBack)
		default:
			base.Output.Info("%s", messages.InitCancelled)
			return clicontext.Context{}, nil
		}
	}
}

// selectionResult holds the outcome of one interactive selection cycle.
type selectionResult struct {
	serviceConfigs []types.ServiceConfig
	validation     map[string]bool
	sharing        *clicontext.SharingSpec
	advanced       *clicontext.AdvancedSpec
	action         string
}

// runSelectionCycle runs one pass of: service selection → advanced options → validation → confirm.
func (h *InitHandler) runSelectionCycle(projectName string, base *base.BaseCommand) (*selectionResult, error) {
	selectedConfigs, err := h.promptManager.PromptForServiceConfigs()
	if err != nil {
		return nil, err
	}

	serviceNames := services.ExtractServiceNames(selectedConfigs)
	serviceConfigs, err := services.ResolveUpServices(serviceNames, nil)
	if err != nil {
		return nil, err
	}

	manager, err := services.New()
	if err != nil {
		return nil, err
	}
	if err := services.NewValidationService(manager).ValidateResolvedServices(serviceConfigs); err != nil {
		return nil, err
	}

	validation, sharing, advanced, err := h.promptManager.PromptForAdvancedOptions(serviceConfigs)
	if err != nil {
		return nil, err
	}

	if err := h.validationManager.RunValidations(validation, h, serviceConfigs, base); err != nil {
		return nil, err
	}

	action, err := h.promptManager.ConfirmInitialization(projectName, serviceNames, validation, base)
	if err != nil {
		return nil, err
	}

	// Return original user selection — dependencies are resolved at runtime, not stored in config.
	return &selectionResult{
		serviceConfigs: selectedConfigs,
		validation:     validation,
		sharing:        sharing,
		advanced:       advanced,
		action:         action,
	}, nil
}

// gatherNonInteractive builds the context directly from command-line flags.
func (h *InitHandler) gatherNonInteractive(initFlags *core.InitFlags) (clicontext.Context, error) {
	if initFlags.Services == "" {
		return clicontext.Context{}, pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldServiceName, messages.ValidationServicesRequiredNonInteractive, nil)
	}

	if initFlags.ProjectName == "" {
		return clicontext.Context{}, pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldProjectName, messages.ValidationProjectNameRequiredNonInteractive, nil)
	}

	serviceNames := parseServiceNames(initFlags.Services)
	serviceConfigs, err := services.ResolveUpServices(serviceNames, nil)
	if err != nil {
		return clicontext.Context{}, err
	}

	manager, err := services.New()
	if err != nil {
		return clicontext.Context{}, err
	}
	if err := services.NewValidationService(manager).ValidateResolvedServices(serviceConfigs); err != nil {
		return clicontext.Context{}, err
	}

	sharingConfig, err := h.buildSharingConfig(initFlags, serviceConfigs)
	if err != nil {
		return clicontext.Context{}, err
	}

	return clicontext.NewBuilder().
		WithProject(initFlags.ProjectName, "").
		WithServices(serviceNames, serviceConfigs).
		WithValidation(defaultValidation()).
		WithAdvanced(map[string]bool{}).
		WithAdvancedSpec(&clicontext.AdvancedSpec{AutoStart: initFlags.AutoStart}).
		WithRuntimeFlags(initFlags, false).
		WithSharing(sharingConfig).
		Build(), nil
}

func parseServiceNames(s string) []string {
	parts := strings.Split(s, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func defaultValidation() map[string]bool {
	validation := make(map[string]bool)
	for key := range ValidationRegistry {
		validation[key] = true
	}
	return validation
}

func (h *InitHandler) handlePanic() {
	if r := recover(); r != nil {
		logger.Error(logger.LogMsgOperationFailed, logger.LogFieldOperation, logger.OperationInit, logger.LogFieldError, fmt.Errorf("panic: %v", r))
		panic(r)
	}
}

func (h *InitHandler) logDebugInfo(base *base.BaseCommand, cmd *cobra.Command, initFlags *core.InitFlags) {
	if base.GetVerbose(cmd) {
		logger.Debug("Initializing project", "projectName", initFlags.ProjectName, "force", initFlags.Force)
	}
}

func (h *InitHandler) setDefaultProjectName(initFlags *core.InitFlags, base *base.BaseCommand, cmd *cobra.Command) error {
	if initFlags.ProjectName != "" {
		return nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldServiceName, messages.ErrorsCurrentDirectoryFailed, err)
	}

	initFlags.ProjectName = filepath.Base(cwd)
	if base.GetVerbose(cmd) {
		logger.Debug("Using current directory as project name", "projectName", initFlags.ProjectName)
	}
	return nil
}

func (h *InitHandler) executeInit(projectCtx clicontext.Context, base *base.BaseCommand) error {
	_, dirExistedErr := os.Stat(core.OttoStackDir)

	err := h.projectManager.CreateProjectStructure(projectCtx, base)
	if err != nil && os.IsNotExist(dirExistedErr) {
		// Roll back the partially-created directory tree so the user can retry cleanly.
		_ = os.RemoveAll(core.OttoStackDir)
	}
	return err
}

func (h *InitHandler) validateInitFlags(cmd *cobra.Command) error {
	initFlags, err := core.ParseInitFlags(cmd)
	if err != nil {
		logger.Error(logger.LogMsgOperationFailed, logger.LogFieldOperation, logger.OperationInit, logger.LogFieldError, err)
		return err
	}

	return h.validateNonInteractiveMode(cmd, initFlags)
}

func (h *InitHandler) validateNonInteractiveMode(cmd *cobra.Command, initFlags *core.InitFlags) error {
	ciFlags := ci.GetFlags(cmd)
	if ciFlags.NonInteractive && initFlags.Services == "" {
		return pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldServiceName, messages.ValidationNoServicesSelected, nil)
	}
	return nil
}

func (h *InitHandler) displaySuccessMessage(_ string, base *base.BaseCommand) {
	base.Output.Success("%s", messages.SuccessInit)
	base.Output.Info("%s", messages.InitNextSteps)
	h.displayNextSteps(base)
}

func (h *InitHandler) displayNextSteps(base *base.BaseCommand) {
	base.Output.Info(messages.InitStepReviewConfig, core.OttoStackDir, core.ConfigFileName)
	base.Output.Info(messages.InitStepStartStack, core.AppName)
	base.Output.Info(messages.InitStepCheckStatus, core.AppName)
}

// ValidateArgs validates the command arguments
func (h *InitHandler) ValidateArgs(args []string) error {
	return nil
}

// GetRequiredFlags returns required flags for this command
func (h *InitHandler) GetRequiredFlags() []string {
	return []string{}
}

// validateProjectName validates a project name
func (h *InitHandler) validateProjectName(name string) error {
	return h.validationManager.ValidateProjectName(name)
}

// buildSharingConfig creates sharing configuration from init flags
func (h *InitHandler) buildSharingConfig(initFlags *core.InitFlags, serviceConfigs []types.ServiceConfig) (*clicontext.SharingSpec, error) {
	if initFlags.NoSharedContainers {
		return &clicontext.SharingSpec{Enabled: false}, nil
	}

	sharingSpec := &clicontext.SharingSpec{
		Enabled:  true,
		Services: make(map[string]bool),
	}

	if initFlags.SharedServices == "" {
		return sharingSpec, nil
	}

	for svc := range strings.SplitSeq(initFlags.SharedServices, ",") {
		svc = strings.TrimSpace(svc)
		if svc == "" {
			continue
		}

		if err := h.validateServiceShareable(svc, serviceConfigs); err != nil {
			return nil, err
		}
		sharingSpec.Services[svc] = true
	}

	return sharingSpec, nil
}

// validateServiceShareable checks if a service can be shared
func (h *InitHandler) validateServiceShareable(serviceName string, serviceConfigs []types.ServiceConfig) error {
	for _, cfg := range serviceConfigs {
		if cfg.Name != serviceName {
			continue
		}
		if cfg.Shareable {
			return nil
		}
		if h.forceOverwrite {
			h.base.Output.Warning(messages.ValidationServiceNotShareableWarning, serviceName)
			return nil
		}
		return pkgerrors.NewValidationErrorf(
			pkgerrors.ErrCodeInvalid,
			"shared-services",
			messages.ValidationServiceNotShareable,
			serviceName,
		)
	}
	return nil
}

// autoStartServices starts the project services after a successful init.
// It loads the freshly-written config so the service manager sees the correct project name.
// Shared services are excluded — they run under their own compose project and must be
// started separately via `otto-stack up <service>` in shared mode.
func (h *InitHandler) autoStartServices(ctx context.Context, projectCtx clicontext.Context) error {
	svc, err := common.NewServiceManager(false)
	if err != nil {
		return err
	}

	var sharingEnabled bool
	var sharingServices map[string]bool
	if projectCtx.Sharing != nil {
		sharingEnabled = projectCtx.Sharing.Enabled
		sharingServices = projectCtx.Sharing.Services
	}
	projectConfigs := FilterProjectServices(projectCtx.Services.Configs, sharingEnabled, sharingServices)

	return svc.Start(ctx, services.StartRequest{
		Project:        projectCtx.Project.Name,
		ServiceConfigs: projectConfigs,
		Detach:         true,
	})
}

func (h *InitHandler) checkAlreadyInitialized(force bool) error {
	if force {
		return nil
	}

	if _, err := os.Stat(core.OttoStackDir); err == nil {
		return pkgerrors.NewSystemError(pkgerrors.ErrCodeAlreadyExists, messages.MiddlewareProjectAlreadyInitialized, nil)
	}

	return nil
}
