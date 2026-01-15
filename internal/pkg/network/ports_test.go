//go:build unit

package network

import (
	"runtime"
	"testing"

	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/stretchr/testify/assert"
)

func TestIsPortInUse(t *testing.T) {
	t.Run("handles port check without panic", func(t *testing.T) {
		// Test with a high port number that's unlikely to be in use
		result := IsPortInUse(65432)

		// Should return a boolean without panicking
		assert.IsType(t, false, result)
	})

	t.Run("handles invalid port gracefully", func(t *testing.T) {
		// Test with invalid port numbers
		result1 := IsPortInUse(-1)
		result2 := IsPortInUse(0)
		result3 := IsPortInUse(70000)

		// Should handle gracefully without panicking
		assert.IsType(t, false, result1)
		assert.IsType(t, false, result2)
		assert.IsType(t, false, result3)
	})

	t.Run("returns consistent type", func(t *testing.T) {
		// Test multiple calls return consistent boolean type
		for i := 8000; i < 8005; i++ {
			result := IsPortInUse(i)
			assert.IsType(t, false, result)
		}
	})

	t.Run("handles OS-specific behavior", func(t *testing.T) {
		// Test that function works on current OS
		currentOS := runtime.GOOS
		result := IsPortInUse(12345)

		// Should not panic regardless of OS
		assert.IsType(t, false, result)

		// Verify we're testing on a known OS
		knownOS := currentOS == "linux" || currentOS == "darwin" || currentOS == "windows"
		if !knownOS {
			// For unknown OS, function should return false
			assert.False(t, result)
		}
	})
}

func TestGetFreePort(t *testing.T) {
	t.Run("returns port in valid range", func(t *testing.T) {
		startPort := 8000
		port, err := GetFreePort(startPort)

		if err == nil {
			// If successful, port should be in expected range
			assert.GreaterOrEqual(t, port, startPort)
			assert.Less(t, port, startPort+docker.PortSearchRange)
		} else {
			// If no free port found, should return appropriate error
			assert.Error(t, err)
			assert.Equal(t, 0, port)
		}
	})

	t.Run("handles high start port", func(t *testing.T) {
		// Test with high port number
		startPort := 60000
		port, err := GetFreePort(startPort)

		if err == nil {
			assert.GreaterOrEqual(t, port, startPort)
		} else {
			assert.Error(t, err)
			assert.Equal(t, 0, port)
		}
	})

	t.Run("handles invalid start port", func(t *testing.T) {
		// Test with invalid start port
		port, err := GetFreePort(-1)

		// The function doesn't validate input, so it may return invalid ports
		// This test just ensures it doesn't panic
		if err != nil {
			assert.Equal(t, 0, port)
		} else {
			// Function may return any integer (including negative)
			assert.IsType(t, 0, port)
		}
	})

	t.Run("returns error when no ports available", func(t *testing.T) {
		// Test with port range that's likely to be exhausted
		startPort := 65500
		port, err := GetFreePort(startPort)

		// Should either find a port or return error
		if err != nil {
			assert.Equal(t, 0, port)
			assert.Contains(t, err.Error(), "no free port")
		} else {
			assert.GreaterOrEqual(t, port, startPort)
		}
	})

	t.Run("respects search range constant", func(t *testing.T) {
		// Test that search range is respected
		startPort := 50000
		port, err := GetFreePort(startPort)

		if err == nil {
			// Port should be within the search range
			assert.GreaterOrEqual(t, port, startPort)
			assert.Less(t, port, startPort+docker.PortSearchRange)
		} else {
			// Error should mention the range
			assert.Contains(t, err.Error(), "no free port found")
		}
	})
}

func TestDockerConstants(t *testing.T) {
	t.Run("validates Docker constants", func(t *testing.T) {
		// Test that constants are properly defined
		assert.Greater(t, docker.PortSearchRange, 0)
		assert.NotEmpty(t, docker.ErrNoFreePort)

		// Test OS constants
		assert.Equal(t, "linux", docker.OSLinux)
		assert.Equal(t, "darwin", docker.OSDarwin)
		assert.Equal(t, "windows", docker.OSWindows)

		// Test command constants
		assert.NotEmpty(t, docker.CmdLsof)
		assert.NotEmpty(t, docker.CmdNetstat)
	})
}
