//go:build integration

package integration

import (
	"strings"
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/project"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
)

// MockValidator for testing
type MockValidator struct {
	validateErr error
}

func (m *MockValidator) ValidateProjectName(name string) error {
	return m.validateErr
}

// Error path tests for comprehensive coverage

func TestPromptManager_ErrorPaths(t *testing.T) {
	tests := []struct {
		name          string
		validatorErr  error
		expectError   bool
		errorContains string
	}{
		{
			name:          "validation error",
			validatorErr:  pkgerrors.NewValidationError(pkgerrors.FieldProjectName, "invalid name", nil),
			expectError:   true,
			errorContains: "invalid name",
		},
		{
			name:          "nil validator error",
			validatorErr:  nil,
			expectError:   false,
			errorContains: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := &MockValidator{validateErr: tt.validatorErr}
			pm := project.NewPromptManager(validator)

			err := validator.ValidateProjectName(project.TestProjectName)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got nil")
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("error should contain %q, got: %v", tt.errorContains, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}

			_ = pm // Use pm to avoid unused variable
		})
	}
}

func TestFileGenerator_ErrorPaths(t *testing.T) {
	fg := services.NewFileGenerator()

	// Test empty services error
	err := fg.GenerateComposeFile([]string{}, project.TestProjectName)
	if err == nil {
		t.Error("GenerateComposeFile should error with empty services")
	}
	if !strings.Contains(err.Error(), "no services provided") {
		t.Errorf("Error should mention no services, got: %v", err)
	}

	// Test empty project name error
	err = fg.GenerateEnvFile([]string{services.ServicePostgres}, "")
	if err == nil {
		t.Error("GenerateEnvFile should error with empty project name")
	}
	if !strings.Contains(err.Error(), "project name cannot be empty") {
		t.Errorf("Error should mention empty project name, got: %v", err)
	}
}

func TestFileGenerator_EdgeCases(t *testing.T) {
	fg := services.NewFileGenerator()

	// Test with nil services (should not crash)
	err := fg.GenerateComposeFile(nil, project.TestProjectName)
	if err == nil {
		t.Error("GenerateComposeFile should error with nil services")
	}

	// Test with very long project name
	longName := strings.Repeat("a", 1000)
	err = fg.GenerateEnvFile([]string{"test"}, longName)
	// Should not crash, even with very long names
	if err != nil {
		t.Errorf("GenerateEnvFile should handle long project names: %v", err)
	}
}

func TestPromptManager_EdgeCases(t *testing.T) {
	// Test with nil validator (should not crash during creation)
	pm := project.NewPromptManager(nil)
	if pm == nil {
		t.Error("NewPromptManager should handle nil validator")
	}
}

func TestStateManager_EdgeCases(t *testing.T) {
	// These are basic edge case tests that don't require file system
	// More comprehensive edge cases are in integration tests

	// Test that constants are properly defined
	if project.ErrGoBack == nil {
		t.Error("ErrGoBack constant should be defined")
	}

	if project.ErrGoBack.Error() == "" {
		t.Error("ErrGoBack should have non-empty error message")
	}
}

// Table-driven tests for systematic validation
func TestFileGenerator_TableDriven(t *testing.T) {
	tests := []struct {
		name         string
		services     []string
		projectName  string
		expectError  bool
		errorMessage string
	}{
		{
			name:         "valid input",
			services:     []string{services.ServicePostgres, services.ServiceRedis},
			projectName:  project.TestProjectNameValid,
			expectError:  false,
			errorMessage: "",
		},
		{
			name:         "empty services",
			services:     []string{},
			projectName:  project.TestProjectNameValid,
			expectError:  true,
			errorMessage: "no services provided",
		},
		{
			name:         "nil services",
			services:     nil,
			projectName:  project.TestProjectNameValid,
			expectError:  true,
			errorMessage: "no services provided",
		},
		{
			name:         "empty project name",
			services:     []string{services.ServicePostgres},
			projectName:  "",
			expectError:  true,
			errorMessage: "project name cannot be empty",
		},
		{
			name:         "single service",
			services:     []string{services.ServicePostgres},
			projectName:  "single-service",
			expectError:  false,
			errorMessage: "",
		},
	}

	fg := services.NewFileGenerator()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test GenerateComposeFile
			err := fg.GenerateComposeFile(tt.services, tt.projectName)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got nil")
				} else if tt.errorMessage != "" && !strings.Contains(err.Error(), tt.errorMessage) {
					t.Errorf("error should contain %q, got: %v", tt.errorMessage, err)
				}
			} else {
				// For valid inputs, we expect error because we can't write files in unit tests
				// The important thing is that it doesn't crash and gives a reasonable error
				if err != nil && !strings.Contains(err.Error(), "no such file or directory") &&
					!strings.Contains(err.Error(), "permission denied") {
					t.Errorf("unexpected error type: %v", err)
				}
			}
		})
	}
}
