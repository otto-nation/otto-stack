package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func generateContributingGuide() error {
	data, err := os.ReadFile(contributingPath)
	if err != nil {
		return fmt.Errorf("read CONTRIBUTING.md: %w", err)
	}

	fm := frontmatter{
		Title:       "Contributing",
		Description: "Guide for contributing to otto-stack development",
		Lead:        "Learn how to contribute to otto-stack development",
		Date:        staticDate,
		Lastmod:     today(),
		Draft:       false,
		Weight:      60,
		Toc:         true,
	}
	fmBytes, err := yaml.Marshal(fm)
	if err != nil {
		return err
	}
	out := fmt.Sprintf("---\n%s---\n\n%s", fmBytes, string(data))
	return writeOutput("contributing.md", out)
}
