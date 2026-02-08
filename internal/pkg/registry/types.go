package registry

import "time"

// ContainerInfo represents a shared container in the registry
type ContainerInfo struct {
	Name      string    `yaml:"name" json:"name"`
	Service   string    `yaml:"service" json:"service"`
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

// OrphanInfo represents an orphaned container
type OrphanInfo struct {
	Service   string
	Container string
	Reason    string
}

// ReconcileResult contains the results of a registry reconciliation
type ReconcileResult struct {
	Removed []string // Services removed from registry (container doesn't exist)
	Added   []string // Services added to registry (container exists but not registered)
}

// FindOrphans returns containers with no active projects
func (r *Registry) FindOrphans() []OrphanInfo {
	var orphans []OrphanInfo
	for service, info := range r.Containers {
		if len(info.Projects) == 0 {
			orphans = append(orphans, OrphanInfo{
				Service:   service,
				Container: info.Name,
				Reason:    "no projects using this container",
			})
		}
	}
	return orphans
}
