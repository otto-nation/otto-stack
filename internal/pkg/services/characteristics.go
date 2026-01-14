package services

import (
	"strconv"
	"strings"
	"time"

	servicetypes "github.com/otto-nation/otto-stack/internal/pkg/types"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/core/docker"
)

const expectedEnvParts = 2

// DefaultCharacteristicsResolver implements CharacteristicsResolver using existing docker characteristics
type DefaultCharacteristicsResolver struct {
	resolver *docker.ServiceCharacteristicsResolver
}

// NewDefaultCharacteristicsResolver creates a new characteristics resolver
func NewDefaultCharacteristicsResolver() (*DefaultCharacteristicsResolver, error) {
	resolver, err := docker.NewServiceCharacteristicsResolver()
	if err != nil {
		return nil, err
	}

	return &DefaultCharacteristicsResolver{
		resolver: resolver,
	}, nil
}

// ResolveUpOptions converts characteristics to up options
func (r *DefaultCharacteristicsResolver) ResolveUpOptions(characteristics []string, serviceConfigs []servicetypes.ServiceConfig, base docker.UpOptions) docker.UpOptions {
	base.Services = ExtractServiceNames(serviceConfigs)
	flags := r.resolver.ResolveComposeUpFlags(characteristics)
	return r.applyFlagsToUpOptions(flags, base)
}

// ResolveDownOptions converts characteristics to down options
func (r *DefaultCharacteristicsResolver) ResolveDownOptions(characteristics []string, serviceConfigs []servicetypes.ServiceConfig, base docker.DownOptions) docker.DownOptions {
	base.Services = ExtractServiceNames(serviceConfigs)
	flags := r.resolver.ResolveComposeDownFlags(characteristics)
	return r.applyFlagsToDownOptions(flags, base)
}

// ResolveStopOptions converts characteristics to stop options
func (r *DefaultCharacteristicsResolver) ResolveStopOptions(characteristics []string, serviceConfigs []servicetypes.ServiceConfig, base docker.StopOptions) docker.StopOptions {
	base.Services = ExtractServiceNames(serviceConfigs)
	flags := r.resolver.ResolveComposeDownFlags(characteristics) // Use down flags for stop
	return r.applyFlagsToStopOptions(flags, base)
}

func (r *DefaultCharacteristicsResolver) applyFlagsToUpOptions(flags []string, base docker.UpOptions) docker.UpOptions {
	options := base

	for _, flag := range flags {
		switch {
		case flag == docker.FlagPrefix+core.FlagRemoveOrphans:
			options.RemoveOrphans = true
		case flag == docker.FlagPrefix+docker.FlagBuild:
			options.Build = true
		case flag == docker.FlagPrefix+docker.FlagForceRecreate:
			options.ForceRecreate = true
		case strings.HasPrefix(flag, docker.FlagPrefix+docker.FlagTimeout+"="):
			if timeout, err := r.parseTimeout(flag); err == nil {
				options.Timeout = &timeout
			}
		}
	}

	return options
}

func (r *DefaultCharacteristicsResolver) applyFlagsToDownOptions(flags []string, base docker.DownOptions) docker.DownOptions {
	options := base

	for _, flag := range flags {
		switch {
		case flag == docker.FlagPrefix+core.FlagRemoveOrphans:
			options.RemoveOrphans = true
		case flag == docker.FlagPrefix+core.FlagVolumes || flag == "-v":
			options.RemoveVolumes = true
		case strings.HasPrefix(flag, docker.FlagPrefix+docker.FlagTimeout+"="):
			if timeout, err := r.parseTimeout(flag); err == nil {
				options.Timeout = &timeout
			}
		}
	}

	return options
}

func (r *DefaultCharacteristicsResolver) applyFlagsToStopOptions(flags []string, base docker.StopOptions) docker.StopOptions {
	options := base

	for _, flag := range flags {
		if strings.HasPrefix(flag, docker.FlagPrefix+docker.FlagTimeout+"=") {
			if timeout, err := r.parseTimeout(flag); err == nil {
				options.Timeout = &timeout
			}
		}
	}

	return options
}

func (r *DefaultCharacteristicsResolver) parseTimeout(flag string) (time.Duration, error) {
	parts := strings.Split(flag, "=")
	if len(parts) != expectedEnvParts {
		return 0, nil
	}

	seconds, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, err
	}

	return time.Duration(seconds) * time.Second, nil
}
