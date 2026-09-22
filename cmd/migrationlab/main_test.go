package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAnalyzeValidMigrationFile(t *testing.T) {
	t.Parallel()

	path := writeCLITestFile(t, "001_create_users.sql", `CREATE TABLE users (id bigint);`)

	output, err := executeCLI(t, "analyze", path)
	if err != nil {
		t.Fatalf("analyze command error = %v", err)
	}
	for _, want := range []string{
		"Analyzing migrations...",
		"001_create_users.sql",
		"parsed 1 statement",
		"1 migration",
		"1 statement",
	} {
		if !strings.Contains(output, want) {
			t.Errorf("analyze output missing %q:\n%s", want, output)
		}
	}
}

func TestAnalyzeDirectoryReportsStatementTotalsInOrder(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	writeCLIFile(t, filepath.Join(dir, "020_add_index.sql"), `CREATE INDEX idx_users_email ON users(email);`)
	writeCLIFile(t, filepath.Join(dir, "001_create_users.sql"), `CREATE TABLE users (id bigint);`)
	writeCLIFile(t, filepath.Join(dir, "010_add_status.sql"), `ALTER TABLE users ADD COLUMN status text;
UPDATE users SET status = 'active';`)

	output, err := executeCLI(t, "analyze", dir)
	if err != nil {
		t.Fatalf("analyze command error = %v", err)
	}
	for _, want := range []string{"3 migrations", "4 statements", "parsed 2 statements"} {
		if !strings.Contains(output, want) {
			t.Errorf("analyze output missing %q:\n%s", want, output)
		}
	}

	first := strings.Index(output, "001_create_users.sql")
	second := strings.Index(output, "010_add_status.sql")
	third := strings.Index(output, "020_add_index.sql")
	if first == -1 || second == -1 || third == -1 || !(first < second && second < third) {
		t.Errorf("migration output is not deterministic:\n%s", output)
	}
}

func TestAnalyzeInvalidSQLReturnsError(t *testing.T) {
	t.Parallel()

	path := writeCLITestFile(t, "002_invalid.sql", `ALTER TABLE users ADD COLUMN;`)

	output, err := executeCLI(t, "analyze", path)
	if err == nil {
		t.Fatal("analyze command error = nil, want parser error")
	}
	if !strings.Contains(output, "Analyzing migrations...") {
		t.Errorf("analyze output missing start message:\n%s", output)
	}
	if !strings.Contains(err.Error(), path) {
		t.Errorf("analyze error = %q, want migration path %q", err, path)
	}
	if !strings.Contains(err.Error(), "syntax error") {
		t.Errorf("analyze error = %q, want PostgreSQL parser detail", err)
	}
}

func TestAnalyzePreservesPathValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		path    func(*testing.T) string
		wantErr string
	}{
		{
			name: "nonexistent path",
			path: func(t *testing.T) string {
				return filepath.Join(t.TempDir(), "missing")
			},
			wantErr: "does not exist",
		},
		{
			name: "non-SQL file",
			path: func(t *testing.T) string {
				return writeCLITestFile(t, "notes.txt", "not SQL")
			},
			wantErr: ".sql extension",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			_, err := executeCLI(t, "analyze", test.path(t))
			if err == nil || !strings.Contains(err.Error(), test.wantErr) {
				t.Fatalf("analyze command error = %v, want error containing %q", err, test.wantErr)
			}
		})
	}
}

func executeCLI(t *testing.T, args ...string) (string, error) {
	t.Helper()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command := newRootCommand(&stdout, &stderr)
	command.SetArgs(args)

	err := command.Execute()

	return stdout.String(), err
}

func writeCLITestFile(t *testing.T, name, source string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), name)
	writeCLIFile(t, path, source)

	return path
}

func writeCLIFile(t *testing.T, path, source string) {
	t.Helper()

	if err := os.WriteFile(path, []byte(source), 0o600); err != nil {
		t.Fatalf("write CLI test file %q: %v", path, err)
	}
}
