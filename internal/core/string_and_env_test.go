//go:build unit

package core

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTitleCase(t *testing.T) {
	t.Run("converts hyphenated string to title case", func(t *testing.T) {
		result := TitleCase("hello-world")
		assert.Equal(t, "Hello World", result)
	})

	t.Run("handles single word", func(t *testing.T) {
		result := TitleCase("hello")
		assert.Equal(t, "Hello", result)
	})

	t.Run("handles empty string", func(t *testing.T) {
		result := TitleCase("")
		assert.Equal(t, "", result)
	})

	t.Run("handles multiple hyphens", func(t *testing.T) {
		result := TitleCase("hello-world-test")
		assert.Equal(t, "Hello World Test", result)
	})

	t.Run("handles string without hyphens", func(t *testing.T) {
		result := TitleCase("hello")
		assert.Equal(t, "Hello", result)
	})
}

func TestResolveVar(t *testing.T) {
	t.Run("resolves environment variable with default", func(t *testing.T) {
		// Set a test environment variable
		os.Setenv("TEST_VAR", "test_value")
		defer os.Unsetenv("TEST_VAR")

		result := ResolveVar("${TEST_VAR:-default}")
		assert.Equal(t, "test_value", result)
	})

	t.Run("uses default when env var not set", func(t *testing.T) {
		// Ensure env var is not set
		os.Unsetenv("NONEXISTENT_VAR")

		result := ResolveVar("${NONEXISTENT_VAR:-default_value}")
		assert.Equal(t, "default_value", result)
	})

	t.Run("uses default when env var is empty", func(t *testing.T) {
		// Set empty environment variable
		os.Setenv("EMPTY_VAR", "")
		defer os.Unsetenv("EMPTY_VAR")

		result := ResolveVar("${EMPTY_VAR:-default_value}")
		assert.Equal(t, "default_value", result)
	})

	t.Run("returns original string if not env var syntax", func(t *testing.T) {
		result := ResolveVar("regular_string")
		assert.Equal(t, "regular_string", result)
	})

	t.Run("handles malformed env var syntax", func(t *testing.T) {
		result := ResolveVar("${MALFORMED")
		assert.Equal(t, "${MALFORMED", result)

		result2 := ResolveVar("MALFORMED}")
		assert.Equal(t, "MALFORMED}", result2)
	})

	t.Run("handles env var without default", func(t *testing.T) {
		result := ResolveVar("${VAR_NO_DEFAULT}")
		assert.Equal(t, "${VAR_NO_DEFAULT}", result)
	})
}
