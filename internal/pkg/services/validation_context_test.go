//go:build unit

package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidationContext(t *testing.T) {
	t.Run("user context rejects hidden", func(t *testing.T) {
		ctx := NewUserValidationContext()
		assert.False(t, ctx.AllowHidden)
		assert.True(t, ctx.IsUserRequested)
		assert.False(t, ctx.IsDependency)
	})

	t.Run("dependency context allows hidden", func(t *testing.T) {
		ctx := NewDependencyValidationContext()
		assert.True(t, ctx.AllowHidden)
		assert.False(t, ctx.IsUserRequested)
		assert.True(t, ctx.IsDependency)
	})

	t.Run("internal context allows hidden", func(t *testing.T) {
		ctx := NewInternalValidationContext()
		assert.True(t, ctx.AllowHidden)
		assert.False(t, ctx.IsUserRequested)
		assert.False(t, ctx.IsDependency)
	})
}
