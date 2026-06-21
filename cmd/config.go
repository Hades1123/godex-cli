package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:     "config",
	Aliases: []string{"cfg"},
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
		current := currentPresetName()
		for _, name := range presets {
			if name == current {
				fmt.Fprintf(cmd.OutOrStdout(), "* %s\n", name)
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "  %s\n", name)
			}
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
	Short: "Show the active Claude Code settings",
	RunE: func(cmd *cobra.Command, args []string) error {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("cannot find home dir: %w", err)
		}
		settingsPath := filepath.Join(home, ".claude", "settings.json")
		data, err := os.ReadFile(settingsPath)
		if err != nil {
			return fmt.Errorf("cannot read %s: %w", settingsPath, err)
		}

		var cfg map[string]any
		if err := json.Unmarshal(data, &cfg); err != nil {
			return fmt.Errorf("invalid JSON in %s: %w", settingsPath, err)
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

		// Detect which preset matches.
		if name := currentPresetName(); name != "" {
			fmt.Fprintf(cmd.OutOrStdout(), "Preset: %s\n", name)
		}
		return nil
	},
}

func init() {
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configUseCmd)
	configCmd.AddCommand(configCurrentCmd)
	rootCmd.AddCommand(configCmd)
}

// --- helpers ---

func presetDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot find home dir: %w", err)
	}
	dir := filepath.Join(home, ".config", "godex", "presets")
	// Ensure the directory exists.
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("cannot create preset dir: %w", err)
	}
	return dir, nil
}

func listPresets() ([]string, error) {
	dir, err := presetDir()
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("cannot read preset dir: %w", err)
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		names = append(names, strings.TrimSuffix(e.Name(), ".json"))
	}
	return names, nil
}

func presetPath(name string) (string, error) {
	dir, err := presetDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, name+".json"), nil
}

func switchPreset(name string) error {
	src, err := presetPath(name)
	if err != nil {
		return err
	}

	// Check the preset file exists.
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return fmt.Errorf("preset %q not found at %s", name, src)
	}

	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("cannot read preset: %w", err)
	}

	// Validate JSON.
	var js map[string]any
	if err := json.Unmarshal(data, &js); err != nil {
		return fmt.Errorf("preset %q contains invalid JSON: %w", name, err)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("cannot find home dir: %w", err)
	}

	// Ensure ~/.claude directory exists.
	claudeDir := filepath.Join(home, ".claude")
	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		return fmt.Errorf("cannot create ~/.claude: %w", err)
	}

	target := filepath.Join(claudeDir, "settings.json")

	// Backup existing settings if any.
	if _, err := os.Stat(target); err == nil {
		backup := target + ".bak"
		if err := copyFile(target, backup); err == nil {
			fmt.Printf("Backed up previous settings to %s\n", backup)
		}
	}

	if err := os.WriteFile(target, data, 0644); err != nil {
		return fmt.Errorf("cannot write %s: %w", target, err)
	}

	fmt.Printf("Switched to preset %q → %s\n", name, target)
	return nil
}

func currentPresetName() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	settingsPath := filepath.Join(home, ".claude", "settings.json")
	settingsData, err := os.ReadFile(settingsPath)
	if err != nil {
		return ""
	}
	var current map[string]any
	if err := json.Unmarshal(settingsData, &current); err != nil {
		return ""
	}

	dir, err := presetDir()
	if err != nil {
		return ""
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			continue
		}
		var preset map[string]any
		if err := json.Unmarshal(data, &preset); err != nil {
			continue
		}
		// Compare model and base URL to identify the preset.
		if preset["model"] == current["model"] {
			if presetEnv, ok := preset["env"].(map[string]any); ok {
				if currentEnv, ok := current["env"].(map[string]any); ok {
					if presetEnv["ANTHROPIC_BASE_URL"] == currentEnv["ANTHROPIC_BASE_URL"] {
						return strings.TrimSuffix(e.Name(), ".json")
					}
				}
			}
		}
	}
	return ""
}

func copyFile(src, dst string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()

	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer d.Close()

	_, err = io.Copy(d, s)
	return err
}
