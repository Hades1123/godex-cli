package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func claudeSettingsPath() (string , error) {
	home, err := os.UserHomeDir()

	if err != nil {
		return "", fmt.Errorf("cannot read home dir: %w", err)
	}

	settingPath := filepath.Join(home, ".claude", "settings.json")

	return settingPath, nil
}

func readClaudeSettings(path string) (map[string]any, error) {
	data, err := os.ReadFile(path)
	if err != nil{
		return nil, fmt.Errorf("cannot read file %s : %w", path, err)
	}

	var cfg map[string]any 

	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("invalid JSON in %s : %w", path, err)
	}

	return cfg, nil
}

func writeClaudeSettings(cfg map[string]any, settingPath string, mode fs.FileMode) error {
	jsonData, err := json.MarshalIndent(cfg, "", " ")
	if err != nil {
		return fmt.Errorf("cannot parse json data")
	}

	if err := os.MkdirAll(filepath.Dir(settingPath), 0755); err != nil{
		return fmt.Errorf("cannot create ~/.claude dir: %w", err)
	}

	err = os.WriteFile(settingPath, jsonData, mode) 
	if err != nil {
		return fmt.Errorf("cannot write setting to file at dir %s : %w", settingPath, err)
	}

	return nil
}

func presetDir() (string, error) {
	// os.UserConfigDir is cross-platform:
	//   Linux/macOS: ~/.config        (or $XDG_CONFIG_HOME)
	//   Windows:     %AppData%        (e.g. C:\Users\<user>\AppData\Roaming)
	base, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("cannot find config dir: %w", err)
	}
	
	dir := filepath.Join(base, "godex", "presets")
	// Ensure the directory exists.
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("cannot create preset dir: %w", err)
	}
	
	return dir, nil
}

func presetPath(name string) (string, error) {
	dir, err := presetDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, name + ".json"), nil
}

func switchPreset(name string) error {
	src, err := presetPath(name)
	if err != nil {
		return err
	}

	if _, err := os.Stat(src); os.IsNotExist(err) {
		return fmt.Errorf("preset %q not found at %s", name, src)
	}

	cfg, err := readClaudeSettings(src)
	if err != nil {
		return err
	}

	cfg["godex"] = name

	target, err := claudeSettingsPath()
	if err != nil {
		return err
	}

	if err := writeClaudeSettings(cfg, target, 0644); err != nil {
		return err
	}

	fmt.Printf("Switched to preset %q → %s\n", name, target)
	return nil
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

// ---- template helpers ----

const templateBaseURL = "https://raw.githubusercontent.com/Hades1123/godex-cli/main/templates"

// templates is the curated list of presets available for download.
var templates = []string{"glm", "deepseek"}

func listTemplates() []string {
	return templates
}

func templateURL(name string) string {
	return templateBaseURL + "/" + name + ".json"
}

// installTemplate downloads a template JSON from GitHub and saves it to the
// local preset directory (~/.config/godex/presets/<name>.json).
func installTemplate(name string) error {
	if !slices.Contains(templates, name) {
		return fmt.Errorf("unknown template %q — run: godex claude template list", name)
	}

	url := templateURL(name)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("cannot download template: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("template %q not found at %s (HTTP %d)", name, url, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("cannot read response: %w", err)
	}

	// Validate JSON.
	var cfg map[string]any
	if err := json.Unmarshal(body, &cfg); err != nil {
		return fmt.Errorf("template %q is not valid JSON: %w", name, err)
	}

	dst, err := presetPath(name)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return fmt.Errorf("cannot create preset dir: %w", err)
	}

	if err := os.WriteFile(dst, body, 0644); err != nil {
		return fmt.Errorf("cannot write template to %s: %w", dst, err)
	}

	return nil
}
