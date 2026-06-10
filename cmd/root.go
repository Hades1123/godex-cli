package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "devvm",
	Short: "Manage local Java and Node.js versions",
	Long:  "devvm manages local Java and Node.js versions and can print shell exports for switching versions.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(javaCmd)
	rootCmd.AddCommand(nodeCmd)
	rootCmd.AddCommand(guiCmd)
}
