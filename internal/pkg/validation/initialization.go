package validation

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/messages"
)

// CheckInitialization verifies that otto-stack is initialized
func CheckInitialization() error {
	configPath := filepath.Join(core.OttoStackDir, core.ConfigFileName)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return errors.New(messages.ErrorsNotInitialized)
	}
	return nil
}
