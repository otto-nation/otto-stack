//go:build unit

package validation

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestValidateUpArgs_Empty(t *testing.T) {
	err := ValidateUpArgs([]string{})
	assert.NoError(t, err)
}

func TestValidateUpArgs_MultipleServices(t *testing.T) {
	args := []string{"postgres", "redis"}
	err := ValidateUpArgs(args)
	assert.NoError(t, err)
}

func TestValidateUpArgs_SingleService(t *testing.T) {
	args := []string{"postgres"}
	err := ValidateUpArgs(args)
	assert.NoError(t, err)
}

func TestValidateUpFlags_Valid(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().Bool("build", false, "build images")
	cmd.Flags().Bool("force-recreate", false, "force recreate")
	cmd.Flags().Bool("no-deps", false, "no dependencies")
	cmd.Flags().Bool("detach", false, "detach mode")

	err := ValidateUpFlags(cmd)
	if err != nil {
		assert.Contains(t, err.Error(), "validation error")
	} else {
		assert.NoError(t, err)
	}
}

func TestValidateUpFlags_MissingFlags(t *testing.T) {
	cmd := &cobra.Command{}

	err := ValidateUpFlags(cmd)
	if err != nil {
		assert.Contains(t, err.Error(), "validation error")
	}
}
