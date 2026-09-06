package migration

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestDiscoverSQLFile(t *testing.T) {
	t.Parallel()

	path := writeTestFile(t, "001_create_users.sql")

	got, err := Discover(path)
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}
	want := []string{path}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Discover() = %v, want %v", got, want)
	}
}

func TestDiscoverDirectory(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "001_create_users.sql"))
	writeFile(t, filepath.Join(dir, "README.md"))
	writeFile(t, filepath.Join(dir, "nested", "002_add_email.sql"))

	got, err := Discover(dir)
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}
	want := []string{
		filepath.Join(dir, "001_create_users.sql"),
		filepath.Join(dir, "nested", "002_add_email.sql"),
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Discover() = %v, want %v", got, want)
	}
}

func TestDiscoverSortsDeterministically(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	for _, name := range []string{"020_last.sql", "003_middle.sql", "001_first.sql"} {
		writeFile(t, filepath.Join(dir, name))
	}

	got, err := Discover(dir)
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}
	want := []string{
		filepath.Join(dir, "001_first.sql"),
		filepath.Join(dir, "003_middle.sql"),
		filepath.Join(dir, "020_last.sql"),
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Discover() = %v, want %v", got, want)
	}
}

func TestDiscoverNonexistentPath(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "missing")

	_, err := Discover(path)
	if err == nil || !strings.Contains(err.Error(), "does not exist") {
		t.Fatalf("Discover() error = %v, want a does-not-exist error", err)
	}
}

func TestDiscoverRejectsNonSQLFile(t *testing.T) {
	t.Parallel()

	path := writeTestFile(t, "notes.txt")

	_, err := Discover(path)
	if err == nil || !strings.Contains(err.Error(), ".sql extension") {
		t.Fatalf("Discover() error = %v, want a .sql extension error", err)
	}
}

func TestDiscoverEmptyDirectory(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "README.md"))

	_, err := Discover(dir)
	if err == nil || !strings.Contains(err.Error(), "no SQL migrations") {
		t.Fatalf("Discover() error = %v, want a no-migrations error", err)
	}
}

func writeTestFile(t *testing.T, name string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), name)
	writeFile(t, path)

	return path
}

func writeFile(t *testing.T, path string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create test directory: %v", err)
	}
	if err := os.WriteFile(path, []byte("-- test migration\n"), 0o600); err != nil {
		t.Fatalf("write test file: %v", err)
	}
}
