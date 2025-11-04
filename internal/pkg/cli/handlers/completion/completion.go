package completion

import (
	"context"
	"fmt"
	"os"

	"github.com/otto-nation/otto-stack/internal/pkg/constants"
	"github.com/otto-nation/otto-stack/internal/pkg/types"
	"github.com/spf13/cobra"
)

type CompletionHandler struct{}

func NewCompletionHandler() *CompletionHandler {
	return &CompletionHandler{}
}

func (h *CompletionHandler) ValidateArgs(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf(constants.MsgCompletionRequiresOneArg.Content,
			fmt.Sprintf("%v", types.AllShellTypeStrings()))
	}

	shell := types.ShellType(args[0])
	if !shell.IsValid() {
		return fmt.Errorf(constants.MsgUnsupportedShell.Content,
			args[0], pkgTypes.AllShellTypeStrings())
	}

	return nil
}

func (h *CompletionHandler) GetRequiredFlags() []string {
	return []string{}
}

func (h *CompletionHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *types.BaseCommand) error {
	shell := types.ShellType(args[0])

	// Get the root command to generate completion for
	rootCmd := cmd.Root()

	switch shell {
	case types.ShellTypeBash:
		return rootCmd.GenBashCompletion(os.Stdout)
	case types.ShellTypeZsh:
		return rootCmd.GenZshCompletion(os.Stdout)
	case types.ShellTypeFish:
		return rootCmd.GenFishCompletion(os.Stdout, true)
	case types.ShellTypePowerShell:
		return rootCmd.GenPowerShellCompletion(os.Stdout)
	default:
		return fmt.Errorf("unsupported shell: %s", shell)
	}
}
