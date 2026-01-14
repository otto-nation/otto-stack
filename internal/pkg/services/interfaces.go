package services

import (
	"github.com/compose-spec/compose-go/v2/types"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	servicetypes "github.com/otto-nation/otto-stack/internal/pkg/types"
)

// CharacteristicsResolver converts service characteristics to compose options
type CharacteristicsResolver interface {
	ResolveUpOptions(characteristics []string, serviceConfigs []servicetypes.ServiceConfig, base docker.UpOptions) docker.UpOptions
	ResolveDownOptions(characteristics []string, serviceConfigs []servicetypes.ServiceConfig, base docker.DownOptions) docker.DownOptions
	ResolveStopOptions(characteristics []string, serviceConfigs []servicetypes.ServiceConfig, base docker.StopOptions) docker.StopOptions
}

// ProjectLoader loads compose projects
type ProjectLoader interface {
	Load(projectName string) (*types.Project, error)
}
