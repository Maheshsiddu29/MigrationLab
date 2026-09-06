# Migration safety rules

MigrationLab rule identifiers use the `ML` namespace followed by a three-digit number. The namespace is reserved here so future diagnostics can remain stable across CLI and report formats.

No safety rules are implemented in Phase 1. The following initial rule set is planned:

- **ML001 — CREATE INDEX without CONCURRENTLY:** identify indexes that may block writes while they are built.
- **ML002 — DROP TABLE:** flag destructive table removal.
- **ML003 — DROP COLUMN:** flag destructive column removal.
- **ML004 — ALTER COLUMN TYPE:** identify type changes that may rewrite a table or break compatibility.
- **ML005 — SET NOT NULL:** identify constraint changes that may scan a table or reject existing data.
- **ML006 — expensive/default column addition:** identify column additions whose defaults may cause expensive work.
- **ML007 — foreign-key validation:** identify foreign-key operations that may scan data or hold disruptive locks.
- **ML008 — RENAME COLUMN:** flag compatibility risks from column renames.
- **ML009 — RENAME TABLE:** flag compatibility risks from table renames.
- **ML010 — TRUNCATE:** flag destructive data removal and its locking implications.
- **ML011 — DROP INDEX:** identify index removal that may affect live query performance or block activity.
- **ML012 — multiple high-lock DDL operations:** identify migrations that combine several operations with strong lock requirements.

The descriptions are planning notes, not implemented detection guarantees. Parsing, rule evaluation, severity levels, suppressions, and remediation guidance belong to later phases.
