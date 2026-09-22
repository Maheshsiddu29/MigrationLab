package migration

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestLoadDiscoveredMigration(t *testing.T) {
	t.Parallel()

	source := `CREATE TABLE users (
    id bigint PRIMARY KEY
);`
	path := writeMigrationFile(t, t.TempDir(), "001_create_users.sql", source)

	migrations, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if len(migrations) != 1 {
		t.Fatalf("Load() returned %d migrations, want 1", len(migrations))
	}
	if migrations[0].Path != path {
		t.Errorf("migration path = %q, want %q", migrations[0].Path, path)
	}
	if migrations[0].Source != source {
		t.Errorf("migration source = %q, want %q", migrations[0].Source, source)
	}
	if len(migrations[0].Statements) != 1 {
		t.Fatalf("migration has %d statements, want 1", len(migrations[0].Statements))
	}
	if migrations[0].Statements[0].Node.GetCreateStmt() == nil {
		t.Fatal("migration statement AST is not a CREATE statement")
	}
}

func TestLoadMultipleMigrationsPreservesDiscoveryOrder(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	writeMigrationFile(t, dir, "020_add_index.sql", `CREATE INDEX idx_users_email ON users(email);`)
	writeMigrationFile(t, dir, "001_create_users.sql", `CREATE TABLE users (id bigint, email text);`)
	writeMigrationFile(t, dir, "010_add_status.sql", `ALTER TABLE users ADD COLUMN status text;
UPDATE users SET status = 'active';`)

	migrations, err := Load(dir)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if len(migrations) != 3 {
		t.Fatalf("Load() returned %d migrations, want 3", len(migrations))
	}

	gotPaths := make([]string, 0, len(migrations))
	gotStatementCounts := make([]int, 0, len(migrations))
	for _, migration := range migrations {
		gotPaths = append(gotPaths, filepath.Base(migration.Path))
		gotStatementCounts = append(gotStatementCounts, len(migration.Statements))
	}
	wantPaths := []string{"001_create_users.sql", "010_add_status.sql", "020_add_index.sql"}
	if !reflect.DeepEqual(gotPaths, wantPaths) {
		t.Errorf("migration paths = %v, want %v", gotPaths, wantPaths)
	}
	wantStatementCounts := []int{1, 2, 1}
	if !reflect.DeepEqual(gotStatementCounts, wantStatementCounts) {
		t.Errorf("statement counts = %v, want %v", gotStatementCounts, wantStatementCounts)
	}
}

func TestLoadInvalidMigrationIdentifiesFilename(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := writeMigrationFile(t, dir, "002_invalid.sql", `ALTER TABLE users ADD COLUMN;`)

	_, err := Load(dir)
	if err == nil {
		t.Fatal("Load() error = nil, want parser error")
	}
	if !strings.Contains(err.Error(), path) {
		t.Errorf("Load() error = %q, want migration path %q", err, path)
	}
	if !strings.Contains(err.Error(), "failed to parse migration") {
		t.Errorf("Load() error = %q, want parse context", err)
	}
	if !strings.Contains(err.Error(), "syntax error") {
		t.Errorf("Load() error = %q, want PostgreSQL parser detail", err)
	}
}

func TestLoadReadErrorIncludesFilename(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "001_unreadable.sql")
	if err := os.Symlink(filepath.Join(dir, "missing.sql"), path); err != nil {
		t.Skipf("create broken migration symlink: %v", err)
	}

	_, err := Load(dir)
	if err == nil {
		t.Fatal("Load() error = nil, want read error")
	}
	if !strings.Contains(err.Error(), path) {
		t.Errorf("Load() error = %q, want migration path %q", err, path)
	}
	if !strings.Contains(err.Error(), "read migration") {
		t.Errorf("Load() error = %q, want read context", err)
	}
}

func writeMigrationFile(t *testing.T, dir, name, source string) string {
	t.Helper()

	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(source), 0o600); err != nil {
		t.Fatalf("write migration %q: %v", path, err)
	}

	return path
}
