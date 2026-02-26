package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

func generateHomepage() error {
	data, err := os.ReadFile(readmePath)
	if err != nil {
		return fmt.Errorf("read README.md: %w", err)
	}

	content := extractFromFirstHeading(string(data))
	content = fixHomepageLinks(content)

	fm := frontmatter{
		Title:       "otto-stack",
		Description: "A powerful development stack management tool built in Go for streamlined local development automation",
		Lead:        "Streamline your local development with powerful CLI tools and automated service management",
		Date:        staticDate,
		Lastmod:     today(),
		Draft:       false,
		Weight:      50,
		Toc:         true,
	}
	fmBytes, err := yaml.Marshal(fm)
	if err != nil {
		return err
	}
	out := fmt.Sprintf("---\n%s---\n\n%s", fmBytes, content)
	return writeOutput("_index.md", out)
}

// extractFromFirstHeading trims content before the first top-level heading.
func extractFromFirstHeading(content string) string {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "# ") {
			return strings.Join(lines[i:], "\n")
		}
	}
	return content
}

// fixHomepageLinks rewrites README links to work under the Hugo site's baseURL.
func fixHomepageLinks(content string) string {
	// Convert docs-site/content/file.md → file/
	content = regexp.MustCompile(`docs-site/content/([^)]+)\.md`).
		ReplaceAllString(content, "$1/")

	// Convert remaining .md links to Hugo directory format
	content = regexp.MustCompile(`\]\(([^)]+)\.md\)`).
		ReplaceAllString(content, "]($1/)")

	// Convert docs-site root links to Hugo ref shortcode
	reSiteRoot := regexp.MustCompile(`\[([^\]]+)\]\(docs-site/\)`)
	content = reSiteRoot.ReplaceAllStringFunc(content, func(match string) string {
		inner := reSiteRoot.FindStringSubmatch(match)
		if len(inner) > 1 {
			return fmt.Sprintf(`[%s]({{< ref "/" >}})`, inner[1])
		}
		return match
	})

	// Expand bare LICENSE link to full GitHub URL
	content = regexp.MustCompile(`\[([^\]]+)\]\(LICENSE\)`).
		ReplaceAllString(content, "[$1](https://github.com/otto-nation/otto-stack/blob/main/LICENSE)")

	return content
}
