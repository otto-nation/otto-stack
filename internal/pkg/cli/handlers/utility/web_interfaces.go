package utility

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/ci"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/common"
	"github.com/otto-nation/otto-stack/internal/pkg/display"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/spf13/cobra"
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

	setup, cleanup, err := common.SetupCoreCommand(ctx, base)
	if err != nil {
		return ci.FormatError(ciFlags, err)
	}
	defer cleanup()

	// Resolve services - for utility handlers, use all enabled services if no args provided
	var serviceConfigs []types.ServiceConfig
	if len(args) > 0 {
		serviceConfigs, err = services.ResolveUpServices(args, setup.Config)
	} else {
		serviceConfigs, err = services.ResolveUpServices(setup.Config.Stack.Enabled, setup.Config)
	}
	if err != nil {
		return ci.FormatError(ciFlags, err)
	}

	interfaces, err := h.collectInterfaces(setup, serviceConfigs, flags.All)
	if err != nil {
		return ci.FormatError(ciFlags, err)
	}

	h.outputResults(interfaces, ciFlags, base)
	return nil
}

// collectInterfaces gathers web interfaces from services
func (h *WebInterfacesHandler) collectInterfaces(setup *common.CoreSetup, serviceConfigs []types.ServiceConfig, showAll bool) ([]WebInterface, error) {
	runningServices, err := h.getRunningServices(setup, serviceConfigs, showAll)
	if err != nil {
		return nil, err
	}

	return h.extractWebInterfaces(serviceConfigs, runningServices, showAll), nil
}

// getRunningServices gets the status of services if not showing all
func (h *WebInterfacesHandler) getRunningServices(setup *common.CoreSetup, serviceConfigs []types.ServiceConfig, showAll bool) (map[string]bool, error) {
	if showAll {
		return nil, nil
	}

	serviceNames := services.ExtractServiceNames(serviceConfigs)

	stackService, err := common.NewServiceManager(false)
	if err != nil {
		return nil, pkgerrors.NewServiceError(common.ComponentStack, common.MsgFailedCreateStackService, err)
	}

	statuses, err := stackService.DockerClient.GetServiceStatus(context.Background(), setup.Config.Project.Name, serviceNames)
	if err != nil {
		return nil, pkgerrors.NewServiceError("stack", "get service status", err)
	}

	running := make(map[string]bool)
	for _, status := range statuses {
		running[status.Name] = status.State == docker.HealthStatusRunning
	}
	return running, nil
}

// extractWebInterfaces extracts web interfaces from service definitions
func (h *WebInterfacesHandler) extractWebInterfaces(serviceConfigs []types.ServiceConfig, runningServices map[string]bool, showAll bool) []WebInterface {
	var interfaces []WebInterface

	for _, config := range serviceConfigs {
		if !h.shouldIncludeService(config.Name, runningServices, showAll) {
			continue
		}

		// Extract web interface info directly from ServiceConfig
		interfaces = append(interfaces, h.createWebInterfaces(config.Name, config.Documentation.WebInterfaces)...)
	}

	return interfaces
}

// createWebInterfaces creates WebInterface structs for a service
func (h *WebInterfacesHandler) createWebInterfaces(serviceName string, webInterfaces []types.WebInterface) []WebInterface {
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
		h.outputJSON(interfaces, ciFlags)
		return
	}

	if len(interfaces) == 0 {
		base.Output.Info(core.MsgWeb_interfaces_no_interfaces_found)
		return
	}

	h.printTable(interfaces)
}

func (h *WebInterfacesHandler) outputJSON(interfaces []WebInterface, ciFlags ci.Flags) {
	ci.OutputResult(ciFlags, map[string]any{
		"interfaces": interfaces,
		"count":      len(interfaces),
	}, core.ExitSuccess)
}

// printTable prints interfaces in table format
func (h *WebInterfacesHandler) printTable(interfaces []WebInterface) {
	headers := []string{display.HeaderService, display.HeaderInterface, display.HeaderURL, display.HeaderStatus}
	rows := make([][]string, len(interfaces))

	for i, iface := range interfaces {
		status := h.checkStatus(iface.URL)
		rows[i] = []string{iface.Service, iface.Name, iface.URL, status}
	}

	display.RenderTable(os.Stdout, headers, rows)
}

// checkStatus checks if a web interface is accessible
func (h *WebInterfacesHandler) checkStatus(url string) string {
	client := &http.Client{Timeout: core.DefaultHTTPTimeoutSeconds * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return h.formatStatus(false)
	}
	defer func() { _ = resp.Body.Close() }()

	return h.formatStatusFromResponse(resp.StatusCode)
}

func (h *WebInterfacesHandler) formatStatus(available bool) string {
	if available {
		return core.IconHealth_healthy + " " + core.MsgWeb_interfaces_available
	}
	return core.IconHealth_unhealthy + " " + core.MsgWeb_interfaces_not_available
}

func (h *WebInterfacesHandler) formatStatusFromResponse(statusCode int) string {
	return h.formatStatus(statusCode < core.HTTPOKStatusThreshold)
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
