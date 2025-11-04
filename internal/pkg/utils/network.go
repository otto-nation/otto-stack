package utils

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/otto-nation/otto-stack/internal/pkg/constants"
)

// IsPortInUse checks if a port is in use
func IsPortInUse(port int) bool {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case OSLinux, OSDarwin:
		args := []string{"-i", fmt.Sprintf(":%d", port)}
		cmd = exec.Command("lsof", args...)
	case OSWindows:
		cmd = exec.Command("netstat", "-an")
	default:
		return false
	}

	output, err := cmd.Output()
	if err != nil {
		return false
	}

	if runtime.GOOS == OSWindows {
		return strings.Contains(string(output), fmt.Sprintf(":%d", port))
	}

	return len(output) > 0
}

// GetFreePort finds an available port starting from the given port
func GetFreePort(startPort int) (int, error) {
	for port := startPort; port < startPort+constants.PortSearchRange; port++ {
		if !IsPortInUse(port) {
			return port, nil
		}
	}
	return 0, fmt.Errorf("no free port found in range %d-%d", startPort, startPort+constants.PortSearchRange)
}
