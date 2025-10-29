package constants

// Command names
const (
	CmdNameUp         = "up"
	CmdNameDown       = "down"
	CmdNameRestart    = "restart"
	CmdNameStatus     = "status"
	CmdNameInit       = "init"
	CmdNameDoctor     = "doctor"
	CmdNameCompletion = "completion"
	CmdNameServices   = "services"
	CmdNameDeps       = "deps"
	CmdNameConflicts  = "conflicts"
	CmdNameLogs       = "logs"
	CmdNameExec       = "exec"
	CmdNameConnect    = "connect"
	CmdNameCleanup    = "cleanup"
	CmdNameValidate   = "validate"
)

// Service names for connection
const (
	ServicePostgres   = "postgres"
	ServicePostgreSQL = "postgresql"
	ServiceMySQL      = "mysql"
	ServiceRedis      = "redis"
	ServiceMongoDB    = "mongodb"
	ServiceMongo      = "mongo"
	ServiceLocalhost  = "localhost"
)

// Database client commands
const (
	ClientPSQL    = "psql"
	ClientMySQL   = "mysql"
	ClientRedis   = "redis-cli"
	ClientMongoDB = "mongosh"
)

// Database connection flags
const (
	PostgresUserFlag  = "-U"
	PostgresDBFlag    = "-d"
	PostgresHostFlag  = "-h"
	PostgresPortFlag  = "-p"
	MySQLUserFlag     = "-u"
	MySQLHostFlag     = "-h"
	MySQLPortFlag     = "-P"
	MySQLPasswordFlag = "-p"
	RedisHostFlag     = "-h"
	RedisPortFlag     = "-p"
	RedisDBFlag       = "-n"
)

// Default database users
const (
	DefaultPostgresUser = "postgres"
	DefaultMySQLUser    = "root"
)

// Flag names
const (
	FlagFollow       = "follow"
	FlagTail         = "tail"
	FlagTimestamps   = "timestamps"
	FlagInteractive  = "interactive"
	FlagTTY          = "tty"
	FlagUser         = "user"
	FlagWorkdir      = "workdir"
	FlagDatabase     = "database"
	FlagHost         = "host"
	FlagPort         = "port"
	FlagReadOnly     = "read-only"
	FlagAll          = "all"
	FlagVolumes      = "volumes"
	FlagImages       = "images"
	FlagNetworks     = "networks"
	FlagForce        = "force"
	FlagDryRun       = "dry-run"
	FlagOutput       = "output"
	FlagFormat       = "format"
	FlagFull         = "full"
	FlagCheckUpdates = "check-updates"
)

// ServiceConnectionConfig defines connection configuration for a service
type ServiceConnectionConfig struct {
	Client      string
	DefaultUser string
	UserFlag    string
	HostFlag    string
	PortFlag    string
	DBFlag      string
	ExtraFlags  []string
}

// ServiceConnections maps service names to their connection configurations
var ServiceConnections = map[string]ServiceConnectionConfig{
	ServicePostgres: {
		Client:      ClientPSQL,
		DefaultUser: DefaultPostgresUser,
		UserFlag:    PostgresUserFlag,
		HostFlag:    PostgresHostFlag,
		PortFlag:    PostgresPortFlag,
		DBFlag:      PostgresDBFlag,
	},
	ServicePostgreSQL: {
		Client:      ClientPSQL,
		DefaultUser: DefaultPostgresUser,
		UserFlag:    PostgresUserFlag,
		HostFlag:    PostgresHostFlag,
		PortFlag:    PostgresPortFlag,
		DBFlag:      PostgresDBFlag,
	},
	ServiceMySQL: {
		Client:      ClientMySQL,
		DefaultUser: DefaultMySQLUser,
		UserFlag:    MySQLUserFlag,
		HostFlag:    MySQLHostFlag,
		PortFlag:    MySQLPortFlag,
		ExtraFlags:  []string{MySQLPasswordFlag},
	},
	ServiceRedis: {
		Client:   ClientRedis,
		HostFlag: RedisHostFlag,
		PortFlag: RedisPortFlag,
		DBFlag:   RedisDBFlag,
	},
	ServiceMongoDB: {
		Client: ClientMongoDB,
	},
	ServiceMongo: {
		Client: ClientMongoDB,
	},
}

// Shell types for completion
const (
	ShellBash       = "bash"
	ShellZsh        = "zsh"
	ShellFish       = "fish"
	ShellPowerShell = "powershell"
)

// Docker commands
const (
	DockerCmd        = "docker"
	DockerInfoCmd    = "info"
	DockerComposeCmd = "compose"
	DockerVersionCmd = "version"
)

// URLs
const (
	DockerInstallURL = "https://docs.docker.com/get-docker/"
)

// Command reference builders
func CmdRef(cmdName string) string {
	return AppName + " " + cmdName
}

// Common command references
var (
	CmdUp     = CmdRef(CmdNameUp)
	CmdDown   = CmdRef(CmdNameDown)
	CmdStatus = CmdRef(CmdNameStatus)
	CmdInit   = CmdRef(CmdNameInit)
)

// Error messages
var (
	ErrNotInitialized = AppName + " not initialized. Run '" + CmdInit + "' first"
)

// Status formatting
const (
	StatusHeaderService = "SERVICE"
	StatusHeaderState   = "STATE"
	StatusHeaderHealth  = "HEALTH"
	StatusSeparator     = "-"
)
