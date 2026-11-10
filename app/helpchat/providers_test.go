package helpchat

import (
	"net/http"
	"strings"
	"testing"
)

// useProviderURL points a provider at a local server for one test.
func useProviderURL(t *testing.T, name, url string) {
	t.Helper()
	old := providers[name]
	changed := old
	changed.baseURL = url
	providers[name] = changed
	t.Cleanup(func() { providers[name] = old })
}

// The setting picks the provider, and the next question uses it. This test
// goes through the real route and both real clients, to two local servers.
func TestTheSettingPicksTheProviderForTheNextQuestion(t *testing.T) {
	openai := newFakeChat(t, chatTextReply)
	claude := newFakeAPI(t, textReply)
	useProviderURL(t, providerOpenAI, openai.url())
	useProviderURL(t, providerAnthropic, claude.server.URL)
	t.Setenv("OPENAI_API_KEY", "openai-key")
	t.Setenv("ANTHROPIC_API_KEY", "anthropic-key")
	t.Setenv("OPENROUTER_API_KEY", "")

	f := newFixture(t, map[string]any{"enabled": true, "provider": "openai", "model": "gpt-test", "effort": "minimal"})
	f.service.Model = nil // use the real clients

	list := events(t, f.ask(t, staffEmail, map[string]any{"question": "How do I submit my week?"}))
	if done := eventData[map[string]string](t, list, "done"); done["status"] != StatusAnswered {
		t.Fatalf("done = %v", done)
	}
	if len(openai.bodies) != 1 || len(claude.bodies) != 0 {
		t.Fatalf("openai got %d requests and claude got %d", len(openai.bodies), len(claude.bodies))
	}
	if openai.bodies[0]["model"] != "gpt-test" || openai.bodies[0]["reasoning_effort"] != "minimal" || openai.headers[0].Get("Authorization") != "Bearer openai-key" {
		t.Errorf("the openai request: model %v, effort %v", openai.bodies[0]["model"], openai.bodies[0]["reasoning_effort"])
	}
	// Both providers get the same rules, handbook and tools.
	if shared := at(t, openai.bodies[0], "messages.0.content"); shared != sharedPrompt(testHandbook) {
		t.Error("the openai request must start with the shared prompt")
	}
	first := f.onlyTurn(t)
	if first.GetString("provider") != "openai" || first.GetString("model") != "gpt-test" || first.GetString("effort") != "minimal" ||
		first.GetString("answer") != "Select **Submit**." || first.GetInt("tokens_in") != 9400 || first.GetInt("tokens_out") != 12 || first.GetString("request_id") != "req_chat_1" {
		t.Errorf("the openai turn: %v", first.PublicExport())
	}

	// An admin changes the row. The next question goes to Claude, with no restart.
	f.setSettings(t, map[string]any{"enabled": true, "provider": "anthropic"})
	events(t, f.ask(t, otherEmail, map[string]any{"question": "And me?"}))
	if len(openai.bodies) != 1 || len(claude.bodies) != 1 {
		t.Fatalf("openai got %d requests and claude got %d", len(openai.bodies), len(claude.bodies))
	}
	if claude.bodies[0]["model"] != "claude-haiku-5-5" || claude.headers[0].Get("Authorization") != "Bearer anthropic-key" ||
		at(t, claude.bodies[0], "system.0.text") != sharedPrompt(testHandbook) {
		t.Errorf("the claude request: model %v", claude.bodies[0]["model"])
	}
	second := f.turns(t)[1]
	if second.GetString("provider") != "anthropic" || second.GetString("model") != "claude-haiku-5-5" || second.GetInt("tokens_cache_read") != 9000 {
		t.Errorf("the claude turn: %v", second.PublicExport())
	}

	// A provider with no key on the server, and settings that name no real
	// provider or no model, leave the chat off.
	for _, settings := range []map[string]any{
		{"enabled": true, "provider": "openrouter", "model": "openai/gpt-test"},
		{"enabled": true, "provider": "azure", "model": "gpt-test"},
		{"enabled": true, "provider": "openai"},
		{"enabled": true, "provider": "openai", "model": "gpt-test", "effort": "max"},
	} {
		f.setSettings(t, settings)
		recorder := f.ask(t, staffEmail, map[string]any{"question": "Now?"})
		mustStatus(t, recorder, http.StatusServiceUnavailable)
		if strings.Contains(recorder.Body.String(), "KEY") || strings.Contains(recorder.Body.String(), "openrouter") {
			t.Errorf("the browser learns about the server's setup: %s", recorder.Body.String())
		}
	}
	if len(f.turns(t)) != 2 {
		t.Fatal("a turn started with no model")
	}
}

func TestAKeyIsReadFromTheProvidersOwnSecret(t *testing.T) {
	for _, name := range providerNames() {
		t.Setenv(providers[name].keyEnv, "")
	}
	for _, name := range providerNames() {
		if _, err := newModel(name); err == nil || !strings.Contains(err.Error(), providers[name].keyEnv) {
			t.Errorf("%s with no key: %v", name, err)
		}
		t.Setenv(providers[name].keyEnv, "  a-key\n")
		if model, err := newModel(name); err != nil || model == nil {
			t.Errorf("%s with a key: %v", name, err)
		}
	}
	if _, err := newModel("nope"); err == nil {
		t.Error("an unknown provider must fail")
	}
}
