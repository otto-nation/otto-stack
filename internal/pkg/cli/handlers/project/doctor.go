package project

import (
	"context"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/common"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/logger"
	"github.com/otto-nation/otto-stack/internal/pkg/messages"
	"github.com/otto-nation/otto-stack/internal/pkg/ui"
	"github.com/spf13/cobra"
)

type DoctorHandler struct {
	common.BaseHandler
	output             *ui.Output
	healthCheckManager *HealthCheckManager
}

func NewDoctorHandler() *DoctorHandler {
	return &DoctorHandler{
		output:             ui.NewOutput(),
		healthCheckManager: NewHealthCheckManager(),
	}
}

func (h *DoctorHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
	logger.Info(logger.LogMsgProjectAction, logger.LogFieldAction, core.CommandDoctor, logger.LogFieldProject, "health_check")

	verbose, _ := cmd.Flags().GetBool("verbose")
	logger.Debug("Running doctor command", "verbose", verbose)

	base.Output.Header(messages.DoctorHealthCheckHeader, core.AppName)
	logger.Info("Starting health checks")

	allGood := h.healthCheckManager.RunAllChecks(ctx, base)

	if allGood {
		base.Output.Success(messages.SuccessAllChecksPassed, core.AppName)
		logger.Info("All health checks passed")
		return nil
	}

	base.Output.Error(messages.DoctorSomeIssues)
	logger.Error("Health checks failed")
	return pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, "health", messages.DoctorHealthCheckFailed, nil)
}
