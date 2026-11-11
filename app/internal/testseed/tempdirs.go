package testseed

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// templateDirPrefix names the seeded template directory that TemplateDir
// creates in os.TempDir().
const templateDirPrefix = "tybalt-testseed-template-"

// staleTempDirAge is how old a leftover test directory must be before
// SweepStaleTempDirs removes it. A test binary is killed by go test's default
// timeout long before this, so a directory this old belongs to a run that
// crashed, was interrupted, or timed out, never to one still running in
// another package or clone.
const staleTempDirAge = 6 * time.Hour

// staleTempDirPrefixes are the directories that test runs leave in
// os.TempDir() when they stop before their cleanup runs: PocketBase's TestApp
// clones (tests.TempDirClone) and this package's seeded templates.
var staleTempDirPrefixes = []string{"pb_test_", templateDirPrefix}

var (
	mainActive atomic.Bool
	sweepOnce  sync.Once
)

// Main runs a package's tests and then removes the seeded template that
// NewSeededTestApp clones. Every package whose tests use NewSeededTestApp (or
// testutils.SetupTestApp) calls it from TestMain:
//
//	func TestMain(m *testing.M) { testseed.Main(m) }
func Main(m *testing.M) {
	mainActive.Store(true)
	sweepOnce.Do(func() { SweepStaleTempDirs(os.TempDir(), time.Now()) })

	code := m.Run()

	if templateDir != "" {
		os.RemoveAll(templateDir)
	}
	os.Exit(code)
}

// SweepStaleTempDirs removes test directories in root that a crashed or
// interrupted run left behind. It removes only directories whose names start
// with one of staleTempDirPrefixes and that were last modified more than
// staleTempDirAge before now. Errors are ignored: another test process may be
// sweeping the same directory at the same time.
func SweepStaleTempDirs(root string, now time.Time) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() || !hasStalePrefix(entry.Name()) {
			continue
		}
		info, err := entry.Info()
		if err != nil || now.Sub(info.ModTime()) < staleTempDirAge {
			continue
		}
		os.RemoveAll(filepath.Join(root, entry.Name()))
	}
}

func hasStalePrefix(name string) bool {
	for _, prefix := range staleTempDirPrefixes {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}
	return false
}
