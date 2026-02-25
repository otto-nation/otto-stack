package cli

import (
	"slices"
	"sort"

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

// completeServiceName returns catalog service names for tab completion.
// Hidden services (internal dependencies) are excluded since users never
// address them directly.
func completeServiceName(_ *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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

	if toComplete == "" {
		return names, cobra.ShellCompDirectiveNoFileComp
	}

	// Filter to names that start with the typed prefix.
	var filtered []string
	for _, name := range names {
		if len(name) >= len(toComplete) && name[:len(toComplete)] == toComplete {
			filtered = append(filtered, name)
		}
	}
	return filtered, cobra.ShellCompDirectiveNoFileComp
}
