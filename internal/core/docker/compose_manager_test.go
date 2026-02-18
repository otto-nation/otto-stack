package docker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleLogConsumer_Log(t *testing.T) {
	consumer := &SimpleLogConsumer{}
	consumer.Log("container1", "test message")
	// No assertion needed - just testing it doesn't panic
}

func TestSimpleLogConsumer_Err(t *testing.T) {
	consumer := &SimpleLogConsumer{}
	consumer.Err("container1", "error message")
	// No assertion needed - just testing it doesn't panic
}

func TestSimpleLogConsumer_Status(t *testing.T) {
	consumer := &SimpleLogConsumer{}
	consumer.Status("container1", "status message")
	// No assertion needed - just testing it doesn't panic
}

func TestManager_GetService(t *testing.T) {
	manager := &Manager{service: nil}
	service := manager.GetService()
	assert.Nil(t, service)
}
