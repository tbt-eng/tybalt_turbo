package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"testing"

	"tybalt/helpchat"
	"tybalt/helpcorpus"
)

// stubModel writes one rule for each file, and records what it was asked.
type stubModel struct {
	mu      sync.Mutex
	asked   []string
	systems []string
	reply   func(question string) string
	// phraseAsked holds the requests of the second step, for search phrases.
	phraseAsked []string
	phraseReply func(question string) string
}

var numberedRule = regexp.MustCompile(`(?m)^\d+\. `)

func (m *stubModel) Warm(context.Context, helpchat.Request) error { return nil }
func (m *stubModel) Start(request helpchat.Request) helpchat.Session {
	return &stubSession{model: m, request: request}
}

type stubSession struct {
	model   *stubModel
	request helpchat.Request
}

func (s *stubSession) AddToolResults([]helpchat.ToolResult) {}
func (s *stubSession) Next(context.Context, func(string)) (*helpchat.Reply, error) {
	m := s.model
	m.mu.Lock()
	defer m.mu.Unlock()
	if s.request.Shared == phraseInstructions {
		m.phraseAsked = append(m.phraseAsked, s.request.Question)
		// Two phrases for each rule in the request.
		var entries []string
		for n := range numberedRule.FindAllString(s.request.Question, -1) {
			entries = append(entries, fmt.Sprintf(`{"n":%d,"phrases":["my own manager","self approval"]}`, n+1))
		}
		text := `{"asks":[` + strings.Join(entries, ",") + `]}`
		if m.phraseReply != nil {
			text = m.phraseReply(s.request.Question)
		}
		return &helpchat.Reply{Text: text, StopReason: "end_turn", Usage: helpchat.Usage{Input: 4, Output: 2}}, nil
	}
	m.asked = append(m.asked, s.request.Question)
	m.systems = append(m.systems, s.request.Shared)
	text := `Here are the rules: {"rules":[{"title":"A rule","text":"Turbo refuses it.","start":1,"end":2}]} Done.`
	if m.reply != nil {
		text = m.reply(s.request.Question)
	}
	return &helpchat.Reply{Text: text, StopReason: "end_turn", Usage: helpchat.Usage{Input: 10, Output: 5, CacheRead: 100}}, nil
}

func writeRepo(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for name, content := range files {
		full := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func testService(t *testing.T, model helpchat.Model) *helpchat.Service {
	t.Helper()
	index, err := helpcorpus.NewIndex(&helpcorpus.Corpus{Format: helpcorpus.FormatVersion})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { index.Close() })
	return &helpchat.Service{Model: model, Index: index}
}

var repoFiles = map[string]string{
	"ui/src/docs/glossary.md":    "# Glossary\n\n| Word | Meaning |\n| --- | --- |\n| Time sheet | Your time for one week. |\n",
	"ui/src/docs/time.md":        "# Time\n",
	"descriptions/notes.md":      "# Notes\n",
	"app/main.go":                "package main\n",
	"app/hooks/expenses.go":      "package hooks\n\nfunc submit() {}\n",
	"app/routes/approve.go":      "package routes\n\nfunc approve() {}\n",
	"app/routes/approve_test.go": "package routes\n",
	"ui/src/lib/utilities.ts":    "export const a = 1;\nexport const b = 2;\n",
	"ui/src/routes/+page.svelte": "<p>page</p>\n",
}

func TestWritesRulesOnlyForSourcesThatChanged(t *testing.T) {
	repo := writeRepo(t, repoFiles)
	model := &stubModel{}
	service := testService(t, model)
	o := options{repo: repo, model: "claude-sonnet-5-5", effort: "medium", phraseModel: "claude-haiku-5-5", phraseEffort: "low", parallel: 2}
	appDir := filepath.Join(repo, "app")

	// With no rules yet, a check fails and calls no model.
	check := o
	check.check = true
	if err := run(context.Background(), check, service, io.Discard); err == nil || len(model.asked) != 0 {
		t.Fatalf("check: err = %v, asked %d", err, len(model.asked))
	}

	if err := run(context.Background(), o, service, io.Discard); err != nil {
		t.Fatal(err)
	}
	// One request for each file that gets rules: not the tests, the pages or main.go.
	if len(model.asked) != 3 {
		t.Fatalf("asked %d times: %v", len(model.asked), model.asked)
	}
	for _, system := range model.systems {
		// The instructions and the glossary are the same for every file, so the
		// model caches them.
		if system != model.systems[0] || !strings.Contains(system, "Time sheet | Your time for one week.") {
			t.Fatal("every request must have the same instructions, with the glossary")
		}
	}
	if !strings.Contains(strings.Join(model.asked, "\n"), "File: app/hooks/expenses.go\n\n<source>\n    1| package hooks\n    2| \n    3| func submit() {}\n</source>") {
		t.Errorf("the request must hold the source with line numbers:\n%s", model.asked[0])
	}
	file, ok, err := helpcorpus.ReadRuleFile(appDir, "app/hooks/expenses.go")
	if err != nil || !ok || len(file.Rules) != 1 || file.Rules[0].Title != "A rule" || file.Model != "claude-sonnet-5-5" ||
		file.SHA256 != helpcorpus.SourceHash(repoFiles["app/hooks/expenses.go"]) {
		t.Fatalf("rule file = %+v, %v, %v", file, ok, err)
	}
	// The second step wrote the search phrases of each rule, from the rule and
	// not from the code.
	if len(model.phraseAsked) != 3 || !reflect.DeepEqual(file.Rules[0].Asks, []string{"my own manager", "self approval"}) {
		t.Fatalf("phrase requests = %d, asks = %q", len(model.phraseAsked), file.Rules[0].Asks)
	}
	if got := strings.Join(model.phraseAsked, "\n"); !strings.Contains(got, "The rules of app/hooks/expenses.go:\n\n1. A rule\nTurbo refuses it.\n") || strings.Contains(got, "func submit") {
		t.Errorf("a phrase request holds the rules and not the code:\n%s", got)
	}
	// The corpus takes the rules.
	_, stats, err := helpcorpus.Generate(helpcorpus.Options{RepoDir: repo})
	if err != nil || stats.Rules != 3 || len(stats.NoRules)+len(stats.StaleRules) != 0 {
		t.Fatalf("stats = %+v, %v", stats, err)
	}

	// A second run asks nothing, and a check passes.
	if err := run(context.Background(), o, service, io.Discard); err != nil || len(model.asked) != 3 || len(model.phraseAsked) != 3 {
		t.Fatalf("second run: %v, asked %d and %d", err, len(model.asked), len(model.phraseAsked))
	}
	if err := run(context.Background(), check, service, io.Discard); err != nil {
		t.Errorf("check after a full run: %v", err)
	}

	// One source changes and one is deleted: one request, and one file removed.
	if err := os.WriteFile(filepath.Join(appDir, "hooks", "expenses.go"), []byte("package hooks\n\nfunc submit() {}\n\nfunc recall() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(appDir, "routes", "approve.go")); err != nil {
		t.Fatal(err)
	}
	if err := run(context.Background(), check, service, io.Discard); err == nil {
		t.Error("check must fail when a source changed")
	}
	if err := run(context.Background(), o, service, io.Discard); err != nil {
		t.Fatal(err)
	}
	if len(model.asked) != 4 || !strings.Contains(model.asked[3], "File: app/hooks/expenses.go") {
		t.Errorf("asked %d times, last: %.60s", len(model.asked), model.asked[len(model.asked)-1])
	}
	// New rules have no phrases, so the second step runs for that one file.
	if len(model.phraseAsked) != 4 {
		t.Errorf("phrase requests = %d", len(model.phraseAsked))
	}
	if _, err := os.Stat(helpcorpus.RuleFilePath(appDir, "app/routes/approve.go")); !os.IsNotExist(err) {
		t.Error("the rules of a deleted source must be removed")
	}

	// --only and --force write the rules of one file again.
	again := o
	again.only, again.force = "utilities.ts", true
	if err := run(context.Background(), again, service, io.Discard); err != nil || len(model.asked) != 5 {
		t.Errorf("only and force: %v, asked %d", err, len(model.asked))
	}
}

func TestABadReplyWritesNoRuleFile(t *testing.T) {
	repo := writeRepo(t, repoFiles)
	bad := []string{
		"I can't do that.",
		`{"rules":[{"title":"","text":"No title.","start":1,"end":2}]}`,
		// The file has 3 lines.
		`{"rules":[{"title":"Wrong lines","text":"Text.","start":2,"end":9}]}`,
		`{"rules":[{"title":"Wrong lines","text":"Text.","start":3,"end":2}]}`,
	}
	for _, reply := range bad {
		model := &stubModel{reply: func(string) string { return reply }}
		o := options{repo: repo, model: "claude-sonnet-5-5", effort: "medium", phraseModel: "claude-haiku-5-5", phraseEffort: "low", parallel: 1, only: "app/hooks/expenses.go"}
		if err := run(context.Background(), o, testService(t, model), io.Discard); err == nil {
			t.Errorf("%s: must be an error", reply)
		}
		// It asks one more time, then gives up.
		if len(model.asked) != 2 {
			t.Errorf("%s: asked %d times", reply, len(model.asked))
		}
		if _, ok, _ := helpcorpus.ReadRuleFile(filepath.Join(repo, "app"), "app/hooks/expenses.go"); ok {
			t.Errorf("%s: a rule file was written", reply)
		}
	}

	// A file with no rules is a good reply, and is not asked about again.
	model := &stubModel{reply: func(string) string { return `{"rules":[]}` }}
	o := options{repo: repo, model: "claude-sonnet-5-5", effort: "medium", phraseModel: "claude-haiku-5-5", phraseEffort: "low", parallel: 1, only: "app/hooks/expenses.go"}
	service := testService(t, model)
	for i := 0; i < 2; i++ {
		if err := run(context.Background(), o, service, io.Discard); err != nil {
			t.Fatal(err)
		}
	}
	if len(model.asked) != 1 {
		t.Errorf("a file with no rules was asked about %d times", len(model.asked))
	}
}

func TestAsksAgainAboutThePartsWithNoRule(t *testing.T) {
	// A file of 61 lines. The first reply covers lines 1-2 only.
	var source strings.Builder
	source.WriteString("package hooks\n")
	for i := 0; i < 30; i++ {
		source.WriteString("// a comment is not a line of code\nfunc check() {}\n")
	}
	files := map[string]string{"ui/src/docs/glossary.md": "# Glossary\n", "app/hooks/expenses.go": source.String()}
	repo := writeRepo(t, files)

	model := &stubModel{}
	model.reply = func(question string) string {
		switch {
		case !strings.Contains(question, "have no rule yet"):
			return `{"rules":[{"title":"First","text":"One.","start":1,"end":2}]}`
		case !strings.Contains(question, "Second"):
			// The follow-up names the rule that is written and the lines that are not.
			if !strings.Contains(question, "- lines 1-2: First") || !strings.Contains(question, "These lines have no rule yet: 3-61.") {
				t.Errorf("follow-up = %s", question[strings.Index(question, "These rules are written"):])
			}
			return `{"rules":[{"title":"Second","text":"Two.","start":40,"end":61},{"title":"first","text":"Again.","start":1,"end":2}]}`
		default:
			return `{"rules":[]}`
		}
	}
	o := options{repo: repo, model: "claude-sonnet-5-5", effort: "medium", phraseModel: "claude-haiku-5-5", phraseEffort: "low", parallel: 1}
	if err := run(context.Background(), o, testService(t, model), io.Discard); err != nil {
		t.Fatal(err)
	}
	// The first request and one follow-up. What is left with no rule after that
	// is too small for another.
	if len(model.asked) != 2 {
		t.Errorf("asked %d times", len(model.asked))
	}
	file, _, err := helpcorpus.ReadRuleFile(filepath.Join(repo, "app"), "app/hooks/expenses.go")
	if err != nil || len(file.Rules) != 2 || file.Rules[0].Title != "First" || file.Rules[1].Title != "Second" {
		t.Errorf("rules = %+v, %v", file.Rules, err)
	}

	// Comments and blank lines are not a gap.
	if gaps, uncovered := gapsIn("// a\n\n// b\n"+strings.Repeat("// c\n", 40), nil); len(gaps) != 0 || uncovered != 0 {
		t.Errorf("gaps = %v, %d", gaps, uncovered)
	}
}

func TestRulesWithNoPhrasesGetThem(t *testing.T) {
	// A rule file from before the phrases: its source is as it was, so the rules
	// stay, and only the phrases are written.
	files := map[string]string{}
	for name, content := range repoFiles {
		files[name] = content
	}
	source := repoFiles["app/hooks/expenses.go"]
	files["app/helpcorpus/rules/app/hooks/expenses.go.md"] = helpcorpus.RuleFile{
		Source: "app/hooks/expenses.go", SHA256: helpcorpus.SourceHash(source), Model: "old-model",
		Rules: []helpcorpus.Rule{{Title: "One", Text: "First.", Start: 1, End: 1}, {Title: "Two", Text: "Second.", Start: 3, End: 3}},
	}.Format()
	repo := writeRepo(t, files)
	appDir := filepath.Join(repo, "app")
	o := options{repo: repo, model: "claude-sonnet-5-5", effort: "medium", phraseModel: "claude-haiku-5-5", phraseEffort: "low", parallel: 1, only: "app/hooks/expenses.go"}

	// A check reports it and calls no model. With --no-phrases there is nothing to do.
	model := &stubModel{}
	service := testService(t, model)
	check := o
	check.check = true
	if err := run(context.Background(), check, service, io.Discard); err == nil {
		t.Error("check must fail for rules with no phrases")
	}
	without := o
	without.noPhrases = true
	if err := run(context.Background(), without, service, io.Discard); err != nil || len(model.asked)+len(model.phraseAsked) != 0 {
		t.Fatalf("no-phrases: %v, %d requests", err, len(model.asked)+len(model.phraseAsked))
	}

	// A reply that skips a rule is refused, and the file stays as it was.
	model.phraseReply = func(string) string { return `{"asks":[{"n":1,"phrases":["a b"]},{"n":7,"phrases":["c d"]}]}` }
	if err := run(context.Background(), o, service, io.Discard); err == nil || len(model.phraseAsked) != 2 {
		t.Errorf("a rule with no phrases: err = %v, requests = %d", err, len(model.phraseAsked))
	}
	if file, _, _ := helpcorpus.ReadRuleFile(appDir, "app/hooks/expenses.go"); file.Rules[0].Asks != nil {
		t.Error("a bad reply must not change the file")
	}

	model.phraseReply = func(string) string {
		return `Sure: {"asks":[{"n":2,"phrases":["  second   thing ","` + strings.Repeat("x", 90) + `"]},{"n":1,"phrases":["a","b","c","d","e","f","g","h"]}]}`
	}
	if err := run(context.Background(), o, service, io.Discard); err != nil {
		t.Fatal(err)
	}
	file, _, err := helpcorpus.ReadRuleFile(appDir, "app/hooks/expenses.go")
	if err != nil || len(model.asked) != 0 {
		t.Fatalf("%v, rule requests = %d", err, len(model.asked))
	}
	// The rules, their lines, the hash and the model of the rules are as they were.
	if file.Model != "old-model" || file.SHA256 != helpcorpus.SourceHash(source) || file.Rules[1].Text != "Second." || file.Rules[1].Start != 3 {
		t.Errorf("the rules changed: %+v", file)
	}
	// At most 6 phrases, with the spaces cleaned, and no phrase that is too long.
	if len(file.Rules[0].Asks) != 6 || !reflect.DeepEqual(file.Rules[1].Asks, []string{"second thing"}) {
		t.Errorf("asks = %q and %q", file.Rules[0].Asks, file.Rules[1].Asks)
	}
	if err := run(context.Background(), check, service, io.Discard); err != nil {
		t.Errorf("check after the phrases: %v", err)
	}
}
