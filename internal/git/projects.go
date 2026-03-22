package git

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetProjectsInDir(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var projects []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		_, err := os.Stat(filepath.Join(dir, entry.Name(), ".git"))
		if err != nil {
			continue
		}

		projects = append(projects, entry.Name())
	}

	if len(projects) == 0 {
		return nil, fmt.Errorf("no git projects found in %s", dir)
	}

	return projects, nil
}
