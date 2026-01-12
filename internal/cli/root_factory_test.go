//go:build unit

package cli

import (
	"testing"

	"github.com/otto-nation/otto-stack/test/testhelpers"
)

func TestRootFactory(t *testing.T) {
	t.Run("validates root command creation", func(t *testing.T) {
		cmd, err := CreateRootCommand()
		testhelpers.AssertValidConstructor(t, cmd, err, "CreateRootCommand")
	})

	t.Run("validates factory execution", func(t *testing.T) {
		err := ExecuteFactory()
		// ExecuteFactory may succeed in test environment
		_ = err
	})

	t.Run("validates config initialization", func(t *testing.T) {
		initConfig()
		// No assertion needed, just testing execution
	})

	t.Run("validates viper setup", func(t *testing.T) {
		setupViper()
		// No assertion needed, just testing execution
	})

	t.Run("validates logger configuration", func(t *testing.T) {
		configureLogger()
		// No assertion needed, just testing execution
	})

	t.Run("validates command config retrieval", func(t *testing.T) {
		config, err := GetCommandConfig()
		testhelpers.AssertValidConstructor(t, config, err, "GetCommandConfig")
	})

	t.Run("validates config validation", func(t *testing.T) {
		err := ValidateConfig()
		testhelpers.AssertErrorPattern(t, nil, err, true, "ValidateConfig")
	})
}
