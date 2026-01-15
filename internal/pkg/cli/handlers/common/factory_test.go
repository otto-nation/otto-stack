//go:build unit

package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServiceManager(t *testing.T) {
	// Reset caches for test isolation
	cacheMutex.Lock()
	stackServiceCache = nil
	dockerManagerCache = nil
	resolverCache = nil
	cacheMutex.Unlock()

	t.Run("creates service manager successfully", func(t *testing.T) {
		service, err := NewServiceManager(false)

		if err != nil {
			// Docker might not be available in test environment
			t.Skipf("Skipping Docker test: %v", err)
		}

		assert.NotNil(t, service)
	})

	t.Run("returns cached instance on second call", func(t *testing.T) {
		service1, err1 := NewServiceManager(false)
		if err1 != nil {
			t.Skipf("Skipping Docker test: %v", err1)
		}

		service2, err2 := NewServiceManager(true) // Different debug flag
		if err2 != nil {
			t.Skipf("Skipping Docker test: %v", err2)
		}

		assert.Same(t, service1, service2, "Should return cached instance")
	})
}

func TestGetDockerManager(t *testing.T) {
	// Reset cache for test
	cacheMutex.Lock()
	dockerManagerCache = nil
	cacheMutex.Unlock()

	manager1, err1 := getDockerManager()
	if err1 != nil {
		t.Skipf("Skipping Docker test: %v", err1)
	}

	manager2, err2 := getDockerManager()
	if err2 != nil {
		t.Skipf("Skipping Docker test: %v", err2)
	}

	assert.Same(t, manager1, manager2, "Should return cached instance")
}

func TestGetCharacteristicsResolver(t *testing.T) {
	// Reset cache for test
	cacheMutex.Lock()
	resolverCache = nil
	cacheMutex.Unlock()

	resolver1, err1 := getCharacteristicsResolver()
	assert.NoError(t, err1)
	assert.NotNil(t, resolver1)

	resolver2, err2 := getCharacteristicsResolver()
	assert.NoError(t, err2)
	assert.Same(t, resolver1, resolver2, "Should return cached instance")
}
