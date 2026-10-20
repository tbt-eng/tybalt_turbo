package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"tybalt/internal/testutils"

	"github.com/pocketbase/pocketbase/tests"
)

// Fixtures: time@test.com (rzr98oadsp9qc11, "Tester Time") owns the active
// One-Time PO poa1ctvbrnch001 (2026-0901, 132.10, vendor z66xe6vqhwtokt4).
// book@keeper.com (tqqf7q0f3378rvp) holds book_keeper.

func vendorInvoiceExpenseFields(description, total string, invoice map[string]string) map[string]string {
	fields := bookKeeperExpenseFields(map[string]string{"description": description, "total": total})
	for key, value := range invoice {
		fields[key] = value
	}
	return fields
}

func postVendorInvoiceExpense(t *testing.T, app *tests.TestApp, token string, fields map[string]string, content []byte) (int, string, string) {
	t.Helper()
	body, contentType := mustMultipartExpenseWithContent(t, fields, "invoice.png", content)
	res := performTestAPIRequest(t, app, http.MethodPost, "/api/expenses", body, map[string]string{
		"Authorization": token,
		"Content-Type":  contentType,
	})
	var created struct {
		ID string `json:"id"`
	}
	_ = json.Unmarshal(res.Body.Bytes(), &created)
	return res.Code, res.Body.String(), created.ID
}

func pngContent(seed byte) []byte {
	return []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x40, seed}
}

func TestVendorInvoiceDetails(t *testing.T) {
	ownerToken, err := testutils.GenerateRecordToken("users", "time@test.com")
	if err != nil {
		t.Fatal(err)
	}
	bookkeeperToken, err := testutils.GenerateRecordToken("users", "book@keeper.com")
	if err != nil {
		t.Fatal(err)
	}
	invoice := func(number, date, total string) map[string]string {
		return map[string]string{"vendor_invoice_number": number, "vendor_invoice_date": date, "vendor_invoice_total": total}
	}

	t.Run("an on-account expense needs its invoice details", func(t *testing.T) {
		app := testutils.SetupTestApp(t)
		t.Cleanup(app.Cleanup)
		status, body, _ := postVendorInvoiceExpense(t, app, ownerToken,
			vendorInvoiceExpenseFields("no invoice details", "25", invoice("", "", "")), pngContent(1))
		if status != http.StatusBadRequest || !strings.Contains(body, `"vendor_invoice_number":{"code":"required"`) {
			t.Fatalf("expected required invoice details, got %d %s", status, body)
		}
	})

	t.Run("details are stored on the document with a normalized key", func(t *testing.T) {
		app := testutils.SetupTestApp(t)
		t.Cleanup(app.Cleanup)
		status, body, id := postVendorInvoiceExpense(t, app, ownerToken,
			vendorInvoiceExpenseFields("owner entered invoice", "25", invoice(" INV-0042 ", "2024-08-30", "25")), pngContent(2))
		if status != http.StatusOK {
			t.Fatalf("create failed: %d %s", status, body)
		}
		expense, err := app.FindRecordById("expenses", id)
		if err != nil {
			t.Fatal(err)
		}
		document, err := app.FindRecordById("expense_documents", expense.GetString("attachment_document"))
		if err != nil {
			t.Fatal(err)
		}
		for field, want := range map[string]string{
			"vendor":                    "z66xe6vqhwtokt4",
			"vendor_invoice_number":     "INV-0042",
			"vendor_invoice_number_key": "inv0042",
			"vendor_invoice_date":       "2024-08-30",
			"vendor_invoice_currency":   expense.GetString("currency"),
		} {
			if got := document.GetString(field); got != want {
				t.Errorf("%s = %q; want %q", field, got, want)
			}
		}
		if got := document.GetFloat("vendor_invoice_total"); got != 25 {
			t.Errorf("vendor_invoice_total = %v; want 25", got)
		}
	})

	t.Run("the same invoice entered by the PO owner and a bookkeeper is refused", func(t *testing.T) {
		app := testutils.SetupTestApp(t)
		t.Cleanup(app.Cleanup)
		status, body, _ := postVendorInvoiceExpense(t, app, ownerToken,
			vendorInvoiceExpenseFields("owner entered invoice", "25", invoice("INV-0042", "2024-08-30", "25")), pngContent(3))
		if status != http.StatusOK {
			t.Fatalf("owner create failed: %d %s", status, body)
		}
		// A different file of the same invoice, typed differently.
		status, body, _ = postVendorInvoiceExpense(t, app, bookkeeperToken,
			vendorInvoiceExpenseFields("bookkeeper entered invoice", "25", invoice("inv 0042", "2024-08-30", "25")), pngContent(4))
		if status != http.StatusBadRequest || !strings.Contains(body, "duplicate_invoice") ||
			!strings.Contains(body, "2026-0901") || !strings.Contains(body, "Tester Time") {
			t.Fatalf("expected a duplicate invoice error naming the PO and who entered it, got %d %s", status, body)
		}
	})

	t.Run("expenses that exactly fill an invoice are accepted, and one cent more is refused", func(t *testing.T) {
		app := testutils.SetupTestApp(t)
		t.Cleanup(app.Cleanup)
		details := invoice("INV-1627", "2024-08-30", "16.27")
		content := pngContent(5)
		status, body, _ := postVendorInvoiceExpense(t, app, bookkeeperToken,
			vendorInvoiceExpenseFields("first part", "16.26", details), content)
		if status != http.StatusOK {
			t.Fatalf("first part failed: %d %s", status, body)
		}
		// The bookkeeper reuses the same file for the rest of the invoice.
		status, body, _ = postVendorInvoiceExpense(t, app, bookkeeperToken,
			vendorInvoiceExpenseFields("last cent", "0.01", details), content)
		if status != http.StatusOK {
			t.Fatalf("an exact fill should be accepted: %d %s", status, body)
		}
		status, body, _ = postVendorInvoiceExpense(t, app, bookkeeperToken,
			vendorInvoiceExpenseFields("one cent over", "0.01", details), content)
		if status != http.StatusBadRequest || !strings.Contains(body, "exceeds_vendor_invoice") {
			t.Fatalf("expected the invoice total to be exceeded, got %d %s", status, body)
		}
	})

	t.Run("details lock once another expense on the invoice is submitted", func(t *testing.T) {
		app := testutils.SetupTestApp(t)
		t.Cleanup(app.Cleanup)
		content := pngContent(6)
		status, body, firstID := postVendorInvoiceExpense(t, app, bookkeeperToken,
			vendorInvoiceExpenseFields("first part", "10", invoice("INV-LOCK", "2024-08-30", "30")), content)
		if status != http.StatusOK {
			t.Fatalf("first part failed: %d %s", status, body)
		}
		submit := performTestAPIRequest(t, app, http.MethodPost, "/api/expenses/"+firstID+"/submit", nil, map[string]string{
			"Authorization": bookkeeperToken,
		})
		mustStatus(t, submit, http.StatusOK)

		status, body, _ = postVendorInvoiceExpense(t, app, bookkeeperToken,
			vendorInvoiceExpenseFields("second part", "10", invoice("INV-LOCK", "2024-08-30", "40")), content)
		if status != http.StatusBadRequest || !strings.Contains(body, `"code":"locked"`) {
			t.Fatalf("expected the details to be locked, got %d %s", status, body)
		}
		// Reusing the invoice without changing its details still works.
		status, body, _ = postVendorInvoiceExpense(t, app, bookkeeperToken,
			vendorInvoiceExpenseFields("second part", "10", invoice("INV-LOCK", "2024-08-30", "30")), content)
		if status != http.StatusOK {
			t.Fatalf("reusing unchanged details should work: %d %s", status, body)
		}
	})

	t.Run("the uploader can still change an invoice a bookkeeper reused", func(t *testing.T) {
		app := testutils.SetupTestApp(t)
		t.Cleanup(app.Cleanup)
		content := pngContent(7)
		status, body, ownerExpenseID := postVendorInvoiceExpense(t, app, ownerToken,
			vendorInvoiceExpenseFields("owner part", "10", invoice("INV-SHARED", "2024-08-30", "30")), content)
		if status != http.StatusOK {
			t.Fatalf("owner part failed: %d %s", status, body)
		}
		status, body, _ = postVendorInvoiceExpense(t, app, bookkeeperToken,
			vendorInvoiceExpenseFields("bookkeeper part", "10", invoice("INV-SHARED", "2024-08-30", "30")), content)
		if status != http.StatusOK {
			t.Fatalf("bookkeeper reuse failed: %d %s", status, body)
		}
		update := performTestAPIRequest(t, app, http.MethodPatch, "/api/expenses/"+ownerExpenseID,
			strings.NewReader(`{"vendor_invoice_total": 50}`), map[string]string{
				"Authorization": ownerToken,
				"Content-Type":  "application/json",
			})
		mustStatus(t, update, http.StatusOK)
		var total float64
		if err := app.DB().NewQuery("SELECT ed.vendor_invoice_total FROM expenses e JOIN expense_documents ed ON ed.id = e.attachment_document WHERE e.id = {:id}").
			Bind(map[string]any{"id": ownerExpenseID}).Row(&total); err != nil {
			t.Fatal(err)
		}
		if total != 50 {
			t.Fatalf("vendor_invoice_total = %v; want 50", total)
		}
	})

	t.Run("only a bookkeeper fills in a blank invoice that is already in use", func(t *testing.T) {
		app := testutils.SetupTestApp(t)
		t.Cleanup(app.Cleanup)
		details := invoice("INV-OLD", "2024-08-30", "30")
		content := pngContent(12)
		status, body, ownerExpenseID := postVendorInvoiceExpense(t, app, ownerToken,
			vendorInvoiceExpenseFields("owner part", "10", details), content)
		if status != http.StatusOK {
			t.Fatalf("owner part failed: %d %s", status, body)
		}
		status, body, bookkeeperExpenseID := postVendorInvoiceExpense(t, app, bookkeeperToken,
			vendorInvoiceExpenseFields("bookkeeper part", "10", details), content)
		if status != http.StatusOK {
			t.Fatalf("bookkeeper part failed: %d %s", status, body)
		}
		// Make it an invoice uploaded before vendor invoice details existed: no
		// details, and an expense on it already submitted.
		if _, err := app.DB().NewQuery(`UPDATE expense_documents SET vendor = '', vendor_invoice_number = '',
			vendor_invoice_number_key = '', vendor_invoice_date = '', vendor_invoice_total = 0, vendor_invoice_currency = ''
			WHERE id = (SELECT attachment_document FROM expenses WHERE id = {:id})`).
			Bind(map[string]any{"id": ownerExpenseID}).Execute(); err != nil {
			t.Fatal(err)
		}
		submit := performTestAPIRequest(t, app, http.MethodPost, "/api/expenses/"+bookkeeperExpenseID+"/submit", nil, map[string]string{
			"Authorization": bookkeeperToken,
		})
		mustStatus(t, submit, http.StatusOK)

		// The uploader is not a bookkeeper, so the blank invoice stays locked.
		update := performTestAPIRequest(t, app, http.MethodPatch, "/api/expenses/"+ownerExpenseID,
			strings.NewReader(`{"vendor_invoice_number": "INV-OLD", "vendor_invoice_date": "2024-08-30", "vendor_invoice_total": 30}`),
			map[string]string{"Authorization": ownerToken, "Content-Type": "application/json"})
		if update.Code != http.StatusBadRequest || !strings.Contains(update.Body.String(), `"code":"locked"`) {
			t.Fatalf("expected the uploader to be locked out of a blank shared invoice, got %d %s", update.Code, update.Body.String())
		}

		// A bookkeeper reusing the invoice for another part may fill it in.
		status, body, _ = postVendorInvoiceExpense(t, app, bookkeeperToken,
			vendorInvoiceExpenseFields("third part", "10", details), content)
		if status != http.StatusOK {
			t.Fatalf("a bookkeeper should fill in a blank invoice: %d %s", status, body)
		}
		var number string
		if err := app.DB().NewQuery("SELECT ed.vendor_invoice_number FROM expenses e JOIN expense_documents ed ON ed.id = e.attachment_document WHERE e.id = {:id}").
			Bind(map[string]any{"id": ownerExpenseID}).Row(&number); err != nil {
			t.Fatal(err)
		}
		if number != "INV-OLD" {
			t.Fatalf("vendor_invoice_number = %q; want INV-OLD", number)
		}

		// Once filled in, the details are locked for bookkeepers too.
		status, body, _ = postVendorInvoiceExpense(t, app, bookkeeperToken,
			vendorInvoiceExpenseFields("fourth part", "0.01", invoice("INV-OLD", "2024-08-30", "40")), content)
		if status != http.StatusBadRequest || !strings.Contains(body, `"code":"locked"`) {
			t.Fatalf("expected filled-in details to be locked, got %d %s", status, body)
		}
	})

	t.Run("a shared invoice keeps its vendor", func(t *testing.T) {
		app := testutils.SetupTestApp(t)
		t.Cleanup(app.Cleanup)
		details := invoice("INV-VENDOR", "2024-08-30", "30")
		content := pngContent(9)
		status, body, _ := postVendorInvoiceExpense(t, app, bookkeeperToken,
			vendorInvoiceExpenseFields("first part", "10", details), content)
		if status != http.StatusOK {
			t.Fatalf("first part failed: %d %s", status, body)
		}
		// The bookkeeper reuses the invoice for an expense with another vendor.
		fields := vendorInvoiceExpenseFields("other vendor part", "10", details)
		fields["vendor"] = "2zqxtsmymf670ha"
		status, body, _ = postVendorInvoiceExpense(t, app, bookkeeperToken, fields, content)
		if status != http.StatusBadRequest || !strings.Contains(body, `"vendor":{"code":"vendor_mismatch"`) {
			t.Fatalf("expected a vendor mismatch, got %d %s", status, body)
		}
	})

	t.Run("a shared invoice keeps its currency", func(t *testing.T) {
		app := testutils.SetupTestApp(t)
		t.Cleanup(app.Cleanup)
		details := invoice("INV-CURRENCY", "2024-08-30", "30")
		content := pngContent(10)
		status, body, _ := postVendorInvoiceExpense(t, app, bookkeeperToken,
			vendorInvoiceExpenseFields("CAD part", "10", details), content)
		if status != http.StatusOK {
			t.Fatalf("CAD part failed: %d %s", status, body)
		}
		// The same file reused on the owner's USD PO 26-PRINT-004, keeping the
		// invoice's vendor so only the currency differs. That fixture PO has no
		// branch, so it borrows the CAD PO's for this test.
		if _, err := app.DB().NewQuery(`UPDATE purchase_orders SET branch = (SELECT branch FROM purchase_orders WHERE id = 'poa1ctvbrnch001') WHERE id = 'poprint00000004'`).Execute(); err != nil {
			t.Fatal(err)
		}
		fields := vendorInvoiceExpenseFields("USD part", "10", details)
		fields["purchase_order"] = "poprint00000004"
		fields["currency"] = "usdcurr00000001"
		fields["date"] = "2024-11-06" // the PO starts on 2024-11-05
		status, body, _ = postVendorInvoiceExpense(t, app, bookkeeperToken, fields, content)
		if status != http.StatusBadRequest || !strings.Contains(body, `"currency":{"code":"currency_mismatch"`) {
			t.Fatalf("expected a currency mismatch, got %d %s", status, body)
		}
	})

	t.Run("an invoice used by one expense follows that expense's vendor", func(t *testing.T) {
		app := testutils.SetupTestApp(t)
		t.Cleanup(app.Cleanup)
		status, body, id := postVendorInvoiceExpense(t, app, ownerToken,
			vendorInvoiceExpenseFields("owner entered invoice", "10", invoice("INV-FIX", "2024-08-30", "10")), pngContent(11))
		if status != http.StatusOK {
			t.Fatalf("create failed: %d %s", status, body)
		}
		update := performTestAPIRequest(t, app, http.MethodPatch, "/api/expenses/"+id,
			strings.NewReader(`{"vendor": "2zqxtsmymf670ha"}`), map[string]string{
				"Authorization": ownerToken,
				"Content-Type":  "application/json",
			})
		mustStatus(t, update, http.StatusOK)
		var vendor string
		if err := app.DB().NewQuery("SELECT ed.vendor FROM expenses e JOIN expense_documents ed ON ed.id = e.attachment_document WHERE e.id = {:id}").
			Bind(map[string]any{"id": id}).Row(&vendor); err != nil {
			t.Fatal(err)
		}
		if vendor != "2zqxtsmymf670ha" {
			t.Fatalf("document vendor = %q; want the corrected vendor", vendor)
		}
	})

	t.Run("details need an attached invoice", func(t *testing.T) {
		app := testutils.SetupTestApp(t)
		t.Cleanup(app.Cleanup)
		// A personal reimbursement has no attachment, so it has no document to
		// hold invoice details.
		payload := `{
			"uid": "rzr98oadsp9qc11",
			"date": "2024-09-01",
			"division": "vccd5fo56ctbigh",
			"description": "reimbursement with invoice details",
			"payment_type": "PersonalReimbursement",
			"kind": "l3vtlbqg529m52j",
			"total": 20,
			"vendor_invoice_number": "INV-NODOC",
			"vendor_invoice_date": "2024-08-30",
			"vendor_invoice_total": 20
		}`
		res := performTestAPIRequest(t, app, http.MethodPost, "/api/expenses", strings.NewReader(payload), map[string]string{
			"Authorization": ownerToken,
			"Content-Type":  "application/json",
		})
		if res.Code != http.StatusBadRequest || !strings.Contains(res.Body.String(), `"vendor_invoice_number":{"code":"document_required"`) {
			t.Fatalf("expected details without a document to be refused, got %d %s", res.Code, res.Body.String())
		}
	})

	t.Run("a corporate card expense needs only its receipt", func(t *testing.T) {
		app := testutils.SetupTestApp(t)
		t.Cleanup(app.Cleanup)
		status, body, _ := postVendorInvoiceExpense(t, app, ownerToken, map[string]string{
			"uid":              "rzr98oadsp9qc11",
			"date":             "2024-09-01",
			"division":         "vccd5fo56ctbigh",
			"description":      "card receipt",
			"payment_type":     "CorporateCreditCard",
			"cc_last_4_digits": "1234",
			"kind":             "l3vtlbqg529m52j",
			"total":            "20",
			"vendor":           "2zqxtsmymf670ha",
		}, pngContent(8))
		if status != http.StatusOK {
			t.Fatalf("a card expense without invoice details should save: %d %s", status, body)
		}
	})
}
