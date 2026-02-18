package core

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestParseDepsFlags(t *testing.T) {
	cmd := &cobra.Command{}
	flags, err := ParseDepsFlags(cmd)
	assert.NoError(t, err)
	assert.NotNil(t, flags)
}

func TestParseValidateFlags(t *testing.T) {
	cmd := &cobra.Command{}
	flags, err := ParseValidateFlags(cmd)
	assert.NoError(t, err)
	assert.NotNil(t, flags)
}

func TestParseRestartFlags(t *testing.T) {
	cmd := &cobra.Command{}
	flags, err := ParseRestartFlags(cmd)
	assert.NoError(t, err)
	assert.NotNil(t, flags)
}

func TestParseDoctorFlags(t *testing.T) {
	cmd := &cobra.Command{}
	flags, err := ParseDoctorFlags(cmd)
	assert.NoError(t, err)
	assert.NotNil(t, flags)
}

func TestParseCleanupFlags(t *testing.T) {
	cmd := &cobra.Command{}
	flags, err := ParseCleanupFlags(cmd)
	assert.NoError(t, err)
	assert.NotNil(t, flags)
}

func TestParseConflictsFlags(t *testing.T) {
	cmd := &cobra.Command{}
	flags, err := ParseConflictsFlags(cmd)
	assert.NoError(t, err)
	assert.NotNil(t, flags)
}

func TestParseWebInterfacesFlags(t *testing.T) {
	cmd := &cobra.Command{}
	flags, err := ParseWebInterfacesFlags(cmd)
	assert.NoError(t, err)
	assert.NotNil(t, flags)
}
