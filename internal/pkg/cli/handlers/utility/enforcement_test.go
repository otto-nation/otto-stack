package utility

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnforcementHandler_HandleEnforce(t *testing.T) {
	handler := &EnforcementHandler{}
	assert.NotNil(t, handler)
}
