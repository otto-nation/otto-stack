package constants

// File names
const (
	ConfigFileName           = "otto-stack-config.yml"
	ConfigFileNameYAML       = "otto-stack-config.yaml"
	ConfigFileNameHidden     = ".otto-stack-config.yml"
	ConfigFileNameHiddenYAML = ".otto-stack-config.yaml"
	DockerComposeFileName    = "docker-compose.yml"
	EnvGeneratedFileName     = ".env.generated"
	GitignoreFileName        = ".gitignore"
	ReadmeFileName           = "README.md"
	ServiceConfigExtension   = ".yaml"
)

// Directory names
const (
	DevStackDir = "otto-stack"
	DataDir     = "data"
	LogsDir     = "logs"
	TmpDir      = "tmp"
	ServicesDir = "internal/config/services"
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
	DevStackDir + "/" + EnvGeneratedFileName,
	DevStackDir + "/" + DataDir + "/",
	DevStackDir + "/" + LogsDir + "/",
	DevStackDir + "/" + TmpDir + "/",
}
