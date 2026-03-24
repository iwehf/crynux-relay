## Migration Authoring Rules

### Versioning Contract

Migrations are versioned and executed once by the migration framework.

Migration code in this directory MUST NOT add defensive re-execution guards such as:

- `HasTable` checks before `CreateTable` or `DropTable`
- `HasColumn` checks before `AddColumn` or `DropColumn`
- `HasIndex` checks before `CreateIndex` or `DropIndex`

Write migration steps as direct, deterministic schema transitions for the target version.
