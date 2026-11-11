package helpchat

import (
	"context"
	"encoding/json"
)

// Model is the language model. The chat talks to it only through this
// interface, so tests run with a stub and each API's shapes stay in one file.
type Model interface {
	// Start begins one turn. Nothing is sent until Next.
	Start(request Request) Session
	// Warm writes the shared part of the prompt to the model's cache, so the
	// next question starts faster. It asks for no answer.
	Warm(ctx context.Context, request Request) error
}

// Session is one turn in progress. It holds what the model said and the tool
// results so far, and sends all of it again on each Next.
type Session interface {
	// Next sends the conversation and streams the reply. It calls onText for
	// each piece of answer text.
	Next(ctx context.Context, onText func(text string)) (*Reply, error)
	// AddToolResults answers the tool calls of the last reply.
	AddToolResults(results []ToolResult)
}

// Request is everything the model gets for a turn.
type Request struct {
	Model     string
	Effort    string
	MaxTokens int
	// Shared is the system rules and the handbook. It is the same bytes for
	// every user, so the model caches it. Only a deploy changes it.
	Shared string
	// Common is what the database holds for every user alike: all the claims
	// in Turbo, and Turbo's settings. It follows Shared and has its own cache
	// mark, so a change to a claim or a setting rewrites this part only.
	Common string
	// Context describes this user and this request. It follows Shared and
	// Common, so it does not break the cache.
	Context string
	Tools   []Tool
	// History is the earlier turns of this chat, oldest first.
	History  []Exchange
	Question string
}

// Tool describes one tool to the model.
type Tool struct {
	Name        string
	Description string
	// Properties and Required are the JSON schema of the input object.
	Properties map[string]any
	Required   []string
}

// Exchange is one earlier question and its answer.
type Exchange struct {
	Question string
	Answer   string
}

// Reply is one model reply. It ends the turn unless it holds tool calls.
type Reply struct {
	Text       string
	ToolCalls  []ToolCall
	StopReason string
	Usage      Usage
	RequestID  string
}

// ToolCall is the model's request to run a tool.
type ToolCall struct {
	ID    string
	Name  string
	Input json.RawMessage
}

// ToolResult is the server's answer to a tool call.
type ToolResult struct {
	ID      string
	Content string
	IsError bool
}

// Usage counts tokens. Cached tokens cost less than Input tokens.
type Usage struct {
	Input      int
	CacheRead  int
	CacheWrite int
	Output     int
}

func (u *Usage) add(other Usage) {
	u.Input += other.Input
	u.CacheRead += other.CacheRead
	u.CacheWrite += other.CacheWrite
	u.Output += other.Output
}

// The reasons a reply stops that the turn loop acts on.
const (
	stopToolUse = "tool_use"
	stopRefusal = "refusal"
)
