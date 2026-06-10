package runtime

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Install struct {
	Name string
	Path string
}

func expandHome(path string) string {
	if path == "" || path[0] != '~' {
		return path
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}
	if path == "~" {
		return home
	}
	return filepath.Join(home, strings.TrimPrefix(path, "~/"))
}

func listDirs(roots ...string) ([]Install, error) {
	var installs []Install
	for _, root := range roots {
		root = expandHome(root)
		entries, err := os.ReadDir(root)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, err
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			installs = append(installs, Install{
				Name: entry.Name(),
				Path: filepath.Join(root, entry.Name()),
			})
		}
	}
	sort.Slice(installs, func(i, j int) bool {
		return installs[i].Name < installs[j].Name
	})
	return installs, nil
}

func findInstall(query string, installs []Install) (Install, error) {
	queryPath := expandHome(query)
	if stat, err := os.Stat(queryPath); err == nil && stat.IsDir() {
		return Install{Name: filepath.Base(queryPath), Path: queryPath}, nil
	}
	for _, install := range installs {
		if install.Name == query || strings.Contains(install.Name, query) {
			return install, nil
		}
	}
	return Install{}, errors.New("version or path not found: " + query)
}
