package testseed

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSweepStaleTempDirs(t *testing.T) {
	root := t.TempDir()
	now := time.Now()
	old := now.Add(-staleTempDirAge - time.Minute)
	recent := now.Add(-staleTempDirAge + time.Minute)

	mkdir := func(name string, modTime time.Time) {
		t.Helper()
		path := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Join(path, "storage"), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.Chtimes(path, modTime, modTime); err != nil {
			t.Fatal(err)
		}
	}

	mkdir("pb_test_old", old)
	mkdir(templateDirPrefix+"old", old)
	mkdir("pb_test_recent", recent)
	mkdir(templateDirPrefix+"recent", recent)
	mkdir("other_old", old)

	oldFile := filepath.Join(root, "pb_test_file")
	if err := os.WriteFile(oldFile, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(oldFile, old, old); err != nil {
		t.Fatal(err)
	}

	SweepStaleTempDirs(root, now)

	for name, wantKept := range map[string]bool{
		"pb_test_old":                false,
		templateDirPrefix + "old":    false,
		"pb_test_recent":             true,
		templateDirPrefix + "recent": true,
		"other_old":                  true,
		"pb_test_file":               true,
	} {
		_, err := os.Lstat(filepath.Join(root, name))
		if kept := err == nil; kept != wantKept {
			t.Errorf("%s: kept=%v, want %v", name, kept, wantKept)
		}
	}
}
