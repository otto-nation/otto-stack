//go:build unit

package context

import (
	"path/filepath"
	"testing"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestNewBuilder(t *testing.T) {
	builder := NewBuilder()

	assert.NotNil(t, builder)
	assert.NotNil(t, builder.ctx.Options.Validation)
	assert.NotNil(t, builder.ctx.Options.Advanced)
}

func TestBuilder_FluentAPI(t *testing.T) {
	flags := &core.InitFlags{Force: true}

	ctx := NewBuilder().
		WithProject("test-project", "/test/path").
		WithServices([]string{"postgres", "redis"}, []types.ServiceConfig{
			{Name: "postgres"},
			{Name: "redis"},
		}).
		WithValidation(map[string]bool{"docker": true}).
		WithAdvanced(map[string]bool{"networking": true}).
		WithRuntimeFlags(flags, false).
		Build()

	assert.Equal(t, "test-project", ctx.Project.Name)
	assert.Equal(t, "/test/path", ctx.Project.Path)
	assert.Equal(t, []string{"postgres", "redis"}, ctx.Services.Names)
	assert.Len(t, ctx.Services.Configs, 2)
	assert.True(t, ctx.Options.Validation["docker"])
	assert.True(t, ctx.Options.Advanced["networking"])
	assert.True(t, ctx.Runtime.Force)
	assert.False(t, ctx.Runtime.Interactive)
}

func TestBuilder_EmptyContext(t *testing.T) {
	ctx := NewBuilder().Build()

	assert.Empty(t, ctx.Project.Name)
	assert.Empty(t, ctx.Services.Names)
	assert.Empty(t, ctx.Options.Validation)
	assert.Empty(t, ctx.Options.Advanced)
	assert.False(t, ctx.Runtime.Force)
}

func TestBuilder_WithSharing(t *testing.T) {
	sharing := &SharingSpec{
		Enabled:  true,
		Services: map[string]bool{"redis": true, "postgres": true},
	}

	ctx := NewBuilder().
		WithSharing(sharing).
		Build()

	assert.NotNil(t, ctx.Sharing)
	assert.True(t, ctx.Sharing.Enabled)
	assert.True(t, ctx.Sharing.Services["redis"])
	assert.True(t, ctx.Sharing.Services["postgres"])
}

func TestNewProjectInfo(t *testing.T) {
	configDir := filepath.Join(string(filepath.Separator), "test", ".otto-stack")
	info := NewProjectInfo(configDir)

	expectedRoot := filepath.Join(string(filepath.Separator), "test")
	assert.Equal(t, expectedRoot, info.Root)
	assert.Equal(t, configDir, info.ConfigDir)
	assert.Contains(t, info.ConfigFile, "config.yaml")
}
