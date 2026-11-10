// Command helpeval measures how well Help chat answers hard questions.
//
// It reads questions.json: questions about Turbo's behavior, each with an ideal
// answer, the key points that a correct answer states, and the lines of source
// that settle it. It asks each question the way the app does, and a second
// model (the judge) checks the answer against the key points.
//
// Run it from the app/ directory. Each provider's key is a variable in the
// shell, as for the app (see descriptions/help_chat.md).
//
//	go run ./cmd/helpeval --mode search      no model and no cost: what the search finds
//	go run ./cmd/helpeval --mode replay --out <run>   no cost: the searches of a run, sent again
//	go run ./cmd/helpeval --mode ask         the chat as it is
//	go run ./cmd/helpeval --mode oracle      the chat with the right source given to it
//	go run ./cmd/helpeval --mode notools     the chat with the handbook and no tools
//	go run ./cmd/helpeval --mode judge-check does the judge fail bad answers?
//
// A run writes results.jsonl, summary.md and one trace file for each question to
// cmd/helpeval/runs/<name>. A second run with the same name asks only the
// questions that have no result yet.
package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"sync"
	"time"

	"tybalt/helpchat"
)

// Case is one question of the benchmark.
type Case struct {
	ID       string `json:"id"`
	Category string `json:"category"`
	Question string `json:"question"`
	// Held names the claims of the person who asks. A question that starts
	// with "I'm a committer" is asked as a committer, so the prompt and the
	// question agree. Every question has "time", as all staff do.
	Held        []string `json:"held"`
	IdealAnswer string   `json:"ideal_answer"`
	KeyPoints   []string `json:"key_points"`
	// ExtraPoints are true and in the ideal answer, but an answer can be correct
	// without them. Nothing grades them.
	ExtraPoints []string      `json:"extra_points,omitempty"`
	WrongIf     []string      `json:"wrong_if"`
	Sources     []SourceRange `json:"sources"`
	// InUserDocs says if the user docs hold the answer: no, partly or yes.
	InUserDocs string `json:"in_user_docs"`
	WhyHard    string `json:"why_hard,omitempty"`
}

// SourceRange is the lines of one file that settle an answer.
type SourceRange struct {
	Path  string `json:"path"`
	Start int    `json:"start"`
	End   int    `json:"end"`
	Why   string `json:"why,omitempty"`
}

// The grades of an answer.
const (
	gradeCorrect   = "correct"   // every key point, and nothing wrong
	gradePartial   = "partial"   // some key points, and nothing wrong
	gradeWrong     = "wrong"     // it says something that the reference contradicts
	gradeAbstained = "abstained" // it says that it does not know
	gradeError     = "error"     // the run failed: not a result of the chat
)

// Row is the result of one question in one run.
type Row struct {
	ID         string `json:"id"`
	Category   string `json:"category"`
	InUserDocs string `json:"in_user_docs"`
	Mode       string `json:"mode"`
	Model      string `json:"model"`
	Effort     string `json:"effort"`

	Status string `json:"status"`
	Grade  string `json:"grade,omitempty"`
	// Stated is how many of the key points the answer states.
	Stated      int    `json:"stated"`
	Points      int    `json:"points"`
	WrongReason string `json:"wrong_reason,omitempty"`
	// Flag is a statement that the judge takes for an error and that is not a
	// key point or a listed wrong claim. It does not change the grade.
	Flag    string   `json:"flag,omitempty"`
	Verdict *Verdict `json:"verdict,omitempty"`

	// What the model's tool calls returned.
	Searches   int  `json:"searches"`
	Reads      int  `json:"reads"`
	FileFound  bool `json:"file_found"`
	RangeFound bool `json:"range_found"`
	AllRanges  bool `json:"all_ranges_found"`

	Usage      helpchat.Usage `json:"usage"`
	JudgeUsage helpchat.Usage `json:"judge_usage"`
	MSFirst    int64          `json:"ms_first_token"`
	MSTotal    int64          `json:"ms_total"`
	Error      string         `json:"error,omitempty"`
}

// trace is everything about one question, for a person to read.
type trace struct {
	Case   Case                 `json:"case"`
	Row    Row                  `json:"row"`
	Given  string               `json:"given,omitempty"`
	Result *helpchat.EvalResult `json:"result,omitempty"`
}

type options struct {
	questions   string
	claims      string
	mode        string
	provider    string
	model       string
	effort      string
	judgeModel  string
	judgeEffort string
	noJudge     bool
	out         string
	only        string
	limit       int
	parallel    int
	force       bool
	held        string
	judgeCases  string
}

func main() {
	var o options
	defaults := helpchat.DefaultSettings()
	flag.StringVar(&o.questions, "questions", filepath.Join("cmd", "helpeval", "questions.json"), "the questions file")
	flag.StringVar(&o.claims, "claims", filepath.Join("test_seed_data", "data", "claims.csv"), "a CSV of the claims that the prompt lists (columns name and description)")
	flag.StringVar(&o.mode, "mode", "search", "search, replay, ask, oracle, notools, regrade, rescore or judge-check")
	flag.StringVar(&o.provider, "provider", defaults.Provider, "the provider of the model that answers")
	flag.StringVar(&o.model, "model", defaults.Model, "the model that answers")
	flag.StringVar(&o.effort, "effort", defaults.Effort, "the effort of the model that answers")
	flag.StringVar(&o.judgeModel, "judge-model", "claude-sonnet-5-5", "the model that grades (the anthropic provider)")
	flag.StringVar(&o.judgeEffort, "judge-effort", "medium", "the effort of the model that grades")
	flag.StringVar(&o.judgeCases, "judge-cases", filepath.Join("cmd", "helpeval", "judge_cases.json"), "answers with a known grade, for judge-check")
	flag.BoolVar(&o.noJudge, "no-judge", false, "ask the questions and do not grade the answers")
	flag.StringVar(&o.out, "out", "", "the folder for the results (default cmd/helpeval/runs/<mode>-<model>-<effort>)")
	flag.StringVar(&o.only, "only", "", "ids to run, with commas between them")
	flag.IntVar(&o.limit, "limit", 0, "run at most this many questions")
	flag.IntVar(&o.parallel, "parallel", 4, "questions to ask at the same time")
	flag.BoolVar(&o.force, "force", false, "ask again the questions that have a result")
	flag.StringVar(&o.held, "held", "", "ask every question as a person with these claims, with commas between them, and not with the claims that each question names")
	flag.Parse()

	if err := run(context.Background(), o, helpchat.NewService(), os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "helpeval:", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, o options, service *helpchat.Service, out *os.File) error {
	cases, err := loadCases(o.questions)
	if err != nil {
		return err
	}
	all := cases
	if o.mode == "judge-check" && o.limit > 0 {
		// A sample from every area, not the first questions of the file.
		cases, o.limit = spread(pick(cases, o.only, 0), o.limit), 0
	}
	cases = pick(cases, o.only, o.limit)
	if len(cases) == 0 {
		return errors.New("no questions to run")
	}

	settings := helpchat.DefaultSettings()
	settings.Provider, settings.Model, settings.Effort = o.provider, o.model, o.effort
	judge := helpchat.DefaultSettings()
	judge.Model, judge.Effort = o.judgeModel, o.judgeEffort
	judge.MaxOutputTokens, judge.TurnTimeoutSeconds = 4000, 120

	if o.mode == "search" {
		return searchMode(service, cases, "", out)
	}
	if o.mode == "replay" {
		if o.out == "" {
			return errors.New("replay needs --out, the folder of the run whose searches to send again")
		}
		return searchMode(service, cases, o.out, out)
	}
	if o.mode == "rescore" {
		if o.out == "" {
			return errors.New("rescore needs --out, the folder of the run")
		}
		return rescore(o, cases, out)
	}
	// A run that calls a model stops here if a key is missing, before any work.
	if _, err := service.CorpusHash(settings); err != nil && o.mode != "regrade" && o.mode != "judge-check" {
		return keyHint(err)
	}
	if !o.noJudge {
		if _, err := service.CorpusHash(judge); err != nil {
			return keyHint(err)
		}
	}
	if o.mode == "judge-check" {
		// The report also goes to a file, so it is there after the terminal is gone.
		var report strings.Builder
		both := io.MultiWriter(out, &report)
		if err := judgeCheck(ctx, service, judge, cases, o.parallel, both); err != nil {
			return err
		}
		if err := checkJudgeCases(ctx, service, judge, o.judgeCases, all, o.parallel, both); err != nil {
			return err
		}
		dir := filepath.Join("cmd", "helpeval", "runs")
		if o.out != "" {
			dir = o.out
		}
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
		return os.WriteFile(filepath.Join(dir, "judge-check.txt"), []byte(report.String()), 0o644)
	}
	if o.mode != "ask" && o.mode != "oracle" && o.mode != "notools" && o.mode != "regrade" {
		return fmt.Errorf("there is no mode %q", o.mode)
	}

	dir := o.out
	if dir == "" {
		if o.mode == "regrade" {
			return errors.New("regrade needs --out, the folder of the run to grade again")
		}
		dir = filepath.Join("cmd", "helpeval", "runs", fmt.Sprintf("%s-%s-%s", o.mode, o.model, o.effort))
	}
	if err := os.MkdirAll(filepath.Join(dir, "traces"), 0o755); err != nil {
		return err
	}
	rows, err := loadRows(filepath.Join(dir, "results.jsonl"))
	if err != nil {
		return err
	}
	claims, err := loadClaims(o.claims)
	if err != nil {
		return err
	}
	var held []string
	for _, name := range strings.Split(o.held, ",") {
		if name = strings.TrimSpace(name); name != "" {
			held = append(held, name)
		}
	}
	corpus, err := service.CorpusHash(settings)
	if err != nil {
		return err
	}
	// A claim that Turbo does not have would make a prompt that no user gets.
	known := map[string]bool{}
	for _, claim := range claims {
		known[claim.Name] = true
	}
	for _, c := range cases {
		for _, name := range append(slices.Clone(c.Held), held...) {
			if !known[name] {
				return fmt.Errorf("%s: there is no claim %q in %s", c.ID, name, o.claims)
			}
		}
	}

	var mu sync.Mutex
	save := func(c Case, row Row, given string, result *helpchat.EvalResult) {
		mu.Lock()
		defer mu.Unlock()
		rows[row.ID] = row
		if err := writeRows(filepath.Join(dir, "results.jsonl"), rows); err != nil {
			fmt.Fprintln(os.Stderr, "helpeval: can't save the results:", err)
		}
		if result != nil {
			data, _ := json.MarshalIndent(trace{Case: c, Row: row, Given: given, Result: result}, "", "  ")
			if err := os.WriteFile(filepath.Join(dir, "traces", row.ID+".json"), data, 0o644); err != nil {
				fmt.Fprintln(os.Stderr, "helpeval: can't save a trace:", err)
			}
		}
		fmt.Fprintf(out, "%-16s %-9s %d/%d  searches %d  reads %d  source %-5v %5d ms  %s\n",
			row.ID, row.Grade, row.Stated, row.Points, row.Searches, row.Reads, row.RangeFound, row.MSTotal, row.Error)
	}

	one := func(c Case) {
		row := Row{ID: c.ID, Category: c.Category, InUserDocs: c.InUserDocs, Mode: o.mode, Model: o.model, Effort: o.effort, Points: len(c.KeyPoints)}
		var result *helpchat.EvalResult
		given := ""

		if o.mode == "regrade" {
			old, ok := readTrace(filepath.Join(dir, "traces", c.ID+".json"))
			if !ok {
				return
			}
			row, result, given = old.Row, old.Result, old.Given
			row.Points = len(c.KeyPoints)
		} else {
			// --held asks every question as one kind of person.
			asker := c.Held
			if len(held) > 0 {
				asker = held
			}
			request := helpchat.EvalRequest{Settings: settings, AllClaims: claims, Held: asker, Question: c.Question, NoTools: o.mode == "notools"}
			if o.mode == "oracle" {
				for _, source := range c.Sources {
					text, err := service.SourceText(settings, source.Path, source.Start, source.End)
					if err != nil {
						row.Grade, row.Status, row.Error = gradeError, "error", err.Error()
						save(c, row, "", nil)
						return
					}
					given += text
				}
				request.Given = given
			}
			var got helpchat.EvalResult
			// One more try for a failure of the API. It is not a result of the chat.
			for try := 0; try < 2; try++ {
				got, err = service.Evaluate(ctx, request)
				if err != nil {
					row.Grade, row.Status, row.Error = gradeError, "error", err.Error()
					save(c, row, given, nil)
					return
				}
				if got.Status != helpchat.StatusError {
					break
				}
			}
			result = &got
			row.Status, row.Error = got.Status, got.Error
			row.Usage, row.MSFirst, row.MSTotal = got.Usage, got.FirstTokenMS, got.TotalMS
			scoreRetrieval(&row, c, got.ToolCalls)
		}

		switch {
		case row.Status != helpchat.StatusAnswered && row.Status != helpchat.StatusGap:
			row.Grade = gradeError
		case o.noJudge:
			row.Grade = ""
		default:
			verdict, usage, err := grade(ctx, service, judge, c, result.Answer)
			row.JudgeUsage = usage
			if err != nil {
				row.Grade, row.Error = gradeError, "the judge failed: "+err.Error()
			} else {
				row.Verdict = verdict
				row.Grade, row.Stated, row.WrongReason = verdict.grade(c)
				row.Flag = verdict.flag()
			}
		}
		save(c, row, given, result)
	}

	var todo []Case
	for _, c := range cases {
		old, done := rows[c.ID]
		if o.mode == "regrade" || o.force || !done || old.Grade == gradeError || (old.Grade == "" && !o.noJudge) {
			todo = append(todo, c)
		}
	}
	fmt.Fprintf(out, "%s: %d questions to run, %d have a result. Corpus %.12s. Results in %s\n\n", o.mode, len(todo), len(cases)-len(todo), corpus, dir)
	if len(todo) > 0 {
		// The first question writes the model's cache. The rest then read it.
		one(todo[0])
		each(todo[1:], o.parallel, one)
	}

	var final []Row
	for _, c := range cases {
		if row, ok := rows[c.ID]; ok {
			final = append(final, row)
		}
	}
	text := summary(o, corpus, final)
	fmt.Fprintln(out, "\n"+text)
	return os.WriteFile(filepath.Join(dir, "summary.md"), []byte(text), 0o644)
}

// rescore makes the grades of a run again from the verdicts that the judge
// gave, with no model. Use it after a change to how a verdict becomes a grade.
// An answer whose question has other key points now keeps its grade, and needs
// regrade.
func rescore(o options, cases []Case, out *os.File) error {
	path := filepath.Join(o.out, "results.jsonl")
	rows, err := loadRows(path)
	if err != nil {
		return err
	}
	var final []Row
	var stale []string
	for _, c := range cases {
		row, ok := rows[c.ID]
		if !ok {
			continue
		}
		switch {
		case row.Verdict == nil:
		case len(row.Verdict.Points) != len(c.KeyPoints):
			stale = append(stale, c.ID)
		default:
			row.Points = len(c.KeyPoints)
			row.Grade, row.Stated, row.WrongReason = row.Verdict.grade(c)
			row.Flag = row.Verdict.flag()
			rows[c.ID] = row
		}
		final = append(final, rows[c.ID])
	}
	if err := writeRows(path, rows); err != nil {
		return err
	}
	text := summary(o, "", final)
	if len(stale) > 0 {
		text += fmt.Sprintf("\n%d answers keep their old grade, because their questions have other key points now. Grade them again with --mode regrade: %s\n", len(stale), strings.Join(stale, ", "))
	}
	fmt.Fprintln(out, text)
	return os.WriteFile(filepath.Join(o.out, "summary.md"), []byte(text), 0o644)
}

// each runs work for every case, a few at a time.
func each(cases []Case, parallel int, work func(Case)) {
	if parallel < 1 {
		parallel = 1
	}
	queue := make(chan Case)
	var wait sync.WaitGroup
	for i := 0; i < parallel; i++ {
		wait.Add(1)
		go func() {
			defer wait.Done()
			for c := range queue {
				work(c)
			}
		}()
	}
	for _, c := range cases {
		queue <- c
	}
	close(queue)
	wait.Wait()
}

func loadCases(path string) ([]Case, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cases []Case
	if err := json.Unmarshal(data, &cases); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	seen := map[string]bool{}
	for _, c := range cases {
		switch {
		case c.ID == "" || seen[c.ID]:
			return nil, fmt.Errorf("%s: the id %q is empty or is there twice", path, c.ID)
		case strings.TrimSpace(c.Question) == "" || strings.TrimSpace(c.IdealAnswer) == "":
			return nil, fmt.Errorf("%s: %s needs a question and an ideal answer", path, c.ID)
		case len(c.KeyPoints) == 0 || len(c.Sources) == 0:
			return nil, fmt.Errorf("%s: %s needs key points and sources", path, c.ID)
		}
		seen[c.ID] = true
	}
	return cases, nil
}

// keyHint adds what to do to the error for a missing key.
func keyHint(err error) error {
	if strings.Contains(err.Error(), "is not set") {
		return fmt.Errorf("%w. Set the key in this shell and run the command again in the same shell", err)
	}
	return err
}

// spread returns n cases at even steps through the list.
func spread(cases []Case, n int) []Case {
	if n >= len(cases) {
		return cases
	}
	out := make([]Case, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, cases[i*len(cases)/n])
	}
	return out
}

func pick(cases []Case, only string, limit int) []Case {
	if only != "" {
		want := map[string]bool{}
		for _, id := range strings.Split(only, ",") {
			want[strings.TrimSpace(id)] = true
		}
		var kept []Case
		for _, c := range cases {
			if want[c.ID] {
				kept = append(kept, c)
			}
		}
		cases = kept
	}
	if limit > 0 && len(cases) > limit {
		cases = cases[:limit]
	}
	return cases
}

// loadClaims reads the names and descriptions of the claims from a CSV, such as
// the one in the test fixtures. The app reads them from its database.
func loadClaims(path string) ([]helpchat.Claim, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	if len(records) == 0 {
		return nil, fmt.Errorf("%s is empty", path)
	}
	name, description := -1, -1
	for i, column := range records[0] {
		switch column {
		case "name":
			name = i
		case "description":
			description = i
		}
	}
	if name < 0 || description < 0 {
		return nil, fmt.Errorf("%s needs the columns name and description", path)
	}
	var claims []helpchat.Claim
	for _, record := range records[1:] {
		claims = append(claims, helpchat.Claim{Name: record[name], Description: record[description]})
	}
	sort.Slice(claims, func(i, j int) bool { return claims[i].Name < claims[j].Name })
	return claims, nil
}

func loadRows(path string) (map[string]Row, error) {
	rows := map[string]Row{}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return rows, nil
	}
	if err != nil {
		return nil, err
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var row Row
		if err := json.Unmarshal([]byte(line), &row); err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		rows[row.ID] = row
	}
	return rows, nil
}

func writeRows(path string, rows map[string]Row) error {
	ids := make([]string, 0, len(rows))
	for id := range rows {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	var out strings.Builder
	for _, id := range ids {
		data, err := json.Marshal(rows[id])
		if err != nil {
			return err
		}
		out.Write(data)
		out.WriteByte('\n')
	}
	return os.WriteFile(path, []byte(out.String()), 0o644)
}

func readTrace(path string) (trace, bool) {
	var t trace
	data, err := os.ReadFile(path)
	if err != nil || json.Unmarshal(data, &t) != nil || t.Result == nil {
		return t, false
	}
	return t, true
}

// overlaps reports if a chunk holds any line of a source range.
func overlaps(chunk helpchat.EvalChunk, source SourceRange) bool {
	return chunk.Path == source.Path && chunk.Start <= source.End && chunk.End >= source.Start
}

// found compares the chunks that the tools returned with the lines that settle
// the answer: the right file, any of the right lines, and all of them.
func found(c Case, chunks []helpchat.EvalChunk) (file, any, all bool) {
	all = true
	for _, source := range c.Sources {
		hit := false
		for _, chunk := range chunks {
			if chunk.Path == source.Path {
				file = true
			}
			if overlaps(chunk, source) {
				hit = true
			}
		}
		any = any || hit
		all = all && hit
	}
	return file, any, all
}

func scoreRetrieval(row *Row, c Case, calls []helpchat.EvalToolCall) {
	var chunks []helpchat.EvalChunk
	for _, call := range calls {
		switch call.Tool {
		case "search_turbo":
			row.Searches++
		case "read_source":
			row.Reads++
		}
		chunks = append(chunks, call.Chunks...)
	}
	row.FileFound, row.RangeFound, row.AllRanges = found(c, chunks)
}

// searchMode measures the search with no model. It costs nothing and gives the
// same result each time, so it is the fast check of a change to the search or
// to the chunks. With no run, each question is one search in the user's own
// words. With a run, it sends again the searches that the model wrote in that
// run, which shows what the model would get from the search as it is now.
func searchMode(service *helpchat.Service, cases []Case, run string, out *os.File) error {
	type tally struct{ n, file, any, all int }
	total := tally{}
	by := map[string]*tally{}
	var missed, none []string
	for _, c := range cases {
		queries := []string{c.Question}
		if run != "" {
			queries = nil
			old, ok := readTrace(filepath.Join(run, "traces", c.ID+".json"))
			if !ok {
				continue
			}
			for _, call := range old.Result.ToolCalls {
				var input struct {
					Query string `json:"query"`
				}
				if call.Tool == "search_turbo" && json.Unmarshal(call.Input, &input) == nil && input.Query != "" {
					queries = append(queries, input.Query)
				}
			}
			if len(queries) == 0 {
				none = append(none, c.ID)
			}
		}
		var chunks []helpchat.EvalChunk
		for _, query := range queries {
			hits, err := service.SearchOnly(query)
			if err != nil {
				return err
			}
			chunks = append(chunks, hits...)
		}
		file, any, all := found(c, chunks)
		for _, t := range []*tally{&total, at(by, c.Category, &tally{})} {
			t.n++
			t.file += count(file)
			t.any += count(any)
			t.all += count(all)
		}
		if !any && len(queries) > 0 {
			missed = append(missed, fmt.Sprintf("%s (%s)", c.ID, c.Sources[0].Path))
		}
	}
	if run == "" {
		fmt.Fprint(out, "Search only: the question itself as one search, no model.\n\n")
	} else {
		fmt.Fprintf(out, "Replay: the searches that the model wrote in %s, sent to the search as it is now.\n\n", run)
	}
	fmt.Fprintln(out, "| Area | Questions | Right file | Right lines | All the right lines |")
	fmt.Fprintln(out, "| --- | ---: | ---: | ---: | ---: |")
	for _, name := range keys(by) {
		t := by[name]
		fmt.Fprintf(out, "| %s | %d | %s | %s | %s |\n", name, t.n, pct(t.file, t.n), pct(t.any, t.n), pct(t.all, t.n))
	}
	fmt.Fprintf(out, "| **All** | %d | %s | %s | %s |\n", total.n, pct(total.file, total.n), pct(total.any, total.n), pct(total.all, total.n))
	fmt.Fprintf(out, "\nNot found (%d): %s\n", len(missed), strings.Join(missed, ", "))
	if len(none) > 0 {
		fmt.Fprintf(out, "The model did not search (%d): %s\n", len(none), strings.Join(none, ", "))
	}
	return nil
}

func at[T any](m map[string]*T, key string, zero *T) *T {
	if m[key] == nil {
		m[key] = zero
	}
	return m[key]
}

func keys[T any](m map[string]T) []string {
	names := make([]string, 0, len(m))
	for name := range m {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func count(yes bool) int {
	if yes {
		return 1
	}
	return 0
}

func pct(part, whole int) string {
	if whole == 0 {
		return "-"
	}
	return fmt.Sprintf("%d%% (%d)", (part*100+whole/2)/whole, part)
}

// prices are US dollars for one million tokens: input, output, cache read and
// cache write. They are only for the estimate in the summary.
var prices = map[string][4]float64{
	"claude-haiku-5-5":  {0.10, 0.50, 0.01, 0.125},
	"claude-sonnet-5-5": {2.00, 10.00, 0.10, 2.50},
}

func cost(model string, usage helpchat.Usage) (float64, bool) {
	p, ok := prices[model]
	if !ok {
		return 0, false
	}
	return (float64(usage.Input)*p[0] + float64(usage.Output)*p[1] + float64(usage.CacheRead)*p[2] + float64(usage.CacheWrite)*p[3]) / 1e6, true
}

func median(values []int64) int64 {
	if len(values) == 0 {
		return 0
	}
	sort.Slice(values, func(i, j int) bool { return values[i] < values[j] })
	return values[len(values)/2]
}

// summary is the report of a run. Every number comes from the rows.
func summary(o options, corpus string, rows []Row) string {
	type tally struct{ n, correct, partial, wrong, abstained, stated, points int }
	add := func(t *tally, row Row) {
		t.n++
		t.stated += row.Stated
		t.points += row.Points
		switch row.Grade {
		case gradeCorrect:
			t.correct++
		case gradePartial:
			t.partial++
		case gradeWrong:
			t.wrong++
		case gradeAbstained:
			t.abstained++
		}
	}
	line := func(name string, t *tally) string {
		return fmt.Sprintf("| %s | %d | %s | %s | %s | %s | %s |\n", name, t.n,
			pct(t.correct, t.n), pct(t.partial, t.n), pct(t.wrong, t.n), pct(t.abstained, t.n), pct(t.stated, t.points))
	}

	total := &tally{}
	byArea, byDocs, bySource := map[string]*tally{}, map[string]*tally{}, map[string]*tally{}
	var errorsN, searched, fileFound, rangeFound, allRanges, toolCalls int
	var usage, judgeUsage helpchat.Usage
	var first, whole []int64
	var failed []string
	for _, row := range rows {
		if row.Grade == gradeError || row.Grade == "" {
			errorsN++
			failed = append(failed, row.ID)
			continue
		}
		add(total, row)
		add(at(byArea, row.Category, &tally{}), row)
		add(at(byDocs, row.InUserDocs, &tally{}), row)
		source := "the tools did not return the right lines"
		if row.RangeFound {
			source = "the tools returned the right lines"
		}
		add(at(bySource, source, &tally{}), row)
		searched += count(row.Searches+row.Reads > 0)
		fileFound += count(row.FileFound)
		rangeFound += count(row.RangeFound)
		allRanges += count(row.AllRanges)
		toolCalls += row.Searches + row.Reads
		usage.Input += row.Usage.Input
		usage.Output += row.Usage.Output
		usage.CacheRead += row.Usage.CacheRead
		usage.CacheWrite += row.Usage.CacheWrite
		judgeUsage.Input += row.JudgeUsage.Input
		judgeUsage.Output += row.JudgeUsage.Output
		first = append(first, row.MSFirst)
		whole = append(whole, row.MSTotal)
	}

	// A regrade reports the run that it grades again.
	if len(rows) > 0 {
		o.mode, o.model, o.effort = rows[0].Mode, rows[0].Model, rows[0].Effort
	}
	var out strings.Builder
	fmt.Fprintf(&out, "# Help chat benchmark: %s, %s, effort %s\n\n", o.mode, o.model, o.effort)
	fmt.Fprintf(&out, "Run on %s. Corpus %.12s. Judge %s, effort %s.\n\n", time.Now().Format("2006-01-02"), corpus, o.judgeModel, o.judgeEffort)
	if total.n == 0 {
		fmt.Fprintf(&out, "No question has a grade. Failed: %s\n", strings.Join(failed, ", "))
		return out.String()
	}
	out.WriteString("Correct: every key point and nothing wrong. Partial: some key points and nothing wrong. " +
		"Wrong: it contradicts a key point, or makes a wrong claim that the question lists. Abstained: the chat said that it does not know.\n\n")
	out.WriteString("| | Questions | Correct | Partial | Wrong | Abstained | Key points stated |\n| --- | ---: | ---: | ---: | ---: | ---: | ---: |\n")
	out.WriteString(line("**All**", total))
	for _, name := range keys(byArea) {
		out.WriteString(line(name, byArea[name]))
	}
	out.WriteString("\nBy whether the user docs hold the answer:\n\n| In the user docs | Questions | Correct | Partial | Wrong | Abstained | Key points stated |\n| --- | ---: | ---: | ---: | ---: | ---: | ---: |\n")
	for _, name := range keys(byDocs) {
		out.WriteString(line(name, byDocs[name]))
	}
	if o.mode == "ask" {
		out.WriteString("\nBy what the model's searches and reads returned:\n\n| | Questions | Correct | Partial | Wrong | Abstained | Key points stated |\n| --- | ---: | ---: | ---: | ---: | ---: | ---: |\n")
		for _, name := range keys(bySource) {
			out.WriteString(line(name, bySource[name]))
		}
		fmt.Fprintf(&out, "\nThe model used a tool in %s of the questions, %.1f calls for each question. "+
			"The tools returned the right file in %s, some of the right lines in %s, and all of them in %s.\n",
			pct(searched, total.n), float64(toolCalls)/float64(total.n), pct(fileFound, total.n), pct(rangeFound, total.n), pct(allRanges, total.n))
	}
	fmt.Fprintf(&out, "\nTime: first word after %d ms and the whole answer after %d ms (medians).\n", median(first), median(whole))
	fmt.Fprintf(&out, "Tokens for the answers: %d input, %d from the cache, %d written to the cache, %d output.", usage.Input, usage.CacheRead, usage.CacheWrite, usage.Output)
	if dollars, ok := cost(o.model, usage); ok {
		fmt.Fprintf(&out, " About $%.3f.", dollars)
	}
	fmt.Fprintf(&out, "\nTokens for the judge: %d input, %d output.", judgeUsage.Input, judgeUsage.Output)
	if dollars, ok := cost(o.judgeModel, judgeUsage); ok {
		fmt.Fprintf(&out, " About $%.3f.", dollars)
	}
	out.WriteString("\n")
	var flagged []string
	for _, row := range rows {
		if row.Flag != "" && row.Grade != gradeError && row.Grade != "" {
			flagged = append(flagged, fmt.Sprintf("- %s (%s): %s", row.ID, row.Grade, row.Flag))
		}
	}
	if len(flagged) > 0 {
		fmt.Fprintf(&out, "\n%d answers state something else that the judge takes for an error. This does not change their grade, because the judge is often wrong about it: it reads only a few lines of the source. Read each one against the code:\n%s\n",
			len(flagged), strings.Join(flagged, "\n"))
	}
	if errorsN > 0 {
		fmt.Fprintf(&out, "\n%d questions have no grade, because the run failed for them, and are not in the numbers: %s\n", errorsN, strings.Join(failed, ", "))
	}
	fmt.Fprintf(&out, "\nWith %d questions asked one time each, a difference of less than about %d points between two runs can be chance.\n",
		total.n, noise(total.n))
	return out.String()
}

// noise is about two standard errors of a rate near one half, in points.
func noise(n int) int {
	for points := 1; ; points++ {
		if points*points*n >= 10000 {
			return points
		}
	}
}
