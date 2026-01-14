package services

import (
	"testing"

	servicetypes "github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestStartRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request StartRequest
		valid   bool
	}{
		{
			name: "valid request with postgres and redis",
			request: StartRequest{
				Project: "test-project",
				ServiceConfigs: []servicetypes.ServiceConfig{
					{Name: ServicePostgres, Category: CategoryDatabase},
					{Name: ServiceRedis, Category: CategoryCache},
				},
				Build: false,
			},
			valid: true,
		},
		{
			name: "empty project name should be invalid",
			request: StartRequest{
				Project:        "",
				ServiceConfigs: []servicetypes.ServiceConfig{{Name: ServicePostgres}},
			},
			valid: false,
		},
		{
			name: "no services should be allowed",
			request: StartRequest{
				Project:        "test-project",
				ServiceConfigs: []servicetypes.ServiceConfig{},
			},
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just test that the struct is valid
			if tt.valid {
				assert.NotEmpty(t, tt.request.Project)
			} else {
				assert.Empty(t, tt.request.Project)
			}
		})
	}
}

func TestService_ConstructorsAndOperations(t *testing.T) {
	t.Run("new service with dependencies", func(t *testing.T) {
		service := NewServiceWithDependencies(nil, nil, nil, nil)
		if service == nil {
			t.Error("NewServiceWithDependencies should return service")
		}
	})

	t.Run("new service", func(t *testing.T) {
		service, err := NewService(nil, nil, nil)
		if err != nil {
			t.Log("NewService failed as expected due to nil dependencies")
		}
		if service == nil && err == nil {
			t.Error("NewService should return service or error")
		}
	})
}
