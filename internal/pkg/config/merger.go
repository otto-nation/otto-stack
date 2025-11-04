package config

import (
	"reflect"

	"github.com/otto-nation/otto-stack/internal/pkg/types"
)

// ConfigMerger handles merging of base and local configurations
type ConfigMerger struct{}

// NewConfigMerger creates a new config merger
func NewConfigMerger() *ConfigMerger {
	return &ConfigMerger{}
}

// Merge combines base config with local overrides (local takes precedence)
func (m *ConfigMerger) Merge(base, local *types.OttoStackConfig) *types.OttoStackConfig {
	if local == nil {
		return base
	}
	if base == nil {
		return local
	}

	merged := *base // Copy base config

	// Merge project settings
	m.mergeProjectConfig(&merged.Project, &local.Project)

	// Merge stack settings
	m.mergeStackConfig(&merged.Stack, &local.Stack)

	// Override service configuration if specified in local
	if local.ServiceConfiguration != nil {
		merged.ServiceConfiguration = local.ServiceConfiguration
	}

	return &merged
}

// mergeProjectConfig merges project configuration settings
func (m *ConfigMerger) mergeProjectConfig(base, local *types.ProjectInfo) {
	if local.Name != "" {
		base.Name = local.Name
	}
	if local.Environment != "" {
		base.Environment = local.Environment
	}
}

// mergeStackConfig merges stack configuration settings
func (m *ConfigMerger) mergeStackConfig(base, local *types.StackConfig) {
	// Use reflection to merge non-zero values
	baseValue := reflect.ValueOf(base).Elem()
	localValue := reflect.ValueOf(local).Elem()

	for i := 0; i < localValue.NumField(); i++ {
		localField := localValue.Field(i)
		if !localField.IsZero() {
			baseValue.Field(i).Set(localField)
		}
	}
}
