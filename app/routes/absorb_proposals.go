package routes

import (
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"tybalt/absorb"
	"tybalt/notifications"
	"tybalt/utilities"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

// Absorb proposals: any signed-in user proposes an absorb of clients, client
// contacts or vendors, and an absorb holder carries it out with the normal
// absorb (POST /api/{collection}/{id}/absorb with the proposal id) or
// declines it. The requester may delete a withdrawn proposal. See
// descriptions/Absorb.md.

//go:embed absorb_proposals.sql
var absorbProposalsQuery string

const absorbProposalDecidedListLimit = 200

type absorbProposalRecordRow struct {
	Proposal   string `db:"proposal" json:"-"`
	RecordID   string `db:"record_id" json:"record_id"`
	RecordName string `db:"record_name" json:"record_name"`
}

type absorbProposalRow struct {
	ID              string                    `db:"id" json:"id"`
	CollectionName  string                    `db:"collection_name" json:"collection_name"`
	TargetID        string                    `db:"target_id" json:"target_id"`
	TargetName      string                    `db:"target_name" json:"target_name"`
	ClientID        string                    `db:"client_id" json:"client_id"`
	ClientName      string                    `db:"client_name" json:"client_name"`
	Requester       string                    `db:"requester" json:"requester"`
	Reason          string                    `db:"reason" json:"reason"`
	Status          string                    `db:"status" json:"status"`
	DecidedAt       string                    `db:"decided_at" json:"decided_at"`
	DecisionComment string                    `db:"decision_comment" json:"decision_comment"`
	Created         string                    `db:"created" json:"created"`
	RequesterName   string                    `db:"requester_name" json:"requester_name"`
	DecidedByName   string                    `db:"decided_by_name" json:"decided_by_name"`
	InProgress      bool                      `db:"in_progress" json:"in_progress"`
	Records         []absorbProposalRecordRow `db:"-" json:"records"`
}

func respondAbsorbProposalError(e *core.RequestEvent, err error) error {
	var refused *absorb.ProposalError
	if errors.As(err, &refused) {
		return e.JSON(refused.Status, map[string]string{"code": refused.Code, "message": refused.Message})
	}
	return e.Error(http.StatusInternalServerError, "failed to update the absorb proposal", err)
}

func absorbProposalFailure(status int, code string, message string) *absorb.ProposalError {
	return &absorb.ProposalError{Status: status, Code: code, Message: message}
}

// holdsAbsorbClaim reports whether the user reviews absorb proposals.
func holdsAbsorbClaim(app core.App, auth *core.Record) (bool, error) {
	return utilities.HasClaim(app, auth, "absorb")
}

// findAbsorbProposalRows returns the proposals that match where, each with
// its records.
func findAbsorbProposalRows(app core.App, where string, orderBy string, params dbx.Params) ([]absorbProposalRow, error) {
	rows := []absorbProposalRow{}
	query := absorbProposalsQuery + "\nWHERE " + where + "\nORDER BY " + orderBy
	if err := app.DB().NewQuery(query).Bind(params).All(&rows); err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return rows, nil
	}
	ids := make([]any, len(rows))
	byID := map[string]int{}
	for i := range rows {
		ids[i] = rows[i].ID
		byID[rows[i].ID] = i
		rows[i].Records = []absorbProposalRecordRow{}
	}
	records := []absorbProposalRecordRow{}
	if err := app.DB().Select("proposal", "record_id", "record_name").
		From(absorb.ProposalRecordsCollection).
		Where(dbx.In("proposal", ids...)).
		OrderBy("record_name", "record_id").
		All(&records); err != nil {
		return nil, err
	}
	for _, record := range records {
		i := byID[record.Proposal]
		rows[i].Records = append(rows[i].Records, record)
	}
	return rows, nil
}

// createGetAbsorbProposalsHandler lists proposals. An absorb holder sees
// everyone's, and anyone else sees their own. ?status=decided lists the
// closed ones, newest first.
func createGetAbsorbProposalsHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		reviewer, err := holdsAbsorbClaim(app, e.Auth)
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to check claims", err)
		}
		params := dbx.Params{"status": absorb.ProposalPending, "uid": e.Auth.Id}
		mine := ""
		if !reviewer {
			mine = " AND p.requester = {:uid}"
		}

		var rows []absorbProposalRow
		switch e.Request.URL.Query().Get("status") {
		case "", absorb.ProposalPending:
			rows, err = findAbsorbProposalRows(app, "p.status = {:status}"+mine, "p.created", params)
		case "decided":
			rows, err = findAbsorbProposalRows(app, "p.status != {:status}"+mine,
				fmt.Sprintf("p.decided_at DESC LIMIT %d", absorbProposalDecidedListLimit), params)
		default:
			return respondAbsorbProposalError(e, absorbProposalFailure(http.StatusBadRequest, "invalid_status", "status must be pending or decided"))
		}
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to load absorb proposals", err)
		}
		return e.JSON(http.StatusOK, rows)
	}
}

// createGetAbsorbProposalHandler returns one proposal to an absorb holder or
// to its requester.
func createGetAbsorbProposalHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		reviewer, err := holdsAbsorbClaim(app, e.Auth)
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to check claims", err)
		}
		rows, err := findAbsorbProposalRows(app, "p.id = {:id}", "p.id", dbx.Params{"id": e.Request.PathValue("id")})
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to load the absorb proposal", err)
		}
		if len(rows) == 0 || (!reviewer && rows[0].Requester != e.Auth.Id) {
			return respondAbsorbProposalError(e, absorbProposalFailure(http.StatusNotFound, "proposal_not_found", "absorb proposal not found"))
		}
		return e.JSON(http.StatusOK, rows[0])
	}
}

// createCreateAbsorbProposalHandler lets any signed-in user propose an
// absorb.
func createCreateAbsorbProposalHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		var body struct {
			CollectionName string   `json:"collection_name"`
			TargetID       string   `json:"target_id"`
			RecordIDs      []string `json:"record_ids"`
			Reason         string   `json:"reason"`
		}
		if err := e.BindBody(&body); err != nil {
			return respondAbsorbProposalError(e, absorbProposalFailure(http.StatusBadRequest, "invalid_request_body", "invalid request body"))
		}
		reason := strings.TrimSpace(body.Reason)

		var proposalID string
		err := app.RunInTransaction(func(txApp core.App) error {
			target, records, err := absorb.ValidateProposal(txApp, body.CollectionName, body.TargetID, body.RecordIDs, "")
			if err != nil {
				return err
			}
			if reason == "" {
				return absorbProposalFailure(http.StatusBadRequest, "reason_required", "say why these records should be absorbed")
			}
			if len(reason) > absorb.ProposalTextMaxLength {
				return absorbProposalFailure(http.StatusBadRequest, "reason_too_long", "the reason is too long")
			}

			proposals, err := txApp.FindCollectionByNameOrId(absorb.ProposalsCollection)
			if err != nil {
				return err
			}
			proposal := core.NewRecord(proposals)
			proposal.Set("collection_name", body.CollectionName)
			proposal.Set("target_id", target.Id)
			proposal.Set("target_name", absorb.RecordName(target))
			if body.CollectionName == "client_contacts" {
				clientID := target.GetString("client")
				proposal.Set("client_id", clientID)
				if client, err := txApp.FindRecordById("clients", clientID); err == nil {
					proposal.Set("client_name", absorb.RecordName(client))
				}
			}
			proposal.Set("requester", e.Auth.Id)
			proposal.Set("reason", reason)
			proposal.Set("status", absorb.ProposalPending)
			if err := txApp.Save(proposal); err != nil {
				return err
			}

			items, err := txApp.FindCollectionByNameOrId(absorb.ProposalRecordsCollection)
			if err != nil {
				return err
			}
			for _, record := range records {
				item := core.NewRecord(items)
				item.Set("proposal", proposal.Id)
				item.Set("collection_name", body.CollectionName)
				item.Set("record_id", record.Id)
				item.Set("record_name", absorb.RecordName(record))
				item.Set("open", true)
				if err := txApp.Save(item); err != nil {
					return err
				}
			}
			proposalID = proposal.Id
			return nil
		})
		if err != nil {
			return respondAbsorbProposalError(e, err)
		}
		return e.JSON(http.StatusOK, map[string]string{"id": proposalID})
	}
}

// findDecidableAbsorbProposal returns a pending proposal that no pending
// absorb action is carrying out.
func findDecidableAbsorbProposal(app core.App, id string) (*core.Record, error) {
	proposal, err := app.FindRecordById(absorb.ProposalsCollection, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, absorbProposalFailure(http.StatusNotFound, "proposal_not_found", "absorb proposal not found")
	}
	if err != nil {
		return nil, err
	}
	if proposal.GetString("status") != absorb.ProposalPending {
		return nil, absorbProposalFailure(http.StatusBadRequest, "proposal_not_pending", "this absorb proposal has already been decided")
	}
	action, err := absorb.FindProposalAction(app, proposal.Id)
	if err != nil {
		return nil, err
	}
	if action != nil {
		return nil, absorbProposalFailure(http.StatusConflict, "proposal_in_progress", "this absorb proposal has been absorbed; commit or undo the absorb first")
	}
	return proposal, nil
}

// createWithdrawAbsorbProposalHandler lets the requester withdraw their
// pending proposal.
func createWithdrawAbsorbProposalHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		err := app.RunInTransaction(func(txApp core.App) error {
			proposal, err := findDecidableAbsorbProposal(txApp, e.Request.PathValue("id"))
			if err != nil {
				return err
			}
			if proposal.GetString("requester") != e.Auth.Id {
				return absorbProposalFailure(http.StatusForbidden, "unauthorized", "only the requester can withdraw this absorb proposal")
			}
			return absorb.DecideProposal(txApp, proposal, absorb.ProposalWithdrawn, e.Auth.Id, "")
		})
		if err != nil {
			return respondAbsorbProposalError(e, err)
		}
		return e.JSON(http.StatusOK, map[string]string{"message": "absorb proposal withdrawn"})
	}
}

// createDeleteAbsorbProposalHandler lets the requester delete their proposal
// once it is withdrawn. A proposal with any other status stays, as the record
// of what was proposed and decided. The delete cascades to the proposal's
// absorb_proposal_records rows in the same transaction.
func createDeleteAbsorbProposalHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		err := app.RunInTransaction(func(txApp core.App) error {
			proposal, err := txApp.FindRecordById(absorb.ProposalsCollection, e.Request.PathValue("id"))
			if errors.Is(err, sql.ErrNoRows) {
				return absorbProposalFailure(http.StatusNotFound, "proposal_not_found", "absorb proposal not found")
			}
			if err != nil {
				return err
			}
			if proposal.GetString("requester") != e.Auth.Id {
				return absorbProposalFailure(http.StatusForbidden, "unauthorized", "only the requester can delete this absorb proposal")
			}
			if proposal.GetString("status") != absorb.ProposalWithdrawn {
				return absorbProposalFailure(http.StatusBadRequest, "proposal_not_withdrawn", "only a withdrawn absorb proposal can be deleted")
			}
			return txApp.Delete(proposal)
		})
		if err != nil {
			return respondAbsorbProposalError(e, err)
		}
		return e.JSON(http.StatusOK, map[string]string{"message": "absorb proposal deleted"})
	}
}

// createDeclineAbsorbProposalHandler lets an absorb holder decline a pending
// proposal with a comment, which the requester is sent.
func createDeclineAbsorbProposalHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		reviewer, err := holdsAbsorbClaim(app, e.Auth)
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to check claims", err)
		}
		if !reviewer {
			return respondAbsorbProposalError(e, absorbProposalFailure(http.StatusForbidden, "unauthorized", "you are not authorized to decline absorb proposals"))
		}

		var body struct {
			Comment string `json:"comment"`
		}
		if err := e.BindBody(&body); err != nil {
			return respondAbsorbProposalError(e, absorbProposalFailure(http.StatusBadRequest, "invalid_request_body", "invalid request body"))
		}
		comment := strings.TrimSpace(body.Comment)

		var proposal *core.Record
		err = app.RunInTransaction(func(txApp core.App) error {
			var err error
			proposal, err = findDecidableAbsorbProposal(txApp, e.Request.PathValue("id"))
			if err != nil {
				return err
			}
			if comment == "" {
				return absorbProposalFailure(http.StatusBadRequest, "comment_required", "say why the proposal is declined")
			}
			if len(comment) > absorb.ProposalTextMaxLength {
				return absorbProposalFailure(http.StatusBadRequest, "comment_too_long", "the comment is too long")
			}
			return absorb.DecideProposal(txApp, proposal, absorb.ProposalDeclined, e.Auth.Id, comment)
		})
		if err != nil {
			return respondAbsorbProposalError(e, err)
		}

		if err := notifications.QueueAbsorbProposalOutcomeNotification(app, proposal); err != nil {
			app.Logger().Error("error queueing absorb proposal notification", "proposal_id", proposal.Id, "error", err)
		}
		return e.JSON(http.StatusOK, map[string]string{"message": "absorb proposal declined"})
	}
}
