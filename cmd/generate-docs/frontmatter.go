package main

import (
	"fmt"
	"strings"
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
		Date:        docs.Date,
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

// htmlComment renders lines as an HTML comment block, suitable for embedding
// in generated markdown files to warn editors that content is auto-generated.
func htmlComment(lines ...string) string {
	var sb strings.Builder
	sb.WriteString("<!--\n")
	for _, line := range lines {
		sb.WriteString("  " + line + "\n")
	}
	sb.WriteString("-->\n\n")
	return sb.String()
}

// writePage formats content with frontmatter for the named page and writes it to disk.
func writePage(pageID, content string) error {
	out, err := formatDocument(pageFM(pageID), content)
	if err != nil {
		return err
	}
	return writeOutput(pageOutput(pageID), out)
}

// codeBlock returns content inside a fenced code block with the given language specifier.
// content is trimmed of trailing newlines so the closing fence is always on its own line.
func codeBlock(lang, content string) string {
	const fence = "```"
	return fence + lang + "\n" + strings.TrimRight(content, "\n") + "\n" + fence + "\n\n"
}
