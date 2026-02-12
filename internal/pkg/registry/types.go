package registry

import (
	"time"

	"github.com/otto-nation/otto-stack/internal/pkg/messages"
)

// ContainerInfo represents a shared container in the registry
type ContainerInfo struct {
	Name      string    `yaml:"name" json:"name"`
	Projects  []string  `yaml:"projects" json:"projects"`
	CreatedAt time.Time `yaml:"created_at" json:"created_at"`
	UpdatedAt time.Time `yaml:"updated_at" json:"updated_at"`
}

// Registry represents the shared container registry
type Registry struct {
	Containers map[string]*ContainerInfo `yaml:"shared_containers" json:"shared_containers"`
}

// NewRegistry creates a new empty registry
func NewRegistry() *Registry {
	return &Registry{
		Containers: make(map[string]*ContainerInfo),
	}
}

// OrphanSeverity indicates how safe it is to remove an orphan
type OrphanSeverity string

const (
	OrphanSeveritySafe     OrphanSeverity = "safe"     // Safe to remove
	OrphanSeverityWarning  OrphanSeverity = "warning"  // Needs investigation
	OrphanSeverityCritical OrphanSeverity = "critical" // Running but not registered
)

// Container states for orphan detection
const (
	ContainerStateRunning  = "running"
	ContainerStateStopped  = "stopped"
	ContainerStateNotFound = "not_found"
	ContainerStateUnknown  = "unknown"
)

// OrphanInfo represents an orphaned container
type OrphanInfo struct {
	Service        string
	Container      string
	Reason         string
	Severity       OrphanSeverity
	ContainerState string   // running, stopped, not_found
	ProjectsFound  []string // Projects that still reference this
}

// ReconcileResult contains the results of a registry reconciliation
type ReconcileResult struct {
	Removed []string // Services removed from registry (container doesn't exist)
	Added   []string // Services added to registry (container exists but not registered)
}

// FindOrphans returns containers with no active projects (basic check)
func (r *Registry) FindOrphans() []OrphanInfo {
	var orphans []OrphanInfo
	for service, info := range r.Containers {
		if len(info.Projects) == 0 {
			orphans = append(orphans, OrphanInfo{
				Service:        service,
				Container:      info.Name,
				Reason:         messages.OrphanReasonNoProjects,
				Severity:       OrphanSeveritySafe,
				ContainerState: ContainerStateUnknown,
			})
		}
	}
	return orphans
}
