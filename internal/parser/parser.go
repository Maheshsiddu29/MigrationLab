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
// zero-based within the migration.
type Statement struct {
	Position int
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
		parsed.Statements = append(parsed.Statements, Statement{
			Position: position,
			Node:     rawStatement.GetStmt(),
		})
	}

	return parsed, nil
}
