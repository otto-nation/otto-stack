package version

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/otto-nation/otto-stack/internal/pkg/constants"
)

// UpdateNotification represents an update notification
type UpdateNotification struct {
	CurrentVersion  Version `json:"current_version"`
	LatestVersion   Version `json:"latest_version"`
	UpdateType      string  `json:"update_type"` // ChangeType*
	Severity        string  `json:"severity"`    // Severity*
	BreakingChanges bool    `json:"breaking_changes"`
	Message         string  `json:"message"`
}

// NotificationConfig controls update notification behavior
type NotificationConfig struct {
	Enabled         bool          `json:"enabled"`
	Frequency       string        `json:"frequency"`      // NotificationFrequency*
	MinSeverity     string        `json:"min_severity"`   // Severity*
	CheckInterval   time.Duration `json:"check_interval"` // How often to check for updates
	LastCheck       time.Time     `json:"last_check"`
	AutoCheck       bool          `json:"auto_check"`       // Check on command execution
	ShowPrerelease  bool          `json:"show_prerelease"`  // Include prerelease versions
	SuppressedUntil time.Time     `json:"suppressed_until"` // Suppress notifications until this time
}

// Notification frequencies with their check functions
var NotificationFrequencies = map[string]func(*UpdateNotifier) bool{
	"never":  func(n *UpdateNotifier) bool { return false },
	"always": func(n *UpdateNotifier) bool { return true },
	"daily":  func(n *UpdateNotifier) bool { return time.Since(n.config.LastCheck) >= 24*time.Hour },
	"weekly": func(n *UpdateNotifier) bool { return time.Since(n.config.LastCheck) >= 7*24*time.Hour },
}

// Update categorization rules
var updateCategorizationRules = []struct {
	condition       func(current, latest Version) bool
	updateType      string
	severity        string
	breakingChanges bool
	messageTemplate string
}{
	{
		func(c, l Version) bool { return l.Major > c.Major },
		ChangeTypeMajor,
		"medium",
		true,
		"Major update available: %s → %s (may contain breaking changes)",
	},
	{
		func(c, l Version) bool { return l.Minor > c.Minor },
		ChangeTypeMinor,
		"medium",
		false,
		"Minor update available: %s → %s (new features)",
	},
	{
		func(c, l Version) bool { return l.Patch > c.Patch },
		ChangeTypePatch,
		"low",
		false,
		"Patch update available: %s → %s (bug fixes)",
	},
}

// UpdateNotifier handles update notifications
type UpdateNotifier struct {
	manager    VersionManager
	config     NotificationConfig
	configPath string
}

// NewUpdateNotifier creates a new update notifier
func NewUpdateNotifier(manager VersionManager, configPath string) *UpdateNotifier {
	notifier := &UpdateNotifier{
		manager:    manager,
		configPath: configPath,
		config: NotificationConfig{
			Enabled:        DefaultNotifyUpdates,
			Frequency:      "daily",
			MinSeverity:    "medium",
			CheckInterval:  DefaultCheckInterval,
			AutoCheck:      DefaultAutoCheck,
			ShowPrerelease: DefaultShowPrerelease,
		},
	}

	notifier.loadConfig()
	return notifier
}

// CheckForUpdates checks for available updates
func (n *UpdateNotifier) CheckForUpdates() (*UpdateNotification, error) {
	if !n.config.Enabled || !n.shouldCheck() {
		return nil, nil
	}

	currentVersion, latestVersion, err := n.getVersions()
	if err != nil {
		return nil, err
	}

	if latestVersion == nil || currentVersion.Version.Compare(*latestVersion) >= 0 {
		n.updateLastCheck()
		return nil, nil
	}

	notification := n.createNotification(currentVersion.Version, *latestVersion)
	if !n.shouldNotify(notification) {
		return nil, nil
	}

	n.updateLastCheck()
	return notification, nil
}

// getVersions retrieves current and latest available versions
func (n *UpdateNotifier) getVersions() (*InstalledVersion, *Version, error) {
	currentVersion, err := n.manager.GetActiveVersion()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get current version: %w", err)
	}

	availableVersions, err := n.manager.ListAvailableVersions()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list available versions: %w", err)
	}

	latestVersion := n.findLatestVersion(availableVersions)
	return currentVersion, latestVersion, nil
}

// findLatestVersion finds the latest version from available versions
func (n *UpdateNotifier) findLatestVersion(versions []Version) *Version {
	var latest *Version
	for _, version := range versions {
		if !n.config.ShowPrerelease && version.PreRelease != "" {
			continue
		}
		if latest == nil || version.Compare(*latest) > 0 {
			latest = &version
		}
	}
	return latest
}

// createNotification creates and categorizes an update notification
func (n *UpdateNotifier) createNotification(current, latest Version) *UpdateNotification {
	notification := &UpdateNotification{
		CurrentVersion: current,
		LatestVersion:  latest,
	}
	n.categorizeUpdate(notification)
	return notification
}

// NotifyIfNeeded checks for updates and shows notification if needed
func (n *UpdateNotifier) NotifyIfNeeded() error {
	notification, err := n.CheckForUpdates()
	if err != nil {
		return err
	}

	if notification != nil {
		n.showNotification(notification)
	}

	return nil
}

// SuppressNotifications suppresses notifications for a duration
func (n *UpdateNotifier) SuppressNotifications(duration time.Duration) error {
	n.config.SuppressedUntil = time.Now().Add(duration)
	return n.saveConfig()
}

// SetConfig updates the notification configuration
func (n *UpdateNotifier) SetConfig(config NotificationConfig) error {
	n.config = config
	return n.saveConfig()
}

// GetConfig returns the current notification configuration
func (n *UpdateNotifier) GetConfig() NotificationConfig {
	return n.config
}

func (n *UpdateNotifier) shouldCheck() bool {
	if time.Now().Before(n.config.SuppressedUntil) {
		return false
	}

	if checkFunc, exists := NotificationFrequencies[n.config.Frequency]; exists {
		return checkFunc(n)
	}

	// Default fallback
	return time.Since(n.config.LastCheck) >= n.config.CheckInterval
}

func (n *UpdateNotifier) shouldNotify(notification *UpdateNotification) bool {
	minLevel := SeverityLevels[n.config.MinSeverity]
	notificationLevel := SeverityLevels[notification.Severity]

	return notificationLevel >= minLevel
}

func (n *UpdateNotifier) categorizeUpdate(notification *UpdateNotification) {
	current := notification.CurrentVersion
	latest := notification.LatestVersion

	// Check categorization rules in order of precedence
	for _, rule := range updateCategorizationRules {
		if rule.condition(current, latest) {
			notification.UpdateType = rule.updateType
			notification.Severity = rule.severity
			notification.BreakingChanges = rule.breakingChanges
			notification.Message = fmt.Sprintf(rule.messageTemplate, current, latest)
			return
		}
	}

	// Default to prerelease if no other rule matches
	notification.UpdateType = ChangeTypePrerelease
	notification.Severity = "low"
	notification.Message = fmt.Sprintf("Prerelease update available: %s → %s", current, latest)
}

func (n *UpdateNotifier) showNotification(notification *UpdateNotification) {
	fmt.Printf("\n🔔 Update Available!\n")
	fmt.Printf("   Current: %s\n", notification.CurrentVersion)
	fmt.Printf("   Latest:  %s\n", notification.LatestVersion)
	fmt.Printf("   Type:    %s (%s)\n", notification.UpdateType, notification.Severity)

	if notification.BreakingChanges {
		fmt.Printf("   ⚠️  May contain breaking changes\n")
	}

	fmt.Printf("\n   %s\n", notification.Message)
	fmt.Printf("\n   To update: %s version install %s\n", constants.AppName, notification.LatestVersion)
	fmt.Printf("\n   To suppress: %s version suppress 7d\n\n", constants.AppName)
}

func (n *UpdateNotifier) updateLastCheck() {
	n.config.LastCheck = time.Now()
	_ = n.saveConfig() // Ignore errors for non-critical operation
}

func (n *UpdateNotifier) loadConfig() {
	if n.configPath == "" {
		return
	}

	data, err := os.ReadFile(n.configPath)
	if err != nil {
		return // Use defaults
	}

	_ = json.Unmarshal(data, &n.config)
}

func (n *UpdateNotifier) saveConfig() error {
	if n.configPath == "" {
		return nil
	}

	// Ensure directory exists
	dir := filepath.Dir(n.configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(n.config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	return os.WriteFile(n.configPath, data, 0644)
}
