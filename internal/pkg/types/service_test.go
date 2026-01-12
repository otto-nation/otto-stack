//go:build unit

package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServiceConfig_Fields(t *testing.T) {
	t.Run("validates ServiceConfig structure", func(t *testing.T) {
		config := ServiceConfig{
			Name:           "test-service",
			Category:       "database",
			Container:      ContainerConfig{},
			Documentation:  DocumentationConfig{},
			AllEnvironment: map[string]string{"KEY": "value"},
		}

		assert.Equal(t, "test-service", config.Name)
		assert.Equal(t, "database", config.Category)
		assert.NotNil(t, config.Container)
		assert.NotNil(t, config.Documentation)
		assert.Equal(t, "value", config.AllEnvironment["KEY"])
	})
}

func TestContainerConfig_Fields(t *testing.T) {
	t.Run("validates ContainerConfig structure", func(t *testing.T) {
		config := ContainerConfig{
			Image:       "postgres:13",
			Ports:       []PortConfig{{Host: 5432, Container: 5432}},
			Environment: map[string]string{"DB_NAME": "test"},
			Restart:     "unless-stopped",
		}

		assert.Equal(t, "postgres:13", config.Image)
		assert.Len(t, config.Ports, 1)
		assert.Equal(t, "test", config.Environment["DB_NAME"])
		assert.Equal(t, "unless-stopped", config.Restart)
	})
}

func TestPortConfig_Fields(t *testing.T) {
	t.Run("validates PortConfig structure", func(t *testing.T) {
		port := PortConfig{
			Host:      8080,
			Container: 80,
			Protocol:  "tcp",
		}

		assert.Equal(t, 8080, port.Host)
		assert.Equal(t, 80, port.Container)
		assert.Equal(t, "tcp", port.Protocol)
	})
}

func TestInitContainerConfig_Fields(t *testing.T) {
	t.Run("validates InitContainerConfig structure", func(t *testing.T) {
		config := InitContainerConfig{
			Image:   "busybox",
			Command: []string{"sh", "-c", "echo init"},
		}

		assert.Equal(t, "busybox", config.Image)
		assert.Len(t, config.Command, 3)
		assert.Equal(t, "sh", config.Command[0])
	})
}

func TestDocumentationConfig_Fields(t *testing.T) {
	t.Run("validates DocumentationConfig structure", func(t *testing.T) {
		config := DocumentationConfig{
			WebInterfaces: []WebInterfaceConfig{
				{Name: "Admin", URL: "http://localhost:8080", Port: 8080},
			},
		}

		assert.Len(t, config.WebInterfaces, 1)
		assert.Equal(t, "Admin", config.WebInterfaces[0].Name)
	})
}

func TestWebInterfaceConfig_Fields(t *testing.T) {
	t.Run("validates WebInterfaceConfig structure", func(t *testing.T) {
		config := WebInterfaceConfig{
			Name: "Dashboard",
			URL:  "http://localhost:3000",
			Port: 3000,
		}

		assert.Equal(t, "Dashboard", config.Name)
		assert.Equal(t, "http://localhost:3000", config.URL)
		assert.Equal(t, 3000, config.Port)
	})
}
