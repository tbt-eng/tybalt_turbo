package main

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"tybalt/helpchat"
	"tybalt/helpcorpus"
)

// stubModel answers a question after one search, and grades as a judge when it
// gets the judge's prompt.
type stubModel struct {
	mu       sync.Mutex
	asked    []string
	judged   int
	sawGiven bool
	contexts map[string]string
	sawTools []int
}

func (m *stubModel) Warm(context.Context, helpchat.Request) error { return nil }

func (m *stubModel) Start(request helpchat.Request) helpchat.Session {
	return &stubSession{model: m, request: request}
}

type stubSession struct {
	model   *stubModel
	request helpchat.Request
	results int
}

func (s *stubSession) AddToolResults(results []helpchat.ToolResult) { s.results += len(results) }

func (s *stubSession) Next(_ context.Context, onText func(string)) (*helpchat.Reply, error) {
	m := s.model
	m.mu.Lock()
	defer m.mu.Unlock()
	usage := helpchat.Usage{Input: 100, Output: 20, CacheRead: 1000}

	if s.request.Shared == judgePrompt {
		m.judged++
		verdict := `{"points":[{"n":1,"verdict":"stated","why":"says it"}],"wrong_claims":[],"other_error":"","abstained":false}`
		switch {
		case strings.Contains(s.request.Question, "<candidate_answer>\nWRONG"):
			verdict = `{"points":[{"n":1,"verdict":"contradicted","why":"says the opposite"}],"wrong_claims":[],"other_error":"","abstained":false}`
		case strings.Contains(s.request.Question, "<candidate_answer>\nNOTHING"):
			verdict = "Here is my grade:\n" + `{"points":[{"n":1,"verdict":"missing","why":""}],"wrong_claims":[],"other_error":"","abstained":true}`
		}
		return &helpchat.Reply{Text: verdict, StopReason: "end_turn", Usage: usage}, nil
	}

	if s.results == 0 && len(s.request.Tools) > 0 {
		m.asked = append(m.asked, s.request.Question)
		m.sawTools = append(m.sawTools, len(s.request.Tools))
		m.sawGiven = m.sawGiven || strings.Contains(s.request.Context, "# Source for this question")
		if m.contexts == nil {
			m.contexts = map[string]string{}
		}
		m.contexts[s.request.Question] = s.request.Context
		return &helpchat.Reply{
			ToolCalls:  []helpchat.ToolCall{{ID: "t1", Name: "search_turbo", Input: json.RawMessage(`{"query":"runTurn turnResult tool loop"}`)}},
			StopReason: "tool_use",
			Usage:      usage,
		}, nil
	}
	if len(s.request.Tools) == 0 {
		m.asked = append(m.asked, s.request.Question)
		m.sawTools = append(m.sawTools, 0)
	}
	answer := "RIGHT"
	switch {
	case strings.Contains(s.request.Question, "second"):
		answer = "WRONG"
	case strings.Contains(s.request.Question, "third"):
		answer = "NOTHING"
	}
	onText(answer)
	return &helpchat.Reply{Text: answer, StopReason: "end_turn", Usage: usage}, nil
}

func testService(t *testing.T, model helpchat.Model) *helpchat.Service {
	t.Helper()
	corpus, err := helpcorpus.Open(filepath.Join("..", "..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	index, err := helpcorpus.NewIndex(corpus)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { index.Close() })
	return &helpchat.Service{Model: model, Index: index}
}

func writeCases(t *testing.T, dir string) string {
	t.Helper()
	cases := []Case{
		{ID: "a-01", Category: "a", Question: "The first question?", Held: []string{"commit", "time"}, IdealAnswer: "Yes.", KeyPoints: []string{"Yes"},
			Sources: []SourceRange{{Path: "app/helpchat/turn.go", Start: 125, End: 140}}, InUserDocs: "no"},
		{ID: "a-02", Category: "a", Question: "The second question?", IdealAnswer: "No.", KeyPoints: []string{"No"},
			Sources: []SourceRange{{Path: "app/no/such/file.go", Start: 1, End: 2}}, InUserDocs: "no"},
		{ID: "b-01", Category: "b", Question: "The third question?", IdealAnswer: "Maybe.", KeyPoints: []string{"Maybe"},
			Sources: []SourceRange{{Path: "app/helpchat/turn.go", Start: 125, End: 140}}, InUserDocs: "partly"},
	}
	data, _ := json.Marshal(cases)
	path := filepath.Join(dir, "questions.json")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func testOptions(t *testing.T, dir, mode string) options {
	return options{
		questions: writeCases(t, dir),
		claims:    filepath.Join("..", "..", "test_seed_data", "data", "claims.csv"),
		mode:      mode, provider: "anthropic", model: "claude-haiku-5-5", effort: "low",
		judgeModel: "claude-sonnet-5-5", judgeEffort: "medium",
		out: filepath.Join(dir, "run"), parallel: 2,
	}
}

func TestAskGradesEachAnswerAndResumes(t *testing.T) {
	dir := t.TempDir()
	model := &stubModel{}
	service := testService(t, model)
	o := testOptions(t, dir, "ask")
	devnull, _ := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	defer devnull.Close()

	if err := run(context.Background(), o, service, devnull); err != nil {
		t.Fatal(err)
	}
	rows, err := loadRows(filepath.Join(o.out, "results.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]string{"a-01": gradeCorrect, "a-02": gradeWrong, "b-01": gradeAbstained}
	for id, grade := range want {
		if rows[id].Grade != grade {
			t.Errorf("%s: grade %q, want %q", id, rows[id].Grade, grade)
		}
		if rows[id].Searches != 1 || rows[id].Usage.CacheRead != 2000 || rows[id].JudgeUsage.Input != 100 {
			t.Errorf("%s: searches %d, usage %+v, judge %+v", id, rows[id].Searches, rows[id].Usage, rows[id].JudgeUsage)
		}
	}
	if !rows["a-01"].RangeFound || !rows["a-01"].FileFound {
		t.Errorf("a-01: the search for runTurn must return turn.go lines 125-140: %+v", rows["a-01"])
	}
	if rows["a-02"].FileFound || rows["a-02"].RangeFound {
		t.Errorf("a-02: no tool can return a file that is not there")
	}
	if rows["a-02"].WrongReason != "contradicts: No" {
		t.Errorf("a-02: wrong reason %q", rows["a-02"].WrongReason)
	}
	for _, name := range []string{"summary.md", "traces/a-01.json", "traces/b-01.json"} {
		if _, err := os.Stat(filepath.Join(o.out, name)); err != nil {
			t.Error(err)
		}
	}
	text, _ := os.ReadFile(filepath.Join(o.out, "summary.md"))
	for _, part := range []string{"| **All** | 3 | 33% (1) | 0% (0) | 33% (1) | 33% (1) | 33% (1) |", "the tools returned the right lines | 2 |", "claude-haiku-5-5"} {
		if !strings.Contains(string(text), part) {
			t.Errorf("the summary lacks %q:\n%s", part, text)
		}
	}

	// Each question is asked as the person that it names.
	if got := model.contexts["The first question?"]; !strings.Contains(got, "The user holds these claims: commit, time\n") {
		t.Errorf("the first question must be asked as a committer:\n%.300s", got)
	}
	if got := model.contexts["The second question?"]; !strings.Contains(got, "The user holds no claims.") {
		t.Errorf("the second question names no claims:\n%.200s", got)
	}

	// A second run asks nothing, because every question has a grade.
	asked, judged := len(model.asked), model.judged
	if err := run(context.Background(), o, service, devnull); err != nil {
		t.Fatal(err)
	}
	if len(model.asked) != asked || model.judged != judged {
		t.Errorf("the second run asked %d and judged %d more", len(model.asked)-asked, model.judged-judged)
	}

	// regrade asks the judge again and not the chat.
	o.mode = "regrade"
	if err := run(context.Background(), o, service, devnull); err != nil {
		t.Fatal(err)
	}
	if len(model.asked) != asked || model.judged != judged+3 {
		t.Errorf("regrade asked %d and judged %d more", len(model.asked)-asked, model.judged-judged)
	}
	text, _ = os.ReadFile(filepath.Join(o.out, "summary.md"))
	for _, part := range []string{"# Help chat benchmark: ask, claude-haiku-5-5", "the tools returned the right lines | 2 |"} {
		if !strings.Contains(string(text), part) {
			t.Errorf("the summary of a regrade lacks %q:\n%s", part, text)
		}
	}

	// rescore makes the grades again from the saved verdicts, and asks nobody.
	before, _ := loadRows(filepath.Join(o.out, "results.jsonl"))
	asked, judged = len(model.asked), model.judged
	o.mode = "rescore"
	if err := run(context.Background(), o, service, devnull); err != nil {
		t.Fatal(err)
	}
	after, _ := loadRows(filepath.Join(o.out, "results.jsonl"))
	if len(model.asked) != asked || model.judged != judged {
		t.Error("rescore must call no model")
	}
	for id, grade := range want {
		if after[id].Grade != grade || after[id].Stated != before[id].Stated {
			t.Errorf("%s: rescore gave %q", id, after[id].Grade)
		}
	}
}

func TestOracleGivesTheSourceAndNotoolsTakesTheToolsAway(t *testing.T) {
	devnull, _ := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	defer devnull.Close()

	model := &stubModel{}
	o := testOptions(t, t.TempDir(), "oracle")
	o.only = "a-01"
	if err := run(context.Background(), o, testService(t, model), devnull); err != nil {
		t.Fatal(err)
	}
	if !model.sawGiven {
		t.Error("oracle must put the source in the prompt")
	}
	tr, ok := readTrace(filepath.Join(o.out, "traces", "a-01.json"))
	if !ok || !strings.Contains(tr.Given, "func runTurn") {
		t.Errorf("the trace must hold the source that the model got: %q", tr.Given)
	}

	// A source that the corpus does not hold is an error of the run.
	o.only, o.out = "a-02", filepath.Join(t.TempDir(), "run")
	if err := run(context.Background(), o, testService(t, model), devnull); err != nil {
		t.Fatal(err)
	}
	rows, _ := loadRows(filepath.Join(o.out, "results.jsonl"))
	if rows["a-02"].Grade != gradeError {
		t.Errorf("a-02: grade %q, want an error", rows["a-02"].Grade)
	}

	model = &stubModel{}
	o = testOptions(t, t.TempDir(), "notools")
	o.only = "a-01"
	if err := run(context.Background(), o, testService(t, model), devnull); err != nil {
		t.Fatal(err)
	}
	if len(model.sawTools) != 1 || model.sawTools[0] != 0 {
		t.Errorf("notools must send no tools: %v", model.sawTools)
	}
}

func TestVerdictGrade(t *testing.T) {
	c := Case{KeyPoints: []string{"one", "two"}, WrongIf: []string{"bad"}}
	parse := func(text string) *Verdict {
		v, err := parseVerdict(text, c)
		if err != nil {
			t.Fatal(err)
		}
		return v
	}
	tests := []struct {
		verdict string
		grade   string
		stated  int
	}{
		{`{"points":[{"n":1,"verdict":"stated"},{"n":2,"verdict":"stated"}],"wrong_claims":[{"n":1,"made":false}]}`, gradeCorrect, 2},
		{`{"points":[{"n":1,"verdict":"stated"},{"n":2,"verdict":"missing"}],"wrong_claims":[{"n":1,"made":false}]}`, gradePartial, 1},
		{`{"points":[{"n":1,"verdict":"stated"},{"n":2,"verdict":"stated"}],"wrong_claims":[{"n":1,"made":true}]}`, gradeWrong, 2},
		{`{"points":[{"n":1,"verdict":"stated"},{"n":2,"verdict":"contradicted"}],"wrong_claims":[]}`, gradeWrong, 1},
		// A statement that the judge takes for an error is a flag, not a grade.
		{`{"points":[{"n":1,"verdict":"stated"},{"n":2,"verdict":"stated"}],"other_error":"It says 5 and the reference says 6."}`, gradeCorrect, 2},
		{`{"points":[{"n":1,"verdict":"missing"},{"n":2,"verdict":"missing"}],"abstained":true}`, gradeAbstained, 0},
		// An answer that says nothing right and does not say "I don't know".
		{`{"points":[{"n":1,"verdict":"missing"},{"n":2,"verdict":"missing"}],"abstained":false}`, gradePartial, 0},
	}
	for _, test := range tests {
		grade, stated, _ := parse(test.verdict).grade(c)
		if grade != test.grade || stated != test.stated {
			t.Errorf("%s: got %s %d, want %s %d", test.verdict, grade, stated, test.grade, test.stated)
		}
	}
	// Words before and after the object, with a brace in them, do no harm.
	after := "My grade:\n" + `{"points":[{"n":1,"verdict":"stated"},{"n":2,"verdict":"missing"}]}` + "\nCheck: {done}"
	if grade, stated, _ := parse(after).grade(c); grade != gradePartial || stated != 1 {
		t.Errorf("text around the JSON: got %s %d", grade, stated)
	}
	// One verdict too many: the numbers say which ones are for the key points.
	extra := `{"points":[{"n":1,"verdict":"stated"},{"n":2,"verdict":"missing"},{"n":3,"verdict":"contradicted"}]}`
	if grade, stated, _ := parse(extra).grade(c); grade != gradePartial || stated != 1 {
		t.Errorf("one verdict too many: got %s %d", grade, stated)
	}
	if flag := parse(`{"points":[{"n":1,"verdict":"stated"},{"n":2,"verdict":"stated"}],"other_error":" It says 5. "}`).flag(); flag != "It says 5." {
		t.Errorf("flag = %q", flag)
	}
	for _, bad := range []string{"no json", `{"points":[{"n":1,"verdict":"stated"}]}`, `{"points":[{"n":1,"verdict":"yes"},{"n":2,"verdict":"stated"}]}`} {
		if _, err := parseVerdict(bad, c); err == nil {
			t.Errorf("%s: must be an error", bad)
		}
	}
}

func TestFound(t *testing.T) {
	c := Case{Sources: []SourceRange{{Path: "a.go", Start: 10, End: 20}, {Path: "b.go", Start: 5, End: 6}}}
	tests := []struct {
		chunks         []helpchat.EvalChunk
		file, any, all bool
	}{
		{nil, false, false, false},
		{[]helpchat.EvalChunk{{Path: "a.go", Start: 30, End: 40}}, true, false, false},
		{[]helpchat.EvalChunk{{Path: "a.go", Start: 1, End: 10}}, true, true, false},
		{[]helpchat.EvalChunk{{Path: "a.go", Start: 20, End: 25}, {Path: "b.go", Start: 1, End: 9}}, true, true, true},
		{[]helpchat.EvalChunk{{Path: "c.go", Start: 10, End: 20}}, false, false, false},
	}
	for _, test := range tests {
		file, any, all := found(c, test.chunks)
		if file != test.file || any != test.any || all != test.all {
			t.Errorf("%v: got %v %v %v", test.chunks, file, any, all)
		}
	}
}

func TestTheRealQuestionsNameRealClaims(t *testing.T) {
	cases, err := loadCases("questions.json")
	if err != nil {
		t.Fatal(err)
	}
	claims, err := loadClaims(filepath.Join("..", "..", "test_seed_data", "data", "claims.csv"))
	if err != nil {
		t.Fatal(err)
	}
	known := map[string]bool{}
	for _, claim := range claims {
		known[claim.Name] = true
	}
	for _, c := range cases {
		staff := false
		for _, name := range c.Held {
			staff = staff || name == "time"
			if !known[name] {
				t.Errorf("%s: there is no claim %q", c.ID, name)
			}
		}
		if !staff {
			t.Errorf("%s: every question is asked by a member of staff, who holds time", c.ID)
		}
	}
	if len(cases) != 104 {
		t.Errorf("%d questions", len(cases))
	}
}

func TestAClaimThatDoesNotExistStopsTheRun(t *testing.T) {
	devnull, _ := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	defer devnull.Close()
	o := testOptions(t, t.TempDir(), "ask")
	o.held = "no_such_claim"
	model := &stubModel{}
	if err := run(context.Background(), o, testService(t, model), devnull); err == nil || len(model.asked) != 0 {
		t.Errorf("err = %v, asked %d", err, len(model.asked))
	}
}

func TestTheJudgeGetsTheSourceOfTheQuestion(t *testing.T) {
	service := testService(t, &stubModel{})
	c := Case{ID: "x", Question: "Q?", IdealAnswer: "A.", KeyPoints: []string{"One"}, WrongIf: []string{"Bad"},
		Sources: []SourceRange{{Path: "app/helpchat/turn.go", Start: 125, End: 140}, {Path: "app/no/such/file.go", Start: 1, End: 2}}}
	source := sourceOf(service, helpchat.DefaultSettings(), c)
	if !strings.Contains(source, "func runTurn") || !strings.Contains(source, `path="app/helpchat/turn.go"`) {
		t.Fatalf("source = %.200s", source)
	}
	input := judgeInput(c, "The answer.", source)
	for _, part := range []string{"<source>\n<chunk", "func runTurn", "</source>\n\n<candidate_answer>\nThe answer.", `Give exactly 1 entries in "points" and exactly 1 in "wrong_claims".`} {
		if !strings.Contains(input, part) {
			t.Errorf("the judge's input lacks %q", part)
		}
	}
	if !strings.Contains(judgeInput(c, "", ""), "<source>\n(not given)\n</source>") {
		t.Error("with no source, the input says so")
	}
}

func TestJudgeCheckComparesWithKnownGrades(t *testing.T) {
	dir := t.TempDir()
	o := testOptions(t, dir, "judge-check")
	set := []judgeCase{
		{Name: "right", Case: "a-01", Expect: []string{gradeCorrect}, Why: "It is right.", Answer: "RIGHT"},
		{Name: "caught", Case: "a-01", Expect: []string{gradeWrong}, Why: "It is wrong.", Answer: "WRONG"},
		{Name: "missed", Case: "a-01", Expect: []string{gradeWrong, gradePartial}, Why: "The judge must fail this.", Answer: "RIGHT"},
	}
	data, _ := json.Marshal(set)
	o.judgeCases = filepath.Join(dir, "judge_cases.json")
	if err := os.WriteFile(o.judgeCases, data, 0o644); err != nil {
		t.Fatal(err)
	}
	out, err := os.Create(filepath.Join(dir, "out.txt"))
	if err != nil {
		t.Fatal(err)
	}
	defer out.Close()
	if err := run(context.Background(), o, testService(t, &stubModel{}), out); err != nil {
		t.Fatal(err)
	}
	text, _ := os.ReadFile(out.Name())
	for _, part := range []string{"the judge agrees on 2 of 3.", "WRONG  missed: the judge says correct, and it is wrong or partial. The judge must fail this."} {
		if !strings.Contains(string(text), part) {
			t.Errorf("the output lacks %q:\n%s", part, text)
		}
	}

	// A case for a question that is not there is an error.
	set[0].Case = "no-such-question"
	data, _ = json.Marshal(set)
	os.WriteFile(o.judgeCases, data, 0o644)
	if err := run(context.Background(), o, testService(t, &stubModel{}), out); err == nil {
		t.Error("a case for an unknown question must be an error")
	}
}

func TestTheRealJudgeCasesNameRealQuestions(t *testing.T) {
	cases, err := loadCases("questions.json")
	if err != nil {
		t.Fatal(err)
	}
	known := map[string]bool{}
	for _, c := range cases {
		known[c.ID] = true
	}
	data, err := os.ReadFile("judge_cases.json")
	if err != nil {
		t.Fatal(err)
	}
	var set []judgeCase
	if err := json.Unmarshal(data, &set); err != nil {
		t.Fatal(err)
	}
	grades := map[string]bool{gradeCorrect: true, gradePartial: true, gradeWrong: true, gradeAbstained: true}
	for _, jc := range set {
		if !known[jc.Case] || strings.TrimSpace(jc.Answer) == "" || jc.Why == "" || len(jc.Expect) == 0 {
			t.Errorf("%s: needs a real question, an answer, a reason and a grade", jc.Name)
		}
		for _, grade := range jc.Expect {
			if !grades[grade] {
				t.Errorf("%s: %q is not a grade", jc.Name, grade)
			}
		}
	}
	if len(set) < 10 {
		t.Errorf("%d cases", len(set))
	}
}
