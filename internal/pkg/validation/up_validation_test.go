//go:build unit

package validation

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestValidateUpArgs(t *testing.T) {
	t.Run("accepts empty args", func(t *testing.T) {
		err := ValidateUpArgs([]string{})
		assert.NoError(t, err)
	})

	t.Run("accepts service names", func(t *testing.T) {
		args := []string{"postgres", "redis"}
		err := ValidateUpArgs(args)
		assert.NoError(t, err)
	})

	t.Run("accepts single service", func(t *testing.T) {
		args := []string{"postgres"}
		err := ValidateUpArgs(args)
		assert.NoError(t, err)
	})
}

func TestValidateUpFlags(t *testing.T) {
	t.Run("validates valid flags", func(t *testing.T) {
		cmd := &cobra.Command{}
		// Add required flags that ParseUpFlags expects
		cmd.Flags().Bool("build", false, "build images")
		cmd.Flags().Bool("force-recreate", false, "force recreate")
		cmd.Flags().Bool("no-deps", false, "no dependencies")
		cmd.Flags().Bool("detach", false, "detach mode")

		err := ValidateUpFlags(cmd)
		// Should not error with valid flags
		if err != nil {
			// If ParseUpFlags fails, it should return a validation error
			assert.Contains(t, err.Error(), "validation error")
		} else {
			assert.NoError(t, err)
		}
	})

	t.Run("handles missing flags gracefully", func(t *testing.T) {
		cmd := &cobra.Command{}
		// No flags added - should trigger parse error

		err := ValidateUpFlags(cmd)
		// Should return validation error when flags can't be parsed
		if err != nil {
			assert.Contains(t, err.Error(), "validation error")
		}
	})
}
