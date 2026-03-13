package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	commandsYAMLPath = "internal/config/commands.yaml"
	schemaYAMLPath   = "internal/config/schema.yaml"
	servicesDirPath  = "internal/config/services"
	docsConfigPath   = "docs-site/docs.yaml"
	contributingPath = "CONTRIBUTING.md"
	readmePath       = "README.md"
	outputDirPath    = "docs-site/content"

	permDir  = 0o755
	permFile = 0o644
)

// loadYAML reads path and unmarshals its contents into out.
func loadYAML(path string, out any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	if err := yaml.Unmarshal(data, out); err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}
	return nil
}

func writeOutput(filename, content string) error {
	outPath := filepath.Join(outputDirPath, filename)
	if err := os.MkdirAll(filepath.Dir(outPath), permDir); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}
	if err := os.WriteFile(outPath, []byte(content), permFile); err != nil {
		return fmt.Errorf("write %s: %w", outPath, err)
	}
	fmt.Printf("generated %s\n", outPath)
	return nil
}

func main() {
	generatorFlag := flag.String("generator", "", "Run a specific generator by name")
	flag.Parse()

	if err := loadDocsConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "load docs.yaml: %v\n", err)
		os.Exit(1)
	}

	type generatorFn struct {
		name string
		run  func() error
	}

	allGenerators := []generatorFn{
		{"cli-reference", generateCLIReference},
		{"services-guide", generateServicesGuide},
		{"configuration-guide", generateConfigurationGuide},
		{"homepage", generateHomepage},
		{"contributing-guide", generateContributingGuide},
	}

	var toRun []generatorFn
	if *generatorFlag != "" {
		for _, g := range allGenerators {
			if g.name == *generatorFlag {
				toRun = append(toRun, g)
				break
			}
		}
		if len(toRun) == 0 {
			fmt.Fprintf(os.Stderr, "unknown generator: %s\n", *generatorFlag)
			os.Exit(1)
		}
	} else {
		toRun = allGenerators
	}

	failed := false
	for _, g := range toRun {
		if err := g.run(); err != nil {
			fmt.Fprintf(os.Stderr, "generator %s failed: %v\n", g.name, err)
			failed = true
		}
	}
	if failed {
		os.Exit(1)
	}
}
