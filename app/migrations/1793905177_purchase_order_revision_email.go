package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

// A revised purchase order goes back to its approver with the number it already
// has. po_approval_required says the purchase order was just created, which
// would mislead the approver, so a revision gets its own email.
var purchaseOrderRevisionEmail = []projectWorkflowEmail{
	{"porevisionappr0", "po_revision_approval_required", "Asks the approver to approve a purchase order again after its total was revised",
		"Approval needed: revised purchase order {{.PONumber}}",
		`Hello {{.RecipientName}},

{{.UserName}} revised the total of purchase order {{.PONumber}}, and you are its approver. It was approved before, and it needs approval again for the new total.

Description: {{.Description}}
Vendor: {{.VendorName}}
Revised total: {{.Total}}

No expenses can be added to it until it is approved. You may review it, then approve or reject it here:

{{.ActionURL}}`},
}

func init() {
	m.Register(func(app core.App) error {
		return replaceProjectWorkflowEmails(app, nil, purchaseOrderRevisionEmail)
	}, func(app core.App) error {
		return replaceProjectWorkflowEmails(app, purchaseOrderRevisionEmail, nil)
	})
}
