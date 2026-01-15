//go:build unit

package ci

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestShellType_Validation(t *testing.T) {
	t.Run("validates supported shell types", func(t *testing.T) {
		assert.True(t, ShellTypeBash.IsValid())
		assert.True(t, ShellTypeZsh.IsValid())
		assert.True(t, ShellTypeFish.IsValid())
		assert.True(t, ShellTypePowerShell.IsValid())
	})

	t.Run("rejects invalid shell types", func(t *testing.T) {
		invalidShell := ShellType("invalid")
		assert.False(t, invalidShell.IsValid())

		emptyShell := ShellType("")
		assert.False(t, emptyShell.IsValid())
	})
}

func TestShellType_Constants(t *testing.T) {
	t.Run("validates shell type constants", func(t *testing.T) {
		assert.Equal(t, "bash", string(ShellTypeBash))
		assert.Equal(t, "zsh", string(ShellTypeZsh))
		assert.Equal(t, "fish", string(ShellTypeFish))
		assert.Equal(t, "powershell", string(ShellTypePowerShell))
	})
}

func TestAllShellTypeStrings(t *testing.T) {
	t.Run("returns all shell types", func(t *testing.T) {
		shells := AllShellTypeStrings()

		assert.Len(t, shells, 4)
		assert.Contains(t, shells, "bash")
		assert.Contains(t, shells, "zsh")
		assert.Contains(t, shells, "fish")
		assert.Contains(t, shells, "powershell")
	})
}

func TestGetFlags(t *testing.T) {
	t.Run("extracts flags from command", func(t *testing.T) {
		cmd := &cobra.Command{}
		cmd.Flags().Bool(core.FlagQuiet, true, "quiet")
		cmd.Flags().Bool(core.FlagJSON, false, "json")
		cmd.Flags().Bool(core.FlagNoColor, true, "no-color")
		cmd.Flags().Bool(core.FlagNonInteractive, false, "non-interactive")
		cmd.Flags().Bool(core.FlagDryRun, true, "dry-run")

		flags := GetFlags(cmd)

		assert.True(t, flags.Quiet)
		assert.False(t, flags.JSON)
		assert.True(t, flags.NoColor)
		assert.False(t, flags.NonInteractive)
		assert.True(t, flags.DryRun)
	})

	t.Run("handles missing flags gracefully", func(t *testing.T) {
		cmd := &cobra.Command{}
		// No flags added

		flags := GetFlags(cmd)

		// Should default to false values
		assert.False(t, flags.Quiet)
		assert.False(t, flags.JSON)
		assert.False(t, flags.NoColor)
		assert.False(t, flags.NonInteractive)
		assert.False(t, flags.DryRun)
	})
}

func TestFlags_Structure(t *testing.T) {
	t.Run("validates Flags structure", func(t *testing.T) {
		flags := Flags{
			Quiet:          true,
			JSON:           false,
			NoColor:        true,
			NonInteractive: false,
			DryRun:         true,
		}

		assert.True(t, flags.Quiet)
		assert.False(t, flags.JSON)
		assert.True(t, flags.NoColor)
		assert.False(t, flags.NonInteractive)
		assert.True(t, flags.DryRun)
	})
}

func TestFormatError(t *testing.T) {
	t.Run("formats error without exiting", func(t *testing.T) {
		flags := Flags{Quiet: false, JSON: false}
		testErr := assert.AnError

		result := FormatError(flags, testErr)

		// Should return the same error without calling os.Exit()
		assert.Equal(t, testErr, result)
	})

	t.Run("handles quiet mode", func(t *testing.T) {
		flags := Flags{Quiet: true, JSON: false}
		testErr := assert.AnError

		result := FormatError(flags, testErr)

		// Should return error without output in quiet mode
		assert.Equal(t, testErr, result)
	})

	t.Run("handles JSON mode", func(t *testing.T) {
		flags := Flags{Quiet: false, JSON: true}
		testErr := assert.AnError

		result := FormatError(flags, testErr)

		// Should return error and output JSON
		assert.Equal(t, testErr, result)
	})
}
