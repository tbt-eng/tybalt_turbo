package migrations

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

const (
	expenseUncommitReasonsCollectionID  = "pbc_unc_reasons"
	expenseUncommitRequestsCollectionID = "pbc_unc_requests"
	expenseUncommitPendingIndexName     = "idx_expense_uncommit_requests_pending"

	payablesAdminConvertActive = "convert active One-Time POs to Cumulative"
	payablesAdminConvert       = "convert One-Time POs to Cumulative"
)

type expenseUncommitReason struct {
	id, name, defaultOutcome, notePrompt string
	noteRequired                         bool
}

// expenseUncommitReasons were distilled from the uncommit requests made in
// Teams between May and November 2026, most common first, with Other last.
var expenseUncommitReasons = []expenseUncommitReason{
	{"uncreason000001", "Duplicate PO, invoice is on another PO", "uncommit_and_reject", "Which PO has the invoice?", true},
	{"uncreason000002", "Amount doesn't match the invoice", "uncommit_and_reject", "Correct amount (tax, tip, shipping, over the 5%)", true},
	{"uncreason000003", "Vendor revised the invoice or price", "uncommit_and_reject", "What changed?", false},
	{"uncreason000004", "Wrong currency", "uncommit_and_reject", "Correct currency", false},
	{"uncreason000005", "Wrong payment type", "uncommit_and_reject", "Correct payment type", true},
	{"uncreason000006", "Expense is on the wrong PO", "uncommit_and_reject", "Which PO should it be on?", true},
	{"uncreason000007", "Wrong document attached", "uncommit_and_reject", "For example, order confirmation instead of invoice", false},
	{"uncreason000008", "Order cancelled or billed in error", "uncommit_and_reject", "What happened?", false},
	{"uncreason000009", "Receipt doesn't meet policy", "uncommit_and_reject", "For example, not itemized", false},
	{"uncreason000010", "Committed by mistake", "uncommit_only", "What should have happened?", false},
	{"uncreason000011", "Wrong vendor, job, branch or coding", "uncommit_and_reject", "What should it be?", true},
	{"uncreason000012", "Other", "uncommit_and_reject", "Explain", true},
}

var expenseUncommitEmails = []projectWorkflowEmail{
	{"uncommitapprove", "expense_uncommit_approved", "Tells the requester their expense uncommit request was approved",
		"Uncommit approved: {{.ExpenseLabel}}",
		`Hello {{.RecipientName}},

{{.DeciderName}} approved your uncommit request.

Expense: {{.ExpenseLabel}}
Reason: {{.Reason}}
Result: {{.Outcome}}

{{.ActionURL}}`},
	{"uncommitdecline", "expense_uncommit_declined", "Tells the requester their expense uncommit request was declined",
		"Uncommit declined: {{.ExpenseLabel}}",
		`Hello {{.RecipientName}},

{{.DeciderName}} declined your uncommit request. The expense is still committed.

Expense: {{.ExpenseLabel}}
Reason: {{.Reason}}
Why it was declined: {{.Comment}}

{{.ActionURL}}`},
	{"uncommitreminde", "expense_uncommit_reminder", "Reminds admin holders that expense uncommit requests are waiting",
		"Uncommit requests waiting: {{.Count}}",
		`Hello {{.RecipientName}},

{{.Count}} expense uncommit requests are waiting for a decision:

{{.ActionURL}}`},
}

// Add expense uncommit requests.
//
// Holders of the commit claim ask for a committed expense to be uncommitted,
// and holders of the admin claim approve or decline from a queue. Requests are
// read and written only through the uncommit request routes, so neither
// collection has API rules. See descriptions/uncommit_requests.md.
//
// expense is deliberately not required. PocketBase refuses to delete a record
// that a required relation points at, and a rejected expense on a duplicate PO
// has to be deleted before the PO can be cancelled. The snapshot fields keep
// the history readable once the expense is gone.
//
// A partial unique index allows one pending request per expense.
func init() {
	m.Register(func(app core.App) error {
		users, err := app.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}
		expenses, err := app.FindCollectionByNameOrId("expenses")
		if err != nil {
			return err
		}

		outcomes := []string{"uncommit_and_reject", "uncommit_only"}

		reasons := core.NewBaseCollection("expense_uncommit_reasons")
		reasons.Id = expenseUncommitReasonsCollectionID
		reasons.Fields.Add(
			&core.TextField{Name: "name", Required: true},
			&core.SelectField{Name: "default_outcome", Values: outcomes, MaxSelect: 1, Required: true},
			&core.TextField{Name: "note_prompt"},
			&core.BoolField{Name: "note_required"},
			&core.NumberField{Name: "sort"},
			&core.BoolField{Name: "active"},
			&core.AutodateField{Name: "created", OnCreate: true},
			&core.AutodateField{Name: "updated", OnCreate: true, OnUpdate: true},
		)
		if err := app.Save(reasons); err != nil {
			return err
		}

		for i, reason := range expenseUncommitReasons {
			record := core.NewRecord(reasons)
			record.Id = reason.id
			record.Set("name", reason.name)
			record.Set("default_outcome", reason.defaultOutcome)
			record.Set("note_prompt", reason.notePrompt)
			record.Set("note_required", reason.noteRequired)
			record.Set("sort", i+1)
			record.Set("active", true)
			if err := app.Save(record); err != nil {
				return err
			}
		}

		requests := core.NewBaseCollection("expense_uncommit_requests")
		requests.Id = expenseUncommitRequestsCollectionID
		requests.Fields.Add(
			&core.RelationField{Name: "expense", CollectionId: expenses.Id, MaxSelect: 1},
			&core.RelationField{Name: "requester", CollectionId: users.Id, MaxSelect: 1, Required: true},
			// Blank only on rows written by a direct uncommit.
			&core.RelationField{Name: "reason", CollectionId: reasons.Id, MaxSelect: 1},
			&core.TextField{Name: "note", Max: 1000},
			&core.SelectField{Name: "outcome", Values: outcomes, MaxSelect: 1, Required: true},
			&core.SelectField{Name: "status", Values: []string{"pending", "approved", "declined", "cancelled"}, MaxSelect: 1, Required: true},
			&core.RelationField{Name: "decided_by", CollectionId: users.Id, MaxSelect: 1},
			&core.DateField{Name: "decided_at"},
			&core.TextField{Name: "decision_comment", Max: 1000},
			&core.TextField{Name: "employee_name"},
			&core.TextField{Name: "po_number"},
			&core.TextField{Name: "vendor_name"},
			&core.TextField{Name: "description"},
			&core.TextField{Name: "currency_code"},
			&core.TextField{Name: "payment_type"},
			&core.NumberField{Name: "total"},
			&core.DateField{Name: "committed"},
			&core.RelationField{Name: "committer", CollectionId: users.Id, MaxSelect: 1},
			&core.TextField{Name: "committed_week_ending"},
			&core.TextField{Name: "pay_period_ending"},
			&core.AutodateField{Name: "created", OnCreate: true},
			&core.AutodateField{Name: "updated", OnCreate: true, OnUpdate: true},
		)
		requests.AddIndex(expenseUncommitPendingIndexName, true, "`expense`", "`status` = 'pending'")
		if err := app.Save(requests); err != nil {
			return err
		}

		if err := setPayablesAdminConvertDescription(app, payablesAdminConvertActive, payablesAdminConvert); err != nil {
			return err
		}

		return replaceProjectWorkflowEmails(app, nil, expenseUncommitEmails)
	}, func(app core.App) error {
		if err := replaceProjectWorkflowEmails(app, expenseUncommitEmails, nil); err != nil {
			return err
		}

		if err := setPayablesAdminConvertDescription(app, payablesAdminConvert, payablesAdminConvertActive); err != nil {
			return err
		}

		for _, id := range []string{expenseUncommitRequestsCollectionID, expenseUncommitReasonsCollectionID} {
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

// setPayablesAdminConvertDescription keeps the payables_admin description true
// now that make_cumulative also accepts a system-closed One-Time PO.
func setPayablesAdminConvertDescription(app core.App, from, to string) error {
	claim, err := app.FindFirstRecordByData("claims", "name", "payables_admin")
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	if err != nil {
		return err
	}
	claim.Set("description", strings.ReplaceAll(claim.GetString("description"), from, to))
	return app.Save(claim)
}
