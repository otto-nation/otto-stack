package docker

import (
	"testing"

	"github.com/otto-nation/otto-stack/test/testhelpers"
)

func TestClient_GetCliMethod(t *testing.T) {
	client := &Client{}

	// Test getter method
	cli := client.GetCli()

	// Should return the internal cli (even if nil)
	if cli != client.cli {
		t.Error("GetCli should return internal cli")
	}
}

func TestNewResourceManagerConstructor(t *testing.T) {
	client := &Client{}
	manager := NewResourceManager(client)
	testhelpers.AssertValidConstructor(t, manager, nil, "ResourceManager")
}

func TestContainerInfoBasic(t *testing.T) {
	info := &ContainerInfo{
		ID:     "test-id",
		Name:   "test-name",
		Status: "running",
	}

	// Basic validation
	if info.ID != "test-id" {
		t.Errorf("Expected ID test-id, got %s", info.ID)
	}
	if info.Name != "test-name" {
		t.Errorf("Expected Name test-name, got %s", info.Name)
	}
}
