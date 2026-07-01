package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "godex",
	Short: "Go developer toolbox — Claude Code presets and network port manager",
	Long:  "godex is a developer toolbox for managing Claude Code presets and monitoring network ports.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
