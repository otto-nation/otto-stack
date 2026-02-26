package main

import (
	"fmt"
	"time"

	"gopkg.in/yaml.v3"
)

type frontmatter struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Lead        string `yaml:"lead"`
	Date        string `yaml:"date"`
	Lastmod     string `yaml:"lastmod"`
	Draft       bool   `yaml:"draft"`
	Weight      int    `yaml:"weight"`
	Toc         bool   `yaml:"toc"`
}

func today() string {
	return time.Now().Format("2006-01-02")
}

func newFrontmatter(title, description, lead string, weight int) frontmatter {
	return frontmatter{
		Title:       title,
		Description: description,
		Lead:        lead,
		Date:        staticDate,
		Lastmod:     today(),
		Draft:       false,
		Weight:      weight,
		Toc:         true,
	}
}

func formatDocument(fm frontmatter, content string) (string, error) {
	fmBytes, err := yaml.Marshal(fm)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("---\n%s---\n\n%s", fmBytes, content), nil
}
