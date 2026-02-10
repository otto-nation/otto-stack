package display

import "time"

// ServiceStatus represents the status of a service
type ServiceStatus struct {
	Name      string        `json:"name" yaml:"name"`
	State     string        `json:"state" yaml:"state"`
	Health    string        `json:"health" yaml:"health"`
	Provider  string        `json:"provider,omitempty" yaml:"provider,omitempty"`
	Ports     []string      `json:"ports" yaml:"ports"`
	CreatedAt time.Time     `json:"created_at" yaml:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" yaml:"updated_at"`
	Uptime    time.Duration `json:"uptime" yaml:"uptime"`
}

// SharedContainerStatus represents the status of a shared container with usage info
type SharedContainerStatus struct {
	Name      string    `json:"name" yaml:"name"`
	Service   string    `json:"service" yaml:"service"`
	State     string    `json:"state" yaml:"state"`
	Projects  []string  `json:"projects" yaml:"projects"`
	CreatedAt time.Time `json:"created_at" yaml:"created_at"`
	UpdatedAt time.Time `json:"updated_at" yaml:"updated_at"`
}

// SharedStatusResponse represents the response for shared container status queries
type SharedStatusResponse struct {
	SharedContainers []SharedContainerStatus `json:"shared_containers" yaml:"shared_containers"`
	Count            int                     `json:"count" yaml:"count"`
}

// ProjectSharedStatusResponse represents the response for project-specific shared container queries
type ProjectSharedStatusResponse struct {
	Project          string                  `json:"project" yaml:"project"`
	SharedContainers []SharedContainerStatus `json:"shared_containers" yaml:"shared_containers"`
	Count            int                     `json:"count" yaml:"count"`
}

// ServiceStatusResponse represents the response for service status queries
type ServiceStatusResponse struct {
	Services []any `json:"services" yaml:"services"`
	Count    int   `json:"count" yaml:"count"`
}

// ValidationResult represents validation results
type ValidationResult struct {
	Valid    bool              `json:"valid" yaml:"valid"`
	Errors   []ValidationIssue `json:"errors,omitempty" yaml:"errors,omitempty"`
	Warnings []ValidationIssue `json:"warnings,omitempty" yaml:"warnings,omitempty"`
	Summary  map[string]int    `json:"summary" yaml:"summary"`
}

// ValidationIssue represents a validation error or warning
type ValidationIssue struct {
	Type       string `json:"type" yaml:"type"`
	Field      string `json:"field" yaml:"field"`
	Message    string `json:"message" yaml:"message"`
	Severity   string `json:"severity" yaml:"severity"`
	Suggestion string `json:"suggestion,omitempty" yaml:"suggestion,omitempty"`
}

// VersionInfo represents version information
type VersionInfo struct {
	Version   string            `json:"version" yaml:"version"`
	BuildInfo map[string]string `json:"build_info" yaml:"build_info"`
	GoVersion string            `json:"go_version" yaml:"go_version"`
	Platform  string            `json:"platform" yaml:"platform"`
}

// HealthReport represents health check results
type HealthReport struct {
	Overall HealthStatus  `json:"overall" yaml:"overall"`
	Checks  []HealthCheck `json:"checks" yaml:"checks"`
}

// HealthStatus represents overall health status
type HealthStatus struct {
	Status  string `json:"status" yaml:"status"`
	Message string `json:"message" yaml:"message"`
}

// HealthCheck represents individual health check
type HealthCheck struct {
	Name       string `json:"name" yaml:"name"`
	Status     string `json:"status" yaml:"status"`
	Message    string `json:"message" yaml:"message"`
	Suggestion string `json:"suggestion,omitempty" yaml:"suggestion,omitempty"`
	Category   string `json:"category" yaml:"category"`
}

// ServiceCatalog represents available services
type ServiceCatalog struct {
	Categories map[string][]ServiceInfo `json:"categories" yaml:"categories"`
	Total      int                      `json:"total" yaml:"total"`
}

// ServiceInfo represents a service in the catalog
type ServiceInfo struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
	Category    string `json:"category" yaml:"category"`
}

// Options for formatting
type Options struct {
	Format          string
	Quiet           bool
	Compact         bool
	Verbose         bool
	Full            bool
	ShowSummary     bool
	GroupByCategory bool
	ShowProvider    bool
}
