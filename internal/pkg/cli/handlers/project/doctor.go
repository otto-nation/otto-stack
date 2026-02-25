package project

import (
	"context"
	"encoding/json"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/logger"
	"github.com/otto-nation/otto-stack/internal/pkg/messages"
	"github.com/spf13/cobra"
)

type DoctorHandler struct {
	healthCheckManager *HealthCheckManager
}

func NewDoctorHandler() *DoctorHandler {
	return &DoctorHandler{
		healthCheckManager: NewHealthCheckManager(),
	}
}

func (h *DoctorHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
	flags, err := core.ParseDoctorFlags(cmd)
	if err != nil {
		return pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldFlags, messages.ValidationFailedParseFlags, err)
	}

	logger.Info(logger.LogMsgProjectAction, logger.LogFieldAction, core.CommandDoctor, logger.LogFieldProject, "health_check")

	if flags.Format == "json" {
		results := h.healthCheckManager.collectResults(ctx)
		allPassed := true
		for _, r := range results {
			if !r.Passed {
				allPassed = false
				break
			}
		}
		_ = json.NewEncoder(base.Output.Writer()).Encode(results)
		if !allPassed {
			return pkgerrors.ErrSilentExit
		}
		return nil
	}

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
	return pkgerrors.NewSystemError(pkgerrors.ErrCodeInvalid, messages.DoctorHealthCheckFailed, nil)
}
