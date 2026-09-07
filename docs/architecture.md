# Architecture

MigrationLab is planned as a staged PostgreSQL migration safety pipeline. Phase 1 implements only migration loading: accepting a file or directory, validating it, discovering SQL files, and sorting them deterministically.

## Planned pipeline

```text
Migration Loader
      ↓
PostgreSQL Parser (Phase 2 foundation implemented)
      ↓
Static Analyzer (planned)
      ↓
Schema Model (planned)
      ↓
PostgreSQL Sandbox (planned)
      ↓
Workload Generator (planned)
      ↓
Lock/Performance Collector (planned)
      ↓
Rollback Verifier (planned)
      ↓
Risk Report (planned)
```

### Migration Loader — implemented in Phase 1

The loader accepts one filesystem path. A SQL file is returned directly; a directory is searched recursively for files with a `.sql` extension. Results are sorted lexicographically by path to make local and CI output reproducible. It does not read or parse SQL contents.

### PostgreSQL Parser — Phase 2 foundation implemented

The parser converts migration source into PostgreSQL AST nodes, preserves statement order, and retains each parser-defined source span. It validates syntax but does not decide whether a statement is safe.

### Future components — planned

- **Static Analyzer:** evaluate parsed statements against versioned safety rules.
- **Schema Model:** track the schema changes produced by a migration sequence.
- **PostgreSQL Sandbox:** apply migrations to an isolated PostgreSQL instance.
- **Workload Generator:** create controlled concurrent database activity.
- **Lock/Performance Collector:** measure locks, waits, and latency changes.
- **Rollback Verifier:** exercise and validate rollback behavior.
- **Risk Report:** combine evidence into human-readable and CI-friendly results.

## Design boundaries

The initial codebase has one internal package because file discovery is the only domain behavior implemented. New packages should be added when their corresponding pipeline stage is implemented, rather than as empty placeholders. PostgreSQL is available for local development but is not used by Phase 1 commands or tests.

### Parser boundary invariant

MigrationLab never determines PostgreSQL statement boundaries by raw semicolon splitting. PostgreSQL permits semicolons inside string literals, dollar-quoted function bodies, and PL/pgSQL blocks, so a textual split can corrupt otherwise valid statements.

The parser package uses PostgreSQL's statement location and length metadata to retain each statement's original SQL. PostgreSQL's spans associate leading comments and inter-statement whitespace with the following statement; trailing whitespace outside the final statement span remains available in the migration-level source. A terminating semicolon is retained only when it appears exactly at a parser-reported boundary.
