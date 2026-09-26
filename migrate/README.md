# gFly DB Migrate

    Copyright © 2026, gFly
    https://www.gfly.dev
    All rights reserved.

Runs versioned `*.sql` migration files (golang-migrate's `NNNNNN_name.up.sql` / `.down.sql`
naming convention) against PostgreSQL or MySQL, tracking each file's state — applied/rolled
back, batch, run count, checksum — in a `migrations` table it manages itself.

## Install

```bash
go get -u github.com/gflydev/db/migrate@latest

# PostgreSQL
go get -u github.com/gflydev/db/migrate/postgres@latest
# MySQL
go get -u github.com/gflydev/db/migrate/mysql@latest
```

## Usage

Wire one `case` into your app's CLI entrypoint:

```go
import (
	"os"

	dbmigrate "github.com/gflydev/db/migrate"
	dbmigratePg "github.com/gflydev/db/migrate/postgres"
)

func main() {
	args := os.Args[1:]
	switch {
	case len(args) > 0 && args[0] == "db:migrate":
		os.Exit(dbmigrate.RunCLI(args[1:], dbmigratePg.Dialect{}, "database/migrations"))
	}
}
```

Connection settings are read from `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USERNAME`,
`DB_PASSWORD`, `DB_SSL_MODE` (the same variables `github.com/gflydev/db/psql` and
`github.com/gflydev/db/mysql` already use).

```bash
./build/artisan db:migrate                    # apply the single next pending migration
./build/artisan db:migrate --all              # apply every pending migration (one batch)
./build/artisan db:migrate --down             # roll back the single most recent migration
./build/artisan db:migrate --down --all       # roll back every migration in the latest batch
./build/artisan db:migrate --status           # show every migration's state
./build/artisan db:migrate --dry-run          # combine with the above: print, don't execute
./build/artisan db:migrate --baseline=000023  # mark files up to 000023 as already applied
./build/artisan db:migrate help               # same as -h/--help: print the flag list and examples above
```

## Design notes

- Each migration file runs in its own transaction on PostgreSQL; MySQL's DDL auto-commits, so a
  mid-file failure there cannot be rolled back and is reported as such.
- One advisory lock (`pg_advisory_lock` / `GET_LOCK`) is held for a whole invocation, so two
  concurrent `db:migrate` runs serialize rather than race.
- Before any command runs, every migration currently recorded as `up` has its on-disk checksum
  re-verified against what was recorded; a mismatch aborts (`--force` to proceed anyway).
- `--baseline` is a one-time bootstrap for a database that already has this schema outside the
  tool's tracking — it never executes SQL, and refuses to run on a non-empty tracking table
  without `--force`.
- Rolling back a migration keeps its tracking row (status flips to `down`, `run_count` and
  `batch` stay as history) rather than deleting it, so `run_count` survives repeated
  rollback/reapply cycles.

## Testing

```bash
go test ./...                     # this package: no database required
go test ./postgres/...            # SQL-string unit tests: no database required
go test -tags=integration ./postgres/...  # full up/down/status cycle against a real Postgres
go test ./mysql/...                # SQL-string unit tests: no database required
go test -tags=integration ./mysql/...      # full up/down/status cycle against a real MySQL
```

The `integration` build tag needs `DB_*` environment variables pointed at a disposable database
— it creates and drops its own `migrations`/`widgets` tables and skips (not fails) if it can't
reach a database.
