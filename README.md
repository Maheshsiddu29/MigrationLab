# MigrationLab

MigrationLab is an early-stage PostgreSQL migration safety and simulation engine. Its goal is to give engineering teams reliable evidence about migration behavior before database changes reach production.

> **Project status:** Phase 2 PostgreSQL parsing is implemented. MigrationLab is under active development and validates syntax, but it does not yet determine whether a migration is safe.

## The problem

PostgreSQL migrations that look harmless can take disruptive locks, trigger expensive table work, degrade application latency, or fail to roll back safely. Code review alone often lacks the schema context and runtime evidence needed to evaluate those risks consistently.

MigrationLab is intended to combine PostgreSQL-aware static analysis, controlled execution, workload simulation, and measured evidence in one local and CI-friendly workflow. Only discovery and syntax parsing are implemented today.

## Implemented functionality

- ✓ Migration file and directory discovery
- ✓ Recursive `.sql` file discovery
- ✓ Deterministic migration ordering
- ✓ Migration file loading with contextual errors
- ✓ PostgreSQL AST parsing through `pg_query_go/v6`
- ✓ Multi-statement parsing in source order
- ✓ PostgreSQL syntax validation
- ✓ Original migration and statement SQL preservation
- ✓ CLI parsing validation with per-file and aggregate statement counts

PostgreSQL's parser is the source of truth for both syntax and statement boundaries. MigrationLab does not split SQL on raw semicolons, use regular expressions to identify statement types, or maintain a handwritten SQL parser.

## Planned functionality

- ○ Migration safety rules, including `ML001`–`ML012`
- ○ Risk severity evaluation
- ○ Schema-state modeling
- ○ PostgreSQL sandbox execution
- ○ Concurrent workload replay
- ○ Lock analysis
- ○ Performance regression analysis
- ○ Rollback verification
- ○ Developer-friendly and CI-friendly risk reports

These are roadmap items, not current product capabilities. A successfully parsed migration may still be unsafe to deploy.

## Architecture overview

The implemented Phase 2 pipeline is:

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

Future phases will use the AST as input to static analysis, schema modeling, and database simulation. Parsing answers what PostgreSQL syntax represents; analysis will decide whether that operation is dangerous.

See [docs/architecture.md](docs/architecture.md) for component ownership and parser invariants, and [docs/rules.md](docs/rules.md) for the planned rule namespace.

## Requirements

- Go 1.27.x
- A C toolchain supported by Go's cgo, required by `pg_query_go/v6`
- Docker with Docker Compose, only for the local PostgreSQL 18 development database

## Local development

```bash
make build          # Build bin/migrationlab
make test           # Run unit tests
make lint           # Verify formatting and run go vet
make run            # Run CLI help
make postgres-up    # Start the PostgreSQL 18 development database
make postgres-down  # Stop the development database
```

The parser and its tests do not require a running PostgreSQL server. Compose credentials are intentionally simple local-development defaults and must not be used in production.

## CLI examples

Build the CLI:

```bash
make build
```

Inspect commands and validate migration syntax:

```bash
./bin/migrationlab --help
./bin/migrationlab version
./bin/migrationlab analyze testdata/migrations
./bin/migrationlab analyze testdata/migrations/001_create_users.sql
```

Example result:

```text
Analyzing migrations...

001_create_users.sql
  parsed 1 statement

1 migration
1 statement
```

Invalid PostgreSQL syntax returns a non-zero exit status and identifies the affected file. The command does not report migration risk yet and does not print full ASTs by default.

## Roadmap

1. **Repository foundation — implemented:** CLI shell, deterministic discovery, tests, CI, documentation, and PostgreSQL development service.
2. **PostgreSQL AST parsing — implemented:** migration loading, syntax validation, safe multi-statement boundaries, original SQL preservation, and CLI integration.
3. **Static safety analysis — planned:** implement `ML001`–`ML012` using PostgreSQL AST nodes.
4. **Schema and simulation — planned:** model schema transitions and execute migrations in a PostgreSQL sandbox under concurrent workloads.
5. **Verification and reporting — planned:** measure locks and latency, verify rollback behavior, and emit CI-friendly risk reports.

## Delivery model

MigrationLab does not have an automated release pipeline yet. A later release phase is expected to implement:

```text
Git tag
   ↓
GitHub Actions release workflow
   ↓
Cross-platform Go binaries
   ↓
Checksums
   ↓
GitHub Release
```

No binaries or packages are currently published by CI.

## Contributing

Run `make lint`, `make test`, and `go test -race ./...` before proposing changes. Keep additions scoped to implemented roadmap work and avoid introducing infrastructure without a current requirement.
