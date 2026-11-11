package helpchat

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// fakeChat is a server that speaks the OpenAI chat completions format.
type fakeChat struct {
	server  *httptest.Server
	bodies  []map[string]any
	headers []http.Header
	replies []chatReply
}

type chatReply struct {
	status int
	body   string
}

func newFakeChat(t *testing.T, replies ...chatReply) *fakeChat {
	fake := &fakeChat{replies: replies}
	fake.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		var body map[string]any
		if err := json.Unmarshal(raw, &body); err != nil {
			t.Errorf("the request body is not JSON: %v", err)
		}
		fake.bodies = append(fake.bodies, body)
		fake.headers = append(fake.headers, r.Header.Clone())
		if r.URL.Path != "/v1/chat/completions" || len(fake.replies) == 0 {
			http.Error(w, `{"error":{"message":"no reply left"}}`, http.StatusNotFound)
			return
		}
		reply := fake.replies[0]
		fake.replies = fake.replies[1:]
		w.Header().Set("x-request-id", fmt.Sprintf("req_chat_%d", len(fake.bodies)))
		if reply.status != 0 && reply.status != http.StatusOK {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(reply.status)
			io.WriteString(w, reply.body)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		io.WriteString(w, reply.body)
	}))
	t.Cleanup(fake.server.Close)
	return fake
}

func (fake *fakeChat) url() string { return fake.server.URL + "/v1" }

func (fake *fakeChat) model(tokensField string) Model {
	return NewOpenAIModel("test-key", fake.url(), tokensField)
}

// chatStream writes chunks as the API does, and ends the stream.
func chatStream(chunks ...string) chatReply {
	var out strings.Builder
	for _, chunk := range chunks {
		if strings.HasPrefix(chunk, ":") {
			out.WriteString(chunk + "\n\n")
			continue
		}
		fmt.Fprintf(&out, "data: %s\n\n", chunk)
	}
	out.WriteString("data: [DONE]\n\n")
	return chatReply{body: out.String()}
}

var chatToolReply = chatStream(
	`: OPENROUTER PROCESSING`,
	`{"id":"chatcmpl-1","choices":[{"index":0,"delta":{"role":"assistant","content":null,"tool_calls":[{"index":0,"id":"call_1","type":"function","function":{"name":"search_turbo","arguments":""}}]},"finish_reason":null}]}`,
	`{"id":"chatcmpl-1","choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"function":{"arguments":"{\"query\":"}}]},"finish_reason":null}]}`,
	`{"id":"chatcmpl-1","choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"function":{"arguments":" \"time entry\"}"}}]},"finish_reason":null}]}`,
	`{"id":"chatcmpl-1","choices":[{"index":0,"delta":{},"finish_reason":"tool_calls"}]}`,
	`{"id":"chatcmpl-1","choices":[],"usage":{"prompt_tokens":9120,"completion_tokens":25,"total_tokens":9145,"prompt_tokens_details":{"cached_tokens":9000}}}`,
)

var chatTextReply = chatStream(
	`{"id":"chatcmpl-2","choices":[{"index":0,"delta":{"role":"assistant","content":""},"finish_reason":null}]}`,
	`{"id":"chatcmpl-2","choices":[{"index":0,"delta":{"content":"Select "},"finish_reason":null}]}`,
	`{"id":"chatcmpl-2","choices":[{"index":0,"delta":{"content":"**Submit**."},"finish_reason":null}]}`,
	`{"id":"chatcmpl-2","choices":[{"index":0,"delta":{},"finish_reason":"stop"}]}`,
	`{"id":"chatcmpl-2","choices":[],"usage":{"prompt_tokens":9400,"completion_tokens":12,"total_tokens":9412}}`,
)

func openAIRequest() Request {
	request := testRequest()
	request.Model = "gpt-test"
	return request
}

func TestTheOpenAIRequestHasTheSharedPrefixAndTheTools(t *testing.T) {
	fake := newFakeChat(t, chatToolReply, chatTextReply)
	session := fake.model("max_completion_tokens").Start(openAIRequest())

	var streamed []string
	reply, err := session.Next(context.Background(), func(text string) { streamed = append(streamed, text) })
	if err != nil {
		t.Fatal(err)
	}
	if reply.StopReason != stopToolUse || reply.Text != "" || len(streamed) != 0 || len(reply.ToolCalls) != 1 {
		t.Fatalf("reply = %+v", reply)
	}
	toolCall := reply.ToolCalls[0]
	if toolCall.ID != "call_1" || toolCall.Name != toolSearch || string(toolCall.Input) != `{"query": "time entry"}` {
		t.Errorf("tool call = %+v", toolCall)
	}
	// Cached tokens are part of prompt_tokens in this format. The log keeps them apart.
	if reply.Usage != (Usage{Input: 120, CacheRead: 9000, Output: 25}) || reply.RequestID != "req_chat_1" {
		t.Errorf("usage %+v, request id %q", reply.Usage, reply.RequestID)
	}

	body := fake.bodies[0]
	for path, want := range map[string]any{
		"model":                            "gpt-test",
		"stream":                           true,
		"stream_options.include_usage":     true,
		"max_completion_tokens":            float64(8000),
		"reasoning_effort":                 "low",
		"messages.0.role":                  "system",
		"messages.0.content":               "the rules and the handbook",
		"messages.1.role":                  "system",
		"messages.1.content":               "the user holds the time claim",
		"messages.2.role":                  "user",
		"messages.2.content":               "An earlier question?",
		"messages.3.role":                  "assistant",
		"messages.3.content":               "An earlier answer.",
		"messages.4.role":                  "user",
		"messages.4.content":               "How do I submit my week?",
		"tools.0.type":                     "function",
		"tools.0.function.name":            "search_turbo",
		"tools.0.function.strict":          true,
		"tools.1.function.name":            "read_source",
		"tools.2.function.name":            "report_gap",
		"tools.0.function.parameters.type": "object",
		"tools.0.function.parameters.additionalProperties": false,
		"tools.0.function.parameters.required.0":           "query",
	} {
		if got := at(t, body, path); got != want {
			t.Errorf("%s = %v, want %v", path, got, want)
		}
	}
	if len(at(t, body, "messages").([]any)) != 5 || len(at(t, body, "tools").([]any)) != 3 {
		t.Error("the body has parts that the chat did not set")
	}
	for _, key := range []string{"max_tokens", "temperature", "top_p", "tool_choice"} {
		if _, set := body[key]; set {
			t.Errorf("the body sets %s", key)
		}
	}
	if got := fake.headers[0].Get("Authorization"); got != "Bearer test-key" {
		t.Errorf("authorization header = %q", got)
	}

	// The second request sends the model's tool call back, then the result.
	session.AddToolResults([]ToolResult{
		{ID: "call_1", Content: "<chunk ref=\"S1\">text</chunk>"},
	})
	reply, err = session.Next(context.Background(), func(text string) { streamed = append(streamed, text) })
	if err != nil {
		t.Fatal(err)
	}
	if reply.Text != "Select **Submit**." || strings.Join(streamed, "|") != "Select |**Submit**." || reply.StopReason != "stop" || len(reply.ToolCalls) != 0 {
		t.Errorf("reply = %+v, streamed %q", reply, streamed)
	}
	if reply.Usage != (Usage{Input: 9400, Output: 12}) {
		t.Errorf("usage = %+v", reply.Usage)
	}
	body = fake.bodies[1]
	for path, want := range map[string]any{
		"messages.0.content":                         "the rules and the handbook",
		"messages.5.role":                            "assistant",
		"messages.5.content":                         nil,
		"messages.5.tool_calls.0.id":                 "call_1",
		"messages.5.tool_calls.0.type":               "function",
		"messages.5.tool_calls.0.function.name":      "search_turbo",
		"messages.5.tool_calls.0.function.arguments": `{"query": "time entry"}`,
		"messages.6.role":                            "tool",
		"messages.6.tool_call_id":                    "call_1",
		"messages.6.content":                         "<chunk ref=\"S1\">text</chunk>",
	} {
		if got := at(t, body, path); got != want {
			t.Errorf("second request: %s = %v, want %v", path, got, want)
		}
	}
}

func TestOpenAIFormatDifferencesBetweenServices(t *testing.T) {
	fake := newFakeChat(t, chatTextReply)
	request := openAIRequest()
	request.Effort = ""
	request.Context = ""
	if _, err := fake.model("max_tokens").Start(request).Next(context.Background(), func(string) {}); err != nil {
		t.Fatal(err)
	}
	body := fake.bodies[0]
	if body["max_tokens"] != float64(8000) {
		t.Errorf("max_tokens = %v", body["max_tokens"])
	}
	for _, key := range []string{"max_completion_tokens", "reasoning_effort"} {
		if _, set := body[key]; set {
			t.Errorf("the body sets %s", key)
		}
	}
	if at(t, body, "messages.1.role") != "user" {
		t.Error("with no context, the earlier turns follow the shared message")
	}
}

func TestTheOpenAISharedPartIsTheSameForEveryUser(t *testing.T) {
	fake := newFakeChat(t, chatTextReply, chatTextReply)
	prefix := func(request Request) string {
		if _, err := fake.model("max_completion_tokens").Start(request).Next(context.Background(), func(string) {}); err != nil {
			t.Fatal(err)
		}
		body := fake.bodies[len(fake.bodies)-1]
		tools, _ := json.Marshal(body["tools"])
		shared, _ := json.Marshal(at(t, body, "messages.0"))
		return string(tools) + string(shared)
	}
	one := openAIRequest()
	two := openAIRequest()
	two.Context = "another user, with the admin claim"
	two.History = nil
	two.Question = "Another question?"
	if prefix(one) != prefix(two) {
		t.Error("two users must send the same tools and the same first message, or the cache serves neither")
	}
}

func TestOpenAIStopReasonsAndToolErrors(t *testing.T) {
	// Some services say "stop" when the reply holds a tool call.
	stopWithCall := chatStream(
		`{"id":"c","choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"id":"call_a","function":{"name":"report_gap","arguments":"{}"}},{"index":1,"id":"call_b","function":{"name":"search_turbo","arguments":"{\"query\":\"x\"}"}}]},"finish_reason":"stop"}]}`,
	)
	refused := chatStream(`{"id":"c","choices":[{"index":0,"delta":{},"finish_reason":"content_filter"}]}`)
	cut := chatStream(`{"id":"c","choices":[{"index":0,"delta":{"content":"Half an ans"},"finish_reason":"length"}]}`)
	fake := newFakeChat(t, stopWithCall, refused, cut)
	session := fake.model("max_tokens").Start(openAIRequest())

	reply, err := session.Next(context.Background(), func(string) {})
	if err != nil || reply.StopReason != stopToolUse || len(reply.ToolCalls) != 2 || reply.ToolCalls[1].ID != "call_b" || reply.RequestID == "" {
		t.Fatalf("reply %+v, err %v", reply, err)
	}
	session.AddToolResults([]ToolResult{{ID: "call_a", Content: "Recorded."}, {ID: "call_b", Content: "The search failed.", IsError: true}})

	reply, err = session.Next(context.Background(), func(string) {})
	if err != nil || reply.StopReason != stopRefusal {
		t.Fatalf("reply %+v, err %v", reply, err)
	}
	// Each tool result is its own message, and an error says that it is one.
	body := fake.bodies[1]
	if at(t, body, "messages.6.tool_call_id") != "call_a" || at(t, body, "messages.7.tool_call_id") != "call_b" || at(t, body, "messages.7.content") != "Error: The search failed." {
		t.Errorf("tool messages: %v", at(t, body, "messages"))
	}

	reply, err = session.Next(context.Background(), func(string) {})
	if err != nil || reply.StopReason != "length" || reply.Text != "Half an ans" {
		t.Fatalf("reply %+v, err %v", reply, err)
	}
}

func TestOpenAIFailures(t *testing.T) {
	next := func(fake *fakeChat) (*Reply, error) {
		return fake.model("max_tokens").Start(openAIRequest()).Next(context.Background(), func(string) {})
	}

	t.Run("a bad key is an error with the service's message, and is not tried again", func(t *testing.T) {
		fake := newFakeChat(t, chatReply{status: 401, body: `{"error":{"message":"Incorrect API key provided.","type":"invalid_request_error"}}`}, chatTextReply)
		_, err := next(fake)
		if err == nil || !strings.Contains(err.Error(), "401") || !strings.Contains(err.Error(), "Incorrect API key provided.") || strings.Contains(err.Error(), "test-key") {
			t.Fatalf("err = %v", err)
		}
		if len(fake.bodies) != 1 {
			t.Errorf("the request was sent %d times", len(fake.bodies))
		}
	})

	t.Run("a busy service gets one more try", func(t *testing.T) {
		fake := newFakeChat(t, chatReply{status: 503, body: `overloaded`}, chatTextReply)
		reply, err := next(fake)
		if err != nil || reply.Text != "Select **Submit**." || len(fake.bodies) != 2 {
			t.Fatalf("reply %+v, err %v, sent %d times", reply, err, len(fake.bodies))
		}
		busy := newFakeChat(t, chatReply{status: 429, body: `{"error":{"message":"Rate limit reached."}}`}, chatReply{status: 429, body: `{"error":{"message":"Rate limit reached."}}`}, chatTextReply)
		if _, err := next(busy); err == nil || !strings.Contains(err.Error(), "Rate limit reached.") || len(busy.bodies) != 2 {
			t.Fatalf("err %v, sent %d times", err, len(busy.bodies))
		}
	})

	t.Run("an error inside the stream", func(t *testing.T) {
		fake := newFakeChat(t, chatStream(
			`{"id":"c","choices":[{"index":0,"delta":{"content":"Sel"},"finish_reason":null}]}`,
			`{"error":{"message":"The provider stopped.","code":502}}`,
		))
		if _, err := next(fake); err == nil || !strings.Contains(err.Error(), "The provider stopped.") {
			t.Fatalf("err = %v", err)
		}
	})

	t.Run("a stream that stops early", func(t *testing.T) {
		fake := newFakeChat(t, chatStream(`{"id":"c","choices":[{"index":0,"delta":{"content":"Sel"},"finish_reason":null}]}`))
		if _, err := next(fake); err == nil {
			t.Fatal("a stream with no end must be an error")
		}
		bad := newFakeChat(t, chatReply{body: "data: {not json\n\n"})
		if _, err := next(bad); err == nil {
			t.Fatal("bad JSON must be an error")
		}
	})

	t.Run("the user leaves", func(t *testing.T) {
		fake := newFakeChat(t, chatReply{status: 503, body: `overloaded`}, chatTextReply)
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()
		started := time.Now()
		_, err := fake.model("max_tokens").Start(openAIRequest()).Next(ctx, func(string) {})
		if err == nil || time.Since(started) > 400*time.Millisecond {
			t.Fatalf("err %v after %s", err, time.Since(started))
		}
	})

	if err := NewOpenAIModel("k", "http://127.0.0.1:1", "max_tokens").Warm(context.Background(), openAIRequest()); err != nil {
		t.Errorf("warm does nothing for this format, got %v", err)
	}
}

func TestTheOpenAICommonPartFollowsTheSharedPart(t *testing.T) {
	fake := newFakeChat(t, chatTextReply)
	request := openAIRequest()
	request.Common = "all the claims, and Turbo's settings"
	if _, err := fake.model("max_completion_tokens").Start(request).Next(context.Background(), func(string) {}); err != nil {
		t.Fatal(err)
	}
	body := fake.bodies[0]
	// These services cache the start of a prompt, so what is the same for
	// every user comes before what is of one user.
	for path, want := range map[string]any{
		"messages.0.role":    "system",
		"messages.1.role":    "system",
		"messages.1.content": "all the claims, and Turbo's settings",
		"messages.2.role":    "system",
		"messages.2.content": openAIRequest().Context,
	} {
		if got := at(t, body, path); got != want {
			t.Errorf("%s = %v, want %v", path, got, want)
		}
	}
}
