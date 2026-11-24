package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

const (
	absorbProposalsCollectionID       = "pbc_absorb_prop"
	absorbProposalRecordsCollectionID = "pbc_absorb_precs"
	absorbProposalOpenIndexName       = "idx_absorb_proposal_records_open"
)

var absorbProposalEmails = []projectWorkflowEmail{
	{"absorbpropaccep", "absorb_proposal_accepted", "Tells the requester their absorb proposal was carried out",
		"Absorb done: {{.TargetLabel}}",
		`Hello {{.RecipientName}},

{{.DeciderName}} absorbed the records you proposed.

Kept: {{.TargetLabel}}
Absorbed: {{.Records}}

{{.ActionURL}}`},
	{"absorbpropdecli", "absorb_proposal_declined", "Tells the requester their absorb proposal was declined",
		"Absorb proposal declined: {{.TargetLabel}}",
		`Hello {{.RecipientName}},

{{.DeciderName}} declined your absorb proposal. Nothing was absorbed.

Kept: {{.TargetLabel}}
Records: {{.Records}}
Why it was declined: {{.Comment}}

{{.ActionURL}}`},
	{"absorbpropobsol", "absorb_proposal_obsolete", "Tells the requester their absorb proposal was overtaken by another absorb",
		"Absorb proposal closed: {{.TargetLabel}}",
		`Hello {{.RecipientName}},

Another absorb removed a record in your absorb proposal, so the proposal is closed.

Kept: {{.TargetLabel}}
Records: {{.Records}}
{{.Comment}}

{{.ActionURL}}`},
}

// Add absorb proposals.
//
// Any signed-in user proposes that records of clients, client_contacts or
// vendors be absorbed into a target, and an absorb holder reviews the
// proposal and carries it out with the normal absorb. Proposals are read and
// written only through the absorb proposal routes, so neither collection has
// API rules. See descriptions/Absorb.md.
//
// The names of the target and the records are copied onto the rows, because
// committing an absorb deletes the absorbed records.
//
// A record may be an item of only one open proposal. open mirrors a pending
// proposal on each of its records, so a partial unique index can enforce it.
//
// absorb_actions.proposal links a pending absorb action to the proposal it
// carries out. Its commit accepts the proposal.
func init() {
	m.Register(func(app core.App) error {
		users, err := app.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}

		collections := []string{"clients", "client_contacts", "vendors"}

		proposals := core.NewBaseCollection("absorb_proposals")
		proposals.Id = absorbProposalsCollectionID
		proposals.Fields.Add(
			&core.SelectField{Name: "collection_name", Values: collections, MaxSelect: 1, Required: true},
			&core.TextField{Name: "target_id", Required: true},
			&core.TextField{Name: "target_name"},
			// For client_contacts, the client that the contacts belong to.
			&core.TextField{Name: "client_id"},
			&core.TextField{Name: "client_name"},
			&core.RelationField{Name: "requester", CollectionId: users.Id, MaxSelect: 1, Required: true},
			&core.TextField{Name: "reason", Required: true, Max: 1000},
			&core.SelectField{Name: "status", Values: []string{"pending", "accepted", "declined", "withdrawn", "obsolete"}, MaxSelect: 1, Required: true},
			&core.RelationField{Name: "decided_by", CollectionId: users.Id, MaxSelect: 1},
			&core.DateField{Name: "decided_at"},
			&core.TextField{Name: "decision_comment", Max: 1000},
			&core.AutodateField{Name: "created", OnCreate: true},
			&core.AutodateField{Name: "updated", OnCreate: true, OnUpdate: true},
		)
		proposals.AddIndex("idx_absorb_proposals_status", false, "`status`", "")
		if err := app.Save(proposals); err != nil {
			return err
		}

		records := core.NewBaseCollection("absorb_proposal_records")
		records.Id = absorbProposalRecordsCollectionID
		records.Fields.Add(
			&core.RelationField{Name: "proposal", CollectionId: proposals.Id, MaxSelect: 1, Required: true, CascadeDelete: true},
			&core.SelectField{Name: "collection_name", Values: collections, MaxSelect: 1, Required: true},
			&core.TextField{Name: "record_id", Required: true},
			&core.TextField{Name: "record_name"},
			// True while the proposal is pending.
			&core.BoolField{Name: "open"},
			&core.AutodateField{Name: "created", OnCreate: true},
			&core.AutodateField{Name: "updated", OnCreate: true, OnUpdate: true},
		)
		records.AddIndex(absorbProposalOpenIndexName, true, "`collection_name`, `record_id`", "`open` = 1")
		records.AddIndex("idx_absorb_proposal_records_proposal", false, "`proposal`", "")
		if err := app.Save(records); err != nil {
			return err
		}

		actions, err := app.FindCollectionByNameOrId("absorb_actions")
		if err != nil {
			return err
		}
		actions.Fields.Add(&core.TextField{Name: "proposal"})
		if err := app.Save(actions); err != nil {
			return err
		}

		return replaceProjectWorkflowEmails(app, nil, absorbProposalEmails)
	}, func(app core.App) error {
		if err := replaceProjectWorkflowEmails(app, absorbProposalEmails, nil); err != nil {
			return err
		}

		actions, err := app.FindCollectionByNameOrId("absorb_actions")
		if err != nil {
			return err
		}
		actions.Fields.RemoveByName("proposal")
		if err := app.Save(actions); err != nil {
			return err
		}

		for _, id := range []string{absorbProposalRecordsCollectionID, absorbProposalsCollectionID} {
			collection, err := app.FindCollectionByNameOrId(id)
			if err != nil {
				continue
			}
			if err := app.Delete(collection); err != nil {
				return err
			}
		}
		return nil
	})
}
