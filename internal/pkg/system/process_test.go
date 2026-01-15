package system

import (
	"errors"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRunCommand(t *testing.T) {
	t.Run("executes simple command successfully", func(t *testing.T) {
		var cmd, arg string
		if runtime.GOOS == "windows" {
			cmd, arg = "cmd", "/c echo test"
		} else {
			cmd, arg = "echo", "test"
		}

		output, err := RunCommand(cmd, arg)
		assert.NoError(t, err)
		assert.Contains(t, output, "test")
	})

	t.Run("returns error for invalid command", func(t *testing.T) {
		_, err := RunCommand("nonexistent-command-12345")
		assert.Error(t, err)
	})
}

func TestRunCommandWithDir(t *testing.T) {
	t.Run("executes command in specified directory", func(t *testing.T) {
		tempDir := t.TempDir()

		var cmd string
		var args []string
		if runtime.GOOS == "windows" {
			cmd = "cmd"
			args = []string{"/c", "cd"}
		} else {
			cmd = "ls"
			args = []string{"-la"}
		}

		output, err := RunCommandWithDir(tempDir, cmd, args...)
		assert.NoError(t, err)
		assert.NotEmpty(t, output)
	})

	t.Run("returns error for invalid directory", func(t *testing.T) {
		var cmd, arg string
		if runtime.GOOS == "windows" {
			cmd, arg = "cmd", "/c echo test"
		} else {
			cmd, arg = "echo", "test"
		}

		_, err := RunCommandWithDir("/nonexistent/directory", cmd, arg)
		assert.Error(t, err)
	})
}

func TestRunCommandQuiet(t *testing.T) {
	t.Run("executes command without output", func(t *testing.T) {
		var cmd, arg string
		if runtime.GOOS == "windows" {
			cmd, arg = "cmd", "/c echo test"
		} else {
			cmd, arg = "echo", "test"
		}

		err := RunCommandQuiet(cmd, arg)
		assert.NoError(t, err)
	})

	t.Run("returns error for invalid command", func(t *testing.T) {
		err := RunCommandQuiet("nonexistent-command-12345")
		assert.Error(t, err)
	})
}

func TestIsCommandAvailable(t *testing.T) {
	t.Run("returns true for available command", func(t *testing.T) {
		var cmd string
		if runtime.GOOS == "windows" {
			cmd = "cmd"
		} else {
			cmd = "echo"
		}

		available := IsCommandAvailable(cmd)
		assert.True(t, available)
	})

	t.Run("returns false for unavailable command", func(t *testing.T) {
		available := IsCommandAvailable("nonexistent-command-12345")
		assert.False(t, available)
	})
}

func TestGetProcessPID(t *testing.T) {
	t.Run("handles nonexistent process", func(t *testing.T) {
		pid, err := GetProcessPID("nonexistent-process-12345")
		assert.Error(t, err)
		assert.Equal(t, 0, pid)
	})

	// Note: Testing actual process lookup is platform-dependent and unreliable
	// in CI environments, so we focus on error handling
}

func TestKillProcess(t *testing.T) {
	t.Run("handles invalid PID gracefully", func(t *testing.T) {
		// Use a PID that's very unlikely to exist
		err := KillProcess(999999)
		// Should return an error (process not found) rather than panic
		assert.Error(t, err)
	})
}

func TestRetry(t *testing.T) {
	t.Run("succeeds on first attempt", func(t *testing.T) {
		attempts := 0
		fn := func() error {
			attempts++
			return nil
		}

		err := Retry(3, time.Millisecond, fn)
		assert.NoError(t, err)
		assert.Equal(t, 1, attempts)
	})

	t.Run("succeeds after retries", func(t *testing.T) {
		attempts := 0
		fn := func() error {
			attempts++
			if attempts < 3 {
				return errors.New("temporary error")
			}
			return nil
		}

		err := Retry(3, time.Millisecond, fn)
		assert.NoError(t, err)
		assert.Equal(t, 3, attempts)
	})

	t.Run("fails after max attempts", func(t *testing.T) {
		attempts := 0
		fn := func() error {
			attempts++
			return errors.New("persistent error")
		}

		err := Retry(2, time.Millisecond, fn)
		assert.Error(t, err)
		assert.Equal(t, 2, attempts)
		assert.Contains(t, err.Error(), "failed after")
	})
}

func TestTimeout(t *testing.T) {
	t.Run("completes within timeout", func(t *testing.T) {
		fn := func() error {
			time.Sleep(10 * time.Millisecond)
			return nil
		}

		err := Timeout(100*time.Millisecond, fn)
		assert.NoError(t, err)
	})

	t.Run("times out when function takes too long", func(t *testing.T) {
		fn := func() error {
			time.Sleep(100 * time.Millisecond)
			return nil
		}

		err := Timeout(10*time.Millisecond, fn)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "timed out")
	})

	t.Run("returns function error when completed", func(t *testing.T) {
		expectedErr := errors.New("function error")
		fn := func() error {
			return expectedErr
		}

		err := Timeout(100*time.Millisecond, fn)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}
