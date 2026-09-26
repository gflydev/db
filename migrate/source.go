package migrate

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/gflydev/core/errors"
	"github.com/gflydev/core/utils"
)

// filenamePattern matches a single migration file and captures its three parts: the timestamp
// (timestampLayout's shape), the snake_case description, and its direction (up/down). A file
// that doesn't match this — including a bare "*.sql" with no direction suffix, or one still
// using the old NNNNNN sequential-number convention — is rejected by Load rather than silently
// skipped, since an unrecognized file can't be placed in the ascending/descending order the rest
// of this package depends on.
var filenamePattern = regexp.MustCompile(`^(\d{8}_\d{6})_([a-z0-9_]+)\.(up|down)\.sql$`)

// descriptionPattern is the same charset filenamePattern requires for a description, used to
// validate a name before New creates files from it.
var descriptionPattern = regexp.MustCompile(`^[a-z0-9_]+$`)

// Load reads every *.sql file directly inside dir (not recursively), validates that each one
// matches "YYYYMMDD_HHMMSS_description.(up|down).sql" with a real, parseable timestamp, and that
// every "up" file has a matching "down" file and vice versa, and returns the migrations sorted
// ascending by version (equivalently, by timestamp, equivalently, by creation order).
//
// Returns an error listing every offending filename if any file fails the naming pattern, has an
// invalid timestamp (e.g. month 13), or is missing its up/down counterpart — Load fails
// all-or-nothing rather than returning a partial, silently-incomplete set.
//
// Load lists dir with the standard library's os.ReadDir rather than gflydev/storage: the
// migrations directory is part of the deployed source tree (like the framework's own view
// templates), not a caller-configurable storage disk, and IStorage has no directory-listing
// method to begin with.
func Load(dir string) ([]Migration, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, errors.New("reading migrations directory %s: %w", dir, err)
	}

	ups := map[string]string{}   // name -> path
	downs := map[string]string{} // name -> path
	var badNames []string

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		filename := entry.Name()
		if !strings.HasSuffix(filename, ".sql") {
			continue
		}
		match := filenamePattern.FindStringSubmatch(filename)
		if match == nil {
			badNames = append(badNames, filename)
			continue
		}
		version, desc, direction := match[1], match[2], match[3]
		if _, err := time.Parse(timestampLayout, version); err != nil {
			badNames = append(badNames, filename)
			continue
		}
		name := version + "_" + desc
		path := filepath.Join(dir, filename)
		if direction == "up" {
			ups[name] = path
		} else {
			downs[name] = path
		}
	}

	if len(badNames) > 0 {
		sort.Strings(badNames)
		return nil, errors.New("filenames must match YYYYMMDD_HHMMSS_name.(up|down).sql with a valid UTC timestamp: %s", strings.Join(badNames, ", "))
	}

	var orphans []string
	for name := range ups {
		if _, ok := downs[name]; !ok {
			orphans = append(orphans, name+".up.sql (missing "+name+".down.sql)")
		}
	}
	for name := range downs {
		if _, ok := ups[name]; !ok {
			orphans = append(orphans, name+".down.sql (missing "+name+".up.sql)")
		}
	}
	if len(orphans) > 0 {
		sort.Strings(orphans)
		return nil, errors.New("unpaired migration file(s): %s", strings.Join(orphans, ", "))
	}

	migrations := make([]Migration, 0, len(ups))
	for name, upPath := range ups {
		migrations = append(migrations, Migration{
			Version:  name[:versionLength],
			Name:     name,
			UpPath:   upPath,
			DownPath: downs[name],
		})
	}
	sort.Slice(migrations, func(i, j int) bool { return migrations[i].Name < migrations[j].Name })

	return migrations, nil
}

// checksum returns the hex-encoded sha256 of path's contents. Migrator uses it to detect when
// an applied migration's .up.sql file has changed since it ran (ErrChecksumMismatch).
//
// It reads the file through gflydev/core/utils.ReadFileAsString rather than crypto/sha256's
// usual io.Copy-from-an-open-file idiom: migration files are always small (SQL scripts, not
// data dumps), so the simpler one-shot read costs nothing and keeps this package's file access
// on the same helper Migrator.runFile uses.
func checksum(path string) (string, error) {
	contents, err := utils.ReadFileAsString(path)
	if err != nil {
		return "", errors.New("checksumming %s: %w", path, err)
	}
	sum := sha256.Sum256([]byte(contents))
	return hex.EncodeToString(sum[:]), nil
}

// New creates a fresh migration pair in dir, named "{current UTC timestamp}_{description}", and
// returns it as a Migration (with UpPath/DownPath already pointing at the new files). Both files
// are created with a one-line placeholder comment and nothing else — valid, no-op SQL until
// filled in.
//
// description must already be snake_case ([a-z0-9_]+); RunCLI's --new flag validates this before
// calling New, and New itself re-validates so a direct library caller gets the same guarantee.
// Returns an error without creating anything if description is invalid, or if either target
// file already exists (only possible if New is called twice for the same description within the
// same second).
func New(dir, description string) (Migration, error) {
	if !descriptionPattern.MatchString(description) {
		return Migration{}, errors.New("migration description %q must be snake_case (lowercase letters, digits, underscore)", description)
	}

	version := time.Now().UTC().Format(timestampLayout)
	name := version + "_" + description
	upPath := filepath.Join(dir, name+".up.sql")
	downPath := filepath.Join(dir, name+".down.sql")

	for _, path := range []string{upPath, downPath} {
		if utils.FileExists(path) {
			return Migration{}, errors.New("%s already exists", path)
		}
	}

	if err := utils.WriteStringToFile(upPath, "-- write the \"up\" SQL for "+description+" here\n", 0o644); err != nil {
		return Migration{}, errors.New("creating %s: %w", upPath, err)
	}
	if err := utils.WriteStringToFile(downPath, "-- write the \"down\" SQL for "+description+" here\n", 0o644); err != nil {
		return Migration{}, errors.New("creating %s: %w", downPath, err)
	}

	return Migration{Version: version, Name: name, UpPath: upPath, DownPath: downPath}, nil
}
