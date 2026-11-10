package helpchat

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"tybalt/helpcorpus"
)

// This file lets the benchmark (cmd/helpeval) ask a question the way the route
// does: the same prompt, tools and tool loop. It uses no database, so it has
// no rate limit and writes no chat_turns row.

// EvalRequest is one benchmark question.
type EvalRequest struct {
	Settings Settings
	// AllClaims is the list of claims that the prompt shows the model.
	AllClaims []Claim
	// Held names the claims of the person who asks.
	Held     []string
	Question string
	// Given is source text that the model gets with the question. It measures
	// the answer when the search has nothing left to find.
	Given string
	// NoTools takes the tools away, so the model has the handbook and no more.
	NoTools bool
}

// EvalChunk is a chunk that a tool call returned.
type EvalChunk struct {
	Ref   string `json:"ref,omitempty"`
	Path  string `json:"path"`
	Start int    `json:"start"`
	End   int    `json:"end"`
}

// EvalToolCall is one tool call of a turn.
type EvalToolCall struct {
	Tool   string          `json:"tool"`
	Input  json.RawMessage `json:"input"`
	Chunks []EvalChunk     `json:"chunks,omitempty"`
	Error  string          `json:"error,omitempty"`
	MS     int64           `json:"ms"`
}

// EvalResult is what a turn gave.
type EvalResult struct {
	Status       string         `json:"status"`
	Answer       string         `json:"answer"`
	Sources      []Source       `json:"sources"`
	Gap          *Gap           `json:"gap,omitempty"`
	ToolCalls    []EvalToolCall `json:"tool_calls"`
	StopReason   string         `json:"stop_reason"`
	Error        string         `json:"error,omitempty"`
	RequestID    string         `json:"request_id,omitempty"`
	Usage        Usage          `json:"usage"`
	FirstTokenMS int64          `json:"ms_first_token"`
	TotalMS      int64          `json:"ms_total"`
}

// noEvents drops the progress of a turn.
type noEvents struct{}

func (noEvents) Status(string, string) {}
func (noEvents) Text(string)           {}

// CorpusHash names the corpus that the answers come from.
func (s *Service) CorpusHash(settings Settings) (string, error) {
	if _, err := s.ready(settings); err != nil {
		return "", err
	}
	return s.Index.Hash(), nil
}

// Evaluate answers one question as handleAsk does.
func (s *Service) Evaluate(ctx context.Context, eval EvalRequest) (EvalResult, error) {
	if err := eval.Settings.validate(); err != nil {
		return EvalResult{}, err
	}
	model, err := s.ready(eval.Settings)
	if err != nil {
		return EvalResult{}, err
	}
	held := helpcorpus.Claims{}
	for _, name := range eval.Held {
		held[name] = true
	}

	request := s.request(eval.Settings)
	request.Context = contextPrompt(eval.AllClaims, held, "")
	if eval.Given != "" {
		request.Context += "\n# Source for this question\n\n" +
			"The Turbo team gives you these parts of Turbo's source. They hold the answer.\n\n" + eval.Given
	}
	if eval.NoTools {
		request.Tools = nil
	}
	request.Question = eval.Question

	ctx, cancel := context.WithTimeout(ctx, time.Duration(eval.Settings.TurnTimeoutSeconds)*time.Second)
	defer cancel()
	result := runTurn(ctx, model, s.Index, turnInput{settings: eval.Settings, claims: held, request: request}, noEvents{})

	out := EvalResult{
		Status:       result.Status,
		Answer:       result.Answer,
		Sources:      result.Sources,
		Gap:          result.Gap,
		ToolCalls:    make([]EvalToolCall, len(result.ToolCalls)),
		StopReason:   result.StopReason,
		Error:        result.Error,
		RequestID:    result.RequestID,
		Usage:        result.Usage,
		FirstTokenMS: result.FirstTokenMS,
		TotalMS:      result.TotalMS,
	}
	for i, call := range result.ToolCalls {
		entry := EvalToolCall{Tool: call.Tool, Input: call.Input, Error: call.Error, MS: call.MS}
		for _, hit := range call.Hits {
			entry.Chunks = append(entry.Chunks, s.chunkByID(hit.Ref, hit.ID))
		}
		out.ToolCalls[i] = entry
	}
	return out, nil
}

// chunkByID finds the lines of a chunk from its id, which is path:start.
func (s *Service) chunkByID(ref, id string) EvalChunk {
	chunk := EvalChunk{Ref: ref, Path: id}
	cut := strings.LastIndex(id, ":")
	if cut < 0 {
		return chunk
	}
	start, err := strconv.Atoi(id[cut+1:])
	if err != nil {
		return chunk
	}
	chunk.Path, chunk.Start, chunk.End = id[:cut], start, start
	for _, found := range s.Index.Read(helpcorpus.Claims{}, chunk.Path, start, start) {
		if found.ID == id {
			chunk.End = found.End
		}
	}
	return chunk
}

// SearchOnly runs one search as the search_turbo tool does, with no model.
func (s *Service) SearchOnly(query string) ([]EvalChunk, error) {
	if err := s.loadIndex(); err != nil {
		return nil, err
	}
	hits, err := searchCorpus(s.Index, helpcorpus.Claims{}, truncate(query, maxQueryChars))
	if err != nil {
		return nil, err
	}
	chunks := make([]EvalChunk, len(hits))
	for i, hit := range hits {
		chunks[i] = EvalChunk{Path: hit.Path, Start: hit.Start, End: hit.End}
	}
	return chunks, nil
}

// SourceText returns the chunks that hold the lines of a file, in the form that
// a tool result has.
func (s *Service) SourceText(settings Settings, path string, from, to int) (string, error) {
	if _, err := s.ready(settings); err != nil {
		return "", err
	}
	chunks := s.Index.Read(helpcorpus.Claims{}, path, from, to)
	if len(chunks) == 0 {
		return "", fmt.Errorf("the corpus has nothing at %s lines %d-%d", path, from, to)
	}
	var out strings.Builder
	for _, chunk := range chunks {
		writeChunk(&out, "", chunk, 0)
	}
	return out.String(), nil
}

// Complete sends one plain request with no tools and returns the reply. The
// benchmark's judge uses it.
func (s *Service) Complete(ctx context.Context, settings Settings, system, user string) (string, Usage, error) {
	if err := settings.validate(); err != nil {
		return "", Usage{}, err
	}
	model, err := s.ready(settings)
	if err != nil {
		return "", Usage{}, err
	}
	ctx, cancel := context.WithTimeout(ctx, time.Duration(settings.TurnTimeoutSeconds)*time.Second)
	defer cancel()
	reply, err := model.Start(Request{
		Model:     settings.Model,
		Effort:    settings.Effort,
		MaxTokens: settings.MaxOutputTokens,
		Shared:    system,
		Question:  user,
	}).Next(ctx, func(string) {})
	if err != nil {
		return "", Usage{}, err
	}
	if reply.StopReason == stopRefusal {
		return "", reply.Usage, fmt.Errorf("the model refused")
	}
	return reply.Text, reply.Usage, nil
}
