package project

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateHandler_ValidateArgs(t *testing.T) {
	handler := &ValidateHandler{}
	assert.NoError(t, handler.ValidateArgs([]string{}))
}

func TestValidateHandler_GetRequiredFlags(t *testing.T) {
	handler := &ValidateHandler{}
	flags := handler.GetRequiredFlags()
	assert.Empty(t, flags)
}

func TestDoctorHandler_Creation(t *testing.T) {
	handler := &DoctorHandler{}
	assert.NotNil(t, handler)
}
