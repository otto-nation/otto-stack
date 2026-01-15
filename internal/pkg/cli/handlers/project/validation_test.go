//go:build unit

package project

import (
	"runtime"
	"strings"
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestValidateProjectName(t *testing.T) {
	handler := NewInitHandler()

	tests := []struct {
		name        string
		projectName string
		expectError bool
	}{
		{"valid name", TestProjectName, false},
		{"valid with underscore", "my_project", false},
		{"valid with numbers", "project123", false},
		{"empty name", "", true},
		{"too short", "a", true},
		{"too long", strings.Repeat("a", 51), true},
		{"invalid characters", TestProjectNameInvalid, true},
		{"starts with hyphen", "-project", true},
		{"starts with underscore", "_project", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.validateProjectName(tt.projectName)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateServiceConfigs(t *testing.T) {
	handler := NewInitHandler()

	tests := []struct {
		name           string
		serviceConfigs []types.ServiceConfig
		expectError    bool
	}{
		{"empty services", []types.ServiceConfig{}, true},
		{"nil services", nil, true},
		{"valid services", []types.ServiceConfig{{Name: services.ServicePostgres}, {Name: services.ServiceRedis}}, false},
		{"duplicate services", []types.ServiceConfig{{Name: services.ServicePostgres}, {Name: services.ServicePostgres}}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.validateServiceConfigs(tt.serviceConfigs)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIsCommandAvailable(t *testing.T) {
	// Use commands that definitely exist on each platform
	existingCommand := "echo" // Available on Unix systems
	if runtime.GOOS == "windows" {
		existingCommand = "where" // Built-in Windows command
	}

	tests := []struct {
		name     string
		command  string
		expected bool
	}{
		{"empty command", "", false},
		{"existing command", existingCommand, true},
		{"nonexistent command", "nonexistent-command-12345", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isCommandAvailable(tt.command)
			assert.Equal(t, tt.expected, result)
		})
	}
}
