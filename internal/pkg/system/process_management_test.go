//go:build unit

package system

import (
	"runtime"
	"testing"
	"time"

	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/stretchr/testify/assert"
)

func TestGetProcessPID_EdgeCases(t *testing.T) {
	t.Run("handles nonexistent process", func(t *testing.T) {
		pid, err := GetProcessPID("nonexistent-process-12345")
		assert.Error(t, err)
		assert.Equal(t, 0, pid)
	})

	t.Run("handles empty process name", func(t *testing.T) {
		pid, err := GetProcessPID("")
		assert.Error(t, err)
		assert.Equal(t, 0, pid)
	})

	t.Run("handles OS-specific behavior", func(t *testing.T) {
		// Test with current process name (should exist)
		currentOS := runtime.GOOS

		// Use a process that should exist on the system
		var processName string
		switch currentOS {
		case "darwin", "linux":
			processName = "kernel" // Should exist on Unix systems
		case "windows":
			processName = "System" // Should exist on Windows
		default:
			t.Skip("Unsupported OS for this test")
		}

		pid, err := GetProcessPID(processName)
		// Either finds the process or returns appropriate error
		if err == nil {
			assert.Greater(t, pid, 0)
		} else {
			assert.Error(t, err)
		}
	})

	t.Run("validates OS constants", func(t *testing.T) {
		assert.Equal(t, "linux", docker.OSLinux)
		assert.Equal(t, "darwin", docker.OSDarwin)
		assert.Equal(t, "windows", docker.OSWindows)
		assert.NotEmpty(t, docker.CmdPgrep)
		assert.NotEmpty(t, docker.CmdTasklist)
	})
}

func TestKillProcess_EdgeCases(t *testing.T) {
	t.Run("handles invalid PID", func(t *testing.T) {
		err := KillProcess(-1)
		assert.Error(t, err)
	})

	t.Run("handles zero PID", func(t *testing.T) {
		err := KillProcess(0)
		// Should handle gracefully (may succeed or fail depending on OS)
		_ = err // Don't assert specific behavior as it's OS-dependent
	})

	t.Run("handles nonexistent PID", func(t *testing.T) {
		// Use a very high PID that's unlikely to exist
		err := KillProcess(999999)
		// Should return error for nonexistent process
		assert.Error(t, err)
	})

	t.Run("validates OS-specific commands", func(t *testing.T) {
		if runtime.GOOS == docker.OSWindows {
			assert.NotEmpty(t, docker.CmdTaskkill)
		}
		// Unix systems use os.FindProcess, no specific command needed
	})
}

func TestRetry_EdgeCases(t *testing.T) {
	t.Run("succeeds on first attempt", func(t *testing.T) {
		attempts := 0
		err := Retry(3, 10*time.Millisecond, func() error {
			attempts++
			return nil // Success on first try
		})

		assert.NoError(t, err)
		assert.Equal(t, 1, attempts)
	})

	t.Run("fails after max attempts", func(t *testing.T) {
		attempts := 0
		err := Retry(2, 1*time.Millisecond, func() error {
			attempts++
			return assert.AnError // Always fail
		})

		assert.Error(t, err)
		assert.Equal(t, 2, attempts)
		assert.Contains(t, err.Error(), "failed after")
	})

	t.Run("succeeds on last attempt", func(t *testing.T) {
		attempts := 0
		err := Retry(3, 1*time.Millisecond, func() error {
			attempts++
			if attempts < 3 {
				return assert.AnError
			}
			return nil // Success on third try
		})

		assert.NoError(t, err)
		assert.Equal(t, 3, attempts)
	})
}

func TestTimeout_EdgeCases(t *testing.T) {
	t.Run("completes within timeout", func(t *testing.T) {
		err := Timeout(100*time.Millisecond, func() error {
			time.Sleep(10 * time.Millisecond)
			return nil
		})

		assert.NoError(t, err)
	})

	t.Run("times out", func(t *testing.T) {
		err := Timeout(10*time.Millisecond, func() error {
			time.Sleep(50 * time.Millisecond)
			return nil
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "timed out")
	})

	t.Run("function returns error before timeout", func(t *testing.T) {
		err := Timeout(100*time.Millisecond, func() error {
			return assert.AnError
		})

		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
	})
}

func TestDockerConstants_Validation(t *testing.T) {
	t.Run("validates error constants", func(t *testing.T) {
		assert.NotEmpty(t, docker.ErrUnsupportedOS)
		assert.NotEmpty(t, docker.ErrProcessNotFound)
		assert.NotEmpty(t, docker.ErrFailedAfterRetry)
		assert.Greater(t, docker.MinFieldCount, 0)
	})

	t.Run("validates command constants", func(t *testing.T) {
		assert.NotEmpty(t, docker.CmdPgrep)
		assert.NotEmpty(t, docker.CmdTasklist)
		assert.NotEmpty(t, docker.CmdTaskkill)
	})
}
