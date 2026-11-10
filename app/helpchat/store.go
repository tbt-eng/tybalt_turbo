package helpchat

import (
	"database/sql"
	"net/http"
	"regexp"
	"slices"
	"sort"
	"strings"
	"time"

	"tybalt/helpcorpus"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

const turnsCollection = "chat_turns"

// maxAnswerChars matches the size of chat_turns.answer.
const maxAnswerChars = 60000

// askError is a refusal to start a turn, with the HTTP status for it.
type askError struct {
	status  int
	message string
}

func (e *askError) Error() string { return e.message }

var (
	errTooMany = &askError{http.StatusTooManyRequests, "You have asked a lot of questions in the last hour. Try again later."}
	errBusy    = &askError{http.StatusTooManyRequests, "Help is still answering your last question."}
	errNoChat  = &askError{http.StatusNotFound, "Turbo can't find that chat."}
)

// askBody is what the browser sends. It holds no earlier turns and no claims:
// the server reads those from the database.
type askBody struct {
	// Chat continues a chat. Empty starts a new one.
	Chat     string `json:"chat"`
	Question string `json:"question"`
	// Route is the page that the user is on, as a route pattern such as
	// /expenses/[id]/edit. The model gets it.
	Route string `json:"route"`
	// Record is the collection and id of the record on that page, such as
	// expenses/abc123. It goes in the log and not to the model.
	Record string `json:"record"`
	// Surface is where the user asked: "help" or "popover".
	Surface string `json:"surface"`
}

var (
	routePattern  = regexp.MustCompile(`^/[A-Za-z0-9_\-/\[\]().=+@]{0,150}$`)
	recordPattern = regexp.MustCompile(`^[a-z_]{1,50}/[a-z0-9]{1,30}$`)
	chatPattern   = regexp.MustCompile(`^[a-z0-9]{15}$`)
)

// clean drops the optional values that don't have the right shape, so text from
// the browser can't reach the prompt through them.
func (b *askBody) clean() {
	b.Question = strings.TrimSpace(b.Question)
	if !routePattern.MatchString(b.Route) {
		b.Route = ""
	}
	if !recordPattern.MatchString(b.Record) {
		b.Record = ""
	}
	if b.Surface != "popover" {
		b.Surface = "help"
	}
}

// loadClaims returns every claim in Turbo, and the names of the claims that
// this user holds. It reads them on each request, so a change to a user's
// claims applies to their next question.
func loadClaims(app core.App, uid string) ([]Claim, helpcorpus.Claims, error) {
	var all []struct {
		Name        string `db:"name"`
		Description string `db:"description"`
	}
	if err := app.DB().NewQuery(`SELECT name, description FROM claims ORDER BY name`).All(&all); err != nil {
		return nil, nil, err
	}
	var names []struct {
		Name string `db:"name"`
	}
	if err := app.DB().NewQuery(`
		SELECT c.name FROM user_claims uc JOIN claims c ON c.id = uc.cid WHERE uc.uid = {:uid}
	`).Bind(dbx.Params{"uid": uid}).All(&names); err != nil {
		return nil, nil, err
	}
	claims := make([]Claim, len(all))
	for i, claim := range all {
		claims[i] = Claim{Name: claim.Name, Description: claim.Description}
	}
	held := helpcorpus.Claims{}
	for _, name := range names {
		held[name.Name] = true
	}
	return claims, held, nil
}

func claimNames(held helpcorpus.Claims) []string {
	names := make([]string, 0, len(held))
	for name := range held {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

type turnStart struct {
	settings Settings
	body     askBody
	held     helpcorpus.Claims
	build    string
	corpus   string
}

// startTurn checks the user's limits and writes the row for a new turn, in one
// transaction, so two requests at once can't both pass a limit. It returns the
// row and the earlier turns of the chat.
//
// The row exists before the model runs. It gives the turn its id, it counts
// toward the limits at once, and it stays as a trace if the server stops.
func startTurn(app core.App, uid string, start turnStart) (*core.Record, []Exchange, error) {
	var record *core.Record
	var history []Exchange
	err := app.RunInTransaction(func(tx core.App) error {
		now := time.Now().UTC()
		count := func(where string, since time.Time) (int, error) {
			var row struct {
				N int `db:"n"`
			}
			err := tx.DB().NewQuery(`SELECT COUNT(*) AS n FROM chat_turns WHERE uid = {:uid} AND created >= {:since}` + where).
				Bind(dbx.Params{"uid": uid, "since": since.Format(types.DefaultDateLayout)}).One(&row)
			return row.N, err
		}
		asked, err := count("", now.Add(-time.Hour))
		if err != nil {
			return err
		}
		if asked >= start.settings.QuestionsPerHour {
			return errTooMany
		}
		// A running row older than the time limit is from a server that stopped.
		stale := time.Duration(start.settings.TurnTimeoutSeconds+15) * time.Second
		running, err := count(" AND status = 'running'", now.Add(-stale))
		if err != nil {
			return err
		}
		if running > 0 {
			return errBusy
		}

		chat, number := start.body.Chat, 1
		if chat == "" {
			chat = core.GenerateDefaultRandomId()
		} else {
			// The chat must be this user's. An id from another user finds no rows.
			if !chatPattern.MatchString(chat) {
				return errNoChat
			}
			var last struct {
				Turn sql.NullInt64 `db:"turn"`
			}
			if err := tx.DB().NewQuery(`SELECT MAX(turn) AS turn FROM chat_turns WHERE chat = {:chat} AND uid = {:uid}`).
				Bind(dbx.Params{"chat": chat, "uid": uid}).One(&last); err != nil {
				return err
			}
			if !last.Turn.Valid {
				return errNoChat
			}
			number = int(last.Turn.Int64) + 1
			if history, err = loadHistory(tx, chat, uid, start.settings.HistoryTurns); err != nil {
				return err
			}
		}

		collection, err := tx.FindCollectionByNameOrId(turnsCollection)
		if err != nil {
			return err
		}
		record = core.NewRecord(collection)
		record.Set("uid", uid)
		record.Set("chat", chat)
		record.Set("turn", number)
		record.Set("question", start.body.Question)
		record.Set("model", start.settings.Model)
		record.Set("effort", start.settings.Effort)
		record.Set("build", start.build)
		record.Set("corpus", start.corpus)
		record.Set("claims", claimNames(start.held))
		record.Set("route", start.body.Route)
		record.Set("record", start.body.Record)
		record.Set("surface", start.body.Surface)
		record.Set("status", StatusRunning)
		return tx.Save(record)
	})
	return record, history, err
}

// loadHistory returns the last answered turns of one user's chat, oldest first.
func loadHistory(app core.App, chat, uid string, limit int) ([]Exchange, error) {
	if limit <= 0 {
		return nil, nil
	}
	var rows []struct {
		Question string `db:"question"`
		Answer   string `db:"answer"`
	}
	err := app.DB().NewQuery(`
		SELECT question, answer FROM chat_turns
		WHERE chat = {:chat} AND uid = {:uid} AND status IN ('answered', 'gap') AND answer != ''
		ORDER BY turn DESC LIMIT {:limit}
	`).Bind(dbx.Params{"chat": chat, "uid": uid, "limit": limit}).All(&rows)
	if err != nil {
		return nil, err
	}
	history := make([]Exchange, len(rows))
	for i, row := range rows {
		history[len(rows)-1-i] = Exchange{Question: row.Question, Answer: row.Answer}
	}
	return history, nil
}

// finishTurn writes the result of a turn to its row.
func finishTurn(app core.App, record *core.Record, result turnResult) error {
	record.Set("status", result.Status)
	record.Set("answer", truncate(result.Answer, maxAnswerChars))
	record.Set("sources", result.Sources)
	record.Set("tool_calls", result.ToolCalls)
	record.Set("gap", result.Gap)
	record.Set("stop_reason", result.StopReason)
	record.Set("error", result.Error)
	record.Set("request_id", result.RequestID)
	record.Set("tokens_in", result.Usage.Input)
	record.Set("tokens_cache_read", result.Usage.CacheRead)
	record.Set("tokens_cache_write", result.Usage.CacheWrite)
	record.Set("tokens_out", result.Usage.Output)
	record.Set("ms_first_token", result.FirstTokenMS)
	record.Set("ms_total", result.TotalMS)
	return app.Save(record)
}

// feedbackBody is what the browser sends about a turn. Each part is optional.
type feedbackBody struct {
	// Rating is "up" or "down". An empty value clears the rating.
	Rating *string `json:"rating"`
	Note   *string `json:"note"`
	// Opened is a source that the user opened: a chunk's ref or a doc's url.
	Opened string `json:"opened"`
	// Copied is true when the user copied the answer.
	Copied *bool `json:"copied"`
}

const (
	maxNoteChars     = 1000
	maxOpenedSources = 50
)

// saveFeedback changes only the feedback fields of a turn, and only for the
// user who asked. A turn of another user is "not found", the same as no turn.
func saveFeedback(app core.App, id, uid string, body feedbackBody) error {
	record, err := app.FindRecordById(turnsCollection, id)
	if err != nil || record.GetString("uid") != uid {
		return &askError{http.StatusNotFound, "Turbo can't find that answer."}
	}
	if body.Rating != nil {
		if !slices.Contains([]string{"", "up", "down"}, *body.Rating) {
			return &askError{http.StatusBadRequest, `The rating must be "up" or "down".`}
		}
		record.Set("rating", *body.Rating)
	}
	if body.Note != nil {
		record.Set("rating_note", truncate(strings.TrimSpace(*body.Note), maxNoteChars))
	}
	if opened := truncate(strings.TrimSpace(body.Opened), 200); opened != "" {
		var sources []string
		_ = record.UnmarshalJSONField("opened_sources", &sources)
		if !slices.Contains(sources, opened) && len(sources) < maxOpenedSources {
			record.Set("opened_sources", append(sources, opened))
		}
	}
	if body.Copied != nil && *body.Copied {
		record.Set("copied", true)
	}
	return app.Save(record)
}
