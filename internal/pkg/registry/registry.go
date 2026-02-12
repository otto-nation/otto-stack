package registry

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/messages"
)

// Manager handles shared container registry operations
type Manager struct {
	registryPath   string
	orphanDetector *OrphanDetector
}

// NewManager creates a new registry manager
func NewManager(sharedDir string) *Manager {
	m := &Manager{
		registryPath: filepath.Join(sharedDir, core.SharedRegistryFile),
	}
	m.orphanDetector = NewOrphanDetector(m)
	return m
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
		return pkgerrors.NewSystemError(pkgerrors.ErrCodeOperationFail, messages.ErrorsRegistryMarshalFailed, err)
	}

	// Prepend header comment
	dataWithHeader := append([]byte(core.RegistryHeader), data...)

	// Ensure directory exists
	dir := filepath.Dir(m.registryPath)
	if err := os.MkdirAll(dir, core.PermReadWriteExec); err != nil {
		return pkgerrors.NewSystemError(pkgerrors.ErrCodeOperationFail, messages.ErrorsDirectoryCreateFailed, err)
	}

	// Atomic write: write to temp file, then rename
	tempPath := m.registryPath + ".tmp"
	f, err := os.OpenFile(tempPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, core.PermReadWrite)
	if err != nil {
		return pkgerrors.NewSystemError(pkgerrors.ErrCodeOperationFail, messages.ErrorsFileWriteFailed, err)
	}
	defer func() { _ = os.Remove(tempPath) }() // Clean up temp file on error

	// Acquire exclusive lock
	if err := lockFile(f); err != nil {
		_ = f.Close()
		return pkgerrors.NewSystemError(pkgerrors.ErrCodeOperationFail, messages.ErrorsRegistryLockFailed, err)
	}

	// Write data
	if _, err := f.Write(dataWithHeader); err != nil {
		_ = unlockFile(f)
		_ = f.Close()
		return pkgerrors.NewSystemError(pkgerrors.ErrCodeOperationFail, messages.ErrorsFileWriteFailed, err)
	}

	// Sync to disk
	if err := f.Sync(); err != nil {
		_ = unlockFile(f)
		_ = f.Close()
		return pkgerrors.NewSystemError(pkgerrors.ErrCodeOperationFail, messages.ErrorsRegistrySyncFailed, err)
	}

	_ = unlockFile(f)
	_ = f.Close()

	// Atomic rename
	if err := os.Rename(tempPath, m.registryPath); err != nil {
		return pkgerrors.NewSystemError(pkgerrors.ErrCodeOperationFail, messages.ErrorsRegistrySaveFailed, err)
	}
	return nil
}

// Register adds or updates a container in the registry
func (m *Manager) Register(service, containerName, projectName string) error {
	registry, err := m.Load()
	if err != nil {
		return pkgerrors.NewSystemError(pkgerrors.ErrCodeOperationFail, messages.ErrorsRegistryLoadFailed, err)
	}

	m.cleanupOrphans(registry)

	now := time.Now()
	container, exists := registry.Containers[service]

	if !exists {
		container = &ContainerInfo{
			Name:      containerName,
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

	if err := m.Save(registry); err != nil {
		return err
	}

	return m.createOrUpdateSharedReadme(registry)
}

// Unregister removes a project from a container's usage list
func (m *Manager) Unregister(service, projectName string) error {
	registry, err := m.Load()
	if err != nil {
		return pkgerrors.NewSystemError(pkgerrors.ErrCodeOperationFail, messages.ErrorsRegistryLoadFailed, err)
	}

	m.cleanupOrphans(registry)

	container, exists := registry.Containers[service]
	if !exists {
		return nil
	}

	container.Projects = remove(container.Projects, projectName)
	container.UpdatedAt = time.Now()

	if len(container.Projects) == 0 {
		delete(registry.Containers, service)
	}

	if err := m.Save(registry); err != nil {
		return err
	}

	return m.createOrUpdateSharedReadme(registry)
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
func (m *Manager) List() (map[string]*ContainerInfo, error) {
	registry, err := m.Load()
	if err != nil {
		return nil, err
	}

	return registry.Containers, nil
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
	return m.orphanDetector.FindOrphans()
}

// FindOrphansWithChecks performs enhanced orphan detection with filesystem and Docker checks
func (m *Manager) FindOrphansWithChecks(ctx context.Context, dockerClient *docker.Client) ([]OrphanInfo, error) {
	return m.orphanDetector.FindOrphansWithChecks(ctx, dockerClient)
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

// cleanupOrphans removes projects that no longer exist from the registry
func (m *Manager) cleanupOrphans(registry *Registry) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return // Can't determine home, skip cleanup
	}

	expectedSharedPath := filepath.Join(homeDir, core.OttoStackDir, core.SharedDir)
	actualSharedPath := filepath.Dir(m.registryPath)

	// Only run cleanup if we're in the real shared directory
	if actualSharedPath != expectedSharedPath {
		return // Skip cleanup in test environments
	}

	for service, container := range registry.Containers {
		validProjects := make([]string, 0, len(container.Projects))
		for _, project := range container.Projects {
			if m.projectExists(project, homeDir) {
				validProjects = append(validProjects, project)
			}
		}
		container.Projects = validProjects
		if len(container.Projects) == 0 {
			delete(registry.Containers, service)
		}
	}
}

// projectExists checks if a project directory exists
func (m *Manager) projectExists(projectName, homeDir string) bool {
	// Check common project locations relative to home
	searchPaths := []string{
		filepath.Join(homeDir, "git", "*", projectName, core.OttoStackDir),
		filepath.Join(homeDir, "projects", projectName, core.OttoStackDir),
		filepath.Join(homeDir, projectName, core.OttoStackDir),
	}

	for _, pattern := range searchPaths {
		matches, _ := filepath.Glob(pattern)
		if len(matches) > 0 {
			return true
		}
	}

	return false
}

// createOrUpdateSharedReadme generates/updates the README in the shared directory
func (m *Manager) createOrUpdateSharedReadme(registry *Registry) error {
	sharedDir := filepath.Dir(m.registryPath)
	readmePath := filepath.Join(sharedDir, "README.md")

	content := m.buildReadmeContent(registry)
	return os.WriteFile(readmePath, []byte(content), core.PermReadWrite)
}

// buildReadmeContent generates the README content
func (m *Manager) buildReadmeContent(registry *Registry) string {
	var b strings.Builder

	b.WriteString("# Otto Stack - Shared Services\n\n")
	b.WriteString("This directory manages shared containers across multiple otto-stack projects.\n\n")
	b.WriteString("## Files\n\n")
	b.WriteString("- `containers.yaml` - Registry tracking which projects use which shared containers\n")
	b.WriteString("- `docker-compose.yml` - Generated compose file for shared containers (created on demand)\n\n")
	b.WriteString("## Active Shared Containers\n\n")

	if len(registry.Containers) == 0 {
		b.WriteString("No shared containers currently registered.\n\n")
	} else {
		for service, info := range registry.Containers {
			b.WriteString(fmt.Sprintf("### %s\n", service))
			b.WriteString(fmt.Sprintf("- Container: `%s`\n", info.Name))
			b.WriteString(fmt.Sprintf("- Projects: %s\n", strings.Join(info.Projects, ", ")))
			b.WriteString(fmt.Sprintf("- Created: %s\n", info.CreatedAt.Format("2006-01-02 15:04:05")))
			b.WriteString(fmt.Sprintf("- Updated: %s\n\n", info.UpdatedAt.Format("2006-01-02 15:04:05")))
		}
	}

	b.WriteString("## Management\n\n")
	b.WriteString("- View status: `otto-stack status --shared`\n")
	b.WriteString("- Shared containers are automatically managed by otto-stack\n")
	b.WriteString("- Orphaned projects are automatically cleaned up\n\n")
	b.WriteString("## Important\n\n")
	b.WriteString("⚠️  Do not manually edit files in this directory. Use otto-stack CLI commands.\n\n")
	b.WriteString(fmt.Sprintf("*Last updated: %s*\n", time.Now().Format("2006-01-02 15:04:05")))

	return b.String()
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
