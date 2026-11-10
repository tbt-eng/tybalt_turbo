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

// openAIModel calls an API that speaks the OpenAI chat completions format:
// OpenAI itself, OpenRouter, and others. It uses plain HTTP, because the
// format is small and the services that copy it differ in small ways.
type openAIModel struct {
	client      *http.Client
	baseURL     string
	apiKey      string
	tokensField string
}

// NewOpenAIModel returns an OpenAI-format API as a Model. baseURL ends before
// /chat/completions, such as https://api.openai.com/v1.
func NewOpenAIModel(apiKey, baseURL, tokensField string) Model {
	return &openAIModel{
		// The turn's context sets the time limit.
		client:      &http.Client{},
		baseURL:     strings.TrimRight(baseURL, "/"),
		apiKey:      apiKey,
		tokensField: tokensField,
	}
}

// openAIMessage is one message of the conversation, in the API's shape.
type openAIMessage struct {
	Role       string           `json:"role"`
	Content    *string          `json:"content"`
	ToolCalls  []openAIToolCall `json:"tool_calls,omitempty"`
	ToolCallID string           `json:"tool_call_id,omitempty"`
}

type openAIToolCall struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

func textMessage(role, text string) openAIMessage {
	return openAIMessage{Role: role, Content: &text}
}

func (m *openAIModel) Start(request Request) Session {
	// The shared text is the first message, and it is the same bytes for every
	// user. OpenAI caches the start of a prompt by itself, so this order is all
	// that the cache needs. The user's own context follows it.
	messages := []openAIMessage{textMessage("system", request.Shared)}
	if request.Context != "" {
		messages = append(messages, textMessage("system", request.Context))
	}
	for _, exchange := range request.History {
		messages = append(messages, textMessage("user", exchange.Question), textMessage("assistant", exchange.Answer))
	}
	messages = append(messages, textMessage("user", request.Question))

	tools := make([]map[string]any, len(request.Tools))
	for i, tool := range request.Tools {
		tools[i] = map[string]any{
			"type": "function",
			"function": map[string]any{
				"name":        tool.Name,
				"description": tool.Description,
				// The API then sends only arguments that fit the schema.
				"strict": true,
				"parameters": map[string]any{
					"type":                 "object",
					"properties":           tool.Properties,
					"required":             tool.Required,
					"additionalProperties": false,
				},
			},
		}
	}

	body := map[string]any{
		"model":          request.Model,
		"stream":         true,
		"stream_options": map[string]any{"include_usage": true},
		m.tokensField:    request.MaxTokens,
	}
	if len(tools) > 0 {
		body["tools"] = tools
	}
	if request.Effort != "" {
		body["reasoning_effort"] = request.Effort
	}
	return &openAISession{model: m, body: body, messages: messages}
}

// Warm does nothing. OpenAI caches a prompt when it is used, and has no
// request that only fills the cache.
func (m *openAIModel) Warm(ctx context.Context, request Request) error { return nil }

type openAISession struct {
	model    *openAIModel
	body     map[string]any
	messages []openAIMessage
}

// openAIChunk is one event of the stream.
type openAIChunk struct {
	ID      string `json:"id"`
	Choices []struct {
		Delta struct {
			Content   string `json:"content"`
			ToolCalls []struct {
				Index    int    `json:"index"`
				ID       string `json:"id"`
				Function struct {
					Name      string `json:"name"`
					Arguments string `json:"arguments"`
				} `json:"function"`
			} `json:"tool_calls"`
		} `json:"delta"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage *struct {
		PromptTokens        int `json:"prompt_tokens"`
		CompletionTokens    int `json:"completion_tokens"`
		PromptTokensDetails struct {
			CachedTokens int `json:"cached_tokens"`
		} `json:"prompt_tokens_details"`
	} `json:"usage"`
	// Some services report a failure inside the stream.
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (s *openAISession) Next(ctx context.Context, onText func(string)) (*Reply, error) {
	response, err := s.send(ctx)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	reply := &Reply{RequestID: response.Header.Get("x-request-id")}
	var calls []openAIToolCall
	finished := false

	lines := bufio.NewScanner(response.Body)
	lines.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	for lines.Scan() {
		// A line that is not data is a comment or a blank line between events.
		data, isData := strings.CutPrefix(lines.Text(), "data:")
		data = strings.TrimSpace(data)
		if !isData || data == "" {
			continue
		}
		if data == "[DONE]" {
			break
		}
		var chunk openAIChunk
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			return nil, fmt.Errorf("the model's stream holds bad JSON: %w", err)
		}
		if chunk.Error != nil {
			return nil, fmt.Errorf("the model's API failed: %s", truncate(chunk.Error.Message, maxErrorChars))
		}
		if reply.RequestID == "" {
			reply.RequestID = chunk.ID
		}
		if chunk.Usage != nil {
			cached := chunk.Usage.PromptTokensDetails.CachedTokens
			reply.Usage = Usage{
				Input:     chunk.Usage.PromptTokens - cached,
				CacheRead: cached,
				Output:    chunk.Usage.CompletionTokens,
			}
		}
		for _, choice := range chunk.Choices {
			if choice.Delta.Content != "" {
				reply.Text += choice.Delta.Content
				onText(choice.Delta.Content)
			}
			// A tool call arrives in pieces. The first piece has the id and the
			// name, and the arguments grow piece by piece.
			for _, piece := range choice.Delta.ToolCalls {
				for len(calls) <= piece.Index {
					calls = append(calls, openAIToolCall{Type: "function"})
				}
				call := &calls[piece.Index]
				if piece.ID != "" {
					call.ID = piece.ID
				}
				call.Function.Name += piece.Function.Name
				call.Function.Arguments += piece.Function.Arguments
			}
			if choice.FinishReason != "" {
				reply.StopReason = choice.FinishReason
				finished = true
			}
		}
	}
	if err := lines.Err(); err != nil {
		return nil, err
	}
	if !finished {
		return nil, errors.New("the model's stream ended before the answer did")
	}

	// The model's reply goes back to it on the next request.
	answer := openAIMessage{Role: "assistant", ToolCalls: calls}
	if reply.Text != "" {
		answer.Content = &reply.Text
	}
	s.messages = append(s.messages, answer)

	for _, call := range calls {
		reply.ToolCalls = append(reply.ToolCalls, ToolCall{ID: call.ID, Name: call.Function.Name, Input: json.RawMessage(call.Function.Arguments)})
	}
	// The turn loop knows two reasons by name. Some services say "stop" when
	// the reply holds tool calls, so the calls decide.
	switch {
	case len(calls) > 0:
		reply.StopReason = stopToolUse
	case reply.StopReason == "content_filter":
		reply.StopReason = stopRefusal
	}
	return reply, nil
}

func (s *openAISession) AddToolResults(results []ToolResult) {
	for _, result := range results {
		content := result.Content
		if result.IsError {
			content = "Error: " + content
		}
		s.messages = append(s.messages, openAIMessage{Role: "tool", ToolCallID: result.ID, Content: &content})
	}
}

// send posts the conversation. It tries one more time when the service is busy
// and has sent nothing yet.
func (s *openAISession) send(ctx context.Context) (*http.Response, error) {
	s.body["messages"] = s.messages
	payload, err := json.Marshal(s.body)
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
		request, err := http.NewRequestWithContext(ctx, http.MethodPost, s.model.baseURL+"/chat/completions", bytes.NewReader(payload))
		if err != nil {
			return nil, err
		}
		request.Header.Set("Authorization", "Bearer "+s.model.apiKey)
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Accept", "text/event-stream")
		response, err := s.model.client.Do(request)
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

// apiErrorMessage reads the message from an error body, which is
// {"error": {"message": "..."}} in this format.
func apiErrorMessage(body []byte) string {
	var parsed struct {
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if json.Unmarshal(body, &parsed) == nil && parsed.Error.Message != "" {
		return truncate(parsed.Error.Message, maxErrorChars)
	}
	return truncate(strings.TrimSpace(string(body)), maxErrorChars)
}
