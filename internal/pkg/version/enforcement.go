package version

import (
	"fmt"
	"time"

	"github.com/otto-nation/otto-stack/internal/pkg/constants"
)

// Version drift classification rules
var driftClassificationRules = []struct {
	condition            func(required, active Version) bool
	changeType, severity string
}{
	{func(r, a Version) bool { return a.Major != r.Major }, ChangeTypeMajor, "high"},
	{func(r, a Version) bool { return a.Minor != r.Minor }, ChangeTypeMinor, "medium"},
	{func(r, a Version) bool { return a.Patch != r.Patch }, ChangeTypePatch, "low"},
}

// Policy enforcement rules
var policyEnforcementRules = []struct {
	condition   func(e *VersionEnforcer, drift *DriftDetection) bool
	action      string
	messageFunc func(activeVersion, requiredVersion Version, drift *DriftDetection, enforcer *VersionEnforcer) string
}{
	{
		func(e *VersionEnforcer, d *DriftDetection) bool { return e.policy.StrictMode },
		EnforcementActionSwitch,
		func(a, r Version, d *DriftDetection, e *VersionEnforcer) string {
			return fmt.Sprintf("Strict mode: must switch to version %s", r)
		},
	},
	{
		func(e *VersionEnforcer, d *DriftDetection) bool { return !e.policy.AllowDrift },
		EnforcementActionWarn,
		func(a, r Version, d *DriftDetection, e *VersionEnforcer) string {
			return fmt.Sprintf("Version drift detected: using %s, required %s", a, r)
		},
	},
	{
		func(e *VersionEnforcer, d *DriftDetection) bool { return d.DriftDuration > e.policy.MaxDriftDuration },
		EnforcementActionSwitch,
		func(a, r Version, d *DriftDetection, e *VersionEnforcer) string {
			return fmt.Sprintf("Drift duration exceeded: %v > %v", d.DriftDuration, e.policy.MaxDriftDuration)
		},
	},
}

// EnforcementPolicy defines how strict version enforcement should be
type EnforcementPolicy struct {
	StrictMode       bool          `json:"strict_mode"`
	AllowDrift       bool          `json:"allow_drift"`
	MaxDriftDuration time.Duration `json:"max_drift_duration"`
	AutoSync         bool          `json:"auto_sync"`
}

// DriftDetection represents version drift information
type DriftDetection struct {
	ProjectPath     string        `json:"project_path"`
	RequiredVersion Version       `json:"required_version"`
	ActiveVersion   Version       `json:"active_version"`
	DriftDuration   time.Duration `json:"drift_duration"`
	DriftType       string        `json:"drift_type"` // ChangeType*
	Severity        string        `json:"severity"`   // Severity*
}

// EnforcementResult represents the result of version enforcement
type EnforcementResult struct {
	Compliant bool            `json:"compliant"`
	Drift     *DriftDetection `json:"drift,omitempty"`
	Action    string          `json:"action"` // version.EnforcementAction*
	Message   string          `json:"message"`
	ExitCode  int             `json:"exit_code"`
}

// VersionEnforcer handles strict version checking and enforcement
type VersionEnforcer struct {
	manager VersionManager
	policy  EnforcementPolicy
}

// NewVersionEnforcer creates a new version enforcer
func NewVersionEnforcer(manager VersionManager, policy EnforcementPolicy) *VersionEnforcer {
	return &VersionEnforcer{
		manager: manager,
		policy:  policy,
	}
}

// CheckCompliance checks if current version complies with project requirements
func (e *VersionEnforcer) CheckCompliance(projectPath string) (*EnforcementResult, error) {
	// Get project version requirements
	projectConfig, err := e.manager.GetProjectConfig(projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get project config: %w", err)
	}

	// Get active version
	activeVersion, err := e.manager.GetActiveVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to get active version: %w", err)
	}

	// Check if active version satisfies requirements
	if projectConfig.Required.Satisfies(activeVersion.Version) {
		return &EnforcementResult{
			Compliant: true,
			Action:    EnforcementActionNone,
			Message:   "Version compliance satisfied",
			ExitCode:  constants.ExitSuccess,
		}, nil
	}

	// Detect drift and determine enforcement action
	drift := e.detectDrift(projectPath, projectConfig.Required.Version, activeVersion.Version, projectConfig.LastUsed)
	action, message := e.determineAction(drift, activeVersion.Version, projectConfig.Required.Version)

	exitCode := constants.ExitSuccess
	if action == EnforcementActionSwitch {
		exitCode = constants.ExitError
	}

	return &EnforcementResult{
		Compliant: false,
		Drift:     drift,
		Action:    action,
		Message:   message,
		ExitCode:  exitCode,
	}, nil
}

// determineAction determines the enforcement action based on policy
func (e *VersionEnforcer) determineAction(drift *DriftDetection, activeVersion, requiredVersion Version) (string, string) {
	for _, rule := range policyEnforcementRules {
		if rule.condition(e, drift) {
			return rule.action, rule.messageFunc(activeVersion, requiredVersion, drift, e)
		}
	}

	return EnforcementActionNone, "Drift within acceptable limits"
}

// EnforceCompliance enforces version compliance based on policy
func (e *VersionEnforcer) EnforceCompliance(projectPath string) error {
	result, err := e.CheckCompliance(projectPath)
	if err != nil {
		return err
	}

	switch result.Action {
	case EnforcementActionSwitch:
		if e.policy.AutoSync {
			return e.autoSwitchVersion(projectPath)
		}
		return fmt.Errorf("version enforcement required: %s", result.Message)
	case EnforcementActionWarn:
		fmt.Printf("⚠️  Warning: %s\n", result.Message)
		return nil
	default:
		return nil
	}
}

// DetectAllDrift detects version drift across all projects
func (e *VersionEnforcer) DetectAllDrift() ([]DriftDetection, error) {
	projects, err := e.manager.ListProjectConfigs()
	if err != nil {
		return nil, fmt.Errorf("failed to list project configs: %w", err)
	}

	activeVersion, err := e.manager.GetActiveVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to get active version: %w", err)
	}

	var drifts []DriftDetection
	for _, project := range projects {
		if !project.Required.Satisfies(activeVersion.Version) {
			drift := e.detectDrift(project.ProjectPath, project.Required.Version, activeVersion.Version, project.LastUsed)
			drifts = append(drifts, *drift)
		}
	}

	return drifts, nil
}

func (e *VersionEnforcer) detectDrift(projectPath string, required, active Version, lastUsed time.Time) *DriftDetection {
	driftType, severity := e.classifyVersionDrift(required, active)

	return &DriftDetection{
		ProjectPath:     projectPath,
		RequiredVersion: required,
		ActiveVersion:   active,
		DriftDuration:   time.Since(lastUsed),
		DriftType:       driftType,
		Severity:        severity,
	}
}

// classifyVersionDrift determines the type and severity of version drift
func (e *VersionEnforcer) classifyVersionDrift(required, active Version) (string, string) {
	if active.Compare(required) == 0 {
		return ChangeTypeNone, "low"
	}

	for _, rule := range driftClassificationRules {
		if rule.condition(required, active) {
			return rule.changeType, rule.severity
		}
	}

	return ChangeTypePrerelease, "low"
}

func (e *VersionEnforcer) autoSwitchVersion(projectPath string) error {
	projectConfig, err := e.manager.GetProjectConfig(projectPath)
	if err != nil {
		return fmt.Errorf("failed to get project config: %w", err)
	}

	// Try to resolve to an installed version first
	if installedVersion, err := e.manager.ResolveVersion(projectConfig.Required); err == nil {
		return e.manager.SwitchToVersion(installedVersion.Version)
	}

	// Find and install the best matching available version
	bestMatch, err := e.findBestMatchingVersion(projectConfig.Required)
	if err != nil {
		return err
	}

	// Install and switch to the best match
	if err := e.manager.InstallVersion(*bestMatch); err != nil {
		return fmt.Errorf("failed to install version %s: %w", bestMatch, err)
	}

	return e.manager.SwitchToVersion(*bestMatch)
}

// findBestMatchingVersion finds the best available version that satisfies the constraint
func (e *VersionEnforcer) findBestMatchingVersion(constraint VersionConstraint) (*Version, error) {
	availableVersions, err := e.manager.ListAvailableVersions()
	if err != nil {
		return nil, fmt.Errorf("failed to list available versions: %w", err)
	}

	var bestMatch *Version
	for _, version := range availableVersions {
		if constraint.Satisfies(version) && (bestMatch == nil || version.Compare(*bestMatch) > 0) {
			bestMatch = &version
		}
	}

	if bestMatch == nil {
		return nil, fmt.Errorf("no available version satisfies constraint %s", constraint.Original)
	}

	return bestMatch, nil
}
