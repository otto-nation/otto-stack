package services

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/constants"
	"github.com/stretchr/testify/assert"
)

func TestConnectionConfig(t *testing.T) {
	config := ConnectionConfig{
		Client:      constants.ClientPsql,
		DefaultUser: "postgres",
		DefaultPort: constants.DefaultPortPOSTGRES_port,
		UserFlag:    "-U",
		HostFlag:    "-h",
		PortFlag:    "-p",
		DBFlag:      "-d",
		ExtraFlags:  []string{"--no-password"},
	}

	assert.Equal(t, constants.ClientPsql, config.Client)
	assert.Equal(t, "postgres", config.DefaultUser)
	assert.Equal(t, constants.DefaultPortPOSTGRES_port, config.DefaultPort)
	assert.Equal(t, "-U", config.UserFlag)
	assert.Contains(t, config.ExtraFlags, "--no-password")
}

func TestGetConnectionConfig(t *testing.T) {
	// Test error case for non-existent service
	_, err := GetConnectionConfig("nonexistent-service")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to load services")
}
