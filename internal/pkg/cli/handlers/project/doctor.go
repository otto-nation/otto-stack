package project

import (
	"context"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	basehandler "github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/base"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/logger"
	"github.com/otto-nation/otto-stack/internal/pkg/ui"
	"github.com/spf13/cobra"
)

type DoctorHandler struct {
	basehandler.BaseHandler
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

	base.Output.Header(core.MsgDoctor_health_check_header, core.AppName)
	logger.Info("Starting health checks")

	allGood := h.healthCheckManager.RunAllChecks(base)

	if allGood {
		base.Output.Success(core.MsgSuccess_all_checks_passed, core.AppName)
		logger.Info("All health checks passed")
		return nil
	}

	base.Output.Error(core.MsgDoctor_some_issues)
	logger.Error("Health checks failed")
	return pkgerrors.NewValidationError(FieldHealth, core.MsgDoctor_health_check_failed, nil)
}
