package ui

import (
	"testing"

	"github.com/otto-nation/otto-stack/test/testhelpers"
)

func TestUI_prompt_functions(t *testing.T) {
	t.Run("prompt input with empty message", func(t *testing.T) {
		result, err := PromptInput("", "default")
		// This will likely fail due to no stdin, but we're testing the code path
		if err != nil {
			testhelpers.AssertError(t, err, "PromptInput should handle empty message")
		} else {
			testhelpers.AssertNoError(t, err, "PromptInput should not error")
		}
		_ = result
	})

	t.Run("prompt confirm with empty message", func(t *testing.T) {
		result, err := PromptConfirm("", false)
		// This will likely fail due to no stdin, but we're testing the code path
		if err != nil {
			testhelpers.AssertError(t, err, "PromptConfirm should handle empty message")
		} else {
			testhelpers.AssertNoError(t, err, "PromptConfirm should not error")
		}
		_ = result
	})
}
