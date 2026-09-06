// Package migration locates migration files for analysis.
package migration

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
)

// Discover returns the SQL migration files represented by path. Directories
// are searched recursively so repositories may organize migrations in nested
// folders.
func Discover(path string) ([]string, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("migration path %q does not exist", path)
		}

		return nil, fmt.Errorf("inspect migration path %q: %w", path, err)
	}

	if !info.IsDir() {
		if filepath.Ext(path) != ".sql" {
			return nil, fmt.Errorf("migration file %q must have a .sql extension", path)
		}

		return []string{filepath.Clean(path)}, nil
	}

	var migrations []string
	err = filepath.WalkDir(path, func(currentPath string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return fmt.Errorf("walk %q: %w", currentPath, walkErr)
		}
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".sql" {
			migrations = append(migrations, filepath.Clean(currentPath))
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("discover migrations in %q: %w", path, err)
	}
	if len(migrations) == 0 {
		return nil, fmt.Errorf("no SQL migrations found in %q", path)
	}

	sort.Strings(migrations)

	return migrations, nil
}
