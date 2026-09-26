package migrate

import "errors"

// Sentinel errors returned by Migrator. Wrap these with fmt.Errorf("%w: ...", ErrX) rather than
// building a new, unmatched error string for the same condition — callers (including this
// package's own CLI and tests) must be able to distinguish them with errors.Is instead of
// parsing error text.
var (
	// ErrChecksumMismatch means a migration recorded as applied ("up") no longer matches the
	// sha256 checksum stored when it ran — its .up.sql file changed on disk since then. Up,
	// Down and Baseline all refuse to proceed past this unless force is true.
	ErrChecksumMismatch = errors.New("checksum mismatch: an applied migration changed on disk since it ran")

	// ErrNonEmptyTable means Baseline was called with force == false against a tracking table
	// that already has at least one row. Baseline is a one-time bootstrap, not a merge tool.
	ErrNonEmptyTable = errors.New("baseline refuses to run on a non-empty tracking table")

	// ErrUnknownVersion means the version passed to Baseline does not match any migration
	// discovered by Load.
	ErrUnknownVersion = errors.New("baseline version does not match any known migration")

	// ErrFilesMissing means the tracking table has a row for a migration whose .sql files are
	// no longer present in the migrations directory, so Down has nothing to execute for it.
	ErrFilesMissing = errors.New("migration is recorded as applied but its files are missing from disk")

	// ErrConflictingFlags means RunCLI was given --baseline together with --down or --all,
	// which the CLI does not support (baseline never applies or rolls back a range of files).
	ErrConflictingFlags = errors.New("--baseline cannot be combined with --down or --all")
)
