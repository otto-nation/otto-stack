//go:build unit

package core

import (
	"testing"

	"github.com/otto-nation/otto-stack/test/testhelpers"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestYAMLUtilities(t *testing.T) {
	t.Run("validates YAML file detection", func(t *testing.T) {
		assert.True(t, IsYAMLFile("config.yaml"))
		assert.True(t, IsYAMLFile("config.yml"))
		assert.False(t, IsYAMLFile("config.json"))
	})

	t.Run("trims YAML extensions", func(t *testing.T) {
		assert.Equal(t, "config", TrimYAMLExt("config.yaml"))
		assert.Equal(t, "config", TrimYAMLExt("config.yml"))
		assert.Equal(t, "config.json", TrimYAMLExt("config.json"))
	})

	t.Run("finds YAML files", func(t *testing.T) {
		result, err := FindYAMLFile("/tmp", "nonexistent")
		assert.Error(t, err)
		assert.Empty(t, result)
	})
}

func TestFlagParsers(t *testing.T) {
	t.Run("parses up flags", func(t *testing.T) {
		cmd := &cobra.Command{}
		cmd.Flags().Bool("build", false, "build flag")

		flags, err := ParseUpFlags(cmd)
		testhelpers.AssertValidConstructor(t, flags, err, "ParseUpFlags")
	})

	t.Run("parses down flags", func(t *testing.T) {
		cmd := &cobra.Command{}
		cmd.Flags().Bool("volumes", false, "volumes flag")

		flags, err := ParseDownFlags(cmd)
		testhelpers.AssertValidConstructor(t, flags, err, "ParseDownFlags")
	})

	t.Run("parses status flags", func(t *testing.T) {
		cmd := &cobra.Command{}
		cmd.Flags().String("format", "", "format flag")

		flags, err := ParseStatusFlags(cmd)
		testhelpers.AssertValidConstructor(t, flags, err, "ParseStatusFlags")
	})

	t.Run("parses logs flags", func(t *testing.T) {
		cmd := &cobra.Command{}
		cmd.Flags().Bool("follow", false, "follow flag")

		flags, err := ParseLogsFlags(cmd)
		testhelpers.AssertValidConstructor(t, flags, err, "ParseLogsFlags")
	})
}

func TestServiceFlagParsers(t *testing.T) {
	t.Run("parses services flags", func(t *testing.T) {
		cmd := &cobra.Command{}
		cmd.Flags().String("format", "", "format flag")

		flags, err := ParseServicesFlags(cmd)
		testhelpers.AssertValidConstructor(t, flags, err, "ParseServicesFlags")
	})

	t.Run("parses init flags", func(t *testing.T) {
		cmd := &cobra.Command{}
		cmd.Flags().Bool("force", false, "force flag")

		flags, err := ParseInitFlags(cmd)
		testhelpers.AssertValidConstructor(t, flags, err, "ParseInitFlags")
	})

	t.Run("parses version flags", func(t *testing.T) {
		cmd := &cobra.Command{}
		cmd.Flags().String("format", "", "format flag")

		flags, err := ParseVersionFlags(cmd)
		testhelpers.AssertValidConstructor(t, flags, err, "ParseVersionFlags")
	})
}
