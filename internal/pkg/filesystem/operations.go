package filesystem

import (
	"os"
	"path/filepath"

	"github.com/otto-nation/otto-stack/internal/core"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/messages"
)

// EnsureDir creates a directory if it doesn't exist
func EnsureDir(dirname string) error {
	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		if err := os.MkdirAll(dirname, core.PermReadWriteExec); err != nil {
			return pkgerrors.NewSystemError(pkgerrors.ErrCodeOperationFail, messages.ErrorsDirectoryCreateFailed, err)
		}
	}
	return nil
}

// WriteFile writes content to a file, creating directories if needed
func WriteFile(filename string, content []byte, perm os.FileMode) error {
	if err := EnsureDir(filepath.Dir(filename)); err != nil {
		return err
	}
	if err := os.WriteFile(filename, content, perm); err != nil {
		return pkgerrors.NewSystemError(pkgerrors.ErrCodeOperationFail, messages.ErrorsFileWriteFailed, err)
	}
	return nil
}
