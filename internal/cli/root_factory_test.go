//go:build unit

package cli

import (
	"testing"
)

func TestRootFactory_Execute(t *testing.T) {
	err := ExecuteFactory()
	_ = err
}

func TestRootFactory_InitConfig(t *testing.T) {
	initConfig()
}

func TestRootFactory_SetupViper(t *testing.T) {
	setupViper()
}

func TestRootFactory_ConfigureLogger(t *testing.T) {
	configureLogger()
}
