# Architecture

MigrationLab is planned as a staged PostgreSQL migration safety pipeline. Phase 1 implements only migration loading: accepting a file or directory, validating it, discovering SQL files, and sorting them deterministically.

## Planned pipeline

```text
Migration Loader
      ↓
PostgreSQL Parser (planned)
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

### Future components — planned

- **PostgreSQL Parser:** convert PostgreSQL SQL into a syntax-aware representation.
- **Static Analyzer:** evaluate parsed statements against versioned safety rules.
- **Schema Model:** track the schema changes produced by a migration sequence.
- **PostgreSQL Sandbox:** apply migrations to an isolated PostgreSQL instance.
- **Workload Generator:** create controlled concurrent database activity.
- **Lock/Performance Collector:** measure locks, waits, and latency changes.
- **Rollback Verifier:** exercise and validate rollback behavior.
- **Risk Report:** combine evidence into human-readable and CI-friendly results.

## Design boundaries

The initial codebase has one internal package because file discovery is the only domain behavior implemented. New packages should be added when their corresponding pipeline stage is implemented, rather than as empty placeholders. PostgreSQL is available for local development but is not used by Phase 1 commands or tests.
