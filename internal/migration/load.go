package migration

import (
	"fmt"
	"os"

	"github.com/Maheshsiddu29/MigrationLab/internal/parser"
)

// File is a loaded migration and its ordered PostgreSQL statements.
type File struct {
	Path       string
	Source     string
	Statements []parser.Statement
}

// Load discovers SQL migrations at path, reads them, and parses each one using
// PostgreSQL's parser. Files retain the deterministic order returned by
// Discover.
func Load(path string) ([]File, error) {
	paths, err := Discover(path)
	if err != nil {
		return nil, err
	}

	migrations := make([]File, 0, len(paths))
	for _, migrationPath := range paths {
		source, err := os.ReadFile(migrationPath)
		if err != nil {
			return nil, fmt.Errorf("read migration %q: %w", migrationPath, err)
		}

		parsed, err := parser.Parse(string(source))
		if err != nil {
			return nil, fmt.Errorf("failed to parse migration %q: %w", migrationPath, err)
		}

		migrations = append(migrations, File{
			Path:       migrationPath,
			Source:     parsed.Source,
			Statements: parsed.Statements,
		})
	}

	return migrations, nil
}
