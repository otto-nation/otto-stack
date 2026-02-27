package main

import (
	"fmt"
	"os"
)

func generateContributingGuide() error {
	data, err := os.ReadFile(contributingPath)
	if err != nil {
		return fmt.Errorf("read CONTRIBUTING.md: %w", err)
	}

	return writePage(pageContributing, string(data))
}
