package config

// ConfigService provides configuration management functionality
type ConfigService interface {
	LoadConfig() (*Config, error)
	ValidateConfig(cfg *Config) error
	GetConfigHash(cfg *Config) (string, error)
}
