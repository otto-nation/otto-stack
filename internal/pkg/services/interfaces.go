package services

import (
	"time"

	"github.com/compose-spec/compose-go/v2/types"
	"github.com/docker/compose/v5/pkg/api"
	"github.com/otto-nation/otto-stack/internal/pkg/config"
)

// CharacteristicsResolver converts service characteristics to compose options
type CharacteristicsResolver interface {
	ResolveUpOptions(characteristics []string, base UpOptions) UpOptions
	ResolveDownOptions(characteristics []string, base DownOptions) DownOptions
	ResolveStopOptions(characteristics []string, base StopOptions) StopOptions
}

// ProjectLoader loads compose projects
type ProjectLoader interface {
	Load(projectName string) (*types.Project, error)
}

// ConfigService provides configuration management functionality
type ConfigService interface {
	LoadConfig() (*config.Config, error)
	SaveConfig(cfg *config.Config) error
	ValidateConfig(cfg *config.Config) error
	GetConfigHash(cfg *config.Config) (string, error)
}

// UpOptions defines options for compose up operations
type UpOptions struct {
	Build         bool
	ForceRecreate bool
	RemoveOrphans bool
	Timeout       time.Duration
}

// DownOptions defines options for compose down operations
type DownOptions struct {
	RemoveVolumes bool
	RemoveOrphans bool
	Timeout       time.Duration
}

// StopOptions defines options for compose stop operations
type StopOptions struct {
	Services []string
	Timeout  time.Duration
}

// LogOptions defines options for compose logs operations
type LogOptions struct {
	Services   []string
	Follow     bool
	Timestamps bool
	Tail       string
}

// ToSDK converts our options to official SDK options
func (o UpOptions) ToSDK() api.UpOptions {
	upOptions := api.UpOptions{
		Create: api.CreateOptions{
			RemoveOrphans: o.RemoveOrphans,
		},
		Start: api.StartOptions{},
	}

	if o.Build {
		upOptions.Create.Build = &api.BuildOptions{}
	}

	if o.ForceRecreate {
		upOptions.Create.Recreate = api.RecreateForce
	}

	if o.Timeout > 0 {
		upOptions.Create.Timeout = &o.Timeout
	}

	return upOptions
}

func (o DownOptions) ToSDK() api.DownOptions {
	downOptions := api.DownOptions{
		RemoveOrphans: o.RemoveOrphans,
		Volumes:       o.RemoveVolumes,
	}

	if o.Timeout > 0 {
		downOptions.Timeout = &o.Timeout
	}

	return downOptions
}

func (o StopOptions) ToSDK() api.StopOptions {
	stopOptions := api.StopOptions{
		Services: o.Services,
	}

	if o.Timeout > 0 {
		stopOptions.Timeout = &o.Timeout
	}

	return stopOptions
}

func (o LogOptions) ToSDK() api.LogOptions {
	return api.LogOptions{
		Services:   o.Services,
		Follow:     o.Follow,
		Timestamps: o.Timestamps,
		Tail:       o.Tail,
	}
}
