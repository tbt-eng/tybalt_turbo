package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"tybalt/internal/testutils"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

// Fixtures: time@test.com (rzr98oadsp9qc11) is active and gets the invoicing
// claim in these tests. fatt@mac.com holds report only. Job 24-321
// (cjf0kt0defhq480) belongs to Mobilia SA (lb0fnenkeyitsny). On it, expense
// su3hyft6n9rlt7d is committed, expense 31vvfz3z77n9628 is not, and PO
// 2plsetqdxht7esg (2024-0008) is Active. PO powip0000000001 is Active on a job
// of another client.

const (
	invoiceClient     = "lb0fnenkeyitsny"
	invoiceJob        = "cjf0kt0defhq480"
	committedExpense  = "su3hyft6n9rlt7d"
	uncommittedExpnse = "31vvfz3z77n9628"
	activePO          = "2plsetqdxht7esg"
	otherClientPO     = "powip0000000001"
)

func grantInvoicingClaim(t *testing.T, app *tests.TestApp, uid string) {
	t.Helper()
	claim, err := app.FindFirstRecordByData("claims", "name", "invoicing")
	if err != nil {
		t.Fatalf("invoicing claim missing: %v", err)
	}
	collection, err := app.FindCollectionByNameOrId("user_claims")
	if err != nil {
		t.Fatal(err)
	}
	record := core.NewRecord(collection)
	record.Set("uid", uid)
	record.Set("cid", claim.Id)
	if err := app.Save(record); err != nil {
		t.Fatal(err)
	}
}

func invoiceRequest(t *testing.T, app *tests.TestApp, token, method, path string, body any) (int, string) {
	t.Helper()
	var reader *strings.Reader
	if body == nil {
		reader = strings.NewReader("")
	} else {
		data, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
		reader = strings.NewReader(string(data))
	}
	res := performTestAPIRequest(t, app, method, path, reader, map[string]string{
		"Authorization": token,
		"Content-Type":  "application/json",
	})
	return res.Code, res.Body.String()
}

func createInvoice(t *testing.T, app *tests.TestApp, token string, fields map[string]any) string {
	t.Helper()
	status, body := invoiceRequest(t, app, token, http.MethodPost, "/api/collections/invoices/records", fields)
	if status != http.StatusOK {
		t.Fatalf("create invoice failed: %d %s", status, body)
	}
	var created struct {
		ID string `json:"id"`
	}
	_ = json.Unmarshal([]byte(body), &created)
	return created.ID
}

func issuedInvoice(t *testing.T, app *tests.TestApp, token, client, number string) string {
	t.Helper()
	id := createInvoice(t, app, token, map[string]any{"client": client, "number": number, "date": "2026-09-15", "total": 1200})
	if status, body := invoiceRequest(t, app, token, http.MethodPost, "/api/invoices/"+id+"/issue", nil); status != http.StatusOK {
		t.Fatalf("issue failed: %d %s", status, body)
	}
	return id
}

func TestClientInvoices(t *testing.T) {
	billerToken, err := testutils.GenerateRecordToken("users", "time@test.com")
	if err != nil {
		t.Fatal(err)
	}
	reportToken, err := testutils.GenerateRecordToken("users", "fatt@mac.com")
	if err != nil {
		t.Fatal(err)
	}
	setup := func(t *testing.T) *tests.TestApp {
		app := testutils.SetupTestApp(t)
		t.Cleanup(app.Cleanup)
		grantInvoicingClaim(t, app, "rzr98oadsp9qc11")
		return app
	}

	t.Run("only invoicing holders create invoices", func(t *testing.T) {
		app := setup(t)
		status, body := invoiceRequest(t, app, reportToken, http.MethodPost, "/api/collections/invoices/records",
			map[string]any{"client": invoiceClient})
		if status == http.StatusOK {
			t.Fatalf("report holder created an invoice: %s", body)
		}
	})

	t.Run("a new invoice is a draft owned by its creator", func(t *testing.T) {
		app := setup(t)
		id := createInvoice(t, app, billerToken, map[string]any{
			"client": invoiceClient, "status": "Issued", "issued_by": "rzr98oadsp9qc11",
		})
		invoice, err := app.FindRecordById("invoices", id)
		if err != nil {
			t.Fatal(err)
		}
		if invoice.GetString("status") != "Draft" || invoice.GetString("creator") != "rzr98oadsp9qc11" || invoice.GetString("issued_by") != "" {
			t.Fatalf("unexpected draft: status=%q creator=%q issued_by=%q",
				invoice.GetString("status"), invoice.GetString("creator"), invoice.GetString("issued_by"))
		}
	})

	t.Run("numbers are YYMMNNN with a real month", func(t *testing.T) {
		app := setup(t)
		for _, number := range []string{"261300", "2613001", "26090010", "INV-1"} {
			status, body := invoiceRequest(t, app, billerToken, http.MethodPost, "/api/collections/invoices/records",
				map[string]any{"client": invoiceClient, "number": number})
			if status != http.StatusBadRequest || !strings.Contains(body, `"invalid_number"`) {
				t.Fatalf("number %q: expected invalid_number, got %d %s", number, status, body)
			}
		}
	})

	t.Run("numbers are unique, and deleting a draft frees its number", func(t *testing.T) {
		app := setup(t)
		first := createInvoice(t, app, billerToken, map[string]any{"client": invoiceClient, "number": "2609001"})
		status, body := invoiceRequest(t, app, billerToken, http.MethodPost, "/api/collections/invoices/records",
			map[string]any{"client": invoiceClient, "number": "2609001"})
		if status != http.StatusBadRequest || !strings.Contains(body, `"number_taken"`) {
			t.Fatalf("expected number_taken, got %d %s", status, body)
		}
		if status, body := invoiceRequest(t, app, billerToken, http.MethodDelete, "/api/collections/invoices/records/"+first, nil); status != http.StatusNoContent {
			t.Fatalf("delete draft failed: %d %s", status, body)
		}
		createInvoice(t, app, billerToken, map[string]any{"client": invoiceClient, "number": "2609001"})
	})

	t.Run("issuing needs a number, a date and a total", func(t *testing.T) {
		app := setup(t)
		id := createInvoice(t, app, billerToken, map[string]any{"client": invoiceClient})
		status, body := invoiceRequest(t, app, billerToken, http.MethodPost, "/api/invoices/"+id+"/issue", nil)
		if status != http.StatusBadRequest || !strings.Contains(body, `"number":{"code":"required"`) {
			t.Fatalf("expected number required, got %d %s", status, body)
		}
		invoiceRequest(t, app, billerToken, http.MethodPatch, "/api/collections/invoices/records/"+id,
			map[string]any{"number": "2609002", "date": "2026-09-15"})
		status, body = invoiceRequest(t, app, billerToken, http.MethodPost, "/api/invoices/"+id+"/issue", nil)
		if status != http.StatusBadRequest || !strings.Contains(body, `"total":{"code":"required"`) {
			t.Fatalf("expected total required, got %d %s", status, body)
		}
	})

	t.Run("an issued invoice cannot be edited or deleted", func(t *testing.T) {
		app := setup(t)
		id := issuedInvoice(t, app, billerToken, invoiceClient, "2609003")
		invoice, err := app.FindRecordById("invoices", id)
		if err != nil {
			t.Fatal(err)
		}
		if invoice.GetString("status") != "Issued" || invoice.GetString("issued_by") != "rzr98oadsp9qc11" || invoice.GetString("issued_at") == "" {
			t.Fatalf("issue did not record who and when: %v", invoice.FieldsData())
		}
		if status, _ := invoiceRequest(t, app, billerToken, http.MethodPatch, "/api/collections/invoices/records/"+id,
			map[string]any{"total": 1}); status == http.StatusOK {
			t.Fatal("edited an issued invoice")
		}
		if status, _ := invoiceRequest(t, app, billerToken, http.MethodDelete, "/api/collections/invoices/records/"+id, nil); status == http.StatusNoContent {
			t.Fatal("deleted an issued invoice")
		}
		if status, _ := invoiceRequest(t, app, billerToken, http.MethodPost, "/api/invoices/"+id+"/issue", nil); status == http.StatusOK {
			t.Fatal("issued an invoice twice")
		}
	})

	t.Run("an edit cannot issue a draft", func(t *testing.T) {
		app := setup(t)
		id := createInvoice(t, app, billerToken, map[string]any{"client": invoiceClient})
		invoiceRequest(t, app, billerToken, http.MethodPatch, "/api/collections/invoices/records/"+id,
			map[string]any{"status": "Issued", "issued_by": "rzr98oadsp9qc11"})
		invoice, _ := app.FindRecordById("invoices", id)
		if invoice.GetString("status") != "Draft" || invoice.GetString("issued_by") != "" {
			t.Fatalf("edit changed status: %q issued_by %q", invoice.GetString("status"), invoice.GetString("issued_by"))
		}
	})

	addLine := func(t *testing.T, app *tests.TestApp, token string, fields map[string]any) (int, string) {
		t.Helper()
		return invoiceRequest(t, app, token, http.MethodPost, "/api/collections/invoice_lines/records", fields)
	}

	t.Run("lines bill committed expenses and active POs once", func(t *testing.T) {
		app := setup(t)
		id := issuedInvoice(t, app, billerToken, invoiceClient, "2609004")
		if status, body := addLine(t, app, billerToken, map[string]any{"invoice": id, "kind": "expense", "expense": committedExpense, "amount": 800}); status != http.StatusOK {
			t.Fatalf("expense line failed: %d %s", status, body)
		}
		if status, body := addLine(t, app, billerToken, map[string]any{"invoice": id, "kind": "po", "purchase_order": activePO, "amount": 50}); status != http.StatusOK {
			t.Fatalf("PO line failed: %d %s", status, body)
		}
		other := issuedInvoice(t, app, billerToken, invoiceClient, "2609005")
		status, body := addLine(t, app, billerToken, map[string]any{"invoice": other, "kind": "expense", "expense": committedExpense, "amount": 800})
		if status != http.StatusBadRequest || !strings.Contains(body, `"already_billed"`) || !strings.Contains(body, "2609004") {
			t.Fatalf("expected already_billed on 2609004, got %d %s", status, body)
		}

		status, body = invoiceRequest(t, app, reportToken, http.MethodGet, "/api/invoices/"+id+"/lines", nil)
		if status != http.StatusOK {
			t.Fatalf("lines failed: %d %s", status, body)
		}
		var lines []map[string]any
		if err := json.Unmarshal([]byte(body), &lines); err != nil {
			t.Fatal(err)
		}
		if len(lines) != 2 || lines[0]["job_number"] != "24-321" {
			t.Fatalf("unexpected lines: %s", body)
		}
	})

	t.Run("a line needs an issued invoice of the cost's client", func(t *testing.T) {
		app := setup(t)
		draft := createInvoice(t, app, billerToken, map[string]any{"client": invoiceClient, "number": "2609006"})
		status, body := addLine(t, app, billerToken, map[string]any{"invoice": draft, "kind": "expense", "expense": committedExpense, "amount": 10})
		if status != http.StatusBadRequest || !strings.Contains(body, `"invoice_not_issued"`) {
			t.Fatalf("expected invoice_not_issued, got %d %s", status, body)
		}
		id := issuedInvoice(t, app, billerToken, invoiceClient, "2609007")
		cases := []struct {
			fields map[string]any
			code   string
		}{
			{map[string]any{"invoice": id, "kind": "po", "purchase_order": otherClientPO, "amount": 10}, "client_mismatch"},
			{map[string]any{"invoice": id, "kind": "expense", "expense": uncommittedExpnse, "amount": 10}, "expense_not_committed"},
			{map[string]any{"invoice": id, "kind": "expense", "purchase_order": activePO, "amount": 10}, "required"},
			{map[string]any{"invoice": id, "kind": "expense", "expense": committedExpense, "purchase_order": activePO, "amount": 10}, "not_allowed"},
		}
		for _, c := range cases {
			status, body := addLine(t, app, billerToken, c.fields)
			if status != http.StatusBadRequest || !strings.Contains(body, `"`+c.code+`"`) {
				t.Fatalf("%v: expected %s, got %d %s", c.fields, c.code, status, body)
			}
		}
	})

	t.Run("report holders cannot add lines", func(t *testing.T) {
		app := setup(t)
		id := issuedInvoice(t, app, billerToken, invoiceClient, "2609008")
		if status, body := addLine(t, app, reportToken, map[string]any{"invoice": id, "kind": "po", "purchase_order": activePO, "amount": 10}); status == http.StatusOK {
			t.Fatalf("report holder added a line: %s", body)
		}
	})

	t.Run("billable lists the client's unbilled costs", func(t *testing.T) {
		app := setup(t)
		id := issuedInvoice(t, app, billerToken, invoiceClient, "2609009")
		status, body := invoiceRequest(t, app, billerToken, http.MethodGet, "/api/invoices/"+id+"/billable", nil)
		if status != http.StatusOK || !strings.Contains(body, `"id":"`+invoiceJob+`"`) {
			t.Fatalf("expected job %s in billable jobs, got %d %s", invoiceJob, status, body)
		}
		status, body = invoiceRequest(t, app, billerToken, http.MethodGet, "/api/invoices/"+id+"/billable?job="+invoiceJob, nil)
		if status != http.StatusOK || !strings.Contains(body, committedExpense) || !strings.Contains(body, activePO) || strings.Contains(body, uncommittedExpnse) {
			t.Fatalf("unexpected billable costs: %d %s", status, body)
		}
		addLine(t, app, billerToken, map[string]any{"invoice": id, "kind": "po", "purchase_order": activePO, "amount": 10})
		_, body = invoiceRequest(t, app, billerToken, http.MethodGet, "/api/invoices/"+id+"/billable?job="+invoiceJob, nil)
		if strings.Contains(body, activePO) {
			t.Fatalf("billed PO still listed: %s", body)
		}
	})

	t.Run("an expense's details show the invoice it was billed on", func(t *testing.T) {
		app := setup(t)
		// The biller cannot open someone else's committed expense; a report
		// holder can, but cannot bill it.
		if status, body := invoiceRequest(t, app, billerToken, http.MethodGet, "/api/expenses/client_invoice/"+committedExpense, nil); status != http.StatusNotFound {
			t.Fatalf("expected the biller to be refused, got %d %s", status, body)
		}
		status, body := invoiceRequest(t, app, reportToken, http.MethodGet, "/api/expenses/client_invoice/"+committedExpense, nil)
		if status != http.StatusOK || !strings.Contains(body, `"line":null`) || !strings.Contains(body, `"can_bill":false`) {
			t.Fatalf("expected an unbilled expense the report holder cannot bill, got %d %s", status, body)
		}
		id := issuedInvoice(t, app, billerToken, invoiceClient, "2609011")
		addLine(t, app, billerToken, map[string]any{"invoice": id, "kind": "expense", "expense": committedExpense, "amount": 900})
		status, body = invoiceRequest(t, app, reportToken, http.MethodGet, "/api/expenses/client_invoice/"+committedExpense, nil)
		if status != http.StatusOK || !strings.Contains(body, `"invoice_number":"2609011"`) || !strings.Contains(body, `"can_open_invoice":true`) {
			t.Fatalf("expected the billing invoice, got %d %s", status, body)
		}
	})

	t.Run("a PO's details show the invoice it was billed on", func(t *testing.T) {
		app := setup(t)
		status, body := invoiceRequest(t, app, billerToken, http.MethodGet, "/api/purchase_orders/client_invoice/"+activePO, nil)
		if status != http.StatusOK || !strings.Contains(body, `"line":null`) || !strings.Contains(body, `"can_bill":true`) {
			t.Fatalf("expected an unbilled, billable PO, got %d %s", status, body)
		}
		id := issuedInvoice(t, app, billerToken, invoiceClient, "2609010")
		addLine(t, app, billerToken, map[string]any{"invoice": id, "kind": "po", "purchase_order": activePO, "amount": 10})
		status, body = invoiceRequest(t, app, reportToken, http.MethodGet, "/api/purchase_orders/client_invoice/"+activePO, nil)
		if status != http.StatusOK || !strings.Contains(body, `"invoice_number":"2609010"`) || !strings.Contains(body, `"can_bill":false`) {
			t.Fatalf("expected the billing invoice, got %d %s", status, body)
		}
	})
}
