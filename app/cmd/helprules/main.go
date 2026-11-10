// Command helprules writes the rule files that Help chat searches.
//
// A rule file says in plain words what one source file does. A model reads the
// file and writes the rules. The rule files are in the repo (helpcorpus/rules),
// so the app's build needs no model. See helpcorpus/rules.go.
//
// Run it from the app/ directory, with the provider's key in the shell:
//
//	go run ./cmd/helprules --check    no model: list the files with no rules, or with old rules
//	go run ./cmd/helprules            write the rules of those files, then the search phrases of the rules
//	go run ./cmd/helprules --only app/hooks/expenses.go --force
//
// A file gets new rules only when its source changed, so a second run costs
// nothing. Commit the rule files with the change to the source.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"tybalt/helpchat"
	"tybalt/helpcorpus"
)

type options struct {
	repo, appDir  string
	only          string
	limit         int
	force, check  bool
	model, effort string
	parallel      int
	// The second step writes the search phrases of each rule.
	phraseModel, phraseEffort string
	noPhrases                 bool
}

func main() {
	var o options
	flag.StringVar(&o.repo, "repo", "..", "the folder that holds ui/src and descriptions")
	flag.StringVar(&o.appDir, "app-dir", "", "the Go app folder, if it is not <repo>/app")
	flag.StringVar(&o.only, "only", "", "write rules only for source paths that hold one of these texts, with commas between them")
	flag.IntVar(&o.limit, "limit", 0, "write the rules of at most this many files")
	flag.BoolVar(&o.force, "force", false, "write the rules again for files that have them")
	flag.BoolVar(&o.check, "check", false, "call no model: list the files whose rules are missing or old, and fail if there are any")
	flag.StringVar(&o.model, "model", "claude-sonnet-5-5", "the model that writes the rules (the anthropic provider)")
	flag.StringVar(&o.effort, "effort", "medium", "the effort of that model")
	flag.IntVar(&o.parallel, "parallel", 4, "files to work on at the same time")
	flag.StringVar(&o.phraseModel, "phrase-model", "claude-haiku-5-5", "the model that writes the search phrases of the rules (the anthropic provider)")
	flag.StringVar(&o.phraseEffort, "phrase-effort", "low", "the effort of that model")
	flag.BoolVar(&o.noPhrases, "no-phrases", false, "write the rules and not their search phrases")
	flag.Parse()

	if err := run(context.Background(), o, helpchat.NewService(), os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "helprules:", err)
		os.Exit(1)
	}
}

// job is one source file and the state of its rule file.
type job struct {
	source helpcorpus.SourceFile
	state  string // "new", "old" or "ok"
}

func run(ctx context.Context, o options, service *helpchat.Service, out io.Writer) error {
	corpusOptions := helpcorpus.Options{RepoDir: o.repo, AppDir: o.appDir}
	appDir := o.appDir
	if appDir == "" {
		appDir = filepath.Join(o.repo, "app")
	}
	sources, err := helpcorpus.RuleSources(corpusOptions)
	if err != nil {
		return err
	}

	var jobs []job
	counts := map[string]int{}
	bytes := 0
	for _, source := range sources {
		state := "new"
		existing, ok, err := helpcorpus.ReadRuleFile(appDir, source.Path)
		switch {
		case err != nil:
			fmt.Fprintln(out, "can't read the rules, so they will be written again:", err)
		case ok && existing.SHA256 == helpcorpus.SourceHash(source.Content):
			state = "ok"
		case ok:
			state = "old"
		}
		counts[state]++
		if (state != "ok" || o.force) && selected(source.Path, o.only) {
			jobs = append(jobs, job{source: source, state: state})
			bytes += len(source.Content)
		}
	}
	if o.limit > 0 && len(jobs) > o.limit {
		jobs = jobs[:o.limit]
	}
	extra, err := orphans(appDir, sources)
	if err != nil {
		return err
	}
	fmt.Fprintf(out, "%d source files get rules: %d have them, %d have rules for an older source, %d have none.\n",
		len(sources), counts["ok"], counts["old"], counts["new"])

	if o.check {
		for _, j := range jobs {
			fmt.Fprintf(out, "  %-4s %s\n", j.state, j.source.Path)
		}
		for _, name := range extra {
			fmt.Fprintf(out, "  gone %s\n", name)
		}
		var phrases []job
		if !o.noPhrases {
			phrases = needPhrases(appDir, sources, o.only, out)
			for _, j := range phrases {
				fmt.Fprintf(out, "  asks %s\n", j.source.Path)
			}
		}
		if len(jobs)+len(extra)+len(phrases) > 0 {
			return fmt.Errorf("%d rule files to write, %d to remove and %d with no search phrases: run go run ./cmd/helprules", len(jobs), len(extra), len(phrases))
		}
		return nil
	}

	for _, name := range extra {
		if err := os.Remove(name); err != nil {
			return err
		}
		fmt.Fprintln(out, "removed", name, "(its source is gone)")
	}
	var failed []string
	if len(jobs) > 0 {
		if failed, err = writeAll(ctx, o, service, appDir, jobs, bytes, out); err != nil {
			return err
		}
	}
	if !o.noPhrases {
		// Read the rule files again: the step above wrote some of them.
		if phrases := needPhrases(appDir, sources, o.only, out); len(phrases) > 0 {
			more, err := writeAllPhrases(ctx, o, service, appDir, phrases, out)
			if err != nil {
				return err
			}
			failed = append(failed, more...)
		} else if len(jobs) == 0 {
			fmt.Fprintln(out, "Nothing to write.")
		}
	} else if len(jobs) == 0 {
		fmt.Fprintln(out, "Nothing to write.")
	}
	if len(failed) > 0 {
		sort.Strings(failed)
		return fmt.Errorf("%d files are not done, run the command again for them: %s", len(failed), strings.Join(failed, ", "))
	}
	return nil
}

// writeAll writes the rules of the jobs, and returns the paths that failed.
func writeAll(ctx context.Context, o options, service *helpchat.Service, appDir string, jobs []job, bytes int, out io.Writer) ([]string, error) {
	settings := helpchat.DefaultSettings()
	settings.Model, settings.Effort = o.model, o.effort
	settings.MaxOutputTokens, settings.TurnTimeoutSeconds = 32000, 300
	system, err := systemPrompt(o.repo)
	if err != nil {
		return nil, err
	}
	fmt.Fprintf(out, "Writing the rules of %d files (%d KB of source) with %s, effort %s.\n\n", len(jobs), bytes/1024, o.model, o.effort)

	var mu sync.Mutex
	var usage helpchat.Usage
	rules := 0
	var failed []string
	one := func(j job) {
		file, used, err := write(ctx, service, settings, system, appDir, j.source)
		mu.Lock()
		defer mu.Unlock()
		usage.Input += used.Input
		usage.Output += used.Output
		usage.CacheRead += used.CacheRead
		usage.CacheWrite += used.CacheWrite
		if err != nil {
			failed = append(failed, j.source.Path)
			fmt.Fprintf(out, "FAILED %s: %v\n", j.source.Path, err)
			return
		}
		rules += len(file.Rules)
		fmt.Fprintf(out, "%3d rules  %s\n", len(file.Rules), j.source.Path)
	}
	// The first file writes the model's cache of the instructions. The rest read it.
	one(jobs[0])
	each(jobs[1:], o.parallel, one)

	fmt.Fprintf(out, "\n%d rules in %d files. Tokens: %d input, %d from the cache, %d written to the cache, %d output.\n\n",
		rules, len(jobs)-len(failed), usage.Input, usage.CacheRead, usage.CacheWrite, usage.Output)
	return failed, nil
}

// needPhrases lists the sources whose rule file is for the source as it is and
// has a rule with no search phrases.
func needPhrases(appDir string, sources []helpcorpus.SourceFile, only string, out io.Writer) []job {
	var jobs []job
	for _, source := range sources {
		file, ok, err := helpcorpus.ReadRuleFile(appDir, source.Path)
		if err != nil || !ok || file.SHA256 != helpcorpus.SourceHash(source.Content) || !selected(source.Path, only) {
			continue
		}
		for _, rule := range file.Rules {
			if len(rule.Asks) == 0 {
				jobs = append(jobs, job{source: source, state: "asks"})
				break
			}
		}
	}
	return jobs
}

// phraseInstructions is the same for every file.
const phraseInstructions = `You get rules about Turbo, a company's internal app for time sheets, expenses, purchase orders (POs), jobs and clients. Each rule says what Turbo does in one case. Turbo's Help chat finds a rule by a keyword search, and the search fails when a person asks with other words than the rule has.

For each rule, write 3 to 5 short phrases that a member of staff could type when they ask about that case.
- Use the words of a person who is in the situation and does not know the rule: what they see, what they want to do, what went wrong. For a rule about an expense whose owner is also its approver, good phrases are "my own manager", "approve my own expense" and "self approval".
- Use words that the rule does not have: the everyday word, the other name for the thing, the symptom.
- A phrase is 2 to 6 words. It is not a sentence, and it does not repeat the title.
- Write only phrases that this rule answers.

Reply with one JSON object and no other text. Give one entry for each rule, with the number of the rule:
{"asks":[{"n":1,"phrases":["...","..."]}]}`

const maxPhrases = 6

// writeAllPhrases writes the search phrases of the rules that have none, and
// returns the paths that failed.
func writeAllPhrases(ctx context.Context, o options, service *helpchat.Service, appDir string, jobs []job, out io.Writer) ([]string, error) {
	settings := helpchat.DefaultSettings()
	settings.Model, settings.Effort = o.phraseModel, o.phraseEffort
	settings.MaxOutputTokens, settings.TurnTimeoutSeconds = 16000, 300
	fmt.Fprintf(out, "Writing the search phrases of the rules of %d files with %s, effort %s.\n\n", len(jobs), o.phraseModel, o.phraseEffort)

	var mu sync.Mutex
	var usage helpchat.Usage
	phrases := 0
	var failed []string
	one := func(j job) {
		n, used, err := writePhrases(ctx, service, settings, appDir, j.source)
		mu.Lock()
		defer mu.Unlock()
		usage.Input += used.Input
		usage.Output += used.Output
		usage.CacheRead += used.CacheRead
		usage.CacheWrite += used.CacheWrite
		if err != nil {
			failed = append(failed, j.source.Path)
			fmt.Fprintf(out, "FAILED %s: %v\n", j.source.Path, err)
			return
		}
		phrases += n
		fmt.Fprintf(out, "%3d phrases  %s\n", n, j.source.Path)
	}
	one(jobs[0])
	each(jobs[1:], o.parallel, one)
	fmt.Fprintf(out, "\n%d phrases in %d files. Tokens: %d input, %d from the cache, %d written to the cache, %d output.\n",
		phrases, len(jobs)-len(failed), usage.Input, usage.CacheRead, usage.CacheWrite, usage.Output)
	return failed, nil
}

// writePhrases asks for the phrases of every rule of one file and saves the
// rule file with them. The request holds the rules and not the code: a phrase
// is how a person asks, and the rule already says what the code does.
func writePhrases(ctx context.Context, service *helpchat.Service, settings helpchat.Settings, appDir string, source helpcorpus.SourceFile) (int, helpchat.Usage, error) {
	var total helpchat.Usage
	file, ok, err := helpcorpus.ReadRuleFile(appDir, source.Path)
	if err != nil || !ok {
		return 0, total, fmt.Errorf("the rule file is gone: %v", err)
	}
	var question strings.Builder
	fmt.Fprintf(&question, "The rules of %s:\n", source.Path)
	for i, rule := range file.Rules {
		fmt.Fprintf(&question, "\n%d. %s\n%s\n", i+1, rule.Title, rule.Text)
	}

	var lastErr error
	for try := 0; try < 2; try++ {
		reply, usage, err := service.Complete(ctx, settings, phraseInstructions, question.String())
		total.Input += usage.Input
		total.Output += usage.Output
		total.CacheRead += usage.CacheRead
		total.CacheWrite += usage.CacheWrite
		if err != nil {
			lastErr = err
			continue
		}
		asks, err := parsePhrases(reply, len(file.Rules))
		if err != nil {
			lastErr = err
			continue
		}
		count := 0
		for i := range file.Rules {
			file.Rules[i].Asks = asks[i]
			count += len(asks[i])
		}
		return count, total, os.WriteFile(helpcorpus.RuleFilePath(appDir, source.Path), []byte(file.Format()), 0o644)
	}
	return 0, total, lastErr
}

// parsePhrases reads the model's reply: the phrases of each rule, by the number
// of the rule. Every rule must have a phrase, or the file is asked about again
// on every run.
func parsePhrases(reply string, rules int) ([][]string, error) {
	from := strings.Index(reply, "{")
	if from < 0 {
		return nil, errors.New("the reply has no JSON object")
	}
	var parsed struct {
		Asks []struct {
			N       int      `json:"n"`
			Phrases []string `json:"phrases"`
		} `json:"asks"`
	}
	if err := json.NewDecoder(strings.NewReader(reply[from:])).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("the reply is not valid JSON: %w", err)
	}
	out := make([][]string, rules)
	for _, entry := range parsed.Asks {
		if entry.N < 1 || entry.N > rules || out[entry.N-1] != nil {
			continue
		}
		for _, phrase := range entry.Phrases {
			phrase = strings.Join(strings.Fields(phrase), " ")
			if phrase != "" && len(phrase) <= 80 && len(out[entry.N-1]) < maxPhrases {
				out[entry.N-1] = append(out[entry.N-1], phrase)
			}
		}
	}
	for i, phrases := range out {
		if len(phrases) == 0 {
			return nil, fmt.Errorf("rule %d has no phrases", i+1)
		}
	}
	return out, nil
}

func selected(sourcePath, only string) bool {
	if strings.TrimSpace(only) == "" {
		return true
	}
	for _, part := range strings.Split(only, ",") {
		if part = strings.TrimSpace(part); part != "" && strings.Contains(sourcePath, part) {
			return true
		}
	}
	return false
}

// orphans lists the rule files whose source file is gone or no longer gets rules.
func orphans(appDir string, sources []helpcorpus.SourceFile) ([]string, error) {
	want := map[string]bool{}
	for _, source := range sources {
		want[helpcorpus.RuleFilePath(appDir, source.Path)] = true
	}
	var extra []string
	err := filepath.WalkDir(helpcorpus.RulesDir(appDir), func(name string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() && strings.HasSuffix(name, ".md") && !want[name] {
			extra = append(extra, name)
		}
		return nil
	})
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	return extra, err
}

// each runs work for every job, a few at a time.
func each(jobs []job, parallel int, work func(job)) {
	if parallel < 1 {
		parallel = 1
	}
	queue := make(chan job)
	var wait sync.WaitGroup
	for i := 0; i < parallel; i++ {
		wait.Add(1)
		go func() {
			defer wait.Done()
			for j := range queue {
				work(j)
			}
		}()
	}
	for _, j := range jobs {
		queue <- j
	}
	close(queue)
	wait.Wait()
}

// instructions is the same for every file, so the model caches it with the
// glossary that follows it.
const instructions = `You write the rule summaries for Turbo, a company's internal app for time sheets, expenses, purchase orders (POs), jobs and clients. Turbo's Help chat searches these summaries to answer questions from staff. Staff act on the answers, so each rule must say exactly what the code does.

You get one source file with its line numbers. Write one rule for each thing in the file that a person who uses or administers Turbo can observe:
- who may do an action, and who may not;
- when Turbo refuses an action, with the message that it gives;
- a limit, a default or a threshold, with its number;
- what Turbo calculates or sets by itself, and from what;
- what a person can do with a record in each state, and what changes the state;
- an order that actions must follow;
- a case that gives a result a person would not expect: a fallback, a comparison that includes or excludes the equal value, a check that runs at one step and not at another, something that Turbo does not do.

How to write a rule:
- "title": at most 12 words that name the case, in plain words. A person searches with these words.
- "text": 1 to 4 short sentences in the active voice and the present tense, with one fact in a sentence. Say what Turbo does and when. Give the exact numbers, statuses and messages from the code.
- Use the words that staff use. The glossary below has them. Name the people by what they do, and add the name of the claim in parentheses, such as "a committer (the commit claim)". After the plain words, add the code's name for a field, a status or a setting in parentheses when a developer could search for it.
- A rule must be true and clear by itself, because a reader sees one rule and not the others. Say which record and which action it is about.
- "start" and "end": the first and the last line of the code that makes the rule true. Use the smallest range that holds it.
- Write only what the lines of this file show. Do not guess what another file does. When a rule depends on a function or a value from another file, name it and do not state its result.
- Leave out how the code is built: reading a request, wrapping an error, logging, types, and queries that only fetch a name.
- Write every rule that the file has. Go through the file from its first line to its last, one check and one branch at a time, and do not stop after the main rules. A long file can have 40 rules. A file with nothing that a person can observe gets an empty list.
- Give the message that Turbo shows. Do not give the HTTP status number.

Reply with one JSON object and no other text:
{"rules":[{"title":"...","text":"...","start":1,"end":9}]}`

func systemPrompt(repo string) (string, error) {
	glossary, err := os.ReadFile(filepath.Join(repo, "ui", "src", "docs", "glossary.md"))
	if err != nil {
		return "", fmt.Errorf("the rules use the words of the glossary: %w", err)
	}
	return instructions + "\n\n# The glossary of Turbo's user docs\n\n" + strings.TrimSpace(string(glossary)), nil
}

// numbered returns the source with a line number before each line.
func numbered(content string) (string, int) {
	lines := strings.Split(strings.TrimRight(content, "\n"), "\n")
	var out strings.Builder
	for i, line := range lines {
		fmt.Fprintf(&out, "%5d| %s\n", i+1, line)
	}
	return out.String(), len(lines)
}

// write asks the model for the rules of one source file and saves the rule
// file. It asks one more time when the reply is not what it must be.
//
// A model can stop early: in a test run it wrote 2 rules for a file of 400
// lines that has 30. So write looks for long parts of the file that no rule
// covers, and asks about those parts again.
func write(ctx context.Context, service *helpchat.Service, settings helpchat.Settings, system, appDir string, source helpcorpus.SourceFile) (helpcorpus.RuleFile, helpchat.Usage, error) {
	text, lines := numbered(source.Content)
	question := fmt.Sprintf("File: %s\n\n<source>\n%s</source>", source.Path, text)
	file := helpcorpus.RuleFile{Source: source.Path, SHA256: helpcorpus.SourceHash(source.Content), Model: settings.Model}

	var total helpchat.Usage
	ask := func(question string) ([]helpcorpus.Rule, error) {
		var lastErr error
		for try := 0; try < 2; try++ {
			reply, usage, err := service.Complete(ctx, settings, system, question)
			total.Input += usage.Input
			total.Output += usage.Output
			total.CacheRead += usage.CacheRead
			total.CacheWrite += usage.CacheWrite
			if err != nil {
				lastErr = err
				continue
			}
			rules, err := parseRules(reply, lines)
			if err != nil {
				lastErr = err
				continue
			}
			return rules, nil
		}
		return nil, lastErr
	}

	rules, err := ask(question)
	if err != nil {
		return file, total, err
	}
	file.Rules = rules
	for pass := 0; pass < maxFollowUps; pass++ {
		gaps, uncovered := gapsIn(source.Content, file.Rules)
		if uncovered < minUncoveredLines {
			break
		}
		// A failed follow-up keeps the rules that are there.
		more, err := ask(question + followUp(file.Rules, gaps))
		if err != nil || len(more) == 0 {
			break
		}
		file.Rules = merge(file.Rules, more)
	}
	sort.SliceStable(file.Rules, func(a, b int) bool { return file.Rules[a].Start < file.Rules[b].Start })

	name := helpcorpus.RuleFilePath(appDir, source.Path)
	if err := os.MkdirAll(filepath.Dir(name), 0o755); err != nil {
		return file, total, err
	}
	return file, total, os.WriteFile(name, []byte(file.Format()), 0o644)
}

const (
	// A part of a file with at least this many lines of code and no rule is a gap.
	minGapLines = 12
	// With fewer lines of code than this in all the gaps, there is no follow-up.
	minUncoveredLines = 30
	maxFollowUps      = 2
)

// gap is a run of lines that no rule covers.
type gap struct{ start, end int }

// gapsIn returns the long parts of a file that no rule covers, and how many
// lines of code they hold. A blank line or a comment is not a line of code.
func gapsIn(content string, rules []helpcorpus.Rule) (gaps []gap, uncovered int) {
	lines := strings.Split(strings.TrimRight(content, "\n"), "\n")
	covered := make([]bool, len(lines)+2)
	for _, rule := range rules {
		for i := rule.Start; i <= rule.End && i <= len(lines); i++ {
			covered[i] = true
		}
	}
	isCode := func(line string) bool {
		line = strings.TrimSpace(line)
		return line != "" && !strings.HasPrefix(line, "//") && !strings.HasPrefix(line, "--") &&
			!strings.HasPrefix(line, "/*") && !strings.HasPrefix(line, "*")
	}
	for i := 1; i <= len(lines); i++ {
		if covered[i] {
			continue
		}
		start, code := i, 0
		for ; i <= len(lines) && !covered[i]; i++ {
			if isCode(lines[i-1]) {
				code++
			}
		}
		if code >= minGapLines {
			gaps = append(gaps, gap{start, i - 1})
			uncovered += code
		}
	}
	return gaps, uncovered
}

// followUp is what the second request adds to the first.
func followUp(rules []helpcorpus.Rule, gaps []gap) string {
	var out strings.Builder
	out.WriteString("\n\nThese rules are written for this file:\n")
	for _, rule := range rules {
		fmt.Fprintf(&out, "- lines %d-%d: %s\n", rule.Start, rule.End, rule.Title)
	}
	if len(rules) == 0 {
		out.WriteString("(none)\n")
	}
	out.WriteString("\nThese lines have no rule yet:")
	for _, g := range gaps {
		fmt.Fprintf(&out, " %d-%d", g.start, g.end)
	}
	out.WriteString(".\nRead those lines again, one check and one branch at a time. Write a rule for each thing in them that a person can observe: " +
		"each condition that refuses an action, each value that must match, each limit, and each thing that Turbo sets. " +
		"Do not write again a rule from the list. Lines that hold only imports, types, or code that fetches data have no rules: " +
		"if no line has a rule, reply {\"rules\":[]}.")
	return out.String()
}

// merge adds the new rules, except one that the file has already.
func merge(rules, more []helpcorpus.Rule) []helpcorpus.Rule {
	for _, rule := range more {
		known := false
		for _, old := range rules {
			if strings.EqualFold(old.Title, rule.Title) || (old.Start == rule.Start && old.End == rule.End) {
				known = true
			}
		}
		if !known {
			rules = append(rules, rule)
		}
	}
	return rules
}

// parseRules reads the model's reply. A rule with no text, or with lines that
// the file does not have, makes the whole reply an error: a rule that points at
// the wrong code is worse than no rule.
func parseRules(reply string, lines int) ([]helpcorpus.Rule, error) {
	from := strings.Index(reply, "{")
	if from < 0 {
		return nil, errors.New("the reply has no JSON object")
	}
	var parsed struct {
		Rules []helpcorpus.Rule `json:"rules"`
	}
	if err := json.NewDecoder(strings.NewReader(reply[from:])).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("the reply is not valid JSON: %w", err)
	}
	for i, rule := range parsed.Rules {
		switch {
		case strings.TrimSpace(rule.Title) == "" || strings.TrimSpace(rule.Text) == "":
			return nil, fmt.Errorf("rule %d has no title or no text", i+1)
		case rule.Start < 1 || rule.End < rule.Start || rule.End > lines:
			return nil, fmt.Errorf("rule %d names the lines %d-%d of a file with %d lines", i+1, rule.Start, rule.End, lines)
		}
	}
	return parsed.Rules, nil
}
