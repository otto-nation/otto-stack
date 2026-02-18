//go:build unit

package core

import (
	"testing"

	"github.com/otto-nation/otto-stack/test/testhelpers"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestYAMLUtilities(t *testing.T) {
	assert.True(t, IsYAMLFile("config.yaml"))
	assert.True(t, IsYAMLFile("config.yml"))
	assert.False(t, IsYAMLFile("config.json"))

	assert.Equal(t, "config", TrimYAMLExt("config.yaml"))
	assert.Equal(t, "config", TrimYAMLExt("config.yml"))
	assert.Equal(t, "config.json", TrimYAMLExt("config.json"))
}

func TestFlagParsers(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().Bool("build", false, "build flag")

	flags, err := ParseUpFlags(cmd)
	testhelpers.AssertValidConstructor(t, flags, err, "ParseUpFlags")

	cmd = &cobra.Command{}
	cmd.Flags().Bool("volumes", false, "volumes flag")

	flags2, err := ParseDownFlags(cmd)
	testhelpers.AssertValidConstructor(t, flags2, err, "ParseDownFlags")

	cmd = &cobra.Command{}
	cmd.Flags().String("format", "", "format flag")

	flags3, err := ParseStatusFlags(cmd)
	testhelpers.AssertValidConstructor(t, flags3, err, "ParseStatusFlags")

	cmd = &cobra.Command{}
	cmd.Flags().Bool("follow", false, "follow flag")

	flags4, err := ParseLogsFlags(cmd)
	testhelpers.AssertValidConstructor(t, flags4, err, "ParseLogsFlags")
}

func TestServiceFlagParsers(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String("format", "", "format flag")

	flags, err := ParseServicesFlags(cmd)
	testhelpers.AssertValidConstructor(t, flags, err, "ParseServicesFlags")

	cmd = &cobra.Command{}
	cmd.Flags().Bool("force", false, "force flag")

	flags2, err := ParseInitFlags(cmd)
	testhelpers.AssertValidConstructor(t, flags2, err, "ParseInitFlags")

	cmd = &cobra.Command{}
	cmd.Flags().String("format", "", "format flag")

	flags3, err := ParseVersionFlags(cmd)
	testhelpers.AssertValidConstructor(t, flags3, err, "ParseVersionFlags")
}
