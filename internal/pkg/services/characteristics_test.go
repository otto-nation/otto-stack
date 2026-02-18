//go:build unit

package services

import (
	"testing"
	"time"

	"github.com/otto-nation/otto-stack/internal/core/docker"
	servicetypes "github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestDefaultCharacteristicsResolver_ResolveUpOptions(t *testing.T) {
	resolver := &DefaultCharacteristicsResolver{}
	serviceConfigs := []servicetypes.ServiceConfig{
		{Name: ServicePostgres},
		{Name: ServiceRedis},
	}
	baseOptions := docker.UpOptions{}
	characteristics := []string{}

	result := resolver.ResolveUpOptions(characteristics, serviceConfigs, baseOptions)
	expected := []string{ServicePostgres, ServiceRedis}
	assert.Equal(t, expected, result.Services)
}

func TestDefaultCharacteristicsResolver_ResolveDownOptions(t *testing.T) {
	resolver := &DefaultCharacteristicsResolver{}
	serviceConfigs := []servicetypes.ServiceConfig{
		{Name: ServicePostgres},
	}
	baseOptions := docker.DownOptions{}
	characteristics := []string{}

	result := resolver.ResolveDownOptions(characteristics, serviceConfigs, baseOptions)
	expected := []string{ServicePostgres}
	assert.Equal(t, expected, result.Services)
}

func TestDefaultCharacteristicsResolver_ResolveStopOptions(t *testing.T) {
	resolver := &DefaultCharacteristicsResolver{}
	serviceConfigs := []servicetypes.ServiceConfig{
		{Name: ServiceRedis},
	}
	baseOptions := docker.StopOptions{}
	characteristics := []string{}

	result := resolver.ResolveStopOptions(characteristics, serviceConfigs, baseOptions)
	expected := []string{ServiceRedis}
	assert.Equal(t, expected, result.Services)
}

func TestCharacteristicsResolver_ApplyFlags_Timeout(t *testing.T) {
	resolver := &DefaultCharacteristicsResolver{}
	flags := []string{"--timeout=30"}
	base := docker.DownOptions{}

	result := resolver.applyFlagsToDownOptions(flags, base)
	assert.NotNil(t, result.Timeout)
	if result.Timeout != nil {
		assert.Equal(t, 30*time.Second, *result.Timeout)
	}
}

func TestCharacteristicsResolver_ApplyFlags_InvalidTimeout(t *testing.T) {
	resolver := &DefaultCharacteristicsResolver{}
	flags := []string{"--timeout=invalid"}
	base := docker.DownOptions{}

	result := resolver.applyFlagsToDownOptions(flags, base)
	assert.Equal(t, base, result)
}

func TestCharacteristicsResolver_ApplyFlags_Volumes(t *testing.T) {
	resolver := &DefaultCharacteristicsResolver{}
	flags := []string{"--volumes"}
	base := docker.DownOptions{}

	result := resolver.applyFlagsToDownOptions(flags, base)
	assert.True(t, result.RemoveVolumes)
}

func TestServiceConstantsValidation_ExtractNames(t *testing.T) {
	serviceConfigs := []servicetypes.ServiceConfig{
		{Name: ServicePostgres},
		{Name: ServiceRedis},
		{Name: ServiceMysql},
	}

	names := ExtractServiceNames(serviceConfigs)
	expected := []string{ServicePostgres, ServiceRedis, ServiceMysql}
	assert.Equal(t, expected, names)
}

func TestServiceConstantsValidation_EmptyConfigs(t *testing.T) {
	names := ExtractServiceNames([]servicetypes.ServiceConfig{})
	assert.Empty(t, names)
}

func TestApplyFlagsToUpOptions_RemoveOrphans(t *testing.T) {
	resolver := &DefaultCharacteristicsResolver{}
	flags := []string{"--remove-orphans"}
	base := docker.UpOptions{}
	result := resolver.applyFlagsToUpOptions(flags, base)
	assert.True(t, result.RemoveOrphans)
}

func TestApplyFlagsToUpOptions_Build(t *testing.T) {
	resolver := &DefaultCharacteristicsResolver{}
	flags := []string{"--build"}
	base := docker.UpOptions{}
	result := resolver.applyFlagsToUpOptions(flags, base)
	assert.True(t, result.Build)
}

func TestApplyFlagsToUpOptions_ForceRecreate(t *testing.T) {
	resolver := &DefaultCharacteristicsResolver{}
	flags := []string{"--force-recreate"}
	base := docker.UpOptions{}
	result := resolver.applyFlagsToUpOptions(flags, base)
	assert.True(t, result.ForceRecreate)
}

func TestApplyFlagsToUpOptions_Timeout(t *testing.T) {
	resolver := &DefaultCharacteristicsResolver{}
	flags := []string{"--timeout=30"}
	base := docker.UpOptions{}
	result := resolver.applyFlagsToUpOptions(flags, base)
	assert.NotNil(t, result.Timeout)
}

func TestApplyFlagsToUpOptions_Multiple(t *testing.T) {
	resolver := &DefaultCharacteristicsResolver{}
	flags := []string{"--build", "--remove-orphans"}
	base := docker.UpOptions{}
	result := resolver.applyFlagsToUpOptions(flags, base)
	assert.True(t, result.Build)
	assert.True(t, result.RemoveOrphans)
}

func TestApplyFlagsToDownOptions_RemoveOrphans(t *testing.T) {
	resolver := &DefaultCharacteristicsResolver{}
	flags := []string{"--remove-orphans"}
	base := docker.DownOptions{}
	result := resolver.applyFlagsToDownOptions(flags, base)
	assert.True(t, result.RemoveOrphans)
}

func TestApplyFlagsToDownOptions_Volumes(t *testing.T) {
	resolver := &DefaultCharacteristicsResolver{}
	flags := []string{"--volumes"}
	base := docker.DownOptions{}
	result := resolver.applyFlagsToDownOptions(flags, base)
	assert.True(t, result.RemoveVolumes)
}

func TestApplyFlagsToDownOptions_VolumeShort(t *testing.T) {
	resolver := &DefaultCharacteristicsResolver{}
	flags := []string{"-v"}
	base := docker.DownOptions{}
	result := resolver.applyFlagsToDownOptions(flags, base)
	assert.True(t, result.RemoveVolumes)
}

func TestApplyFlagsToStopOptions_Timeout(t *testing.T) {
	resolver := &DefaultCharacteristicsResolver{}
	flags := []string{"--timeout=30"}
	base := docker.StopOptions{}
	result := resolver.applyFlagsToStopOptions(flags, base)
	assert.NotNil(t, result)
}

func TestApplyFlagsToStopOptions_NoFlags(t *testing.T) {
	resolver := &DefaultCharacteristicsResolver{}
	flags := []string{}
	base := docker.StopOptions{}
	result := resolver.applyFlagsToStopOptions(flags, base)
	assert.Equal(t, base, result)
}

func TestApplyFlagsToStopOptions_InvalidTimeout(t *testing.T) {
	resolver := &DefaultCharacteristicsResolver{}
	flags := []string{"--timeout=invalid"}
	base := docker.StopOptions{}
	result := resolver.applyFlagsToStopOptions(flags, base)
	assert.Equal(t, base, result)
}

func TestNewDefaultCharacteristicsResolver(t *testing.T) {
	resolver, err := NewDefaultCharacteristicsResolver()
	assert.NoError(t, err)
	assert.NotNil(t, resolver)
}
