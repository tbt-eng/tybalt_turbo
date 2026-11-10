// Command helpcollections writes one file for each PocketBase collection, with
// its access rules and its fields, for Help chat to search.
//
// The UI reads and writes many collections through PocketBase's own API, and
// there the collection's rules decide who may list, view, create, update and
// delete a record. Those rules are in the migrations, which the corpus leaves
// out, because an old migration says what a rule was and not what it is. This
// command applies every migration to an empty database and writes what the
// collections are at the end.
//
// Run it from the app/ directory after a migration that changes a collection:
//
//	go run ./cmd/helpcollections           write helpcorpus/collections/*.md
//	go run ./cmd/helpcollections --check   fail if a file is not what the migrations give
//
// It calls no model and needs no key. A changed file makes its rule file old:
// run go run ./cmd/helprules after it. See helpcorpus/rules.go.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pocketbase/pocketbase/core"

	_ "tybalt/migrations"
)

func main() {
	dir := flag.String("dir", filepath.Join("helpcorpus", "collections"), "the folder of the collection files")
	check := flag.Bool("check", false, "write nothing, and fail if a file is not what the migrations give")
	flag.Parse()

	files, err := Render()
	if err == nil {
		if *check {
			err = Check(*dir, files)
		} else {
			err = write(*dir, files)
		}
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "helpcollections:", err)
		os.Exit(1)
	}
}

// Render applies every migration to an empty database and returns the file of
// each collection, by file name.
func Render() (map[string]string, error) {
	data, err := os.MkdirTemp("", "helpcollections")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(data)
	app := core.NewBaseApp(core.BaseAppConfig{DataDir: data})
	if err := app.Bootstrap(); err != nil {
		return nil, err
	}
	defer app.ResetBootstrapState()
	if err := app.RunAllMigrations(); err != nil {
		return nil, err
	}
	collections, err := app.FindAllCollections()
	if err != nil {
		return nil, err
	}
	names := map[string]string{}
	for _, c := range collections {
		names[c.Id] = c.Name
	}
	files := map[string]string{}
	for _, c := range collections {
		// PocketBase's own collections, such as _superusers, are not Turbo's.
		if c.System {
			continue
		}
		files[c.Name+".md"] = render(c, names)
	}
	return files, nil
}

const header = "<!-- `go run ./cmd/helpcollections` wrote this file from the migrations. Do not edit it. -->\n\n"

func render(c *core.Collection, names map[string]string) string {
	var out strings.Builder
	out.WriteString(header)
	fmt.Fprintf(&out, "# Collection: %s\n\n", c.Name)
	switch c.Type {
	case core.CollectionTypeView:
		fmt.Fprintf(&out, "`%s` is a view: a saved query that Turbo reads and never writes. ", c.Name)
	case core.CollectionTypeAuth:
		fmt.Fprintf(&out, "`%s` holds the accounts that people sign in with. ", c.Name)
	default:
		fmt.Fprintf(&out, "`%s` is a collection of records. ", c.Name)
	}
	fmt.Fprintf(&out, "Turbo's screens reach it through PocketBase's API (`/api/collections/%s/records`), and the rules below decide who may do each action there. "+
		"A rule is a filter: a request passes when the filter is true. In a rule, `@request.auth` is the signed-in person, `@request.body` is what the request sends, "+
		"and a plain name such as `uid` is a field of the record. These rules do not limit Turbo's own server code, such as a hook or a custom route.\n", c.Name)

	rules := []struct {
		title string
		rule  *string
	}{
		{"Who can list the records (list rule)", c.ListRule},
		{"Who can view one record (view rule)", c.ViewRule},
		{"Who can create a record (create rule)", c.CreateRule},
		{"Who can update a record (update rule)", c.UpdateRule},
		{"Who can delete a record (delete rule)", c.DeleteRule},
	}
	if c.Type == core.CollectionTypeView {
		rules = rules[:2]
	}
	for _, r := range rules {
		fmt.Fprintf(&out, "\n## %s: %s\n\n", r.title, c.Name)
		switch {
		case r.rule == nil:
			out.WriteString("No rule. Nobody can do this through the API, except a PocketBase superuser.\n")
		case strings.TrimSpace(*r.rule) == "":
			out.WriteString("An empty rule. Everybody can do this, also a person who is not signed in.\n")
		default:
			fmt.Fprintf(&out, "```\n%s\n```\n", strings.TrimSpace(*r.rule))
		}
	}

	if c.Type == core.CollectionTypeView && strings.TrimSpace(c.ViewQuery) != "" {
		fmt.Fprintf(&out, "\n## The query of the view: %s\n\n```sql\n%s\n```\n", c.Name, strings.TrimSpace(c.ViewQuery))
	}

	fmt.Fprintf(&out, "\n## Fields: %s\n\n| Field | Type | Details |\n| --- | --- | --- |\n", c.Name)
	for _, field := range c.Fields {
		fmt.Fprintf(&out, "| `%s` | %s | %s |\n", field.GetName(), field.Type(), details(field, names))
	}

	if len(c.Indexes) > 0 {
		fmt.Fprintf(&out, "\n## Indexes: %s\n\nA unique index means that two records can't have the same values in its fields.\n\n", c.Name)
		for _, index := range c.Indexes {
			fmt.Fprintf(&out, "- `%s`\n", strings.Join(strings.Fields(index), " "))
		}
	}
	return out.String()
}

// detailKeys are the settings of a field that change what a person can enter,
// in the order that they are written.
var detailKeys = []string{"required", "min", "max", "onlyInt", "pattern", "values", "minSelect", "maxSelect", "collectionId", "cascadeDelete", "maxSize", "mimeTypes", "onCreate", "onUpdate", "hidden"}

func details(field core.Field, names map[string]string) string {
	raw, err := json.Marshal(field)
	if err != nil {
		return ""
	}
	settings := map[string]any{}
	if err := json.Unmarshal(raw, &settings); err != nil {
		return ""
	}
	var parts []string
	for _, key := range detailKeys {
		value, set := settings[key]
		if !set || isZero(value) {
			continue
		}
		switch key {
		case "collectionId":
			id, _ := value.(string)
			name := names[id]
			if name == "" {
				name = id
			}
			parts = append(parts, "to `"+name+"`")
		case "required", "onlyInt", "cascadeDelete", "onCreate", "onUpdate", "hidden":
			parts = append(parts, key)
		default:
			text, _ := json.Marshal(value)
			parts = append(parts, fmt.Sprintf("%s %s", key, strings.ReplaceAll(string(text), "|", "\\|")))
		}
	}
	return strings.Join(parts, "; ")
}

func isZero(value any) bool {
	switch v := value.(type) {
	case nil:
		return true
	case bool:
		return !v
	case float64:
		return v == 0
	case string:
		return v == ""
	case []any:
		return len(v) == 0
	}
	return false
}

func write(dir string, files map[string]string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
			return err
		}
	}
	// A collection that is gone has no file.
	old, err := filepath.Glob(filepath.Join(dir, "*.md"))
	if err != nil {
		return err
	}
	for _, path := range old {
		if _, ok := files[filepath.Base(path)]; !ok {
			if err := os.Remove(path); err != nil {
				return err
			}
		}
	}
	fmt.Printf("Wrote %d collection files to %s\n", len(files), dir)
	return nil
}

// Check compares the files on disk with what the migrations give.
func Check(dir string, files map[string]string) error {
	var wrong []string
	for name, content := range files {
		have, err := os.ReadFile(filepath.Join(dir, name))
		switch {
		case os.IsNotExist(err):
			wrong = append(wrong, name+" is missing")
		case err != nil:
			return err
		case string(have) != content:
			wrong = append(wrong, name+" is not what the migrations give")
		}
	}
	old, err := filepath.Glob(filepath.Join(dir, "*.md"))
	if err != nil {
		return err
	}
	for _, path := range old {
		if _, ok := files[filepath.Base(path)]; !ok {
			wrong = append(wrong, filepath.Base(path)+" is for a collection that is gone")
		}
	}
	if len(wrong) == 0 {
		return nil
	}
	sort.Strings(wrong)
	return fmt.Errorf("the collection files in %s do not match the migrations:\n  %s\nIn app/, run: go run ./cmd/helpcollections\nThen run go run ./cmd/helprules, with the key of the model provider in the shell, and commit the files",
		dir, strings.Join(wrong, "\n  "))
}
