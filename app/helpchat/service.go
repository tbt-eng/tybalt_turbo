package helpchat

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"tybalt/constants"
	"tybalt/helpcorpus"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

// warmEvery is the least time between two cache warm-ups. Claude's cache lasts
// 5 minutes after its last use.
const warmEvery = 4 * time.Minute

// Service holds what the routes need: the search index and the model clients.
type Service struct {
	// Model and Index are set directly in tests. In the app they are nil. The
	// first request loads the index, and each provider's client is made when a
	// setting first names that provider.
	Model Model
	Index *helpcorpus.Index

	indexOnce sync.Once
	indexErr  error

	modelsMu sync.Mutex
	models   map[string]Model

	warmMu   sync.Mutex
	lastWarm time.Time
}

// NewService returns the service for the app. It does no work until a user asks.
func NewService() *Service { return &Service{} }

var errNotSetUp = errors.New("help chat is not set up")

// ready returns the model for the settings, and loads the index on the first
// call. A server with no key for the provider, or with no corpus, has no chat,
// and the rest of Turbo is not affected.
func (s *Service) ready(settings Settings) (Model, error) {
	s.indexOnce.Do(func() {
		if s.Index != nil {
			return
		}
		// In development the server runs in app/, and the repo is one folder up.
		// Open then reads the files on disk. In production it finds no repo and
		// uses the corpus in the binary.
		corpus, err := helpcorpus.Open("..")
		if err != nil {
			s.indexErr = fmt.Errorf("%w: %w", errNotSetUp, err)
			return
		}
		if s.Index, err = helpcorpus.NewIndex(corpus); err != nil {
			s.indexErr = fmt.Errorf("%w: %w", errNotSetUp, err)
		}
	})
	if s.indexErr != nil {
		return nil, s.indexErr
	}
	if s.Model != nil {
		return s.Model, nil
	}

	s.modelsMu.Lock()
	defer s.modelsMu.Unlock()
	if model, ok := s.models[settings.Provider]; ok {
		return model, nil
	}
	// A missing key is not remembered, so a key that arrives with a restart of
	// the app, or in a test, is found.
	model, err := newModel(settings.Provider)
	if err != nil {
		return nil, err
	}
	if s.models == nil {
		s.models = map[string]Model{}
	}
	s.models[settings.Provider] = model
	return model, nil
}

// build names the code that answers: the app version and its git commit.
func build() string {
	return constants.FullVersion + "+" + constants.GitCommitShort
}

// AddRoutes adds the Help chat routes. Every route needs a signed-in user.
func AddRoutes(app core.App, service *Service) {
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		group := se.Router.Group("/api/help")
		group.Bind(apis.RequireAuth("users"))
		group.POST("/ask", service.handleAsk(app))
		group.POST("/turns/{id}/feedback", service.handleFeedback(app))
		group.POST("/warm", service.handleWarm(app))
		return se.Next()
	})
}

// enabledSettings returns the settings, or the error to send when the chat is
// off or its settings are bad.
func enabledSettings(app core.App) (Settings, error) {
	settings, err := LoadSettings(app)
	if err != nil {
		app.Logger().Error("help chat settings are not valid", "error", err)
		return settings, apis.NewApiError(http.StatusServiceUnavailable, "Help chat is not set up.", nil)
	}
	if !settings.Enabled {
		return settings, apis.NewForbiddenError("Help chat is off.", nil)
	}
	return settings, nil
}

func (s *Service) checkReady(app core.App, settings Settings) (Model, error) {
	model, err := s.ready(settings)
	if err != nil {
		app.Logger().Error("help chat can't start", "provider", settings.Provider, "error", err)
		return nil, apis.NewApiError(http.StatusServiceUnavailable, "Help chat is not set up.", nil)
	}
	return model, nil
}

func (s *Service) request(settings Settings) Request {
	return Request{
		Model:     settings.Model,
		Effort:    settings.Effort,
		MaxTokens: settings.MaxOutputTokens,
		Shared:    sharedPrompt(s.Index.Handbook()),
		Tools:     tools,
	}
}

// handleAsk answers one question as a stream of server-sent events:
//
//	meta     the chat id and the turn id
//	status   what the server is doing, such as searching
//	text     the next piece of the answer
//	sources  what the answer cites
//	done     the result of the turn
func (s *Service) handleAsk(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		settings, err := enabledSettings(app)
		if err != nil {
			return err
		}
		var body askBody
		if err := e.BindBody(&body); err != nil {
			return apis.NewBadRequestError("Turbo can't read that request.", nil)
		}
		body.clean()
		if body.Question == "" {
			return apis.NewBadRequestError("Type a question.", nil)
		}
		if len([]rune(body.Question)) > settings.MaxQuestionChars {
			return apis.NewBadRequestError(fmt.Sprintf("The question is too long. The limit is %d characters.", settings.MaxQuestionChars), nil)
		}
		model, err := s.checkReady(app, settings)
		if err != nil {
			return err
		}

		// Who the user is, and what they hold, comes from the auth token and
		// the database. Nothing in the request body can change it.
		uid := e.Auth.Id
		allClaims, held, err := loadClaims(app, uid)
		if err != nil {
			return err
		}
		record, history, err := startTurn(app, uid, turnStart{
			settings: settings,
			body:     body,
			held:     held,
			build:    build(),
			corpus:   s.Index.Hash(),
		})
		var refused *askError
		if errors.As(err, &refused) {
			return apis.NewApiError(refused.status, refused.message, nil)
		}
		if err != nil {
			return err
		}

		request := s.request(settings)
		request.Context = contextPrompt(allClaims, held, body.Route)
		request.History = history
		request.Question = body.Question

		stream := newEventStream(e)
		stream.send("meta", map[string]any{"chat": record.GetString("chat"), "turn": record.Id})

		ctx, cancel := context.WithTimeout(e.Request.Context(), time.Duration(settings.TurnTimeoutSeconds)*time.Second)
		defer cancel()
		result := runTurn(ctx, model, s.Index, turnInput{settings: settings, claims: held, request: request}, stream)

		// The row is saved even when the user has left.
		if err := finishTurn(app, record, result); err != nil {
			app.Logger().Error("help chat can't save a turn", "turn", record.Id, "error", err)
		}
		if result.Status == StatusError {
			app.Logger().Error("help chat turn failed", "turn", record.Id, "error", result.Error)
		}
		stream.send("sources", result.Sources)
		stream.send("done", map[string]any{"status": result.Status, "message": doneMessage(result.Status)})
		return nil
	}
}

// doneMessage is what the user reads when a turn gives no answer.
func doneMessage(status string) string {
	switch status {
	case StatusError:
		return "Help can't answer now. Try again in a minute."
	case StatusRefused:
		return "Help can't answer that question."
	}
	return ""
}

func (s *Service) handleFeedback(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		var body feedbackBody
		if err := e.BindBody(&body); err != nil {
			return apis.NewBadRequestError("Turbo can't read that request.", nil)
		}
		err := saveFeedback(app, e.Request.PathValue("id"), e.Auth.Id, body)
		var refused *askError
		if errors.As(err, &refused) {
			return apis.NewApiError(refused.status, refused.message, nil)
		}
		if err != nil {
			return err
		}
		return e.NoContent(http.StatusNoContent)
	}
}

// handleWarm asks the model to cache the shared prompt, so the user's first
// question starts faster. The browser calls it when the user opens Help. It
// does the work at most once in 4 minutes, for all users together.
func (s *Service) handleWarm(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		settings, err := enabledSettings(app)
		if err != nil {
			return err
		}
		model, err := s.checkReady(app, settings)
		if err != nil {
			return err
		}
		s.warmMu.Lock()
		due := time.Since(s.lastWarm) >= warmEvery
		if due {
			s.lastWarm = time.Now()
		}
		s.warmMu.Unlock()
		if due {
			request := s.request(settings)
			go func() {
				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				defer cancel()
				if err := model.Warm(ctx, request); err != nil {
					app.Logger().Warn("help chat can't warm the model's cache", "error", err)
				}
			}()
		}
		return e.NoContent(http.StatusNoContent)
	}
}

// eventStream writes server-sent events to one response. It is the Events of a
// turn, so the answer reaches the browser as the model writes it.
type eventStream struct {
	e *core.RequestEvent
}

func newEventStream(e *core.RequestEvent) *eventStream {
	header := e.Response.Header()
	header.Set("Content-Type", "text/event-stream")
	header.Set("Cache-Control", "no-store")
	// Tells a proxy in front of the app not to hold the stream back.
	header.Set("X-Accel-Buffering", "no")
	e.Response.WriteHeader(http.StatusOK)
	return &eventStream{e: e}
}

// send writes one event. A write to a browser that has gone fails, and the
// turn's context then stops the model.
func (s *eventStream) send(event string, data any) {
	payload, err := json.Marshal(data)
	if err != nil {
		return
	}
	fmt.Fprintf(s.e.Response, "event: %s\ndata: %s\n\n", event, payload)
	_ = s.e.Flush()
}

func (s *eventStream) Status(state, detail string) {
	s.send("status", map[string]string{"state": state, "detail": detail})
}

func (s *eventStream) Text(delta string) {
	s.send("text", map[string]string{"delta": delta})
}
