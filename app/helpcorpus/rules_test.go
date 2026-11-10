package helpcorpus

import (
	"maps"
	"reflect"
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
			{Title: "Second", Text: "Lines 3-4 are not the lines of this rule.", Start: 20, End: 20,
				Asks: []string{"my own manager", "approve; my own\nexpense", " "}},
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
	if !reflect.DeepEqual(got.Rules[0], want) {
		t.Errorf("rule 1 = %+v", got.Rules[0])
	}
	if got.Rules[1].Start != 20 || got.Rules[1].Text != "Lines 3-4 are not the lines of this rule." {
		t.Errorf("rule 2 = %+v", got.Rules[1])
	}
	// The phrases come back, each on its own, and the text does not hold them.
	if !reflect.DeepEqual(got.Rules[1].Asks, []string{"my own manager", "approve, my own expense"}) || got.Rules[0].Asks != nil {
		t.Errorf("asks = %q and %q", got.Rules[0].Asks, got.Rules[1].Asks)
	}
	if !strings.Contains(text, "\nAsked as: my own manager; approve, my own expense\n") {
		t.Errorf("the phrases are one line of the file:\n%s", text)
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
	yes := []string{"app/hooks/expenses.go", "app/routes/po_visibility_base.sql", "app/utilities/config.go", "ui/src/lib/poVisibility.ts", "app/helpcorpus/collections/profiles.md"}
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
	if strings.Join(paths, " ") != "app/hooks/expenses.go app/routes/expenses.sql app/helpcorpus/collections/profiles.md ui/src/lib/utilities.ts" {
		t.Errorf("RuleSources = %v", paths)
	}
}

func TestTheSearchMatchesAPhraseThatTheModelDoesNotSee(t *testing.T) {
	const source = "package routes\n\nfunc submit() {}\n"
	files := maps.Clone(fixture)
	files["app/routes/submit.go"] = source
	files["app/helpcorpus/rules/app/routes/submit.go.md"] = RuleFile{
		Source: "app/routes/submit.go", SHA256: SourceHash(source), Model: "test-model",
		Rules: []Rule{{
			Title: "Submit approves an expense when the owner is also the approver",
			Text:  "Turbo approves the expense at submit.", Start: 3, End: 3,
			Asks: []string{"my own manager", "self approval"},
		}},
	}.Format()
	corpus, _, err := Generate(Options{RepoDir: writeRepo(t, files, false)})
	if err != nil {
		t.Fatal(err)
	}
	index, err := NewIndex(corpus)
	if err != nil {
		t.Fatal(err)
	}
	defer index.Close()

	// "manager" is in the phrases only.
	hits, err := index.Search(Claims{}, "my own manager", SearchOptions{Kinds: []string{KindRules}})
	if err != nil || len(hits) != 1 || hits[0].Start != 3 {
		t.Fatalf("hits = %+v, %v", hits, err)
	}
	if hits[0].Text != "Turbo approves the expense at submit." || hits[0].Asks != "my own manager. self approval" {
		t.Errorf("text = %q, asks = %q", hits[0].Text, hits[0].Asks)
	}
}

// The rules of a source file are written by a model, with a key, so nothing
// writes them in a build or in a test. This test is what keeps them current: it
// fails when a source file changed and its rule file did not.
func TestTheRuleFilesMatchTheSource(t *testing.T) {
	status, err := CheckRules(Options{RepoDir: "../.."})
	if err != nil {
		t.Fatal(err)
	}
	if err := status.Error(); err != nil {
		t.Error(err)
	}
}

func TestCheckRulesNamesEachProblem(t *testing.T) {
	files := maps.Clone(fixture)
	files["app/hooks/expenses.go"] = expensesSource
	// Rules for an older source, a rule with no phrases, and a rule file with no source.
	files["app/helpcorpus/rules/app/hooks/expenses.go.md"] = ruleFixture("an older hash")
	sql := fixture["app/routes/expenses.sql"]
	files["app/helpcorpus/rules/app/routes/expenses.sql.md"] = RuleFile{
		Source: "app/routes/expenses.sql", SHA256: SourceHash(sql),
		Rules: []Rule{{Title: "One", Text: "Text.", Start: 1, End: 1}},
	}.Format()
	files["app/helpcorpus/rules/app/routes/deleted.go.md"] = RuleFile{Source: "app/routes/deleted.go", SHA256: "abc"}.Format()
	repo := writeRepo(t, files, false)

	status, err := CheckRules(Options{RepoDir: repo})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(status.Stale, []string{"app/hooks/expenses.go"}) ||
		!reflect.DeepEqual(status.NoPhrases, []string{"app/routes/expenses.sql"}) ||
		!reflect.DeepEqual(status.Missing, []string{"app/helpcorpus/collections/profiles.md", "ui/src/lib/utilities.ts"}) ||
		len(status.Gone) != 1 || !strings.HasSuffix(status.Gone[0], "app/helpcorpus/rules/app/routes/deleted.go.md") {
		t.Errorf("status = %+v", status)
	}
	message := status.Error().Error()
	for _, part := range []string{"app/hooks/expenses.go changed after its rules were written", "go run ./cmd/helprules", "commit the rule files"} {
		if !strings.Contains(message, part) {
			t.Errorf("the error lacks %q:\n%s", part, message)
		}
	}

	// With every rule file as it must be, there is no error.
	files["app/helpcorpus/rules/app/hooks/expenses.go.md"] = RuleFile{Source: "app/hooks/expenses.go", SHA256: SourceHash(expensesSource)}.Format()
	files["app/helpcorpus/rules/app/routes/expenses.sql.md"] = RuleFile{Source: "app/routes/expenses.sql", SHA256: SourceHash(sql)}.Format()
	files["app/helpcorpus/rules/ui/src/lib/utilities.ts.md"] = RuleFile{Source: "ui/src/lib/utilities.ts", SHA256: SourceHash(fixture["ui/src/lib/utilities.ts"])}.Format()
	const collection = "app/helpcorpus/collections/profiles.md"
	files["app/helpcorpus/rules/"+collection+".md"] = RuleFile{Source: collection, SHA256: SourceHash(fixture[collection])}.Format()
	delete(files, "app/helpcorpus/rules/app/routes/deleted.go.md")
	status, err = CheckRules(Options{RepoDir: writeRepo(t, files, false)})
	if err != nil || !status.OK() || status.Error() != nil {
		t.Errorf("status = %+v, %v", status, err)
	}
}

func TestARuleOfACollectionNamesItsCollectionAndItsAction(t *testing.T) {
	// A model wrote the title in plain words. The chat searches for
	// "profiles update rule", so the heading starts with those words.
	collection := file{path: CollectionsRoot + "/profiles.md", content: strings.Join([]string{
		"# Collection: profiles", // line 1
		"",
		"## Who can list the records (list rule): profiles", // 3
		"",
		"uid = @request.auth.id",
		"",
		"## Who can view one record (view rule): profiles", // 7
		"",
		"@request.auth.id != \"\"",
		"",
		"## Who can update a record (update rule): profiles", // 11
		"",
		"uid = @request.auth.id",
		"",
		"## Fields: profiles", // 15
		"",
		"| `surname` | text | required |",
	}, "\n")}
	for _, test := range []struct {
		start, end int
		want       string
	}{
		{13, 13, "profiles collection, update rule: "},
		{5, 9, "profiles collection, list rule, view rule: "},
		{17, 17, "profiles collection, fields: "},
		{1, 1, "profiles collection: "},
	} {
		if got := collectionPrefix(collection, Rule{Start: test.start, End: test.end}); got != test.want {
			t.Errorf("lines %d-%d: got %q, want %q", test.start, test.end, got, test.want)
		}
	}
	// A rule about a field says that the code can ask for more.
	collection.path = CollectionsRoot + "/client_notes.md"
	rules := RuleFile{Source: collection.path, SHA256: SourceHash(collection.content), Rules: []Rule{
		{Title: "The job is optional", Text: "The job is optional.", Start: 17, End: 17, Asks: []string{"note with no job"}},
		{Title: "Who can update a note", Text: "Only the owner.", Start: 13, End: 13, Asks: []string{"edit a note"}},
	}}
	repo := writeRepo(t, map[string]string{
		collection.path: collection.content,
		"app/helpcorpus/rules/" + collection.path + ".md": rules.Format(),
	}, false)
	chunks, stale, err := ruleChunks(Options{RepoDir: repo}, collection)
	if err != nil || stale || len(chunks) != 2 {
		t.Fatalf("chunks = %+v, stale = %v, err = %v", chunks, stale, err)
	}
	if !strings.HasSuffix(chunks[0].Text, fieldCaveat) || chunks[0].Heading != "client_notes collection, fields: The job is optional" {
		t.Errorf("the rule of a field: %q, %q", chunks[0].Heading, chunks[0].Text)
	}
	if chunks[1].Text != "Only the owner." {
		t.Errorf("an access rule must keep its text: %q", chunks[1].Text)
	}

	// A rule of the code keeps its title as it is.
	if got := collectionPrefix(file{path: "app/hooks/profiles.go", content: "package hooks"}, Rule{Start: 1, End: 1}); got != "" {
		t.Errorf("a rule of the code got the prefix %q", got)
	}
}
