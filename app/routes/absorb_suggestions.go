package routes

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"tybalt/absorb"
	"tybalt/absorb/suggest"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

// Absorb suggestions: buckets of probable duplicate clients or vendors, best
// first, for the users who may absorb them.
//
//   - GET /api/absorb_suggestions/{collection}?level=strict|normal|loose lists
//     the buckets, computed on each request (app/absorb/suggest).
//   - POST /api/absorb_suggestions/{collection}/dismiss marks the records of
//     a bucket "Not duplicates", for every user.
//
// Errors are JSON with code and message. See descriptions/Absorb.md.

const absorbDismissalsCollection = "absorb_dismissals"

// suggestionRefs are the references counted for each record: the uses that
// an absorb moves to the kept record. The suggested target is the member
// with the most.
var suggestionRefs = map[string][]absorb.RefConfig{
	"clients": {{Table: "jobs", Column: "client"}, {Table: "jobs", Column: "job_owner"}},
	"vendors": {{Table: "purchase_orders", Column: "vendor"}, {Table: "expenses", Column: "vendor"}},
}

type absorbSuggestionMember struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Alias      string  `json:"alias"`
	References int     `json:"references"`
	Score      float64 `json:"score"`
	// Proposal is the id of an open absorb proposal that names the record,
	// as a record to absorb or as its target. Blank when there is none.
	Proposal string `json:"proposal"`
}

type absorbSuggestionBucket struct {
	// Key is the suggested target's id. A record is in one bucket only.
	Key             string                   `json:"key"`
	TargetID        string                   `json:"target_id"`
	Confidence      float64                  `json:"confidence"`
	MovedReferences int                      `json:"moved_references"`
	Rank            float64                  `json:"rank"`
	Members         []absorbSuggestionMember `json:"members"`
}

type absorbSuggestionPendingAction struct {
	ID         string   `json:"id"`
	TargetID   string   `json:"target_id"`
	TargetName string   `json:"target_name"`
	Absorbed   []string `json:"absorbed"`
	Created    string   `json:"created"`
	// AbsorbedBy is the user who absorbed and AbsorbedByName their name.
	// Both are blank for an action made before the field existed.
	AbsorbedBy     string `json:"absorbed_by"`
	AbsorbedByName string `json:"absorbed_by_name"`
}

type absorbSuggestionsResponse struct {
	Collection string                   `json:"collection"`
	Level      string                   `json:"level"`
	Buckets    []absorbSuggestionBucket `json:"buckets"`
	// Dismissed is the number of buckets hidden by "Not duplicates".
	Dismissed int `json:"dismissed"`
	// PendingAction is the collection's absorb action that waits for a
	// commit or an undo.
	PendingAction *absorbSuggestionPendingAction `json:"pending_action"`
}

func suggestionFailure(e *core.RequestEvent, status int, code string, message string) error {
	return e.JSON(status, map[string]string{"code": code, "message": message})
}

// checkSuggestionAccess refuses a collection without suggestions and a user
// who may not absorb it. When it refuses, ok is false and the refusal is
// already written.
func checkSuggestionAccess(app core.App, e *core.RequestEvent, collection string) (ok bool, err error) {
	if _, ok := suggestionRefs[collection]; !ok {
		return false, suggestionFailure(e, http.StatusBadRequest, "invalid_collection", "absorb suggestions are for clients and vendors only")
	}
	allowed, err := absorb.HasRequiredClaims(app, e.Auth, collection)
	if err != nil {
		return false, e.Error(http.StatusInternalServerError, "failed to check claims", err)
	}
	if !allowed {
		return false, suggestionFailure(e, http.StatusForbidden, "unauthorized", "you may not absorb these records")
	}
	return true, nil
}

func createGetAbsorbSuggestionsHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		collection := e.Request.PathValue("collection")
		if ok, err := checkSuggestionAccess(app, e, collection); !ok {
			return err
		}
		level, err := suggest.ParseLevel(e.Request.URL.Query().Get("level"))
		if err != nil {
			return suggestionFailure(e, http.StatusBadRequest, "invalid_level", "level must be strict, normal or loose")
		}
		response, err := absorbSuggestions(app, collection, level)
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to find absorb suggestions", err)
		}
		return e.JSON(http.StatusOK, response)
	}
}

// absorbSuggestions finds the buckets of a collection and hides the
// dismissed ones.
func absorbSuggestions(app core.App, collection string, level suggest.Level) (*absorbSuggestionsResponse, error) {
	records, err := loadSuggestionRecords(app, collection)
	if err != nil {
		return nil, err
	}
	proposals, err := openProposalRecords(app, collection)
	if err != nil {
		return nil, err
	}
	dismissals, err := loadDismissals(app, collection)
	if err != nil {
		return nil, err
	}
	pending, err := pendingSuggestionAction(app, collection)
	if err != nil {
		return nil, err
	}

	response := &absorbSuggestionsResponse{
		Collection:    collection,
		Level:         string(level),
		Buckets:       []absorbSuggestionBucket{},
		PendingAction: pending,
	}
	for _, bucket := range suggest.Find(records, suggest.LevelParams(level, collection)) {
		if dismissals.Covers(bucket.IDs()) {
			response.Dismissed++
			continue
		}
		out := absorbSuggestionBucket{
			Key:             bucket.TargetID(),
			TargetID:        bucket.TargetID(),
			Confidence:      bucket.Confidence,
			MovedReferences: bucket.MovedRefs,
			Rank:            bucket.Rank,
			Members:         make([]absorbSuggestionMember, 0, len(bucket.Members)),
		}
		for _, member := range bucket.Members {
			out.Members = append(out.Members, absorbSuggestionMember{
				ID:         member.ID,
				Name:       member.Name,
				Alias:      member.Alias,
				References: member.Refs,
				Score:      member.Score,
				Proposal:   proposals[member.ID],
			})
		}
		response.Buckets = append(response.Buckets, out)
	}
	return response, nil
}

// loadSuggestionRecords reads every record of the collection with its
// reference count.
func loadSuggestionRecords(app core.App, collection string) ([]suggest.Record, error) {
	refs := suggestionRefs[collection]
	uses := make([]string, len(refs))
	for i, ref := range refs {
		uses[i] = fmt.Sprintf("SELECT %[2]s AS ref FROM %[1]s WHERE %[2]s != ''", ref.Table, ref.Column)
	}
	query := fmt.Sprintf(`
		SELECT r.id AS id, r.name AS name, COALESCE(r.alias, '') AS alias, COALESCE(u.refs, 0) AS refs
		FROM %s r
		LEFT JOIN (
			SELECT ref, COUNT(*) AS refs FROM (%s) GROUP BY ref
		) u ON u.ref = r.id`, collection, strings.Join(uses, " UNION ALL "))
	rows := []struct {
		ID    string `db:"id"`
		Name  string `db:"name"`
		Alias string `db:"alias"`
		Refs  int    `db:"refs"`
	}{}
	if err := app.DB().NewQuery(query).All(&rows); err != nil {
		return nil, fmt.Errorf("loading %s: %w", collection, err)
	}
	records := make([]suggest.Record, len(rows))
	for i, row := range rows {
		records[i] = suggest.Record{ID: row.ID, Name: row.Name, Alias: row.Alias, Refs: row.Refs}
	}
	return records, nil
}

// openProposalRecords maps each record that an open absorb proposal names, as
// a record to absorb or as its target, to that proposal.
func openProposalRecords(app core.App, collection string) (map[string]string, error) {
	rows := []struct {
		ID       string `db:"id"`
		Proposal string `db:"proposal"`
	}{}
	err := app.DB().NewQuery(`
		SELECT record_id AS id, proposal FROM absorb_proposal_records
		WHERE open = 1 AND collection_name = {:collection}
		UNION ALL
		SELECT target_id AS id, id AS proposal FROM absorb_proposals
		WHERE status = {:pending} AND collection_name = {:collection}
	`).Bind(dbx.Params{"collection": collection, "pending": absorb.ProposalPending}).All(&rows)
	if err != nil {
		return nil, fmt.Errorf("loading open absorb proposals: %w", err)
	}
	out := map[string]string{}
	for _, row := range rows {
		if _, ok := out[row.ID]; !ok {
			out[row.ID] = row.Proposal
		}
	}
	return out, nil
}

func loadDismissals(app core.App, collection string) (*suggest.Dismissals, error) {
	rows := []struct {
		RecordIDs string `db:"record_ids"`
	}{}
	err := app.DB().Select("record_ids").From(absorbDismissalsCollection).
		Where(dbx.HashExp{"collection_name": collection}).All(&rows)
	if err != nil {
		return nil, fmt.Errorf("loading absorb dismissals: %w", err)
	}
	sets := make([][]string, 0, len(rows))
	for _, row := range rows {
		ids := []string{}
		if err := json.Unmarshal([]byte(row.RecordIDs), &ids); err != nil {
			return nil, fmt.Errorf("reading absorb dismissal: %w", err)
		}
		sets = append(sets, ids)
	}
	return suggest.NewDismissals(sets), nil
}

// pendingSuggestionAction describes the collection's pending absorb action,
// or returns nil.
func pendingSuggestionAction(app core.App, collection string) (*absorbSuggestionPendingAction, error) {
	action, err := app.FindFirstRecordByData("absorb_actions", "collection_name", collection)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	out := &absorbSuggestionPendingAction{
		ID:       action.Id,
		TargetID: action.GetString("target_id"),
		Absorbed: []string{},
		Created:  action.GetDateTime("created").String(),
	}
	out.AbsorbedBy = action.GetString("absorbed_by")
	out.AbsorbedByName = userDisplayName(app, out.AbsorbedBy)
	if target, err := app.FindRecordById(collection, out.TargetID); err == nil {
		out.TargetName = absorb.RecordName(target)
	} else if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	absorbed := []map[string]any{}
	if err := json.Unmarshal([]byte(action.GetString("absorbed_records")), &absorbed); err != nil {
		return nil, fmt.Errorf("reading absorbed records: %w", err)
	}
	for _, record := range absorbed {
		name, _ := record["name"].(string)
		if alias, _ := record["alias"].(string); strings.TrimSpace(alias) != "" {
			name += " (" + strings.TrimSpace(alias) + ")"
		}
		out.Absorbed = append(out.Absorbed, name)
	}
	return out, nil
}

func createDismissAbsorbSuggestionHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		collection := e.Request.PathValue("collection")
		if ok, err := checkSuggestionAccess(app, e, collection); !ok {
			return err
		}
		var body struct {
			RecordIDs []string `json:"record_ids"`
		}
		if err := e.BindBody(&body); err != nil {
			return suggestionFailure(e, http.StatusBadRequest, "invalid_body", "the request body is not valid")
		}
		ids := append([]string(nil), body.RecordIDs...)
		sort.Strings(ids)
		for i := 1; i < len(ids); i++ {
			if ids[i] == ids[i-1] {
				return suggestionFailure(e, http.StatusBadRequest, "duplicate_records", "a record is listed twice")
			}
		}
		if len(ids) < 2 {
			return suggestionFailure(e, http.StatusBadRequest, "too_few_records", "choose at least two records")
		}
		var found int
		err := app.DB().Select("COUNT(*)").From(collection).
			Where(dbx.In("id", stringsToAny(ids)...)).Row(&found)
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to check the records", err)
		}
		if found != len(ids) {
			return suggestionFailure(e, http.StatusNotFound, "record_not_found", "a record no longer exists")
		}

		collectionRecord, err := app.FindCollectionByNameOrId(absorbDismissalsCollection)
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to save the dismissal", err)
		}
		dismissal := core.NewRecord(collectionRecord)
		dismissal.Set("collection_name", collection)
		dismissal.Set("record_ids", ids)
		dismissal.Set("dismissed_by", e.Auth.Id)
		if err := app.Save(dismissal); err != nil {
			return e.Error(http.StatusInternalServerError, "failed to save the dismissal", err)
		}
		return e.JSON(http.StatusOK, map[string]string{"id": dismissal.Id})
	}
}
