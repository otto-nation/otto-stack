//go:build unit

package docker

import (
	"testing"

	"github.com/otto-nation/otto-stack/test/testhelpers"
)

func TestClient_GetResourceManager(t *testing.T) {
	logger := testhelpers.MockLogger()
	mockDocker := &testhelpers.MockDockerClient{}
	client := NewClientWithDependencies(mockDocker, nil, logger)

	rm := NewResourceManager(client)
	if rm == nil {
		t.Error("NewResourceManager should return non-nil manager")
	}

	if rm.client != client {
		t.Error("ResourceManager should reference the client")
	}
}
