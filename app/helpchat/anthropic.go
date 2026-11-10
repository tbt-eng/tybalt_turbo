package helpchat

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/anthropics/anthropic-sdk-go/v2"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
)

// anthropicModel calls the Claude API. This is the only file that uses the SDK.
type anthropicModel struct {
	client anthropic.Client
}

// NewAnthropicModel returns the Claude API as a Model.
func NewAnthropicModel(apiKey string, options ...option.RequestOption) Model {
	options = append([]option.RequestOption{
		option.WithAPIKey(apiKey),
		// A person waits for the answer, so one retry is enough.
		option.WithMaxRetries(1),
	}, options...)
	return &anthropicModel{client: anthropic.NewClient(options...)}
}

// params builds the request. The order is fixed by the API: tools, then
// system, then messages. The cache mark is on the shared system block, so the
// tools and that block are one cached prefix for all users.
func (m *anthropicModel) params(request Request) anthropic.MessageNewParams {
	system := []anthropic.TextBlockParam{{
		Text:         request.Shared,
		CacheControl: anthropic.NewCacheControlEphemeralParam(),
	}}
	if request.Context != "" {
		system = append(system, anthropic.TextBlockParam{Text: request.Context})
	}

	tools := make([]anthropic.ToolUnionParam, len(request.Tools))
	for i, tool := range request.Tools {
		tools[i] = anthropic.ToolUnionParam{OfTool: &anthropic.ToolParam{
			Name:        tool.Name,
			Description: anthropic.String(tool.Description),
			InputSchema: anthropic.ToolInputSchemaParam{
				Properties:  tool.Properties,
				Required:    tool.Required,
				ExtraFields: map[string]any{"additionalProperties": false},
			},
			// The API then sends only input that fits the schema.
			Strict: anthropic.Bool(true),
		}}
	}

	var messages []anthropic.MessageParam
	for _, exchange := range request.History {
		messages = append(messages,
			anthropic.NewUserMessage(anthropic.NewTextBlock(exchange.Question)),
			anthropic.NewAssistantMessage(anthropic.NewTextBlock(exchange.Answer)),
		)
	}
	messages = append(messages, anthropic.NewUserMessage(anthropic.NewTextBlock(request.Question)))

	return anthropic.MessageNewParams{
		Model:        request.Model,
		MaxTokens:    int64(request.MaxTokens),
		System:       system,
		Tools:        tools,
		Messages:     messages,
		OutputConfig: anthropic.OutputConfigParam{Effort: anthropic.OutputConfigEffort(request.Effort)},
	}
}

func (m *anthropicModel) Start(request Request) Session {
	return &anthropicSession{client: m.client, params: m.params(request)}
}

// Warm sends the shared prefix with no room for an answer. The API writes the
// cache and returns at once.
func (m *anthropicModel) Warm(ctx context.Context, request Request) error {
	request.Context = ""
	request.History = nil
	request.Question = "Hello"
	params := m.params(request)
	params.MaxTokens = 0
	_, err := m.client.Messages.New(ctx, params)
	return err
}

type anthropicSession struct {
	client anthropic.Client
	params anthropic.MessageNewParams
}

func (s *anthropicSession) Next(ctx context.Context, onText func(string)) (*Reply, error) {
	var response *http.Response
	stream := s.client.Messages.NewStreaming(ctx, s.params, option.WithResponseInto(&response))
	message := anthropic.Message{}
	for stream.Next() {
		event := stream.Current()
		if err := message.Accumulate(event); err != nil {
			return nil, err
		}
		if delta, ok := event.AsAny().(anthropic.ContentBlockDeltaEvent); ok {
			if text, ok := delta.Delta.AsAny().(anthropic.TextDelta); ok && text.Text != "" {
				onText(text.Text)
			}
		}
	}
	if err := stream.Err(); err != nil {
		return nil, err
	}

	// The reply goes back to the model unchanged on the next request. The model
	// needs its own thinking blocks to continue after a tool call.
	s.params.Messages = append(s.params.Messages, message.ToParam())

	reply := &Reply{
		StopReason: string(message.StopReason),
		RequestID:  message.ID,
		Usage: Usage{
			Input:      int(message.Usage.InputTokens),
			CacheRead:  int(message.Usage.CacheReadInputTokens),
			CacheWrite: int(message.Usage.CacheCreationInputTokens),
			Output:     int(message.Usage.OutputTokens),
		},
	}
	if response != nil && response.Header.Get("request-id") != "" {
		reply.RequestID = response.Header.Get("request-id")
	}
	for _, block := range message.Content {
		switch block := block.AsAny().(type) {
		case anthropic.TextBlock:
			reply.Text += block.Text
		case anthropic.ToolUseBlock:
			reply.ToolCalls = append(reply.ToolCalls, ToolCall{
				ID:    block.ID,
				Name:  block.Name,
				Input: json.RawMessage(block.JSON.Input.Raw()),
			})
		}
	}
	return reply, nil
}

func (s *anthropicSession) AddToolResults(results []ToolResult) {
	blocks := make([]anthropic.ContentBlockParamUnion, len(results))
	for i, result := range results {
		blocks[i] = anthropic.NewToolResultBlock(result.ID, result.Content, result.IsError)
	}
	// All the results of one reply go in one user message.
	s.params.Messages = append(s.params.Messages, anthropic.NewUserMessage(blocks...))
}
