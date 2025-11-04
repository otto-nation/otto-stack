package utils

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/otto-nation/otto-stack/internal/pkg/logger"
)

// AskConfirmation asks for user confirmation
func AskConfirmation(message string) bool {
	logger.Debug("Ask confirmation", "message", message)
	fmt.Printf("%s (y/N): ", message)
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		response := strings.ToLower(strings.TrimSpace(scanner.Text()))
		result := response == "y" || response == "yes"
		logger.Debug("Confirmation result", "response", response, "result", result)
		return result
	}
	return false
}

// PromptInput prompts for user input with a message
func PromptInput(message string) (string, error) {
	logger.Debug("Prompt input", "message", message)
	fmt.Printf("%s: ", message)
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		result := strings.TrimSpace(scanner.Text())
		logger.Debug("Input result", "result", result)
		return result, nil
	}
	return "", scanner.Err()
}

// PromptSelect prompts user to select from options
func PromptSelect(message string, options []string) (int, error) {
	logger.Debug("Prompt select", "message", message, "options", options)
	fmt.Println(message)
	for i, option := range options {
		fmt.Printf("  %d) %s\n", i+1, option)
	}

	for {
		fmt.Printf("Select (1-%d): ", len(options))
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			choice, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
			if err == nil && choice >= 1 && choice <= len(options) {
				result := choice - 1
				logger.Debug("Select result", "choice", choice, "result", result, "option", options[result])
				return result, nil
			}
		}
		logger.Debug("Invalid choice", "input", scanner.Text(), "valid_range", fmt.Sprintf("1-%d", len(options)))
		fmt.Printf("Invalid choice. Please select 1-%d.\n", len(options))
	}
}
