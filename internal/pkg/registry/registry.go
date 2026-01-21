package registry

import (
	"os"
	"path/filepath"
	"slices"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/otto-nation/otto-stack/internal/core"
)

// Manager handles shared container registry operations
type Manager struct {
	registryPath string
}

// NewManager creates a new registry manager
func NewManager(sharedDir string) *Manager {
	return &Manager{
		registryPath: filepath.Join(sharedDir, core.SharedRegistryFile),
	}
}

// Load reads the registry from disk
func (m *Manager) Load() (*Registry, error) {
	data, err := os.ReadFile(m.registryPath)
	if os.IsNotExist(err) {
		return NewRegistry(), nil
	}
	if err != nil {
		return nil, err
	}

	var registry Registry
	if err := yaml.Unmarshal(data, &registry); err != nil {
		return nil, err
	}

	if registry.Containers == nil {
		registry.Containers = make(map[string]*ContainerInfo)
	}

	return &registry, nil
}

// Save writes the registry to disk
func (m *Manager) Save(registry *Registry) error {
	data, err := yaml.Marshal(registry)
	if err != nil {
		return err
	}

	return os.WriteFile(m.registryPath, data, core.PermReadWrite)
}

// Register adds or updates a container in the registry
func (m *Manager) Register(service, containerName, projectName string) error {
	registry, err := m.Load()
	if err != nil {
		return err
	}

	now := time.Now()
	container, exists := registry.Containers[service]

	if !exists {
		container = &ContainerInfo{
			Name:      containerName,
			Service:   service,
			Projects:  []string{projectName},
			CreatedAt: now,
			UpdatedAt: now,
		}
		registry.Containers[service] = container
	} else {
		if !slices.Contains(container.Projects, projectName) {
			container.Projects = append(container.Projects, projectName)
		}
		container.UpdatedAt = now
	}

	return m.Save(registry)
}

// Unregister removes a project from a container's usage list
func (m *Manager) Unregister(service, projectName string) error {
	registry, err := m.Load()
	if err != nil {
		return err
	}

	container, exists := registry.Containers[service]
	if !exists {
		return nil
	}

	container.Projects = remove(container.Projects, projectName)
	container.UpdatedAt = time.Now()

	if len(container.Projects) == 0 {
		delete(registry.Containers, service)
	}

	return m.Save(registry)
}

// Get retrieves container info for a service
func (m *Manager) Get(service string) (*ContainerInfo, error) {
	registry, err := m.Load()
	if err != nil {
		return nil, err
	}

	return registry.Containers[service], nil
}

// List returns all registered containers
func (m *Manager) List() ([]*ContainerInfo, error) {
	registry, err := m.Load()
	if err != nil {
		return nil, err
	}

	containers := make([]*ContainerInfo, 0, len(registry.Containers))
	for _, container := range registry.Containers {
		containers = append(containers, container)
	}

	return containers, nil
}

// IsShared checks if a service has a shared container
func (m *Manager) IsShared(service string) (bool, error) {
	container, err := m.Get(service)
	if err != nil {
		return false, err
	}
	return container != nil, nil
}

// remove removes a value from a string slice
func remove(slice []string, value string) []string {
	result := make([]string, 0, len(slice))
	for _, item := range slice {
		if item != value {
			result = append(result, item)
		}
	}
	return result
}
