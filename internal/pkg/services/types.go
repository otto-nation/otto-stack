package services

// Service represents a service definition
type Service struct {
	Name         string            `yaml:"name"`
	Description  string            `yaml:"description"`
	Category     string            `yaml:"category"`
	Type         string            `yaml:"type,omitempty"`
	Docker       DockerConfig      `yaml:"docker,omitempty"`
	Connection   ConnectionConfig  `yaml:"connection,omitempty"`
	Dependencies []string          `yaml:"dependencies,omitempty"`
	Environment  map[string]string `yaml:"environment,omitempty"`
}

// DockerConfig represents Docker configuration
type DockerConfig struct {
	Image         string   `yaml:"image,omitempty"`
	Ports         []string `yaml:"ports,omitempty"`
	Environment   []string `yaml:"environment,omitempty"`
	Volumes       []string `yaml:"volumes,omitempty"`
	SimpleVolumes []string `yaml:"simple_volumes,omitempty"`
	Command       []string `yaml:"command,omitempty"`
	Restart       string   `yaml:"restart,omitempty"`
	DependsOn     []string `yaml:"depends_on,omitempty"`
}

// ConnectionConfig represents connection configuration
type ConnectionConfig struct {
	Client      string   `yaml:"client,omitempty"`
	DefaultUser string   `yaml:"default_user,omitempty"`
	DefaultPort int      `yaml:"default_port,omitempty"`
	UserFlag    string   `yaml:"user_flag,omitempty"`
	HostFlag    string   `yaml:"host_flag,omitempty"`
	PortFlag    string   `yaml:"port_flag,omitempty"`
	DBFlag      string   `yaml:"database_flag,omitempty"`
	ExtraFlags  []string `yaml:"extra_flags,omitempty"`
}
