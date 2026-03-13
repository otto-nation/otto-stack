package cli

import (
	"slices"
	"sort"

	clicontext "github.com/otto-nation/otto-stack/internal/pkg/cli/context"
	"github.com/otto-nation/otto-stack/internal/pkg/registry"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/spf13/cobra"
)

// serviceNameCommands lists the commands that accept service names as positional args.
var serviceNameCommands = []string{"up", "down", "restart", "status", "logs", "deps", "conflicts"}

// RegisterCompletions wires service-name tab completion onto commands that accept
// service names as positional arguments.
func RegisterCompletions(rootCmd *cobra.Command) {
	for _, sub := range rootCmd.Commands() {
		if slices.Contains(serviceNameCommands, sub.Name()) {
			sub.ValidArgsFunction = completeServiceName
		}
	}
}

// completeServiceName returns context-aware service name suggestions:
//   - SharedMode + up: shareable catalog services (any service that can be shared)
//   - SharedMode + down/restart/status/logs: registered shared container names
//   - ProjectMode (or detection failure): full catalog, hidden services excluded
func completeServiceName(cmd *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if sharedMode, ok := detectSharedMode(); ok {
		if cmd.Name() == "up" {
			return completeShareableCatalogServices(toComplete)
		}
		return completeRegisteredServices(sharedMode.Shared.Root, toComplete)
	}
	return completeCatalogServices(toComplete)
}

// detectSharedMode returns the SharedMode execution context if the current working
// directory is outside any otto-stack project. Returns false on any error or when
// inside a project directory.
func detectSharedMode() (*clicontext.SharedMode, bool) {
	detector, err := clicontext.NewDetector()
	if err != nil {
		return nil, false
	}
	execCtx, err := detector.DetectContext()
	if err != nil {
		return nil, false
	}
	sharedMode, ok := execCtx.(*clicontext.SharedMode)
	return sharedMode, ok
}

// completeCatalogServices returns all non-hidden catalog service names.
func completeCatalogServices(toComplete string) ([]string, cobra.ShellCompDirective) {
	manager, err := services.New()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	all := manager.GetAllServices()
	names := make([]string, 0, len(all))
	for name, cfg := range all {
		if !cfg.Hidden {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	return filterByPrefix(names, toComplete), cobra.ShellCompDirectiveNoFileComp
}

// completeShareableCatalogServices returns catalog services that support sharing.
func completeShareableCatalogServices(toComplete string) ([]string, cobra.ShellCompDirective) {
	manager, err := services.New()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	all := manager.GetAllServices()
	names := make([]string, 0, len(all))
	for name, cfg := range all {
		if !cfg.Hidden && cfg.Shareable {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	return filterByPrefix(names, toComplete), cobra.ShellCompDirectiveNoFileComp
}

// completeRegisteredServices returns names from the shared container registry.
// Falls back to shareable catalog services if the registry is empty or unreadable.
func completeRegisteredServices(sharedRoot string, toComplete string) ([]string, cobra.ShellCompDirective) {
	reg := registry.NewManager(sharedRoot)
	containers, err := reg.List()
	if err != nil || len(containers) == 0 {
		// Silent fallback: registry missing or empty — suggest what could be registered.
		return completeShareableCatalogServices(toComplete)
	}

	names := make([]string, 0, len(containers))
	for name := range containers {
		names = append(names, name)
	}
	sort.Strings(names)
	return filterByPrefix(names, toComplete), cobra.ShellCompDirectiveNoFileComp
}

// filterByPrefix returns names that start with the given prefix.
// If the prefix is empty all names are returned.
func filterByPrefix(names []string, prefix string) []string {
	if prefix == "" {
		return names
	}
	var filtered []string
	for _, name := range names {
		if len(name) >= len(prefix) && name[:len(prefix)] == prefix {
			filtered = append(filtered, name)
		}
	}
	return filtered
}
