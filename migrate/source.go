package migrate

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// filenamePattern matches a single migration file and captures its three parts: the 6-digit
// version, the snake_case description, and its direction (up/down). A file that doesn't match
// this — including a bare "*.sql" with no direction suffix — is rejected by Load rather than
// silently skipped, since an unrecognized file can't be placed in the ascending/descending order
// the rest of this package depends on.
var filenamePattern = regexp.MustCompile(`^(\d{6})_([a-z0-9_]+)\.(up|down)\.sql$`)

// Load reads every *.sql file directly inside dir (not recursively), validates that each one
// matches "NNNNNN_description.(up|down).sql" and that every "up" file has a matching "down"
// file and vice versa, and returns the migrations sorted ascending by version.
//
// Returns an error listing every offending filename if any file fails the naming pattern or is
// missing its up/down counterpart — Load fails all-or-nothing rather than returning a partial,
// silently-incomplete set.
func Load(dir string) ([]Migration, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading migrations directory %s: %w", dir, err)
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
		return nil, fmt.Errorf("filenames must match NNNNNN_name.(up|down).sql: %s", strings.Join(badNames, ", "))
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
		return nil, fmt.Errorf("unpaired migration file(s): %s", strings.Join(orphans, ", "))
	}

	migrations := make([]Migration, 0, len(ups))
	for name, upPath := range ups {
		migrations = append(migrations, Migration{
			Version:  name[:6],
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
func checksum(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("checksumming %s: %w", path, err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("checksumming %s: %w", path, err)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
