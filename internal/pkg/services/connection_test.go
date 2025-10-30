package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetServiceConnectionConfig(t *testing.T) {
	tests := []struct {
		name           string
		serviceName    string
		expectedClient string
		expectedUser   string
		expectedPort   int
		expectError    bool
	}{
		{
			name:           "postgres service",
			serviceName:    "postgres",
			expectedClient: "psql",
			expectedUser:   "postgres",
			expectedPort:   5432,
			expectError:    false,
		},
		{
			name:           "redis service",
			serviceName:    "redis",
			expectedClient: "redis-cli",
			expectedPort:   6379,
			expectError:    false,
		},
		{
			name:        "unknown service",
			serviceName: "unknown",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := GetServiceConnectionConfig(tt.serviceName)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expectedClient, config.Client)
			if tt.expectedUser != "" {
				assert.Equal(t, tt.expectedUser, config.DefaultUser)
			}
			assert.Equal(t, tt.expectedPort, config.DefaultPort)
		})
	}
}
