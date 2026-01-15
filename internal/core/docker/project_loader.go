package docker

import (
	"context"
	"os"
	"path/filepath"

	"github.com/compose-spec/compose-go/v2/types"
	"github.com/otto-nation/otto-stack/internal/pkg/logger"
)

// DefaultProjectLoader implements ProjectLoader using the compose manager
type DefaultProjectLoader struct {
	manager *Manager
}

// NewDefaultProjectLoader creates a new project loader
func NewDefaultProjectLoader() (*DefaultProjectLoader, error) {
	manager, err := NewManager()
	if err != nil {
		return nil, err
	}

	return &DefaultProjectLoader{
		manager: manager,
	}, nil
}

// Load loads a compose project by name
func (l *DefaultProjectLoader) Load(projectName string) (*types.Project, error) {
	logger.Debug("Loading compose project", "project", projectName)

	composePath := DockerComposeFilePath
	logger.Debug("Checking compose file path", "path", composePath)

	// Check if file exists at the default path, if not try current working directory
	if _, err := os.Stat(composePath); os.IsNotExist(err) {
		wd, _ := os.Getwd()
		composePath = filepath.Join(wd, DockerComposeFilePath)
		logger.Debug("File not found at default path, trying working directory", "path", composePath)
	}

	project, err := l.manager.LoadProject(context.Background(), composePath, projectName)
	if err != nil {
		logger.Debug("Failed to load compose project", "error", err, "path", composePath)
		return nil, err
	}

	logger.Debug("Successfully loaded compose project", "project", project.Name, "services", len(project.Services))
	return project, nil
}
