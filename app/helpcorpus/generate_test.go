package helpcorpus

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// fixture is a small repo. The keys are repo paths.
var fixture = map[string]string{
	"ui/src/docs/start-here.md": "# Start here\n\nTurbo shows only the features that you have access to.\n",
	"ui/src/docs/time.md":       "# Time\n\nA time sheet holds your entries for one week.\n",
	"ui/src/docs/whats-new.md":  "# What's new\n\n## November 1 to 7, 2026\n\n- Help has a search bar.\n",
	"ui/src/docs/expenses.md":   "# Expenses\n\nAn expense records what you spent.\n",
	"ui/src/docs/README.md":     "# How to write the docs\n\nNot for users.\n",

	"descriptions/purchase_orders.md": "# Purchase orders\n\nThe authority matrix.\n",

	"app/main.go":                       "package main\n\nfunc main() {}\n",
	"app/hooks/expenses.go":             "package hooks\n\n// commit is final for staff.\nfunc commit() {}\n\n// uncommit reverses it.\nfunc uncommit() {}\n",
	"app/hooks/expenses_test.go":        "package hooks\n\nfunc TestCommit() {}\n",
	"app/routes/expenses.sql":           "SELECT id FROM expenses\n",
	"app/constants/version.go":          "package constants\n\nconst GitCommit = \"abc123\"\n",
	"app/migrations/1700000000_init.go": "package migrations\n",
	"app/cmd/tool/main.go":              "package main\n",
	"app/pb_public/index.html":          "<html></html>\n",
	"app/pb_data/notes.go":              "package data\n",
	"app/helpcorpus/data/corpus.json":   "{}\n",
	// The file of one PocketBase collection, which cmd/helpcollections writes.
	"app/helpcorpus/collections/profiles.md": "# Collection: profiles\n\n## Who can update a record (update rule): profiles\n\n```\nuid = @request.auth.id\n```\n",
	"app/.env":                               "ANTHROPIC_API_KEY=secret\n",
	"app/hooks/testdata/sample.go":           "package sample\n",

	"ui/src/routes/time/+page.svelte":   "<script lang=\"ts\">\n  let hours = $state(0);\n</script>\n\n<p>{hours}</p>\n",
	"ui/src/lib/utilities.ts":           "export function shortDate(value: string) {\n  return value;\n}\n",
	"ui/src/lib/utilities.test.ts":      "test();\n",
	"ui/src/lib/pocketbase-types.ts":    "export type Generated = string;\n",
	"ui/src/lib/version.ts":             "export const version = \"1\";\n",
	"ui/src/app.d.ts":                   "declare global {}\n",
	"ui/src/lib/logo.svg":               "<svg></svg>\n",
	"ui/node_modules/pkg/index.ts":      "export {};\n",
	"ui/src/node_modules/pkg/index.ts":  "export {};\n",
	"ui/tests/fixtures/docsHelp.svelte": "<p>fixture</p>\n",
}

func writeRepo(t *testing.T, files map[string]string, reverse bool) string {
	t.Helper()
	root := t.TempDir()
	paths := make([]string, 0, len(files))
	for path := range files {
		paths = append(paths, path)
	}
	// Map order is random, and reverse changes the order again. The corpus must
	// not depend on the order that files were made in.
	if reverse {
		for a, b := 0, len(paths)-1; a < b; a, b = a+1, b-1 {
			paths[a], paths[b] = paths[b], paths[a]
		}
	}
	for _, path := range paths {
		name := filepath.Join(root, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(name), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(name, []byte(files[path]), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func generate(t *testing.T, options Options) (*Corpus, []byte) {
	t.Helper()
	corpus, _, err := Generate(options)
	if err != nil {
		t.Fatal(err)
	}
	data, err := corpus.Encode()
	if err != nil {
		t.Fatal(err)
	}
	return corpus, data
}

func chunkPaths(corpus *Corpus) map[string]bool {
	paths := map[string]bool{}
	for _, chunk := range corpus.Chunks {
		paths[chunk.Path] = true
	}
	return paths
}

func TestGenerateReadsOnlyTheSources(t *testing.T) {
	corpus, _ := generate(t, Options{RepoDir: writeRepo(t, fixture, false)})
	want := []string{
		"ui/src/docs/expenses.md", "ui/src/docs/start-here.md", "ui/src/docs/time.md", "ui/src/docs/whats-new.md",
		"descriptions/purchase_orders.md",
		"app/hooks/expenses.go", "app/main.go", "app/routes/expenses.sql",
		// The collection files are Markdown, and they are a source of their own.
		"app/helpcorpus/collections/profiles.md",
		"ui/src/lib/utilities.ts", "ui/src/routes/time/+page.svelte",
	}
	got := chunkPaths(corpus)
	for _, path := range want {
		if !got[path] {
			t.Errorf("%s is missing from the corpus", path)
		}
		delete(got, path)
	}
	for path := range got {
		t.Errorf("%s must not be in the corpus", path)
	}
	for _, chunk := range corpus.Chunks {
		if strings.Contains(chunk.Text, "secret") {
			t.Errorf("%s holds text from an .env file", chunk.ID)
		}
	}
}

func TestGenerateTagsEachChunk(t *testing.T) {
	corpus, _ := generate(t, Options{RepoDir: writeRepo(t, fixture, false)})
	kinds := map[string]string{}
	for _, chunk := range corpus.Chunks {
		kinds[chunk.Path] = chunk.Kind
		if !chunk.ForEveryone() {
			t.Errorf("%s: the repo is public, so every chunk is for everyone, got %v", chunk.ID, chunk.Requires)
		}
		if chunk.Requires == nil {
			t.Errorf("%s: requires must be a list, so the file always holds one", chunk.ID)
		}
		if chunk.ID != fmt.Sprintf("%s:%d", chunk.Path, chunk.Start) || chunk.End < chunk.Start {
			t.Errorf("%s has a bad id or line range", chunk.ID)
		}
	}
	for path, kind := range map[string]string{
		"ui/src/docs/time.md":             KindDocs,
		"descriptions/purchase_orders.md": KindNotes,
		"app/hooks/expenses.go":           KindCode,
		"ui/src/lib/utilities.ts":         KindCode,
	} {
		if kinds[path] != kind {
			t.Errorf("%s has kind %q, want %q", path, kinds[path], kind)
		}
	}
}

func TestSameFilesGiveTheSameBytes(t *testing.T) {
	first, firstBytes := generate(t, Options{RepoDir: writeRepo(t, fixture, false)})
	_, secondBytes := generate(t, Options{RepoDir: writeRepo(t, fixture, true)})
	if !bytes.Equal(firstBytes, secondBytes) {
		t.Fatal("two copies of the same files gave different bytes")
	}
	if strings.Contains(string(firstBytes), os.TempDir()) {
		t.Fatal("the corpus holds a path from this machine")
	}

	changed := map[string]string{}
	for path, content := range fixture {
		changed[path] = content
	}
	changed["ui/src/lib/utilities.ts"] += "// one more line\n"
	third, thirdBytes := generate(t, Options{RepoDir: writeRepo(t, changed, false)})
	if bytes.Equal(firstBytes, thirdBytes) || first.Hash == third.Hash {
		t.Fatal("a changed file must change the bytes and the hash")
	}
}

func TestAppDirCanBeOutsideTheRepoDir(t *testing.T) {
	root := writeRepo(t, fixture, false)
	_, want := generate(t, Options{RepoDir: root})

	moved := filepath.Join(t.TempDir(), "elsewhere")
	if err := os.Rename(filepath.Join(root, "app"), moved); err != nil {
		t.Fatal(err)
	}
	_, got := generate(t, Options{RepoDir: root, AppDir: moved})
	if !bytes.Equal(got, want) {
		t.Fatal("the Docker layout gave a different corpus")
	}
}

func TestEncodeAndDecode(t *testing.T) {
	corpus, data := generate(t, Options{RepoDir: writeRepo(t, fixture, false)})
	if lines := bytes.Count(data, []byte("\n")); lines != len(corpus.Chunks)+2 {
		t.Fatalf("the file has %d lines for %d chunks", lines, len(corpus.Chunks))
	}
	if bytes.Contains(data, []byte("\\"+"u003c")) {
		t.Fatal("the file escapes <, which makes code hard to read")
	}
	decoded, err := Decode(data)
	if err != nil {
		t.Fatal(err)
	}
	again, err := decoded.Encode()
	if err != nil || !bytes.Equal(again, data) {
		t.Fatal("a decoded corpus must encode to the same bytes")
	}

	tampered := bytes.Replace(data, []byte("A time sheet holds"), []byte("A time sheet hides"), 1)
	if _, err := Decode(tampered); err == nil || !strings.Contains(err.Error(), "hash") {
		t.Fatalf("a changed file must fail its hash, got %v", err)
	}
	otherFormat := bytes.Replace(data, []byte(`{"format":1,`), []byte(`{"format":99,`), 1)
	if _, err := Decode(otherFormat); err == nil || !strings.Contains(err.Error(), "format") {
		t.Fatalf("another format must fail, got %v", err)
	}
}

func TestHandbookHoldsTheUserDocsInOrder(t *testing.T) {
	corpus, _ := generate(t, Options{RepoDir: writeRepo(t, fixture, false)})
	order := []string{`<doc url="/docs/start-here">`, `<doc url="/docs/expenses">`, `<doc url="/docs/time">`, `<doc url="/docs/whats-new">`}
	position := -1
	for _, tag := range order {
		next := strings.Index(corpus.Handbook, tag)
		if next <= position {
			t.Fatalf("%s is missing or out of order", tag)
		}
		position = next
	}
	for _, absent := range []string{"How to write the docs", "authority matrix", "package hooks"} {
		if strings.Contains(corpus.Handbook, absent) {
			t.Errorf("the handbook holds %q, which is not in the user docs", absent)
		}
	}
	if !strings.Contains(corpus.Handbook, "A time sheet holds your entries for one week.\n</doc>") {
		t.Error("the handbook must hold each page's text")
	}
}

func TestARuleCanRestrictAFolderOrAWord(t *testing.T) {
	useRules(t,
		[]pathRule{
			{Prefix: "ui/src/", Requires: everyone},
			{Prefix: "app/", Requires: everyone},
			{Prefix: "descriptions/", Requires: []string{"report"}},
		},
		[]termRule{{Term: "uncommit", Requires: []string{"commit", "admin"}}},
	)
	corpus, _ := generate(t, Options{RepoDir: writeRepo(t, fixture, false)})
	staff := map[string]bool{"time": true}
	committer := map[string]bool{"time": true, "commit": true}
	sawWord := false
	for _, chunk := range corpus.Chunks {
		hasWord := strings.Contains(strings.ToLower(chunk.Text), "uncommit")
		sawWord = sawWord || hasWord
		switch {
		case chunk.Kind == KindNotes:
			if chunk.Allows(committer) || !chunk.Allows(map[string]bool{"report": true}) {
				t.Errorf("%s: only the report claim opens the notes, got %v", chunk.ID, chunk.Requires)
			}
		case hasWord:
			if chunk.Allows(staff) || !chunk.Allows(committer) {
				t.Errorf("%s: only commit or admin opens the word, got %v", chunk.ID, chunk.Requires)
			}
		case !chunk.ForEveryone():
			t.Errorf("%s must be for everyone, got %v", chunk.ID, chunk.Requires)
		}
	}
	if !sawWord {
		t.Fatal("the fixture must hold the restricted word")
	}
}

func TestTheHandbookMustBeForEveryone(t *testing.T) {
	useRules(t, []pathRule{{Prefix: "ui/src/docs/", Requires: []string{"admin"}}}, nil)
	_, _, err := Generate(Options{RepoDir: writeRepo(t, fixture, false)})
	if err == nil || !strings.Contains(err.Error(), "is not for everyone") {
		t.Fatalf("got %v", err)
	}
}

func TestUncommitInTheUserDocsStopsTheBuild(t *testing.T) {
	files := map[string]string{}
	for path, content := range fixture {
		files[path] = content
	}
	files["ui/src/docs/expenses.md"] += "\nAn uncommitted expense is not final.\n"
	if _, _, err := Generate(Options{RepoDir: writeRepo(t, files, false)}); err != nil {
		t.Fatalf("\"uncommitted\" is fine in the user docs, got %v", err)
	}
	files["ui/src/docs/expenses.md"] += "\nAsk a committer to Uncommit it.\n"
	_, _, err := Generate(Options{RepoDir: writeRepo(t, files, false)})
	if err == nil || !strings.Contains(err.Error(), "ui/src/docs/expenses.md") || !strings.Contains(err.Error(), `"Uncommit"`) {
		t.Fatalf("got %v", err)
	}
}

func TestAMissingSourceStopsTheBuild(t *testing.T) {
	files := map[string]string{}
	for path, content := range fixture {
		if !strings.HasPrefix(path, "descriptions/") {
			files[path] = content
		}
	}
	if _, _, err := Generate(Options{RepoDir: writeRepo(t, files, false)}); err == nil || !strings.Contains(err.Error(), "descriptions") {
		t.Fatalf("got %v", err)
	}
}

func TestBigAndBinaryFilesAreSkippedAndReported(t *testing.T) {
	files := map[string]string{}
	for path, content := range fixture {
		files[path] = content
	}
	files["app/routes/huge.sql"] = strings.Repeat("SELECT 1;\n", maxFileBytes/10+1)
	files["app/routes/binary.go"] = "package routes\n\xff\xfe\n"
	corpus, stats, err := Generate(Options{RepoDir: writeRepo(t, files, false)})
	if err != nil {
		t.Fatal(err)
	}
	paths := chunkPaths(corpus)
	if paths["app/routes/huge.sql"] || paths["app/routes/binary.go"] {
		t.Fatal("a big or binary file is in the corpus")
	}
	skipped := strings.Join(stats.Skipped, "\n")
	if !strings.Contains(skipped, "app/routes/huge.sql (too big)") || !strings.Contains(skipped, "app/routes/binary.go (not text)") {
		t.Fatalf("skipped = %q", skipped)
	}
}

func TestTheCorpusFileCannotGoInThePublicFolder(t *testing.T) {
	for _, out := range []string{"pb_public/corpus.json", "../app/pb_public/help/corpus.json", "/srv/app/pb_public/corpus.json"} {
		if err := CheckOutPath(out); err == nil {
			t.Errorf("%s is in the public folder, and it was accepted", out)
		}
	}
	for _, out := range []string{"helpcorpus/data/corpus.json", "/tmp/pb_public_copy/corpus.json"} {
		if err := CheckOutPath(out); err != nil {
			t.Errorf("%s: %v", out, err)
		}
	}
}

// The real repo is two folders up from this package.
func TestTheRealRepo(t *testing.T) {
	repo := filepath.Join("..", "..")
	if _, err := os.Stat(filepath.Join(repo, "ui", "src", "docs")); err != nil {
		t.Skip("the repo's UI sources are not here")
	}
	corpus, first := generate(t, Options{RepoDir: repo})
	_, second := generate(t, Options{RepoDir: repo})
	if !bytes.Equal(first, second) {
		t.Fatal("two runs on the repo gave different bytes")
	}
	if _, err := Decode(first); err != nil {
		t.Fatal(err)
	}
	if len(corpus.Chunks) < 1000 {
		t.Fatalf("the repo gave only %d chunks", len(corpus.Chunks))
	}

	ids := map[string]bool{}
	docs := 0
	for _, chunk := range corpus.Chunks {
		if ids[chunk.ID] {
			t.Fatalf("%s is in the corpus twice", chunk.ID)
		}
		ids[chunk.ID] = true
		if chunk.Kind == KindDocs {
			docs++
		}
		if !chunk.ForEveryone() {
			t.Errorf("%s is not for everyone: %v", chunk.ID, chunk.Requires)
		}
		for _, banned := range []string{"/migrations/", "/pb_public/", "/pb_data/", "_test.go", "pocketbase-types.ts", "/node_modules/", "app/cmd/"} {
			if strings.Contains(chunk.Path, banned) {
				t.Errorf("%s must not be in the corpus", chunk.Path)
			}
		}
		if len(chunk.Text) > chunkMax+chunkMin && strings.Count(chunk.Text, "\n") > 0 {
			t.Errorf("%s has %d bytes", chunk.ID, len(chunk.Text))
		}
	}
	if docs == 0 {
		t.Fatal("the corpus has no user docs")
	}
	if banned := handbookBans.FindString(corpus.Handbook); banned != "" {
		t.Errorf("the handbook holds %q", banned)
	}
	// One prompt must stay well under the size where the price goes up.
	if len(corpus.Handbook) > 150_000 {
		t.Errorf("the handbook has %d bytes, which is too much for every prompt", len(corpus.Handbook))
	}
}
