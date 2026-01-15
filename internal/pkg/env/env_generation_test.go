package env

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/otto-nation/otto-stack/test/testhelpers"
)

func TestEnv_generate_edge_cases(t *testing.T) {
	t.Run("generate file with empty services", func(t *testing.T) {
		tempDir := testhelpers.CreateTempDir(t)
		defer os.RemoveAll(tempDir)

		envFile := filepath.Join(tempDir, ".env")
		err := GenerateFile("test-project", nil, envFile)
		testhelpers.AssertNoError(t, err, "GenerateFile with empty services should not error")

		// Check if file was created
		if _, err := os.Stat(envFile); os.IsNotExist(err) {
			t.Error("GenerateFile should create .env file")
		}
	})
}
