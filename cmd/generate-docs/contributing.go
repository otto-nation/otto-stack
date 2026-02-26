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

	out, err := formatDocument(pageFM(pageContributing), string(data))
	if err != nil {
		return err
	}
	return writeOutput(pageOutput(pageContributing), out)
}
