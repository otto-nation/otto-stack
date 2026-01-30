package testutil

import (
	"net"
	"sync"

	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
)

type PortManager struct {
	allocatedPorts map[int]bool
	mu             sync.Mutex
}

func NewPortManager() *PortManager {
	return &PortManager{
		allocatedPorts: make(map[int]bool),
	}
}

func (pm *PortManager) AllocatePort() (int, error) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	defer func() { _ = listener.Close() }()

	port := listener.Addr().(*net.TCPAddr).Port
	pm.allocatedPorts[port] = true
	return port, nil
}

func (pm *PortManager) ReleasePort(port int) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	delete(pm.allocatedPorts, port)
}

func (pm *PortManager) AllocateServicePorts(services []string) (map[string]int, error) {
	ports := make(map[string]int)
	for _, service := range services {
		port, err := pm.AllocatePort()
		if err != nil {
			return nil, pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, "test", "allocate port", err)
		}
		ports[service] = port
	}
	return ports, nil
}
