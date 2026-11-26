package routes

import (
	"database/sql"
	_ "embed"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"tybalt/utilities"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

// Absorb history: the committed absorbs, kept for
// ABSORB_HISTORY_RETENTION_DAYS so that a bad absorb can be traced.
//
//   - GET /api/absorb_history lists them, newest first, each with who
//     absorbed, who committed, the kept record, the names of the absorbed
//     records and the number of references moved.
//   - GET /api/absorb_history/{id} returns one with the saved copies of the
//     absorbed records and the references that moved.
//
// An absorb holder sees every row. An admin without the claim sees the rate
// sheet rows only, as with absorb_actions. Errors are JSON with code and
// message. See descriptions/Absorb.md.

//go:embed absorb_history.sql
var absorbHistoryQuery string

type absorbHistoryRow struct {
	ID                string `db:"id" json:"id"`
	CollectionName    string `db:"collection_name" json:"collection_name"`
	TargetID          string `db:"target_id" json:"target_id"`
	TargetName        string `db:"target_name" json:"target_name"`
	AbsorbedRecords   string `db:"absorbed_records" json:"-"`
	UpdatedReferences string `db:"updated_references" json:"-"`
	Proposal          string `db:"proposal" json:"proposal"`
	AbsorbedBy        string `db:"absorbed_by" json:"absorbed_by"`
	AbsorbedByName    string `db:"absorbed_by_name" json:"absorbed_by_name"`
	AbsorbedAt        string `db:"absorbed_at" json:"absorbed_at"`
	CommittedBy       string `db:"committed_by" json:"committed_by"`
	CommittedByName   string `db:"committed_by_name" json:"committed_by_name"`
	CommittedAt       string `db:"committed_at" json:"committed_at"`

	// Absorbed holds the names of the absorbed records.
	Absorbed []string `json:"absorbed"`
	// ReferencesMoved is how many references the absorb moved to the kept
	// record.
	ReferencesMoved int `json:"references_moved"`
}

type absorbHistoryDetail struct {
	absorbHistoryRow
	AbsorbedRecordsJSON   json.RawMessage `json:"absorbed_records"`
	UpdatedReferencesJSON json.RawMessage `json:"updated_references"`
}

// absorbHistoryScope reports which rows the user may read: all of them, the
// rate sheet rows only, or none.
func absorbHistoryScope(app core.App, auth *core.Record) (all bool, rateSheets bool, err error) {
	all, err = utilities.HasClaim(app, auth, "absorb")
	if err != nil || all {
		return all, all, err
	}
	rateSheets, err = utilities.HasClaim(app, auth, "admin")
	return false, rateSheets, err
}

// absorbedRecordName names one saved copy of an absorbed record: its name
// and alias, a contact's given name and surname, or its id.
func absorbedRecordName(record map[string]any) string {
	text := func(key string) string {
		value, _ := record[key].(string)
		return strings.TrimSpace(value)
	}
	name := text("name")
	if name == "" {
		name = strings.TrimSpace(text("given_name") + " " + text("surname"))
	}
	if name == "" {
		return text("id")
	}
	if alias := text("alias"); alias != "" {
		name += " (" + alias + ")"
	}
	return name
}

// summarize fills the names of the absorbed records and the number of
// references moved from the saved JSON. JSON that cannot be read leaves them
// empty: the row is still worth showing.
func (row *absorbHistoryRow) summarize() {
	row.Absorbed = []string{}
	records := []map[string]any{}
	if err := json.Unmarshal([]byte(row.AbsorbedRecords), &records); err == nil {
		for _, record := range records {
			row.Absorbed = append(row.Absorbed, absorbedRecordName(record))
		}
	}
	// updated_references[table][column][row id] = old value
	references := map[string]map[string]map[string]any{}
	if err := json.Unmarshal([]byte(row.UpdatedReferences), &references); err == nil {
		for _, columns := range references {
			for _, rows := range columns {
				row.ReferencesMoved += len(rows)
			}
		}
	}
}

func createGetAbsorbHistoryHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		all, rateSheets, err := absorbHistoryScope(app, e.Auth)
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to check claims", err)
		}
		if !all && !rateSheets {
			return e.JSON(http.StatusForbidden, map[string]string{"code": "unauthorized", "message": "you may not see the absorb history"})
		}
		query := absorbHistoryQuery
		if !all {
			query += " WHERE h.collection_name = 'rate_sheets'"
		}
		query += " ORDER BY h.created DESC, h.id"
		rows := []absorbHistoryRow{}
		if err := app.DB().NewQuery(query).All(&rows); err != nil {
			return e.Error(http.StatusInternalServerError, "failed to read the absorb history", err)
		}
		for i := range rows {
			rows[i].summarize()
		}
		return e.JSON(http.StatusOK, rows)
	}
}

func createGetAbsorbHistoryItemHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		all, rateSheets, err := absorbHistoryScope(app, e.Auth)
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to check claims", err)
		}
		if !all && !rateSheets {
			return e.JSON(http.StatusForbidden, map[string]string{"code": "unauthorized", "message": "you may not see the absorb history"})
		}
		var row absorbHistoryRow
		err = app.DB().NewQuery(absorbHistoryQuery + " WHERE h.id = {:id}").
			Bind(dbx.Params{"id": e.Request.PathValue("id")}).One(&row)
		// A row outside the user's scope is reported as missing.
		if errors.Is(err, sql.ErrNoRows) || (err == nil && !all && row.CollectionName != "rate_sheets") {
			return e.JSON(http.StatusNotFound, map[string]string{"code": "history_not_found", "message": "this absorb is not in the history"})
		}
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to read the absorb history", err)
		}
		row.summarize()
		detail := absorbHistoryDetail{
			absorbHistoryRow:      row,
			AbsorbedRecordsJSON:   rawJSONOrNull(row.AbsorbedRecords),
			UpdatedReferencesJSON: rawJSONOrNull(row.UpdatedReferences),
		}
		return e.JSON(http.StatusOK, detail)
	}
}

// rawJSONOrNull passes stored JSON through, or null when it is not valid.
func rawJSONOrNull(stored string) json.RawMessage {
	if json.Valid([]byte(stored)) {
		return json.RawMessage(stored)
	}
	return json.RawMessage("null")
}

// userDisplayName is the given name and surname of a user's profile, or "".
func userDisplayName(app core.App, uid string) string {
	if uid == "" {
		return ""
	}
	var name string
	err := app.DB().NewQuery("SELECT COALESCE(TRIM(given_name || ' ' || surname), '') FROM profiles WHERE uid = {:uid}").
		Bind(dbx.Params{"uid": uid}).Row(&name)
	if err != nil {
		return ""
	}
	return name
}
