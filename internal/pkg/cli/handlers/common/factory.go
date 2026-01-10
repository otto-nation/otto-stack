package common

import (
	"sync"

	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
)

var (
	stackServiceCache  *services.Service
	dockerManagerCache *docker.Manager
	resolverCache      services.CharacteristicsResolver
	cacheMutex         sync.RWMutex
)

// NewServiceManager creates a new stack service with all dependencies
func NewServiceManager(debug bool) (*services.Service, error) {
	cacheMutex.RLock()
	if stackServiceCache != nil {
		defer cacheMutex.RUnlock()
		return stackServiceCache, nil
	}
	cacheMutex.RUnlock()

	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	// Double-check after acquiring write lock
	if stackServiceCache != nil {
		return stackServiceCache, nil
	}

	// Create compose manager
	manager, err := getDockerManager()
	if err != nil {
		return nil, err
	}

	// Create characteristics resolver
	resolver, err := getCharacteristicsResolver()
	if err != nil {
		return nil, err
	}

	// Create project loader
	loader, err := docker.NewDefaultProjectLoader()
	if err != nil {
		return nil, err
	}

	// Create and cache service
	service, err := services.NewService(manager.GetService(), resolver, loader)
	if err != nil {
		return nil, err
	}

	stackServiceCache = service
	return service, nil
}

func getDockerManager() (*docker.Manager, error) {
	if dockerManagerCache != nil {
		return dockerManagerCache, nil
	}

	manager, err := docker.NewManager()
	if err != nil {
		return nil, err
	}

	dockerManagerCache = manager
	return manager, nil
}

func getCharacteristicsResolver() (services.CharacteristicsResolver, error) {
	if resolverCache != nil {
		return resolverCache, nil
	}

	resolver, err := services.NewDefaultCharacteristicsResolver()
	if err != nil {
		return nil, err
	}

	resolverCache = resolver
	return resolver, nil
}
