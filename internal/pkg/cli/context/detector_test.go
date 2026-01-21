package context

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/otto-nation/otto-stack/internal/core"
)

func TestDetector_Detect_GlobalContext(t *testing.T) {
	// Create temp directory for test
	tmpDir := t.TempDir()

	// Change to temp directory
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	// Create detector with temp home
	detector := &Detector{homeDir: tmpDir}

	ctx, err := detector.Detect()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ctx.Type != Global {
		t.Errorf("expected Global context, got %s", ctx.Type)
	}
	if ctx.Project != nil {
		t.Error("expected nil Project in global context")
	}
	if ctx.Shared == nil {
		t.Fatal("expected Shared to be populated")
	}

	// Verify shared directory was created
	expectedShared := filepath.Join(tmpDir, core.OttoStackDir, core.SharedDir)
	if ctx.Shared.Root != expectedShared {
		t.Errorf("expected shared root %s, got %s", expectedShared, ctx.Shared.Root)
	}
	if _, err := os.Stat(expectedShared); os.IsNotExist(err) {
		t.Error("shared directory was not created")
	}
}

func TestDetector_Detect_ProjectContext(t *testing.T) {
	// Create temp directory structure
	tmpDir := t.TempDir()
	projectDir := filepath.Join(tmpDir, "myproject")
	configDir := filepath.Join(projectDir, ".otto-stack")
	configFile := filepath.Join(configDir, "config.yaml")

	os.MkdirAll(configDir, core.PermReadWriteExec)
	os.WriteFile(configFile, []byte("project_name: test\n"), core.PermReadWrite)

	// Change to project directory
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(projectDir)

	// Create detector with temp home
	detector := &Detector{homeDir: tmpDir}

	ctx, err := detector.Detect()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ctx.Type != Project {
		t.Errorf("expected Project context, got %s", ctx.Type)
	}
	if ctx.Project == nil {
		t.Fatal("expected Project to be populated")
	}

	// Resolve symlinks for comparison (macOS /var -> /private/var)
	expectedRoot, _ := filepath.EvalSymlinks(projectDir)
	actualRoot, _ := filepath.EvalSymlinks(ctx.Project.Root)
	if actualRoot != expectedRoot {
		t.Errorf("expected project root %s, got %s", expectedRoot, actualRoot)
	}

	expectedConfigDir, _ := filepath.EvalSymlinks(configDir)
	actualConfigDir, _ := filepath.EvalSymlinks(ctx.Project.ConfigDir)
	if actualConfigDir != expectedConfigDir {
		t.Errorf("expected config dir %s, got %s", expectedConfigDir, actualConfigDir)
	}

	expectedConfigFile, _ := filepath.EvalSymlinks(configFile)
	actualConfigFile, _ := filepath.EvalSymlinks(ctx.Project.ConfigFile)
	if actualConfigFile != expectedConfigFile {
		t.Errorf("expected config file %s, got %s", expectedConfigFile, actualConfigFile)
	}

	if ctx.Shared == nil {
		t.Fatal("expected Shared to be populated")
	}
}

func TestDetector_Detect_ProjectContextInSubdirectory(t *testing.T) {
	// Create temp directory structure
	tmpDir := t.TempDir()
	projectDir := filepath.Join(tmpDir, "myproject")
	configDir := filepath.Join(projectDir, ".otto-stack")
	configFile := filepath.Join(configDir, "config.yaml")
	subDir := filepath.Join(projectDir, "src", "app")

	os.MkdirAll(configDir, core.PermReadWriteExec)
	os.WriteFile(configFile, []byte("project_name: test\n"), core.PermReadWrite)
	os.MkdirAll(subDir, core.PermReadWriteExec)

	// Change to subdirectory
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(subDir)

	// Create detector with temp home
	detector := &Detector{homeDir: tmpDir}

	ctx, err := detector.Detect()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ctx.Type != Project {
		t.Errorf("expected Project context, got %s", ctx.Type)
	}
	if ctx.Project == nil {
		t.Fatal("expected Project to be populated")
	}

	// Resolve symlinks for comparison
	expectedRoot, _ := filepath.EvalSymlinks(projectDir)
	actualRoot, _ := filepath.EvalSymlinks(ctx.Project.Root)
	if actualRoot != expectedRoot {
		t.Errorf("expected project root %s, got %s", expectedRoot, actualRoot)
	}
}

func TestDetector_findProjectRoot_NoProject(t *testing.T) {
	tmpDir := t.TempDir()

	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	detector := &Detector{homeDir: tmpDir}

	info, err := detector.findProjectRoot()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info != nil {
		t.Error("expected nil when no project found")
	}
}

func TestDetectContext(t *testing.T) {
	// Just verify it doesn't error
	ctx, err := DetectContext()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ctx == nil {
		t.Fatal("expected context to be returned")
	}
	if ctx.Shared == nil {
		t.Error("expected Shared to be populated")
	}
}
