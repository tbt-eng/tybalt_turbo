package helpchat

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const statusPath = "/api/help/status"
const testKeyPath = "/api/help/status/test_key"

func (f *fixture) status(t *testing.T) (Status, string) {
	t.Helper()
	recorder := f.do(t, context.Background(), http.MethodGet, statusPath, adminEmail, "")
	mustStatus(t, recorder, http.StatusOK)
	var status Status
	if err := json.Unmarshal(recorder.Body.Bytes(), &status); err != nil {
		t.Fatal(err)
	}
	return status, recorder.Body.String()
}

func (f *fixture) testKey(t *testing.T) (KeyTest, string) {
	t.Helper()
	recorder := f.do(t, context.Background(), http.MethodPost, testKeyPath, adminEmail, "")
	mustStatus(t, recorder, http.StatusOK)
	var result KeyTest
	if err := json.Unmarshal(recorder.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	return result, recorder.Body.String()
}

func TestOnlyAnAdminSeesTheStatus(t *testing.T) {
	f := newFixture(t, on())
	for _, path := range []string{statusPath, testKeyPath} {
		method := http.MethodGet
		if path == testKeyPath {
			method = http.MethodPost
		}
		mustStatus(t, f.do(t, context.Background(), method, path, "", ""), http.StatusUnauthorized)
		for _, email := range []string{staffEmail, otherEmail} {
			mustStatus(t, f.do(t, context.Background(), method, path, email, ""), http.StatusForbidden)
		}
		mustStatus(t, f.do(t, context.Background(), method, path, adminEmail, ""), http.StatusOK)
	}
}

// The status works when the chat is off or its setting is bad, and says why.
func TestTheStatusSaysWhyTheChatIsOff(t *testing.T) {
	f := newFixture(t, nil)
	status, _ := f.status(t)
	if status.HasRow || status.SettingsError != "" || status.Settings.Enabled {
		t.Errorf("with no row: %+v", status)
	}
	if status.Corpus.State != "not_loaded" {
		t.Errorf("the corpus state is %q before the first question", status.Corpus.State)
	}

	f.setSettings(t, map[string]any{"enabled": true, "deploy_grp": "abc"})
	status, _ = f.status(t)
	if !status.HasRow || !strings.Contains(status.SettingsError, "deploy_grp") {
		t.Errorf("a bad property: has_row %v, error %q", status.HasRow, status.SettingsError)
	}

	f.setSettings(t, map[string]any{"enabled": true, "questions_per_hour": 0})
	status, _ = f.status(t)
	if !strings.Contains(status.SettingsError, "questions_per_hour") {
		t.Errorf("a bad limit: error %q", status.SettingsError)
	}
}

// The status names the secret of each provider and says whether the server has
// it. It never holds a key.
func TestTheStatusSaysWhichKeysAreSetAndShowsNone(t *testing.T) {
	const key = "sk-test-secret-value-1234567890"
	f := newFixture(t, on())
	t.Setenv("ANTHROPIC_API_KEY", key)
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("OPENROUTER_API_KEY", "")

	status, body := f.status(t)
	if strings.Contains(body, key) || strings.Contains(body, key[len(key)-4:]) {
		t.Fatal("the status holds a part of the key")
	}
	if len(status.Keys) != 3 {
		t.Fatalf("keys: %+v", status.Keys)
	}
	for _, entry := range status.Keys {
		wantSet := entry.Provider == providerAnthropic
		if entry.Set != wantSet || entry.InUse != wantSet {
			t.Errorf("%s: set %v, in use %v", entry.Provider, entry.Set, entry.InUse)
		}
		if entry.Secret != providers[entry.Provider].keyEnv {
			t.Errorf("%s: secret %q", entry.Provider, entry.Secret)
		}
	}

	t.Setenv("ANTHROPIC_API_KEY", "")
	status, _ = f.status(t)
	for _, entry := range status.Keys {
		if entry.Set {
			t.Errorf("%s is set with no key", entry.Provider)
		}
	}
}

func TestTheStatusDescribesTheDeployGroup(t *testing.T) {
	f := newFixture(t, on())
	status, _ := f.status(t)
	if status.DeployGroup != nil || len(status.DeployGroups) != 0 {
		t.Fatalf("with no group: %+v", status)
	}

	group := f.deployGroup(t, "help chat testers", staffEmail, otherEmail)
	empty := f.deployGroup(t, "another group")
	f.setSettings(t, map[string]any{"enabled": true, "deploy_group": group})
	status, _ = f.status(t)
	if got := status.DeployGroup; got == nil || !got.Found || got.Name != "help chat testers" || got.Members != 2 {
		t.Errorf("the chat's group: %+v", got)
	}
	if len(status.DeployGroups) != 2 || status.DeployGroups[0].ID != empty || status.DeployGroups[0].Members != 0 {
		t.Errorf("all groups: %+v", status.DeployGroups)
	}

	f.setSettings(t, map[string]any{"enabled": true, "deploy_group": "nosuchgroup00000"})
	status, _ = f.status(t)
	if got := status.DeployGroup; got == nil || got.Found || got.Members != 0 {
		t.Errorf("a group that does not exist: %+v", got)
	}
}

// The status counts the turns of the last day and shows the last error. It
// holds no question and no answer.
func TestTheStatusCountsTurnsAndShowsNoChat(t *testing.T) {
	f := newFixture(t, on(), step{text: "The secret answer."})
	events(t, f.ask(t, staffEmail, map[string]any{"question": "The staff secret question?"}))
	events(t, f.ask(t, otherEmail, map[string]any{"question": "The other secret question?"}))

	status, body := f.status(t)
	if strings.Contains(body, "secret question") || strings.Contains(body, "secret answer") {
		t.Fatalf("the status holds a chat: %s", body)
	}
	if status.LastDay.Questions != 2 || status.LastDay.Users != 2 || status.LastDay.Errors != 0 {
		t.Errorf("last day: %+v", status.LastDay)
	}
	if status.LastAnswered == "" || status.LastError != nil {
		t.Errorf("last answered %q, last error %+v", status.LastAnswered, status.LastError)
	}
	if status.Corpus.State != "loaded" || status.Corpus.Chunks != 4 || status.Corpus.Hash != "corpus-hash-for-tests" {
		t.Errorf("corpus: %+v", status.Corpus)
	}

	// A turn that ended in an error shows with its message.
	turn := f.turns(t)[0]
	turn.Set("status", string(StatusError))
	turn.Set("error", "the provider is busy")
	if err := f.app.Save(turn); err != nil {
		t.Fatal(err)
	}
	status, _ = f.status(t)
	if status.LastDay.Errors != 1 || status.LastError == nil || status.LastError.Error != "the provider is busy" {
		t.Errorf("after an error: last day %+v, last error %+v", status.LastDay, status.LastError)
	}
}

// fakeProvider answers the requests of a key test and records them.
type fakeProvider struct {
	server *httptest.Server
	status int
	body   string
	paths  []string
	header http.Header
}

func newFakeProvider(t *testing.T) *fakeProvider {
	t.Helper()
	fake := &fakeProvider{status: http.StatusOK, body: `{"id":"a-model"}`}
	fake.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fake.paths = append(fake.paths, r.Method+" "+r.URL.Path)
		fake.header = r.Header.Clone()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(fake.status)
		_, _ = w.Write([]byte(fake.body))
	}))
	t.Cleanup(fake.server.Close)
	return fake
}

// A key test sends one request with no prompt: it asks for the model, or for
// the key at OpenRouter.
func TestTheKeyTestAsksTheProviderAndUsesNoTokens(t *testing.T) {
	const key = "sk-test-secret-value-1234567890"
	cases := []struct {
		provider string
		model    string
		wantPath string
		auth     string
	}{
		{providerAnthropic, "claude-haiku-5-5", "GET /v1/models/claude-haiku-5-5", "x-api-key"},
		{providerOpenAI, "gpt-5-mini", "GET /models/gpt-5-mini", "Authorization"},
		{providerOpenRouter, "openai/gpt-5-mini", "GET /key", "Authorization"},
	}
	for _, c := range cases {
		t.Run(c.provider, func(t *testing.T) {
			fake := newFakeProvider(t)
			useProviderURL(t, c.provider, fake.server.URL)
			t.Setenv(providers[c.provider].keyEnv, key)
			f := newFixture(t, map[string]any{"enabled": false, "provider": c.provider, "model": c.model})

			result, body := f.testKey(t)
			if !result.OK || result.Provider != c.provider || result.Model != c.model {
				t.Errorf("result: %+v", result)
			}
			if len(fake.paths) != 1 || fake.paths[0] != c.wantPath {
				t.Errorf("requests: %v, want one %q", fake.paths, c.wantPath)
			}
			if !strings.Contains(fake.header.Get(c.auth), key) {
				t.Errorf("the request has no key in %s", c.auth)
			}
			if strings.Contains(body, key) {
				t.Error("the result holds the key")
			}
			if len(f.model.requests) != 0 || f.model.warmed != 0 || len(f.turns(t)) != 0 {
				t.Error("the key test ran the model")
			}
		})
	}
}

func TestTheKeyTestSaysWhatIsWrong(t *testing.T) {
	const key = "sk-test-secret-value-1234567890"
	fake := newFakeProvider(t)
	useProviderURL(t, providerAnthropic, fake.server.URL)
	f := newFixture(t, on())

	// No key: no request goes out.
	t.Setenv("ANTHROPIC_API_KEY", "")
	result, _ := f.testKey(t)
	if result.OK || !strings.Contains(result.Message, "ANTHROPIC_API_KEY") || !strings.Contains(result.Message, "fly secrets set") || len(fake.paths) != 0 {
		t.Errorf("with no key: %+v, requests %v", result, fake.paths)
	}

	// A key that the provider refuses. The provider's message repeats the key
	// here, and the result must not.
	t.Setenv("ANTHROPIC_API_KEY", key)
	fake.status = http.StatusUnauthorized
	fake.body = `{"type":"error","error":{"type":"authentication_error","message":"invalid x-api-key ` + key + `"}}`
	result, body := f.testKey(t)
	if result.OK || !strings.Contains(result.Message, "refused the key") || !strings.Contains(result.Message, "invalid x-api-key") {
		t.Errorf("a refused key: %+v", result)
	}
	if strings.Contains(body, key) {
		t.Error("the result repeats the key")
	}

	// A model that the provider does not know.
	fake.status = http.StatusNotFound
	fake.body = `{"type":"error","error":{"type":"not_found_error","message":"model: claude-haiku-5-5"}}`
	result, _ = f.testKey(t)
	if result.OK || !strings.Contains(result.Message, "does not know the model") {
		t.Errorf("an unknown model: %+v", result)
	}

	// The provider is down.
	fake.status = http.StatusServiceUnavailable
	fake.body = `not json`
	result, _ = f.testKey(t)
	if result.OK || !strings.Contains(result.Message, "HTTP 503") {
		t.Errorf("a provider error: %+v", result)
	}

	// A bad setting: no request goes out.
	sent := len(fake.paths)
	f.setSettings(t, map[string]any{"enabled": true, "deploy_grp": "abc"})
	result, _ = f.testKey(t)
	if result.OK || !strings.Contains(result.Message, "not valid") || len(fake.paths) != sent {
		t.Errorf("a bad setting: %+v", result)
	}
}
