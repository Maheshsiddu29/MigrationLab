package parser

import (
	"strings"
	"testing"
)

func TestParseSingleStatement(t *testing.T) {
	t.Parallel()

	source := `CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email TEXT NOT NULL
);`

	parsed, err := Parse(source)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if parsed.Source != source {
		t.Fatal("Parse() did not retain the original source")
	}
	if len(parsed.Statements) != 1 {
		t.Fatalf("Parse() returned %d statements, want 1", len(parsed.Statements))
	}
	if parsed.Statements[0].Position != 0 {
		t.Fatalf("statement position = %d, want 0", parsed.Statements[0].Position)
	}
	if parsed.Statements[0].Node.GetCreateStmt() == nil {
		t.Fatal("statement AST is not a CREATE statement")
	}
}

func TestParseMultipleStatementsPreservesOrder(t *testing.T) {
	t.Parallel()

	source := `ALTER TABLE users
ADD COLUMN status TEXT;

CREATE INDEX idx_users_status
ON users(status);`

	parsed, err := Parse(source)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if len(parsed.Statements) != 2 {
		t.Fatalf("Parse() returned %d statements, want 2", len(parsed.Statements))
	}
	if parsed.Statements[0].Position != 0 || parsed.Statements[0].Node.GetAlterTableStmt() == nil {
		t.Fatal("first statement is not the ALTER TABLE statement at position 0")
	}
	if parsed.Statements[1].Position != 1 || parsed.Statements[1].Node.GetIndexStmt() == nil {
		t.Fatal("second statement is not the CREATE INDEX statement at position 1")
	}
}

func TestParsePostgreSQLConcurrentIndex(t *testing.T) {
	t.Parallel()

	source := `CREATE INDEX CONCURRENTLY idx_users_email
ON users(email);`

	parsed, err := Parse(source)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if len(parsed.Statements) != 1 {
		t.Fatalf("Parse() returned %d statements, want 1", len(parsed.Statements))
	}
	indexStatement := parsed.Statements[0].Node.GetIndexStmt()
	if indexStatement == nil {
		t.Fatal("statement AST is not an index statement")
	}
	if !indexStatement.GetConcurrent() {
		t.Fatal("index statement AST does not retain CONCURRENTLY")
	}
}

func TestParseAlterTableConstraint(t *testing.T) {
	t.Parallel()

	source := `ALTER TABLE users
ADD CONSTRAINT users_email_unique UNIQUE (email);`

	parsed, err := Parse(source)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if len(parsed.Statements) != 1 {
		t.Fatalf("Parse() returned %d statements, want 1", len(parsed.Statements))
	}
	if parsed.Statements[0].Node.GetAlterTableStmt() == nil {
		t.Fatal("statement AST is not an ALTER TABLE statement")
	}
}

func TestParseInvalidSyntax(t *testing.T) {
	t.Parallel()

	_, err := Parse(`ALTER TABLE users
ADD COLUMN;`)
	if err == nil {
		t.Fatal("Parse() error = nil, want invalid syntax error")
	}
	if !strings.Contains(err.Error(), "parse PostgreSQL migration") {
		t.Fatalf("Parse() error = %q, want parser context", err)
	}
}

func TestParseEmptyMigration(t *testing.T) {
	t.Parallel()

	parsed, err := Parse("")
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if parsed.Source != "" {
		t.Fatalf("Parse() source = %q, want empty source", parsed.Source)
	}
	if len(parsed.Statements) != 0 {
		t.Fatalf("Parse() returned %d statements, want 0", len(parsed.Statements))
	}
}

func TestParseWhitespaceOnlyMigration(t *testing.T) {
	t.Parallel()

	source := " \n\t\r\n"

	parsed, err := Parse(source)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if parsed.Source != source {
		t.Fatal("Parse() did not retain whitespace-only source")
	}
	if len(parsed.Statements) != 0 {
		t.Fatalf("Parse() returned %d statements, want 0", len(parsed.Statements))
	}
}
