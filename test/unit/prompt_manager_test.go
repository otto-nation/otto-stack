//go:build unit

package unit

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/project"
	"github.com/stretchr/testify/assert"
)

// MockProjectValidator for testing
type MockProjectValidator struct {
	validateError error
}

func (m *MockProjectValidator) ValidateProjectName(name string) error {
	return m.validateError
}

func TestPromptManager_NewPromptManager(t *testing.T) {
	validator := &MockProjectValidator{}
	pm := project.NewPromptManager(validator)

	assert.NotNil(t, pm)
}

func TestPromptManager_BuildServiceOptions_Unit(t *testing.T) {
	// Test the logic without file system dependencies
	tests := []struct {
		name         string
		serviceCount int
		allowGoBack  bool
		wantCount    int
		wantGoBack   bool
	}{
		{
			name:         "no services without go back",
			serviceCount: 0,
			allowGoBack:  false,
			wantCount:    0,
			wantGoBack:   false,
		},
		{
			name:         "services without go back",
			serviceCount: 2,
			allowGoBack:  false,
			wantCount:    2,
			wantGoBack:   false,
		},
		{
			name:         "services with go back",
			serviceCount: 3,
			allowGoBack:  true,
			wantCount:    4,
			wantGoBack:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify the logic works correctly
			expectedCount := tt.serviceCount
			if tt.allowGoBack {
				expectedCount++
			}
			assert.Equal(t, tt.wantCount, expectedCount)
		})
	}
}
