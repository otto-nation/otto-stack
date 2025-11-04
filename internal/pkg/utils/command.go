package utils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/otto-nation/otto-stack/internal/pkg/constants"
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
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case constants.OSLinux, constants.OSDarwin:
		cmd = exec.Command(constants.CmdPgrep, "-f", name)
	case constants.OSWindows:
		args := []string{"/FI", fmt.Sprintf("IMAGENAME eq %s.exe", name), "/FO", "CSV", "/NH"}
		cmd = exec.Command(constants.CmdTasklist, args...)
	default:
		return 0, fmt.Errorf(constants.ErrUnsupportedOS, runtime.GOOS)
	}

	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	outputStr := strings.TrimSpace(string(output))
	if outputStr == "" {
		return 0, fmt.Errorf(constants.ErrProcessNotFound, name)
	}

	if runtime.GOOS == constants.OSWindows {
		lines := strings.Split(outputStr, "\n")
		if len(lines) > 0 {
			fields := strings.Split(lines[0], ",")
			if len(fields) >= constants.MinFieldCount {
				pidStr := strings.Trim(fields[1], "\"")
				return strconv.Atoi(pidStr)
			}
		}
		return 0, fmt.Errorf("failed to parse PID from tasklist output")
	}

	pidStr := strings.Split(outputStr, "\n")[0]
	return strconv.Atoi(pidStr)
}

// KillProcess kills a process by PID
func KillProcess(pid int) error {
	if runtime.GOOS == constants.OSWindows {
		args := []string{"/F", "/PID", strconv.Itoa(pid)}
		cmd := exec.Command(constants.CmdTaskkill, args...)
		return cmd.Run()
	}

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
	return fmt.Errorf(constants.ErrFailedAfterRetry, attempts, err)
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
		return fmt.Errorf(constants.ErrOperationTimeout, timeout)
	}
}
