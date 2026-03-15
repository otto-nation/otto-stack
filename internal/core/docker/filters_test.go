package docker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProjectFilter(t *testing.T) {
	t.Run("with project name applies label filter", func(t *testing.T) {
		filter := NewProjectFilter("test-project")
		assert.True(t, filter.Len() > 0)
	})

	t.Run("with empty project name returns no filter", func(t *testing.T) {
		filter := NewProjectFilter("")
		assert.Equal(t, 0, filter.Len())
	})
}
