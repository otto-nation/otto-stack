//go:build unit

package registry

import (
	"io"
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/stretchr/testify/assert"
)

type mockOutput struct {
	messages []string
}

func (m *mockOutput) Info(format string, args ...interface{}) {
	m.messages = append(m.messages, "info")
}
func (m *mockOutput) Success(format string, args ...interface{}) {
	m.messages = append(m.messages, "success")
}
func (m *mockOutput) Warning(format string, args ...interface{}) {
	m.messages = append(m.messages, "warning")
}
func (m *mockOutput) Error(format string, args ...interface{}) {
	m.messages = append(m.messages, "error")
}
func (m *mockOutput) Debug(format string, args ...interface{}) {
	m.messages = append(m.messages, "debug")
}
func (m *mockOutput) Header(format string, args ...interface{}) {
	m.messages = append(m.messages, "header")
}
func (m *mockOutput) Muted(format string, args ...interface{}) {
	m.messages = append(m.messages, "muted")
}
func (m *mockOutput) Writer() io.Writer { return nil }

func TestNewOrphanDisplay(t *testing.T) {
	output := &mockOutput{}
	display := NewOrphanDisplay(output)
	assert.NotNil(t, display)
}

func TestOrphanDisplay_Display(t *testing.T) {
	output := &mockOutput{}
	display := NewOrphanDisplay(output)

	t.Run("handles empty orphans", func(t *testing.T) {
		display.Display([]OrphanInfo{})
		assert.Empty(t, output.messages)
	})

	t.Run("displays critical orphans", func(t *testing.T) {
		output.messages = nil
		orphans := []OrphanInfo{
			{Service: "test", Severity: OrphanSeverityCritical, Reason: "test reason"},
		}
		display.Display(orphans)
		assert.NotEmpty(t, output.messages)
	})

	t.Run("displays warning orphans", func(t *testing.T) {
		output.messages = nil
		orphans := []OrphanInfo{
			{Service: "test", Severity: OrphanSeverityWarning, Reason: "test reason", ProjectsFound: []string{"p1"}},
		}
		display.Display(orphans)
		assert.NotEmpty(t, output.messages)
	})

	t.Run("displays safe orphans", func(t *testing.T) {
		output.messages = nil
		orphans := []OrphanInfo{
			{Service: "test", Severity: OrphanSeveritySafe, Reason: "test reason"},
		}
		display.Display(orphans)
		assert.NotEmpty(t, output.messages)
	})

	t.Run("displays mixed severity orphans", func(t *testing.T) {
		output.messages = nil
		orphans := []OrphanInfo{
			{Service: "test1", Severity: OrphanSeverityCritical, Reason: "critical"},
			{Service: "test2", Severity: OrphanSeverityWarning, Reason: "warning"},
			{Service: "test3", Severity: OrphanSeveritySafe, Reason: "safe"},
		}
		display.Display(orphans)
		assert.NotEmpty(t, output.messages)
	})
}

func TestOrphanDisplay_GroupBySeverity(t *testing.T) {
	output := &mockOutput{}
	display := NewOrphanDisplay(output)

	orphans := []OrphanInfo{
		{Service: "s1", Severity: OrphanSeveritySafe},
		{Service: "s2", Severity: OrphanSeverityWarning},
		{Service: "s3", Severity: OrphanSeverityCritical},
		{Service: "s4", Severity: OrphanSeveritySafe},
	}

	safe, warning, critical := display.groupBySeverity(orphans)
	assert.Len(t, safe, 2)
	assert.Len(t, warning, 1)
	assert.Len(t, critical, 1)
}

var _ base.Output = (*mockOutput)(nil)
