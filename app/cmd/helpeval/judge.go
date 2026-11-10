package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"sync"

	"tybalt/helpchat"
)

// The judge is a second model. It does not decide the grade. It says, for each
// key point, if the answer states it, and the code below makes the grade from
// that. So a grade always traces back to a sentence of the reference.
const judgePrompt = `You grade one answer from the Help chat of Turbo, a company's internal business app, against a reference.

The reference is correct. An engineer wrote it from the source code. You also get the lines of the source that the reference comes from. Grade only against the reference and the source. Do not use your own idea of how such an app works.

You get the question, the reference answer, numbered key points, numbered wrong claims, the source, and the candidate answer.

For each key point, give one verdict:
- "stated": the candidate says this fact, in any words. A statement that means the same, or that is more exact, counts.
- "missing": the candidate does not say it, or says only a part of it. Text so vague that it could mean the key point or not is "missing".
- "contradicted": the candidate says something that cannot be true if the key point is true. A candidate that says less than the key point is "missing", not "contradicted".

For each wrong claim, say if the candidate makes that claim: true or false. A candidate makes a claim only if it states it as a fact. If it says that it could not confirm the claim, it does not make it.

"other_error": use this only when the candidate states, as a fact, something that a sentence of the reference or a line of the source shows to be false. Quote the candidate's words, then the words of the reference or the source that show it. Use an empty string in every other case. These are not errors:
- a statement that the source supports, even when the reference does not mention it;
- a statement about something that the reference and the source do not cover, even if you doubt it;
- a statement that leaves out a case or a detail;
- more detail than the reference has;
- a statement that the candidate marks as not confirmed, such as "I could not confirm this" or "the notes describe a plan";
- the length, the style, the marks such as [S3], the links, or advice to ask a person.

"abstained": true if the candidate gives no answer of substance: it says that it does not know or cannot find the answer, or it only sends the user to a person or a page.

Reply with one JSON object and no other text:
{"points":[{"n":1,"verdict":"stated","why":"a few words"}],"wrong_claims":[{"n":1,"made":false}],"other_error":"","abstained":false}

Give one entry in "points" for each key point, in order. Give one entry in "wrong_claims" for each wrong claim, in order, or an empty list if there are none.`

// Verdict is the judge's reply.
type Verdict struct {
	Points []struct {
		N       int    `json:"n"`
		Verdict string `json:"verdict"`
		Why     string `json:"why"`
	} `json:"points"`
	WrongClaims []struct {
		N    int  `json:"n"`
		Made bool `json:"made"`
	} `json:"wrong_claims"`
	OtherError string `json:"other_error"`
	Abstained  bool   `json:"abstained"`
}

// grade turns the verdict into a grade, the number of key points stated, and
// the reason for a "wrong". An answer is wrong when it contradicts a key point
// or makes one of the listed wrong claims.
func (v *Verdict) grade(c Case) (grade string, stated int, wrongReason string) {
	var reasons []string
	for i, point := range v.Points {
		switch point.Verdict {
		case "stated":
			stated++
		case "contradicted":
			if i < len(c.KeyPoints) {
				reasons = append(reasons, "contradicts: "+c.KeyPoints[i])
			}
		}
	}
	for i, claim := range v.WrongClaims {
		if claim.Made && i < len(c.WrongIf) {
			reasons = append(reasons, "claims: "+c.WrongIf[i])
		}
	}
	switch {
	case len(reasons) > 0:
		return gradeWrong, stated, strings.Join(reasons, "; ")
	case stated == len(c.KeyPoints):
		return gradeCorrect, stated, ""
	case stated == 0 && v.Abstained:
		return gradeAbstained, 0, ""
	default:
		return gradePartial, stated, ""
	}
}

// flag is a statement of the answer that the judge takes for an error, and that
// is not a key point or a listed wrong claim. It does not change the grade. In
// the first run of 100 questions the judge gave three, and each statement was
// true: the judge reads a few lines of the source, and the statement was about
// lines that it did not have. So a person reads the flagged answers.
func (v *Verdict) flag() string {
	return strings.TrimSpace(v.OtherError)
}

// judgeInput is what the judge reads about one answer. source is the lines of
// the code that settle the question. With them the judge can tell a true
// statement that the reference does not have from a false one.
func judgeInput(c Case, answer, source string) string {
	var out strings.Builder
	fmt.Fprintf(&out, "<question>\n%s\n</question>\n\n<reference_answer>\n%s\n</reference_answer>\n\n<key_points>\n", c.Question, c.IdealAnswer)
	for i, point := range c.KeyPoints {
		fmt.Fprintf(&out, "%d. %s\n", i+1, point)
	}
	out.WriteString("</key_points>\n\n<wrong_claims>\n")
	for i, claim := range c.WrongIf {
		fmt.Fprintf(&out, "%d. %s\n", i+1, claim)
	}
	if len(c.WrongIf) == 0 {
		out.WriteString("(none)\n")
	}
	out.WriteString("</wrong_claims>\n\n<source>\n")
	if strings.TrimSpace(source) == "" {
		source = "(not given)\n"
	}
	out.WriteString(source)
	out.WriteString("</source>\n\n")
	if strings.TrimSpace(answer) == "" {
		answer = "(the candidate gave no text)"
	}
	fmt.Fprintf(&out, "<candidate_answer>\n%s\n</candidate_answer>\n\n", answer)
	// The judge has given one verdict too many, for a wrong claim, without this.
	fmt.Fprintf(&out, "Give exactly %d entries in \"points\" and exactly %d in \"wrong_claims\".", len(c.KeyPoints), len(c.WrongIf))
	return out.String()
}

// grade asks the judge about one answer. It asks one more time if the reply is
// not the JSON that it must be.
func grade(ctx context.Context, service *helpchat.Service, judge helpchat.Settings, c Case, answer string) (*Verdict, helpchat.Usage, error) {
	var total helpchat.Usage
	var lastErr error
	input := judgeInput(c, answer, sourceOf(service, judge, c))
	for try := 0; try < 2; try++ {
		text, usage, err := service.Complete(ctx, judge, judgePrompt, input)
		total.Input += usage.Input
		total.Output += usage.Output
		total.CacheRead += usage.CacheRead
		total.CacheWrite += usage.CacheWrite
		if err != nil {
			lastErr = err
			continue
		}
		verdict, err := parseVerdict(text, c)
		if err != nil {
			lastErr = err
			continue
		}
		return verdict, total, nil
	}
	return nil, total, lastErr
}

// sourceOf returns the chunks of the corpus that hold the lines that settle a
// question. A source that the corpus does not hold is left out.
func sourceOf(service *helpchat.Service, settings helpchat.Settings, c Case) string {
	var out strings.Builder
	for _, source := range c.Sources {
		if text, err := service.SourceText(settings, source.Path, source.Start, source.End); err == nil {
			out.WriteString(text)
		}
	}
	return out.String()
}

func parseVerdict(text string, c Case) (*Verdict, error) {
	from := strings.Index(text, "{")
	if from < 0 {
		return nil, fmt.Errorf("the judge's reply has no JSON object")
	}
	// Decode reads one JSON value and stops, so words after it do no harm.
	var verdict Verdict
	if err := json.NewDecoder(strings.NewReader(text[from:])).Decode(&verdict); err != nil {
		return nil, fmt.Errorf("the judge's reply is not valid JSON: %w", err)
	}
	// The judge sometimes adds a verdict, for a wrong claim, to the key points.
	// Each verdict has the number of its key point, so take one for each number.
	if len(verdict.Points) > len(c.KeyPoints) {
		kept := verdict.Points[:0:0]
		for n := 1; n <= len(c.KeyPoints); n++ {
			for _, point := range verdict.Points {
				if point.N == n {
					kept = append(kept, point)
					break
				}
			}
		}
		verdict.Points = kept
	}
	if len(verdict.Points) != len(c.KeyPoints) {
		return nil, fmt.Errorf("the judge gave %d verdicts for %d key points", len(verdict.Points), len(c.KeyPoints))
	}
	for _, point := range verdict.Points {
		if point.Verdict != "stated" && point.Verdict != "missing" && point.Verdict != "contradicted" {
			return nil, fmt.Errorf("the judge gave the verdict %q", point.Verdict)
		}
	}
	return &verdict, nil
}

// judgeCheck tests the judge on answers whose grade is known. A judge that
// passes a bad answer, or fails the reference itself, can't be trusted with
// the real ones.
// judgeCase is one answer whose right grade a person has worked out. The set is
// in judge_cases.json. It holds the answers that the judge got wrong in earlier
// runs, with answers that are wrong and that it must still catch.
type judgeCase struct {
	Name string `json:"name"`
	// Case is the id of the question.
	Case string `json:"case"`
	// Expect lists the grades that are right for the answer.
	Expect []string `json:"expect"`
	Why    string   `json:"why"`
	Answer string   `json:"answer"`
}

// checkJudgeCases grades each answer of the set and compares with the grade
// that a person gave. all is every question, by id.
func checkJudgeCases(ctx context.Context, service *helpchat.Service, judge helpchat.Settings, path string, all []Case, parallel int, out io.Writer) error {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	var set []judgeCase
	if err := json.Unmarshal(data, &set); err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	byID := map[string]Case{}
	for _, c := range all {
		byID[c.ID] = c
	}
	type result struct {
		name, grade, reason string
		ok                  bool
		err                 error
	}
	results := make([]result, len(set))
	var jobs []Case
	index := map[string]int{}
	for i, jc := range set {
		if _, ok := byID[jc.Case]; !ok {
			return fmt.Errorf("%s: %s names the question %q, which is not in the questions file", path, jc.Name, jc.Case)
		}
		id := fmt.Sprintf("%d", i)
		index[id] = i
		jobs = append(jobs, Case{ID: id})
	}
	each(jobs, parallel, func(j Case) {
		i := index[j.ID]
		c := byID[set[i].Case]
		verdict, _, err := grade(ctx, service, judge, c, set[i].Answer)
		r := result{name: set[i].Name, err: err}
		if err == nil {
			r.grade, _, r.reason = verdict.grade(c)
			if flag := verdict.flag(); flag != "" {
				r.reason = strings.TrimSpace(r.reason + " flagged: " + flag)
			}
			r.ok = slices.Contains(set[i].Expect, r.grade)
		}
		results[i] = r
	})
	agree := 0
	for _, r := range results {
		if r.ok {
			agree++
		}
	}
	fmt.Fprintf(out, "\nAnswers with a known grade (%s): the judge agrees on %d of %d.\n", path, agree, len(set))
	for i, r := range results {
		switch {
		case r.err != nil:
			fmt.Fprintf(out, "  ERROR  %s: %v\n", r.name, r.err)
		case !r.ok:
			fmt.Fprintf(out, "  WRONG  %s: the judge says %s, and it is %s. %s\n         the judge's reason: %s\n",
				r.name, r.grade, strings.Join(set[i].Expect, " or "), set[i].Why, r.reason)
		}
	}
	return nil
}

func judgeCheck(ctx context.Context, service *helpchat.Service, judge helpchat.Settings, cases []Case, parallel int, out io.Writer) error {
	type probe struct {
		name   string
		answer func(i int) string
		pass   func(grade string, stated int) bool
	}
	probes := []probe{
		{"the reference answer itself must be correct",
			func(i int) string { return cases[i].IdealAnswer },
			func(grade string, _ int) bool { return grade == gradeCorrect }},
		{"no text must not be correct or partial",
			func(int) string { return "" },
			func(grade string, stated int) bool { return stated == 0 && grade != gradeCorrect }},
		{"\"I don't know\" must be abstained",
			func(int) string {
				return "The Turbo help does not cover this. Ask your manager or the Turbo team."
			},
			func(grade string, _ int) bool { return grade == gradeAbstained }},
		{"the reference answer of another question must not be correct",
			func(i int) string { return cases[(i+len(cases)/2)%len(cases)].IdealAnswer },
			func(grade string, _ int) bool { return grade != gradeCorrect }},
	}

	type result struct {
		probe, id, grade string
		ok               bool
		err              error
	}
	var mu sync.Mutex
	var results []result
	type job struct{ p, i int }
	var jobs []Case // one Case for each job, with the job in its ID, to reuse each()
	index := map[string]job{}
	for p := range probes {
		for i := range cases {
			id := fmt.Sprintf("%d/%s", p, cases[i].ID)
			index[id] = job{p, i}
			jobs = append(jobs, Case{ID: id})
		}
	}
	each(jobs, parallel, func(j Case) {
		p, i := index[j.ID].p, index[j.ID].i
		verdict, _, err := grade(ctx, service, judge, cases[i], probes[p].answer(i))
		r := result{probe: probes[p].name, id: cases[i].ID, err: err}
		if err == nil {
			var stated int
			r.grade, stated, _ = verdict.grade(cases[i])
			r.ok = probes[p].pass(r.grade, stated)
		}
		mu.Lock()
		results = append(results, r)
		mu.Unlock()
	})

	fmt.Fprintf(out, "Judge check: %s, effort %s, %d questions.\n\n", judge.Model, judge.Effort, len(cases))
	for _, p := range probes {
		passed, n := 0, 0
		var failed []string
		for _, r := range results {
			if r.probe != p.name {
				continue
			}
			n++
			switch {
			case r.err != nil:
				failed = append(failed, r.id+" (error: "+r.err.Error()+")")
			case r.ok:
				passed++
			default:
				failed = append(failed, r.id+" ("+r.grade+")")
			}
		}
		fmt.Fprintf(out, "%-62s %d/%d", p.name, passed, n)
		if len(failed) > 0 {
			fmt.Fprintf(out, "  failed: %s", strings.Join(failed, ", "))
		}
		fmt.Fprintln(out)
	}
	return nil
}
