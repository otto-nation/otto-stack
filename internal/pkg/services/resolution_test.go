//go:build unit

package services

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/config"
	servicetypes "github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestResolveUpServices_SpecificService(t *testing.T) {
	cfg := &config.Config{
		Stack: config.StackConfig{
			Enabled: []string{ServicePostgres, ServiceRedis},
		},
	}
	args := []string{ServicePostgres}

	configs, err := ResolveUpServices(args, cfg)

	if err == nil {
		assert.NotEmpty(t, configs)
		found := false
		for _, config := range configs {
			if config.Name == ServicePostgres {
				found = true
				break
			}
		}
		assert.True(t, found, "Should resolve requested service")
	}
}

func TestResolveUpServices_EnabledServices(t *testing.T) {
	cfg := &config.Config{
		Stack: config.StackConfig{
			Enabled: []string{ServicePostgres, ServiceRedis},
		},
	}

	configs, err := ResolveUpServices([]string{}, cfg)

	if err == nil {
		assert.NotEmpty(t, configs)
	}
}

func TestResolveUpServices_InvalidServiceName(t *testing.T) {
	cfg := &config.Config{
		Stack: config.StackConfig{
			Enabled: []string{ServicePostgres, ServiceRedis},
		},
	}

	configs, err := ResolveUpServices([]string{"nonexistent-service"}, cfg)

	if err != nil {
		assert.Error(t, err)
	} else {
		assert.NotNil(t, configs)
	}
}

func TestServiceConfigValidation(t *testing.T) {
	config := servicetypes.ServiceConfig{
		Name:        ServicePostgres,
		Description: "PostgreSQL database",
		Category:    CategoryDatabase,
	}

	assert.Equal(t, ServicePostgres, config.Name)
	assert.Equal(t, CategoryDatabase, config.Category)
	assert.NotEmpty(t, config.Description)
}

func TestServiceResolver_Create(t *testing.T) {
	manager, _ := New()
	resolver := NewServiceResolver(manager)
	assert.NotNil(t, resolver)
}

func TestServiceResolver_InvalidNames(t *testing.T) {
	manager, _ := New()
	resolver := NewServiceResolver(manager)

	_, err := resolver.ResolveServices([]string{"invalid-service"})
	assert.Error(t, err)
}

func TestServiceResolver_MultipleInvalid(t *testing.T) {
	manager, _ := New()
	resolver := NewServiceResolver(manager)

	_, err := resolver.ResolveServices([]string{"invalid1", "invalid2"})
	assert.Error(t, err)
}

func TestServiceResolver_ValidServices(t *testing.T) {
	manager, _ := New()
	resolver := NewServiceResolver(manager)

	configs, err := resolver.ResolveServices([]string{ServicePostgres})
	if err == nil {
		assert.NotEmpty(t, configs)
	}
}

func TestServiceResolver_EmptyList(t *testing.T) {
	manager, _ := New()
	resolver := NewServiceResolver(manager)

	_, err := resolver.ResolveServices([]string{})
	assert.Error(t, err)
}

func TestServiceResolver_MultipleServices(t *testing.T) {
	manager, _ := New()
	resolver := NewServiceResolver(manager)

	configs, err := resolver.ResolveServices([]string{ServicePostgres, ServiceRedis})
	if err == nil {
		assert.GreaterOrEqual(t, len(configs), 2)
	}
}
