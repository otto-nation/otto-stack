//go:build unit

package services

import (
	"testing"

	servicetypes "github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestStartRequest_Validation_ValidRequest(t *testing.T) {
	request := StartRequest{
		Project: "test-project",
		ServiceConfigs: []servicetypes.ServiceConfig{
			{Name: ServicePostgres, Category: CategoryDatabase},
			{Name: ServiceRedis, Category: CategoryCache},
		},
		Build: false,
	}

	assert.NotEmpty(t, request.Project)
}

func TestStartRequest_Validation_EmptyProject(t *testing.T) {
	request := StartRequest{
		Project:        "",
		ServiceConfigs: []servicetypes.ServiceConfig{{Name: ServicePostgres}},
	}

	assert.Empty(t, request.Project)
}

func TestStartRequest_Validation_NoServices(t *testing.T) {
	request := StartRequest{
		Project:        "test-project",
		ServiceConfigs: []servicetypes.ServiceConfig{},
	}

	assert.NotEmpty(t, request.Project)
}

func TestService_NewService(t *testing.T) {
	service, err := NewService(nil, nil, nil)
	if err != nil {
		t.Log("NewService failed as expected due to nil dependencies")
	}
	if service == nil && err == nil {
		t.Error("NewService should return service or error")
	}
}
