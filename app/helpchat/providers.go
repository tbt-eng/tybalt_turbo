package helpchat

import (
	"fmt"
	"os"
	"regexp"
	"slices"
	"sort"
	"strings"

	"github.com/anthropics/anthropic-sdk-go/v2/option"
)

// A provider is a company whose API runs the model. The setting names one, and
// the server reads that provider's key from its own secret.
//
// To add a provider that speaks the OpenAI format, add one entry here with its
// name, the secret that holds its key, and its address. To add a provider with
// another format, write a Model for it and return it from newModel.
type provider struct {
	// keyEnv is the secret that holds the API key. It is a Fly secret in
	// production. A key is never in the repo, in app_config or in the browser.
	keyEnv string
	// baseURL is the API's address. Empty means the SDK's own address.
	baseURL string
	// models lists the model ids that the setting may name. Empty means any id.
	models []string
	// defaultModel is used when the setting names no model.
	defaultModel string
	// efforts lists the values that the setting may give for effort.
	efforts []string
	// tokensField names the output limit in an OpenAI-format request. OpenAI's
	// own API takes max_completion_tokens. Others know only max_tokens.
	tokensField string
}

const (
	providerAnthropic  = "anthropic"
	providerOpenAI     = "openai"
	providerOpenRouter = "openrouter"
)

// openAIEfforts are OpenAI's reasoning_effort values. An empty value sends no
// effort, for a model that has no such setting.
var openAIEfforts = []string{"", "minimal", "low", "medium", "high"}

// providers is a variable so that tests can point a provider at a local server.
var providers = map[string]provider{
	providerAnthropic: {
		keyEnv:       "ANTHROPIC_API_KEY",
		models:       []string{"claude-haiku-5-5", "claude-sonnet-5-5"},
		defaultModel: "claude-haiku-5-5",
		efforts:      []string{"low", "medium", "high"},
	},
	providerOpenAI: {
		keyEnv:      "OPENAI_API_KEY",
		baseURL:     "https://api.openai.com/v1",
		efforts:     openAIEfforts,
		tokensField: "max_completion_tokens",
	},
	// OpenRouter passes a request on to many companies' models. Use the
	// anthropic provider for Claude: this client does not mark the shared
	// prompt for Claude's cache, so each question would pay for all of it.
	providerOpenRouter: {
		keyEnv:      "OPENROUTER_API_KEY",
		baseURL:     "https://openrouter.ai/api/v1",
		efforts:     openAIEfforts,
		tokensField: "max_tokens",
	},
}

func providerNames() []string {
	names := make([]string, 0, len(providers))
	for name := range providers {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// modelID is the shape of a model id, such as gpt-5-mini or openai/gpt-5-mini.
var modelID = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:/-]{0,99}$`)

// checkModel reports if a provider can use a model id and an effort.
func (p provider) check(model, effort string) error {
	switch {
	case len(p.models) > 0 && !slices.Contains(p.models, model):
		return fmt.Errorf("model must be one of %v", p.models)
	case !modelID.MatchString(model):
		return fmt.Errorf("model must be a model id of the provider")
	case !slices.Contains(p.efforts, effort):
		return fmt.Errorf("effort must be one of %q", p.efforts)
	}
	return nil
}

// newModel returns the client for a provider. It fails when the server has no
// key for that provider.
func newModel(name string) (Model, error) {
	p, ok := providers[name]
	if !ok {
		return nil, fmt.Errorf("%w: no provider named %q", errNotSetUp, name)
	}
	key := strings.TrimSpace(os.Getenv(p.keyEnv))
	if key == "" {
		return nil, fmt.Errorf("%w: %s is not set", errNotSetUp, p.keyEnv)
	}
	if name == providerAnthropic {
		if p.baseURL != "" {
			return NewAnthropicModel(key, option.WithBaseURL(p.baseURL)), nil
		}
		return NewAnthropicModel(key), nil
	}
	return NewOpenAIModel(key, p.baseURL, p.tokensField), nil
}
