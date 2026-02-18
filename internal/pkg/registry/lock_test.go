package registry

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLockFile(t *testing.T) {
	tempFile, err := os.CreateTemp("", "lock-test-*.txt")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	err = lockFile(tempFile)
	assert.NoError(t, err)

	err = unlockFile(tempFile)
	assert.NoError(t, err)
}
