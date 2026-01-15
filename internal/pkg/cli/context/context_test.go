//go:build unit

package context

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/services"
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
	ctx := NewBuilder().
		WithProject("test-project", "/test/path").
		WithServices([]string{services.ServicePostgres, services.ServiceRedis}, []types.ServiceConfig{
			{Name: services.ServicePostgres},
			{Name: services.ServiceRedis},
		}).
		WithValidation(map[string]bool{"docker": true}).
		WithAdvanced(map[string]bool{"networking": true}).
		WithRuntime(true, false, false).
		Build()

	assert.Equal(t, "test-project", ctx.Project.Name)
	assert.Equal(t, "/test/path", ctx.Project.Path)
	assert.Equal(t, []string{services.ServicePostgres, services.ServiceRedis}, ctx.Services.Names)
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
