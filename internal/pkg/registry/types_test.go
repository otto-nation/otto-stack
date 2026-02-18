//go:build unit

package registry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegistry_FindOrphans_EmptyProjects(t *testing.T) {
	registry := &Registry{
		Containers: map[string]*ContainerInfo{
			"orphan": {Name: "otto-stack-orphan", Projects: []string{}},
		},
	}
	orphans := registry.FindOrphans()
	assert.Len(t, orphans, 1)
	assert.Equal(t, "orphan", orphans[0].Service)
}

func TestRegistry_FindOrphans_NoOrphans(t *testing.T) {
	registry := &Registry{
		Containers: map[string]*ContainerInfo{
			"active": {Name: "otto-stack-active", Projects: []string{"p1"}},
		},
	}
	orphans := registry.FindOrphans()
	assert.Empty(t, orphans)
}

func TestRegistry_FindOrphans_NilProjects(t *testing.T) {
	registry := &Registry{
		Containers: map[string]*ContainerInfo{
			"orphan": {Name: "otto-stack-orphan", Projects: nil},
		},
	}
	orphans := registry.FindOrphans()
	assert.Len(t, orphans, 1)
}

func TestContainerInfo(t *testing.T) {
	info := &ContainerInfo{
		Name:     "otto-stack-test",
		Projects: []string{"p1", "p2"},
	}
	assert.Equal(t, "otto-stack-test", info.Name)
	assert.Len(t, info.Projects, 2)
}

func TestOrphanInfo(t *testing.T) {
	orphan := OrphanInfo{
		Service:   "test",
		Container: "otto-stack-test",
		Severity:  OrphanSeverityCritical,
		Reason:    "test reason",
	}
	assert.Equal(t, "test", orphan.Service)
	assert.Equal(t, OrphanSeverityCritical, orphan.Severity)
}
