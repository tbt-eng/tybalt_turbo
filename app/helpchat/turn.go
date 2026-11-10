package helpchat

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"tybalt/helpcorpus"
)

// The results of a turn. They are the values of chat_turns.status, with
// "running" for a turn that has not ended.
const (
	StatusRunning   = "running"
	StatusAnswered  = "answered"
	StatusGap       = "gap"
	StatusRefused   = "refused"
	StatusCancelled = "cancelled"
	StatusError     = "error"
)

const (
	// One search returns rules, code and design notes as three groups. In one
	// list, the notes push the code out, because their headings are in plain
	// words. The code is what says how Turbo behaves today, and a rule says in
	// plain words what some lines of the code do (see helpcorpus/rules.go).
	searchRuleHits = 5
	searchCodeHits = 6
	searchNoteHits = 2
	// A search hit shows this much of a chunk. read_source shows all of it.
	searchHitChars = 1600
	maxQueryChars  = 300
	maxErrorChars  = 500
)

// Events receives the progress of a turn, to stream to the browser.
type Events interface {
	// Status says what the server is doing, such as "searching".
	Status(state, detail string)
	// Text is the next piece of the answer.
	Text(delta string)
}

type turnInput struct {
	settings Settings
	claims   helpcorpus.Claims
	request  Request
}

// Source is something that the answer cites: a page of the user docs, or a
// chunk that a tool returned in this turn.
type Source struct {
	Kind string `json:"kind"` // "doc" or "chunk"
	// URL is the page for a doc, such as /docs/time#submit-a-week.
	URL string `json:"url,omitempty"`
	// Ref is the mark in the answer for a chunk, such as S3.
	Ref     string `json:"ref,omitempty"`
	ID      string `json:"id,omitempty"`
	Path    string `json:"path,omitempty"`
	Start   int    `json:"start,omitempty"`
	End     int    `json:"end,omitempty"`
	Heading string `json:"heading,omitempty"`
	Type    string `json:"type,omitempty"` // the chunk's kind: notes or code
}

// Gap is the model's report that Help could not answer.
type Gap struct {
	Kind  string `json:"kind"`
	Topic string `json:"topic"`
}

// toolLog records one tool call. It holds chunk ids, not chunk text: the corpus
// hash on the row says which text the ids mean.
type toolLog struct {
	Tool  string          `json:"tool"`
	Input json.RawMessage `json:"input"`
	Hits  []toolHit       `json:"hits,omitempty"`
	Error string          `json:"error,omitempty"`
	MS    int64           `json:"ms"`
}

type toolHit struct {
	Ref   string  `json:"ref"`
	ID    string  `json:"id"`
	Score float64 `json:"score,omitempty"`
}

type turnResult struct {
	Status       string
	Answer       string
	Sources      []Source
	Gap          *Gap
	ToolCalls    []toolLog
	StopReason   string
	Error        string
	RequestID    string
	Usage        Usage
	FirstTokenMS int64
	TotalMS      int64
}

// turn is the state of one turn while it runs.
type turn struct {
	index  *helpcorpus.Index
	input  turnInput
	events Events
	result turnResult
	// refs gives each chunk that a tool returned a short name for the answer to
	// cite. A chunk keeps its name if a later tool call returns it again.
	refs    map[string]string
	chunks  map[string]helpcorpus.Chunk
	started time.Time
}

// runTurn asks the model, runs the tools that it calls, and returns the answer
// with everything that the log needs. It writes nothing to the database.
func runTurn(ctx context.Context, model Model, index *helpcorpus.Index, input turnInput, events Events) turnResult {
	t := &turn{
		index:   index,
		input:   input,
		events:  events,
		refs:    map[string]string{},
		chunks:  map[string]helpcorpus.Chunk{},
		started: time.Now(),
	}
	t.result.ToolCalls = []toolLog{}
	t.result.Sources = []Source{}
	t.run(ctx, model.Start(input.request))
	t.result.TotalMS = max(time.Since(t.started).Milliseconds(), t.result.FirstTokenMS)
	return t.result
}

func (t *turn) run(ctx context.Context, session Session) {
	var answer []string
	defer func() {
		t.result.Answer = strings.TrimSpace(strings.Join(answer, "\n\n"))
		t.result.Sources = t.cited(t.result.Answer)
	}()

	onText := func(text string) {
		if t.result.FirstTokenMS == 0 {
			t.result.FirstTokenMS = max(time.Since(t.started).Milliseconds(), 1)
		}
		t.events.Text(text)
	}

	// The model gets one more request than it has tool calls, to write the answer.
	for calls := 0; ; {
		reply, err := session.Next(ctx, onText)
		if err != nil {
			t.fail(ctx, err)
			return
		}
		t.result.Usage.add(reply.Usage)
		t.result.RequestID = reply.RequestID
		t.result.StopReason = reply.StopReason
		if text := strings.TrimSpace(reply.Text); text != "" {
			answer = append(answer, text)
		}

		if reply.StopReason == stopRefusal {
			t.result.Status = StatusRefused
			return
		}
		if reply.StopReason != stopToolUse || len(reply.ToolCalls) == 0 {
			break
		}
		if calls > t.input.settings.MaxToolCalls {
			t.result.Status = StatusError
			t.result.Error = "the model kept calling tools after the limit"
			return
		}

		results := make([]ToolResult, len(reply.ToolCalls))
		for i, call := range reply.ToolCalls {
			// report_gap is free. The limit is on the calls that read the corpus.
			if call.Name != toolGap {
				calls++
			}
			if calls > t.input.settings.MaxToolCalls {
				results[i] = ToolResult{ID: call.ID, Content: "You have used all your tool calls for this question. Answer now with what you have."}
				continue
			}
			results[i] = t.callTool(call)
		}
		session.AddToolResults(results)
	}

	switch {
	case t.result.Gap != nil:
		t.result.Status = StatusGap
	case len(answer) == 0:
		t.result.Status = StatusError
		t.result.Error = "the model gave no answer"
	default:
		t.result.Status = StatusAnswered
	}
}

func (t *turn) fail(ctx context.Context, err error) {
	switch {
	case errors.Is(ctx.Err(), context.Canceled):
		// The user closed the page or stopped the answer.
		t.result.Status = StatusCancelled
	case errors.Is(ctx.Err(), context.DeadlineExceeded):
		t.result.Status = StatusError
		t.result.Error = "the answer took too long"
	default:
		t.result.Status = StatusError
		t.result.Error = truncate(err.Error(), maxErrorChars)
	}
}

// callTool runs one tool for the model. Every read of the corpus passes the
// user's claims, so the model gets only what the user may read.
func (t *turn) callTool(call ToolCall) ToolResult {
	started := time.Now()
	entry := toolLog{Tool: call.Name, Input: call.Input}
	if !json.Valid(entry.Input) {
		entry.Input = json.RawMessage(`null`)
	}
	fail := func(message string) ToolResult {
		entry.Error = message
		entry.MS = time.Since(started).Milliseconds()
		t.result.ToolCalls = append(t.result.ToolCalls, entry)
		return ToolResult{ID: call.ID, Content: message, IsError: true}
	}
	done := func(content string) ToolResult {
		entry.MS = time.Since(started).Milliseconds()
		t.result.ToolCalls = append(t.result.ToolCalls, entry)
		return ToolResult{ID: call.ID, Content: content}
	}

	switch call.Name {
	case toolSearch:
		var input struct {
			Query string `json:"query"`
		}
		if err := json.Unmarshal(call.Input, &input); err != nil || strings.TrimSpace(input.Query) == "" {
			return fail("The input needs a query.")
		}
		query := truncate(input.Query, maxQueryChars)
		t.events.Status("searching", query)
		hits, err := searchCorpus(t.index, t.input.claims, query)
		if err != nil {
			return fail("The search failed.")
		}
		if len(hits) == 0 {
			return done("No matches. Try other words.")
		}
		var out strings.Builder
		for _, hit := range hits {
			ref := t.ref(hit.Chunk)
			entry.Hits = append(entry.Hits, toolHit{Ref: ref, ID: hit.ID, Score: hit.Score})
			writeChunk(&out, ref, hit.Chunk, searchHitChars)
		}
		return done(out.String())

	case toolRead:
		var input struct {
			Path string `json:"path"`
			From int    `json:"from"`
			To   int    `json:"to"`
		}
		if err := json.Unmarshal(call.Input, &input); err != nil || input.Path == "" {
			return fail("The input needs a path.")
		}
		t.events.Status("reading", "")
		chunks := t.index.Read(t.input.claims, input.Path, input.From, input.To)
		if len(chunks) == 0 {
			return done("Nothing at that path and those lines. Use a path from a search hit.")
		}
		var out strings.Builder
		for _, chunk := range chunks {
			ref := t.ref(chunk)
			entry.Hits = append(entry.Hits, toolHit{Ref: ref, ID: chunk.ID})
			writeChunk(&out, ref, chunk, 0)
		}
		return done(out.String())

	case toolGap:
		var input Gap
		if err := json.Unmarshal(call.Input, &input); err != nil || (input.Kind != gapNotCovered && input.Kind != gapOffTopic) {
			return fail(`The input needs a kind: "not_covered" or "off_topic".`)
		}
		input.Topic = truncate(strings.TrimSpace(input.Topic), 200)
		t.result.Gap = &input
		return done("Recorded. Now tell the user in one or two sentences.")
	}
	return fail("There is no tool with that name.")
}

// searchCorpus is one search for the model: the best rules, then the best code
// chunks, then the best chunks of the design notes. The handbook already holds
// the user docs, so a search leaves them out.
func searchCorpus(index *helpcorpus.Index, claims helpcorpus.Claims, query string) ([]helpcorpus.Hit, error) {
	var all []helpcorpus.Hit
	for _, group := range []struct {
		kind  string
		limit int
	}{
		{helpcorpus.KindRules, searchRuleHits},
		{helpcorpus.KindCode, searchCodeHits},
		{helpcorpus.KindNotes, searchNoteHits},
	} {
		hits, err := index.Search(claims, query, helpcorpus.SearchOptions{Limit: group.limit, Kinds: []string{group.kind}})
		if err != nil {
			return nil, err
		}
		all = append(all, hits...)
	}
	return all, nil
}

// ref returns the short name of a chunk in this turn, such as S3.
func (t *turn) ref(chunk helpcorpus.Chunk) string {
	if ref, ok := t.refs[chunk.ID]; ok {
		return ref
	}
	ref := fmt.Sprintf("S%d", len(t.refs)+1)
	t.refs[chunk.ID] = ref
	t.chunks[ref] = chunk
	return ref
}

func writeChunk(out *strings.Builder, ref string, chunk helpcorpus.Chunk, maxChars int) {
	fmt.Fprintf(out, "<chunk ref=%q type=%q path=%q lines=\"%d-%d\"", ref, chunkType(chunk), chunk.Path, chunk.Start, chunk.End)
	if chunk.Heading != "" {
		fmt.Fprintf(out, " in=%q", chunk.Heading)
	}
	out.WriteString(">\n")
	text := chunk.Text
	if maxChars > 0 && len(text) > maxChars {
		text = truncate(text, maxChars) + "\n[cut: use read_source for the rest]"
	}
	out.WriteString(text)
	out.WriteString("\n</chunk>\n")
}

// chunkType tells the model what a chunk is, because it trusts the two kinds
// differently.
func chunkType(chunk helpcorpus.Chunk) string {
	switch chunk.Kind {
	case helpcorpus.KindNotes:
		return "design note"
	case helpcorpus.KindRules:
		return "rule"
	}
	return chunk.Kind
}

var (
	// [S3] or [S3, S4]
	chunkCitation = regexp.MustCompile(`\[(S\d+(?:\s*,\s*S\d+)*)\]`)
	chunkRef      = regexp.MustCompile(`S\d+`)
	// [Time](/docs/time) or [Submit a week](/docs/time#submit-a-week)
	docCitation = regexp.MustCompile(`\]\((/docs/([a-z0-9-]+)(?:#[a-z0-9-]+)?)\)`)
)

// cited returns the sources that the answer names, in the order that it names
// them. It keeps a source only if the server knows it: a page of the user docs,
// or a chunk that a tool returned in this turn. The model cannot invent one.
func (t *turn) cited(answer string) []Source {
	type found struct {
		at     int
		source Source
	}
	var all []found
	seen := map[string]bool{}
	for _, match := range docCitation.FindAllStringSubmatchIndex(answer, -1) {
		url, slug := answer[match[2]:match[3]], answer[match[4]:match[5]]
		if seen[url] || len(t.index.Read(t.input.claims, "ui/src/docs/"+slug+".md", 0, 0)) == 0 {
			continue
		}
		seen[url] = true
		all = append(all, found{match[0], Source{Kind: "doc", URL: url}})
	}
	for _, match := range chunkCitation.FindAllStringSubmatchIndex(answer, -1) {
		for _, ref := range chunkRef.FindAllString(answer[match[2]:match[3]], -1) {
			chunk, ok := t.chunks[ref]
			if seen[ref] || !ok {
				continue
			}
			seen[ref] = true
			all = append(all, found{match[0], Source{
				Kind: "chunk", Ref: ref, ID: chunk.ID, Path: chunk.Path,
				Start: chunk.Start, End: chunk.End, Heading: chunk.Heading, Type: chunk.Kind,
			}})
		}
	}
	// A stable sort by position keeps the two kinds in reading order.
	for i := 1; i < len(all); i++ {
		for j := i; j > 0 && all[j].at < all[j-1].at; j-- {
			all[j], all[j-1] = all[j-1], all[j]
		}
	}
	sources := make([]Source, len(all))
	for i, item := range all {
		sources[i] = item.source
	}
	return sources
}

// truncate cuts text to at most max bytes, at the end of a character.
func truncate(text string, max int) string {
	if len(text) <= max {
		return text
	}
	cut := max
	for cut > 0 && cut < len(text) && text[cut]&0xC0 == 0x80 {
		cut--
	}
	return text[:cut]
}
