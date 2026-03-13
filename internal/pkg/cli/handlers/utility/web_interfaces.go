package utility

import (
	"context"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/ci"
	"github.com/otto-nation/otto-stack/internal/pkg/cli/handlers/common"
	"github.com/otto-nation/otto-stack/internal/pkg/display"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/messages"
	"github.com/otto-nation/otto-stack/internal/pkg/services"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/otto-nation/otto-stack/internal/pkg/ui"
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
		return pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, pkgerrors.FieldFlags, messages.ValidationFailedParseFlags, err)
	}

	ciFlags := ci.GetFlags(cmd)
	if !ciFlags.Quiet {
		base.Output.Header("%s", core.TitleCase(core.CommandWebInterfaces))
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

	interfaces, err := h.collectInterfaces(ctx, setup, serviceConfigs, flags.All)
	if err != nil {
		return ci.FormatError(ciFlags, err)
	}

	h.outputResults(interfaces, ciFlags, base)
	return nil
}

// collectInterfaces gathers web interfaces from services and checks their availability concurrently.
func (h *WebInterfacesHandler) collectInterfaces(ctx context.Context, setup *common.CoreSetup, serviceConfigs []types.ServiceConfig, showAll bool) ([]WebInterface, error) {
	runningServices, err := h.getRunningServices(ctx, setup, serviceConfigs, showAll)
	if err != nil {
		return nil, err
	}

	interfaces := h.extractWebInterfaces(serviceConfigs, runningServices, showAll)
	h.checkAvailabilityConcurrent(interfaces)
	return interfaces, nil
}

// getRunningServices gets the status of services if not showing all
func (h *WebInterfacesHandler) getRunningServices(ctx context.Context, setup *common.CoreSetup, serviceConfigs []types.ServiceConfig, showAll bool) (map[string]bool, error) {
	if showAll {
		return nil, nil
	}

	serviceNames := services.ExtractServiceNames(serviceConfigs)

	stackService, err := common.NewServiceManager(false)
	if err != nil {
		return nil, pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentStack, messages.ErrorsStackCreateFailed, err)
	}

	statuses, err := stackService.Status(ctx, services.StatusRequest{
		Project:  setup.Config.Project.Name,
		Services: serviceNames,
	})
	if err != nil {
		return nil, pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, pkgerrors.ComponentStack, messages.ErrorsStackGetStatusFailed, err)
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

// createWebInterfaces creates WebInterface structs for a service.
// Available is left false until checkAvailabilityConcurrent fills it in.
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

// checkAvailabilityConcurrent performs HTTP checks for all interfaces in parallel.
func (h *WebInterfacesHandler) checkAvailabilityConcurrent(interfaces []WebInterface) {
	var wg sync.WaitGroup
	client := &http.Client{Timeout: core.DefaultHTTPTimeoutSeconds * time.Second}

	for i := range interfaces {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			interfaces[idx].Available = h.checkURL(client, interfaces[idx].URL)
		}(i)
	}
	wg.Wait()
}

// checkURL returns true if the URL responds with a successful status code.
func (h *WebInterfacesHandler) checkURL(client *http.Client, url string) bool {
	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	defer func() { _ = resp.Body.Close() }()
	return resp.StatusCode < core.HTTPOKStatusThreshold
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
		base.Output.Info(messages.WebInterfacesNoInterfacesFound)
		return
	}

	h.printTable(interfaces, base.Output.Writer())
}

func (h *WebInterfacesHandler) outputJSON(interfaces []WebInterface, ciFlags ci.Flags) {
	output := ci.InterfacesOutput{
		Interfaces: make([]any, len(interfaces)),
		Count:      len(interfaces),
	}
	for i, iface := range interfaces {
		output.Interfaces[i] = iface
	}
	ci.OutputResult(ciFlags, output, core.ExitSuccess)
}

// printTable prints interfaces in table format
func (h *WebInterfacesHandler) printTable(interfaces []WebInterface, writer io.Writer) {
	headers := []string{display.HeaderService, display.HeaderInterface, display.HeaderURL, display.HeaderStatus}
	rows := make([][]string, len(interfaces))

	for i, iface := range interfaces {
		rows[i] = []string{iface.Service, iface.Name, iface.URL, h.formatStatus(iface.Available)}
	}

	display.RenderTable(writer, headers, rows)
}

func (h *WebInterfacesHandler) formatStatus(available bool) string {
	if available {
		return ui.IconOK + " " + messages.WebInterfacesAvailable
	}
	return ui.IconFail + " " + messages.WebInterfacesNotAvailable
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
