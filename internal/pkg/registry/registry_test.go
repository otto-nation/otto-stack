package registry

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
	"time"

	"github.com/otto-nation/otto-stack/internal/core"
)

func TestNewRegistry(t *testing.T) {
	reg := NewRegistry()
	if reg == nil {
		t.Fatal("expected registry to be created")
	}
	if reg.Containers == nil {
		t.Error("expected Containers map to be initialized")
	}
	if len(reg.Containers) != 0 {
		t.Error("expected empty Containers map")
	}
}

func TestManager_LoadEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	mgr := NewManager(tmpDir)

	registry, err := mgr.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if registry == nil {
		t.Fatal("expected registry to be created")
	}
	if len(registry.Containers) != 0 {
		t.Error("expected empty registry")
	}
}

func TestManager_SaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	mgr := NewManager(tmpDir)

	registry := NewRegistry()
	registry.Containers["postgres"] = &ContainerInfo{
		Name:      "otto-stack-postgres",
		Service:   "postgres",
		Projects:  []string{"project1"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := mgr.Save(registry); err != nil {
		t.Fatalf("unexpected error saving: %v", err)
	}

	loaded, err := mgr.Load()
	if err != nil {
		t.Fatalf("unexpected error loading: %v", err)
	}

	if len(loaded.Containers) != 1 {
		t.Errorf("expected 1 container, got %d", len(loaded.Containers))
	}

	container := loaded.Containers["postgres"]
	if container == nil {
		t.Fatal("expected postgres container")
	}
	if container.Name != "otto-stack-postgres" {
		t.Errorf("expected name otto-stack-postgres, got %s", container.Name)
	}
	if len(container.Projects) != 1 || container.Projects[0] != "project1" {
		t.Errorf("expected projects [project1], got %v", container.Projects)
	}
}

func TestManager_Register(t *testing.T) {
	tmpDir := t.TempDir()
	mgr := NewManager(tmpDir)

	err := mgr.Register("postgres", "otto-stack-postgres", "project1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	container, err := mgr.Get("postgres")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if container == nil {
		t.Fatal("expected container to exist")
	}
	if container.Name != "otto-stack-postgres" {
		t.Errorf("expected name otto-stack-postgres, got %s", container.Name)
	}
	if len(container.Projects) != 1 || container.Projects[0] != "project1" {
		t.Errorf("expected projects [project1], got %v", container.Projects)
	}
}

func TestManager_RegisterMultipleProjects(t *testing.T) {
	tmpDir := t.TempDir()
	mgr := NewManager(tmpDir)

	mgr.Register("postgres", "otto-stack-postgres", "project1")
	mgr.Register("postgres", "otto-stack-postgres", "project2")

	container, err := mgr.Get("postgres")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(container.Projects) != 2 {
		t.Errorf("expected 2 projects, got %d", len(container.Projects))
	}
	if !slices.Contains(container.Projects, "project1") {
		t.Error("expected project1 in projects")
	}
	if !slices.Contains(container.Projects, "project2") {
		t.Error("expected project2 in projects")
	}
}

func TestManager_RegisterDuplicate(t *testing.T) {
	tmpDir := t.TempDir()
	mgr := NewManager(tmpDir)

	mgr.Register("postgres", "otto-stack-postgres", "project1")
	mgr.Register("postgres", "otto-stack-postgres", "project1")

	container, err := mgr.Get("postgres")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(container.Projects) != 1 {
		t.Errorf("expected 1 project, got %d", len(container.Projects))
	}
}

func TestManager_Unregister(t *testing.T) {
	tmpDir := t.TempDir()
	mgr := NewManager(tmpDir)

	mgr.Register("postgres", "otto-stack-postgres", "project1")
	mgr.Register("postgres", "otto-stack-postgres", "project2")

	err := mgr.Unregister("postgres", "project1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	container, err := mgr.Get("postgres")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(container.Projects) != 1 {
		t.Errorf("expected 1 project, got %d", len(container.Projects))
	}
	if container.Projects[0] != "project2" {
		t.Errorf("expected project2, got %s", container.Projects[0])
	}
}

func TestManager_UnregisterLastProject(t *testing.T) {
	tmpDir := t.TempDir()
	mgr := NewManager(tmpDir)

	mgr.Register("postgres", "otto-stack-postgres", "project1")
	mgr.Unregister("postgres", "project1")

	container, err := mgr.Get("postgres")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if container != nil {
		t.Error("expected container to be removed")
	}
}

func TestManager_List(t *testing.T) {
	tmpDir := t.TempDir()
	mgr := NewManager(tmpDir)

	mgr.Register("postgres", "otto-stack-postgres", "project1")
	mgr.Register("redis", "otto-stack-redis", "project1")

	containers, err := mgr.List()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(containers) != 2 {
		t.Errorf("expected 2 containers, got %d", len(containers))
	}
}

func TestManager_IsShared(t *testing.T) {
	tmpDir := t.TempDir()
	mgr := NewManager(tmpDir)

	mgr.Register("postgres", "otto-stack-postgres", "project1")

	shared, err := mgr.IsShared("postgres")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !shared {
		t.Error("expected postgres to be shared")
	}

	shared, err = mgr.IsShared("redis")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if shared {
		t.Error("expected redis to not be shared")
	}
}

func TestManager_RegistryFilePath(t *testing.T) {
	tmpDir := t.TempDir()
	mgr := NewManager(tmpDir)

	mgr.Register("postgres", "otto-stack-postgres", "project1")

	expectedPath := filepath.Join(tmpDir, core.SharedRegistryFile)
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("expected registry file at %s", expectedPath)
	}
}
