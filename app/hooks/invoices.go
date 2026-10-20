package hooks

import (
	"net/http"
	"regexp"
	"time"

	"tybalt/errs"
	"tybalt/utilities"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

// Client invoices and the costs billed back on them. See
// descriptions/invoices.md.

var (
	invoiceNumberPattern = regexp.MustCompile(`^\d{2}(0[1-9]|1[0-2])\d{3}$`)
	invoiceDatePattern   = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
)

func invoiceError(status int, field, code, message string) *errs.HookError {
	err := &errs.HookError{Status: status, Message: message}
	if field != "" {
		err.Data = map[string]errs.CodeError{field: {Code: code, Message: message}}
	}
	return err
}

// RequireActiveInvoicingUser refuses a user who is inactive or lacks the
// invoicing claim. The collection rules check the claim too; routes that save
// with the app rely on this.
func RequireActiveInvoicingUser(app core.App, auth *core.Record) error {
	if auth == nil {
		return invoiceError(http.StatusUnauthorized, "", "", "authentication required")
	}
	active, err := utilities.IsUserActive(app, auth.Id)
	if err != nil {
		return err
	}
	if !active {
		return invoiceError(http.StatusForbidden, "", "user_inactive", "inactive users cannot change invoices")
	}
	has, err := utilities.HasClaim(app, auth, "invoicing")
	if err != nil {
		return err
	}
	if !has {
		return invoiceError(http.StatusForbidden, "", "invoicing_claim_required", "the invoicing claim is required")
	}
	return nil
}

// ProcessInvoiceCreate starts every invoice as a draft owned by its creator.
func ProcessInvoiceCreate(app core.App, e *core.RecordRequestEvent) error {
	if err := RequireActiveInvoicingUser(app, e.Auth); err != nil {
		return err
	}
	e.Record.Set("status", "Draft")
	e.Record.Set("creator", e.Auth.Id)
	e.Record.Set("issued_by", "")
	e.Record.Set("issued_at", "")
	return validateInvoiceDraft(app, e.Record)
}

// ProcessInvoiceUpdate lets a draft's client, number, date and total change.
// Status and the server-set fields keep their saved values; issuing goes
// through the issue route.
func ProcessInvoiceUpdate(app core.App, e *core.RecordRequestEvent) error {
	if err := RequireActiveInvoicingUser(app, e.Auth); err != nil {
		return err
	}
	original := e.Record.Original()
	if original.GetString("status") != "Draft" {
		return invoiceError(http.StatusBadRequest, "status", "invoice_issued", "an issued invoice cannot be edited")
	}
	for _, field := range []string{"status", "creator", "issued_by", "issued_at"} {
		e.Record.Set(field, original.Get(field))
	}
	return validateInvoiceDraft(app, e.Record)
}

// ProcessInvoiceDelete allows deleting a draft only. Deleting it frees its
// number.
func ProcessInvoiceDelete(app core.App, e *core.RecordRequestEvent) error {
	if err := RequireActiveInvoicingUser(app, e.Auth); err != nil {
		return err
	}
	if e.Record.GetString("status") != "Draft" {
		return invoiceError(http.StatusBadRequest, "status", "invoice_issued", "an issued invoice cannot be deleted")
	}
	return nil
}

func validateInvoiceDraft(app core.App, invoice *core.Record) error {
	if invoice.GetString("client") == "" {
		return invoiceError(http.StatusBadRequest, "client", "required", "client is required")
	}
	if err := checkInvoiceNumber(app, invoice); err != nil {
		return err
	}
	if date := invoice.GetString("date"); date != "" && !validInvoiceDate(date) {
		return invoiceError(http.StatusBadRequest, "date", "invalid_date", "date must be a valid YYYY-MM-DD date")
	}
	if invoice.GetFloat("total") < 0 {
		return invoiceError(http.StatusBadRequest, "total", "negative_total", "total cannot be negative")
	}
	return nil
}

func validInvoiceDate(date string) bool {
	if !invoiceDatePattern.MatchString(date) {
		return false
	}
	_, err := time.Parse("2006-01-02", date)
	return err == nil
}

// checkInvoiceNumber checks an invoice's number, when it has one, for its
// YYMMNNN form and that no other invoice uses it.
func checkInvoiceNumber(app core.App, invoice *core.Record) error {
	number := invoice.GetString("number")
	if number == "" {
		return nil
	}
	if !invoiceNumberPattern.MatchString(number) {
		return invoiceError(http.StatusBadRequest, "number", "invalid_number",
			"invoice number must be seven digits, YYMMNNN, with a month from 01 to 12")
	}
	var other struct {
		ID string `db:"id"`
	}
	err := app.DB().NewQuery(`SELECT id FROM invoices WHERE number = {:number} AND id != {:id} LIMIT 1`).
		Bind(dbx.Params{"number": number, "id": invoice.Id}).One(&other)
	if err == nil {
		return invoiceError(http.StatusBadRequest, "number", "number_taken",
			"invoice number "+number+" is already used by another invoice")
	}
	return nil
}

// ValidateInvoiceForIssue checks that a draft has everything an issued
// invoice needs.
func ValidateInvoiceForIssue(app core.App, invoice *core.Record) error {
	if invoice.GetString("status") != "Draft" {
		return invoiceError(http.StatusBadRequest, "status", "invoice_issued", "the invoice is already issued")
	}
	if invoice.GetString("number") == "" {
		return invoiceError(http.StatusBadRequest, "number", "required", "an invoice needs a number before it is issued")
	}
	if invoice.GetString("date") == "" {
		return invoiceError(http.StatusBadRequest, "date", "required", "an invoice needs a date before it is issued")
	}
	if utilities.CurrencyAmountToCents(invoice.GetFloat("total")) <= 0 {
		return invoiceError(http.StatusBadRequest, "total", "required", "an invoice needs a total greater than zero before it is issued")
	}
	return validateInvoiceDraft(app, invoice)
}

// ProcessInvoiceLineCreate records who billed the cost and checks the line
// rules.
func ProcessInvoiceLineCreate(app core.App, e *core.RecordRequestEvent) error {
	if err := RequireActiveInvoicingUser(app, e.Auth); err != nil {
		return err
	}
	e.Record.Set("creator", e.Auth.Id)
	return ValidateInvoiceLine(app, e.Record)
}

// ProcessInvoiceLineDelete lets an invoicing holder remove a line.
func ProcessInvoiceLineDelete(app core.App, e *core.RecordRequestEvent) error {
	return RequireActiveInvoicingUser(app, e.Auth)
}

// ValidateInvoiceLine checks that a line bills one eligible cost of the
// invoice's client on an issued invoice, and that the cost is not billed yet.
func ValidateInvoiceLine(app core.App, line *core.Record) error {
	invoice, err := app.FindRecordById("invoices", line.GetString("invoice"))
	if err != nil {
		return invoiceError(http.StatusBadRequest, "invoice", "not_found", "invoice not found")
	}
	if invoice.GetString("status") != "Issued" {
		return invoiceError(http.StatusBadRequest, "invoice", "invoice_not_issued", "costs can be billed only on an issued invoice")
	}
	if utilities.CurrencyAmountToCents(line.GetFloat("amount")) <= 0 {
		return invoiceError(http.StatusBadRequest, "amount", "required", "the amount billed must be greater than zero")
	}

	var collection, field, otherField, id string
	switch line.GetString("kind") {
	case "expense":
		collection, field, otherField = "expenses", "expense", "purchase_order"
	case "po":
		collection, field, otherField = "purchase_orders", "purchase_order", "expense"
	default:
		return invoiceError(http.StatusBadRequest, "kind", "invalid_kind", "kind must be expense or po")
	}
	id = line.GetString(field)
	if id == "" {
		return invoiceError(http.StatusBadRequest, field, "required", field+" is required for a line of kind "+line.GetString("kind"))
	}
	if line.GetString(otherField) != "" {
		return invoiceError(http.StatusBadRequest, otherField, "not_allowed", "a line bills one expense or one purchase order, not both")
	}

	cost, err := app.FindRecordById(collection, id)
	if err != nil {
		return invoiceError(http.StatusBadRequest, field, "not_found", field+" not found")
	}
	if collection == "expenses" && cost.GetString("committed") == "" {
		return invoiceError(http.StatusBadRequest, field, "expense_not_committed", "only a committed expense can be billed")
	}
	if collection == "purchase_orders" {
		if status := cost.GetString("status"); status != "Active" && status != "Closed" {
			return invoiceError(http.StatusBadRequest, field, "po_not_billable", "only an Active or Closed purchase order can be billed")
		}
	}
	if err := checkCostClient(app, cost.GetString("job"), invoice.GetString("client"), field); err != nil {
		return err
	}

	var other struct {
		Number string `db:"number"`
	}
	err = app.DB().NewQuery(`SELECT COALESCE(i.number, '') AS number FROM invoice_lines l
		JOIN invoices i ON i.id = l.invoice
		WHERE l.` + field + ` = {:cost} AND l.id != {:id} LIMIT 1`).
		Bind(dbx.Params{"cost": id, "id": line.Id}).One(&other)
	if err == nil {
		return invoiceError(http.StatusBadRequest, field, "already_billed", "this cost is already billed on invoice "+other.Number)
	}
	return nil
}

func checkCostClient(app core.App, jobID, clientID, field string) error {
	if jobID == "" {
		return invoiceError(http.StatusBadRequest, field, "no_job", "only a cost on a job can be billed to a client")
	}
	job, err := app.FindRecordById("jobs", jobID)
	if err != nil {
		return invoiceError(http.StatusBadRequest, field, "no_job", "the cost's job was not found")
	}
	if job.GetString("client") != clientID {
		return invoiceError(http.StatusBadRequest, field, "client_mismatch", "the cost's job "+job.GetString("number")+" belongs to a different client than the invoice")
	}
	return nil
}
