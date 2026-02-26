package main

import "fmt"

// docs is the loaded configuration from docs-site/docs.yaml.
// It is populated by loadDocsConfig before any generator runs.
var docs docsConfig

type docsConfig struct {
	Pages      map[string]pageConfig     `yaml:"pages"`
	Categories map[string]categoryConfig `yaml:"categories"`
	Examples   exampleConfig             `yaml:"examples"`
}

// pageConfig holds the frontmatter and output filename for one generated page.
type pageConfig struct {
	Output      string `yaml:"output"`
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Lead        string `yaml:"lead"`
	Weight      int    `yaml:"weight"`
}

// exampleConfig holds the values used in generated code samples.
type exampleConfig struct {
	ProjectName           string   `yaml:"project_name"`
	FullstackProjectName  string   `yaml:"fullstack_project_name"`
	ProjectType           string   `yaml:"project_type"`
	Services              []string `yaml:"services"`
	CompleteServices      []string `yaml:"complete_services"`
	EnvVarDisplayLimit    int      `yaml:"env_var_display_limit"`
	CustomEnvDisplayLimit int      `yaml:"custom_env_display_limit"`
}

var requiredPages = []string{
	"homepage",
	"cli-reference",
	"services",
	"configuration",
	"contributing",
}

func loadDocsConfig() error {
	if err := loadYAML(docsConfigPath, &docs); err != nil {
		return err
	}
	for _, name := range requiredPages {
		if _, ok := docs.Pages[name]; !ok {
			return fmt.Errorf("docs.yaml missing required page: %q", name)
		}
	}
	return nil
}

// pageFM returns the frontmatter for the named page.
func pageFM(name string) frontmatter {
	p := docs.Pages[name]
	return newFrontmatter(p.Title, p.Description, p.Lead, p.Weight)
}

// pageOutput returns the output filename for the named page.
func pageOutput(name string) string {
	return docs.Pages[name].Output
}
