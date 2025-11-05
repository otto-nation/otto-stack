package utils

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/otto-nation/otto-stack/internal/pkg/constants"
)

// CheckInitialization verifies that otto-stack is initialized
func CheckInitialization() error {
	configPath := filepath.Join(constants.OttoStackDir, constants.ConfigFileName)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return errors.New(constants.MsgErrors_not_initialized)
	}
	return nil
}
