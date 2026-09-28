# gFly DB Migrate

    Copyright © 2026, gFly
    https://www.gfly.dev
    All rights reserved.

Runs versioned `*.sql` migration files — named `YYYYMMDD_HHMMSS_description.up.sql` /
`.down.sql`, a UTC timestamp to the second — against PostgreSQL or MySQL, tracking each file's
state — applied/rolled back, batch, run count, checksum — in a `migrations` table it manages
itself.

The timestamp (not a small sequential counter, and not configurable to anything else) is
deliberate: when two people on separate branches each pick "the next number," they either
collide or, worse, silently don't collide but sort in an order neither of them intended once
merged. A per-second timestamp makes that collision astronomically unlikely and keeps the merged
order matching the order each file was actually written in. `--new` (below) generates one.

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
./build/artisan db:migrate                              # apply the single next pending migration
./build/artisan db:migrate --all                        # apply every pending migration (one batch)
./build/artisan db:migrate --down                       # roll back the single most recent migration
./build/artisan db:migrate --down --all                 # roll back every migration in the latest batch
./build/artisan db:migrate --status                     # show every migration's state
./build/artisan db:migrate --dry-run                    # combine with the above: print, don't execute
./build/artisan db:migrate --new=create_widgets_table   # create a new migration pair, timestamped now
./build/artisan db:migrate --baseline=20260101_000000   # mark files up to that timestamp as already applied
./build/artisan db:migrate help                         # same as -h/--help: print the flag list and examples above
```

### Files created by `--new`

The new `.up.sql` and `.down.sql` files are not empty: they hold guidance and numbered section
headers, written only as SQL comments, so an untouched pair still runs as a no-op.

- `.up.sql` — rules to keep in mind (one concern per migration, transaction behaviour per
  dialect, never edit the file once applied), a checklist of the objects it creates, and
  sections in apply order: 1. types/extensions/sequences, 2. tables, 3. changes to existing
  tables, 4. constraints and indexes, 5. data, 6. functions/triggers/views.
- `.down.sql` — the same sections in reverse order (6 → 1), a table of commonly forgotten
  removals (`CREATE TYPE` → `DROP TYPE`, and so on), and a reminder to run up → down → up on a
  scratch database before committing.

Delete the sections a migration does not need.

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
- `--status` also reports any tracking-table row whose `.sql` files are no longer in the
  migrations directory (deleted or renamed after it ran) as `FILE MISSING` — without this check
  such a row would simply vanish from every command's view, staying in the table forever with no
  way to roll it back until the files are restored.
- `Load` accepts exactly one filename shape — `YYYYMMDD_HHMMSS_description.(up|down).sql` with a
  real, parseable UTC timestamp — and refuses the whole directory (not just the offending file)
  if anything else is in it, including golang-migrate's older `NNNNNN_description.up.sql`
  sequential-number style. A migrations directory adopting this module for the first time must
  rename its existing files to the timestamp style before `db:migrate` will run at all.

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
