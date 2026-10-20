package main

import (
	"crypto/sha256"
	"encoding/hex"
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

// TestExpenseVendorInvoiceView checks the list of expenses on an invoice: every
// part is listed, only the parts the viewer may open are linkable, and the
// person who set the details is shown.
func TestExpenseVendorInvoiceView(t *testing.T) {
	ownerToken, err := testutils.GenerateRecordToken("users", "time@test.com")
	if err != nil {
		t.Fatal(err)
	}
	bookkeeperToken, err := testutils.GenerateRecordToken("users", "book@keeper.com")
	if err != nil {
		t.Fatal(err)
	}
	strangerToken, err := testutils.GenerateRecordToken("users", "francesco@mac.com")
	if err != nil {
		t.Fatal(err)
	}
	app := testutils.SetupTestApp(t)
	t.Cleanup(app.Cleanup)

	details := map[string]string{"vendor_invoice_number": "INV-VIEW", "vendor_invoice_date": "2024-08-30", "vendor_invoice_total": "30"}
	content := pngContent(20)
	status, body, ownerExpenseID := postVendorInvoiceExpense(t, app, ownerToken,
		vendorInvoiceExpenseFields("owner part", "10", details), content)
	if status != http.StatusOK {
		t.Fatalf("owner part failed: %d %s", status, body)
	}
	status, body, bookkeeperExpenseID := postVendorInvoiceExpense(t, app, bookkeeperToken,
		vendorInvoiceExpenseFields("bookkeeper part", "12.5", details), content)
	if status != http.StatusOK {
		t.Fatalf("bookkeeper part failed: %d %s", status, body)
	}

	type part struct {
		ID        string  `json:"id"`
		PONumber  string  `json:"po_number"`
		Total     float64 `json:"total"`
		EnteredBy string  `json:"entered_by_name"`
		CanOpen   bool    `json:"can_open"`
	}
	type view struct {
		Number     string  `json:"vendor_invoice_number"`
		Total      float64 `json:"vendor_invoice_total"`
		RecordedBy string  `json:"recorded_by_name"`
		Recorded   string  `json:"vendor_invoice_recorded"`
		Allocated  float64 `json:"allocated"`
		Parts      []part  `json:"parts"`
	}
	load := func(token, expenseID string) (int, view) {
		t.Helper()
		res := performTestAPIRequest(t, app, http.MethodGet, "/api/expenses/invoice/"+expenseID, nil, map[string]string{"Authorization": token})
		var v view
		_ = json.Unmarshal(res.Body.Bytes(), &v)
		return res.Code, v
	}

	// The bookkeeper reused the invoice without changing it, so the owner who
	// set the details is still shown as having recorded them.
	code, v := load(bookkeeperToken, bookkeeperExpenseID)
	if code != http.StatusOK {
		t.Fatalf("bookkeeper view failed: %d", code)
	}
	if v.Number != "INV-VIEW" || v.Total != 30 || v.Allocated != 22.5 || v.RecordedBy != "Tester Time" || v.Recorded == "" {
		t.Fatalf("unexpected invoice view: %+v", v)
	}
	if len(v.Parts) != 2 {
		t.Fatalf("parts = %d; want 2", len(v.Parts))
	}
	canOpen := map[string]bool{}
	for _, p := range v.Parts {
		if p.PONumber != "2026-0901" {
			t.Errorf("part %s PO = %q; want 2026-0901", p.ID, p.PONumber)
		}
		canOpen[p.ID] = p.CanOpen
	}
	// The bookkeeper entered their part but may not open the owner's own part.
	if !canOpen[bookkeeperExpenseID] || canOpen[ownerExpenseID] {
		t.Fatalf("bookkeeper can_open = %v; want only their own part", canOpen)
	}

	// The owner can open both: the bookkeeper's part was entered for them.
	code, v = load(ownerToken, ownerExpenseID)
	if code != http.StatusOK || len(v.Parts) != 2 || !v.Parts[0].CanOpen || !v.Parts[1].CanOpen {
		t.Fatalf("owner view = %d %+v; want both parts openable", code, v)
	}

	// Someone who cannot see the expense cannot see its invoice.
	if code, _ := load(strangerToken, ownerExpenseID); code != http.StatusNotFound {
		t.Fatalf("stranger view = %d; want 404", code)
	}

	// Changing the details records who changed them.
	update := performTestAPIRequest(t, app, http.MethodPatch, "/api/expenses/"+bookkeeperExpenseID,
		strings.NewReader(`{"vendor_invoice_total": 40}`), map[string]string{
			"Authorization": bookkeeperToken,
			"Content-Type":  "application/json",
		})
	mustStatus(t, update, http.StatusOK)
	if _, v = load(ownerToken, ownerExpenseID); v.RecordedBy != "Ultra Chifres" || v.Total != 40 {
		t.Fatalf("after the bookkeeper's change: %+v; want recorded by Ultra Chifres, total 40", v)
	}
}

// TestVendorInvoiceNonUploaderCannotChangeSharedInvoice follows the one way a
// user who neither uploaded an invoice nor is a bookkeeper can end up on a
// shared invoice: the uploader deletes their draft, leaving the document
// orphaned; another user uploads the same file, which relinks that document
// with its original uploader; then a bookkeeper reuses it for a PO. The second
// user may not change the shared invoice's details.
func TestVendorInvoiceNonUploaderCannotChangeSharedInvoice(t *testing.T) {
	ownerToken, err := testutils.GenerateRecordToken("users", "time@test.com")
	if err != nil {
		t.Fatal(err)
	}
	otherToken, err := testutils.GenerateRecordToken("users", "author@soup.com")
	if err != nil {
		t.Fatal(err)
	}
	bookkeeperToken, err := testutils.GenerateRecordToken("users", "book@keeper.com")
	if err != nil {
		t.Fatal(err)
	}
	app := testutils.SetupTestApp(t)
	t.Cleanup(app.Cleanup)
	content := pngContent(30)
	details := map[string]string{"vendor_invoice_number": "INV-ORPHAN", "vendor_invoice_date": "2024-08-30", "vendor_invoice_total": "50"}

	// 1. The PO owner uploads the invoice, then deletes the draft.
	status, body, ownerExpenseID := postVendorInvoiceExpense(t, app, ownerToken,
		vendorInvoiceExpenseFields("owner draft", "10", details), content)
	if status != http.StatusOK {
		t.Fatalf("owner draft failed: %d %s", status, body)
	}
	del := performTestAPIRequest(t, app, http.MethodDelete, "/api/collections/expenses/records/"+ownerExpenseID, nil, map[string]string{"Authorization": ownerToken})
	mustStatus(t, del, http.StatusNoContent)

	// 2. Another user uploads the same file on their own expense. The orphaned
	// document is relinked and keeps its original uploader.
	otherFields := map[string]string{
		"uid":          "f2j5a8vk006baub",
		"date":         "2024-09-01",
		"division":     "vccd5fo56ctbigh",
		"description":  "other user's on-account expense",
		"payment_type": "OnAccount",
		"kind":         "l3vtlbqg529m52j",
		"total":        "10",
		"vendor":       "z66xe6vqhwtokt4",
	}
	for key, value := range details {
		otherFields[key] = value
	}
	status, body, otherExpenseID := postVendorInvoiceExpense(t, app, otherToken, otherFields, content)
	if status != http.StatusOK {
		t.Fatalf("other user's expense failed: %d %s", status, body)
	}
	var uploadedBy string
	if err := app.DB().NewQuery("SELECT ed.uploaded_by FROM expenses e JOIN expense_documents ed ON ed.id = e.attachment_document WHERE e.id = {:id}").
		Bind(map[string]any{"id": otherExpenseID}).Row(&uploadedBy); err != nil {
		t.Fatal(err)
	}
	if uploadedBy != "rzr98oadsp9qc11" {
		t.Fatalf("uploaded_by = %q; the relinked document should keep its original uploader", uploadedBy)
	}

	// 3. A bookkeeper reuses the invoice for the PO owner's PO.
	status, body, _ = postVendorInvoiceExpense(t, app, bookkeeperToken,
		vendorInvoiceExpenseFields("bookkeeper part", "10", details), content)
	if status != http.StatusOK {
		t.Fatalf("bookkeeper reuse failed: %d %s", status, body)
	}

	// The second user is neither the uploader nor a bookkeeper.
	update := performTestAPIRequest(t, app, http.MethodPatch, "/api/expenses/"+otherExpenseID,
		strings.NewReader(`{"vendor_invoice_total": 60}`), map[string]string{
			"Authorization": otherToken,
			"Content-Type":  "application/json",
		})
	if update.Code != http.StatusForbidden || !strings.Contains(update.Body.String(), `"code":"not_editable"`) {
		t.Fatalf("expected not_editable, got %d %s", update.Code, update.Body.String())
	}
	// Saving their expense without changing the invoice still works.
	update = performTestAPIRequest(t, app, http.MethodPatch, "/api/expenses/"+otherExpenseID,
		strings.NewReader(`{"description": "other user's corrected description"}`), map[string]string{
			"Authorization": otherToken,
			"Content-Type":  "application/json",
		})
	mustStatus(t, update, http.StatusOK)
}

// TestVendorInvoiceLookupByFile checks that a bookkeeper choosing a file that is
// already on an expense can see the invoice it will be reused for, and that
// nobody else can look files up.
func TestVendorInvoiceLookupByFile(t *testing.T) {
	ownerToken, err := testutils.GenerateRecordToken("users", "time@test.com")
	if err != nil {
		t.Fatal(err)
	}
	bookkeeperToken, err := testutils.GenerateRecordToken("users", "book@keeper.com")
	if err != nil {
		t.Fatal(err)
	}
	app := testutils.SetupTestApp(t)
	t.Cleanup(app.Cleanup)
	content := pngContent(31)
	status, body, _ := postVendorInvoiceExpense(t, app, ownerToken,
		vendorInvoiceExpenseFields("owner part", "10", map[string]string{
			"vendor_invoice_number": "INV-FILE", "vendor_invoice_date": "2024-08-30", "vendor_invoice_total": "25",
		}), content)
	if status != http.StatusOK {
		t.Fatalf("owner part failed: %d %s", status, body)
	}
	sum := sha256.Sum256(content)
	url := "/api/expenses/invoice-by-file/" + hex.EncodeToString(sum[:])

	res := performTestAPIRequest(t, app, http.MethodGet, url, nil, map[string]string{"Authorization": bookkeeperToken})
	mustStatus(t, res, http.StatusOK)
	var view struct {
		Document string `json:"document"`
		Number   string `json:"vendor_invoice_number"`
		Parts    []struct {
			ID string `json:"id"`
		} `json:"parts"`
	}
	if err := json.Unmarshal(res.Body.Bytes(), &view); err != nil {
		t.Fatal(err)
	}
	if view.Document == "" || view.Number != "INV-FILE" || len(view.Parts) != 1 {
		t.Fatalf("unexpected lookup: %+v", view)
	}

	// A file nobody has uploaded returns no invoice.
	unknown := sha256.Sum256([]byte("not uploaded"))
	res = performTestAPIRequest(t, app, http.MethodGet, "/api/expenses/invoice-by-file/"+hex.EncodeToString(unknown[:]), nil, map[string]string{"Authorization": bookkeeperToken})
	mustStatus(t, res, http.StatusOK)
	if strings.Contains(res.Body.String(), `"document":"`) && !strings.Contains(res.Body.String(), `"document":""`) {
		t.Fatalf("expected no invoice for an unknown file, got %s", res.Body.String())
	}

	// Only bookkeepers can reuse a file, so only they can look one up.
	res = performTestAPIRequest(t, app, http.MethodGet, url, nil, map[string]string{"Authorization": ownerToken})
	mustStatus(t, res, http.StatusForbidden)
}
