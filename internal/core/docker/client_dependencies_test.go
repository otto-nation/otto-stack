package docker

import (
	"testing"

	"github.com/otto-nation/otto-stack/test/testhelpers"
)

func TestNewClientWithDependencies_basic(t *testing.T) {
	logger := testhelpers.MockLogger()
	client := NewClientWithDependencies(nil, nil, logger)

	if client == nil {
		t.Error("NewClientWithDependencies should return non-nil client")
	}

	if client.logger != logger {
		t.Error("Client should use provided logger")
	}
}

func TestClient_Close_basic(t *testing.T) {
	logger := testhelpers.MockLogger()
	client := NewClientWithDependencies(nil, nil, logger)

	// This will panic with nil cli, but that's expected behavior
	// We're testing the code path exists
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic with nil client")
		}
	}()

	client.Close()
}

func TestClient_GetCli_basic(t *testing.T) {
	logger := testhelpers.MockLogger()
	client := NewClientWithDependencies(nil, nil, logger)

	cli := client.GetCli()
	if cli != nil {
		t.Error("Expected nil cli from test client")
	}
}
