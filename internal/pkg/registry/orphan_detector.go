package registry

import (
	"context"
	"os"
	"path/filepath"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/messages"
)

// OrphanDetector handles detection of orphaned containers
type OrphanDetector struct {
	manager *Manager
}

// NewOrphanDetector creates a new orphan detector
func NewOrphanDetector(manager *Manager) *OrphanDetector {
	return &OrphanDetector{manager: manager}
}

// FindOrphans returns containers with no active projects (basic check)
func (d *OrphanDetector) FindOrphans() ([]OrphanInfo, error) {
	registry, err := d.manager.Load()
	if err != nil {
		return nil, err
	}
	return registry.FindOrphans(), nil
}

// FindOrphansWithChecks performs enhanced orphan detection with filesystem and Docker checks
func (d *OrphanDetector) FindOrphansWithChecks(ctx context.Context, dockerClient *docker.Client) ([]OrphanInfo, error) {
	registry, err := d.manager.Load()
	if err != nil {
		return nil, err
	}

	containers, err := dockerClient.ListContainers(ctx, "")
	if err != nil {
		return nil, err
	}

	containerMap := d.buildContainerMap(containers)

	var orphans []OrphanInfo
	orphans = append(orphans, d.checkRegisteredContainers(registry, containerMap)...)
	orphans = append(orphans, d.checkZombieContainers(registry, containers)...)

	return orphans, nil
}

func (d *OrphanDetector) buildContainerMap(containers []docker.ContainerInfo) map[string]bool {
	m := make(map[string]bool)
	for _, c := range containers {
		m[c.Name] = true
	}
	return m
}

func (d *OrphanDetector) checkRegisteredContainers(registry *Registry, containerMap map[string]bool) []OrphanInfo {
	var orphans []OrphanInfo

	for service, info := range registry.Containers {
		if orphan := d.checkContainer(service, info, containerMap); orphan != nil {
			orphans = append(orphans, *orphan)
		}
	}

	return orphans
}

func (d *OrphanDetector) checkContainer(service string, info *ContainerInfo, containerMap map[string]bool) *OrphanInfo {
	// Container not in Docker
	if !containerMap[info.Name] {
		return &OrphanInfo{
			Service:        service,
			Container:      info.Name,
			Reason:         messages.OrphanReasonNotInDocker,
			Severity:       OrphanSeverityWarning,
			ContainerState: ContainerStateNotFound,
		}
	}

	// No projects registered
	if len(info.Projects) == 0 {
		return &OrphanInfo{
			Service:        service,
			Container:      info.Name,
			Reason:         messages.OrphanReasonNoProjects,
			Severity:       OrphanSeveritySafe,
			ContainerState: ContainerStateRunning,
		}
	}

	// Check project filesystem
	return d.checkProjectFilesystem(service, info)
}

func (d *OrphanDetector) checkProjectFilesystem(service string, info *ContainerInfo) *OrphanInfo {
	var existingProjects []string
	for _, project := range info.Projects {
		if projectExists(project) {
			existingProjects = append(existingProjects, project)
		}
	}

	if len(existingProjects) == 0 {
		return &OrphanInfo{
			Service:        service,
			Container:      info.Name,
			Reason:         messages.OrphanReasonAllProjectsDeleted,
			Severity:       OrphanSeveritySafe,
			ContainerState: ContainerStateRunning,
			ProjectsFound:  info.Projects,
		}
	}

	if len(existingProjects) < len(info.Projects) {
		return &OrphanInfo{
			Service:        service,
			Container:      info.Name,
			Reason:         messages.OrphanReasonSomeProjectsDeleted,
			Severity:       OrphanSeverityWarning,
			ContainerState: ContainerStateRunning,
			ProjectsFound:  existingProjects,
		}
	}

	return nil
}

func (d *OrphanDetector) checkZombieContainers(registry *Registry, containers []docker.ContainerInfo) []OrphanInfo {
	var orphans []OrphanInfo

	for _, container := range containers {
		if !isSharedContainer(container.Name) {
			continue
		}

		service := extractServiceName(container.Name)
		if _, exists := registry.Containers[service]; !exists {
			orphans = append(orphans, OrphanInfo{
				Service:        service,
				Container:      container.Name,
				Reason:         messages.OrphanReasonNotRegistered,
				Severity:       OrphanSeverityCritical,
				ContainerState: ContainerStateRunning,
			})
		}
	}

	return orphans
}

// projectExists checks if a project directory exists
func projectExists(projectName string) bool {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	possiblePaths := []string{
		filepath.Join(homeDir, "projects", projectName, core.OttoStackDir),
		filepath.Join(homeDir, "dev", projectName, core.OttoStackDir),
		filepath.Join(homeDir, "workspace", projectName, core.OttoStackDir),
		filepath.Join(homeDir, projectName, core.OttoStackDir),
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return true
		}
	}

	return false
}
