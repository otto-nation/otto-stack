package base

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBaseHandler_ValidateArgs(t *testing.T) {
	handler := &BaseHandler{}

	tests := []struct {
		name string
		args []string
	}{
		{"accepts no arguments", []string{}},
		{"accepts single argument", []string{"arg1"}},
		{"accepts multiple arguments", []string{"arg1", "arg2", "arg3"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.ValidateArgs(tt.args)
			assert.NoError(t, err)
		})
	}
}

func TestBaseHandler_GetRequiredFlags(t *testing.T) {
	handler := &BaseHandler{}

	t.Run("returns empty slice", func(t *testing.T) {
		flags := handler.GetRequiredFlags()
		assert.Empty(t, flags)
		assert.Equal(t, []string{}, flags)
	})
}
