package types

// ServiceConfig represents a complete service configuration
// This is the core domain model used across all packages
type ServiceConfig struct {
	Name           string
	Category       string
	Container      ContainerConfig
	Documentation  DocumentationConfig
	AllEnvironment map[string]string
}

// ContainerConfig represents container-specific configuration
type ContainerConfig struct {
	Image         string
	Ports         []PortConfig
	Environment   map[string]string
	Restart       string
	InitContainer *InitContainerConfig
}

// PortConfig represents port mapping configuration
type PortConfig struct {
	Host      int
	Container int
	Protocol  string
}

// InitContainerConfig represents init container configuration
type InitContainerConfig struct {
	Image   string
	Command []string
}

// DocumentationConfig represents service documentation
type DocumentationConfig struct {
	WebInterfaces []WebInterfaceConfig
}

// WebInterfaceConfig represents web interface configuration
type WebInterfaceConfig struct {
	Name string
	URL  string
	Port int
}
