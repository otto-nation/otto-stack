package constants

// File names
const (
	ConfigFileName         = AppName + "-config.yml"
	LocalConfigFileName    = AppName + "-config.local.yml"
	DockerComposeFileName  = "docker-compose.yml"
	EnvGeneratedFileName   = ".env.generated"
	GitignoreFileName      = ".gitignore"
	ReadmeFileName         = "README.md"
	ServiceConfigExtension = ".yaml"
	KafkaTopicsInitScript  = "kafka-topics-init.sh"
	LocalstackInitScript   = "localstack-init.sh"
	StateFileName          = "state.json"
)

// Directory names
const (
	OttoStackDir        = AppName
	DataDir             = "data"
	LogsDir             = "logs"
	TmpDir              = "tmp"
	ScriptsDir          = "scripts"
	ServicesDir         = "internal/config/services"
	EmbeddedServicesDir = "services" // Directory name in embedded FS
)

// Configuration URLs
const (
	ConfigDocsURL    = "https://github.com/otto-nation/otto-stack/tree/main/docs-site/content/configuration.md"
	ServiceConfigURL = "https://github.com/otto-nation/otto-stack/tree/main/internal/config/services"
)

// Git entries
var GitignoreEntries = []string{
	"",
	"# Otto Stack",
	OttoStackDir + "/" + EnvGeneratedFileName,
	OttoStackDir + "/" + DataDir + "/",
	OttoStackDir + "/" + LogsDir + "/",
	OttoStackDir + "/" + TmpDir + "/",
	"",
	"# Local config overrides",
	OttoStackDir + "/" + LocalConfigFileName,
}
