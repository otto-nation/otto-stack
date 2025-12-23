package services

import (
	"os"
	"strings"
	"testing"

	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFileGenerator(t *testing.T) {
	fg := NewFileGenerator()
	assert.NotNil(t, fg)
}

func TestFileGenerator_GenerateComposeFile(t *testing.T) {
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	os.Chdir(tempDir)

	fg := NewFileGenerator()

	tests := []struct {
		name         string
		services     []string
		projectName  string
		expectError  bool
		errorMessage string
	}{
		{
			name:        "valid services",
			services:    []string{"postgres", "redis"},
			projectName: "test-project",
			expectError: false,
		},
		{
			name:         "empty services",
			services:     []string{},
			projectName:  "test-project",
			expectError:  true,
			errorMessage: "no services provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fg.GenerateComposeFile(tt.services, tt.projectName)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMessage != "" {
					assert.Contains(t, err.Error(), tt.errorMessage)
				}
				return
			}

			require.NoError(t, err)
			assert.FileExists(t, docker.DockerComposeFilePath)

			// Read and verify content
			content, err := os.ReadFile(docker.DockerComposeFilePath)
			require.NoError(t, err)

			contentStr := string(content)
			assert.Contains(t, contentStr, tt.projectName)
			assert.Contains(t, contentStr, "services:")

			for _, service := range tt.services {
				assert.Contains(t, contentStr, service)
			}
		})
	}
}

func TestFileGenerator_GenerateEnvFile(t *testing.T) {
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	os.Chdir(tempDir)

	fg := NewFileGenerator()

	tests := []struct {
		name         string
		services     []string
		projectName  string
		expectError  bool
		errorMessage string
	}{
		{
			name:        "valid input",
			services:    []string{"postgres", "redis"},
			projectName: "test-project",
			expectError: false,
		},
		{
			name:         "empty project name",
			services:     []string{"postgres"},
			projectName:  "",
			expectError:  true,
			errorMessage: "project name cannot be empty",
		},
		{
			name:        "empty services",
			services:    []string{},
			projectName: "test-project",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fg.GenerateEnvFile(tt.services, tt.projectName)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMessage != "" {
					assert.Contains(t, err.Error(), tt.errorMessage)
				}
				return
			}

			require.NoError(t, err)

			// Verify file was created
			envPath := ".env.generated"
			assert.FileExists(t, envPath)

			// Read and verify content
			content, err := os.ReadFile(envPath)
			require.NoError(t, err)

			contentStr := string(content)
			assert.Contains(t, contentStr, tt.projectName)
			assert.Contains(t, contentStr, "PROJECT_NAME="+tt.projectName)
			assert.Contains(t, contentStr, "SERVICES="+strings.Join(tt.services, ","))

			// Cleanup for next test
			os.Remove(envPath)
		})
	}
}

func TestFileGenerator_BuildEnvContent(t *testing.T) {
	fg := NewFileGenerator()

	services := []string{"postgres", "redis"}
	projectName := "test-app"

	content := fg.buildEnvContent(services, projectName)

	assert.Contains(t, content, "# Generated environment file for test-app")
	assert.Contains(t, content, "PROJECT_NAME=test-app")
	assert.Contains(t, content, "SERVICES=postgres,redis")

	// Should be valid environment file format
	lines := strings.Split(content, "\n")
	assert.True(t, len(lines) >= 3)

	// First line should be comment
	assert.True(t, strings.HasPrefix(lines[0], "#"))
}
