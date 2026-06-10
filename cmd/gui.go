package cmd

import (
	"github.com/hades/cli/internal/ui"

	"github.com/spf13/cobra"
)

var guiCmd = &cobra.Command{
	Use:     "tui",
	Aliases: []string{"gui"},
	Short:   "Open the terminal interface",
	RunE: func(cmd *cobra.Command, args []string) error {
		return ui.Run(cmd.OutOrStdout())
	},
}
