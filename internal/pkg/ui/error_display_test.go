//go:build unit

package ui

import (
	"testing"
)

func TestOutput_Error_basic(t *testing.T) {
	output := &Output{}
	output.Error("test error")
	// Should not panic
}

func TestOutput_Muted_basic(t *testing.T) {
	output := &Output{}
	output.Muted("test muted")
	// Should not panic
}
