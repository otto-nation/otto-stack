package cli

import (
	"fmt"
	"os"

	"github.com/otto-nation/otto-stack/internal/pkg/cli"
	"github.com/otto-nation/otto-stack/internal/pkg/config"
	"github.com/otto-nation/otto-stack/internal/pkg/constants"
	"github.com/otto-nation/otto-stack/internal/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// CreateRootCommand creates the root command using the simplified builder
func CreateRootCommand() (*cobra.Command, error) {
	rootCmd, err := cli.BuildRootCommand()
	if err != nil {
		return nil, fmt.Errorf("failed to build root command: %w", err)
	}

	cobra.OnInitialize(initFactoryConfig)

	return rootCmd, nil
}

// ExecuteFactory executes the root command using the functional builder
func ExecuteFactory() error {
	rootCmd, err := CreateRootCommand()
	if err != nil {
		return fmt.Errorf("failed to create CLI: %w", err)
	}

	return rootCmd.Execute()
}

// initFactoryConfig reads in config file and ENV variables if set
func initFactoryConfig() {
	var cfgFile string
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config only in the Otto Stack directory
		viper.AddConfigPath(constants.OttoStackDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName(constants.ConfigFileName[:len(constants.ConfigFileName)-4]) // Remove .yml extension
	}

	// read in environment variables that match
	viper.AutomaticEnv()

	// If a config file is found, read it in
	if err := viper.ReadInConfig(); err == nil {
		if viper.GetBool("verbose") {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
	}

	// Configure logger based on flags
	configureLogger()
}

// configureLogger sets up logger based on command line flags
func configureLogger() {
	config := logger.DefaultConfig()

	if viper.GetBool("verbose") {
		config.Level = logger.LevelInfo
	} else {
		config.Level = logger.LevelWarn
	}

	// Reinitialize logger with new config
	if err := logger.Init(config); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to configure logger: %v\n", err)
	}
}

// GetCommandConfig loads and returns the command configuration
func GetCommandConfig() (*config.CommandConfig, error) {
	loader := config.NewLoader("")
	return loader.Load()
}

// ValidateConfig validates the current configuration
func ValidateConfig() error {
	commandConfig, err := GetCommandConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	result := commandConfig.Validate()
	if !result.Valid {
		logger.Error("Configuration validation failed", "error_count", len(result.Errors))
		for _, err := range result.Errors {
			logger.Error("Validation error", "field", err.Field, "message", err.Message)
		}
		return fmt.Errorf("configuration validation failed")
	}

	logger.Info("Configuration is valid")
	if len(result.Warnings) > 0 {
		logger.Warn("Configuration warnings", "warning_count", len(result.Warnings))
		for _, warning := range result.Warnings {
			logger.Warn("Validation warning", "field", warning.Field, "message", warning.Message)
		}
	}

	return nil
}
