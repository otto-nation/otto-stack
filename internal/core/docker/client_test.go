package docker

import (
	"log/slog"
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/logger"
	"github.com/otto-nation/otto-stack/test/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name        string
		logger      *slog.Logger
		expectError bool
	}{
		{
			name:        "create client with valid logger",
			logger:      logger.GetLogger(),
			expectError: false,
		},
		{
			name:        "create client with nil logger",
			logger:      nil,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.logger)

			if tt.expectError {
				testhelpers.AssertErrorPattern(t, client, err, true, "NewClient")
			} else {
				if err != nil {
					// Docker might not be available in test environment
					t.Skipf("Docker not available: %v", err)
				}
				require.NotNil(t, client)
				assert.NotNil(t, client.GetCli())

				// Clean up
				_ = client.Close()
			}
		})
	}
}

func TestClient_Close_Once(t *testing.T) {
	testLogger := logger.GetLogger()

	client, err := NewClient(testLogger)
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}

	err = client.Close()
	assert.NoError(t, err)
}

func TestClient_Close_Twice(t *testing.T) {
	testLogger := logger.GetLogger()

	client, err := NewClient(testLogger)
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}

	_ = client.Close()
	assert.NotPanics(t, func() {
		_ = client.Close()
	})
}

func TestClient_GetCli(t *testing.T) {
	testLogger := logger.GetLogger()

	client, err := NewClient(testLogger)
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer func() { _ = client.Close() }()

	assert.NotNil(t, client.GetCli())
}

func TestClient_WithNilLogger(t *testing.T) {
	client, err := NewClient(nil)
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer func() {
		if client != nil {
			_ = client.Close()
		}
	}()

	assert.NotNil(t, client)
	assert.NotNil(t, client.GetCli())
}
