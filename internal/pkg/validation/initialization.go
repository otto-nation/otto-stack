package validation

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/otto-nation/otto-stack/internal/core"
)

// CheckInitialization verifies that otto-stack is initialized
func CheckInitialization() error {
	configPath := filepath.Join(core.OttoStackDir, core.ConfigFileName)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return errors.New(core.MsgErrors_not_initialized)
	}
	return nil
}
