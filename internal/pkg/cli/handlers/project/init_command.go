package project

import (
	"context"

	"github.com/otto-nation/otto-stack/internal/pkg/base"
	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
)

// InitCommand implements the Command interface for project initialization
type InitCommand struct {
	projectManager *ProjectManager
}

// NewInitCommand creates a new init command
func NewInitCommand(projectManager *ProjectManager) *InitCommand {
	return &InitCommand{
		projectManager: projectManager,
	}
}

// Execute runs the project initialization
func (c *InitCommand) Execute(ctx context.Context, cliCtx clicontext.Context, base *base.BaseCommand) error {
	return c.projectManager.CreateProjectStructure(cliCtx, base)
}
