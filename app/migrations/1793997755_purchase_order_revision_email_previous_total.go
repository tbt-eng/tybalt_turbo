package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

// A revision can change the currency as well as the total, so the email no
// longer says the total changed. It shows the previous total so that an approver
// can't miss a new currency when the amount stays the same. Notifications queued
// before this change have no PreviousTotal, and index tolerates that where a
// plain field would fail on missingkey=error.
var purchaseOrderRevisionEmailPreviousTotal = []projectWorkflowEmail{
	{"porevisionappr0", "po_revision_approval_required", "Asks the approver to approve a purchase order again after it was revised",
		"Approval needed: revised purchase order {{.PONumber}}",
		`Hello {{.RecipientName}},

{{.UserName}} revised purchase order {{.PONumber}}, and you are its approver. It was approved before, and it needs approval again.

Description: {{.Description}}
Vendor: {{.VendorName}}
Revised total: {{.Total}}{{with index . "PreviousTotal"}}
Previous total: {{.}}{{end}}

No expenses can be added to it until it is approved. You may review it, then approve or reject it here:

{{.ActionURL}}`},
}

func init() {
	m.Register(func(app core.App) error {
		return replaceProjectWorkflowEmails(app, nil, purchaseOrderRevisionEmailPreviousTotal)
	}, func(app core.App) error {
		return replaceProjectWorkflowEmails(app, nil, purchaseOrderRevisionEmail)
	})
}
