//go:build unit

package docker

import (
	"testing"
)

func TestNewDefaultProjectLoader(t *testing.T) {
	loader, err := NewDefaultProjectLoader()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if loader == nil {
		t.Error("Expected non-nil loader")
	}
}
