package docker

import (
	"fmt"

	"github.com/docker/docker/api/types/filters"
)

// NewProjectFilter creates a filter for resources belonging to a specific project
func NewProjectFilter(projectName string) filters.Args {
	f := filters.NewArgs()
	f.Add("label", fmt.Sprintf("%s=%s", ComposeProjectLabel, projectName))
	return f
}
