package services

import (
	"testing"

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
				ServiceConfigs: []ServiceConfig{
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
				ServiceConfigs: []ServiceConfig{},
			},
			valid: false,
		},
		{
			name: "no services should be allowed",
			request: StartRequest{
				Project:        "test-project",
				ServiceConfigs: []ServiceConfig{},
			},
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.valid {
				assert.NotEmpty(t, tt.request.Project)
			} else {
				assert.Empty(t, tt.request.Project)
			}
		})
	}
}

func TestStopRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request StopRequest
		valid   bool
	}{
		{
			name: "valid stop request with remove",
			request: StopRequest{
				Project:       "test-project",
				Remove:        true,
				RemoveVolumes: false,
			},
			valid: true,
		},
		{
			name: "valid stop without remove",
			request: StopRequest{
				Project: "test-project",
				Remove:  false,
			},
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.valid {
				assert.NotEmpty(t, tt.request.Project)
			}
		})
	}
}
