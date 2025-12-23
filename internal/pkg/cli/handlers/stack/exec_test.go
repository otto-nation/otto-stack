package stack

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/stretchr/testify/assert"
)

func TestExecHandler_ValidateArgs(t *testing.T) {
	handler := NewExecHandler()

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "valid args with postgres service and psql command",
			args:    []string{services.ServicePostgres, services.ClientPsql, services.User_flagPostgres, services.DefaultUserPostgres},
			wantErr: false,
		},
		{
			name:    "valid args with redis service and redis-cli command",
			args:    []string{services.ServiceRedis, services.ClientRedisCli},
			wantErr: false,
		},
		{
			name:    "invalid - no args",
			args:    []string{},
			wantErr: true,
		},
		{
			name:    "invalid - only service name",
			args:    []string{services.ServicePostgres},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.ValidateArgs(tt.args)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestExecHandler_GetRequiredFlags(t *testing.T) {
	handler := NewExecHandler()
	flags := handler.GetRequiredFlags()
	assert.Empty(t, flags)
}

func TestNewExecHandler(t *testing.T) {
	handler := NewExecHandler()
	assert.NotNil(t, handler)
}
