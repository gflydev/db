package migrate

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"text/tabwriter"
)

// RunCLI parses args (everything after "db:migrate" on the command line), runs the requested
// operation against dialect's database, and returns a process exit code: 0 on success
// (including "nothing to do"), 1 on any error.
func RunCLI(args []string, dialect Dialect, dir string) int {
	return runCLI(args, dialect, dir, os.Stdout, os.Stderr)
}

func runCLI(args []string, dialect Dialect, dir string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("db:migrate", flag.ContinueOnError)
	fs.SetOutput(stderr)
	down := fs.Bool("down", false, "roll back instead of applying")
	all := fs.Bool("all", false, "apply/roll back every pending file (default: exactly one)")
	status := fs.Bool("status", false, "print every migration's state and exit")
	dryRun := fs.Bool("dry-run", false, "print what would run without executing it")
	baseline := fs.String("baseline", "", "mark every migration up to VERSION as already applied, without running SQL")
	force := fs.Bool("force", false, "overwrite an existing baseline, or proceed past a checksum mismatch")

	if err := fs.Parse(args); err != nil {
		return 1
	}
	if *baseline != "" && (*down || *all) {
		fmt.Fprintln(stderr, "migrate: --baseline cannot be combined with --down or --all")
		return 1
	}

	db, err := dialect.Open()
	if err != nil {
		fmt.Fprintf(stderr, "migrate: connecting to the database: %v\n", err)
		return 1
	}
	defer db.Close()

	m := NewMigrator(db, dialect, dir, "migrations")
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

func printStatus(w io.Writer, rows []StatusRow) {
	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
	fmt.Fprintln(tw, "VERSION\tNAME\tSTATUS\tBATCH\tRUNS\tCHECKSUM")
	for _, row := range rows {
		status, batch, runs, sum := "pending", "-", "0", "-"
		if row.Record != nil {
			status = row.Record.Status
			batch = fmt.Sprintf("%d", row.Record.Batch)
			runs = fmt.Sprintf("%d", row.Record.RunCount)
			if row.Record.Status == "up" {
				if row.ChecksumMatches {
					sum = "ok"
				} else {
					sum = "MISMATCH"
				}
			}
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\n", row.Migration.Version, row.Migration.Name, status, batch, runs, sum)
	}
	tw.Flush()
}
