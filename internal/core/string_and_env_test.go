//go:build unit

package core

import (
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
