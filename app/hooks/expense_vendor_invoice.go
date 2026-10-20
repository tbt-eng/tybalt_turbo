package hooks

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"

	"tybalt/constants"
	"tybalt/errs"
	"tybalt/utilities"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

// Vendor invoice details live on the expense document that substantiates an
// expense: the invoice is the document, whoever enters the expense. See the
// vendor invoice section of descriptions/expense_documents.md.

// The expense write endpoints accept the virtual inputs vendor_invoice_number,
// vendor_invoice_date and vendor_invoice_total alongside "attachment". They are
// not expense fields.

const vendorInvoiceNumberMaxLength = 100

type vendorInvoiceDetails struct {
	Vendor   string
	Number   string
	Key      string
	Date     string
	Total    float64
	Currency string
}

func (d vendorInvoiceDetails) isInvoice() bool {
	return d.Number != ""
}

func vendorInvoiceError(status int, field, code, message string) *errs.HookError {
	return &errs.HookError{
		Status:  status,
		Message: message,
		Data:    map[string]errs.CodeError{field: {Code: code, Message: message}},
	}
}

// VendorInvoiceNumberKey reduces an invoice number to lowercase letters and
// digits, so "INV-0042", "inv 0042" and "INV0042" match. Leading zeros stay.
func VendorInvoiceNumberKey(number string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(number) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func documentVendorInvoice(document *core.Record) vendorInvoiceDetails {
	return vendorInvoiceDetails{
		Vendor:   document.GetString("vendor"),
		Number:   document.GetString("vendor_invoice_number"),
		Key:      document.GetString("vendor_invoice_number_key"),
		Date:     document.GetString("vendor_invoice_date"),
		Total:    document.GetFloat("vendor_invoice_total"),
		Currency: document.GetString("vendor_invoice_currency"),
	}
}

// readVendorInvoiceInputs applies the inputs present in the request body to the
// document's current details. An input left out keeps the current value.
func readVendorInvoiceInputs(body map[string]any, current vendorInvoiceDetails) (vendorInvoiceDetails, bool, error) {
	next := current
	provided := false
	if raw, ok := body["vendor_invoice_number"]; ok {
		provided = true
		next.Number = strings.TrimSpace(stringInput(raw))
	}
	if raw, ok := body["vendor_invoice_date"]; ok {
		provided = true
		next.Date = strings.TrimSpace(stringInput(raw))
	}
	if raw, ok := body["vendor_invoice_total"]; ok {
		provided = true
		total, err := numberInput(raw)
		if err != nil {
			return next, provided, vendorInvoiceError(http.StatusBadRequest, "vendor_invoice_total", "invalid", "Enter the invoice total as a number.")
		}
		next.Total = total
	}
	return next, provided, nil
}

func stringInput(raw any) string {
	switch v := raw.(type) {
	case nil:
		return ""
	case string:
		return v
	case []string:
		if len(v) > 0 {
			return v[0]
		}
		return ""
	default:
		return fmt.Sprint(v)
	}
}

func numberInput(raw any) (float64, error) {
	switch v := raw.(type) {
	case nil:
		return 0, nil
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	default:
		s := strings.TrimSpace(stringInput(v))
		if s == "" {
			return 0, nil
		}
		return strconv.ParseFloat(s, 64)
	}
}

// applyVendorInvoiceDetails writes the request's vendor invoice inputs to the
// expense's document and checks every vendor invoice rule. It runs after the
// document is resolved, in the same transaction as the expense save.
func applyVendorInvoiceDetails(app core.App, e *core.RecordRequestEvent, body map[string]any, hasBookKeeperClaim bool) error {
	expense := e.Record
	documentID := expense.GetString("attachment_document")
	if documentID == "" {
		next, provided, err := readVendorInvoiceInputs(body, vendorInvoiceDetails{})
		if err != nil {
			return err
		}
		if provided && (next.Number != "" || next.Date != "" || next.Total != 0) {
			return vendorInvoiceError(http.StatusBadRequest, "vendor_invoice_number", "document_required", "Attach the invoice before entering its details.")
		}
		return nil
	}

	document, err := app.FindRecordById(constants.ExpenseDocumentsCollectionName, documentID)
	if err != nil {
		return err
	}
	current := documentVendorInvoice(document)
	next, _, err := readVendorInvoiceInputs(body, current)
	if err != nil {
		return err
	}

	others, err := otherDocumentExpenses(app, documentID, expense.Id)
	if err != nil {
		return err
	}

	changed := next.Number != current.Number || next.Date != current.Date ||
		utilities.CurrencyAmountToCents(next.Total) != utilities.CurrencyAmountToCents(current.Total)
	if changed {
		if err := checkVendorInvoiceEditable(document, current, others, e.Auth, hasBookKeeperClaim); err != nil {
			return err
		}
	}

	if err := validateVendorInvoiceFields(expense, &next); err != nil {
		return err
	}

	if next.isInvoice() {
		if err := matchVendorInvoiceVendorAndCurrency(app, expense, others, &next); err != nil {
			return err
		}
		next.Key = VendorInvoiceNumberKey(next.Number)
		if next.Key == "" {
			return vendorInvoiceError(http.StatusBadRequest, "vendor_invoice_number", "invalid", "The invoice number must contain letters or digits.")
		}
		if err := checkVendorInvoiceUnique(app, document.Id, next); err != nil {
			return err
		}
		if err := checkVendorInvoiceTotal(expense, others, next); err != nil {
			return err
		}
	} else {
		// A receipt has no vendor invoice details.
		next = vendorInvoiceDetails{}
	}

	if next == current {
		return nil
	}
	now := types.NowDateTime()
	recordedBy := ""
	if e.Auth != nil {
		recordedBy = e.Auth.Id
	}
	_, err = app.DB().Update(constants.ExpenseDocumentsCollectionName, dbx.Params{
		"vendor":                    next.Vendor,
		"vendor_invoice_number":     next.Number,
		"vendor_invoice_number_key": next.Key,
		"vendor_invoice_date":       next.Date,
		"vendor_invoice_total":      next.Total,
		"vendor_invoice_currency":   next.Currency,
		// Who last set or changed the details, and when. Not a history.
		"vendor_invoice_recorded_by": recordedBy,
		"vendor_invoice_recorded":    now,
		"updated":                    now,
	}, dbx.HashExp{"id": document.Id}).Execute()
	return err
}

// otherDocumentExpenses returns the other expenses that use a document.
func otherDocumentExpenses(app core.App, documentID, expenseID string) ([]*core.Record, error) {
	conditions := []dbx.Expression{dbx.HashExp{"attachment_document": documentID}}
	if expenseID != "" {
		conditions = append(conditions, dbx.Not(dbx.HashExp{"id": expenseID}))
	}
	return app.FindAllRecords("expenses", conditions...)
}

func isRejected(expense *core.Record) bool {
	return !expense.GetDateTime("rejected").IsZero()
}

// checkVendorInvoiceEditable allows the uploader or a bookkeeper to change the
// details, or anyone when no other expense uses the document. The details lock
// once another expense on the document is submitted and not rejected, because
// that expense is waiting for approval against them.
//
// One exception: a bookkeeper may fill in an invoice that has no details at
// all. Invoices uploaded before vendor invoice details existed have none, and a
// bookkeeper reusing one for another PO must add them even though its earlier
// expense is already submitted. Filling blank details changes nothing an
// approver saw. Only bookkeepers reuse documents across POs, so the exception
// is theirs. Existing details stay locked for everyone.
func checkVendorInvoiceEditable(document *core.Record, current vendorInvoiceDetails, others []*core.Record, auth *core.Record, hasBookKeeperClaim bool) error {
	if len(others) > 0 && !hasBookKeeperClaim && (auth == nil || document.GetString("uploaded_by") != auth.Id) {
		return vendorInvoiceError(http.StatusForbidden, "vendor_invoice_number", "not_editable",
			"Only the person who uploaded this invoice or a bookkeeper can change its details.")
	}
	blank := current.Number == "" && current.Date == "" && utilities.CurrencyAmountToCents(current.Total) == 0
	if blank && hasBookKeeperClaim {
		return nil
	}
	for _, other := range others {
		if other.GetBool("submitted") && !isRejected(other) {
			return vendorInvoiceError(http.StatusBadRequest, "vendor_invoice_number", "locked",
				"These invoice details cannot change: another expense on this invoice has been submitted for approval.")
		}
	}
	return nil
}

// validateVendorInvoiceFields checks formats, and that the details are
// complete: required on an OnAccount expense, and all or nothing otherwise.
func validateVendorInvoiceFields(expense *core.Record, next *vendorInvoiceDetails) error {
	if len([]rune(next.Number)) > vendorInvoiceNumberMaxLength {
		return vendorInvoiceError(http.StatusBadRequest, "vendor_invoice_number", "too_long", "Keep the invoice number to 100 characters or fewer.")
	}
	if next.Date != "" {
		if _, err := time.Parse("2006-01-02", next.Date); err != nil {
			return vendorInvoiceError(http.StatusBadRequest, "vendor_invoice_date", "invalid", "Enter a valid invoice date.")
		}
	}
	if next.Total < 0 {
		return vendorInvoiceError(http.StatusBadRequest, "vendor_invoice_total", "invalid", "The invoice total cannot be negative.")
	}

	required := expense.GetString("payment_type") == "OnAccount" ||
		next.Number != "" || next.Date != "" || next.Total != 0
	if !required {
		return nil
	}
	data := map[string]errs.CodeError{}
	if next.Number == "" {
		data["vendor_invoice_number"] = errs.CodeError{Code: "required", Message: "Enter the vendor's invoice number."}
	}
	if next.Date == "" {
		data["vendor_invoice_date"] = errs.CodeError{Code: "required", Message: "Enter the invoice date."}
	}
	if utilities.CurrencyAmountToCents(next.Total) <= 0 {
		data["vendor_invoice_total"] = errs.CodeError{Code: "required", Message: "Enter the invoice total."}
	}
	if len(data) > 0 {
		message := "Enter the vendor invoice details."
		if expense.GetString("payment_type") == "OnAccount" {
			message = "An on-account expense needs its vendor invoice details."
		}
		return &errs.HookError{Status: http.StatusBadRequest, Message: message, Data: data}
	}
	return nil
}

// matchVendorInvoiceVendorAndCurrency takes the vendor and currency from the
// expense while it is the only expense on the document, so correcting the
// expense corrects the invoice. Once the invoice is shared, every expense on it
// must use the document's vendor and currency.
func matchVendorInvoiceVendorAndCurrency(app core.App, expense *core.Record, others []*core.Record, next *vendorInvoiceDetails) error {
	if next.Vendor == "" || len(others) == 0 {
		next.Vendor = expense.GetString("vendor")
		next.Currency = expense.GetString("currency")
	}
	if next.Vendor == "" {
		return vendorInvoiceError(http.StatusBadRequest, "vendor", "required", "Select the vendor that issued this invoice.")
	}
	for _, record := range append([]*core.Record{expense}, others...) {
		if record.GetString("vendor") != next.Vendor {
			name := next.Vendor
			if vendor, err := app.FindRecordById("vendors", next.Vendor); err == nil {
				name = vendor.GetString("name")
			}
			return vendorInvoiceError(http.StatusBadRequest, "vendor", "vendor_mismatch",
				fmt.Sprintf("This invoice is from %s. Every expense on it must use that vendor.", name))
		}
		if record.GetString("currency") != next.Currency {
			return vendorInvoiceError(http.StatusBadRequest, "currency", "currency_mismatch",
				"Every expense on this invoice must use the invoice's currency.")
		}
	}
	return nil
}

// checkVendorInvoiceUnique refuses a second document for the same vendor and
// invoice number. The error names the PO and who entered it, not a link,
// because the user may not be allowed to view that expense.
func checkVendorInvoiceUnique(app core.App, documentID string, next vendorInvoiceDetails) error {
	existing, err := app.FindFirstRecordByFilter(constants.ExpenseDocumentsCollectionName,
		"vendor = {:vendor} && vendor_invoice_number_key = {:key} && id != {:id}",
		dbx.Params{"vendor": next.Vendor, "key": next.Key, "id": documentID})
	if err != nil || existing == nil {
		return nil
	}
	message := fmt.Sprintf("Invoice %s from this vendor has already been entered", existing.GetString("vendor_invoice_number"))
	if expenses, err := app.FindAllRecords("expenses", dbx.HashExp{"attachment_document": existing.Id}); err == nil && len(expenses) > 0 {
		first := expenses[0]
		if po, err := app.FindRecordById("purchase_orders", first.GetString("purchase_order")); err == nil && po.GetString("po_number") != "" {
			message += " against PO " + po.GetString("po_number")
		}
		if name := userDisplayName(app, first.GetString("creator")); name != "" {
			message += " by " + name
		}
	}
	return vendorInvoiceError(http.StatusBadRequest, "vendor_invoice_number", "duplicate_invoice", message+".")
}

func userDisplayName(app core.App, uid string) string {
	if uid == "" {
		return ""
	}
	profile, err := app.FindFirstRecordByData("profiles", "uid", uid)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(profile.GetString("given_name") + " " + profile.GetString("surname"))
}

// checkVendorInvoiceTotal refuses expenses on one invoice that together cost
// more than the invoice. Rejected expenses don't count. Amounts are compared in
// whole cents, so 16.26 + 0.01 matches a 16.27 invoice.
func checkVendorInvoiceTotal(expense *core.Record, others []*core.Record, next vendorInvoiceDetails) error {
	allocated := int64(0)
	if !isRejected(expense) {
		allocated += utilities.CurrencyAmountToCents(expense.GetFloat("total"))
	}
	for _, other := range others {
		if !isRejected(other) {
			allocated += utilities.CurrencyAmountToCents(other.GetFloat("total"))
		}
	}
	invoice := utilities.CurrencyAmountToCents(next.Total)
	if allocated > invoice {
		return vendorInvoiceError(http.StatusBadRequest, "total", "exceeds_vendor_invoice",
			fmt.Sprintf("The expenses on this invoice would total %.2f, more than the invoice total of %.2f.",
				float64(allocated)/100, float64(invoice)/100))
	}
	return nil
}
