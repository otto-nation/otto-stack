//go:build unit

package display

import (
	"bytes"
	"testing"
	"time"

	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStatusFormatter(t *testing.T) {
	buf := &bytes.Buffer{}
	formatter := NewStatusFormatter(buf)
	assert.NotNil(t, formatter)
}

func TestStatusFormatter_FormatTable_Compact(t *testing.T) {
	buf := &bytes.Buffer{}
	formatter := NewStatusFormatter(buf)

	serviceStatuses := []ServiceStatus{
		{
			Name:      services.ServicePostgres,
			State:     "running",
			Health:    "healthy",
			Ports:     []string{"5432:5432"},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Uptime:    time.Hour,
		},
	}

	err := formatter.FormatTable(serviceStatuses, Options{Compact: true})
	require.NoError(t, err)
	output := buf.String()
	assert.Contains(t, output, services.ServicePostgres)
	assert.Contains(t, output, "running")
}

func TestStatusFormatter_FormatTable_Full(t *testing.T) {
	buf := &bytes.Buffer{}
	formatter := NewStatusFormatter(buf)

	serviceStatuses := []ServiceStatus{
		{
			Name:   services.ServiceRedis,
			State:  "running",
			Health: "healthy",
		},
	}

	err := formatter.FormatTable(serviceStatuses, Options{Compact: false})
	require.NoError(t, err)
	output := buf.String()
	assert.Contains(t, output, services.ServiceRedis)
}

func TestStatusFormatter_FormatResourceSummary(t *testing.T) {
	buf := &bytes.Buffer{}
	formatter := NewStatusFormatter(buf)

	services := []ServiceStatus{
		{Name: "service1", State: "running"},
		{Name: "service2", State: "running"},
		{Name: "service3", State: "stopped"},
	}

	formatter.formatResourceSummary(services)
	output := buf.String()
	assert.Contains(t, output, "Summary: 3 total")
	assert.Contains(t, output, "running")
}

func TestStatusFormatter_CreateSummary(t *testing.T) {
	formatter := NewStatusFormatter(&bytes.Buffer{})

	services := []ServiceStatus{
		{State: "running"},
		{State: "running"},
		{State: "stopped"},
	}

	summary := formatter.createSummary(services)
	assert.Equal(t, 2, summary["running"])
	assert.Equal(t, 1, summary["stopped"])
}
