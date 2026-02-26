package main

import "fmt"

// docs is the loaded configuration from docs-site/docs.yaml.
// It is populated by loadDocsConfig before any generator runs.
var docs docsConfig

type docsConfig struct {
	Pages      map[string]pageConfig     `yaml:"pages"`
	Categories map[string]categoryConfig `yaml:"categories"`
	Examples   exampleConfig             `yaml:"examples"`
	Labels     labelsConfig              `yaml:"labels"`
}

// pageConfig holds all content for one generated page: frontmatter, body heading,
// page-level intro text, and page-specific section content.
type pageConfig struct {
	Output        string         `yaml:"output"`
	Title         string         `yaml:"title"`
	Heading       string         `yaml:"heading"`
	Description   string         `yaml:"description"`
	Lead          string         `yaml:"lead"`
	Weight        int            `yaml:"weight"`
	Intro         string         `yaml:"intro"`
	FileStructure string         `yaml:"file_structure"`
	NextSteps     []nextStepLink `yaml:"next_steps"`
	ServiceCount  string         `yaml:"service_count"`
	Sections      pageSections   `yaml:"sections"`
}

// nextStepLink is a labelled link used in "Next Steps" sections.
type nextStepLink struct {
	Label       string `yaml:"label"`
	URL         string `yaml:"url"`
	Description string `yaml:"description"`
}

// pageSections is a union of all page-specific section content.
// Each generator reads only the fields relevant to its page; unused fields
// are zero-valued and ignored.
type pageSections struct {
	// cli-reference sections
	CommandCategories string `yaml:"command_categories"`
	Commands          string `yaml:"commands"`
	GlobalFlags       string `yaml:"global_flags"`
	GlobalFlagsDesc   string `yaml:"global_flags_description"`

	// services sections
	ConfigOptions   string `yaml:"configuration_options"`
	ExampleConfig   string `yaml:"example_configuration"`
	UseCases        string `yaml:"use_cases"`
	ExamplesHeading string `yaml:"examples"`

	// configuration sections (each is a dedicated struct due to nested content)
	FileStructure    configFileStructureSection   `yaml:"file_structure"`
	MainConfig       configMainConfigSection      `yaml:"main_config"`
	Sharing          configSharingSection         `yaml:"sharing"`
	ServiceConfig    configServiceConfigSection   `yaml:"service_config"`
	ServiceMetadata  configServiceMetadataSection `yaml:"service_metadata"`
	CompleteExample  configCompleteExampleSection `yaml:"complete_example"`
	NextStepsSection string                       `yaml:"next_steps_section"`
}

type configFileStructureSection struct {
	Heading string `yaml:"heading"`
	Intro   string `yaml:"intro"`
}

type configMainConfigSection struct {
	Heading   string `yaml:"heading"`
	FileLabel string `yaml:"file_label"`
}

type configSharingSection struct {
	Heading      string   `yaml:"heading"`
	Intro        string   `yaml:"intro"`
	Behaviors    []string `yaml:"behaviors"`
	ExampleLabel string   `yaml:"example_label"`
	// Examples is a block scalar containing the full YAML code block content.
	Examples     string `yaml:"examples"`
	RegistryNote string `yaml:"registry_note"`
}

type configServiceConfigSection struct {
	Heading            string `yaml:"heading"`
	Intro              string `yaml:"intro"`
	EnvGeneratedLabel  string `yaml:"env_generated_label"`
	CustomizingHeading string `yaml:"customizing_heading"`
	CustomizingIntro   string `yaml:"customizing_intro"`
	CustomizingNote    string `yaml:"customizing_note"`
}

type configServiceMetadataSection struct {
	Heading string `yaml:"heading"`
	Intro   string `yaml:"intro"`
	// ExampleLabel is the bold file path shown above the code fence.
	ExampleLabel string `yaml:"example_label"`
	// ExampleContent is a block scalar containing the YAML code block content.
	ExampleContent string `yaml:"example_content"`
	Note           string `yaml:"note"`
}

type configCompleteExampleSection struct {
	Heading     string `yaml:"heading"`
	ConfigLabel string `yaml:"config_label"`
	EnvLabel    string `yaml:"env_label"`
}

// labelsConfig holds rendering labels shared across generators.
type labelsConfig struct {
	// Command detail labels used in CLI reference rendering.
	Usage           string `yaml:"usage"`
	Aliases         string `yaml:"aliases"`
	Examples        string `yaml:"examples"`
	Flags           string `yaml:"flags"`
	Tips            string `yaml:"tips"`
	RelatedCommands string `yaml:"related_commands"`
	CommandsList    string `yaml:"commands_list"`

	// Service schema field labels used in services rendering.
	Items                  string `yaml:"items"`
	Properties             string `yaml:"properties"`
	FieldType              string `yaml:"field_type"`
	FieldDefault           string `yaml:"field_default"`
	FieldRequiredYes       string `yaml:"field_required_yes"`
	FieldRequiredIndicator string `yaml:"field_required_indicator"`
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
