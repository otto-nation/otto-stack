// Package testhelpers provides common test utilities to reduce redundancy
package testhelpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Common test constants used across multiple test files
const (
	TestProjectName = "test-project"
	TestServiceName = "test-service"
)

// AssertValidConstructor validates common constructor patterns
func AssertValidConstructor(t *testing.T, result any, err error, name string) {
	t.Helper()
	assert.NoError(t, err, "%s constructor should not return error", name)
	assert.NotNil(t, result, "%s constructor should return valid instance", name)
}

// AssertErrorPattern validates common error return patterns
func AssertErrorPattern(t *testing.T, result any, err error, expectError bool, operation string) {
	t.Helper()
	if expectError {
		assert.Error(t, err, "%s should return error", operation)
		assert.Nil(t, result, "%s should return nil result on error", operation)
	} else {
		assert.NoError(t, err, "%s should not return error", operation)
		assert.NotNil(t, result, "%s should return valid result", operation)
	}
}

// AssertHandlerInterface validates common handler interface implementations
func AssertHandlerInterface(t *testing.T, handler any) {
	t.Helper()
	assert.NotNil(t, handler, "Handler should not be nil")

	if h, ok := handler.(interface{ GetRequiredFlags() []string }); ok {
		flags := h.GetRequiredFlags()
		assert.NotNil(t, flags, "GetRequiredFlags should return valid slice")
	}
}
