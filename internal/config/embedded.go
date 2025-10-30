package config

import (
	"embed"
)

//go:embed commands.yaml
var EmbeddedCommandsYAML []byte

//go:embed schema.yaml
var EmbeddedSchemaYAML []byte

//go:embed init-settings.yaml
var EmbeddedInitSettingsYAML []byte

//go:embed services
var EmbeddedServicesFS embed.FS
