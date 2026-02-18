package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServiceUtils_Creation(t *testing.T) {
	utils := &ServiceUtils{}
	assert.NotNil(t, utils)
}
