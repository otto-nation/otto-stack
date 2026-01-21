package context

import (
	"os"
	"path/filepath"

	"github.com/otto-nation/otto-stack/internal/core"
)

// Detector detects the current execution context
type Detector struct {
	homeDir string
}

// NewDetector creates a new context detector
func NewDetector() (*Detector, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	return &Detector{homeDir: home}, nil
}

// Detect determines the current execution context
func (d *Detector) Detect() (*ExecutionContext, error) {
	sharedRoot := filepath.Join(d.homeDir, core.OttoStackDir, core.SharedDir)
	if err := os.MkdirAll(sharedRoot, core.PermReadWriteExec); err != nil {
		return nil, err
	}

	shared := &SharedInfo{Root: sharedRoot}
	projectInfo, err := d.findProjectRoot()
	if err != nil {
		return nil, err
	}

	if projectInfo != nil {
		return &ExecutionContext{
			Type:    Project,
			Project: projectInfo,
			Shared:  shared,
		}, nil
	}

	return &ExecutionContext{
		Type:   Global,
		Shared: shared,
	}, nil
}

// findProjectRoot walks up the directory tree to find .otto-stack
func (d *Detector) findProjectRoot() (*ProjectInfo, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	for dir := cwd; ; dir = filepath.Dir(dir) {
		configDir := filepath.Join(dir, core.OttoStackDir)
		configFile := filepath.Join(configDir, core.ConfigFileName)

		if info, err := os.Stat(configDir); err == nil && info.IsDir() {
			if _, err := os.Stat(configFile); err == nil {
				return NewProjectInfo(configDir), nil
			}
		}

		if parent := filepath.Dir(dir); parent == dir {
			break // Reached root
		}
	}

	return nil, nil
}
