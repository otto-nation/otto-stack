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

	dockercontainer "github.com/docker/docker/api/types/container"
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
		// Registry file corrupted — attempt Docker-based rebuild
		return m.rebuildFromDocker(context.Background())
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

// Register adds or updates a container in the registry.
func (m *Manager) Register(service, containerName string, project ProjectRef) error {
	registry, err := m.Load()
	if err != nil {
		return pkgerrors.NewSystemError(pkgerrors.ErrCodeOperationFail, messages.ErrorsRegistryLoadFailed, err)
	}

	// Remove stale project references before registering
	m.cleanupOrphans(registry)

	now := time.Now()
	container, exists := registry.Containers[service]

	if !exists {
		container = &ContainerInfo{
			Name:      containerName,
			Projects:  []ProjectRef{project},
			CreatedAt: now,
			UpdatedAt: now,
		}
		registry.Containers[service] = container
	} else {
		if !slices.ContainsFunc(container.Projects, func(r ProjectRef) bool { return r.Name == project.Name }) {
			container.Projects = append(container.Projects, project)
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

	container, exists := registry.Containers[service]
	if !exists {
		return nil
	}

	container.Projects = removeProject(container.Projects, projectName)
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

// ValidateAgainstDocker compares registry entries against running Docker containers.
// Returns warning strings for any registry entry with no running container.
// Non-blocking — never fails, returns nil if Docker is unavailable.
func (m *Manager) ValidateAgainstDocker(ctx context.Context, dockerClient *docker.Client) []string {
	registry, err := m.Load()
	if err != nil {
		return nil
	}

	containers, err := dockerClient.ListContainers(ctx, "")
	if err != nil {
		return nil // Docker unavailable, skip
	}

	running := make(map[string]bool)
	for _, c := range containers {
		if isSharedContainer(c.Name) {
			running[extractServiceName(c.Name)] = true
		}
	}

	var warnings []string
	for service := range registry.Containers {
		if !running[service] {
			warnings = append(warnings,
				fmt.Sprintf(messages.WarningsRegistryNoRunningContainer, service))
		}
	}
	return warnings
}

// PurgeNonShareable removes registry entries for services not in the shareable set.
// Heals any corruption from previous bugs. Called automatically during registration.
func (m *Manager) PurgeNonShareable(shareableServices map[string]bool) error {
	registry, err := m.Load()
	if err != nil {
		return pkgerrors.NewSystemError(
			pkgerrors.ErrCodeOperationFail, messages.ErrorsRegistryLoadFailed, err,
		)
	}

	changed := false
	for service := range registry.Containers {
		if !shareableServices[service] {
			delete(registry.Containers, service)
			changed = true
		}
	}

	if !changed {
		return nil
	}

	if err := m.Save(registry); err != nil {
		return err
	}
	return m.createOrUpdateSharedReadme(registry)
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

// rebuildFromDocker reconstructs the registry from running Docker containers.
// Called when the registry file is corrupted. Returns an empty registry if
// Docker is unavailable — never returns an error.
func (m *Manager) rebuildFromDocker(ctx context.Context) (*Registry, error) {
	dockerClient, err := docker.NewClient(nil)
	if err != nil {
		return NewRegistry(), nil // Docker unavailable, start fresh
	}
	defer func() { _ = dockerClient.Close() }()

	containers, err := dockerClient.GetDockerClient().ContainerList(ctx, dockercontainer.ListOptions{All: true})
	if err != nil {
		return NewRegistry(), nil
	}

	registry := NewRegistry()
	now := time.Now()

	for _, cont := range containers {
		if len(cont.Names) == 0 {
			continue
		}
		name := strings.TrimPrefix(cont.Names[0], "/")
		if !isSharedContainer(name) {
			continue
		}

		serviceName := extractServiceName(name)
		project := cont.Labels[docker.LabelOttoProject]

		entry, exists := registry.Containers[serviceName]
		if !exists {
			entry = &ContainerInfo{
				Name:      name,
				Projects:  []ProjectRef{},
				CreatedAt: now,
				UpdatedAt: now,
			}
			registry.Containers[serviceName] = entry
		}

		if project != "" && !slices.ContainsFunc(entry.Projects, func(r ProjectRef) bool { return r.Name == project }) {
			// ConfigDir is unknown when rebuilding from Docker labels
			entry.Projects = append(entry.Projects, ProjectRef{Name: project})
		}
	}

	return registry, nil
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
			fmt.Fprintf(&b, "### %s\n", service)
			fmt.Fprintf(&b, "- Container: `%s`\n", info.Name)
			projectNames := make([]string, len(info.Projects))
			for i, ref := range info.Projects {
				projectNames[i] = ref.Name
			}
			fmt.Fprintf(&b, "- Projects: %s\n", strings.Join(projectNames, ", "))
			fmt.Fprintf(&b, "- Created: %s\n", info.CreatedAt.Format("2006-01-02 15:04:05"))
			fmt.Fprintf(&b, "- Updated: %s\n\n", info.UpdatedAt.Format("2006-01-02 15:04:05"))
		}
	}

	b.WriteString("## Management\n\n")
	b.WriteString("- View status: `otto-stack status --shared`\n")
	b.WriteString("- Shared containers are automatically managed by otto-stack\n")
	b.WriteString("- Orphaned projects are automatically cleaned up\n\n")
	b.WriteString("## Important\n\n")
	b.WriteString("⚠️  Do not manually edit files in this directory. Use otto-stack CLI commands.\n\n")
	fmt.Fprintf(&b, "*Last updated: %s*\n", time.Now().Format("2006-01-02 15:04:05"))

	return b.String()
}

// removeProject removes the ProjectRef with the given name from a slice.
func removeProject(projects []ProjectRef, name string) []ProjectRef {
	result := make([]ProjectRef, 0, len(projects))
	for _, ref := range projects {
		if ref.Name != name {
			result = append(result, ref)
		}
	}
	return result
}

// cleanupOrphans removes registry entries for projects whose ConfigDir no longer exists on disk.
// This is called inside Register to keep the registry accurate across project lifecycle events.
func (m *Manager) cleanupOrphans(reg *Registry) {
	for service, info := range reg.Containers {
		info.Projects = filterLiveProjectRefs(info.Projects)
		if len(info.Projects) == 0 {
			delete(reg.Containers, service)
		}
	}
}

// filterLiveProjectRefs returns only refs whose ConfigDir still exists on disk.
// Refs without a ConfigDir (produced by rebuildFromDocker) are kept unconditionally.
func filterLiveProjectRefs(projects []ProjectRef) []ProjectRef {
	live := make([]ProjectRef, 0, len(projects))
	for _, ref := range projects {
		if ref.ConfigDir == "" {
			// ConfigDir is unknown (e.g., rebuilt from Docker labels) — keep the entry.
			live = append(live, ref)
			continue
		}
		if _, err := os.Stat(ref.ConfigDir); err == nil {
			live = append(live, ref)
		}
	}
	return live
}
