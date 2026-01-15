//go:build unit

package docker

import (
	"testing"

	"github.com/otto-nation/otto-stack/test/testhelpers"
)

func TestNewClientWithDependencies_basic(t *testing.T) {
	logger := testhelpers.MockLogger()
	mockDocker := &testhelpers.MockDockerClient{}
	client := NewClientWithDependencies(mockDocker, nil, logger)

	if client == nil {
		t.Error("NewClientWithDependencies should return non-nil client")
	}

	if client.GetLogger() != logger {
		t.Error("Client should use provided logger")
	}

	if client.GetCli() != mockDocker {
		t.Error("Client should use provided Docker client")
	}
}

func TestClient_Close_basic(t *testing.T) {
	logger := testhelpers.MockLogger()
	mockDocker := &testhelpers.MockDockerClient{}
	client := NewClientWithDependencies(mockDocker, nil, logger)

	err := client.Close()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestClient_GetCli_basic(t *testing.T) {
	logger := testhelpers.MockLogger()
	mockDocker := &testhelpers.MockDockerClient{}
	client := NewClientWithDependencies(mockDocker, nil, logger)

	cli := client.GetCli()
	if cli != mockDocker {
		t.Error("Expected mock Docker client")
	}
}
