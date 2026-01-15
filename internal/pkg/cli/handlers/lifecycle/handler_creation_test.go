//go:build unit

package lifecycle

import (
	"testing"

	"github.com/otto-nation/otto-stack/test/testhelpers"
)

func TestNewCleanupHandler(t *testing.T) {
	handler := NewCleanupHandler()
	testhelpers.AssertValidConstructor(t, handler, nil, "CleanupHandler")
}

func TestNewRestartHandler(t *testing.T) {
	handler := NewRestartHandler()
	testhelpers.AssertValidConstructor(t, handler, nil, "RestartHandler")
}
