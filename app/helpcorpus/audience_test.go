package helpcorpus

import (
	"reflect"
	"testing"
)

// useRules swaps the audience rules for one test.
func useRules(t *testing.T, paths []pathRule, terms []termRule) {
	t.Helper()
	oldPaths, oldTerms := pathRules, termRules
	pathRules, termRules = paths, terms
	t.Cleanup(func() { pathRules, termRules = oldPaths, oldTerms })
}

func TestRequiresFor(t *testing.T) {
	useRules(t,
		[]pathRule{
			{Prefix: "ui/src/docs/", Requires: everyone},
			{Prefix: "app/reports/", Requires: []string{"report"}},
			{Prefix: "app/reports/open/", Requires: everyone},
			{Prefix: "app/routes/commit", Requires: []string{"commit"}},
		},
		[]termRule{
			{Term: "uncommit", Requires: []string{"commit", "admin"}},
			{Term: "salary", Requires: []string{"hr"}},
		},
	)
	cases := []struct {
		name, path, text string
		want             [][]string
	}{
		{"user docs are for everyone", "ui/src/docs/time.md", "Submit a week.", [][]string{}},
		{"a path rule names its claim", "app/reports/payroll.sql", "SELECT 1", [][]string{{"report"}}},
		{"the longest prefix wins", "app/reports/open/list.sql", "SELECT 1", [][]string{}},
		{"no rule means admin", "app/new_folder/thing.go", "package thing", [][]string{{"admin"}}},
		{"a term restricts text that is for everyone", "ui/src/docs/time.md", "You can Uncommit it.", [][]string{{"admin", "commit"}}},
		{"a term in the path counts", "app/reports/open/uncommit_list.sql", "SELECT 1", [][]string{{"admin", "commit"}}},
		{"a term adds to a path rule", "app/reports/payroll.sql", "-- uncommitted rows", [][]string{{"report"}, {"admin", "commit"}}},
		{"two terms both apply", "app/reports/open/x.sql", "salary after UNCOMMIT", [][]string{{"hr"}, {"admin", "commit"}}},
		{"a smaller group implies a bigger one", "app/new_folder/thing.go", "uncommit", [][]string{{"admin"}}},
		{"the same group twice is one group", "app/routes/commit_record.go", "uncommit", [][]string{{"commit"}}},
	}
	for _, c := range cases {
		if got := requiresFor(c.path, c.text); !reflect.DeepEqual(got, c.want) {
			t.Errorf("%s: got %v, want %v", c.name, got, c.want)
		}
	}
}

func TestAllows(t *testing.T) {
	claims := func(names ...string) map[string]bool {
		held := map[string]bool{}
		for _, name := range names {
			held[name] = true
		}
		return held
	}
	everyoneChunk := Chunk{Requires: [][]string{}}
	oneGroup := Chunk{Requires: [][]string{{"admin", "commit"}}}
	twoGroups := Chunk{Requires: [][]string{{"report"}, {"admin", "commit"}}}

	if !everyoneChunk.ForEveryone() || oneGroup.ForEveryone() {
		t.Fatal("ForEveryone is wrong")
	}
	checks := []struct {
		name  string
		chunk Chunk
		held  map[string]bool
		want  bool
	}{
		{"everyone, no claims", everyoneChunk, claims(), true},
		{"everyone, nil claims", everyoneChunk, nil, true},
		{"one group, no claims", oneGroup, claims(), false},
		{"one group, another claim", oneGroup, claims("time", "report"), false},
		{"one group, one of its claims", oneGroup, claims("time", "commit"), true},
		{"two groups, only the first", twoGroups, claims("report"), false},
		{"two groups, only the second", twoGroups, claims("admin"), false},
		{"two groups, both", twoGroups, claims("report", "commit"), true},
	}
	for _, c := range checks {
		if got := c.chunk.Allows(c.held); got != c.want {
			t.Errorf("%s: got %v, want %v", c.name, got, c.want)
		}
	}
}

// The repo is public, so the shipped rules open all of it to signed-in users.
// A source that nobody has added a rule for stays closed.
func TestShippedRules(t *testing.T) {
	staff := map[string]bool{"time": true}
	for _, path := range []string{
		"ui/src/docs/expenses.md",
		"descriptions/uncommit.md",
		"app/hooks/expenses.go",
		"app/routes/uncommit_record.go",
		"ui/src/routes/reports/payroll/+page.svelte",
		"ui/src/lib/navConfig.ts",
	} {
		if chunk := (Chunk{Requires: requiresFor(path, "Select Uncommit to reverse it.")}); !chunk.ForEveryone() {
			t.Errorf("%s must be for everyone, got %v", path, chunk.Requires)
		}
	}
	closed := Chunk{Requires: requiresFor("import_data/main.go", "package main")}
	if closed.Allows(staff) || !closed.Allows(map[string]bool{"admin": true}) {
		t.Errorf("a source with no rule must be for admins only, got %v", closed.Requires)
	}
}

func TestHandbookBans(t *testing.T) {
	for _, text := range []string{"Select Uncommit.", "ask about uncommitting it", "An admin uncommits it.", "the uncommit request"} {
		if !handbookBans.MatchString(text) {
			t.Errorf("%q must not pass", text)
		}
	}
	for _, text := range []string{"An uncommitted time sheet.", "Commit it. It is final.", "recommit"} {
		if handbookBans.MatchString(text) {
			t.Errorf("%q must pass", text)
		}
	}
}
