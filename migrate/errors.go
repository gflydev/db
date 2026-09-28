package migrate

import "github.com/gflydev/core/errors"

// Sentinel errors returned by Migrator. Wrap these with errors.New("%w: ...", ErrX) (this
// package's convention — see Migrator's own error sites — is github.com/gflydev/core/errors,
// not the standard library's, matching the rest of the gFly framework) rather than building a
// new, unmatched error string for the same condition. Callers, including this package's own CLI
// and tests, must be able to distinguish them with errors.Is instead of parsing error text.
//
// These stay distinct from github.com/gflydev/core/errors' own sentinels (ItemNotFound,
// InvalidParameter, ...): those are shaped for mapping an HTTP request to a status code, and
// reusing one here would make an unrelated part of a gFly app's API layer match on a condition
// that has nothing to do with it.
var (
	// ErrChecksumMismatch means a migration recorded as applied ("up") no longer matches the
	// sha256 checksum stored when it ran — its .up.sql file changed on disk since then. Up,
	// Down and Baseline all refuse to proceed past this unless force is true.
	ErrChecksumMismatch = errors.New("checksum mismatch: an applied migration changed on disk since it ran")

	// ErrNonEmptyTable means Baseline was called with force == false against a tracking table
	// that has a row this tool applied or rolled back itself. A table holding only an earlier
	// baseline is not refused: Baseline extends it (see ErrBaselineNotAhead).
	ErrNonEmptyTable = errors.New("baseline refuses to run on a non-empty tracking table")

	// ErrBaselineNotAhead means Baseline was called with force == false to extend an existing
	// baseline, but version is not later than the baseline's last migration. A baseline only
	// moves forward; shrinking it would need rows removed, which Baseline never does.
	ErrBaselineNotAhead = errors.New("baseline version must be later than the current baseline")

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
