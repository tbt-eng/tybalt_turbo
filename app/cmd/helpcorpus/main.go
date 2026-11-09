// Command helpcorpus builds the corpus that Help chat searches.
//
// It reads the user docs, the design notes and the app's code from the repo,
// cuts them into chunks, tags each chunk with who may read it, and writes one
// JSON file. The same files always give the same bytes.
//
// Run it from the app/ directory:
//
//	go run ./cmd/helpcorpus
//
// The audience rules are in app/helpcorpus/audience.go.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"tybalt/helpcorpus"
)

func main() {
	repo := flag.String("repo", "..", "the folder that holds ui/src and descriptions")
	app := flag.String("app-dir", "", "the Go app folder, if it is not <repo>/app")
	out := flag.String("out", filepath.Join("helpcorpus", "data", "corpus.json"), "the corpus file to write")
	quiet := flag.Bool("quiet", false, "print no summary")
	flag.Parse()

	if err := run(*repo, *app, *out, *quiet); err != nil {
		fmt.Fprintln(os.Stderr, "helpcorpus:", err)
		os.Exit(1)
	}
}

func run(repo, app, out string, quiet bool) error {
	if err := helpcorpus.CheckOutPath(out); err != nil {
		return err
	}
	corpus, stats, err := helpcorpus.Generate(helpcorpus.Options{RepoDir: repo, AppDir: app})
	if err != nil {
		return err
	}
	data, err := corpus.Encode()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(out, data, 0o644); err != nil {
		return err
	}
	if !quiet {
		printSummary(corpus, stats, out, len(data))
	}
	return nil
}

func printSummary(corpus *helpcorpus.Corpus, stats *helpcorpus.Stats, out string, size int) {
	type count struct{ chunks, bytes int }
	byKind := map[string]*count{}
	byAudience := map[string]*count{}
	add := func(counts map[string]*count, key string, chunk helpcorpus.Chunk) {
		if counts[key] == nil {
			counts[key] = &count{}
		}
		counts[key].chunks++
		counts[key].bytes += len(chunk.Text)
	}
	for _, chunk := range corpus.Chunks {
		add(byKind, chunk.Kind, chunk)
		add(byAudience, audienceLabel(chunk), chunk)
	}
	show := func(title string, counts map[string]*count) {
		fmt.Println(title)
		keys := make([]string, 0, len(counts))
		for key := range counts {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			fmt.Printf("  %-28s %5d chunks %8d bytes\n", key, counts[key].chunks, counts[key].bytes)
		}
	}

	fmt.Printf("Wrote %s: %d bytes, %d files, %d chunks\n", out, size, stats.Files, len(corpus.Chunks))
	fmt.Printf("Hash %s\n", corpus.Hash)
	fmt.Printf("Handbook: %d bytes\n", len(corpus.Handbook))
	show("By kind:", byKind)
	show("By audience:", byAudience)
	for _, skipped := range stats.Skipped {
		fmt.Println("Skipped", skipped)
	}
}

// audienceLabel shows who may read a chunk, such as "admin + (admin or commit)".
func audienceLabel(chunk helpcorpus.Chunk) string {
	if chunk.ForEveryone() {
		return "everyone"
	}
	groups := make([]string, len(chunk.Requires))
	for i, group := range chunk.Requires {
		groups[i] = strings.Join(group, " or ")
		if len(group) > 1 && len(chunk.Requires) > 1 {
			groups[i] = "(" + groups[i] + ")"
		}
	}
	return strings.Join(groups, " + ")
}
