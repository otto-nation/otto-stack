package filesystem

import (
	"os"
	"path/filepath"

	"github.com/otto-nation/otto-stack/internal/core"
)

// EnsureDir creates a directory if it doesn't exist
func EnsureDir(dirname string) error {
	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		return os.MkdirAll(dirname, core.PermReadWriteExec)
	}
	return nil
}

// WriteFile writes content to a file, creating directories if needed
func WriteFile(filename string, content []byte, perm os.FileMode) error {
	if err := EnsureDir(filepath.Dir(filename)); err != nil {
		return err
	}
	return os.WriteFile(filename, content, perm)
}
