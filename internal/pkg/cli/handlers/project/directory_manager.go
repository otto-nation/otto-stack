package project

import (
	"os"

	"github.com/otto-nation/otto-stack/internal/core"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
)

// DirectoryManager handles directory creation operations
type DirectoryManager struct{}

// NewDirectoryManager creates a new directory manager
func NewDirectoryManager() *DirectoryManager {
	return &DirectoryManager{}
}

// CreateDirectoryStructure creates the otto-stack directory structure
func (dm *DirectoryManager) CreateDirectoryStructure() error {
	directories := []string{
		core.OttoStackDir,
		core.OttoStackDir + "/" + core.ServiceConfigsDir,
	}

	for _, dir := range directories {
		if err := os.MkdirAll(dir, core.PermReadWriteExec); err != nil {
			return pkgerrors.NewConfigError(dir, MsgFailedToCreateDirectory, err)
		}
	}

	return nil
}
