package docker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProjectFilter(t *testing.T) {
	filter := NewProjectFilter("test-project")
	assert.NotNil(t, filter)
	assert.True(t, filter.Len() > 0)
}
