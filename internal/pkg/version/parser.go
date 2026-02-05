package version

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"github.com/otto-nation/otto-stack/internal/pkg/messages"
)

var semverRegex = regexp.MustCompile(`^v?(\d+)\.(\d+)\.(\d+)(?:-([0-9A-Za-z\-\.]+))?(?:\+([0-9A-Za-z\-\.]+))?$`)

// Version represents a semantic version
type Version struct {
	Major      int    `json:"major"`
	Minor      int    `json:"minor"`
	Patch      int    `json:"patch"`
	PreRelease string `json:"prerelease,omitempty"`
	Build      string `json:"build,omitempty"`
	Original   string `json:"original"`
}

// String returns the string representation of the version
func (v Version) String() string {
	if v.Original != "" {
		return v.Original
	}
	version := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	if v.PreRelease != "" {
		version += "-" + v.PreRelease
	}
	if v.Build != "" {
		version += "+" + v.Build
	}
	return version
}

// Compare compares two versions, returns VersionOlder/VersionEqual/VersionNewer
func (v Version) Compare(other Version) int {
	if v.Major != other.Major {
		if v.Major < other.Major {
			return VersionOlder
		}
		return VersionNewer
	}

	if v.Minor != other.Minor {
		if v.Minor < other.Minor {
			return VersionOlder
		}
		return VersionNewer
	}

	if v.Patch != other.Patch {
		if v.Patch < other.Patch {
			return VersionOlder
		}
		return VersionNewer
	}

	// Handle pre-release versions
	if v.PreRelease == "" && other.PreRelease != "" {
		return VersionNewer // release > pre-release
	}
	if v.PreRelease != "" && other.PreRelease == "" {
		return VersionOlder // pre-release < release
	}
	if v.PreRelease != other.PreRelease {
		if v.PreRelease < other.PreRelease {
			return VersionOlder
		}
		return VersionNewer
	}

	return VersionEqual
}

// ParseVersion parses a version string into a Version struct
func ParseVersion(versionStr string) (*Version, error) {
	if versionStr == "" {
		return nil, pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, "version", messages.ErrorsEmptyVersion, nil)
	}

	versionStr = strings.TrimSpace(versionStr)

	// Handle special cases
	if versionStr == "latest" || versionStr == "*" {
		return &Version{
			Major:    MaxVersionNumber,
			Minor:    MaxVersionNumber,
			Patch:    MaxVersionNumber,
			Original: versionStr,
		}, nil
	}

	matches := semverRegex.FindStringSubmatch(versionStr)
	if matches == nil {
		return nil, pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, "version", messages.VersionInvalidFormat, nil)
	}

	major, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, "version", messages.VersionInvalidMajor, nil)
	}

	minor, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil, pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, "version", messages.VersionInvalidMinor, nil)
	}

	patch, err := strconv.Atoi(matches[3])
	if err != nil {
		return nil, pkgerrors.NewValidationError(pkgerrors.ErrCodeInvalid, "version", messages.VersionInvalidPatch, nil)
	}

	return &Version{
		Major:      major,
		Minor:      minor,
		Patch:      patch,
		PreRelease: matches[4],
		Build:      matches[5],
		Original:   versionStr,
	}, nil
}

// VersionConstraint represents a simple version constraint
type VersionConstraint struct {
	Operator string  `json:"operator"`
	Version  Version `json:"version"`
	Original string  `json:"original"`
}

// Satisfies checks if a version satisfies the constraint
func (c VersionConstraint) Satisfies(version Version) bool {
	cmp := version.Compare(c.Version)

	switch c.Operator {
	case "=", "==", "":
		return cmp == VersionEqual
	case "!=":
		return cmp != VersionEqual
	case ">":
		return cmp == VersionNewer
	case ">=":
		return cmp == VersionNewer || cmp == VersionEqual
	case "<":
		return cmp == VersionOlder
	case "<=":
		return cmp == VersionOlder || cmp == VersionEqual
	case "*":
		return true
	default:
		return false
	}
}

// ParseVersionConstraint parses a version constraint string
