package registry

import (
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/messages"
)

// OrphanDisplay handles displaying orphan information
type OrphanDisplay struct {
	output base.Output
}

// NewOrphanDisplay creates a new orphan display
func NewOrphanDisplay(output base.Output) *OrphanDisplay {
	return &OrphanDisplay{output: output}
}

// Display shows orphan information grouped by severity
func (d *OrphanDisplay) Display(orphans []OrphanInfo) {
	if len(orphans) == 0 {
		return
	}

	safe, warning, critical := d.groupBySeverity(orphans)

	d.output.Warning(messages.OrphanFound, len(orphans))
	d.displayCritical(critical)
	d.displayWarning(warning)
	d.displaySafe(safe)
	d.output.Info(messages.OrphanRunCleanupHint)
}

func (d *OrphanDisplay) displayCritical(orphans []OrphanInfo) {
	if len(orphans) == 0 {
		return
	}
	d.output.Error(messages.OrphanSeverityCritical, len(orphans))
	for _, o := range orphans {
		d.output.Info("    - %s: %s", o.Service, o.Reason)
	}
}

func (d *OrphanDisplay) displayWarning(orphans []OrphanInfo) {
	if len(orphans) == 0 {
		return
	}
	d.output.Warning(messages.OrphanSeverityWarning, len(orphans))
	for _, o := range orphans {
		d.output.Info("    - %s: %s", o.Service, o.Reason)
		if len(o.ProjectsFound) > 0 {
			d.output.Info("      "+messages.OrphanRemainingProjects, o.ProjectsFound)
		}
	}
}

func (d *OrphanDisplay) displaySafe(orphans []OrphanInfo) {
	if len(orphans) == 0 {
		return
	}
	d.output.Info(messages.OrphanSeveritySafe, len(orphans))
	for _, o := range orphans {
		d.output.Info("    - %s: %s", o.Service, o.Reason)
	}
}

func (d *OrphanDisplay) groupBySeverity(orphans []OrphanInfo) (safe, warning, critical []OrphanInfo) {
	for _, o := range orphans {
		switch o.Severity {
		case OrphanSeveritySafe:
			safe = append(safe, o)
		case OrphanSeverityWarning:
			warning = append(warning, o)
		case OrphanSeverityCritical:
			critical = append(critical, o)
		}
	}
	return
}
