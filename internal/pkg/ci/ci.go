package ci

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/logger"
	"github.com/spf13/cobra"
)

// Flags represents CI-friendly command flags
type Flags struct {
	Quiet          bool
	JSON           bool
	NoColor        bool
	NonInteractive bool
	DryRun         bool
}

// GetFlags extracts CI-friendly flags from command
func GetFlags(cmd *cobra.Command) Flags {
	quiet, _ := cmd.Flags().GetBool(core.FlagQuiet)
	jsonOutput, _ := cmd.Flags().GetBool(core.FlagJSON)
	noColor, _ := cmd.Flags().GetBool(core.FlagNoColor)
	nonInteractive, _ := cmd.Flags().GetBool(core.FlagNonInteractive)
	dryRun, _ := cmd.Flags().GetBool(core.FlagDryRun)

	return Flags{
		Quiet:          quiet,
		JSON:           jsonOutput,
		NoColor:        noColor,
		NonInteractive: nonInteractive,
		DryRun:         dryRun,
	}
}

// OutputResult outputs result in CI-friendly format
func OutputResult(flags Flags, result any, exitCode int) {
	if flags.JSON {
		outputJSON(result, exitCode)
		return
	}

	if !flags.Quiet {
		outputTable(result)
	}

	if exitCode != core.ExitSuccess {
		os.Exit(exitCode)
	}
}

// FormatError formats errors for CI output without exiting
// Use this in handlers that should return errors instead of calling os.Exit()
func FormatError(flags Flags, err error) error {
	if flags.JSON {
		errorResult := map[string]any{
			"error":     err.Error(),
			"exit_code": core.ExitError,
		}
		_ = json.NewEncoder(os.Stdout).Encode(errorResult)
	} else if !flags.Quiet {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
	return err
}

func outputJSON(result any, exitCode int) {
	output := map[string]any{
		"result":    result,
		"exit_code": exitCode,
	}

	_ = json.NewEncoder(os.Stdout).Encode(output)

	if exitCode != core.ExitSuccess {
		os.Exit(exitCode)
	}
}

func outputTable(result any) {
	logger.Debug("Outputting table result", logger.LogFieldResult, result)
	fmt.Printf("%+v\n", result)
}
