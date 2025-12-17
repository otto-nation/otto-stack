package docker

import (
	pkgCore "github.com/otto-nation/otto-stack/internal/core"
)

// Docker constants
const (
	DockerCmd             = "docker"
	ComposeProjectLabel   = "com.docker.compose.project"
	ComposeServiceLabel   = "com.docker.compose.service"
	DockerInstallURL      = "https://docs.docker.com/get-docker/"
	DockerComposeFileName = "docker-compose.yml"
	DockerComposeFilePath = pkgCore.AppName + "/" + pkgCore.GeneratedDir + "/docker-compose.yml"
)

// Docker commands
const (
	DockerInfoCmd    = "docker info"
	DockerComposeCmd = "docker compose"
	DockerVersionCmd = "version"
)

// System command constants
const (
	CmdTaskkill = "taskkill"
	CmdLsof     = "lsof"
	CmdNetstat  = "netstat"
	CmdPgrep    = "pgrep"
	CmdTasklist = "tasklist"
)

// System error templates
const (
	ErrUnsupportedOS    = "unsupported OS: %s"
	ErrProcessNotFound  = "process %s not found"
	ErrFailedAfterRetry = "failed after %d attempts: %w"
	ErrOperationTimeout = "operation timed out after %v"
	ErrNoFreePort       = "no free port found in range %d-%d"
)

// System constants
const (
	MinFieldCount   = 2
	PortSearchRange = 1000
)

// OS constants
const (
	OSLinux   = "linux"
	OSDarwin  = "darwin"
	OSWindows = "windows"
)

// Docker Compose field names
const (
	ComposeFieldServices    = "services"
	ComposeFieldNetworks    = "networks"
	ComposeFieldName        = "name"
	ComposeFieldImage       = "image"
	ComposeFieldEntrypoint  = "entrypoint"
	ComposeFieldPorts       = "ports"
	ComposeFieldEnvironment = "environment"
	ComposeFieldVolumes     = "volumes"
	ComposeFieldRestart     = "restart"
	ComposeFieldCommand     = "command"
	ComposeFieldMemLimit    = "mem_limit"
	ComposeFieldHealthCheck = "healthcheck"
	ComposeFieldLabels      = "labels"
)

// Health check field names
const (
	HealthCheckFieldTest        = "test"
	HealthCheckFieldInterval    = "interval"
	HealthCheckFieldTimeout     = "timeout"
	HealthCheckFieldRetries     = "retries"
	HealthCheckFieldStartPeriod = "start_period"
)

// Docker volume and protocol constants
const (
	VolumeReadOnlySuffix = ":ro"
	ProtocolSeparator    = "/"
)

// State constants
const (
	StateRunning  = "running"
	StateStopped  = "exited"
	StateStarting = "starting"
	StateCreated  = "created"
	StatePaused   = "paused"
)

// Health status constants
const (
	HealthHealthy   = "healthy"
	HealthUnhealthy = "unhealthy"
	HealthStarting  = "starting"
	HealthNone      = "none"
)

// Otto Stack labels
const (
	LabelOttoManaged     = "io.otto-stack.managed"
	LabelOttoProject     = "io.otto-stack.project"
	LabelOttoService     = "io.otto-stack.service"
	LabelOttoVersion     = "io.otto-stack.version"
	LabelOttoSharingMode = "io.otto-stack.sharing-mode"
)

// Network names
const (
	DefaultNetworkName = "default"
)
