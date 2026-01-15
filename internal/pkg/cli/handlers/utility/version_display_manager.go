package utility

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/version"
	"gopkg.in/yaml.v3"
)

// BuildInfo contains build information
type BuildInfo struct {
	Version   string `json:"version" yaml:"version"`
	GoVersion string `json:"go_version" yaml:"go_version"`
	Platform  string `json:"platform" yaml:"platform"`
	BuildDate string `json:"build_date,omitempty" yaml:"build_date,omitempty"`
	GitCommit string `json:"git_commit,omitempty" yaml:"git_commit,omitempty"`
}

// VersionDisplayManager handles version display logic
type VersionDisplayManager struct{}

// NewVersionDisplayManager creates a new version display manager
func NewVersionDisplayManager() *VersionDisplayManager {
	return &VersionDisplayManager{}
}

// DisplayBasic displays basic version information
func (vdm *VersionDisplayManager) DisplayBasic(ver, format string) error {
	switch format {
	case "json":
		return vdm.displayJSON(ver)
	case "yaml":
		return vdm.displayYAML(ver)
	default:
		fmt.Printf("%s version %s\n", core.AppName, ver)
		return nil
	}
}

// DisplayFull displays full version information with build details
func (vdm *VersionDisplayManager) DisplayFull(ver, format string) error {
	buildInfo := vdm.getBuildInfo(ver)

	switch format {
	case "json":
		return vdm.displayFullJSON(buildInfo)
	case "yaml":
		return vdm.displayFullYAML(buildInfo)
	default:
		return vdm.displayFullText(buildInfo)
	}
}

// GetCurrentVersion returns the current version
func (vdm *VersionDisplayManager) GetCurrentVersion() string {
	return version.GetAppVersion()
}

func (vdm *VersionDisplayManager) displayJSON(ver string) error {
	output := map[string]string{"version": ver}
	data, _ := json.MarshalIndent(output, "", "  ") // Simple map marshal cannot fail
	fmt.Println(string(data))
	return nil
}

func (vdm *VersionDisplayManager) displayYAML(ver string) error {
	output := map[string]string{"version": ver}
	data, _ := yaml.Marshal(output) // Simple map marshal cannot fail
	fmt.Print(string(data))
	return nil
}

func (vdm *VersionDisplayManager) displayFullJSON(buildInfo BuildInfo) error {
	data, _ := json.MarshalIndent(buildInfo, "", "  ") // Simple struct marshal cannot fail
	fmt.Println(string(data))
	return nil
}

func (vdm *VersionDisplayManager) displayFullYAML(buildInfo BuildInfo) error {
	data, _ := yaml.Marshal(buildInfo) // Simple struct marshal cannot fail
	fmt.Print(string(data))
	return nil
}

func (vdm *VersionDisplayManager) displayFullText(buildInfo BuildInfo) error {
	fmt.Printf("%s version %s\n", core.AppName, buildInfo.Version)
	fmt.Printf("Go version: %s\n", buildInfo.GoVersion)
	fmt.Printf("Platform: %s\n", buildInfo.Platform)
	if buildInfo.BuildDate != "" {
		fmt.Printf("Build date: %s\n", buildInfo.BuildDate)
	}
	if buildInfo.GitCommit != "" {
		fmt.Printf("Git commit: %s\n", buildInfo.GitCommit)
	}
	return nil
}

func (vdm *VersionDisplayManager) getBuildInfo(ver string) BuildInfo {
	return BuildInfo{
		Version:   ver,
		GoVersion: runtime.Version(),
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		BuildDate: "",
		GitCommit: "",
	}
}
