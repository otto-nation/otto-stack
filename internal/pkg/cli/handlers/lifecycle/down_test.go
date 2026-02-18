package lifecycle

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestDownHandler_serviceNamesToConfigs(t *testing.T) {
	handler := &DownHandler{}
	configs := handler.serviceNamesToConfigs([]string{"postgres", "redis"})
	assert.Len(t, configs, 2)
	assert.Equal(t, "postgres", configs[0].Name)
	assert.Equal(t, "redis", configs[1].Name)
}

func TestDownHandler_filterOutShared(t *testing.T) {
	handler := &DownHandler{}
	serviceConfigs := []types.ServiceConfig{
		{Name: "postgres"},
		{Name: "redis"},
		{Name: "mysql"},
	}
	sharedServices := []string{"redis"}

	filtered := handler.filterOutShared(sharedServices, serviceConfigs)
	assert.Len(t, filtered, 2)
	assert.Equal(t, "postgres", filtered[0].Name)
	assert.Equal(t, "mysql", filtered[1].Name)
}
