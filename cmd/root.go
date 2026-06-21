package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "godex",
	Short: "Go developer toolbox — Java, Node.js, and port manager",
	Long:  "godex is a developer toolbox for managing Java/Node.js versions, monitoring ports, and more.",
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
}
