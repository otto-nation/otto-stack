package version

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/constants"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDetectProjectVersion(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func(t *testing.T) string
		expectError bool
		expectFound bool
	}{
		{
			name: "detect from " + constants.ConfigFileName,
			setupFunc: func(t *testing.T) string {
				tmpDir := t.TempDir()
				configContent := `project:
  name: test
  environment: development

version_config:
  required_version: ">=1.0.0"
`
				configPath := filepath.Join(tmpDir, constants.ConfigFileName)
				require.NoError(t, os.WriteFile(configPath, []byte(configContent), 0644))
				return tmpDir
			},
			expectError: false,
			expectFound: true,
		},
		{
			name: "no version config found",
			setupFunc: func(t *testing.T) string {
				tmpDir := t.TempDir()
				configContent := `project:
  name: test
  environment: development
`
				configPath := filepath.Join(tmpDir, constants.ConfigFileName)
				require.NoError(t, os.WriteFile(configPath, []byte(configContent), 0644))
				return tmpDir
			},
			expectError: false,
			expectFound: false,
		},
		{
			name: "no config file",
			setupFunc: func(t *testing.T) string {
				return t.TempDir()
			},
			expectError: false,
			expectFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectPath := tt.setupFunc(t)

			constraint, err := DetectProjectVersion(projectPath)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, constraint)

			if tt.expectFound {
				assert.NotEqual(t, "*", constraint.Operator)
			} else {
				assert.Equal(t, "*", constraint.Operator)
			}
		})
	}
}

func TestValidateProjectVersion(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func(t *testing.T) string
		expectError bool
	}{
		{
			name: "valid version constraint",
			setupFunc: func(t *testing.T) string {
				tmpDir := t.TempDir()
				configContent := `version_config:
  required_version: ">=1.0.0"
`
				configPath := filepath.Join(tmpDir, constants.ConfigFileName)
				require.NoError(t, os.WriteFile(configPath, []byte(configContent), 0644))
				return tmpDir
			},
			expectError: false,
		},
		{
			name: "no version constraint",
			setupFunc: func(t *testing.T) string {
				return t.TempDir()
			},
			expectError: false, // Default constraint is valid
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectPath := tt.setupFunc(t)

			err := ValidateProjectVersion(projectPath)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
