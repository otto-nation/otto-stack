package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBaseHandler_ValidateArgs(t *testing.T) {
	handler := &BaseHandler{}
	assert.NoError(t, handler.ValidateArgs([]string{}))
	assert.NoError(t, handler.ValidateArgs([]string{"arg1", "arg2"}))
}

func TestBaseHandler_GetRequiredFlags(t *testing.T) {
	handler := &BaseHandler{}
	flags := handler.GetRequiredFlags()
	assert.Empty(t, flags)
}
