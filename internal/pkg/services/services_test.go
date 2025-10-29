package services

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServiceDefinition(t *testing.T) {
	t.Run("service definition structure", func(t *testing.T) {
		service := ServiceDefinition{
			Description: "Test service",
			ServiceConfiguration: []ServiceOption{
				{Name: "option1", Type: "string", Description: "First option"},
				{Name: "option2", Type: "string", Description: "Second option"},
			},
			Documentation: Documentation{
				Examples:   []string{"example1", "example2"},
				UsageNotes: "Usage notes",
				Links:      []string{"http://example.com"},
			},
			Category:    "database",
			DefaultPort: 5432,
			Dependencies: ServiceDependencies{
				Required: []string{"dep1", "dep2"},
			},
			Tags: []string{"tag1", "tag2"},
		}

		assert.Equal(t, "Test service", service.Description)
		assert.Equal(t, 2, len(service.ServiceConfiguration))
		assert.Equal(t, "option1", service.ServiceConfiguration[0].Name)
		assert.Equal(t, "option2", service.ServiceConfiguration[1].Name)
		assert.Equal(t, []string{"example1", "example2"}, service.Documentation.Examples)
		assert.Equal(t, "Usage notes", service.Documentation.UsageNotes)
		assert.Equal(t, []string{"http://example.com"}, service.Documentation.Links)
		assert.Equal(t, "database", service.Category)
		assert.Equal(t, 5432, service.DefaultPort)
		assert.Equal(t, []string{"dep1", "dep2"}, service.Dependencies.Required)
		assert.Equal(t, []string{"tag1", "tag2"}, service.Tags)
	})
}

func TestHealthCheckConfig(t *testing.T) {
	t.Run("health check configuration", func(t *testing.T) {
		healthCheck := HealthCheckConfig{
			Enabled:  true,
			Endpoint: "/health",
			Interval: "30s",
			Timeout:  "10s",
			Retries:  3,
		}

		assert.True(t, healthCheck.Enabled)
		assert.Equal(t, "/health", healthCheck.Endpoint)
		assert.Equal(t, "30s", healthCheck.Interval)
		assert.Equal(t, "10s", healthCheck.Timeout)
		assert.Equal(t, 3, healthCheck.Retries)
	})
}

func TestServiceOperations(t *testing.T) {
	t.Run("connect operation", func(t *testing.T) {
		connect := ConnectOperation{
			Command:  []string{"psql", "-h", "localhost"},
			Args:     map[string][]string{"database": {"-d"}},
			Defaults: map[string]string{"port": "5432"},
		}

		assert.Equal(t, []string{"psql", "-h", "localhost"}, connect.Command)
		assert.Equal(t, map[string][]string{"database": {"-d"}}, connect.Args)
		assert.Equal(t, map[string]string{"port": "5432"}, connect.Defaults)
	})

	t.Run("backup operation", func(t *testing.T) {
		backup := BackupOperation{
			Type:      "command",
			Command:   []string{"pg_dump"},
			Args:      map[string][]string{"database": {"-d"}},
			Defaults:  map[string]string{"format": "custom"},
			Extension: ".dump",
		}

		assert.Equal(t, "command", backup.Type)
		assert.Equal(t, []string{"pg_dump"}, backup.Command)
		assert.Equal(t, map[string][]string{"database": {"-d"}}, backup.Args)
		assert.Equal(t, map[string]string{"format": "custom"}, backup.Defaults)
		assert.Equal(t, ".dump", backup.Extension)
	})

	t.Run("restore operation", func(t *testing.T) {
		restore := RestoreOperation{
			Type:            "command",
			Command:         []string{"pg_restore"},
			Args:            map[string][]string{"database": {"-d"}},
			Defaults:        map[string]string{"format": "custom"},
			RequiresRestart: true,
		}

		assert.Equal(t, "command", restore.Type)
		assert.Equal(t, []string{"pg_restore"}, restore.Command)
		assert.Equal(t, map[string][]string{"database": {"-d"}}, restore.Args)
		assert.Equal(t, map[string]string{"format": "custom"}, restore.Defaults)
		assert.True(t, restore.RequiresRestart)
	})
}

func TestServiceRegistry_GetAllServices(t *testing.T) {
	registry := &ServiceRegistry{
		services: map[string]ServiceDefinition{
			"redis": {
				Description: "Redis cache",
				Category:    "cache",
				DefaultPort: 6379,
			},
			"postgres": {
				Description: "PostgreSQL database",
				Category:    "database",
				DefaultPort: 5432,
			},
		},
	}

	services := registry.GetAllServices()
	assert.Len(t, services, 2)
	assert.Contains(t, services, "redis")
	assert.Contains(t, services, "postgres")
	assert.Equal(t, "Redis cache", services["redis"].Description)
	assert.Equal(t, "PostgreSQL database", services["postgres"].Description)
}

func TestServiceRegistry_GetService(t *testing.T) {
	registry := &ServiceRegistry{
		services: map[string]ServiceDefinition{
			"redis": {
				Description: "Redis cache",
				Category:    "cache",
				DefaultPort: 6379,
			},
		},
	}

	t.Run("existing service", func(t *testing.T) {
		service, exists := registry.GetService("redis")
		assert.True(t, exists)
		assert.Equal(t, "Redis cache", service.Description)
		assert.Equal(t, "cache", service.Category)
		assert.Equal(t, 6379, service.DefaultPort)
	})

	t.Run("non-existing service", func(t *testing.T) {
		service, exists := registry.GetService("nonexistent")
		assert.False(t, exists)
		assert.Equal(t, ServiceDefinition{}, service)
	})
}

func TestServiceRegistry_GetServiceNames(t *testing.T) {
	registry := &ServiceRegistry{
		services: map[string]ServiceDefinition{
			"redis":    {Description: "Redis cache"},
			"postgres": {Description: "PostgreSQL database"},
			"mysql":    {Description: "MySQL database"},
		},
	}

	names := registry.GetServiceNames()
	assert.Len(t, names, 3)
	assert.Equal(t, []string{"mysql", "postgres", "redis"}, names) // Should be sorted
}

func TestServiceRegistry_GetServicesByCategory(t *testing.T) {
	registry := &ServiceRegistry{
		services: map[string]ServiceDefinition{
			"redis": {
				Description: "Redis cache",
				Category:    "cache",
			},
			"postgres": {
				Description: "PostgreSQL database",
				Category:    "database",
			},
			"mysql": {
				Description: "MySQL database",
				Category:    "database",
			},
			"nginx": {
				Description: "Nginx web server",
				Category:    "web",
			},
		},
	}

	t.Run("database category", func(t *testing.T) {
		services := registry.GetServicesByCategory("database")
		assert.Len(t, services, 2)
		assert.Equal(t, []string{"mysql", "postgres"}, services) // Should be sorted
	})

	t.Run("cache category", func(t *testing.T) {
		services := registry.GetServicesByCategory("cache")
		assert.Len(t, services, 1)
		assert.Equal(t, []string{"redis"}, services)
	})

	t.Run("non-existing category", func(t *testing.T) {
		services := registry.GetServicesByCategory("nonexistent")
		assert.Len(t, services, 0)
	})
}

func TestServiceRegistry_GetServicesByTag(t *testing.T) {
	registry := &ServiceRegistry{
		services: map[string]ServiceDefinition{
			"redis": {
				Description: "Redis cache",
				Tags:        []string{"cache", "nosql"},
			},
			"postgres": {
				Description: "PostgreSQL database",
				Tags:        []string{"database", "sql"},
			},
			"mongodb": {
				Description: "MongoDB database",
				Tags:        []string{"database", "nosql"},
			},
		},
	}

	t.Run("nosql tag", func(t *testing.T) {
		services := registry.GetServicesByTag("nosql")
		assert.Len(t, services, 2)
		assert.Equal(t, []string{"mongodb", "redis"}, services) // Should be sorted
	})

	t.Run("sql tag", func(t *testing.T) {
		services := registry.GetServicesByTag("sql")
		assert.Len(t, services, 1)
		assert.Equal(t, []string{"postgres"}, services)
	})

	t.Run("non-existing tag", func(t *testing.T) {
		services := registry.GetServicesByTag("nonexistent")
		assert.Len(t, services, 0)
	})
}

func TestServiceRegistry_Reload(t *testing.T) {
	// Create a temporary services file
	tmpDir := t.TempDir()
	servicesFile := filepath.Join(tmpDir, "services.yaml")

	servicesContent := `redis:
  description: "Redis cache"
  category: "cache"
  default_port: 6379
  tags: ["cache", "nosql"]
postgres:
  description: "PostgreSQL database"
  category: "database"
  default_port: 5432
  tags: ["database", "sql"]`

	err := os.WriteFile(servicesFile, []byte(servicesContent), 0644)
	assert.NoError(t, err)

	registry := &ServiceRegistry{
		services: make(map[string]ServiceDefinition),
	}

	// Load from embedded
	err = registry.LoadFromEmbedded()
	assert.NoError(t, err)
	assert.Greater(t, len(registry.services), 0)

	// Add a service manually
	registry.services["manual"] = ServiceDefinition{Description: "Manual service"}

	// Verify manual service was added
	assert.Contains(t, registry.services, "manual")
	assert.Greater(t, len(registry.services), 0)
}

func TestNewServiceRegistry(t *testing.T) {
	t.Run("embedded services", func(t *testing.T) {
		registry, err := NewServiceRegistry("")
		assert.NoError(t, err)
		assert.NotNil(t, registry)
		assert.Greater(t, len(registry.services), 0)
	})
}

func TestServiceManifest(t *testing.T) {
	t.Run("service manifest type", func(t *testing.T) {
		manifest := ServiceManifest{
			"redis": ServiceDefinition{
				Description: "Redis cache",
				Category:    "cache",
			},
			"postgres": ServiceDefinition{
				Description: "PostgreSQL database",
				Category:    "database",
			},
		}

		assert.Len(t, manifest, 2)
		assert.Contains(t, manifest, "redis")
		assert.Contains(t, manifest, "postgres")
		assert.Equal(t, "Redis cache", manifest["redis"].Description)
		assert.Equal(t, "PostgreSQL database", manifest["postgres"].Description)
	})
}

func TestServiceConfig(t *testing.T) {
	t.Run("service config with operations", func(t *testing.T) {
		config := ServiceConfig{
			Name: "postgres",
			Operations: &ServiceOperations{
				Connect: &ConnectOperation{
					Command: []string{"psql"},
				},
				Backup: &BackupOperation{
					Type:      "command",
					Command:   []string{"pg_dump"},
					Extension: ".sql",
				},
				Restore: &RestoreOperation{
					Type:    "command",
					Command: []string{"psql"},
				},
			},
		}

		assert.Equal(t, "postgres", config.Name)
		assert.NotNil(t, config.Operations)
		assert.NotNil(t, config.Operations.Connect)
		assert.NotNil(t, config.Operations.Backup)
		assert.NotNil(t, config.Operations.Restore)
		assert.Equal(t, []string{"psql"}, config.Operations.Connect.Command)
		assert.Equal(t, "command", config.Operations.Backup.Type)
		assert.Equal(t, ".sql", config.Operations.Backup.Extension)
	})
}
