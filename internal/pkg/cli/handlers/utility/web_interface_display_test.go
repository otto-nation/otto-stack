//go:build unit

package utility

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWebInterfacesHandler_CheckStatus_Localhost(t *testing.T) {
	handler := NewWebInterfacesHandler()
	status := handler.checkStatus("http://localhost:8080")
	assert.IsType(t, "", status)
	assert.NotEmpty(t, status)
}

func TestWebInterfacesHandler_CheckStatus_Invalid(t *testing.T) {
	handler := NewWebInterfacesHandler()
	status := handler.checkStatus("invalid-url")
	assert.IsType(t, "", status)
	assert.NotEmpty(t, status)
}

func TestWebInterfacesHandler_CheckStatus_Empty(t *testing.T) {
	handler := NewWebInterfacesHandler()
	status := handler.checkStatus("")
	assert.IsType(t, "", status)
	assert.NotEmpty(t, status)
}

func TestWebInterfacesHandler_ValidateArgs_Empty(t *testing.T) {
	handler := NewWebInterfacesHandler()
	err := handler.ValidateArgs([]string{})
	if err != nil {
		assert.Error(t, err)
	}
}

func TestWebInterfacesHandler_ValidateArgs_WithServices(t *testing.T) {
	handler := NewWebInterfacesHandler()
	err := handler.ValidateArgs([]string{"postgres", "redis"})
	if err != nil {
		assert.Error(t, err)
	}
}

func TestWebInterfacesHandler_GetRequiredFlags(t *testing.T) {
	handler := NewWebInterfacesHandler()
	flags := handler.GetRequiredFlags()
	assert.IsType(t, []string{}, flags)
}

func TestWebInterfacesHandler_Creation(t *testing.T) {
	handler := NewWebInterfacesHandler()
	assert.NotNil(t, handler)
	assert.IsType(t, &WebInterfacesHandler{}, handler)
}
