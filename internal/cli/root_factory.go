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

// CreateRootCommand creates the root command using the functional builder
func CreateRootCommand() (*cobra.Command, error) {
	loader := config.NewLoader("")
	commandConfig, err := loader.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load command configuration: %w", err)
	}

	validationResult := commandConfig.Validate()
	if !validationResult.Valid {
		fmt.Fprintf(os.Stderr, "Warning: Command configuration has validation errors:\n")
		for _, err := range validationResult.Errors {
			fmt.Fprintf(os.Stderr, "  - %s: %s\n", err.Field, err.Message)
		}
		if len(validationResult.Warnings) > 0 {
			fmt.Fprintf(os.Stderr, "Warnings:\n")
			for _, warning := range validationResult.Warnings {
				fmt.Fprintf(os.Stderr, "  - %s: %s\n", warning.Field, warning.Message)
			}
		}
	}

	rootCmd, err := cli.BuildRootCommand(commandConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to build root command: %w", err)
	}

	cobra.OnInitialize(func() {
		initFactoryConfig(commandConfig)
	})

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
func initFactoryConfig(_ *config.CommandConfig) {
	var cfgFile string
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in multiple locations
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.AddConfigPath(".otto-stack")
		viper.SetConfigType("yaml")

		// Try to find config file with multiple names
		configNames := []string{constants.AppName + "-config", "." + constants.AppName}
		var configFound bool
		for _, name := range configNames {
			viper.SetConfigName(name)
			if err := viper.ReadInConfig(); err == nil {
				configFound = true
				break
			}
		}

		// If no config found, don't call ReadInConfig again
		if configFound {
			return
		}
	}

	// read in environment variables that match
	viper.AutomaticEnv()

	// If a config file is found, read it in (only if not already read above)
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
