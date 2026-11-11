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
)

// fakeAPI is a server that speaks the Claude API's streaming format. It lets
// the tests check the bytes that the chat sends, with no key and no network.
type fakeAPI struct {
	t         *testing.T
	server    *httptest.Server
	bodies    []map[string]any
	headers   []http.Header
	responses []string
}

func newFakeAPI(t *testing.T, responses ...string) *fakeAPI {
	api := &fakeAPI{t: t, responses: responses}
	api.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		var body map[string]any
		if err := json.Unmarshal(raw, &body); err != nil {
			t.Errorf("the request body is not JSON: %v", err)
		}
		api.bodies = append(api.bodies, body)
		api.headers = append(api.headers, r.Header.Clone())
		if r.URL.Path != "/v1/messages" || len(api.responses) == 0 {
			http.Error(w, `{"type":"error","error":{"type":"not_found_error","message":"no response left"}}`, http.StatusNotFound)
			return
		}
		response := api.responses[0]
		api.responses = api.responses[1:]
		w.Header().Set("request-id", fmt.Sprintf("req_fake_%d", len(api.bodies)))
		if body["stream"] == true {
			w.Header().Set("Content-Type", "text/event-stream")
		} else {
			w.Header().Set("Content-Type", "application/json")
		}
		io.WriteString(w, response)
	}))
	t.Cleanup(api.server.Close)
	return api
}

func (api *fakeAPI) model() Model {
	return NewAnthropicModel("test-key", api.server.URL)
}

func sse(events ...string) string {
	var out strings.Builder
	for _, event := range events {
		var typed struct {
			Type string `json:"type"`
		}
		json.Unmarshal([]byte(event), &typed)
		fmt.Fprintf(&out, "event: %s\ndata: %s\n\n", typed.Type, event)
	}
	return out.String()
}

const messageStart = `{"type":"message_start","message":{"id":"msg_1","type":"message","role":"assistant","model":"claude-haiku-5-5","content":[],"stop_reason":null,"stop_sequence":null,"usage":{"input_tokens":120,"cache_creation_input_tokens":30,"cache_read_input_tokens":9000,"output_tokens":1}}}`

func messageEnd(stop string, output int) []string {
	return []string{
		fmt.Sprintf(`{"type":"message_delta","delta":{"stop_reason":%q,"stop_sequence":null},"usage":{"output_tokens":%d}}`, stop, output),
		`{"type":"message_stop"}`,
	}
}

var toolReply = sse(append([]string{
	messageStart,
	`{"type":"content_block_start","index":0,"content_block":{"type":"tool_use","id":"toolu_1","name":"search_turbo","input":{}}}`,
	`{"type":"content_block_delta","index":0,"delta":{"type":"input_json_delta","partial_json":"{\"query\": \"time"}}`,
	`{"type":"content_block_delta","index":0,"delta":{"type":"input_json_delta","partial_json":" entry\"}"}}`,
	`{"type":"content_block_stop","index":0}`,
}, messageEnd("tool_use", 25)...)...)

var textReply = sse(append([]string{
	messageStart,
	`{"type":"content_block_start","index":0,"content_block":{"type":"text","text":""}}`,
	`{"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"Select "}}`,
	`{"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"**Submit**."}}`,
	`{"type":"content_block_stop","index":0}`,
}, messageEnd("end_turn", 12)...)...)

func testRequest() Request {
	return Request{
		Model:     "claude-haiku-5-5",
		Effort:    "low",
		MaxTokens: 8000,
		Shared:    "the rules and the handbook",
		Context:   "the user holds the time claim",
		Tools:     tools,
		History:   []Exchange{{Question: "An earlier question?", Answer: "An earlier answer."}},
		Question:  "How do I submit my week?",
	}
}

// at reads a value from decoded JSON by a path such as "system.0.cache_control.type".
func at(t *testing.T, value any, path string) any {
	t.Helper()
	for _, key := range strings.Split(path, ".") {
		switch typed := value.(type) {
		case map[string]any:
			value = typed[key]
		case []any:
			var index int
			fmt.Sscan(key, &index)
			if index >= len(typed) {
				t.Fatalf("%s: no item %d", path, index)
			}
			value = typed[index]
		default:
			t.Fatalf("%s: nothing at %q", path, key)
		}
	}
	return value
}

func TestTheClaudeRequestHasTheCachedPrefixAndTheTools(t *testing.T) {
	api := newFakeAPI(t, toolReply, textReply)
	session := api.model().Start(testRequest())

	var streamed []string
	reply, err := session.Next(context.Background(), func(text string) { streamed = append(streamed, text) })
	if err != nil {
		t.Fatal(err)
	}
	if reply.StopReason != stopToolUse || reply.Text != "" || len(streamed) != 0 || len(reply.ToolCalls) != 1 {
		t.Fatalf("reply = %+v", reply)
	}
	toolCall := reply.ToolCalls[0]
	if toolCall.ID != "toolu_1" || toolCall.Name != toolSearch || string(toolCall.Input) != `{"query": "time entry"}` {
		t.Errorf("tool call = %+v", toolCall)
	}
	if reply.Usage != (Usage{Input: 120, CacheRead: 9000, CacheWrite: 30, Output: 25}) || reply.RequestID != "req_fake_1" {
		t.Errorf("usage %+v, request id %q", reply.Usage, reply.RequestID)
	}

	body := api.bodies[0]
	for path, want := range map[string]any{
		"model":                       "claude-haiku-5-5",
		"max_tokens":                  float64(8000),
		"stream":                      true,
		"output_config.effort":        "low",
		"system.0.text":               "the rules and the handbook",
		"system.0.cache_control.type": "ephemeral",
		"system.1.text":               "the user holds the time claim",
		"tools.0.name":                "search_turbo",
		"tools.0.strict":              true,
		"tools.0.input_schema.type":   "object",
		"tools.1.name":                "read_source",
		"tools.2.name":                "report_gap",
		"messages.0.role":             "user",
		"messages.0.content.0.text":   "An earlier question?",
		"messages.1.role":             "assistant",
		"messages.1.content.0.text":   "An earlier answer.",
		"messages.2.role":             "user",
		"messages.2.content.0.text":   "How do I submit my week?",
		"tools.0.input_schema.additionalProperties": false,
	} {
		if got := at(t, body, path); got != want {
			t.Errorf("%s = %v, want %v", path, got, want)
		}
	}
	// Only the shared block has a cache mark. A mark on the user's block would
	// write a cache entry for each user.
	if at(t, body, "system.1").(map[string]any)["cache_control"] != nil {
		t.Error("the user's block must have no cache mark")
	}
	if len(at(t, body, "system").([]any)) != 2 || len(at(t, body, "tools").([]any)) != 3 || len(at(t, body, "messages").([]any)) != 3 {
		t.Errorf("the body has parts that the chat did not set")
	}
	// These would be refused by the models that the chat uses.
	for _, key := range []string{"temperature", "top_p", "top_k", "tool_choice", "thinking"} {
		if _, set := body[key]; set {
			t.Errorf("the body sets %s", key)
		}
	}
	if key, version := api.headers[0].Get("x-api-key"), api.headers[0].Get("anthropic-version"); key != "test-key" || version != "2026-10-21" {
		t.Errorf("key header = %q, version header = %q", key, version)
	}

	// The second request sends the model's own reply back, then the tool result.
	session.AddToolResults([]ToolResult{{ID: "toolu_1", Content: "<chunk ref=\"S1\">text</chunk>"}})
	reply, err = session.Next(context.Background(), func(text string) { streamed = append(streamed, text) })
	if err != nil {
		t.Fatal(err)
	}
	if reply.Text != "Select **Submit**." || strings.Join(streamed, "|") != "Select |**Submit**." || reply.StopReason != "end_turn" || len(reply.ToolCalls) != 0 {
		t.Errorf("reply = %+v, streamed %q", reply, streamed)
	}
	body = api.bodies[1]
	for path, want := range map[string]any{
		"messages.3.role":                  "assistant",
		"messages.3.content.0.type":        "tool_use",
		"messages.3.content.0.id":          "toolu_1",
		"messages.3.content.0.input.query": "time entry",
		"messages.4.role":                  "user",
		"messages.4.content.0.type":        "tool_result",
		"messages.4.content.0.tool_use_id": "toolu_1",
		"system.0.cache_control.type":      "ephemeral",
		"system.0.text":                    "the rules and the handbook",
	} {
		if got := at(t, body, path); got != want {
			t.Errorf("second request: %s = %v, want %v", path, got, want)
		}
	}
	if result := fmt.Sprint(at(t, body, "messages.4.content.0.content")); !strings.Contains(result, `<chunk ref="S1">text</chunk>`) {
		t.Errorf("tool result content = %s", result)
	}
}

func TestTheSharedPartOfTheRequestIsTheSameForEveryUser(t *testing.T) {
	api := newFakeAPI(t, textReply, textReply)
	prefix := func(request Request) string {
		if _, err := api.model().Start(request).Next(context.Background(), func(string) {}); err != nil {
			t.Fatal(err)
		}
		body := api.bodies[len(api.bodies)-1]
		tools, _ := json.Marshal(body["tools"])
		shared, _ := json.Marshal(at(t, body, "system.0"))
		return string(tools) + string(shared)
	}
	one := testRequest()
	two := testRequest()
	two.Context = "another user, with the admin claim"
	two.History = nil
	two.Question = "Another question?"
	if prefix(one) != prefix(two) {
		t.Error("two users must send the same tools and the same shared block, or the cache serves neither")
	}
}

func TestWarmAsksForNoAnswer(t *testing.T) {
	api := newFakeAPI(t, `{"id":"msg_1","type":"message","role":"assistant","model":"claude-haiku-5-5","content":[],"stop_reason":"max_tokens","stop_sequence":null,"usage":{"input_tokens":3,"cache_creation_input_tokens":9000,"cache_read_input_tokens":0,"output_tokens":0}}`)
	if err := api.model().Warm(context.Background(), testRequest()); err != nil {
		t.Fatal(err)
	}
	body := api.bodies[0]
	if got, set := body["max_tokens"]; !set || got != float64(0) {
		t.Errorf("max_tokens = %v (set: %v), want 0", got, set)
	}
	if body["stream"] == true {
		t.Error("a warm request must not stream")
	}
	if at(t, body, "system.0.cache_control.type") != "ephemeral" || len(at(t, body, "system").([]any)) != 1 || len(at(t, body, "tools").([]any)) != 3 {
		t.Error("a warm request must hold the tools and the shared block, and nothing about one user")
	}
}

func TestAnErrorFromTheAPIIsAnError(t *testing.T) {
	api := newFakeAPI(t)
	reply, err := api.model().Start(testRequest()).Next(context.Background(), func(string) {})
	if err == nil || reply != nil {
		t.Fatalf("reply %+v, err %v", reply, err)
	}
}

// The model thinks before a tool call. Its thinking blocks must go back with
// the tool call, as they came, or the API refuses the next request.
func TestTheModelsThinkingGoesBackWithItsToolCall(t *testing.T) {
	thinkingReply := sse(append([]string{
		messageStart,
		`{"type":"content_block_start","index":0,"content_block":{"type":"thinking","thinking":"","signature":""}}`,
		`{"type":"content_block_delta","index":0,"delta":{"type":"thinking_delta","thinking":"Search for "}}`,
		`{"type":"content_block_delta","index":0,"delta":{"type":"thinking_delta","thinking":"time entries."}}`,
		`{"type":"content_block_delta","index":0,"delta":{"type":"signature_delta","signature":"sig-abc"}}`,
		`{"type":"content_block_stop","index":0}`,
		`{"type":"content_block_start","index":1,"content_block":{"type":"redacted_thinking","data":"opaque-data"}}`,
		`{"type":"content_block_stop","index":1}`,
		`{"type":"content_block_start","index":2,"content_block":{"type":"text","text":""}}`,
		`{"type":"content_block_stop","index":2}`,
		`{"type":"ping"}`,
		`{"type":"content_block_start","index":3,"content_block":{"type":"tool_use","id":"toolu_9","name":"report_gap","input":{}}}`,
		`{"type":"content_block_stop","index":3}`,
	}, messageEnd("tool_use", 40)...)...)
	api := newFakeAPI(t, thinkingReply, textReply)
	session := api.model().Start(testRequest())
	reply, err := session.Next(context.Background(), func(string) {})
	if err != nil {
		t.Fatal(err)
	}
	// A tool call with no input is an empty object, and the thinking is not text.
	if reply.Text != "" || len(reply.ToolCalls) != 1 || string(reply.ToolCalls[0].Input) != "{}" || reply.Usage.Output != 40 {
		t.Fatalf("reply = %+v", reply)
	}
	session.AddToolResults([]ToolResult{{ID: "toolu_9", Content: "Recorded.", IsError: true}})
	if _, err := session.Next(context.Background(), func(string) {}); err != nil {
		t.Fatal(err)
	}
	body := api.bodies[1]
	for path, want := range map[string]any{
		"messages.3.role":                  "assistant",
		"messages.3.content.0.type":        "thinking",
		"messages.3.content.0.thinking":    "Search for time entries.",
		"messages.3.content.0.signature":   "sig-abc",
		"messages.3.content.1.type":        "redacted_thinking",
		"messages.3.content.1.data":        "opaque-data",
		"messages.3.content.2.type":        "tool_use",
		"messages.3.content.2.id":          "toolu_9",
		"messages.4.content.0.type":        "tool_result",
		"messages.4.content.0.is_error":    true,
		"messages.4.content.0.tool_use_id": "toolu_9",
	} {
		if got := at(t, body, path); got != want {
			t.Errorf("%s = %v, want %v", path, got, want)
		}
	}
	// The empty text block is left out: the API refuses an empty text.
	if blocks := at(t, body, "messages.3.content").([]any); len(blocks) != 3 {
		t.Errorf("the reply went back as %d blocks: %v", len(blocks), blocks)
	}
	if input, ok := at(t, body, "messages.3.content.2.input").(map[string]any); !ok || len(input) != 0 {
		t.Errorf("the input of the tool call = %v", at(t, body, "messages.3.content.2.input"))
	}
}

func TestAStreamThatFailsOrStopsEarlyIsAnError(t *testing.T) {
	failed := sse(messageStart, `{"type":"error","error":{"type":"overloaded_error","message":"Overloaded"}}`)
	cut := sse(messageStart, `{"type":"content_block_start","index":0,"content_block":{"type":"text","text":""}}`,
		`{"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"Select"}}`)
	for name, response := range map[string]string{"overloaded_error: Overloaded": failed, "ended before the answer did": cut, "bad JSON": "data: {not json\n\n"} {
		api := newFakeAPI(t, response)
		reply, err := api.model().Start(testRequest()).Next(context.Background(), func(string) {})
		if err == nil || reply != nil || !strings.Contains(err.Error(), name) {
			t.Errorf("%s: reply %+v, err %v", name, reply, err)
		}
	}
}

func TestABusyAPIGetsOneMoreTry(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if calls == 1 {
			http.Error(w, `{"type":"error","error":{"type":"overloaded_error","message":"Overloaded"}}`, 529)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		io.WriteString(w, textReply)
	}))
	defer server.Close()
	reply, err := NewAnthropicModel("test-key", server.URL).Start(testRequest()).Next(context.Background(), func(string) {})
	if err != nil || reply.Text != "Select **Submit**." || calls != 2 {
		t.Fatalf("reply %+v, err %v, calls %d", reply, err, calls)
	}

	// A request that the API refuses is not sent again, and the error has its message.
	calls = 0
	refuse := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		http.Error(w, `{"type":"error","error":{"type":"invalid_request_error","message":"max_tokens: must be positive"}}`, http.StatusBadRequest)
	}))
	defer refuse.Close()
	_, err = NewAnthropicModel("test-key", refuse.URL).Start(testRequest()).Next(context.Background(), func(string) {})
	if err == nil || calls != 1 || !strings.Contains(err.Error(), "400") || !strings.Contains(err.Error(), "max_tokens: must be positive") {
		t.Errorf("err %v, calls %d", err, calls)
	}
}

// With no tools and no effort, the request leaves both out. The benchmark's
// judge sends such a request.
func TestARequestWithNoToolsLeavesThemOut(t *testing.T) {
	api := newFakeAPI(t, textReply)
	request := Request{Model: "claude-sonnet-5-5", MaxTokens: 4000, Shared: "Grade this.", Question: "The answer."}
	if _, err := api.model().Start(request).Next(context.Background(), func(string) {}); err != nil {
		t.Fatal(err)
	}
	body := api.bodies[0]
	if _, set := body["tools"]; set {
		t.Error("the body has tools")
	}
	if _, set := body["output_config"]; set {
		t.Error("the body has an effort")
	}
	if len(at(t, body, "system").([]any)) != 1 || at(t, body, "messages.0.content.0.text") != "The answer." {
		t.Errorf("body = %v", body)
	}
}

func TestTheCommonPartFollowsTheSharedPartWithItsOwnCacheMark(t *testing.T) {
	api := newFakeAPI(t, textReply)
	request := testRequest()
	request.Common = "all the claims, and Turbo's settings"
	if _, err := api.model().Start(request).Next(context.Background(), func(string) {}); err != nil {
		t.Fatal(err)
	}
	body := api.bodies[0]
	for path, want := range map[string]any{
		"system.0.text":               "the rules and the handbook",
		"system.0.cache_control.type": "ephemeral",
		"system.1.text":               "all the claims, and Turbo's settings",
		"system.1.cache_control.type": "ephemeral",
		"system.2.text":               "the user holds the time claim",
	} {
		if got := at(t, body, path); got != want {
			t.Errorf("%s = %v, want %v", path, got, want)
		}
	}
	// A mark on the user's block would write a cache entry for each user.
	if at(t, body, "system.2").(map[string]any)["cache_control"] != nil || len(at(t, body, "system").([]any)) != 3 {
		t.Error("the user's block must be last and have no cache mark")
	}
}
