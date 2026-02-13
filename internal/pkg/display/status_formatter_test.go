//go:build unit

package display

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStatusFormatter(t *testing.T) {
	buf := &bytes.Buffer{}
	formatter := NewStatusFormatter(buf)
	assert.NotNil(t, formatter)
}

func TestStatusFormatter_FormatTable(t *testing.T) {
	t.Run("formats compact table", func(t *testing.T) {
		buf := &bytes.Buffer{}
		formatter := NewStatusFormatter(buf)

		services := []ServiceStatus{
			{
				Name:      "postgres",
				State:     "running",
				Health:    "healthy",
				Ports:     []string{"5432:5432"},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Uptime:    time.Hour,
			},
		}

		err := formatter.FormatTable(services, Options{Compact: true})
		require.NoError(t, err)
		output := buf.String()
		assert.Contains(t, output, "postgres")
		assert.Contains(t, output, "running")
	})

	t.Run("formats full table", func(t *testing.T) {
		buf := &bytes.Buffer{}
		formatter := NewStatusFormatter(buf)

		services := []ServiceStatus{
			{
				Name:   "redis",
				State:  "running",
				Health: "healthy",
			},
		}

		err := formatter.FormatTable(services, Options{Compact: false})
		require.NoError(t, err)
		output := buf.String()
		assert.Contains(t, output, "redis")
	})
}
