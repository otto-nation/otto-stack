package validation

import (
	"github.com/otto-nation/otto-stack/internal/core"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/spf13/cobra"
)

// ValidateUpArgs validates arguments for the up command
func ValidateUpArgs(args []string) error {
	// Service names are optional - if none provided, all enabled services are used
	return nil
}

// ValidateUpFlags validates flags for the up command
func ValidateUpFlags(cmd *cobra.Command) error {
	flags, err := core.ParseUpFlags(cmd)
	if err != nil {
		return pkgerrors.NewValidationError("flags", "parse_flags", err)
	}

	// Additional flag validation can be added here
	_ = flags
	return nil
}
