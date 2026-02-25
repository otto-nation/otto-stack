//go:build unit

package project

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/display"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/otto-nation/otto-stack/test/fixtures"
	"github.com/stretchr/testify/assert"
)

func TestDepsHandler_JoinOrDash(t *testing.T) {
	assert.Equal(t, depsDash, joinOrDash([]string{}))
	assert.Equal(t, depsDash, joinOrDash(nil))
	assert.Equal(t, "network", joinOrDash([]string{"network"}))
	assert.Equal(t, "network, storage", joinOrDash([]string{"network", "storage"}))
}

func TestDepsHandler_BuildTable_CollapseEmpty(t *testing.T) {
	handler := NewDepsHandler()

	// Service with only Required deps — only SERVICE and REQUIRED columns shown
	postgres := fixtures.NewServiceConfig("postgres").WithRequired("pgvector").Build()
	headers, rows := handler.buildTable([]types.ServiceConfig{postgres})

	assert.Equal(t, []string{display.HeaderService, display.HeaderRequired}, headers)
	assert.Len(t, rows, 1)
	assert.Equal(t, "postgres", rows[0][0])
	assert.Equal(t, "pgvector", rows[0][1])
}

func TestDepsHandler_BuildTable_AllColumns(t *testing.T) {
	handler := NewDepsHandler()

	postgres := fixtures.NewServiceConfig("postgres").
		WithRequired("pgvector").
		WithConflicts("mysql").
		WithProvides("database", "sql").
		Build()
	redis := fixtures.NewServiceConfig("redis").
		WithSoft("sentinel").
		WithProvides("cache").
		Build()

	headers, rows := handler.buildTable([]types.ServiceConfig{postgres, redis})

	assert.Equal(t, []string{
		display.HeaderService,
		display.HeaderRequired,
		display.HeaderSoft,
		display.HeaderConflicts,
		display.HeaderProvides,
	}, headers)
	assert.Len(t, rows, 2)
}

func TestDepsHandler_BuildTable_DashForEmpty(t *testing.T) {
	handler := NewDepsHandler()

	// postgres has Required but no Soft; redis has Soft but no Required
	postgres := fixtures.NewServiceConfig("postgres").WithRequired("pgvector").Build()
	redis := fixtures.NewServiceConfig("redis").WithSoft("sentinel").Build()

	headers, rows := handler.buildTable([]types.ServiceConfig{postgres, redis})

	assert.Contains(t, headers, display.HeaderRequired)
	assert.Contains(t, headers, display.HeaderSoft)

	reqIdx := headerIndex(headers, display.HeaderRequired)
	softIdx := headerIndex(headers, display.HeaderSoft)

	// redis has no Required — should show dash
	assert.Equal(t, depsDash, rows[1][reqIdx])
	// postgres has no Soft — should show dash
	assert.Equal(t, depsDash, rows[0][softIdx])
}

func TestDepsHandler_BuildTable_Empty(t *testing.T) {
	handler := NewDepsHandler()

	headers, rows := handler.buildTable(nil)
	assert.Equal(t, []string{display.HeaderService}, headers)
	assert.Empty(t, rows)
}

func headerIndex(headers []string, target string) int {
	for i, h := range headers {
		if h == target {
			return i
		}
	}
	return -1
}
