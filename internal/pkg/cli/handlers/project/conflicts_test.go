//go:build unit

package project

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/otto-nation/otto-stack/test/fixtures"
	"github.com/stretchr/testify/assert"
)

func TestConflictsHandler_ValidateArgs(t *testing.T) {
	handler := &ConflictsHandler{}

	err := handler.ValidateArgs([]string{})
	assert.NoError(t, err)

	err = handler.ValidateArgs([]string{services.ServicePostgres, services.ServiceRedis})
	assert.NoError(t, err)
}

func TestConflictsHandler_GetRequiredFlags(t *testing.T) {
	handler := &ConflictsHandler{}
	flags := handler.GetRequiredFlags()
	assert.Empty(t, flags)
}

func TestConflictsHandler_ParsePort(t *testing.T) {
	handler := &ConflictsHandler{}

	assert.Equal(t, 8080, handler.parsePort("8080"))
	assert.Equal(t, 0, handler.parsePort("invalid"))
	assert.Equal(t, 0, handler.parsePort(""))
}

func TestConflictsHandler_PairKey(t *testing.T) {
	assert.Equal(t, pairKey("mysql", "postgres"), pairKey("postgres", "mysql"))
	assert.Equal(t, "mysql/postgres", pairKey("postgres", "mysql"))
	assert.Equal(t, "mysql/postgres", pairKey("mysql", "postgres"))
}

func TestConflictsHandler_HasExplicitConflict(t *testing.T) {
	handler := &ConflictsHandler{}

	postgres := fixtures.NewServiceConfig("postgres").WithConflicts("mysql").Build()
	mysql := fixtures.NewServiceConfig("mysql").WithConflicts("postgres").Build()
	redis := fixtures.NewServiceConfig("redis").Build()

	assert.True(t, handler.hasExplicitConflict(postgres, mysql))
	assert.True(t, handler.hasExplicitConflict(mysql, postgres))
	assert.False(t, handler.hasExplicitConflict(postgres, redis))
	assert.False(t, handler.hasExplicitConflict(redis, postgres))
}

func TestConflictsHandler_ProvidesOverlap(t *testing.T) {
	handler := &ConflictsHandler{}

	postgres := fixtures.NewServiceConfig("postgres").WithProvides("database", "sql").Build()
	mysql := fixtures.NewServiceConfig("mysql").WithProvides("database", "sql").Build()
	redis := fixtures.NewServiceConfig("redis").WithProvides("cache").Build()

	overlaps := handler.providesOverlap(postgres, mysql)
	assert.ElementsMatch(t, []string{"database", "sql"}, overlaps)

	overlaps = handler.providesOverlap(postgres, redis)
	assert.Empty(t, overlaps)

	overlaps = handler.providesOverlap(redis, mysql)
	assert.Empty(t, overlaps)
}

func TestConflictsHandler_DetectSemanticConflicts(t *testing.T) {
	handler := &ConflictsHandler{}

	t.Run("explicit conflict detected once", func(t *testing.T) {
		postgres := fixtures.NewServiceConfig("postgres").WithConflicts("mysql").Build()
		mysql := fixtures.NewServiceConfig("mysql").WithConflicts("postgres").Build()

		conflicts := handler.detectSemanticConflicts([]types.ServiceConfig{postgres, mysql})
		assert.Len(t, conflicts, 1)
		assert.Empty(t, conflicts[0].capability, "explicit conflict should not set capability")
	})

	t.Run("provides overlap detected", func(t *testing.T) {
		svcA := fixtures.NewServiceConfig("mariadb").WithProvides("database").Build()
		svcB := fixtures.NewServiceConfig("cockroach").WithProvides("database").Build()

		conflicts := handler.detectSemanticConflicts([]types.ServiceConfig{svcA, svcB})
		assert.Len(t, conflicts, 1)
		assert.Equal(t, "database", conflicts[0].capability)
	})

	t.Run("explicit conflict takes precedence over provides overlap", func(t *testing.T) {
		postgres := fixtures.NewServiceConfig("postgres").
			WithConflicts("mysql").
			WithProvides("database").
			Build()
		mysql := fixtures.NewServiceConfig("mysql").
			WithConflicts("postgres").
			WithProvides("database").
			Build()

		// Reported once as explicit (not twice)
		conflicts := handler.detectSemanticConflicts([]types.ServiceConfig{postgres, mysql})
		assert.Len(t, conflicts, 1)
		assert.Empty(t, conflicts[0].capability)
	})

	t.Run("no conflicts when capabilities differ", func(t *testing.T) {
		postgres := fixtures.NewServiceConfig("postgres").WithProvides("database").Build()
		redis := fixtures.NewServiceConfig("redis").WithProvides("cache").Build()

		conflicts := handler.detectSemanticConflicts([]types.ServiceConfig{postgres, redis})
		assert.Empty(t, conflicts)
	})

	t.Run("empty input", func(t *testing.T) {
		conflicts := handler.detectSemanticConflicts(nil)
		assert.Empty(t, conflicts)
	})
}

func TestDepsHandler_ValidateArgs(t *testing.T) {
	handler := &DepsHandler{}

	err := handler.ValidateArgs([]string{})
	assert.NoError(t, err)

	err = handler.ValidateArgs([]string{services.ServicePostgres})
	assert.NoError(t, err)
}

func TestDepsHandler_GetRequiredFlags(t *testing.T) {
	handler := &DepsHandler{}
	flags := handler.GetRequiredFlags()
	assert.Empty(t, flags)
}
