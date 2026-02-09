package registry

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"slices"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/core/docker"
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
	// Open file for reading with lock
	f, err := os.OpenFile(m.registryPath, os.O_RDONLY|os.O_CREATE, core.PermReadWrite)
	if err != nil {
		if os.IsNotExist(err) {
			return NewRegistry(), nil
		}
		return nil, err
	}
	defer func() { _ = f.Close() }()

	// Acquire shared lock for reading
	if err := lockFile(f); err != nil {
		return nil, err
	}
	defer func() { _ = unlockFile(f) }()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return NewRegistry(), nil
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

	// Ensure directory exists
	dir := filepath.Dir(m.registryPath)
	if err := os.MkdirAll(dir, core.PermReadWriteExec); err != nil {
		return err
	}

	// Atomic write: write to temp file, then rename
	tempPath := m.registryPath + ".tmp"
	f, err := os.OpenFile(tempPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, core.PermReadWrite)
	if err != nil {
		return err
	}
	defer func() { _ = os.Remove(tempPath) }() // Clean up temp file on error

	// Acquire exclusive lock
	if err := lockFile(f); err != nil {
		_ = f.Close()
		return err
	}

	// Write data
	if _, err := f.Write(data); err != nil {
		_ = unlockFile(f)
		_ = f.Close()
		return err
	}

	// Sync to disk
	if err := f.Sync(); err != nil {
		_ = unlockFile(f)
		_ = f.Close()
		return err
	}

	_ = unlockFile(f)
	_ = f.Close()

	// Atomic rename
	return os.Rename(tempPath, m.registryPath)
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

// FindOrphans returns containers with no active projects (basic check)
func (m *Manager) FindOrphans() ([]OrphanInfo, error) {
	registry, err := m.Load()
	if err != nil {
		return nil, err
	}
	return registry.FindOrphans(), nil
}

// FindOrphansWithChecks performs enhanced orphan detection with filesystem and Docker checks
func (m *Manager) FindOrphansWithChecks(ctx context.Context, dockerClient *docker.Client) ([]OrphanInfo, error) {
	registry, err := m.Load()
	if err != nil {
		return nil, err
	}

	var orphans []OrphanInfo

	// Get all running containers
	containers, err := dockerClient.ListContainers(ctx, "")
	if err != nil {
		return nil, err
	}

	containerMap := make(map[string]bool)
	for _, c := range containers {
		containerMap[c.Name] = true
	}

	// Check registered containers
	for service, info := range registry.Containers {
		orphan := OrphanInfo{
			Service:   service,
			Container: info.Name,
		}

		// Check if container exists in Docker
		if !containerMap[info.Name] {
			orphan.Reason = "container not found in Docker"
			orphan.Severity = OrphanSeverityWarning
			orphan.ContainerState = ContainerStateNotFound
			orphans = append(orphans, orphan)
			continue
		}

		// Check if projects still exist
		var existingProjects []string
		for _, project := range info.Projects {
			if projectExists(project) {
				existingProjects = append(existingProjects, project)
			}
		}

		// No projects using it
		if len(info.Projects) == 0 {
			orphan.Reason = "no projects using this container"
			orphan.Severity = OrphanSeveritySafe
			orphan.ContainerState = ContainerStateRunning
			orphans = append(orphans, orphan)
		} else if len(existingProjects) == 0 {
			// Projects listed but none exist on filesystem
			orphan.Reason = "all projects deleted from filesystem"
			orphan.Severity = OrphanSeveritySafe
			orphan.ContainerState = ContainerStateRunning
			orphan.ProjectsFound = info.Projects
			orphans = append(orphans, orphan)
		} else if len(existingProjects) < len(info.Projects) {
			// Some projects deleted
			orphan.Reason = "some projects deleted from filesystem"
			orphan.Severity = OrphanSeverityWarning
			orphan.ContainerState = ContainerStateRunning
			orphan.ProjectsFound = existingProjects
			orphans = append(orphans, orphan)
		}
	}

	// Check for zombie containers (running but not registered)
	for _, container := range containers {
		if isSharedContainer(container.Name) {
			service := extractServiceName(container.Name)
			if _, exists := registry.Containers[service]; !exists {
				orphans = append(orphans, OrphanInfo{
					Service:        service,
					Container:      container.Name,
					Reason:         "container running but not registered",
					Severity:       OrphanSeverityCritical,
					ContainerState: ContainerStateRunning,
				})
			}
		}
	}

	return orphans, nil
}

// projectExists checks if a project directory exists
func projectExists(projectName string) bool {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	// Check common project locations
	// This is a heuristic - projects could be anywhere
	// We check if .otto-stack directory exists in common locations
	possiblePaths := []string{
		filepath.Join(homeDir, "projects", projectName, core.OttoStackDir),
		filepath.Join(homeDir, "dev", projectName, core.OttoStackDir),
		filepath.Join(homeDir, projectName, core.OttoStackDir),
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return true
		}
	}

	return false
}

// isSharedContainer checks if container name matches shared pattern
func isSharedContainer(name string) bool {
	return len(name) > len(core.SharedContainerPrefix) &&
		name[:len(core.SharedContainerPrefix)] == core.SharedContainerPrefix
}

// extractServiceName extracts service name from container name
func extractServiceName(containerName string) string {
	if !isSharedContainer(containerName) {
		return ""
	}
	return containerName[len(core.SharedContainerPrefix):]
}

// CleanOrphans removes orphaned containers from registry
func (m *Manager) CleanOrphans() ([]string, error) {
	registry, err := m.Load()
	if err != nil {
		return nil, err
	}

	var cleaned []string
	for service, info := range registry.Containers {
		if len(info.Projects) == 0 {
			delete(registry.Containers, service)
			cleaned = append(cleaned, service)
		}
	}

	if len(cleaned) > 0 {
		if err := m.Save(registry); err != nil {
			return nil, err
		}
	}

	return cleaned, nil
}

// Reconcile syncs registry with actual Docker container state
func (m *Manager) Reconcile(ctx context.Context, dockerClient *docker.Client) (*ReconcileResult, error) {
	registry, err := m.Load()
	if err != nil {
		return nil, err
	}

	result := &ReconcileResult{
		Removed: []string{},
		Added:   []string{},
	}

	// Get all containers with shared prefix
	containers, err := dockerClient.ListContainers(ctx, "")
	if err != nil {
		return nil, err
	}

	// Build map of existing shared containers
	existingContainers := make(map[string]bool)
	for _, cont := range containers {
		if len(cont.Name) > len(core.SharedContainerPrefix) &&
			cont.Name[:len(core.SharedContainerPrefix)] == core.SharedContainerPrefix {
			serviceName := cont.Name[len(core.SharedContainerPrefix):]
			existingContainers[serviceName] = true
		}
	}

	// Remove registry entries for containers that don't exist
	for service := range registry.Containers {
		if !existingContainers[service] {
			delete(registry.Containers, service)
			result.Removed = append(result.Removed, service)
		}
	}

	// Save if changes were made
	if len(result.Removed) > 0 {
		if err := m.Save(registry); err != nil {
			return nil, err
		}
	}

	return result, nil
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
