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
