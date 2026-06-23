package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:     "claude",
	Aliases: []string{"cld"},
	Short:   "Manage Claude Code settings presets",
	Long:    "Switch between Claude Code settings.json presets (e.g., glm ↔ deepseek).",
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available config presets",
	RunE: func(cmd *cobra.Command, args []string) error {
		presets, err := listPresets()
		if err != nil {
			return err
		}
		if len(presets) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "No presets found in ~/.config/godex/presets/")
			return nil
		}
		for _, name := range presets {
			fmt.Fprintf(cmd.OutOrStdout(), "  %s\n", name)
		}
		return nil
	},
}

var configUseCmd = &cobra.Command{
	Use:   "use <preset>",
	Short: "Switch to a config preset (copies to ~/.claude/settings.json)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		return switchPreset(name)
	},
}

var configCurrentCmd = &cobra.Command{
	Use:   "current",
	Aliases: []string{"cur"},
	Short: "Show the active Claude Code settings",
	RunE: func(cmd *cobra.Command, args []string) error {
	
		settingsPath, err := claudeSettingsPath()
		if err != nil {
			return err
		}

		cfg, err := readClaudeSettings(settingsPath)
		if err != nil {
			return err
		}

		model, _ := cfg["model"].(string)
		if model == "" {
			model = "(not set)"
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Model:  %s\n", model)

		if env, ok := cfg["env"].(map[string]any); ok {
			if baseURL, ok := env["ANTHROPIC_BASE_URL"].(string); ok {
				fmt.Fprintf(cmd.OutOrStdout(), "API:    %s\n", baseURL)
			}
		}

		return nil
	},
}

//! Only change current settings.json
//! We must change the preset
var claudeChangeApiCmd = &cobra.Command{
	Use: "api <key>",
	Short: "Change current api key (.claude/settings.json)",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		settingsPath, err := claudeSettingsPath()
		if err != nil{
			return err
		}

		cfg, err := readClaudeSettings(settingsPath)
		if err != nil {
			return err
		}

		if env, ok := cfg["env"].(map[string]any); ok {
			env["ANTHROPIC_AUTH_TOKEN"] = args[0]
		}

		err = writeClaudeSettings(cfg, settingsPath, 0644)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configUseCmd)
	configCmd.AddCommand(configCurrentCmd)
	configCmd.AddCommand(claudeChangeApiCmd)
	rootCmd.AddCommand(configCmd)
}
