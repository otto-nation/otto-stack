package docker

import (
	"fmt"

	"github.com/docker/docker/api/types/filters"
	"github.com/otto-nation/otto-stack/internal/pkg/constants"
)

// NewProjectFilter creates a filter for resources belonging to a specific project
func NewProjectFilter(projectName string) filters.Args {
	f := filters.NewArgs()
	f.Add("label", fmt.Sprintf("%s=%s", constants.ComposeProjectLabel, projectName))
	return f
}

// NewServiceFilter creates a filter for containers belonging to a specific service
func NewServiceFilter(projectName, serviceName string) filters.Args {
	f := NewProjectFilter(projectName)
	f.Add("label", fmt.Sprintf("%s=%s", constants.ComposeServiceLabel, serviceName))
	return f
}
