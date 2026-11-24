package absorb

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"

	"tybalt/utilities"
)

// Absorb proposals: any signed-in user proposes that records be absorbed into
// a target, and an absorb holder carries the proposal out with the normal
// absorb. See descriptions/Absorb.md.
const (
	ProposalsCollection       = "absorb_proposals"
	ProposalRecordsCollection = "absorb_proposal_records"

	ProposalPending   = "pending"
	ProposalAccepted  = "accepted"
	ProposalDeclined  = "declined"
	ProposalWithdrawn = "withdrawn"
	ProposalObsolete  = "obsolete"

	// ProposalTextMaxLength is the most characters of a reason or a comment.
	ProposalTextMaxLength = 1000
)

// ProposalCollections are the collections whose absorb can be proposed. Rate
// sheets are left out: their absorb has rules of its own.
var ProposalCollections = []string{"clients", "client_contacts", "vendors"}

// IsProposalCollection reports whether an absorb of the collection can be
// proposed.
func IsProposalCollection(collectionName string) bool {
	for _, name := range ProposalCollections {
		if name == collectionName {
			return true
		}
	}
	return false
}

// ProposalError is a rule that refuses a proposal request. Message is for the
// user and Code for the UI.
type ProposalError struct {
	Status  int
	Code    string
	Message string
}

func (e *ProposalError) Error() string {
	return e.Message
}

func proposalError(status int, code string, format string, args ...any) *ProposalError {
	return &ProposalError{Status: status, Code: code, Message: fmt.Sprintf(format, args...)}
}

// RecordName is the name of a client, contact or vendor as a proposal keeps
// it, because a committed absorb deletes the record.
func RecordName(record *core.Record) string {
	switch record.Collection().Name {
	case "client_contacts":
		name := strings.TrimSpace(record.GetString("given_name") + " " + record.GetString("surname"))
		if email := strings.TrimSpace(record.GetString("email")); email != "" {
			if name == "" {
				return email
			}
			return name + " (" + email + ")"
		}
		return name
	case "clients":
		return utilities.ClientLabel(record)
	default:
		name := record.GetString("name")
		if alias := strings.TrimSpace(record.GetString("alias")); alias != "" {
			return name + " (" + alias + ")"
		}
		return name
	}
}

// ValidateProposal checks a target and the ids to absorb into it, for a new
// proposal or for an absorb that carries one out. exceptProposal is the
// proposal being carried out, which may name the same records; it is blank
// for a new proposal. It returns the target and the records, in the order of
// ids. A rule that refuses is a *ProposalError.
func ValidateProposal(app core.App, collectionName string, targetID string, ids []string, exceptProposal string) (*core.Record, []*core.Record, error) {
	if !IsProposalCollection(collectionName) {
		return nil, nil, proposalError(http.StatusBadRequest, "invalid_collection", "an absorb of %s can't be proposed", collectionName)
	}
	if len(ids) == 0 {
		return nil, nil, proposalError(http.StatusBadRequest, "no_records", "choose at least one record to absorb")
	}
	seen := map[string]bool{}
	for _, id := range ids {
		if id == targetID {
			return nil, nil, proposalError(http.StatusBadRequest, "target_in_records", "the record to keep can't also be absorbed")
		}
		if seen[id] {
			return nil, nil, proposalError(http.StatusBadRequest, "duplicate_records", "a record is listed twice")
		}
		seen[id] = true
	}

	target, err := app.FindRecordById(collectionName, targetID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil, proposalError(http.StatusNotFound, "record_not_found", "the record to keep no longer exists")
	}
	if err != nil {
		return nil, nil, err
	}
	records := make([]*core.Record, 0, len(ids))
	for _, id := range ids {
		record, err := app.FindRecordById(collectionName, id)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, proposalError(http.StatusNotFound, "record_not_found", "a record to absorb no longer exists")
		}
		if err != nil {
			return nil, nil, err
		}
		records = append(records, record)
	}

	if parent := Configs[collectionName].ParentConstraint; parent != nil {
		for _, record := range records {
			if record.GetString(parent.Field) != target.GetString(parent.Field) {
				return nil, nil, proposalError(http.StatusBadRequest, "different_parent", "every record must belong to the same %s", strings.TrimSuffix(parent.Collection, "s"))
			}
		}
	}

	// A record may be absorbed by only one open proposal, and a record that
	// another open proposal keeps can't be absorbed by this one.
	all := append([]string{targetID}, ids...)
	names := map[string]string{target.Id: RecordName(target)}
	for _, record := range records {
		names[record.Id] = RecordName(record)
	}
	exprs := []dbx.Expression{
		dbx.HashExp{"open": true, "collection_name": collectionName},
		dbx.In("record_id", stringsToAny(all)...),
	}
	if exceptProposal != "" {
		exprs = append(exprs, dbx.NewExp("proposal != {:except}", dbx.Params{"except": exceptProposal}))
	}
	items, err := app.FindAllRecords(ProposalRecordsCollection, exprs...)
	if err != nil {
		return nil, nil, err
	}
	if len(items) > 0 {
		return nil, nil, proposalError(http.StatusConflict, "record_in_open_proposal", "%s is already in another open absorb proposal", names[items[0].GetString("record_id")])
	}
	exprs = []dbx.Expression{
		dbx.HashExp{"status": ProposalPending, "collection_name": collectionName},
		dbx.In("target_id", stringsToAny(ids)...),
	}
	if exceptProposal != "" {
		exprs = append(exprs, dbx.NewExp("id != {:except}", dbx.Params{"except": exceptProposal}))
	}
	targets, err := app.FindAllRecords(ProposalsCollection, exprs...)
	if err != nil {
		return nil, nil, err
	}
	if len(targets) > 0 {
		return nil, nil, proposalError(http.StatusConflict, "record_kept_by_open_proposal", "another open absorb proposal keeps %s", names[targets[0].GetString("target_id")])
	}

	return target, records, nil
}

// ProposalRecordIDs returns the ids of the records that a proposal would
// absorb.
func ProposalRecordIDs(app core.App, proposalID string) ([]string, error) {
	rows, err := app.FindAllRecords(ProposalRecordsCollection, dbx.HashExp{"proposal": proposalID})
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.GetString("record_id"))
	}
	return ids, nil
}

// DecideProposal closes a pending proposal with a status other than pending.
// It frees its records for other proposals. decidedBy may be blank.
func DecideProposal(app core.App, proposal *core.Record, status string, decidedBy string, comment string) error {
	if _, err := app.DB().NewQuery("UPDATE absorb_proposal_records SET open = 0 WHERE proposal = {:proposal}").
		Bind(dbx.Params{"proposal": proposal.Id}).Execute(); err != nil {
		return err
	}
	proposal.Set("status", status)
	proposal.Set("decided_by", decidedBy)
	proposal.Set("decided_at", types.NowDateTime())
	if comment != "" {
		// Cut by characters, not bytes, so that a name is not left with half
		// a character.
		if runes := []rune(comment); len(runes) > ProposalTextMaxLength {
			comment = string(runes[:ProposalTextMaxLength])
		}
		proposal.Set("decision_comment", comment)
	}
	return app.Save(proposal)
}

// FindProposalAction returns the pending absorb action that carries out a
// proposal, or nil.
func FindProposalAction(app core.App, proposalID string) (*core.Record, error) {
	action, err := app.FindFirstRecordByData("absorb_actions", "proposal", proposalID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return action, err
}

// SettleProposalsOnCommit runs in the commit transaction of an absorb action.
// The proposal that the action carries out is accepted, by the user who
// commits. Every other pending proposal that names an absorbed record, as its
// target or as a record to absorb, can't be carried out any more and becomes
// obsolete. A committed client absorb moves the contacts of the absorbed
// clients, so contact proposals take the kept client's name. It returns the
// proposals that it closed, so their requesters can be told.
func SettleProposalsOnCommit(app core.App, action *core.Record, auth *core.Record) ([]*core.Record, error) {
	collectionName := action.GetString("collection_name")
	if !IsProposalCollection(collectionName) {
		return nil, nil
	}
	deciderID := ""
	if auth != nil {
		deciderID = auth.Id
	}
	closed := []*core.Record{}

	linkedID := action.GetString("proposal")
	if linkedID != "" {
		proposal, err := app.FindRecordById(ProposalsCollection, linkedID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		if proposal != nil && proposal.GetString("status") == ProposalPending {
			if err := DecideProposal(app, proposal, ProposalAccepted, deciderID, ""); err != nil {
				return nil, err
			}
			closed = append(closed, proposal)
		}
	}

	var absorbed []map[string]any
	if err := json.Unmarshal([]byte(action.GetString("absorbed_records")), &absorbed); err != nil {
		return nil, fmt.Errorf("reading absorbed records: %w", err)
	}
	ids := []string{}
	for _, record := range absorbed {
		if id, ok := record["id"].(string); ok && id != "" {
			ids = append(ids, id)
		}
	}
	if len(ids) == 0 {
		return closed, nil
	}

	// names[proposal id] are the absorbed records that the proposal names.
	names := map[string][]string{}
	order := []string{}
	add := func(proposalID, name string) {
		if proposalID == linkedID {
			return
		}
		if _, ok := names[proposalID]; !ok {
			order = append(order, proposalID)
		}
		names[proposalID] = append(names[proposalID], name)
	}
	items, err := app.FindAllRecords(ProposalRecordsCollection,
		dbx.HashExp{"open": true, "collection_name": collectionName},
		dbx.In("record_id", stringsToAny(ids)...),
	)
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		add(item.GetString("proposal"), item.GetString("record_name"))
	}
	targets, err := app.FindAllRecords(ProposalsCollection,
		dbx.HashExp{"status": ProposalPending, "collection_name": collectionName},
		dbx.In("target_id", stringsToAny(ids)...),
	)
	if err != nil {
		return nil, err
	}
	for _, target := range targets {
		add(target.Id, target.GetString("target_name"))
	}

	for _, proposalID := range order {
		proposal, err := app.FindRecordById(ProposalsCollection, proposalID)
		if err != nil {
			return nil, err
		}
		if proposal.GetString("status") != ProposalPending {
			continue
		}
		comment := fmt.Sprintf("Another absorb removed %s, so this proposal can't be carried out.", strings.Join(names[proposalID], ", "))
		if err := DecideProposal(app, proposal, ProposalObsolete, deciderID, comment); err != nil {
			return nil, err
		}
		closed = append(closed, proposal)
	}

	if collectionName == "clients" {
		target, err := app.FindRecordById("clients", action.GetString("target_id"))
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		if target != nil {
			contactProposals, err := app.FindAllRecords(ProposalsCollection,
				dbx.HashExp{"status": ProposalPending, "collection_name": "client_contacts"},
				dbx.In("client_id", stringsToAny(ids)...),
			)
			if err != nil {
				return nil, err
			}
			for _, proposal := range contactProposals {
				proposal.Set("client_id", target.Id)
				proposal.Set("client_name", RecordName(target))
				if err := app.Save(proposal); err != nil {
					return nil, err
				}
			}
		}
	}

	return closed, nil
}

func stringsToAny(values []string) []any {
	result := make([]any, len(values))
	for i, value := range values {
		result[i] = value
	}
	return result
}
