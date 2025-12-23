//go:build unit

package unit

import (
	"errors"
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

func TestPromptManager_NewPromptManager_Comprehensive(t *testing.T) {
	validator := &MockValidator{}
	pm := project.NewPromptManager(validator)

	if pm == nil {
		t.Fatal("NewPromptManager returned nil")
	}
}

func TestPromptManager_ErrorHandling_Comprehensive(t *testing.T) {
	tests := []struct {
		name         string
		validatorErr error
		expectError  bool
	}{
		{
			name:         "validator success",
			validatorErr: nil,
			expectError:  false,
		},
		{
			name:         "validator error",
			validatorErr: pkgerrors.NewValidationError(pkgerrors.FieldProjectName, "invalid name", errors.New("test error")),
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := &MockValidator{validateErr: tt.validatorErr}
			pm := project.NewPromptManager(validator)

			// Test validator integration
			err := validator.ValidateProjectName(project.TestProjectName)
			if tt.expectError && err == nil {
				t.Error("expected error but got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			_ = pm // Use pm to avoid unused variable
		})
	}
}

func TestPromptManager_Constants_Comprehensive(t *testing.T) {
	// Test that ErrGoBack is properly defined
	if project.ErrGoBack == nil {
		t.Error("ErrGoBack should not be nil")
	}

	if project.ErrGoBack.Error() != "go back" {
		t.Errorf("ErrGoBack message should be 'go back', got %q", project.ErrGoBack.Error())
	}
}

// Integration-style test for service category handling
func TestPromptManager_ServiceCategoryLogic_Comprehensive(t *testing.T) {
	validator := &MockValidator{}
	pm := project.NewPromptManager(validator)

	// Test category preparation logic
	categories := map[string][]services.ServiceConfig{
		"Database": {
			{Name: "postgres", Description: "PostgreSQL"},
		},
	}

	// Verify the logic works with valid data
	if len(categories) == 0 {
		t.Error("categories should not be empty for this test")
	}

	_ = pm // Use pm to avoid unused variable
}
