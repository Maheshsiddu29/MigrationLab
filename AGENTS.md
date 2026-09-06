# MigrationLab contributor guidance

## Scope

MigrationLab is a PostgreSQL migration safety and simulation engine under active development. Keep implementation aligned with the current roadmap phase; do not scaffold unused architecture layers.

## Engineering conventions

- Use idiomatic Go and keep packages focused.
- Format Go files with `gofmt` and run `go vet ./...` and `go test ./...` before handing off changes.
- Keep file discovery and other deterministic operations stable across runs.
- Add useful context to returned errors.
- Prefer standard-library and standard Go tooling unless a dependency provides clear current value.
- Do not make tests depend on a running PostgreSQL instance unless they are explicitly integration tests.

## Repository safety

- Do not commit, push, merge, publish releases, create remotes, or rewrite Git history unless the user explicitly requests it.
- Preserve unrelated user changes.
- Do not claim planned migration analysis or simulation capabilities as implemented.
