package docker

import (
	pkgCore "github.com/otto-nation/otto-stack/internal/core"
)

// Docker constants
const (
	DockerCmd                        = "docker"
	ComposeProjectLabel              = "com.docker.compose.project"
	ComposeServiceLabel              = "com.docker.compose.service"
	DockerInstallURL                 = "https://docs.docker.com/get-docker/"
	DockerComposeFileName            = "docker-compose.yml"
	DockerComposeFileNameYaml        = "docker-compose.yaml"
	DockerComposeFilePath            = pkgCore.OttoStackDir + "/" + DockerComposeFileName
	ServiceCharacteristicsConfigPath = "internal/config/docker/service_characteristics.yaml"
	FlagPrefix                       = "--"
)

// Docker commands
const (
	DockerInfoCmd    = "info"
	DockerComposeCmd = "compose"
	DockerVersionCmd = "version"
	DockerExecCmd    = "exec"
	DockerPsCmd      = "ps"
)

// Docker Compose commands
const (
	ComposeUpCmd   = "up"
	ComposeDownCmd = "down"
	ComposeLogsCmd = "logs"
	ComposeExecCmd = "exec"
	ComposeStopCmd = "stop"
)

// Docker flags
const (
	FlagProjectName   = "project-name"
	FlagFile          = "file"
	FlagDetach        = "detach"
	FlagBuild         = "build"
	FlagForceRecreate = "force-recreate"
	FlagRemoveOrphans = "remove-orphans"
	FlagVolumes       = "volumes"
	FlagTimeout       = "timeout"
	FlagUser          = "user"
	FlagWorkdir       = "workdir"
	FlagFilter        = "filter"
	FlagFormat        = "format"
	FlagFollow        = "follow"
	FlagTimestamps    = "timestamps"
	FlagTail          = "tail"
	FlagSince         = "since"
	FlagHealthTimeout = "health-timeout"
	FlagHealthRetries = "health-retries"
	FlagRm            = "rm"
	FlagName          = "name"
	FlagEnv           = "env"
	FlagVolume        = "volume"
	FlagWorkingDir    = "workdir"
	FlagNetwork       = "network"
)

// Shell commands
const (
	ShellSh = "sh"
	ShellC  = "-c"
)

// Container paths
const (
	ContainerConfigPath = "/config"
	ContainerRootPath   = "/"
)

// Docker labels
const (
	LabelComposeService = "com.docker.compose.service"
)

// Docker commands
const (
	DockerRunCmd  = "run"
	DockerStopCmd = "stop"
	DockerRmCmd   = "rm"
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
	NetworkNameSuffix  = "-network"
)

// Init container images
const (
	AlpineLatestImage = "alpine:latest"
)

// Init container constants
const (
	InitServiceEndpointURL = "SERVICE_ENDPOINT_URL"
	InitServiceName        = "INIT_SERVICE_NAME"
	InitConfigDir          = "INIT_CONFIG_DIR"
)
