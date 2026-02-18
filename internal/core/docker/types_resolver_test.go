package docker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServiceCharacteristicsResolver(t *testing.T) {
	resolver, err := NewServiceCharacteristicsResolver()
	assert.NoError(t, err)
	assert.NotNil(t, resolver)
	assert.NotNil(t, resolver.config)
}

func TestLoadServiceCharacteristicsConfig(t *testing.T) {
	config, err := loadServiceCharacteristicsConfig()
	assert.NoError(t, err)
	assert.NotNil(t, config)
}
