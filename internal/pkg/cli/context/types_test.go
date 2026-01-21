package context

import (
	"path/filepath"
	"testing"

	"github.com/otto-nation/otto-stack/internal/core"
)

func TestNewProjectInfo(t *testing.T) {
	configDir := filepath.Join("home", "user", "myproject", core.OttoStackDir)
	info := NewProjectInfo(configDir)

	expectedRoot := filepath.Join("home", "user", "myproject")
	if info.Root != expectedRoot {
		t.Errorf("expected root %s, got %s", expectedRoot, info.Root)
	}
	if info.ConfigDir != configDir {
		t.Errorf("expected configDir %s, got %s", configDir, info.ConfigDir)
	}

	expectedConfigFile := filepath.Join(configDir, core.ConfigFileName)
	if info.ConfigFile != expectedConfigFile {
		t.Errorf("expected configFile %s, got %s", expectedConfigFile, info.ConfigFile)
	}
}

func TestContext_IsProject(t *testing.T) {
	ctx := &ExecutionContext{Type: Project}
	if !ctx.IsProject() {
		t.Error("expected IsProject to return true")
	}
	if ctx.IsGlobal() {
		t.Error("expected IsGlobal to return false")
	}
}

func TestContext_IsGlobal(t *testing.T) {
	ctx := &ExecutionContext{Type: Global}
	if !ctx.IsGlobal() {
		t.Error("expected IsGlobal to return true")
	}
	if ctx.IsProject() {
		t.Error("expected IsProject to return false")
	}
}

func TestContext_GetProjectRoot(t *testing.T) {
	projectRoot := filepath.Join("home", "user", "project")
	tests := []struct {
		name     string
		ctx      *ExecutionContext
		expected string
	}{
		{
			name: "project context",
			ctx: &ExecutionContext{
				Type:    Project,
				Project: &ProjectInfo{Root: projectRoot},
			},
			expected: projectRoot,
		},
		{
			name:     "global context",
			ctx:      &ExecutionContext{Type: Global},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ctx.GetProjectRoot(); got != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, got)
			}
		})
	}
}

func TestContext_GetConfigFile(t *testing.T) {
	configFile := filepath.Join("home", "user", "project", core.OttoStackDir, core.ConfigFileName)
	tests := []struct {
		name     string
		ctx      *ExecutionContext
		expected string
	}{
		{
			name: "project context",
			ctx: &ExecutionContext{
				Type:    Project,
				Project: &ProjectInfo{ConfigFile: configFile},
			},
			expected: configFile,
		},
		{
			name:     "global context",
			ctx:      &ExecutionContext{Type: Global},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ctx.GetConfigFile(); got != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, got)
			}
		})
	}
}

func TestContext_GetSharedRoot(t *testing.T) {
	sharedRoot := filepath.Join("home", "user", core.OttoStackDir, core.SharedDir)
	ctx := &ExecutionContext{
		Shared: &SharedInfo{Root: sharedRoot},
	}
	if got := ctx.GetSharedRoot(); got != sharedRoot {
		t.Errorf("expected %s, got %s", sharedRoot, got)
	}
}
