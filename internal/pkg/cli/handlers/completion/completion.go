package completion

import (
	"context"
	"fmt"
	"os"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/base"
	"github.com/otto-nation/otto-stack/internal/pkg/ci"
	"github.com/spf13/cobra"
)

type CompletionHandler struct{}

func NewCompletionHandler() *CompletionHandler {
	return &CompletionHandler{}
}

func (h *CompletionHandler) ValidateArgs(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf(core.MsgCompletion_requires_one_arg,
			fmt.Sprintf("%v", ci.AllShellTypeStrings()))
	}

	shell := ci.ShellType(args[0])
	if !shell.IsValid() {
		return fmt.Errorf(core.MsgCompletion_unsupported_shell,
			args[0])
	}

	return nil
}

func (h *CompletionHandler) GetRequiredFlags() []string {
	return []string{}
}

func (h *CompletionHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
	if err := h.ValidateArgs(args); err != nil {
		return err
	}

	shell := ci.ShellType(args[0])
	rootCmd := cmd.Root()

	switch shell {
	case ci.ShellTypeBash:
		return rootCmd.GenBashCompletion(os.Stdout)
	case ci.ShellTypeZsh:
		return rootCmd.GenZshCompletion(os.Stdout)
	case ci.ShellTypeFish:
		return rootCmd.GenFishCompletion(os.Stdout, true)
	case ci.ShellTypePowerShell:
		return rootCmd.GenPowerShellCompletion(os.Stdout)
	default:
		return fmt.Errorf(core.MsgCompletion_unsupported_shell, shell)
	}
}
