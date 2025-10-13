package init

import "github.com/otto-nation/otto-stack/internal/pkg/constants"

// Test constants to eliminate magic strings and provide context
const (
	// Test directory patterns
	TestTempDirPattern = "otto-stack-test-*"

	// Test project names
	TestProjectName        = "test-project"
	TestProjectNameValid   = "valid-project"
	TestProjectNameInvalid = "invalid@project"

	// Test services
	TestServicePostgres = "postgres"
	TestServiceRedis    = "redis"
	TestServiceNginx    = "nginx"

	// Test file content
	TestConfigContent    = "test: config"
	TestReadmeContent    = "# Test Project"
	TestGitignoreContent = "*.log\n*.tmp"
	TestExistingContent  = "# Existing content"

	// Test validation messages
	MsgAlreadyInitialized = "already initialized"
	MsgRequiredTool       = "required tool"
	MsgNoServicesSelected = "no services selected"
	MsgInvalidService     = "invalid service"
	MsgDuplicateService   = "duplicate service"
)

// Use constants from the constants package
const (
	// Test environments (use actual constants)
	TestEnvironmentLocal = constants.DefaultEnvironment
	TestEnvironmentDev   = "development"
	TestEnvironmentProd  = "production"

	// Test gitignore entries (use actual constants)
	TestGitignoreEntry = constants.DevStackDir + "/" + constants.EnvGeneratedFileName
)

// Test CLI commands (use variables since they're computed)
var (
	CmdDevStackUp     = constants.CmdUp
	CmdDevStackDown   = constants.CmdDown
	CmdDevStackStatus = constants.CmdStatus
)

// Test file paths (use actual constants for consistency)
var (
	TestConfigFilePath     = constants.DevStackDir + "/" + constants.ConfigFileName
	TestConfigFilePathYAML = constants.DevStackDir + "/" + constants.ConfigFileNameYAML
	TestReadmeFilePath     = constants.DevStackDir + "/" + constants.ReadmeFileName
)
