package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/hades/cli/internal/runtime"
	"github.com/spf13/cobra"
)

var javaCmd = &cobra.Command{
	Use:   "java",
	Short: "Manage Java versions",
}

var javaListCmd = &cobra.Command{
	Use:   "list",
	Short: "List discovered Java installations",
	RunE: func(cmd *cobra.Command, args []string) error {
		versions, err := runtime.ListJava()
		if err != nil {
			return err
		}
		if len(versions) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "No Java installations found.")
			return nil
		}
		for _, version := range versions {
			fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\n", version.Name, version.Path)
		}
		return nil
	},
}

var javaCurrentCmd = &cobra.Command{
	Use:   "current",
	Short: "Print the current Java version",
	RunE: func(cmd *cobra.Command, args []string) error {
		current, err := runtime.CurrentJava()
		if err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), strings.TrimSpace(current))
		return nil
	},
}

var javaUseCmd = &cobra.Command{
	Use:   "use <version-or-path>",
	Short: "Print shell exports for using a Java version",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		install, err := runtime.FindJava(args[0])
		if err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "export JAVA_HOME=%q\n", install.Path)
		fmt.Fprintf(cmd.OutOrStdout(), "export PATH=%q:$PATH\n", install.Path+string(os.PathSeparator)+"bin")
		return nil
	},
}

func init() {
	javaCmd.AddCommand(javaListCmd)
	javaCmd.AddCommand(javaCurrentCmd)
	javaCmd.AddCommand(javaUseCmd)
}
