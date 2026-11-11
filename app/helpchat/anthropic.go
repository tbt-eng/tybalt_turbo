package helpchat

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// anthropicModel calls the Claude Messages API with plain HTTP, as openai.go
// does for the OpenAI format.
//
// It does not use Anthropic's Go SDK. The compiler needs about 2.4 GB of memory
// for that one package, which is more than a small build machine has, and the
// chat uses one endpoint of it.
type anthropicModel struct {
	client  *http.Client
	baseURL string
	apiKey  string
}

const (
	anthropicURL = "https://api.anthropic.com"
	// The version of the API that the request and the stream below are written for.
	anthropicVersion = "2026-10-21"
)

// NewAnthropicModel returns the Claude API as a Model. An empty baseURL is
// Anthropic's own address.
func NewAnthropicModel(apiKey, baseURL string) Model {
	if baseURL == "" {
		baseURL = anthropicURL
	}
	return &anthropicModel{
		// The turn's context sets the time limit.
		client:  &http.Client{},
		baseURL: strings.TrimRight(baseURL, "/"),
		apiKey:  apiKey,
	}
}

// The request, in the API's shape. The fields are in the order that the API
// reads them: tools, then system, then messages.
type claudeRequest struct {
	Model        string          `json:"model"`
	MaxTokens    int             `json:"max_tokens"`
	Tools        []claudeTool    `json:"tools,omitempty"`
	System       []claudeText    `json:"system"`
	Messages     []claudeMessage `json:"messages"`
	OutputConfig *claudeOutput   `json:"output_config,omitempty"`
	Stream       bool            `json:"stream,omitempty"`
}

type claudeText struct {
	Type         string        `json:"type"`
	Text         string        `json:"text"`
	CacheControl *claudeMarker `json:"cache_control,omitempty"`
}

type claudeMarker struct {
	Type string `json:"type"`
}

type claudeTool struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	InputSchema claudeSchema `json:"input_schema"`
	// The API then sends only input that fits the schema.
	Strict bool `json:"strict"`
}

type claudeSchema struct {
	Type                 string         `json:"type"`
	Properties           map[string]any `json:"properties"`
	Required             []string       `json:"required,omitempty"`
	AdditionalProperties bool           `json:"additionalProperties"`
}

type claudeOutput struct {
	Effort string `json:"effort"`
}

// claudeMessage is one message. Each item of Content is a block in the API's
// shape: text, tool_use, tool_result, or a block of the model's thinking.
type claudeMessage struct {
	Role    string `json:"role"`
	Content []any  `json:"content"`
}

func textBlock(text string) map[string]any {
	return map[string]any{"type": "text", "text": text}
}

// body builds the request. The first cache mark is on the shared system block,
// so the tools and that block are one cached prefix for all users. The second
// is on the common block, which is also the same for all users, but changes
// when an admin changes a claim or a setting. Then the API still reads the
// first prefix from the cache and writes the common block again.
func (m *anthropicModel) body(request Request) claudeRequest {
	system := []claudeText{{Type: "text", Text: request.Shared, CacheControl: &claudeMarker{Type: "ephemeral"}}}
	if request.Common != "" {
		system = append(system, claudeText{Type: "text", Text: request.Common, CacheControl: &claudeMarker{Type: "ephemeral"}})
	}
	if request.Context != "" {
		system = append(system, claudeText{Type: "text", Text: request.Context})
	}

	tools := make([]claudeTool, len(request.Tools))
	for i, tool := range request.Tools {
		properties := tool.Properties
		if properties == nil {
			properties = map[string]any{}
		}
		tools[i] = claudeTool{
			Name:        tool.Name,
			Description: tool.Description,
			InputSchema: claudeSchema{Type: "object", Properties: properties, Required: tool.Required},
			Strict:      true,
		}
	}

	var messages []claudeMessage
	for _, exchange := range request.History {
		messages = append(messages,
			claudeMessage{Role: "user", Content: []any{textBlock(exchange.Question)}},
			claudeMessage{Role: "assistant", Content: []any{textBlock(exchange.Answer)}},
		)
	}
	messages = append(messages, claudeMessage{Role: "user", Content: []any{textBlock(request.Question)}})

	body := claudeRequest{
		Model:     request.Model,
		MaxTokens: request.MaxTokens,
		Tools:     tools,
		System:    system,
		Messages:  messages,
	}
	if request.Effort != "" {
		body.OutputConfig = &claudeOutput{Effort: request.Effort}
	}
	return body
}

func (m *anthropicModel) Start(request Request) Session {
	return &anthropicSession{model: m, body: m.body(request)}
}

// Warm sends the shared prefix with no room for an answer. The API writes the
// cache and returns at once.
func (m *anthropicModel) Warm(ctx context.Context, request Request) error {
	request.Context = ""
	request.History = nil
	request.Question = "Hello"
	body := m.body(request)
	body.MaxTokens = 0
	response, err := m.send(ctx, body)
	if err != nil {
		return err
	}
	response.Body.Close()
	return nil
}

type anthropicSession struct {
	model *anthropicModel
	body  claudeRequest
}

// claudeEvent is one event of the stream. The type says which fields it has.
type claudeEvent struct {
	Type    string `json:"type"`
	Index   int    `json:"index"`
	Message *struct {
		ID    string      `json:"id"`
		Usage claudeUsage `json:"usage"`
	} `json:"message"`
	// ContentBlock starts a block: text, tool_use, thinking or redacted_thinking.
	ContentBlock map[string]any `json:"content_block"`
	Delta        *struct {
		Type        string `json:"type"`
		Text        string `json:"text"`
		PartialJSON string `json:"partial_json"`
		Thinking    string `json:"thinking"`
		Signature   string `json:"signature"`
		StopReason  string `json:"stop_reason"`
	} `json:"delta"`
	Usage *claudeUsage `json:"usage"`
	Error *struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error"`
}

type claudeUsage struct {
	InputTokens              int `json:"input_tokens"`
	CacheCreationInputTokens int `json:"cache_creation_input_tokens"`
	CacheReadInputTokens     int `json:"cache_read_input_tokens"`
	OutputTokens             int `json:"output_tokens"`
}

// claudeBlock is one block of the reply while it arrives.
type claudeBlock struct {
	start map[string]any
	// The pieces that the deltas add.
	text, json, thinking, signature strings.Builder
}

func (b *claudeBlock) kind() string {
	kind, _ := b.start["type"].(string)
	return kind
}

// param returns the block in the shape that sends it back to the model. The
// model needs its own reply, with its thinking blocks, to go on after a tool
// call. ok is false for a block that must not go back, which is an empty text.
func (b *claudeBlock) param() (block any, ok bool) {
	switch b.kind() {
	case "text":
		text, _ := b.start["text"].(string)
		text += b.text.String()
		return textBlock(text), text != ""
	case "tool_use":
		return map[string]any{"type": "tool_use", "id": b.start["id"], "name": b.start["name"], "input": b.input()}, true
	case "thinking":
		thinking, _ := b.start["thinking"].(string)
		signature, _ := b.start["signature"].(string)
		return map[string]any{"type": "thinking", "thinking": thinking + b.thinking.String(), "signature": signature + b.signature.String()}, true
	}
	// A redacted_thinking block, or a kind that this code does not know, goes
	// back as it came.
	return b.start, true
}

// input is the input of a tool call, as the JSON that the model wrote.
func (b *claudeBlock) input() json.RawMessage {
	if written := strings.TrimSpace(b.json.String()); written != "" {
		return json.RawMessage(written)
	}
	if given, err := json.Marshal(b.start["input"]); err == nil && b.start["input"] != nil {
		return given
	}
	return json.RawMessage(`{}`)
}

func (s *anthropicSession) Next(ctx context.Context, onText func(string)) (*Reply, error) {
	body := s.body
	body.Stream = true
	response, err := s.model.send(ctx, body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	reply := &Reply{RequestID: response.Header.Get("request-id")}
	var blocks []*claudeBlock
	block := func(index int) *claudeBlock {
		for len(blocks) <= index {
			blocks = append(blocks, &claudeBlock{start: map[string]any{}})
		}
		return blocks[index]
	}
	messageID := ""
	finished := false

	lines := bufio.NewScanner(response.Body)
	lines.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	for lines.Scan() {
		// A line that is not data names the event, or is the blank line between
		// two events. The data holds the name too.
		data, isData := strings.CutPrefix(lines.Text(), "data:")
		data = strings.TrimSpace(data)
		if !isData || data == "" {
			continue
		}
		var event claudeEvent
		if err := json.Unmarshal([]byte(data), &event); err != nil {
			return nil, fmt.Errorf("the model's stream holds bad JSON: %w", err)
		}
		switch event.Type {
		case "message_start":
			if event.Message != nil {
				messageID = event.Message.ID
				reply.Usage = Usage{
					Input:      event.Message.Usage.InputTokens,
					CacheRead:  event.Message.Usage.CacheReadInputTokens,
					CacheWrite: event.Message.Usage.CacheCreationInputTokens,
					Output:     event.Message.Usage.OutputTokens,
				}
			}
		case "content_block_start":
			if event.Index >= 0 && event.ContentBlock != nil {
				block(event.Index).start = event.ContentBlock
			}
		case "content_block_delta":
			if event.Delta == nil || event.Index < 0 {
				continue
			}
			current := block(event.Index)
			switch event.Delta.Type {
			case "text_delta":
				current.text.WriteString(event.Delta.Text)
				if event.Delta.Text != "" {
					onText(event.Delta.Text)
				}
			case "input_json_delta":
				current.json.WriteString(event.Delta.PartialJSON)
			case "thinking_delta":
				current.thinking.WriteString(event.Delta.Thinking)
			case "signature_delta":
				current.signature.WriteString(event.Delta.Signature)
			}
		case "message_delta":
			if event.Delta != nil && event.Delta.StopReason != "" {
				reply.StopReason = event.Delta.StopReason
			}
			// The count of output tokens is the total so far. The other counts
			// are here only in some replies, and then they replace the first ones.
			if event.Usage != nil {
				reply.Usage.Output = event.Usage.OutputTokens
				if event.Usage.InputTokens > 0 {
					reply.Usage.Input = event.Usage.InputTokens
				}
				if event.Usage.CacheReadInputTokens > 0 {
					reply.Usage.CacheRead = event.Usage.CacheReadInputTokens
				}
				if event.Usage.CacheCreationInputTokens > 0 {
					reply.Usage.CacheWrite = event.Usage.CacheCreationInputTokens
				}
			}
		case "message_stop":
			finished = true
		case "error":
			message := "no message"
			if event.Error != nil {
				message = event.Error.Type + ": " + event.Error.Message
			}
			return nil, fmt.Errorf("the model's API failed: %s", truncate(message, maxErrorChars))
		}
	}
	if err := lines.Err(); err != nil {
		return nil, err
	}
	if !finished {
		return nil, errors.New("the model's stream ended before the answer did")
	}
	if reply.RequestID == "" {
		reply.RequestID = messageID
	}

	// The reply goes back to the model on the next request.
	var content []any
	for _, b := range blocks {
		switch b.kind() {
		case "text":
			text, _ := b.start["text"].(string)
			reply.Text += text + b.text.String()
		case "tool_use":
			id, _ := b.start["id"].(string)
			name, _ := b.start["name"].(string)
			reply.ToolCalls = append(reply.ToolCalls, ToolCall{ID: id, Name: name, Input: b.input()})
		}
		if param, ok := b.param(); ok {
			content = append(content, param)
		}
	}
	if len(content) > 0 {
		s.body.Messages = append(s.body.Messages, claudeMessage{Role: "assistant", Content: content})
	}
	return reply, nil
}

func (s *anthropicSession) AddToolResults(results []ToolResult) {
	blocks := make([]any, len(results))
	for i, result := range results {
		blocks[i] = map[string]any{
			"type":        "tool_result",
			"tool_use_id": result.ID,
			"content":     []any{textBlock(result.Content)},
			"is_error":    result.IsError,
		}
	}
	// All the results of one reply go in one user message.
	s.body.Messages = append(s.body.Messages, claudeMessage{Role: "user", Content: blocks})
}

// send posts the request. A person waits for the answer, so it tries one more
// time, and no more, when the service is busy and has sent nothing yet.
func (m *anthropicModel) send(ctx context.Context, body claudeRequest) (*http.Response, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	var lastErr error
	for attempt := 0; attempt < 2; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(500 * time.Millisecond):
			}
		}
		request, err := http.NewRequestWithContext(ctx, http.MethodPost, m.baseURL+"/v1/messages", bytes.NewReader(payload))
		if err != nil {
			return nil, err
		}
		request.Header.Set("x-api-key", m.apiKey)
		request.Header.Set("anthropic-version", anthropicVersion)
		request.Header.Set("Content-Type", "application/json")
		if body.Stream {
			request.Header.Set("Accept", "text/event-stream")
		}
		response, err := m.client.Do(request)
		if err != nil {
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			lastErr = err
			continue
		}
		if response.StatusCode == http.StatusOK {
			return response, nil
		}
		detail, _ := io.ReadAll(io.LimitReader(response.Body, 4096))
		response.Body.Close()
		lastErr = fmt.Errorf("the model's API returned %d: %s", response.StatusCode, apiErrorMessage(detail))
		busy := response.StatusCode == http.StatusTooManyRequests || response.StatusCode >= 500
		if !busy {
			break
		}
	}
	return nil, lastErr
}
