# Architecture

MigrationLab is a staged PostgreSQL migration safety engine. Phase 2 implements migration discovery, loading, PostgreSQL-aware parsing, source preservation, and CLI syntax validation. Safety analysis, schema modeling, and runtime simulation remain planned.

## Implemented parsing pipeline

```text
CLI
 │
 ▼
Migration Discovery
 │
 ▼
Migration Loader
 │
 ▼
PostgreSQL Parser
 │
 ▼
PostgreSQL AST
```

### CLI

`cmd/migrationlab` owns command definitions and user-facing output. The `analyze` command accepts one file or directory, invokes the migration pipeline, reports statement counts, and returns an error for invalid PostgreSQL syntax. It does not contain SQL parsing or safety rules.

### Migration discovery and loading

`internal/migration` owns filesystem concerns:

- validates file and directory paths;
- discovers `.sql` files recursively;
- sorts paths lexicographically for deterministic results;
- reads each file;
- sends source text to the parser; and
- adds the migration path to read and parser errors.

The package returns ordered migration files containing their path, original source, and parsed statements.

### PostgreSQL parser

`internal/parser` owns PostgreSQL syntax parsing. It wraps `github.com/pganalyze/pg_query_go/v6` and returns the smallest representation needed by current callers: original source, ordered top-level statements, parser-defined statement SQL, and PostgreSQL AST nodes.

The parser has no dependency on Cobra, CLI output, Docker, or GitHub Actions. Empty and whitespace-only sources are valid zero-statement migrations. Invalid syntax is returned as a contextual error; expected parser failures do not produce stack traces or panic.

## Why PostgreSQL's parser is authoritative

PostgreSQL syntax cannot be interpreted reliably with regular expressions, a handwritten lexer, or generic SQL assumptions. PostgreSQL-specific constructs and grammar decisions must be represented exactly as PostgreSQL understands them.

MigrationLab therefore treats the PostgreSQL AST as the authoritative syntax representation. Later analyzers will inspect AST nodes rather than searching SQL text for keywords.

### Statement-boundary invariant

> MigrationLab never determines PostgreSQL statement boundaries by raw semicolon splitting.

Semicolons can appear inside string literals, dollar-quoted function bodies, and PL/pgSQL blocks. Splitting source text on `;` can turn one valid statement into several invalid fragments.

The parser instead uses PostgreSQL's byte-based statement location and length metadata. PostgreSQL's spans associate leading comments and inter-statement whitespace with the following statement. MigrationLab preserves those spans, retains a terminating semicolon only at a parser-reported boundary, and keeps the entire original migration source.

## Parsing and analysis are separate responsibilities

Parsing identifies what PostgreSQL syntax represents. Analysis decides whether that operation is dangerous.

```text
Parser:
"This is CREATE INDEX."

Analyzer:
"This CREATE INDEX may be unsafe because it does not use CONCURRENTLY."
```

Phase 2 implements only the parser responsibility. It does not assign findings, severity, risk, or remediation guidance.

## Planned downstream architecture

```text
PostgreSQL AST
 │
 ├── Static Analyzer (planned)
 ├── Schema Model (planned)
 └── Simulation Engine (planned)
       │
       ├── PostgreSQL Sandbox
       ├── Workload Generator
       ├── Lock/Performance Collector
       └── Rollback Verifier
              │
              ▼
          Risk Report
```

- **Static Analyzer:** evaluate AST nodes against versioned migration rules.
- **Schema Model:** track schema transitions across ordered migrations.
- **PostgreSQL Sandbox:** apply migrations to an isolated database.
- **Workload Generator:** create controlled concurrent database activity.
- **Lock/Performance Collector:** measure locks, waits, and latency changes.
- **Rollback Verifier:** exercise and validate rollback behavior.
- **Risk Report:** combine static and runtime evidence for developers and CI.

New packages should be introduced when their pipeline stage is implemented, not as empty placeholders.

## CI and runtime boundaries

Parser tests are deterministic unit and package-integration tests and do not require PostgreSQL. CI verifies formatting, runs `go vet`, executes normal and race-enabled tests, and builds the CLI.

PostgreSQL 18 remains available through Docker Compose for future integration work, but CI does not start a PostgreSQL service during Phase 2.

## Planned release delivery

Release automation is intentionally deferred. The expected future delivery path is:

```text
Tag
 ↓
GitHub Actions release workflow
 ↓
Cross-platform Go binaries
 ↓
Checksums
 ↓
GitHub Release
```

No release workflow, package publication, or production deployment is implemented in Phase 2.
