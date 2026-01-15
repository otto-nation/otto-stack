//go:build unit

package docker

import (
	"context"
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/logger"
)

func TestManager_GetService(t *testing.T) {
	client, err := NewClient(logger.GetLogger())
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer client.Close()

	manager := client.GetComposeManager()
	service := manager.GetService()
	if service == nil {
		t.Error("Expected service, got nil")
	}
}

func TestManager_LoadProject(t *testing.T) {
	client, err := NewClient(logger.GetLogger())
	if err != nil {
		t.Skipf("Docker not available: %v", err)
	}
	defer client.Close()

	manager := client.GetComposeManager()
	ctx := context.Background()

	project, err := manager.LoadProject(ctx, "non-existent-compose.yml", "test-project")
	// Expect error for non-existent file
	if err == nil {
		t.Error("Expected error for non-existent compose file")
	}
	_ = project
}

func TestSimpleLogConsumer_Log(t *testing.T) {
	consumer := &SimpleLogConsumer{}
	consumer.Log("test-container", "Test log message")
	// Should not panic
}

func TestSimpleLogConsumer_Err(t *testing.T) {
	consumer := &SimpleLogConsumer{}
	consumer.Err("test-container", "Test error message")
	// Should not panic
}

func TestSimpleLogConsumer_Status(t *testing.T) {
	consumer := &SimpleLogConsumer{}
	consumer.Status("test-container", "Test status message")
	// Should not panic
}
