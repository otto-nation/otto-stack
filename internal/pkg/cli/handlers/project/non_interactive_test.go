//go:build unit

package project

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/stretchr/testify/assert"
)

func TestParseServices(t *testing.T) {
	result := parseServices(services.ServicePostgres)
	assert.Equal(t, []string{services.ServicePostgres}, result)

	result = parseServices("postgres,redis,mysql")
	assert.Equal(t, []string{services.ServicePostgres, services.ServiceRedis, services.ServiceMysql}, result)

	result = parseServices("postgres , redis , mysql")
	assert.Equal(t, []string{services.ServicePostgres, services.ServiceRedis, services.ServiceMysql}, result)

	result = parseServices("  postgres  ,  redis  ")
	assert.Equal(t, []string{services.ServicePostgres, services.ServiceRedis}, result)

	result = parseServices("")
	assert.Equal(t, []string{""}, result)
}

func TestGetDefaultValidation(t *testing.T) {
	result := getDefaultValidation()
	assert.NotNil(t, result)

	for key, value := range result {
		assert.True(t, value, "validation key %s should be true", key)
	}

	assert.Equal(t, len(ValidationRegistry), len(result))
}

func TestInitHandler_ValidateProjectName(t *testing.T) {
	handler := NewInitHandler()

	err := handler.ValidateProjectName("my-project")
	assert.NoError(t, err)

	err = handler.ValidateProjectName("project123")
	assert.NoError(t, err)

	err = handler.ValidateProjectName("my_project")
	assert.NoError(t, err)
}
