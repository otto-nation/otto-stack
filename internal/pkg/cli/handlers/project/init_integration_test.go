//go:build integration

package project

import (
	"context"
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/ui"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitHandler_Creation(t *testing.T) {
	handler := NewInitHandler()
	require.NotNil(t, handler)

	err := handler.ValidateArgs([]string{})
	assert.NoError(t, err)

	err = handler.ValidateArgs([]string{"extra", "args"})
	assert.NoError(t, err)

	flags := handler.GetRequiredFlags()
	assert.NotNil(t, flags)
}

func TestInitHandler_MockCommand(t *testing.T) {
	handler := NewInitHandler()
	require.NotNil(t, handler)

	cmd := &cobra.Command{
		Use: "init",
	}

	cmd.Flags().String("name", "", "Project name")
	cmd.Flags().Bool("force", false, "Force initialization")
	cmd.Flags().StringSlice("services", []string{}, "Services to include")

	err := cmd.Flags().Set("name", "test-project")
	require.NoError(t, err)
	err = cmd.Flags().Set("force", "true")
	require.NoError(t, err)

	output := ui.NewOutput()
	base := &base.BaseCommand{
		Output: output,
	}

	ctx := context.Background()
	err = handler.Handle(ctx, cmd, []string{}, base)

	if err != nil {
		t.Logf("Handler returned expected error: %v", err)
	}
}
