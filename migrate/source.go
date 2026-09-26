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

var filenamePattern = regexp.MustCompile(`^(\d{6})_([a-z0-9_]+)\.(up|down)\.sql$`)

// Load reads every *.sql file in dir, validates the naming and up/down pairing rules, and
// returns the migrations sorted ascending by version.
func Load(dir string) ([]Migration, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("migrate: reading %s: %w", dir, err)
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
		return nil, fmt.Errorf("migrate: filenames must match NNNNNN_name.(up|down).sql: %s", strings.Join(badNames, ", "))
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
		return nil, fmt.Errorf("migrate: unpaired migration file(s): %s", strings.Join(orphans, ", "))
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

// checksum returns the hex-encoded sha256 of path's contents.
func checksum(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("migrate: checksumming %s: %w", path, err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("migrate: checksumming %s: %w", path, err)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
