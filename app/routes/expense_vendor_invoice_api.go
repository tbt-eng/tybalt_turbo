package routes

import (
	"net/http"

	"tybalt/utilities"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

// The vendor invoice behind an expense and every expense that shares it, so
// the people splitting or approving an invoice can see which expenses a change
// affects. See the vendor invoice section of descriptions/expense_documents.md.

type vendorInvoicePart struct {
	ID            string  `db:"id" json:"id"`
	PONumber      string  `db:"po_number" json:"po_number"`
	Date          string  `db:"date" json:"date"`
	Total         float64 `db:"total" json:"total"`
	ForName       string  `db:"for_name" json:"for_name"`
	EnteredByName string  `db:"entered_by_name" json:"entered_by_name"`
	Submitted     bool    `db:"submitted" json:"submitted"`
	Approved      string  `db:"approved" json:"approved"`
	Committed     string  `db:"committed" json:"committed"`
	Rejected      string  `db:"rejected" json:"rejected"`
	// CanOpen is true when the viewer may open this expense. Other parts are
	// listed without a link.
	CanOpen bool `db:"can_open" json:"can_open"`
}

type vendorInvoiceView struct {
	Document       string              `db:"document" json:"document"`
	Number         string              `db:"vendor_invoice_number" json:"vendor_invoice_number"`
	Date           string              `db:"vendor_invoice_date" json:"vendor_invoice_date"`
	Total          float64             `db:"vendor_invoice_total" json:"vendor_invoice_total"`
	CurrencyCode   string              `db:"currency_code" json:"currency_code"`
	RecordedByName string              `db:"recorded_by_name" json:"recorded_by_name"`
	Recorded       string              `db:"vendor_invoice_recorded" json:"vendor_invoice_recorded"`
	Allocated      float64             `db:"-" json:"allocated"`
	Parts          []vendorInvoicePart `db:"-" json:"parts"`
}

const vendorInvoiceViewQuery = `
SELECT ed.id AS document,
  COALESCE(ed.vendor_invoice_number, '') AS vendor_invoice_number,
  COALESCE(ed.vendor_invoice_date, '') AS vendor_invoice_date,
  COALESCE(ed.vendor_invoice_total, 0) AS vendor_invoice_total,
  COALESCE(c.code, '') AS currency_code,
  TRIM(COALESCE(p.given_name, '') || ' ' || COALESCE(p.surname, '')) AS recorded_by_name,
  COALESCE(ed.vendor_invoice_recorded, '') AS vendor_invoice_recorded
FROM expense_documents ed
LEFT JOIN currencies c ON c.id = ed.vendor_invoice_currency
LEFT JOIN profiles p ON p.uid = ed.vendor_invoice_recorded_by
WHERE ed.id = {:document}
`

// The visibility predicate is the one the expense list and details use.
var vendorInvoicePartsQuery = `
SELECT e.id,
  COALESCE(po.po_number, '') AS po_number,
  e.date,
  CAST(e.total AS REAL) AS total,
  TRIM(COALESCE(pu.given_name, '') || ' ' || COALESCE(pu.surname, '')) AS for_name,
  TRIM(COALESCE(pc.given_name, '') || ' ' || COALESCE(pc.surname, '')) AS entered_by_name,
  e.submitted,
  COALESCE(e.approved, '') AS approved,
  COALESCE(e.committed, '') AS committed,
  COALESCE(e.rejected, '') AS rejected,
  CASE WHEN (` + expenseVisibilityPredicate + ` OR ({:has_admin} = 1 AND e.committed != '')) THEN 1 ELSE 0 END AS can_open
FROM expenses e
LEFT JOIN purchase_orders po ON po.id = e.purchase_order
LEFT JOIN profiles pu ON pu.uid = e.uid
LEFT JOIN profiles pc ON pc.uid = e.creator
WHERE e.attachment_document = {:document}
ORDER BY e.date, e.created
`

// createGetExpenseVendorInvoiceHandler returns the vendor invoice behind an
// expense the viewer can see, and every expense on it.
func createGetExpenseVendorInvoiceHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		auth := e.Auth
		if auth == nil {
			return e.Error(http.StatusUnauthorized, "unauthorized", nil)
		}
		params, err := expenseVisibilityParams(app, auth)
		if err != nil {
			return e.Error(http.StatusInternalServerError, "error checking expense visibility claims", err)
		}
		hasAdmin, err := utilities.HasClaim(app, auth, "admin")
		if err != nil {
			return e.Error(http.StatusInternalServerError, "error checking expense visibility claims", err)
		}
		params["id"] = e.Request.PathValue("id")
		params["has_admin"] = boolToInt(hasAdmin)

		var expense struct {
			Document string `db:"attachment_document"`
		}
		err = app.DB().NewQuery(`SELECT COALESCE(e.attachment_document, '') AS attachment_document FROM expenses e
			WHERE e.id = {:id} AND (` + expenseVisibilityPredicate + ` OR ({:has_admin} = 1 AND e.committed != ''))`).
			Bind(params).One(&expense)
		if err != nil {
			return e.Error(http.StatusNotFound, "expense not found or not authorized", err)
		}
		if expense.Document == "" {
			return e.JSON(http.StatusOK, vendorInvoiceView{Parts: []vendorInvoicePart{}})
		}

		params["document"] = expense.Document
		var view vendorInvoiceView
		if err := app.DB().NewQuery(vendorInvoiceViewQuery).Bind(dbx.Params{"document": expense.Document}).One(&view); err != nil {
			return e.Error(http.StatusInternalServerError, "failed to load the invoice", err)
		}
		view.Parts = []vendorInvoicePart{}
		if err := app.DB().NewQuery(vendorInvoicePartsQuery).Bind(params).All(&view.Parts); err != nil {
			return e.Error(http.StatusInternalServerError, "failed to load the invoice's expenses", err)
		}
		cents := int64(0)
		for _, part := range view.Parts {
			if part.Rejected == "" {
				cents += utilities.CurrencyAmountToCents(part.Total)
			}
		}
		view.Allocated = float64(cents) / 100
		return e.JSON(http.StatusOK, view)
	}
}
