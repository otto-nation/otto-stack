//go:build unit

package docker

import (
	"testing"

	"github.com/otto-nation/otto-stack/test/testhelpers"
)

func TestClient_GetComposeManager_Unit(t *testing.T) {
	mockDocker := &testhelpers.MockDockerClient{}
	mockCompose := &Manager{}

	client := NewClientWithDependencies(mockDocker, mockCompose, testhelpers.MockLogger())

	manager := client.GetComposeManager()
	if manager != mockCompose {
		t.Error("Expected compose manager to match")
	}
}
