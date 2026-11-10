package helpcorpus

import (
	"maps"
	"strings"
	"testing"
)

const expensesSource = "package hooks\n\n// commit is final for staff.\nfunc commit() {}\n\n// uncommit reverses it.\nfunc uncommit() {}\n"

func ruleFixture(sourceHash string) string {
	return RuleFile{
		Source: "app/hooks/expenses.go", SHA256: sourceHash, Model: "test-model",
		Rules: []Rule{
			{Title: "Submit approves your own expense", Text: "When the owner of an expense is also its approver, Turbo approves the expense at submit.", Start: 3, End: 4},
			{Title: "A committer can uncommit", Text: "Uncommit reverses a commit.", Start: 6, End: 7},
		},
	}.Format()
}

func TestRuleFileFormatAndParse(t *testing.T) {
	file := RuleFile{
		Source: "app/hooks/expenses.go", SHA256: "abc", Model: "test-model",
		Rules: []Rule{
			{Title: "A limit\nof 18 hours", Text: "One entry holds at most 18 hours.\n## Not a heading\nMeals count.", Start: 10, End: 12},
			{Title: "Second", Text: "Lines 3-4 are not the lines of this rule.", Start: 20, End: 20},
		},
	}
	text := file.Format()
	if !strings.Contains(text, "Do not edit it") || !strings.Contains(text, "# Rules in app/hooks/expenses.go") {
		t.Errorf("the file must say what it is:\n%s", text)
	}
	got, err := ParseRuleFile(text)
	if err != nil {
		t.Fatal(err)
	}
	if got.Source != file.Source || got.SHA256 != "abc" || got.Model != "test-model" || len(got.Rules) != 2 {
		t.Fatalf("got %+v", got)
	}
	// A line of text can't become a heading, and the text keeps its words.
	want := Rule{Title: "A limit of 18 hours", Text: "One entry holds at most 18 hours. ## Not a heading Meals count.", Start: 10, End: 12}
	if got.Rules[0] != want {
		t.Errorf("rule 1 = %+v", got.Rules[0])
	}
	if got.Rules[1].Start != 20 || got.Rules[1].Text != "Lines 3-4 are not the lines of this rule." {
		t.Errorf("rule 2 = %+v", got.Rules[1])
	}

	// A file with no rules is a valid file.
	empty, err := ParseRuleFile(RuleFile{Source: "app/routes/x.go", SHA256: "abc"}.Format())
	if err != nil || len(empty.Rules) != 0 {
		t.Errorf("empty = %+v, %v", empty, err)
	}
	for _, bad := range []string{"", "# Rules\n\n## A rule\n\nLines 1-2\n\nText.\n", strings.Replace(text, "Lines 10-12", "Lines 12-10", 1)} {
		if _, err := ParseRuleFile(bad); err == nil {
			t.Errorf("must be an error:\n%s", bad)
		}
	}
}

func TestHasRules(t *testing.T) {
	yes := []string{"app/hooks/expenses.go", "app/routes/po_visibility_base.sql", "app/utilities/config.go", "ui/src/lib/poVisibility.ts"}
	no := []string{"app/main.go", "app/helpchat/turn.go", "ui/src/lib/components/DSChat.svelte", "ui/src/lib/stores/helpChat.ts", "ui/src/routes/time/+page.svelte", "descriptions/expenses.md"}
	for _, path := range yes {
		if !HasRules(path) {
			t.Errorf("%s gets rules", path)
		}
	}
	for _, path := range no {
		if HasRules(path) {
			t.Errorf("%s gets no rules", path)
		}
	}
}

func TestTheCorpusTakesRulesOnlyForTheSourceAsItIs(t *testing.T) {
	const rulePath = "app/helpcorpus/rules/app/hooks/expenses.go.md"
	rules := func(corpus *Corpus) []Chunk {
		var out []Chunk
		for _, chunk := range corpus.Chunks {
			if chunk.Kind == KindRules {
				out = append(out, chunk)
			}
		}
		return out
	}

	// Rules that were written from the source as it is.
	files := maps.Clone(fixture)
	files["app/hooks/expenses.go"] = expensesSource
	files[rulePath] = ruleFixture(SourceHash(expensesSource))
	corpus, stats, err := Generate(Options{RepoDir: writeRepo(t, files, false)})
	if err != nil {
		t.Fatal(err)
	}
	got := rules(corpus)
	// The rule about uncommit is left out: the chat treats a commit as final.
	if len(got) != 1 || stats.Rules != 1 || len(stats.StaleRules) != 0 {
		t.Fatalf("rules = %+v, stats = %+v", got, stats)
	}
	want := Chunk{
		ID: "rule:app/hooks/expenses.go:1:3", Path: "app/hooks/expenses.go", Kind: KindRules, Start: 3, End: 4,
		Heading: "Submit approves your own expense", Requires: [][]string{},
		Text: "When the owner of an expense is also its approver, Turbo approves the expense at submit.",
	}
	if got[0].ID != want.ID || got[0].Path != want.Path || got[0].Start != 3 || got[0].End != 4 || got[0].Heading != want.Heading || got[0].Text != want.Text || len(got[0].Requires) != 0 {
		t.Errorf("rule chunk = %+v", got[0])
	}
	// The rule file is not in the corpus as a file, and the other files that get
	// rules are reported.
	if chunkPaths(corpus)[rulePath] {
		t.Error("the rule file itself must not be a source")
	}
	if !strings.Contains(strings.Join(stats.NoRules, " "), "app/routes/expenses.sql") || strings.Contains(strings.Join(stats.NoRules, " "), "app/hooks/expenses.go") {
		t.Errorf("NoRules = %v", stats.NoRules)
	}

	// The source changed after the rules were written: no rules, and a report.
	files["app/hooks/expenses.go"] = expensesSource + "\nfunc recall() {}\n"
	corpus, stats, err = Generate(Options{RepoDir: writeRepo(t, files, false)})
	if err != nil {
		t.Fatal(err)
	}
	if len(rules(corpus)) != 0 || len(stats.StaleRules) != 1 || stats.StaleRules[0] != "app/hooks/expenses.go" {
		t.Errorf("stale rules must stay out: %+v, stats = %+v", rules(corpus), stats)
	}
	if !chunkPaths(corpus)["app/hooks/expenses.go"] {
		t.Error("the code stays in the corpus")
	}

	// A rule file that is not what Format writes stops the build.
	files[rulePath] = "## A rule\n\nNo lines.\n"
	if _, _, err := Generate(Options{RepoDir: writeRepo(t, files, false)}); err == nil {
		t.Error("a broken rule file must be an error")
	}

	// RuleSources lists the files that get rules, and no others.
	sources, err := RuleSources(Options{RepoDir: writeRepo(t, fixture, false)})
	if err != nil {
		t.Fatal(err)
	}
	var paths []string
	for _, source := range sources {
		paths = append(paths, source.Path)
	}
	if strings.Join(paths, " ") != "app/hooks/expenses.go app/routes/expenses.sql ui/src/lib/utilities.ts" {
		t.Errorf("RuleSources = %v", paths)
	}
}
