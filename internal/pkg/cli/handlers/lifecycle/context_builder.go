package lifecycle

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/cobra"

	"github.com/otto-nation/otto-stack/internal/core"
	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/shared"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
)

var (
	serviceManagerCache *services.Manager
	serviceManagerMutex sync.RWMutex
)

// BuildStackContext creates CLI context from command flags and arguments
func BuildStackContext(cmd *cobra.Command, args []string) (clicontext.Context, error) {
	// Get project info from config or defaults
	projectName, enabledServices := getProjectInfo()
	projectPath, _ := os.Getwd()
	if projectPath == "" {
		projectPath = "."
	}

	// Use args or fallback to enabled services
	serviceNames := args
	if len(serviceNames) == 0 {
		serviceNames = enabledServices
	}

	// Resolve service configs
	serviceConfigs, err := resolveServiceConfigs(serviceNames)
	if err != nil {
		return clicontext.Context{}, err
	}

	// Parse flags and build context
	forceFlag, _ := cmd.Flags().GetBool("force")
	return clicontext.NewBuilder().
		WithProject(projectName, projectPath).
		WithServices(serviceNames, serviceConfigs).
		WithRuntime(forceFlag, false, false).
		Build(), nil
}

// getProjectInfo loads project name and enabled services from config
func getProjectInfo() (string, []string) {
	configPath := filepath.Join(core.OttoStackDir, core.ConfigFileName)
	if cfg, err := shared.LoadProjectConfig(configPath); err == nil {
		return cfg.Project.Name, cfg.Stack.Enabled
	}
	return "default-project", []string{}
}

// resolveServiceConfigs converts service names to ServiceConfig objects
func resolveServiceConfigs(serviceNames []string) ([]services.ServiceConfig, error) {
	if len(serviceNames) == 0 {
		return []services.ServiceConfig{}, nil
	}

	manager, err := getServiceManager()
	if err != nil {
		return nil, pkgerrors.NewServiceError(ComponentStack, ActionCreateManager, err)
	}

	var configs []services.ServiceConfig
	for _, name := range serviceNames {
		service, err := manager.GetService(name)
		if err != nil {
			return nil, pkgerrors.NewValidationError(pkgerrors.FieldServiceName, "service not found: "+name, err)
		}
		configs = append(configs, *service)
	}
	return configs, nil
}

// getServiceManager returns cached service manager or creates new one
func getServiceManager() (*services.Manager, error) {
	serviceManagerMutex.RLock()
	if serviceManagerCache != nil {
		defer serviceManagerMutex.RUnlock()
		return serviceManagerCache, nil
	}
	serviceManagerMutex.RUnlock()

	serviceManagerMutex.Lock()
	defer serviceManagerMutex.Unlock()

	// Double-check after acquiring write lock
	if serviceManagerCache != nil {
		return serviceManagerCache, nil
	}

	manager, err := services.New()
	if err != nil {
		return nil, err
	}

	serviceManagerCache = manager
	return manager, nil
}
