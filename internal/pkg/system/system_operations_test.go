//go:build unit

package system

import (
	"testing"

	"github.com/otto-nation/otto-stack/test/testhelpers"
)

func TestProcess_edge_cases(t *testing.T) {
	t.Run("get process PID with empty name", func(t *testing.T) {
		pid, err := GetProcessPID("")
		testhelpers.AssertError(t, err, "GetProcessPID with empty name should error")
		if pid != 0 {
			t.Error("GetProcessPID with empty name should return 0")
		}
	})

	t.Run("kill process with invalid PID", func(t *testing.T) {
		err := KillProcess(-1)
		testhelpers.AssertError(t, err, "KillProcess with invalid PID should error")
	})
}
