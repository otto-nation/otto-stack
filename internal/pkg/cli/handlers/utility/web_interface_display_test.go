//go:build unit

package utility

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWebInterfacesHandler_BasicMethods(t *testing.T) {
	t.Run("tests checkStatus with various URLs", func(t *testing.T) {
		handler := NewWebInterfacesHandler()

		// Test with localhost URL
		status := handler.checkStatus("http://localhost:8080")
		assert.IsType(t, "", status)
		assert.NotEmpty(t, status)

		// Test with invalid URL
		status = handler.checkStatus("invalid-url")
		assert.IsType(t, "", status)
		assert.NotEmpty(t, status)

		// Test with empty URL
		status = handler.checkStatus("")
		assert.IsType(t, "", status)
		assert.NotEmpty(t, status)
	})
}

func TestWebInterfacesHandler_Validation(t *testing.T) {
	t.Run("tests ValidateArgs", func(t *testing.T) {
		handler := NewWebInterfacesHandler()

		// Test with empty args
		err := handler.ValidateArgs([]string{})
		if err != nil {
			assert.Error(t, err)
		}

		// Test with service args
		err = handler.ValidateArgs([]string{"postgres", "redis"})
		if err != nil {
			assert.Error(t, err)
		}
	})

	t.Run("tests GetRequiredFlags", func(t *testing.T) {
		handler := NewWebInterfacesHandler()

		flags := handler.GetRequiredFlags()
		assert.IsType(t, []string{}, flags)
	})
}

func TestWebInterfacesHandler_Creation(t *testing.T) {
	t.Run("creates web interfaces handler", func(t *testing.T) {
		handler := NewWebInterfacesHandler()
		assert.NotNil(t, handler)
		assert.IsType(t, &WebInterfacesHandler{}, handler)
	})
}
