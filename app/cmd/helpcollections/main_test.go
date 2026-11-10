package main

import (
	"path/filepath"
	"strings"
	"testing"
)

// The collection files are written by hand of nobody: go run
// ./cmd/helpcollections makes them from the migrations. This test fails when a
// migration changed a collection and the files were not written again.
func TestTheCollectionFilesMatchTheMigrations(t *testing.T) {
	files, err := Render()
	if err != nil {
		t.Fatal(err)
	}
	if err := Check(filepath.Join("..", "..", "helpcorpus", "collections"), files); err != nil {
		t.Error(err)
	}

	// What a file must hold, for one collection whose rules are known.
	profiles := files["profiles.md"]
	for _, part := range []string{
		"# Collection: profiles",
		"## Who can update a record (update rule): profiles\n\n```\n@request.auth.id != \"\" &&\nuid = @request.auth.id &&\n@request.body.uid:changed = false\n```",
		"## Who can delete a record (delete rule): profiles\n\nNo rule. Nobody can do this through the API, except a PocketBase superuser.",
		"| `manager` | relation | required; maxSelect 1; to `users` |",
		"These rules do not limit Turbo's own server code",
	} {
		if !strings.Contains(profiles, part) {
			t.Errorf("profiles.md lacks %q", part)
		}
	}
	// PocketBase's own collections have no file, and a view has its query and
	// no rules for writing.
	if _, ok := files["_superusers.md"]; ok {
		t.Error("a system collection must have no file")
	}
	for name, content := range files {
		if strings.Contains(content, "is a view") && (strings.Contains(content, "(create rule)") || !strings.Contains(content, "## The query of the view")) {
			t.Errorf("%s: a view has a query and no create rule", name)
		}
	}
}

func TestCheckNamesEachProblem(t *testing.T) {
	dir := t.TempDir()
	if err := write(dir, map[string]string{"a.md": "one\n", "b.md": "two\n", "gone.md": "old\n"}); err != nil {
		t.Fatal(err)
	}
	err := Check(dir, map[string]string{"a.md": "one\n", "b.md": "changed\n", "c.md": "new\n"})
	if err == nil {
		t.Fatal("must be an error")
	}
	for _, part := range []string{"b.md is not what the migrations give", "c.md is missing", "gone.md is for a collection that is gone", "go run ./cmd/helpcollections"} {
		if !strings.Contains(err.Error(), part) {
			t.Errorf("the error lacks %q:\n%v", part, err)
		}
	}
	// Writing again removes the file of a collection that is gone.
	files := map[string]string{"a.md": "one\n", "b.md": "changed\n", "c.md": "new\n"}
	if err := write(dir, files); err != nil {
		t.Fatal(err)
	}
	if err := Check(dir, files); err != nil {
		t.Error(err)
	}
}
