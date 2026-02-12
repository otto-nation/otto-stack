package env

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
)

// GenerateFile generates and writes the env file to disk from service configs
func GenerateFile(projectName string, serviceConfigs []types.ServiceConfig, filePath string) error {
	var content strings.Builder

	fmt.Fprintf(&content, core.EnvGeneratedHeader, time.Now().Format(time.RFC1123))
	fmt.Fprintf(&content, "PROJECT_NAME=%s\n", projectName)
	fmt.Fprintf(&content, "COMPOSE_PROJECT_NAME=%s\n\n", projectName)

	for _, config := range serviceConfigs {
		if len(config.AllEnvironment) > 0 {
			fmt.Fprintf(&content, "# %s\n", strings.ToUpper(config.Name))
			for key, value := range config.AllEnvironment {
				fmt.Fprintf(&content, "%s=%s\n", key, value)
			}
			content.WriteString("\n")
		}
	}

	if err := os.MkdirAll(filepath.Dir(filePath), core.PermReadWriteExec); err != nil {
		return err
	}
	return os.WriteFile(filePath, []byte(content.String()), core.PermReadWrite)
}
