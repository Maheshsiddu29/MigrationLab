# MigrationLab

MigrationLab is an early-stage PostgreSQL migration safety and simulation engine. Its goal is to help engineering teams understand migration risk before database changes reach production.

> **Project status:** Phase 1 (repository foundation) is implemented. MigrationLab is under active development and is not yet a migration safety analyzer.

## The problem

PostgreSQL migrations that look harmless can take disruptive locks, trigger expensive table work, degrade application latency, or fail to roll back safely. Code review alone often lacks the schema context and runtime evidence needed to evaluate those risks consistently.

MigrationLab is intended to bring static checks, controlled execution, workload simulation, and measured evidence into one local and CI-friendly workflow.

## Current functionality

The Phase 1 CLI can:

- accept a SQL file or directory;
- reject missing paths and unsupported files;
- recursively discover `.sql` files in a directory;
- sort discovered files deterministically; and
- print the migration paths it found.

It does **not** parse or analyze SQL yet.

## Planned capabilities

- PostgreSQL-aware SQL parsing
- Static migration safety rules
- Schema-state modeling
- Isolated PostgreSQL migration execution
- Realistic concurrent workload generation
- Lock and performance measurement
- Rollback verification
- Human-readable and CI-friendly risk reports

These capabilities are roadmap items and are not currently implemented.

## Architecture overview

The planned processing pipeline is:

```text
Migration Loader
→ PostgreSQL Parser
→ Static Analyzer
→ Schema Model
→ PostgreSQL Sandbox
→ Workload Generator
→ Lock/Performance Collector
→ Rollback Verifier
→ Risk Report
```

Only the Migration Loader's file-discovery behavior exists in Phase 1. See [docs/architecture.md](docs/architecture.md) for component boundaries and [docs/rules.md](docs/rules.md) for the reserved rule namespace.

## Requirements

- Go 1.27.x
- Docker with Docker Compose (only needed for the local PostgreSQL 18 database)

## Local development

```bash
make build          # Build bin/migrationlab
make test           # Run unit tests
make lint           # Verify formatting and run go vet
make run            # Run CLI help
make postgres-up    # Start the PostgreSQL 18 development database
make postgres-down  # Stop the development database
```

The Compose credentials are intentionally simple local-development defaults and must not be used in production.

## CLI examples

Build the command first:

```bash
make build
```

Then inspect the available commands and discover migrations:

```bash
./bin/migrationlab --help
./bin/migrationlab version
./bin/migrationlab analyze testdata/migrations
./bin/migrationlab analyze testdata/migrations/001_create_users.sql
```

The `analyze` command currently lists migration files only; its name reserves the user workflow for later analysis phases.

## Roadmap

1. **Repository foundation — implemented:** CLI shell, deterministic migration discovery, tests, CI, documentation, and PostgreSQL development service.
2. **SQL parsing — planned:** build a PostgreSQL-aware migration representation.
3. **Static safety analysis — planned:** implement the initial `ML001`–`ML012` rules.
4. **Simulation — planned:** execute migrations against a PostgreSQL sandbox under concurrent workloads.
5. **Verification and reporting — planned:** measure locks and latency, verify rollback behavior, and emit CI-friendly risk reports.

## Contributing

Run `make lint` and `make test` before proposing changes. Keep additions scoped to implemented roadmap work and avoid introducing infrastructure without a current requirement.
