package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

// A revision can change the description, division and category, and the
// attachment while it awaits approval. The email shows the previous value of
// each one that changed, and says how the attachment changed, so that an
// approver can't miss it. Notifications queued before this change have none of
// the new keys, and index tolerates that where a plain field would fail on
// missingkey=error.
var purchaseOrderRevisionEmailDetails = []projectWorkflowEmail{
	{"porevisionappr0", "po_revision_approval_required", "Asks the approver to approve a purchase order again after it was revised",
		"Approval needed: revised purchase order {{.PONumber}}",
		`Hello {{.RecipientName}},

{{.UserName}} revised purchase order {{.PONumber}}, and you are its approver. It was approved before, and it needs approval again.

Description: {{.Description}}{{with index . "PreviousDescription"}}
Previous description: {{.}}{{end}}
Vendor: {{.VendorName}}{{with index . "Job"}}
Job: {{.}}{{end}}{{with index . "PreviousJob"}}
Previous job: {{.}}{{end}}{{with index . "PreviousDivision"}}
Division: {{index $ "Division"}}
Previous division: {{.}}{{end}}{{with index . "PreviousCategory"}}
Category: {{or (index $ "Category") "none"}}
Previous category: {{.}}{{end}}
Revised total: {{.Total}}{{with index . "PreviousTotal"}}
Previous total: {{.}}{{end}}{{with index . "Attachment"}}
Attachment: {{.}}{{end}}

No expenses can be added to it until it is approved. You may review it, then approve or reject it here:

{{.ActionURL}}`},
}

func init() {
	m.Register(func(app core.App) error {
		return replaceProjectWorkflowEmails(app, nil, purchaseOrderRevisionEmailDetails)
	}, func(app core.App) error {
		return replaceProjectWorkflowEmails(app, nil, purchaseOrderRevisionEmailJob)
	})
}
