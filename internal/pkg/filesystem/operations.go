package filesystem

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/otto-nation/otto-stack/internal/core"
)

// EnsureDir creates a directory if it doesn't exist
func EnsureDir(dirname string) error {
	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		return os.MkdirAll(dirname, core.PermReadWriteExec)
	}
	return nil
}

// CopyFile copies a file from src to dst
func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() { _ = sourceFile.Close() }()

	if err := EnsureDir(filepath.Dir(dst)); err != nil {
		return err
	}

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() { _ = destFile.Close() }()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

// WriteFile writes content to a file, creating directories if needed
func WriteFile(filename string, content []byte, perm os.FileMode) error {
	if err := EnsureDir(filepath.Dir(filename)); err != nil {
		return err
	}
	return os.WriteFile(filename, content, perm)
}

// ExpandPath expands ~ and environment variables in a path
func ExpandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			path = filepath.Join(home, path[2:])
		}
	}
	return os.ExpandEnv(path)
}
