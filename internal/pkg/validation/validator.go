package validation

import (
	"github.com/otto-nation/otto-stack/internal/pkg/config"
	"github.com/spf13/cobra"
)

// Validator provides comprehensive validation for CLI-YAML consistency
type Validator struct {
	config             *config.CommandConfig
	metadataValidator  *MetadataValidator
	commandValidator   *CommandValidator
	practicesValidator *BestPracticesValidator
	cliValidator       *CLIValidator
}

// NewValidator creates a new validator instance
func NewValidator(config *config.CommandConfig) *Validator {
	return &Validator{
		config:             config,
		metadataValidator:  NewMetadataValidator(config),
		commandValidator:   NewCommandValidator(config),
		practicesValidator: NewBestPracticesValidator(config),
		cliValidator:       NewCLIValidator(config),
	}
}

// ValidateAll performs comprehensive validation
func (v *Validator) ValidateAll() *ValidationResult {
	result := &ValidationResult{
		Valid: true,
		Summary: ValidationSummary{
			TotalCommands:   len(v.config.Commands),
			TotalCategories: len(v.config.Categories),
		},
	}

	v.metadataValidator.ValidateMetadata(result)
	v.metadataValidator.ValidateGlobalConfiguration(result)
	v.commandValidator.ValidateCategories(result)
	v.commandValidator.ValidateCommands(result)
	v.commandValidator.ValidateReferences(result)
	v.practicesValidator.ValidateBestPractices(result)

	v.calculateSummary(result)
	v.generateSuggestions(result)

	return result
}

// ValidateAgainstCLI validates the configuration against actual CLI implementation
func (v *Validator) ValidateAgainstCLI(rootCmd *cobra.Command) *ValidationResult {
	return v.cliValidator.ValidateAgainstCLI(rootCmd)
}

// calculateSummary calculates validation summary statistics
func (v *Validator) calculateSummary(result *ValidationResult) {
	result.Summary.ErrorCount = len(result.Errors)
	result.Summary.WarningCount = len(result.Warnings)

	for _, err := range result.Errors {
		if err.Severity == "critical" {
			result.Summary.CriticalErrors++
		}
	}

	totalIssues := result.Summary.ErrorCount + result.Summary.WarningCount
	if totalIssues == 0 {
		result.Summary.ConfigurationScore = 100.0
	} else {
		weightedScore := 100.0 - float64(result.Summary.ErrorCount*10+result.Summary.WarningCount*2)
		if weightedScore < 0 {
			weightedScore = 0
		}
		result.Summary.ConfigurationScore = weightedScore
	}

	result.Valid = result.Summary.CriticalErrors == 0 && result.Summary.ErrorCount < 5
}

// generateSuggestions generates improvement suggestions
func (v *Validator) generateSuggestions(result *ValidationResult) {
	if result.Summary.ErrorCount > 0 {
		result.Suggestions = append(result.Suggestions, "Fix validation errors to improve configuration quality")
	}

	if result.Summary.ConfigurationScore < 80 {
		result.Suggestions = append(result.Suggestions, "Consider improving documentation coverage and following best practices")
	}

	if len(v.config.Categories) == 0 {
		result.Suggestions = append(result.Suggestions, "Organize commands into categories for better CLI structure")
	}

}
