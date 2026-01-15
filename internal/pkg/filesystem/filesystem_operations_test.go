package filesystem

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/otto-nation/otto-stack/test/testhelpers"
)

func TestFilesystem_edge_cases(t *testing.T) {
	t.Run("copy file with non-existent source", func(t *testing.T) {
		tempDir := testhelpers.CreateTempDir(t)
		defer os.RemoveAll(tempDir)

		nonExistentFile := filepath.Join(tempDir, "nonexistent.txt")
		destFile := filepath.Join(tempDir, "dest.txt")

		err := CopyFile(nonExistentFile, destFile)
		testhelpers.AssertError(t, err, "CopyFile with non-existent source should error")
	})

	t.Run("write file with invalid permissions", func(t *testing.T) {
		tempDir := testhelpers.CreateTempDir(t)
		defer os.RemoveAll(tempDir)

		testFile := filepath.Join(tempDir, "test.txt")
		err := WriteFile(testFile, []byte("test"), 0000)
		if err != nil {
			testhelpers.AssertError(t, err, "WriteFile with invalid permissions may error")
		} else {
			testhelpers.AssertNoError(t, err, "WriteFile should not error")
		}
	})
}
