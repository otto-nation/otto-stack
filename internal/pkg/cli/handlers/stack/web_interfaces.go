package stack

import (
	"context"
	"fmt"
	"net/http"
	"time"

	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/ci"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/spf13/cobra"
)

const (
	httpTimeout     = 5 * time.Second
	httpOKThreshold = 400
)

// WebInterfacesHandler handles the web interfaces command
type WebInterfacesHandler struct{}

// NewWebInterfacesHandler creates a new web interfaces handler
func NewWebInterfacesHandler() *WebInterfacesHandler {
	return &WebInterfacesHandler{}
}

// Handle executes the web interfaces command
func (h *WebInterfacesHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
	flags, err := core.ParseWebInterfacesFlags(cmd)
	if err != nil {
		return err
	}

	ciFlags := ci.GetFlags(cmd)
	if !ciFlags.Quiet {
		base.Output.Header("%s %s", core.IconCategory_web, core.TitleCase(core.CommandWebInterfaces))
	}

	setup, cleanup, err := SetupCoreCommand(ctx, base)
	if err != nil {
		ci.HandleError(ciFlags, err)
		return nil
	}
	defer cleanup()

	serviceNames := h.getServiceNames(args, setup)
	interfaces, err := h.collectInterfaces(setup, serviceNames, flags.All)
	if err != nil {
		ci.HandleError(ciFlags, err)
		return nil
	}

	h.outputResults(interfaces, ciFlags, base)
	return nil
}

// getServiceNames determines which services to check
func (h *WebInterfacesHandler) getServiceNames(args []string, setup *CoreSetup) []string {
	if len(args) > 0 {
		return args
	}
	return setup.Config.Stack.Enabled
}

// collectInterfaces gathers web interfaces from services
func (h *WebInterfacesHandler) collectInterfaces(setup *CoreSetup, serviceNames []string, showAll bool) ([]WebInterface, error) {
	manager, err := GetServicesManager()
	if err != nil {
		return nil, pkgerrors.NewServiceError(ComponentServiceManager, ActionCreateManager, err)
	}

	resolvedServices, err := manager.ResolveServices(serviceNames)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve services: %w", err)
	}

	runningServices, err := h.getRunningServices(setup, resolvedServices, showAll)
	if err != nil {
		return nil, err
	}

	return h.extractWebInterfaces(manager, resolvedServices, runningServices, showAll), nil
}

// getRunningServices gets the status of services if not showing all
func (h *WebInterfacesHandler) getRunningServices(setup *CoreSetup, services []string, showAll bool) (map[string]bool, error) {
	if showAll {
		return nil, nil
	}

	statuses, err := setup.DockerClient.GetDockerServiceStatus(context.Background(), setup.Config.Project.Name, services)
	if err != nil {
		return nil, fmt.Errorf("failed to get service status: %w", err)
	}

	running := make(map[string]bool)
	for _, status := range statuses {
		running[status.Name] = status.State.IsRunning()
	}
	return running, nil
}

// extractWebInterfaces extracts web interfaces from service definitions
func (h *WebInterfacesHandler) extractWebInterfaces(manager *services.Manager, serviceNames []string, runningServices map[string]bool, showAll bool) []WebInterface {
	var interfaces []WebInterface

	for _, serviceName := range serviceNames {
		if !h.shouldIncludeService(serviceName, runningServices, showAll) {
			continue
		}

		service, err := manager.GetService(serviceName)
		if err != nil {
			continue
		}

		interfaces = append(interfaces, h.createWebInterfaces(serviceName, service.Documentation.WebInterfaces)...)
	}

	return interfaces
}

// createWebInterfaces creates WebInterface structs for a service
func (h *WebInterfacesHandler) createWebInterfaces(serviceName string, webInterfaces []services.WebInterface) []WebInterface {
	interfaces := make([]WebInterface, len(webInterfaces))
	for i, webIface := range webInterfaces {
		interfaces[i] = WebInterface{
			Service:     serviceName,
			Name:        webIface.Name,
			URL:         webIface.URL,
			Description: webIface.Description,
		}
	}
	return interfaces
}

// shouldIncludeService determines if a service should be included in results
func (h *WebInterfacesHandler) shouldIncludeService(serviceName string, runningServices map[string]bool, showAll bool) bool {
	if showAll || runningServices == nil {
		return true
	}
	return runningServices[serviceName]
}

// outputResults outputs the interfaces in the requested format
func (h *WebInterfacesHandler) outputResults(interfaces []WebInterface, ciFlags ci.Flags, base *base.BaseCommand) {
	if ciFlags.JSON {
		ci.OutputResult(ciFlags, map[string]any{
			"interfaces": interfaces,
			"count":      len(interfaces),
		}, core.ExitSuccess)
		return
	}

	if len(interfaces) == 0 {
		base.Output.Info(core.MsgWeb_interfaces_no_interfaces_found)
		return
	}

	h.printTable(interfaces)
}

// printTable prints interfaces in table format
func (h *WebInterfacesHandler) printTable(interfaces []WebInterface) {
	fmt.Printf("%-20s %-28s %-32s %s\n", "SERVICE", "INTERFACE", "URL", "STATUS")
	fmt.Println("--------------------------------------------------------------------------------")

	for _, iface := range interfaces {
		status := h.checkStatus(iface.URL)
		fmt.Printf("%-20s %-28s %-32s %s\n", iface.Service, iface.Name, iface.URL, status)
	}
}

// checkStatus checks if a web interface is accessible
func (h *WebInterfacesHandler) checkStatus(url string) string {
	client := &http.Client{Timeout: httpTimeout}
	resp, err := client.Get(url)
	if err != nil {
		return core.IconHealth_unhealthy + " " + core.MsgWeb_interfaces_not_available
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < httpOKThreshold {
		return core.IconHealth_healthy + " " + core.MsgWeb_interfaces_available
	}
	return core.IconHealth_unhealthy + " " + core.MsgWeb_interfaces_not_available
}

// WebInterface represents a service web interface
type WebInterface struct {
	Service     string `json:"service"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Available   bool   `json:"available"`
}

// ValidateArgs validates the command arguments
func (h *WebInterfacesHandler) ValidateArgs(args []string) error {
	return nil
}

// GetRequiredFlags returns required flags for this command
func (h *WebInterfacesHandler) GetRequiredFlags() []string {
	return []string{}
}
