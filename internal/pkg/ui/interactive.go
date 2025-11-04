package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/otto-nation/otto-stack/internal/pkg/logger"
)

// Confirm prompts the user for yes/no confirmation
func (o *Output) Confirm(message string, defaultYes bool) bool {
	if o.Quiet {
		return defaultYes
	}

	prompt := message
	if defaultYes {
		prompt += " [Y/n]: "
	} else {
		prompt += " [y/N]: "
	}

	logger.Debug("Interactive confirm prompt", "message", message, "default", defaultYes)
	fmt.Print(prompt)

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return defaultYes
	}

	response = strings.TrimSpace(strings.ToLower(response))

	if response == "" {
		return defaultYes
	}

	result := response == "y" || response == "yes"
	logger.Debug("Interactive confirm result", "response", response, "result", result)
	return result
}

// ConfirmDestructive prompts for confirmation of destructive operations
func (o *Output) ConfirmDestructive(operation string) bool {
	logger.Warn("Destructive operation confirmation", "operation", operation)
	o.Warning("This will %s", operation)
	return o.Confirm("Are you sure you want to continue?", false)
}

// SelectFromList prompts user to select from a list of options
func (o *Output) SelectFromList(message string, options []string) (int, error) {
	if o.Quiet {
		return 0, fmt.Errorf("cannot prompt in quiet mode")
	}

	logger.Debug("Interactive select prompt", "message", message, "options", options)
	fmt.Println(message)
	for i, option := range options {
		fmt.Printf("  %d) %s\n", i+1, option)
	}

	fmt.Print("Select an option (1-" + fmt.Sprintf("%d", len(options)) + "): ")

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}

	response = strings.TrimSpace(response)

	var selection int
	if _, err := fmt.Sscanf(response, "%d", &selection); err != nil {
		return 0, fmt.Errorf("invalid selection: %s", response)
	}

	if selection < 1 || selection > len(options) {
		return 0, fmt.Errorf("selection out of range: %d", selection)
	}

	result := selection - 1
	logger.Debug("Interactive select result", "selection", selection, "result", result)
	return result, nil
}
