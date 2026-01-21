package core

import (
	"fmt"
	"os"
	"path/filepath"
)

// Application constants
const (
	AppName      = "otto-stack"
	AppNameTitle = "Otto Stack"
	GitHubOrg    = "otto-nation"
	GitHubRepo   = "otto-stack"
)

// Exit codes
const (
	ExitSuccess = 0
	ExitError   = 1
)

// Timeout constants
const (
	DefaultStartTimeoutSeconds = 30
	DefaultHTTPTimeoutSeconds  = 5
)

// Threshold constants
const (
	HTTPOKStatusThreshold = 400
	MinArgumentCount      = 1
	DefaultLogTailLines   = "100"
	MinProjectNameLength  = 2
	MaxProjectNameLength  = 50
)

// File names
const (
	EnvGeneratedFileName = ".env.generated"
	EnvGeneratedFilePath = OttoStackDir + "/" + GeneratedDir + "/" + EnvGeneratedFileName
	ReadmeFileName       = "README.md"
	GitIgnoreFileName    = ".gitignore"
)

// URLs
const (
	DocsURL = "https://github.com/otto-nation/otto-stack/blob/main/docs"
)

// Actions
const (
	ActionProceed = "proceed"
	ActionBack    = "back"
)

// Prompts
const (
	PromptProjectName = "Enter project name"
	HelpProjectName   = "The name of your project (alphanumeric and hyphens only)"
)

// File permissions
const (
	PermReadWriteExec = 0755
	PermReadWrite     = 0644
)

// File extension constants
const (
	ExtYAML          = ".yaml"
	ExtYML           = ".yml"
	ExtJSON          = ".json"
	ExtENV           = ".env"
	ExtSH            = ".sh"
	ExtMD            = ".md"
	YMLFileExtension = ExtYML // Alias for compatibility
)

// Otto Stack directory and file constants
const (
	OttoStackDir        = ".otto-stack"
	SharedDir           = "shared"
	ConfigFileName      = "config.yaml"
	LocalConfigFileName = "config.local.yaml"
	ServiceConfigsDir   = "services"
	GeneratedDir        = "generated"
	LocalFileExtension  = ".local"
	SharedRegistryFile  = "containers.yaml"
)

// Container naming constants
const (
	SharedContainerPrefix = AppName + "-"
)

// Docker command constants
const (
	DockerCmdPs    = "ps"
	DockerCmdLogs  = "logs"
	DockerCmdExec  = "exec"
	DockerCmdStop  = "stop"
	DockerCmdStart = "start"
	DockerCmdRm    = "rm"
	DockerCmdPull  = "pull"
	DockerCmdBuild = "build"
)

// HTTP status constants
const (
	HTTPStatusOK                  = 200
	HTTPStatusBadRequest          = 400
	HTTPStatusUnauthorized        = 401
	HTTPStatusForbidden           = 403
	HTTPStatusNotFound            = 404
	HTTPStatusInternalServerError = 500
)

// Common environment variable names
const (
	EnvVarPATH            = "PATH"
	EnvVarHOME            = "HOME"
	EnvVarUSER            = "USER"
	EnvVarTERM            = "TERM"
	EnvOttoNonInteractive = "OTTO_NON_INTERACTIVE"
)

// FindYAMLFile finds a YAML file with the given name in the specified directory
// It checks for both .yaml and .yml extensions
func FindYAMLFile(dir, filename string) (string, error) {
	// Try .yaml first
	yamlPath := filepath.Join(dir, filename+ExtYAML)
	if _, err := os.Stat(yamlPath); err == nil {
		return yamlPath, nil
	}

	// Try .yml
	ymlPath := filepath.Join(dir, filename+ExtYML)
	if _, err := os.Stat(ymlPath); err == nil {
		return ymlPath, nil
	}

	return "", fmt.Errorf("YAML file not found: %s", filename)
}

// IsYAMLFile checks if a filename has a YAML extension
func IsYAMLFile(filename string) bool {
	ext := filepath.Ext(filename)
	return ext == ExtYAML || ext == ExtYML
}

// TrimYAMLExt removes the YAML extension from a filename
func TrimYAMLExt(filename string) string {
	if filepath.Ext(filename) == ExtYAML {
		return filename[:len(filename)-len(ExtYAML)]
	}
	if filepath.Ext(filename) == ExtYML {
		return filename[:len(filename)-len(ExtYML)]
	}
	return filename
}
