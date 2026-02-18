//go:build unit

package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServiceManager_Creates(t *testing.T) {
	cacheMutex.Lock()
	stackServiceCache = nil
	dockerManagerCache = nil
	resolverCache = nil
	cacheMutex.Unlock()

	service, err := NewServiceManager(false)
	if err != nil {
		t.Skipf("Skipping Docker test: %v", err)
	}
	assert.NotNil(t, service)
}

func TestNewServiceManager_Cached(t *testing.T) {
	service1, err1 := NewServiceManager(false)
	if err1 != nil {
		t.Skipf("Skipping Docker test: %v", err1)
	}

	service2, err2 := NewServiceManager(true)
	if err2 != nil {
		t.Skipf("Skipping Docker test: %v", err2)
	}

	assert.Same(t, service1, service2, "Should return cached instance")
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

func TestNewServiceManager_Debug(t *testing.T) {
	cacheMutex.Lock()
	stackServiceCache = nil
	cacheMutex.Unlock()

	_, err := NewServiceManager(true)
	if err != nil {
		t.Skipf("Skipping Docker test: %v", err)
	}
}

func TestGetDockerManager_Caching(t *testing.T) {
	cacheMutex.Lock()
	dockerManagerCache = nil
	cacheMutex.Unlock()

	m1, err := getDockerManager()
	if err != nil {
		t.Skipf("Skipping Docker test: %v", err)
	}

	m2, _ := getDockerManager()
	assert.Same(t, m1, m2)
}
