package version

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/otto-nation/otto-stack/internal/core"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/messages"
)

// GitHubRelease represents a GitHub release
type GitHubRelease struct {
	TagName     string    `json:"tag_name"`
	Name        string    `json:"name"`
	Body        string    `json:"body"`
	Draft       bool      `json:"draft"`
	Prerelease  bool      `json:"prerelease"`
	PublishedAt time.Time `json:"published_at"`
	HTMLURL     string    `json:"html_url"`
}

// UpdateChecker checks for updates from GitHub
type UpdateChecker struct {
	currentVersion string
	client         *http.Client
}

// NewUpdateChecker creates a new update checker
func NewUpdateChecker(currentVersion string) *UpdateChecker {
	return &UpdateChecker{
		currentVersion: currentVersion,
		client: &http.Client{
			Timeout: core.DefaultStartTimeoutSeconds * time.Second,
		},
	}
}

// CheckForUpdates checks GitHub for newer releases
func (u *UpdateChecker) CheckForUpdates() (*GitHubRelease, bool, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest",
		core.GitHubOrg, core.GitHubRepo)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, false, err
	}

	req.Header.Set("User-Agent", GetUserAgent())
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := u.client.Do(req)
	if err != nil {
		return nil, false, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, false, pkgerrors.NewServiceError(pkgerrors.ErrCodeOperationFail, "version", messages.VersionGithubApiError, nil)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, false, err
	}

	// Skip drafts and prereleases
	if release.Draft || release.Prerelease {
		return nil, false, nil
	}

	hasUpdate, err := u.isNewer(release.TagName)
	if err != nil {
		return &release, false, err
	}

	return &release, hasUpdate, nil
}

// isNewer checks if the release version is newer than current
func (u *UpdateChecker) isNewer(releaseVersion string) (bool, error) {
	if IsDevBuild() {
		return false, nil // Dev builds don't need updates
	}

	current, err := ParseVersion(u.currentVersion)
	if err != nil {
		return false, err
	}

	release, err := ParseVersion(releaseVersion)
	if err != nil {
		return false, err
	}

	return release.Compare(*current) == VersionNewer, nil
}

// DetectProjectVersion detects required version from project config
func DetectProjectVersion(projectPath string) (*VersionConstraint, error) {
	// For now, return a wildcard constraint
	// This can be enhanced later to read from project config files
	return &VersionConstraint{
		Operator: "*",
		Version:  Version{},
		Original: "*",
	}, nil
}

// ValidateProjectVersion validates current version meets project requirements
func ValidateProjectVersion(projectPath string) error {
	constraint, err := DetectProjectVersion(projectPath)
	if err != nil {
		return err
	}

	// Dev builds always satisfy constraints
	if IsDevBuild() {
		return nil
	}

	currentVersion, err := ParseVersion(GetAppVersion())
	if err != nil {
		return err
	}

	if !constraint.Satisfies(*currentVersion) {
		return pkgerrors.NewValidationErrorf(pkgerrors.ErrCodeInvalid, "version", messages.VersionConstraintNotSatisfied,
			currentVersion, constraint.Original)
	}

	return nil
}
