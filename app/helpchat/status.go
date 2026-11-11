package helpchat

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"tybalt/utilities"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

// The status of Help chat, for an admin.
//
// /api/help/status says whether the chat can run and why not: whether the
// setting is valid, which provider keys the server has, who the deploy group
// lets in, whether the corpus is loaded, and how the last day went. It works
// when the chat is off or its setting is bad, because that is when an admin
// needs it. It sends no request to a model provider.
//
// /api/help/status/test_key sends one request to the provider, which uses no
// tokens: it asks for the description of the model (or of the key, for
// OpenRouter). So it shows that the key works without the cost of a question.
// It runs only when an admin presses the button.
//
// Neither route returns a key or a part of one. Neither returns a question or
// an answer: the admin claim is not a right to read chats. The numbers of the
// last day are counts and token totals.

// keyTestTimeout limits the request of a key test.
const keyTestTimeout = 15 * time.Second

// Status is the body of /api/help/status.
type Status struct {
	// SettingsError is why the help_chat row is not valid. Empty means valid,
	// or no row. With an error the chat is off for everyone.
	SettingsError string `json:"settings_error"`
	// HasRow says whether app_config has a help_chat row.
	HasRow bool `json:"has_row"`
	// Settings are the values in use, with the defaults filled in. With a
	// settings error these are the defaults.
	Settings Settings `json:"settings"`
	// Keys has one entry for each provider.
	Keys []KeyStatus `json:"keys"`
	// DeployGroup is the group that the setting names.
	DeployGroup *GroupStatus `json:"deploy_group"`
	// DeployGroups lists every deploy group, so that an admin can copy an id.
	DeployGroups []GroupStatus `json:"deploy_groups"`
	Corpus       CorpusStatus  `json:"corpus"`
	LastDay      UsageStatus   `json:"last_day"`
	// LastAnswered is when the model last answered a question. Empty means
	// never.
	LastAnswered string `json:"last_answered"`
	// LastError is the last turn that ended in an error.
	LastError *ErrorStatus `json:"last_error"`
}

// KeyStatus says whether the server has the key of one provider. It never
// holds the key.
type KeyStatus struct {
	Provider string `json:"provider"`
	// Secret is the name of the server secret that holds the key.
	Secret string `json:"secret"`
	Set    bool   `json:"set"`
	// InUse marks the provider that the setting names.
	InUse bool `json:"in_use"`
}

// GroupStatus describes one deploy group.
type GroupStatus struct {
	ID string `json:"id"`
	// Name is empty when no group has the id.
	Name    string `json:"name"`
	Found   bool   `json:"found"`
	Members int    `json:"members"`
}

// CorpusStatus says whether the search index is in memory. The first question
// after a start loads it.
type CorpusStatus struct {
	// State is "not_loaded", "loaded" or "failed".
	State  string `json:"state"`
	Chunks int    `json:"chunks"`
	Hash   string `json:"hash"`
	Error  string `json:"error"`
}

// UsageStatus holds the totals of the turns of the last 24 hours.
type UsageStatus struct {
	Questions        int `json:"questions"`
	Users            int `json:"users"`
	Errors           int `json:"errors"`
	TokensIn         int `json:"tokens_in"`
	TokensCacheRead  int `json:"tokens_cache_read"`
	TokensCacheWrite int `json:"tokens_cache_write"`
	TokensOut        int `json:"tokens_out"`
}

// ErrorStatus is a turn that ended in an error. Error is the model client's
// message, not the user's question.
type ErrorStatus struct {
	At    string `json:"at"`
	Error string `json:"error"`
}

// KeyTest is the body of /api/help/status/test_key.
type KeyTest struct {
	OK       bool   `json:"ok"`
	Provider string `json:"provider"`
	Model    string `json:"model"`
	Message  string `json:"message"`
}

// requireAdmin returns an error unless the user holds the admin claim.
func requireAdmin(app core.App, e *core.RequestEvent) error {
	isAdmin, err := utilities.HasClaim(app, e.Auth, "admin")
	if err != nil {
		return err
	}
	if !isAdmin {
		return apis.NewForbiddenError("Only an admin can see the status of Help chat.", nil)
	}
	return nil
}

func (s *Service) handleStatus(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		if err := requireAdmin(app, e); err != nil {
			return err
		}
		status, err := s.status(app)
		if err != nil {
			return err
		}
		return e.JSON(http.StatusOK, status)
	}
}

func (s *Service) status(app core.App) (Status, error) {
	var status Status
	settings, err := LoadSettings(app)
	if err != nil {
		status.SettingsError = err.Error()
	}
	status.Settings = settings
	if _, err := app.FindFirstRecordByData("app_config", "key", settingsKey); err == nil {
		status.HasRow = true
	}

	for _, name := range providerNames() {
		status.Keys = append(status.Keys, KeyStatus{
			Provider: name,
			Secret:   providers[name].keyEnv,
			Set:      strings.TrimSpace(os.Getenv(providers[name].keyEnv)) != "",
			InUse:    name == settings.Provider,
		})
	}

	groups, err := loadGroups(app)
	if err != nil {
		return status, err
	}
	status.DeployGroups = groups
	if settings.DeployGroup != "" {
		chosen := GroupStatus{ID: settings.DeployGroup}
		for _, group := range groups {
			if group.ID == settings.DeployGroup {
				chosen = group
			}
		}
		status.DeployGroup = &chosen
	}

	status.Corpus = s.corpusStatus()

	if err := loadUsage(app, &status); err != nil {
		return status, err
	}
	return status, nil
}

// loadGroups lists the deploy groups with the number of members of each.
func loadGroups(app core.App) ([]GroupStatus, error) {
	var rows []struct {
		ID      string `db:"id"`
		Name    string `db:"name"`
		Members int    `db:"members"`
	}
	err := app.DB().NewQuery(`
		SELECT g.id AS id, g.name AS name, COUNT(m.id) AS members
		FROM deploy_groups g
		LEFT JOIN deploy_group_memberships m ON m.deploy_group = g.id
		GROUP BY g.id
		ORDER BY g.name
	`).All(&rows)
	if err != nil {
		return nil, err
	}
	groups := make([]GroupStatus, 0, len(rows))
	for _, row := range rows {
		groups = append(groups, GroupStatus{ID: row.ID, Name: row.Name, Found: true, Members: row.Members})
	}
	return groups, nil
}

func (s *Service) corpusStatus() CorpusStatus {
	switch s.indexState.Load() {
	case indexLoaded:
		return CorpusStatus{State: "loaded", Chunks: s.Index.Size(), Hash: s.Index.Hash()}
	case indexFailed:
		return CorpusStatus{State: "failed", Error: s.indexErr.Error()}
	}
	return CorpusStatus{State: "not_loaded"}
}

// loadUsage fills the totals of the last 24 hours and the last results.
func loadUsage(app core.App, status *Status) error {
	since := time.Now().UTC().Add(-24 * time.Hour).Format("2006-01-02 15:04:05.000Z")
	var totals struct {
		Questions  int `db:"questions"`
		Users      int `db:"users"`
		Errors     int `db:"errors"`
		In         int `db:"tokens_in"`
		CacheRead  int `db:"tokens_cache_read"`
		CacheWrite int `db:"tokens_cache_write"`
		Out        int `db:"tokens_out"`
	}
	err := app.DB().NewQuery(`
		SELECT COUNT(*) AS questions,
			COUNT(DISTINCT uid) AS users,
			COALESCE(SUM(status = 'error'), 0) AS errors,
			COALESCE(SUM(tokens_in), 0) AS tokens_in,
			COALESCE(SUM(tokens_cache_read), 0) AS tokens_cache_read,
			COALESCE(SUM(tokens_cache_write), 0) AS tokens_cache_write,
			COALESCE(SUM(tokens_out), 0) AS tokens_out
		FROM chat_turns
		WHERE created >= {:since}
	`).Bind(dbx.Params{"since": since}).One(&totals)
	if err != nil {
		return err
	}
	status.LastDay = UsageStatus{
		Questions:        totals.Questions,
		Users:            totals.Users,
		Errors:           totals.Errors,
		TokensIn:         totals.In,
		TokensCacheRead:  totals.CacheRead,
		TokensCacheWrite: totals.CacheWrite,
		TokensOut:        totals.Out,
	}

	// A turn with one of these statuses got a reply from the model.
	var answered []struct {
		Created string `db:"created"`
	}
	err = app.DB().NewQuery(`
		SELECT created FROM chat_turns
		WHERE status IN ('answered', 'gap', 'refused')
		ORDER BY created DESC LIMIT 1
	`).All(&answered)
	if err != nil {
		return err
	}
	if len(answered) > 0 {
		status.LastAnswered = answered[0].Created
	}

	var failed []struct {
		Created string `db:"created"`
		Error   string `db:"error"`
	}
	err = app.DB().NewQuery(`
		SELECT created, error FROM chat_turns
		WHERE status = 'error'
		ORDER BY created DESC LIMIT 1
	`).All(&failed)
	if err != nil {
		return err
	}
	if len(failed) > 0 {
		status.LastError = &ErrorStatus{At: failed[0].Created, Error: failed[0].Error}
	}
	return nil
}

func (s *Service) handleTestKey(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		if err := requireAdmin(app, e); err != nil {
			return err
		}
		settings, err := LoadSettings(app)
		if err != nil {
			return e.JSON(http.StatusOK, KeyTest{Message: "The help_chat setting is not valid: " + err.Error()})
		}
		ctx, cancel := context.WithTimeout(e.Request.Context(), keyTestTimeout)
		defer cancel()
		return e.JSON(http.StatusOK, testKey(ctx, settings))
	}
}

// testKey asks the provider of the settings about the model, or about the key
// for OpenRouter. The request has no prompt and uses no tokens.
func testKey(ctx context.Context, settings Settings) KeyTest {
	result := KeyTest{Provider: settings.Provider, Model: settings.Model}
	p := providers[settings.Provider]
	key := strings.TrimSpace(os.Getenv(p.keyEnv))
	if key == "" {
		result.Message = fmt.Sprintf("The server has no %s. Set the secret and restart the app.", p.keyEnv)
		return result
	}

	var address, about string
	header := http.Header{}
	switch settings.Provider {
	case providerAnthropic:
		base := p.baseURL
		if base == "" {
			base = anthropicURL
		}
		address = strings.TrimRight(base, "/") + "/v1/models/" + url.PathEscape(settings.Model)
		header.Set("x-api-key", key)
		header.Set("anthropic-version", anthropicVersion)
		about = "model"
	case providerOpenRouter:
		// OpenRouter's list of models needs no key, so it can't test one. This
		// route describes the key.
		address = strings.TrimRight(p.baseURL, "/") + "/key"
		header.Set("Authorization", "Bearer "+key)
		about = "key"
	default:
		address = strings.TrimRight(p.baseURL, "/") + "/models/" + url.PathEscape(settings.Model)
		header.Set("Authorization", "Bearer "+key)
		about = "model"
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, address, nil)
	if err != nil {
		result.Message = "Turbo can't build the request: " + err.Error()
		return result
	}
	request.Header = header
	response, err := (&http.Client{}).Do(request)
	if err != nil {
		result.Message = "Turbo can't reach the provider: " + redact(err.Error(), key)
		return result
	}
	defer response.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(response.Body, 4000))

	switch {
	case response.StatusCode == http.StatusOK && about == "model":
		result.OK = true
		result.Message = "The key works and the provider knows the model."
	case response.StatusCode == http.StatusOK:
		result.OK = true
		result.Message = "The key works. This provider has no free check of the model."
	case response.StatusCode == http.StatusUnauthorized || response.StatusCode == http.StatusForbidden:
		result.Message = fmt.Sprintf("The provider refused the key in %s (HTTP %d). %s", p.keyEnv, response.StatusCode, redact(providerMessage(body), key))
	case response.StatusCode == http.StatusNotFound && about == "model":
		result.Message = fmt.Sprintf("The key works, but the provider does not know the model %q. %s", settings.Model, redact(providerMessage(body), key))
	default:
		result.Message = fmt.Sprintf("The provider answered HTTP %d. %s", response.StatusCode, redact(providerMessage(body), key))
	}
	result.Message = strings.TrimSpace(result.Message)
	return result
}

// providerMessage reads the message of a provider's error body. The three
// providers use {"error": {"message": "..."}}.
func providerMessage(body []byte) string {
	var parsed struct {
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if json.Unmarshal(body, &parsed) == nil && parsed.Error.Message != "" {
		return truncate(parsed.Error.Message, 300)
	}
	return ""
}

// redact removes the key from a message, in case a provider or an error
// repeats it.
func redact(message, key string) string {
	return strings.ReplaceAll(message, key, "[key]")
}
