package network

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/otto-nation/otto-stack/internal/core/docker"
)

// IsPortInUse checks if a port is in use
func IsPortInUse(port int) bool {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case docker.OSLinux, docker.OSDarwin:
		args := []string{"-i", fmt.Sprintf(":%d", port)}
		cmd = exec.Command(docker.CmdLsof, args...)
	case docker.OSWindows:
		cmd = exec.Command(docker.CmdNetstat, "-an")
	default:
		return false
	}

	output, err := cmd.Output()
	if err != nil {
		return false
	}

	if runtime.GOOS == docker.OSWindows {
		return strings.Contains(string(output), fmt.Sprintf(":%d", port))
	}

	return len(output) > 0
}

// GetFreePort finds an available port starting from the given port
func GetFreePort(startPort int) (int, error) {
	for port := startPort; port < startPort+docker.PortSearchRange; port++ {
		if !IsPortInUse(port) {
			return port, nil
		}
	}
	return 0, fmt.Errorf(docker.ErrNoFreePort, startPort, startPort+docker.PortSearchRange)
}
