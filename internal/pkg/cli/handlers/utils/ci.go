package utils

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/otto-nation/otto-stack/internal/pkg/constants"
	"github.com/spf13/cobra"
)

// CIFlags represents CI-friendly command flags
type CIFlags struct {
	Quiet          bool
	JSON           bool
	NoColor        bool
	NonInteractive bool
	Strict         bool
}

// GetCIFlags extracts CI-friendly flags from command
func GetCIFlags(cmd *cobra.Command) CIFlags {
	quiet, _ := cmd.Flags().GetBool(constants.FlagQuiet)
	jsonOutput, _ := cmd.Flags().GetBool(constants.FlagJSON)
	noColor, _ := cmd.Flags().GetBool(constants.FlagNoColor)
	nonInteractive, _ := cmd.Flags().GetBool(constants.FlagNonInteractive)
	strict, _ := cmd.Flags().GetBool(constants.FlagStrict)

	return CIFlags{
		Quiet:          quiet,
		JSON:           jsonOutput,
		NoColor:        noColor,
		NonInteractive: nonInteractive,
		Strict:         strict,
	}
}

// OutputResult outputs result in CI-friendly format
func OutputResult(flags CIFlags, result any, exitCode int) {
	if flags.JSON {
		outputJSON(result, exitCode)
		return
	}

	if !flags.Quiet {
		outputTable(result)
	}

	if exitCode != constants.ExitSuccess {
		os.Exit(exitCode)
	}
}

// HandleError handles errors in CI-friendly way
func HandleError(flags CIFlags, err error) {
	if flags.JSON {
		errorResult := map[string]any{
			"error":     err.Error(),
			"exit_code": constants.ExitError,
		}
		_ = json.NewEncoder(os.Stdout).Encode(errorResult)
	} else if !flags.Quiet {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}

	os.Exit(constants.ExitError)
}

func outputJSON(result any, exitCode int) {
	output := map[string]any{
		"result":    result,
		"exit_code": exitCode,
	}

	_ = json.NewEncoder(os.Stdout).Encode(output)

	if exitCode != constants.ExitSuccess {
		os.Exit(exitCode)
	}
}

func outputTable(result any) {
	fmt.Printf("%+v\n", result)
}
