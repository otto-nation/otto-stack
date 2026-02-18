package operations

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestLogsHandler_GetRequiredFlags(t *testing.T) {
	handler := NewLogsHandler()
	flags := handler.GetRequiredFlags()
	assert.Empty(t, flags)
}

func TestLogsHandler_getFlag(t *testing.T) {
	handler := NewLogsHandler()
	cmd := &cobra.Command{}
	cmd.Flags().Bool("follow", false, "")
	cmd.Flags().Set("follow", "true")

	assert.True(t, handler.getFlag(cmd, "follow"))
	assert.False(t, handler.getFlag(cmd, "nonexistent"))
}

func TestLogsHandler_getTailFlag(t *testing.T) {
	handler := NewLogsHandler()
	cmd := &cobra.Command{}
	cmd.Flags().String("tail", "", "")

	assert.Equal(t, core.DefaultLogTailLines, handler.getTailFlag(cmd))

	cmd.Flags().Set("tail", "100")
	assert.Equal(t, "100", handler.getTailFlag(cmd))
}
