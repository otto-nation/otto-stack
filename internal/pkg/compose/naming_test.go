package compose

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/config"
)

func TestNamingStrategy_ContainerName(t *testing.T) {
	tests := []struct {
		name        string
		projectName string
		sharing     *config.SharingConfig
		serviceName string
		expected    string
	}{
		{
			name:        "shared container",
			projectName: "my-app",
			sharing:     &config.SharingConfig{Enabled: true},
			serviceName: "postgres",
			expected:    "otto-stack-postgres",
		},
		{
			name:        "project-specific container",
			projectName: "my-app",
			sharing:     &config.SharingConfig{Enabled: false},
			serviceName: "postgres",
			expected:    "my-app-postgres",
		},
		{
			name:        "per-service override - shared",
			projectName: "my-app",
			sharing: &config.SharingConfig{
				Enabled:  false,
				Services: map[string]bool{"postgres": true},
			},
			serviceName: "postgres",
			expected:    "otto-stack-postgres",
		},
		{
			name:        "per-service override - not shared",
			projectName: "my-app",
			sharing: &config.SharingConfig{
				Enabled:  true,
				Services: map[string]bool{"postgres": false},
			},
			serviceName: "postgres",
			expected:    "my-app-postgres",
		},
		{
			name:        "nil sharing config",
			projectName: "my-app",
			sharing:     nil,
			serviceName: "postgres",
			expected:    "my-app-postgres",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ns := NewNamingStrategy(tt.projectName, tt.sharing)
			got := ns.ContainerName(tt.serviceName)
			if got != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, got)
			}
		})
	}
}

func TestNamingStrategy_VolumeName(t *testing.T) {
	tests := []struct {
		name         string
		projectName  string
		sharing      *config.SharingConfig
		serviceName  string
		volumeSuffix string
		expected     string
	}{
		{
			name:         "shared volume",
			projectName:  "my-app",
			sharing:      &config.SharingConfig{Enabled: true},
			serviceName:  "postgres",
			volumeSuffix: "data",
			expected:     "otto-stack-postgres-data",
		},
		{
			name:         "project-specific volume",
			projectName:  "my-app",
			sharing:      &config.SharingConfig{Enabled: false},
			serviceName:  "postgres",
			volumeSuffix: "data",
			expected:     "my-app-postgres-data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ns := NewNamingStrategy(tt.projectName, tt.sharing)
			got := ns.VolumeName(tt.serviceName, tt.volumeSuffix)
			if got != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, got)
			}
		})
	}
}

func TestNamingStrategy_NetworkName(t *testing.T) {
	tests := []struct {
		name        string
		projectName string
		sharing     *config.SharingConfig
		serviceName string
		expected    string
	}{
		{
			name:        "shared network",
			projectName: "my-app",
			sharing:     &config.SharingConfig{Enabled: true},
			serviceName: "postgres",
			expected:    "otto-stack-shared",
		},
		{
			name:        "project network",
			projectName: "my-app",
			sharing:     &config.SharingConfig{Enabled: false},
			serviceName: "postgres",
			expected:    "my-app-network",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ns := NewNamingStrategy(tt.projectName, tt.sharing)
			got := ns.NetworkName(tt.serviceName)
			if got != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, got)
			}
		})
	}
}

func TestNamingStrategy_IsShared(t *testing.T) {
	tests := []struct {
		name        string
		sharing     *config.SharingConfig
		serviceName string
		expected    bool
	}{
		{
			name:        "globally enabled",
			sharing:     &config.SharingConfig{Enabled: true},
			serviceName: "postgres",
			expected:    true,
		},
		{
			name:        "globally disabled",
			sharing:     &config.SharingConfig{Enabled: false},
			serviceName: "postgres",
			expected:    false,
		},
		{
			name: "per-service override true",
			sharing: &config.SharingConfig{
				Enabled:  false,
				Services: map[string]bool{"postgres": true},
			},
			serviceName: "postgres",
			expected:    true,
		},
		{
			name: "per-service override false",
			sharing: &config.SharingConfig{
				Enabled:  true,
				Services: map[string]bool{"postgres": false},
			},
			serviceName: "postgres",
			expected:    false,
		},
		{
			name:        "nil config",
			sharing:     nil,
			serviceName: "postgres",
			expected:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ns := NewNamingStrategy("my-app", tt.sharing)
			got := ns.IsShared(tt.serviceName)
			if got != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}
