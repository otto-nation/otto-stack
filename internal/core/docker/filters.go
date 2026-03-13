package docker

import (
	"fmt"

	"github.com/docker/docker/api/types/filters"
)

// NewProjectFilter creates a filter for resources belonging to a specific project.
// When projectName is empty, no label filter is applied and all resources are returned.
func NewProjectFilter(projectName string) filters.Args {
	f := filters.NewArgs()
	if projectName != "" {
		f.Add("label", fmt.Sprintf("%s=%s", ComposeProjectLabel, projectName))
	}
	return f
}
