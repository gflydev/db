package migrate

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"text/tabwriter"
)

// helpText is printed by RunCLI's "help" argument, its -h/--help flag, and any flag-parsing
// error (so a typo lands on a full explanation, not just a one-line flag error). Keep it in
// sync with the flag descriptions below and with the package README's own command table.
const helpText = `db:migrate applies or rolls back the *.sql files in a migrations directory
(golang-migrate's NNNNNN_name.up.sql / .down.sql convention), tracking each file's state in a
"migrations" table it manages itself.

Usage:
  db:migrate [flags]

Flags:
  (none)                Apply the single next pending migration.
  --all                 Apply every pending migration (with --down: roll back the whole latest batch).
  --down                Roll back instead of applying.
  --status              Print every migration's state and exit.
  --dry-run             Print what would run without executing it.
  --baseline=VERSION    Mark every migration up to VERSION as already applied, without running SQL.
  --force               Overwrite an existing baseline, or proceed past a checksum mismatch.
  -h, --help            Show this help.

Examples:
  db:migrate                     Apply the next pending migration
  db:migrate --all               Apply every pending migration
  db:migrate --down              Roll back the most recently applied migration
  db:migrate --down --all        Roll back every migration in the latest batch
  db:migrate --status            Show every migration's state
  db:migrate --dry-run --all     Show which migrations --all would apply, without running them
  db:migrate --baseline=000023   Mark 000001..000023 as already applied, without running SQL
`

// RunCLI parses args (everything after "db:migrate" on the command line — see the package
// README for the full flag list), runs the requested operation against dialect's database, and
// returns a process exit code: 0 on success (including "nothing to do"), 1 on any error
// (printed to os.Stderr). A typical caller wires it in as:
//
//	case args[0] == "db:migrate":
//	    os.Exit(migrate.RunCLI(args[1:], postgres.Dialect{}, "database/migrations"))
func RunCLI(args []string, dialect Dialect, dir string) int {
	return runCLI(args, dialect, dir, os.Stdout, os.Stderr)
}

// runCLI is RunCLI with stdout/stderr as parameters instead of the real os.Stdout/os.Stderr, so
// tests can capture output without touching the process's real streams.
func runCLI(args []string, dialect Dialect, dir string, stdout, stderr io.Writer) int {
	// "help" has no leading dash, so the flag package below would otherwise treat it as a bare
	// positional argument and silently fall through to the default (apply) action instead of
	// explaining anything — handle it before fs.Parse ever sees it.
	if len(args) > 0 && args[0] == "help" {
		fmt.Fprint(stdout, helpText)
		return 0
	}

	fs := flag.NewFlagSet("db:migrate", flag.ContinueOnError)
	fs.SetOutput(stderr)                               // a real parse error (e.g. an unknown flag) prints its own line here
	fs.Usage = func() { fmt.Fprint(stdout, helpText) } // -h/--help itself always goes to stdout
	down := fs.Bool("down", false, "roll back instead of applying")
	all := fs.Bool("all", false, "apply/roll back every pending file (default: exactly one)")
	status := fs.Bool("status", false, "print every migration's state and exit")
	dryRun := fs.Bool("dry-run", false, "print what would run without executing it")
	baseline := fs.String("baseline", "", "mark every migration up to VERSION as already applied, without running SQL")
	force := fs.Bool("force", false, "overwrite an existing baseline, or proceed past a checksum mismatch")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 1
	}
	if *baseline != "" && (*down || *all) {
		fmt.Fprintf(stderr, "migrate: %v\n", ErrConflictingFlags)
		return 1
	}

	db, err := dialect.Open()
	if err != nil {
		fmt.Fprintf(stderr, "migrate: connecting to the database: %v\n", err)
		return 1
	}
	defer db.Close()

	m := NewMigrator(db, dialect, dir, DefaultTable)
	ctx := context.Background()

	switch {
	case *baseline != "":
		if err := m.Baseline(ctx, *baseline, *force); err != nil {
			fmt.Fprintf(stderr, "migrate: %v\n", err)
			return 1
		}
		fmt.Fprintf(stdout, "Baselined every migration up to and including %s.\n", *baseline)
		return 0

	case *status:
		rows, err := m.Status(ctx)
		if err != nil {
			fmt.Fprintf(stderr, "migrate: %v\n", err)
			return 1
		}
		printStatus(stdout, rows)
		return 0

	case *down:
		result, err := m.Down(ctx, *all, *dryRun)
		if err != nil {
			fmt.Fprintf(stderr, "migrate: %v\n", err)
			return 1
		}
		printResult(stdout, "roll back", "Rolled back", *dryRun, result)
		return 0

	default:
		result, err := m.Up(ctx, *all, *dryRun)
		if err != nil {
			fmt.Fprintf(stderr, "migrate: %v\n", err)
			return 1
		}
		printResult(stdout, "apply", "Applied", *dryRun, result)
		return 0
	}
}

// printResult writes one line per entry in result.Applied, phrased with pastTense ("Applied",
// "Rolled back") normally or infinitive ("apply", "roll back") prefixed with "Would " under
// dryRun. Writes a single "Nothing to do." line if result.Applied is empty.
func printResult(w io.Writer, infinitive, pastTense string, dryRun bool, result *Result) {
	if len(result.Applied) == 0 {
		fmt.Fprintln(w, "Nothing to do.")
		return
	}
	prefix := "Would " + infinitive
	if !dryRun {
		prefix = pastTense
	}
	for _, name := range result.Applied {
		fmt.Fprintf(w, "%s: %s\n", prefix, name)
	}
}

// printStatus renders rows as a tab-aligned table: one line per migration, showing its version,
// name, status (falling back to "pending" for a migration never applied), batch, run count, and
// whether its on-disk checksum still matches ("ok", "MISMATCH", "FILE MISSING", or "-" if never
// applied). A row with FileMissing prints a trailing note explaining what that means and why
// Down can't touch it, since "FILE MISSING" alone isn't self-explanatory in a status table.
func printStatus(w io.Writer, rows []StatusRow) {
	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
	fmt.Fprintln(tw, "VERSION\tNAME\tSTATUS\tBATCH\tRUNS\tCHECKSUM")
	hasFileMissing := false
	for _, row := range rows {
		status, batch, runs, sum := "pending", "-", "0", "-"
		if row.Record != nil {
			status = row.Record.Status
			batch = fmt.Sprintf("%d", row.Record.Batch)
			runs = fmt.Sprintf("%d", row.Record.RunCount)
			switch {
			case row.FileMissing:
				sum = "FILE MISSING"
				hasFileMissing = true
			case row.Record.Status == StatusUp:
				sum = "ok"
				if !row.ChecksumMatches {
					sum = "MISMATCH"
				}
			}
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\n", row.Migration.Version, row.Migration.Name, status, batch, runs, sum)
	}
	tw.Flush()
	if hasFileMissing {
		fmt.Fprintln(w, "\nFILE MISSING: recorded in the migrations table but its .sql files are no longer in the")
		fmt.Fprintln(w, "migrations directory (deleted or renamed after it ran). Restore the files to roll it back.")
	}
}
