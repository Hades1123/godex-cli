package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/hades/godex/internal/runtime"
	"github.com/spf13/cobra"
)

var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Manage Node.js versions",
}

var nodeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List discovered Node.js installations",
	RunE: func(cmd *cobra.Command, args []string) error {
		versions, err := runtime.ListNode()
		if err != nil {
			return err
		}
		if len(versions) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "No Node.js installations found.")
			return nil
		}
		for _, version := range versions {
			fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\n", version.Name, version.Path)
		}
		return nil
	},
}

var nodeCurrentCmd = &cobra.Command{
	Use:   "current",
	Short: "Print the current Node.js version",
	RunE: func(cmd *cobra.Command, args []string) error {
		current, err := runtime.CurrentNode()
		if err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), strings.TrimSpace(current))
		return nil
	},
}

var nodeUseCmd = &cobra.Command{
	Use:   "use <version-or-path>",
	Short: "Print shell exports for using a Node.js version",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		install, err := runtime.FindNode(args[0])
		if err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "export PATH=%q:$PATH\n", install.Path+string(os.PathSeparator)+"bin")
		return nil
	},
}

func init() {
	nodeCmd.AddCommand(nodeListCmd)
	nodeCmd.AddCommand(nodeCurrentCmd)
	nodeCmd.AddCommand(nodeUseCmd)
}
