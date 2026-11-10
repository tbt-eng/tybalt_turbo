package helpcorpus

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"unicode/utf8"
)

// maxFileBytes keeps generated and data-like files out of the corpus.
const maxFileBytes = 300 * 1024

// Options says where the repo is.
type Options struct {
	// RepoDir holds ui/src and descriptions.
	RepoDir string
	// AppDir holds the Go app. Empty means RepoDir/app. The Docker build keeps
	// the app outside the folder that holds the other sources.
	AppDir string
}

// Stats describes one run, for the command's summary.
type Stats struct {
	Files int
	// Skipped lists the files that were too big or were not text.
	Skipped []string
	// Rules counts the rules in the corpus. NoRules and StaleRules list the
	// source files that should have rules and have none, or have rules for an
	// older source. Neither stops a build: the code is still in the corpus.
	Rules      int
	NoRules    []string
	StaleRules []string
}

// source is one tree of files to read.
type source struct {
	kind string
	// root is a repo path with forward slashes.
	root string
	exts []string
	// skipDirs are repo paths of folders to leave out.
	skipDirs []string
	// skipFiles are repo paths of files to leave out.
	skipFiles []string
	// skipSuffixes leave out files by the end of their name.
	skipSuffixes []string
}

const docsRoot = "ui/src/docs"

// sources lists what the corpus holds. A new folder inside a root is part of
// the corpus from the next build.
var sources = []source{
	{
		kind:      KindDocs,
		root:      docsRoot,
		exts:      []string{".md"},
		skipFiles: []string{docsRoot + "/README.md"},
	},
	{
		kind: KindNotes,
		root: "descriptions",
		exts: []string{".md"},
	},
	{
		kind: KindCode,
		root: "app",
		exts: []string{".go", ".sql"},
		skipDirs: []string{
			"app/cmd",
			"app/helpcorpus/data",
			"app/internal/testseed",
			"app/internal/testutils",
			// Migrations are history. Old rules in them would mislead.
			"app/migrations",
			"app/pb_data",
			"app/pb_public",
			"app/report_testing",
			"app/repros",
			"app/test_pb_data",
			"app/test_seed_data",
			"app/tmp",
		},
		// The build writes version.go, so it differs from one build to the next.
		skipFiles:    []string{"app/constants/version.go"},
		skipSuffixes: []string{"_test.go"},
	},
	{
		kind:     KindCode,
		root:     "ui/src",
		exts:     []string{".svelte", ".ts"},
		skipDirs: []string{docsRoot},
		skipFiles: []string{
			// Generated files.
			"ui/src/lib/pocketbase-types.ts",
			"ui/src/lib/version.ts",
		},
		skipSuffixes: []string{".test.ts", ".d.ts"},
	},
}

// skipDirNames are folders to leave out wherever they are.
var skipDirNames = []string{".git", ".svelte-kit", "node_modules", "testdata"}

// Generate reads the repo and builds the corpus. It reads nothing but the files,
// so the same files always give the same corpus.
func Generate(options Options) (*Corpus, *Stats, error) {
	stats := &Stats{}
	corpus := &Corpus{Format: FormatVersion, Chunks: []Chunk{}}
	var docs []file

	for _, src := range sources {
		files, err := readSource(options, src, stats)
		if err != nil {
			return nil, nil, err
		}
		if len(files) == 0 {
			return nil, nil, fmt.Errorf("found no files in %s", src.root)
		}
		for _, f := range files {
			for _, piece := range splitChunks(f.path, f.content) {
				corpus.Chunks = append(corpus.Chunks, Chunk{
					ID:       fmt.Sprintf("%s:%d", f.path, piece.start),
					Path:     f.path,
					Kind:     src.kind,
					Start:    piece.start,
					End:      piece.end,
					Heading:  piece.heading,
					Requires: requiresFor(f.path, piece.text),
					Text:     piece.text,
				})
			}
		}
		if src.kind == KindDocs {
			docs = files
		}
		for _, f := range files {
			if src.kind != KindCode || !HasRules(f.path) {
				continue
			}
			rules, stale, err := ruleChunks(options, f)
			if err != nil {
				return nil, nil, err
			}
			switch {
			case stale:
				stats.StaleRules = append(stats.StaleRules, f.path)
			case rules == nil:
				stats.NoRules = append(stats.NoRules, f.path)
			}
			stats.Rules += len(rules)
			corpus.Chunks = append(corpus.Chunks, rules...)
		}
	}

	handbook, err := buildHandbook(docs)
	if err != nil {
		return nil, nil, err
	}
	corpus.Handbook = handbook
	if corpus.Hash, err = hashOf(corpus.Handbook, corpus.Chunks); err != nil {
		return nil, nil, err
	}
	return corpus, stats, nil
}

type file struct {
	// path is a repo path with forward slashes.
	path    string
	content string
}

func readSource(options Options, src source, stats *Stats) ([]file, error) {
	dir := filepath.Join(options.RepoDir, filepath.FromSlash(src.root))
	if src.root == "app" && options.AppDir != "" {
		dir = options.AppDir
	}

	var files []file
	err := filepath.WalkDir(dir, func(name string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(dir, name)
		if err != nil {
			return err
		}
		repoPath := path.Join(src.root, filepath.ToSlash(rel))
		if entry.IsDir() {
			if slices.Contains(skipDirNames, entry.Name()) || slices.Contains(src.skipDirs, repoPath) {
				return filepath.SkipDir
			}
			return nil
		}
		if !entry.Type().IsRegular() || !wanted(src, repoPath) {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if info.Size() > maxFileBytes {
			stats.Skipped = append(stats.Skipped, repoPath+" (too big)")
			return nil
		}
		content, err := os.ReadFile(name)
		if err != nil {
			return err
		}
		if !utf8.Valid(content) {
			stats.Skipped = append(stats.Skipped, repoPath+" (not text)")
			return nil
		}
		files = append(files, file{path: repoPath, content: string(content)})
		return nil
	})
	if errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("%s is missing: %w", src.root, err)
	}
	if err != nil {
		return nil, err
	}
	sort.Slice(files, func(a, b int) bool { return files[a].path < files[b].path })
	stats.Files += len(files)
	return files, nil
}

func wanted(src source, repoPath string) bool {
	if !slices.Contains(src.exts, path.Ext(repoPath)) || slices.Contains(src.skipFiles, repoPath) {
		return false
	}
	for _, suffix := range src.skipSuffixes {
		if strings.HasSuffix(repoPath, suffix) {
			return false
		}
	}
	return true
}

// buildHandbook joins the user docs into the text that goes into every prompt.
// Start here comes first and What's new comes last. Each page sits in a tag
// that names its URL, so an answer can point to the page.
//
// The handbook is the same for every user. A page that is not for everyone, or
// that holds text the user docs must not hold, is an error and stops the build.
func buildHandbook(docs []file) (string, error) {
	rank := func(f file) int {
		switch path.Base(f.path) {
		case "start-here.md":
			return 0
		case "whats-new.md":
			return 2
		}
		return 1
	}
	ordered := slices.Clone(docs)
	sort.SliceStable(ordered, func(a, b int) bool { return rank(ordered[a]) < rank(ordered[b]) })

	var handbook strings.Builder
	for _, doc := range ordered {
		if len(requiresFor(doc.path, doc.content)) > 0 {
			return "", fmt.Errorf("%s is not for everyone, and the handbook goes to every user", doc.path)
		}
		if banned := handbookBans.FindString(doc.content); banned != "" {
			return "", fmt.Errorf("%s holds %q, and the user docs treat a commit as final", doc.path, banned)
		}
		slug := strings.TrimSuffix(path.Base(doc.path), ".md")
		fmt.Fprintf(&handbook, "<doc url=\"/docs/%s\">\n%s\n</doc>\n\n", slug, strings.TrimSpace(doc.content))
	}
	return strings.TrimSpace(handbook.String()) + "\n", nil
}

// CheckOutPath refuses a place for the corpus file that the server would serve
// to visitors. The server serves pb_public without a login.
func CheckOutPath(out string) error {
	abs, err := filepath.Abs(out)
	if err != nil {
		return err
	}
	for _, part := range strings.Split(filepath.ToSlash(abs), "/") {
		if part == "pb_public" {
			return fmt.Errorf("%s is inside pb_public, which the server serves without a login", out)
		}
	}
	return nil
}
