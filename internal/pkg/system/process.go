package system

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"

	"github.com/otto-nation/otto-stack/internal/core/docker"
)

// RunCommand executes a command and returns its output
func RunCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// RunCommandWithDir executes a command in a specific directory
func RunCommandWithDir(dir, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// RunCommandQuiet executes a command without capturing output
func RunCommandQuiet(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	return cmd.Run()
}

// IsCommandAvailable checks if a command is available in PATH
func IsCommandAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// GetProcessPID gets the PID of a running process by name
func GetProcessPID(name string) (int, error) {
	if name == "" {
		return 0, fmt.Errorf("process name cannot be empty")
	}

	cmd := createProcessCommand(name)
	if cmd == nil {
		return 0, fmt.Errorf(docker.ErrUnsupportedOS, runtime.GOOS)
	}

	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	return parsePIDFromOutput(string(output), name)
}

func createProcessCommand(name string) *exec.Cmd {
	switch runtime.GOOS {
	case docker.OSLinux, docker.OSDarwin:
		return exec.Command(docker.CmdPgrep, "-f", name)
	case docker.OSWindows:
		args := []string{"/FI", fmt.Sprintf("IMAGENAME eq %s.exe", name), "/FO", "CSV", "/NH"}
		return exec.Command(docker.CmdTasklist, args...)
	default:
		return nil
	}
}

func parsePIDFromOutput(output, processName string) (int, error) {
	outputStr := strings.TrimSpace(output)
	if outputStr == "" {
		return 0, fmt.Errorf(docker.ErrProcessNotFound, processName)
	}

	if runtime.GOOS != docker.OSWindows {
		return parseUnixPID(outputStr)
	}

	return parseWindowsPID(outputStr)
}

func parseUnixPID(output string) (int, error) {
	pidStr := strings.Split(output, "\n")[0]
	return strconv.Atoi(pidStr)
}

func parseWindowsPID(output string) (int, error) {
	lines := strings.Split(output, "\n")
	if len(lines) == 0 {
		return 0, pkgerrors.NewValidationError("input", "failed to parse PID from tasklist output", nil)
	}

	fields := strings.Split(lines[0], ",")
	if len(fields) < docker.MinFieldCount {
		return 0, pkgerrors.NewValidationError("input", "failed to parse PID from tasklist output", nil)
	}

	pidStr := strings.Trim(fields[1], "\"")
	return strconv.Atoi(pidStr)
}

// KillProcess kills a process by PID
func KillProcess(pid int) error {
	if runtime.GOOS == docker.OSWindows {
		return killWindowsProcess(pid)
	}
	return killUnixProcess(pid)
}

func killWindowsProcess(pid int) error {
	args := []string{"/F", "/PID", strconv.Itoa(pid)}
	cmd := exec.Command(docker.CmdTaskkill, args...)
	return cmd.Run()
}

func killUnixProcess(pid int) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return process.Signal(syscall.SIGTERM)
}

// Retry executes a function with retry logic
func Retry(attempts int, delay time.Duration, fn func() error) error {
	var err error
	for i := range attempts {
		if err = fn(); err == nil {
			return nil
		}
		if i < attempts-1 {
			time.Sleep(delay)
		}
	}
	return fmt.Errorf(docker.ErrFailedAfterRetry, attempts, err)
}

// Timeout executes a function with timeout
func Timeout(timeout time.Duration, fn func() error) error {
	done := make(chan error, 1)
	go func() {
		done <- fn()
	}()

	select {
	case err := <-done:
		return err
	case <-time.After(timeout):
		return fmt.Errorf(docker.ErrOperationTimeout, timeout)
	}
}
