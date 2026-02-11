package project

import (
	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
	svc "github.com/otto-nation/otto-stack/internal/pkg/services"
)

// InteractiveProcessor handles interactive mode processing
type InteractiveProcessor struct {
	handler *InitHandler
}

// Process runs interactive prompts to gather project details
func (p *InteractiveProcessor) Process(flags *core.InitFlags, base *base.BaseCommand) (clicontext.Context, error) {
	projectName, err := p.handler.promptManager.PromptForProjectDetails()
	if err != nil {
		return clicontext.Context{}, err
	}

	result, err := p.handler.serviceSelectionManager.RunWorkflow(p.handler, projectName, base)
	if err != nil {
		return clicontext.Context{}, err
	}

	// Extract original service names from configs for interactive mode
	originalServiceNames := svc.ExtractServiceNames(result.ServiceConfigs)

	// For interactive mode, default to sharing enabled
	sharingSpec := &clicontext.SharingSpec{
		Enabled: true,
	}

	ctx := clicontext.NewBuilder().
		WithProject(projectName, "").
		WithServices(originalServiceNames, result.ServiceConfigs).
		WithValidation(result.Validation).
		WithAdvanced(result.Advanced).
		WithRuntimeFlags(flags, true).
		WithSharing(sharingSpec).
		Build()

	return ctx, nil
}
