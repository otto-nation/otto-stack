package operations

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestFilterInitContainers(t *testing.T) {
	configs := []types.ServiceConfig{
		{Name: "service1", Container: types.ContainerSpec{Restart: types.RestartPolicyAlways}},
		{Name: "service2", Container: types.ContainerSpec{Restart: types.RestartPolicyNo}},
		{Name: "service3", Container: types.ContainerSpec{Restart: types.RestartPolicyOnFailure}},
	}

	result := filterInitContainers(configs)
	assert.Len(t, result, 2)
	assert.Contains(t, result, "service1")
	assert.Contains(t, result, "service3")
	assert.NotContains(t, result, "service2")
}
