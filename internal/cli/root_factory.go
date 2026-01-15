package cli

import (
	"fmt"
	"os"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/cli"
	"github.com/otto-nation/otto-stack/internal/pkg/config"
	"github.com/otto-nation/otto-stack/internal/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	// Import handlers to trigger registration
	_ "github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/lifecycle"
	_ "github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/operations"
	_ "github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/project"
	_ "github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/utility"
)

// ExecuteFactory executes the root command using the functional builder
func ExecuteFactory() error {
	rootCmd := cli.BuildRootCommand()
	cobra.OnInitialize(initConfig)
	return rootCmd.Execute()
}

// CreateRootCommand creates the root command using the simplified builder
func CreateRootCommand() (*cobra.Command, error) {
	return cli.BuildRootCommand(), nil
}

// initConfig reads in config file and ENV variables if set
func initConfig() {
	setupViper()
	configureLogger()
}

func setupViper() {
	viper.AddConfigPath(core.OttoStackDir)
	viper.SetConfigType("yaml")
	viper.SetConfigName(core.ConfigFileName[:len(core.ConfigFileName)-4]) // Remove .yml extension
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil && viper.GetBool("verbose") {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

// configureLogger sets up logger based on command line flags
func configureLogger() {
	config := logger.DefaultConfig()
	if viper.GetBool("verbose") {
		config.Level = logger.LogLevelDebug
	} else {
		config.Level = logger.LogLevelWarn
	}

	if err := logger.Init(config); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to configure logger: %v\n", err)
	}
}

// GetCommandConfig loads and returns the command configuration
func GetCommandConfig() (map[string]any, error) {
	return config.LoadCommandConfig()
}

// ValidateConfig validates the current configuration
func ValidateConfig() error {
	if _, err := config.LoadConfig(); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}
	logger.Info("Configuration is valid")
	return nil
}
