// Package parser provides PostgreSQL-aware syntax parsing for migrations.
package parser

import (
	"fmt"

	pg_query "github.com/pganalyze/pg_query_go/v6"
)

// ParsedMigration retains the original source and its ordered PostgreSQL AST
// statements. Empty and whitespace-only sources are valid migrations with no
// statements.
type ParsedMigration struct {
	Source     string
	Statements []Statement
}

// Statement is a top-level PostgreSQL statement in source order. Position is
// zero-based within the migration. SQL is the original byte span PostgreSQL
// assigns to the statement and may include leading comments or whitespace.
type Statement struct {
	Position int
	SQL      string
	Node     *pg_query.Node
}

// Parse parses a migration using PostgreSQL's parser and returns its top-level
// statements. It does not perform migration safety analysis.
func Parse(source string) (*ParsedMigration, error) {
	result, err := pg_query.Parse(source)
	if err != nil {
		return nil, fmt.Errorf("parse PostgreSQL migration: %w", err)
	}

	parsed := &ParsedMigration{
		Source:     source,
		Statements: make([]Statement, 0, len(result.GetStmts())),
	}
	for position, rawStatement := range result.GetStmts() {
		statementSource, err := sourceForStatement(source, rawStatement)
		if err != nil {
			return nil, fmt.Errorf("extract PostgreSQL statement %d: %w", position, err)
		}

		parsed.Statements = append(parsed.Statements, Statement{
			Position: position,
			SQL:      statementSource,
			Node:     rawStatement.GetStmt(),
		})
	}

	return parsed, nil
}

// sourceForStatement uses PostgreSQL's byte offsets as the sole authority for
// statement boundaries. MigrationLab never determines PostgreSQL statement
// boundaries by raw semicolon splitting: semicolons may be valid content inside
// strings, dollar-quoted function bodies, and PL/pgSQL blocks.
func sourceForStatement(source string, statement *pg_query.RawStmt) (string, error) {
	start := int(statement.GetStmtLocation())
	length := int(statement.GetStmtLen())
	if start < 0 || start > len(source) {
		return "", fmt.Errorf("parser returned invalid start offset %d for source length %d", start, len(source))
	}
	if length < 0 || length > len(source)-start {
		return "", fmt.Errorf("parser returned invalid length %d at offset %d for source length %d", length, start, len(source))
	}

	end := len(source)
	if length > 0 {
		end = start + length
		// PostgreSQL's length ends immediately before the delimiter for a
		// non-final statement. Including that byte preserves the original SQL
		// without using it to discover the boundary.
		if end < len(source) && source[end] == ';' {
			end++
		}
	}

	return source[start:end], nil
}
