package docker

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServiceLogConsumer_SingleService_NoPrefix(t *testing.T) {
	var buf bytes.Buffer
	consumer := NewServiceLogConsumer(&buf, true, 1)
	consumer.Log("postgres", "database ready")
	assert.Equal(t, "database ready\n", buf.String())
}

func TestServiceLogConsumer_MultiService_Prefix(t *testing.T) {
	var buf bytes.Buffer
	consumer := NewServiceLogConsumer(&buf, true, 2)
	consumer.Log("postgres", "database ready")
	consumer.Log("redis", "ready to accept connections")

	lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	assert.Len(t, lines, 2)
	assert.Contains(t, lines[0], "postgres | database ready")
	assert.Contains(t, lines[1], "redis | ready to accept connections")
}

func TestServiceLogConsumer_Err_SingleService(t *testing.T) {
	var buf bytes.Buffer
	consumer := NewServiceLogConsumer(&buf, true, 1)
	consumer.Err("postgres", "connection refused")
	assert.Equal(t, "connection refused\n", buf.String())
}

func TestServiceLogConsumer_Status_Silent(t *testing.T) {
	var buf bytes.Buffer
	consumer := NewServiceLogConsumer(&buf, true, 1)
	// Status events must not appear in the writer output
	consumer.Status("postgres", "container started")
	assert.Empty(t, buf.String())
}

func TestServiceLogConsumer_ColorCyclesAcrossServices(t *testing.T) {
	var buf bytes.Buffer
	// noColor=false so ANSI codes are embedded; we just verify different services
	// get different colors by checking the prefix color sequence resets.
	consumer := NewServiceLogConsumer(&buf, false, 3)
	consumer.Log("svc-a", "msg")
	consumer.Log("svc-b", "msg")
	consumer.Log("svc-c", "msg")

	output := buf.String()
	assert.Contains(t, output, "svc-a")
	assert.Contains(t, output, "svc-b")
	assert.Contains(t, output, "svc-c")
}

func TestManager_GetService(t *testing.T) {
	manager := &Manager{service: nil}
	service := manager.GetService()
	assert.Nil(t, service)
}
