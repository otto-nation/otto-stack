package project

import (
	"github.com/otto-nation/otto-stack/internal/core"
)

// Test constants to eliminate magic strings and provide context
const (
	// Test directory patterns
	TestTempDirPattern = core.AppName + "-test-*"

	// Test project names
	TestProjectName        = "test-project"
	TestProjectNameValid   = "valid-project"
	TestProjectNameInvalid = "invalid@project"

	// Test file content
	TestConfigContent    = "test: config"
	TestReadmeContent    = "# Test Project"
	TestGitignoreContent = "*.bak\n*.tmp"
	TestExistingContent  = "# Existing content"

	// Test validation messages
	MsgAlreadyInitialized = "already initialized"
	MsgRequiredTool       = "required tool"
	MsgInvalidService     = "invalid service"
	MsgDuplicateService   = "duplicate service"
)

// Use constants from the constants package
const (
	// Test project types
	TestProjectTypeDocker = "docker"

	// Test gitignore entries (use actual constants)
	TestGitignoreEntry = core.OttoStackDir + "/" + core.EnvGeneratedFileName
)

// Test CLI commands
var (
	CmdDevStackUp     = core.AppName + " up"
	CmdDevStackDown   = core.AppName + " down"
	CmdDevStackStatus = core.AppName + " status"
)

// Test file paths (use actual constants for consistency)
var (
	TestConfigFilePath = core.OttoStackDir + "/" + core.ConfigFileName
	TestReadmeFilePath = core.OttoStackDir + "/" + core.ReadmeFileName
)
