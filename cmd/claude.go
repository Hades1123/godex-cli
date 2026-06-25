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

		currentPreset := cfg["godex"].(string)
		presetPath, err := presetPath(currentPreset)
		if err != nil {
			return err
		}

		presetConfig, err := readClaudeSettings(presetPath)
		if err != nil {
			return err
		}

		if presetEnv, ok := presetConfig["env"].(map[string]any); ok {
			presetEnv["ANTHROPIC_AUTH_TOKEN"] = args[0]
		}

		if env, ok := cfg["env"].(map[string]any); ok {
			env["ANTHROPIC_AUTH_TOKEN"] = args[0]
		}

		if err := writeClaudeSettings(cfg, settingsPath, 0644); err != nil{
			return nil
		}

		if err := writeClaudeSettings(presetConfig, presetPath, 0644) ; err != nil {
			return nil	
		}
		
		return nil
	},
}

// ---- template ----

var configTemplateListCmd = &cobra.Command{
	Use:   "template",
	Short: "List available templates",
	Aliases: []string{"tpl"},
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, t := range listTemplates() {
			fmt.Fprintf(cmd.OutOrStdout(), "  %s\n", t)
		}
		return nil
	},
}

var configTemplateInstallCmd = &cobra.Command{
	Use:   "install <name>",
	Aliases: []string{"i"},
	Short: "Download a template from GitHub and add to presets",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		if err := installTemplate(name); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Template %q installed to presets.\n", name)
		fmt.Fprintf(cmd.OutOrStdout(), "  → edit API key: godex claude api <key>\n")
		fmt.Fprintf(cmd.OutOrStdout(), "  → activate:      godex claude use %s\n", name)
		return nil
	},
}

func init() {
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configUseCmd)
	configCmd.AddCommand(configCurrentCmd)
	configCmd.AddCommand(claudeChangeApiCmd)
	configCmd.AddCommand(configTemplateListCmd)
	configCmd.AddCommand(configTemplateInstallCmd)
	rootCmd.AddCommand(configCmd)
}
