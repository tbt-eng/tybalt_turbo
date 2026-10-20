package routes

import (
	"errors"
	"net/http"
	"time"

	"tybalt/errs"
	"tybalt/hooks"
	"tybalt/utilities"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

// Client invoices and the costs billed back on them. Drafts are created,
// edited and deleted through the invoices collection, and lines through the
// invoice_lines collection; these routes issue invoices and give the screens
// the joined rows the collections cannot. See descriptions/invoices.md.

// invoiceCost is an expense or PO billed, or billable, on an invoice.
type invoiceCost struct {
	Kind         string  `db:"kind" json:"kind"`
	CostID       string  `db:"cost_id" json:"cost_id"`
	JobID        string  `db:"job_id" json:"job_id"`
	JobNumber    string  `db:"job_number" json:"job_number"`
	PONumber     string  `db:"po_number" json:"po_number"`
	Description  string  `db:"description" json:"description"`
	VendorName   string  `db:"vendor_name" json:"vendor_name"`
	Date         string  `db:"date" json:"date"`
	CostTotal    float64 `db:"cost_total" json:"cost_total"`
	CurrencyCode string  `db:"currency_code" json:"currency_code"`
}

type invoiceLineRow struct {
	ID string `db:"id" json:"id"`
	invoiceCost
	Amount      float64 `db:"amount" json:"amount"`
	Created     string  `db:"created" json:"created"`
	CreatorName string  `db:"creator_name" json:"creator_name"`
	// CanOpen is true when the viewer may open the expense or PO.
	CanOpen bool `db:"-" json:"can_open"`
}

type billableJob struct {
	ID            string `db:"id" json:"id"`
	Number        string `db:"number" json:"number"`
	Description   string `db:"description" json:"description"`
	BillableCount int    `db:"billable_count" json:"billable_count"`
}

// The columns every cost row shares, for an expense aliased e and a PO
// aliased po. An expense line shows its PO number when it has one.
const invoiceExpenseCostColumns = `
  'expense' AS kind, e.id AS cost_id, e.job AS job_id, COALESCE(j.number, '') AS job_number,
  COALESCE(ep.po_number, '') AS po_number, COALESCE(e.description, '') AS description,
  COALESCE(v.name, '') AS vendor_name, e.date AS date, CAST(e.total AS REAL) AS cost_total,
  COALESCE(c.code, 'CAD') AS currency_code`

const invoiceExpenseCostJoins = `
LEFT JOIN jobs j ON j.id = e.job
LEFT JOIN purchase_orders ep ON ep.id = e.purchase_order
LEFT JOIN vendors v ON v.id = e.vendor
LEFT JOIN currencies c ON c.id = e.currency`

const invoicePOCostColumns = `
  'po' AS kind, po.id AS cost_id, po.job AS job_id, COALESCE(j.number, '') AS job_number,
  po.po_number AS po_number, COALESCE(po.description, '') AS description,
  COALESCE(v.name, '') AS vendor_name, po.date AS date, CAST(po.total AS REAL) AS cost_total,
  COALESCE(c.code, 'CAD') AS currency_code`

const invoicePOCostJoins = `
LEFT JOIN jobs j ON j.id = po.job
LEFT JOIN vendors v ON v.id = po.vendor
LEFT JOIN currencies c ON c.id = po.currency`

// Unbilled expenses and POs that may be billed: committed expenses, and
// Active or Closed POs. Callers add the job or client filter.
const billableExpensesWhere = `e.committed != ''
  AND NOT EXISTS (SELECT 1 FROM invoice_lines l WHERE l.expense = e.id)`

const billablePOsWhere = `po.status IN ('Active', 'Closed')
  AND NOT EXISTS (SELECT 1 FROM invoice_lines l WHERE l.purchase_order = po.id)`

var invoiceLinesQuery = `
SELECT l.id, CAST(l.amount AS REAL) AS amount, l.created AS created,
  TRIM(COALESCE(p.given_name, '') || ' ' || COALESCE(p.surname, '')) AS creator_name,` + invoiceExpenseCostColumns + `
FROM invoice_lines l
JOIN expenses e ON e.id = l.expense` + invoiceExpenseCostJoins + `
LEFT JOIN profiles p ON p.uid = l.creator
WHERE l.invoice = {:invoice} AND l.kind = 'expense'
UNION ALL
SELECT l.id, CAST(l.amount AS REAL) AS amount, l.created AS created,
  TRIM(COALESCE(p.given_name, '') || ' ' || COALESCE(p.surname, '')) AS creator_name,` + invoicePOCostColumns + `
FROM invoice_lines l
JOIN purchase_orders po ON po.id = l.purchase_order` + invoicePOCostJoins + `
LEFT JOIN profiles p ON p.uid = l.creator
WHERE l.invoice = {:invoice} AND l.kind = 'po'
ORDER BY job_number, date, created
`

var billableCostsQuery = `
SELECT` + invoiceExpenseCostColumns + `
FROM expenses e` + invoiceExpenseCostJoins + `
WHERE e.job = {:job} AND j.client = {:client} AND ` + billableExpensesWhere + `
UNION ALL
SELECT` + invoicePOCostColumns + `
FROM purchase_orders po` + invoicePOCostJoins + `
WHERE po.job = {:job} AND j.client = {:client} AND ` + billablePOsWhere + `
ORDER BY date DESC
`

var billableJobsQuery = `
SELECT id, number, description, billable_count FROM (
  SELECT j.id, j.number, COALESCE(j.description, '') AS description,
    (SELECT COUNT(*) FROM expenses e WHERE e.job = j.id AND ` + billableExpensesWhere + `)
    + (SELECT COUNT(*) FROM purchase_orders po WHERE po.job = j.id AND ` + billablePOsWhere + `) AS billable_count
  FROM jobs j
  WHERE j.client = {:client}
)
WHERE billable_count > 0
ORDER BY number DESC
`

func writeInvoiceError(e *core.RequestEvent, err error) error {
	var hookErr *errs.HookError
	if errors.As(err, &hookErr) {
		status := hookErr.Status
		if status == 0 {
			status = http.StatusBadRequest
		}
		return e.JSON(status, hookErr)
	}
	return e.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
}

func invoiceRouteError(status int, message string) *errs.HookError {
	return &errs.HookError{Status: status, Message: message}
}

// requireInvoiceViewer allows holders of invoicing or report.
func requireInvoiceViewer(app core.App, auth *core.Record) error {
	for _, claim := range []string{"invoicing", "report"} {
		has, err := utilities.HasClaim(app, auth, claim)
		if err != nil {
			return err
		}
		if has {
			return nil
		}
	}
	return invoiceRouteError(http.StatusForbidden, "the invoicing or report claim is required")
}

// createIssueInvoiceHandler issues a draft: it checks the number, date and
// total, and records who issued it and when.
func createIssueInvoiceHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		if err := hooks.RequireActiveInvoicingUser(app, e.Auth); err != nil {
			return writeInvoiceError(e, err)
		}
		var issued *core.Record
		err := app.RunInTransaction(func(tx core.App) error {
			invoice, err := tx.FindRecordById("invoices", e.Request.PathValue("id"))
			if err != nil {
				return invoiceRouteError(http.StatusNotFound, "invoice not found")
			}
			if err := hooks.ValidateInvoiceForIssue(tx, invoice); err != nil {
				return err
			}
			invoice.Set("status", "Issued")
			invoice.Set("issued_by", e.Auth.Id)
			invoice.Set("issued_at", time.Now().UTC().Format("2006-01-02 15:04:05.000Z"))
			if err := tx.Save(invoice); err != nil {
				return err
			}
			issued = invoice
			return nil
		})
		if err != nil {
			return writeInvoiceError(e, err)
		}
		return e.JSON(http.StatusOK, issued)
	}
}

// createGetInvoiceLinesHandler returns an invoice's lines with each cost's
// details.
func createGetInvoiceLinesHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		if err := requireInvoiceViewer(app, e.Auth); err != nil {
			return writeInvoiceError(e, err)
		}
		invoice, err := app.FindRecordById("invoices", e.Request.PathValue("id"))
		if err != nil {
			return writeInvoiceError(e, invoiceRouteError(http.StatusNotFound, "invoice not found"))
		}
		lines := []invoiceLineRow{}
		if err := app.DB().NewQuery(invoiceLinesQuery).Bind(dbx.Params{"invoice": invoice.Id}).All(&lines); err != nil {
			return writeInvoiceError(e, err)
		}
		if err := markOpenableCosts(app, e.Auth, lines); err != nil {
			return writeInvoiceError(e, err)
		}
		return e.JSON(http.StatusOK, lines)
	}
}

// markOpenableCosts sets CanOpen on the lines whose expense or PO the viewer
// may open, using the same visibility as the details pages.
func markOpenableCosts(app core.App, auth *core.Record, lines []invoiceLineRow) error {
	params, err := expenseVisibilityParams(app, auth)
	if err != nil {
		return err
	}
	hasAdmin, err := utilities.HasClaim(app, auth, "admin")
	if err != nil {
		return err
	}
	params["has_admin"] = boolToInt(hasAdmin)
	for i := range lines {
		if lines[i].Kind == "po" {
			row, err := findVisiblePurchaseOrderByID(app, auth.Id, lines[i].CostID)
			if err != nil {
				return err
			}
			lines[i].CanOpen = row != nil
			continue
		}
		params["id"] = lines[i].CostID
		var found struct {
			ID string `db:"id"`
		}
		err := app.DB().NewQuery(`SELECT e.id FROM expenses e WHERE e.id = {:id} AND (` +
			expenseVisibilityPredicate + ` OR ({:has_admin} = 1 AND e.committed != ''))`).Bind(params).One(&found)
		lines[i].CanOpen = err == nil
	}
	return nil
}

// createGetInvoiceBillableHandler lists the invoice client's jobs that have
// unbilled costs or, with ?job=, that job's unbilled costs.
func createGetInvoiceBillableHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		if err := hooks.RequireActiveInvoicingUser(app, e.Auth); err != nil {
			return writeInvoiceError(e, err)
		}
		invoice, err := app.FindRecordById("invoices", e.Request.PathValue("id"))
		if err != nil {
			return writeInvoiceError(e, invoiceRouteError(http.StatusNotFound, "invoice not found"))
		}
		params := dbx.Params{"client": invoice.GetString("client")}
		job := e.Request.URL.Query().Get("job")
		if job == "" {
			jobs := []billableJob{}
			if err := app.DB().NewQuery(billableJobsQuery).Bind(params).All(&jobs); err != nil {
				return writeInvoiceError(e, err)
			}
			return e.JSON(http.StatusOK, jobs)
		}
		params["job"] = job
		costs := []invoiceCost{}
		if err := app.DB().NewQuery(billableCostsQuery).Bind(params).All(&costs); err != nil {
			return writeInvoiceError(e, err)
		}
		return e.JSON(http.StatusOK, costs)
	}
}

// costInvoiceView is the invoice an expense or PO was billed on, for its
// details page.
type costInvoiceView struct {
	Line *costInvoiceLine `json:"line"`
	// CanBill is true when the viewer may bill this unbilled cost now.
	CanBill bool `json:"can_bill"`
	// CanOpenInvoice is true for holders of invoicing or report.
	CanOpenInvoice bool    `json:"can_open_invoice"`
	ClientID       string  `json:"client_id"`
	CostTotal      float64 `json:"cost_total"`
	CurrencyCode   string  `json:"currency_code"`
}

type costInvoiceLine struct {
	ID            string  `db:"id" json:"id"`
	Amount        float64 `db:"amount" json:"amount"`
	InvoiceID     string  `db:"invoice_id" json:"invoice_id"`
	InvoiceNumber string  `db:"invoice_number" json:"invoice_number"`
	InvoiceDate   string  `db:"invoice_date" json:"invoice_date"`
}

const costInvoiceLineQuery = `
SELECT l.id, CAST(l.amount AS REAL) AS amount, i.id AS invoice_id,
  COALESCE(i.number, '') AS invoice_number, COALESCE(i.date, '') AS invoice_date
FROM invoice_lines l
JOIN invoices i ON i.id = l.invoice
WHERE l.expense = {:cost} OR l.purchase_order = {:cost}
LIMIT 1
`

// buildCostInvoiceView loads the line billing a cost the viewer can open,
// and whether the viewer may bill it. billable is the cost's own eligibility:
// committed for an expense, Active or Closed for a PO.
func buildCostInvoiceView(app core.App, auth *core.Record, cost *core.Record, billable bool) (costInvoiceView, error) {
	view := costInvoiceView{CostTotal: cost.GetFloat("total"), CurrencyCode: "CAD"}
	if currency, err := app.FindRecordById("currencies", cost.GetString("currency")); err == nil {
		view.CurrencyCode = currency.GetString("code")
	}
	var line costInvoiceLine
	if err := app.DB().NewQuery(costInvoiceLineQuery).Bind(dbx.Params{"cost": cost.Id}).One(&line); err == nil {
		view.Line = &line
	}
	if err := requireInvoiceViewer(app, auth); err == nil {
		view.CanOpenInvoice = true
	}
	job, err := app.FindRecordById("jobs", cost.GetString("job"))
	if err != nil {
		return view, nil
	}
	view.ClientID = job.GetString("client")
	view.CanBill = view.Line == nil && billable && hooks.RequireActiveInvoicingUser(app, auth) == nil
	return view, nil
}

// createGetExpenseClientInvoiceHandler returns the invoice an expense the
// viewer can open was billed on.
func createGetExpenseClientInvoiceHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		params, err := expenseVisibilityParams(app, e.Auth)
		if err != nil {
			return writeInvoiceError(e, err)
		}
		hasAdmin, err := utilities.HasClaim(app, e.Auth, "admin")
		if err != nil {
			return writeInvoiceError(e, err)
		}
		params["id"] = e.Request.PathValue("id")
		params["has_admin"] = boolToInt(hasAdmin)
		var found struct {
			ID string `db:"id"`
		}
		if err := app.DB().NewQuery(`SELECT e.id FROM expenses e WHERE e.id = {:id} AND (` +
			expenseVisibilityPredicate + ` OR ({:has_admin} = 1 AND e.committed != ''))`).Bind(params).One(&found); err != nil {
			return writeInvoiceError(e, invoiceRouteError(http.StatusNotFound, "expense not found or not authorized"))
		}
		expense, err := app.FindRecordById("expenses", found.ID)
		if err != nil {
			return writeInvoiceError(e, err)
		}
		view, err := buildCostInvoiceView(app, e.Auth, expense, expense.GetString("committed") != "")
		if err != nil {
			return writeInvoiceError(e, err)
		}
		return e.JSON(http.StatusOK, view)
	}
}

// createGetPurchaseOrderClientInvoiceHandler returns the invoice a PO the
// viewer can open was billed on.
func createGetPurchaseOrderClientInvoiceHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		id := e.Request.PathValue("id")
		row, err := findVisiblePurchaseOrderByID(app, e.Auth.Id, id)
		if err != nil {
			return writeInvoiceError(e, err)
		}
		if row == nil {
			return writeInvoiceError(e, invoiceRouteError(http.StatusNotFound, "purchase order not found or not visible"))
		}
		po, err := app.FindRecordById("purchase_orders", id)
		if err != nil {
			return writeInvoiceError(e, err)
		}
		status := po.GetString("status")
		view, err := buildCostInvoiceView(app, e.Auth, po, status == "Active" || status == "Closed")
		if err != nil {
			return writeInvoiceError(e, err)
		}
		return e.JSON(http.StatusOK, view)
	}
}
