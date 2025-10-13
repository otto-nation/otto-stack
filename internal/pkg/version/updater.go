package version

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"

	"github.com/otto-nation/otto-stack/internal/pkg/constants"
)

// GitHubRelease represents a GitHub release
type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	Draft   bool   `json:"draft"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

// UpdateChecker handles checking for updates
type UpdateChecker struct {
	currentVersion string
}

// NewUpdateChecker creates a new update checker
func NewUpdateChecker(currentVersion string) *UpdateChecker {
	return &UpdateChecker{
		currentVersion: currentVersion,
	}
}

// CheckForUpdates checks if a newer version is available
func (u *UpdateChecker) CheckForUpdates() (*GitHubRelease, bool, error) {
	latest, err := u.getLatestRelease()
	if err != nil {
		return nil, false, fmt.Errorf("failed to get latest release: %w", err)
	}

	if latest.Draft {
		return nil, false, nil
	}

	// Compare versions (simple string comparison for now)
	hasUpdate := latest.TagName != u.currentVersion && latest.TagName > u.currentVersion

	return latest, hasUpdate, nil
}

// SelfUpdater handles self-updating the binary
type SelfUpdater struct {
	currentVersion string
}

// NewSelfUpdater creates a new self-updater
func NewSelfUpdater(currentVersion string) *SelfUpdater {
	return &SelfUpdater{
		currentVersion: currentVersion,
	}
}

// Update downloads and installs the latest version
func (u *SelfUpdater) Update(release *GitHubRelease) error {
	// Get current executable path
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Create backup
	backupPath := execPath + ".backup"
	if err := u.copyFile(execPath, backupPath); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	// Find the right asset for current platform
	platform := fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)
	assetName := fmt.Sprintf("%s-%s", constants.AppName, platform)
	if runtime.GOOS == "windows" {
		assetName += ".exe"
	}

	var downloadURL string
	for _, asset := range release.Assets {
		if asset.Name == assetName {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}

	if downloadURL == "" {
		return fmt.Errorf("no asset found for platform %s", platform)
	}

	// Download new version
	tempFile, err := u.downloadFile(downloadURL)
	if err != nil {
		// Restore backup on failure
		_ = u.copyFile(backupPath, execPath)
		_ = os.Remove(backupPath)
		return fmt.Errorf("failed to download new version: %w", err)
	}

	// Replace current executable
	if err := u.copyFile(tempFile, execPath); err != nil {
		// Restore backup on failure
		_ = u.copyFile(backupPath, execPath)
		_ = os.Remove(backupPath)
		_ = os.Remove(tempFile)
		return fmt.Errorf("failed to replace executable: %w", err)
	}

	// Make executable
	if err := os.Chmod(execPath, 0755); err != nil {
		return fmt.Errorf("failed to set executable permissions: %w", err)
	}

	// Cleanup temp file
	_ = os.Remove(tempFile)

	return nil
}

// Rollback restores the previous version from backup
func (u *SelfUpdater) Rollback() error {
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	backupPath := execPath + ".backup"
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("no backup found at %s", backupPath)
	}

	// Replace current with backup
	if err := u.copyFile(backupPath, execPath); err != nil {
		return fmt.Errorf("failed to restore backup: %w", err)
	}

	// Make executable
	if err := os.Chmod(execPath, 0755); err != nil {
		return fmt.Errorf("failed to set executable permissions: %w", err)
	}

	// Remove backup
	_ = os.Remove(backupPath)

	return nil
}

// getLatestRelease gets the latest release from GitHub API
func (u *UpdateChecker) getLatestRelease() (*GitHubRelease, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest",
		constants.GitHubOrg, constants.GitHubRepo)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			_ = closeErr // Ignore close error
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	return &release, nil
}

// downloadFile downloads a file from URL to a temporary location
func (u *SelfUpdater) downloadFile(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			_ = closeErr // Ignore close error
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	// Create temporary file
	tempFile, err := os.CreateTemp("", constants.AppName+"-update-*")
	if err != nil {
		return "", err
	}
	defer func() {
		if closeErr := tempFile.Close(); closeErr != nil {
			_ = closeErr // Ignore close error
		}
	}()

	// Copy response body to temp file
	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		if removeErr := os.Remove(tempFile.Name()); removeErr != nil {
			_ = removeErr // Ignore remove error
		}
		return "", err
	}

	return tempFile.Name(), nil
}

// copyFile copies a file from src to dst
func (u *SelfUpdater) copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := srcFile.Close(); closeErr != nil {
			_ = closeErr // Ignore close error
		}
	}()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := dstFile.Close(); closeErr != nil {
			_ = closeErr // Ignore close error
		}
	}()

	_, err = io.Copy(dstFile, srcFile)
	return err
}
