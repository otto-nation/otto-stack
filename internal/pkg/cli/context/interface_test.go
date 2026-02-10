//go:build unit

package context

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProjectMode_ImplementsExecutionMode(t *testing.T) {
	var _ ExecutionMode = (*ProjectMode)(nil)
}

func TestSharedMode_ImplementsExecutionMode(t *testing.T) {
	var _ ExecutionMode = (*SharedMode)(nil)
}

func TestProjectMode_SharedRoot(t *testing.T) {
	ctx := &ProjectMode{
		Shared: &SharedInfo{Root: "/test/shared"},
	}
	assert.Equal(t, "/test/shared", ctx.SharedRoot())
}

func TestSharedMode_SharedRoot(t *testing.T) {
	ctx := &SharedMode{
		Shared: &SharedInfo{Root: "/test/shared"},
	}
	assert.Equal(t, "/test/shared", ctx.SharedRoot())
}

func TestProjectMode_HasProjectInfo(t *testing.T) {
	project := &ProjectInfo{
		Root:       "/test/project",
		ConfigDir:  "/test/project/.otto-stack",
		ConfigFile: "/test/project/.otto-stack/config.yaml",
	}
	ctx := &ProjectMode{
		Project: project,
		Shared:  &SharedInfo{Root: "/test/shared"},
	}
	assert.NotNil(t, ctx.Project)
	assert.Equal(t, "/test/project", ctx.Project.Root)
}

func TestSharedMode_NoProjectInfo(t *testing.T) {
	ctx := &SharedMode{
		Shared: &SharedInfo{Root: "/test/shared"},
	}
	// SharedMode doesn't have Project field - compile-time safety
	assert.NotNil(t, ctx.Shared)
}
