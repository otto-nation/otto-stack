package services

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/constants"
	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
	service := Service{
		Name:        "postgres",
		Description: "PostgreSQL database",
		Category:    constants.CategoryDatabase,
		Type:        constants.ServiceTypeContainer,
		Docker: DockerConfig{
			Image: "postgres:15",
			Ports: []string{"5432:5432"},
		},
		Connection: ConnectionConfig{
			Client:      constants.ClientPsql,
			DefaultUser: "postgres",
			DefaultPort: constants.DefaultPortPOSTGRES_port,
		},
	}

	assert.Equal(t, "postgres", service.Name)
	assert.Equal(t, constants.CategoryDatabase, service.Category)
	assert.Equal(t, constants.ClientPsql, service.Connection.Client)
}

func TestManager_GetService(t *testing.T) {
	manager := &Manager{
		services: map[string]Service{
			"postgres": {
				Name:     "postgres",
				Category: constants.CategoryDatabase,
			},
		},
	}

	service, err := manager.GetService("postgres")
	assert.NoError(t, err)
	assert.Equal(t, "postgres", service.Name)

	_, err = manager.GetService("nonexistent")
	assert.Error(t, err)
}

func TestManager_GetServicesByCategory(t *testing.T) {
	manager := &Manager{
		services: map[string]Service{
			"postgres": {
				Name:     "postgres",
				Category: constants.CategoryDatabase,
			},
			"redis": {
				Name:     "redis",
				Category: constants.CategoryCache,
			},
		},
	}

	categories := manager.GetServicesByCategory()
	assert.Len(t, categories[constants.CategoryDatabase], 1)
	assert.Len(t, categories[constants.CategoryCache], 1)
	assert.Equal(t, "postgres", categories[constants.CategoryDatabase][0].Name)
}

func TestManager_ValidateServices(t *testing.T) {
	manager := &Manager{
		services: map[string]Service{
			"postgres": {Name: "postgres"},
			"redis":    {Name: "redis"},
		},
	}

	err := manager.ValidateServices([]string{"postgres", "redis"})
	assert.NoError(t, err)

	err = manager.ValidateServices([]string{"postgres", "nonexistent"})
	assert.Error(t, err)
}

func TestManager_BuildConnectCommand(t *testing.T) {
	manager := &Manager{
		services: map[string]Service{
			"postgres": {
				Name: "postgres",
				Connection: ConnectionConfig{
					Client:      constants.ClientPsql,
					DefaultUser: "postgres",
					DefaultPort: constants.DefaultPortPOSTGRES_port,
					UserFlag:    "-U",
					HostFlag:    "-h",
					PortFlag:    "-p",
				},
			},
		},
	}

	cmd, err := manager.BuildConnectCommand("postgres", map[string]string{
		"host": "localhost",
		"user": "testuser",
	})

	assert.NoError(t, err)
	assert.Contains(t, cmd, constants.ClientPsql)
	assert.Contains(t, cmd, "-h")
	assert.Contains(t, cmd, "localhost")
	assert.Contains(t, cmd, "-U")
	assert.Contains(t, cmd, "testuser")
}
