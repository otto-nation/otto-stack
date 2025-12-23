package stack

import (
	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
)

// NewStackService creates a new stack service with all dependencies
func NewStackService(debug bool) (*services.Service, error) {
	// Create compose manager and get the service
	manager, err := docker.NewManager()
	if err != nil {
		return nil, err
	}

	// Create characteristics resolver
	characteristics, err := services.NewDefaultCharacteristicsResolver()
	if err != nil {
		return nil, err
	}

	// Create project loader
	loader, err := docker.NewDefaultProjectLoader()
	if err != nil {
		return nil, err
	}

	// Create and return stack service
	return services.NewService(manager.GetService(), characteristics, loader)
}
