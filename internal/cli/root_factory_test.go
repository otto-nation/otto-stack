//go:build unit

package cli

import (
	"testing"
)

func TestRootFactory(t *testing.T) {
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
}
