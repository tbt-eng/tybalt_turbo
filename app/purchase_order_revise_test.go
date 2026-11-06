package main

import (
	"encoding/json"
	"math"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"tybalt/internal/testutils"
	"tybalt/utilities"

	"github.com/pocketbase/pocketbase/tests"
)

// PO 2024-0008 is an Active One-Time PO with no expenses. time@test.com
// created it and author@soup.com approved it.
const (
	revisePOID     = "2plsetqdxht7esg"
	revisePONumber = "2024-0008"
	reviseCreator  = "rzr98oadsp9qc11"
	reviseApprover = "f2j5a8vk006baub"
)

func reviseTestTokens(t *testing.T) (creator, approver, payablesAdmin string) {
	t.Helper()
	var err error
	if creator, err = testutils.GenerateRecordToken("users", "time@test.com"); err != nil {
		t.Fatal(err)
	}
	if approver, err = testutils.GenerateRecordToken("users", "author@soup.com"); err != nil {
		t.Fatal(err)
	}
	if payablesAdmin, err = testutils.GenerateRecordToken("users", "book@keeper.com"); err != nil {
		t.Fatal(err)
	}
	return creator, approver, payablesAdmin
}

func reviseTestRequest(t *testing.T, app *tests.TestApp, method, path, body, token string, wantStatus int) string {
	t.Helper()
	response := performTestAPIRequest(t, app, method, path, strings.NewReader(body),
		map[string]string{"Authorization": token})
	if response.Code != wantStatus {
		t.Fatalf("%s %s: status=%d, want %d; body=%s", method, path, response.Code, wantStatus, response.Body.String())
	}
	return response.Body.String()
}

func reviseTestPOField(t *testing.T, app *tests.TestApp, id, field string) string {
	t.Helper()
	po, err := app.FindRecordById("purchase_orders", id)
	if err != nil {
		t.Fatal(err)
	}
	return po.GetString(field)
}

func TestRevisePurchaseOrder(t *testing.T) {
	creator, approver, payablesAdmin := reviseTestTokens(t)
	revisePath := "/api/purchase_orders/" + revisePOID + "/revise"

	t.Run("revision returns the PO to Unapproved and approval restores the same number", func(t *testing.T) {
		app := testutils.SetupTestApp(t)
		t.Cleanup(app.Cleanup)

		body := reviseTestRequest(t, app, http.MethodPost, revisePath, `{"total": 150}`, creator, http.StatusOK)
		for _, want := range []string{`"status":"Unapproved"`, `"po_number":"` + revisePONumber + `"`, `"approved":""`, `"total":150`, `"approval_total":150`} {
			if !strings.Contains(body, want) {
				t.Fatalf("response missing %s: %s", want, body)
			}
		}

		var notified struct{ Count int }
		if err := app.DB().NewQuery(`
			SELECT COUNT(*) AS count FROM notifications n
			JOIN notification_templates t ON t.id = n.template
			WHERE t.code = 'po_revision_approval_required' AND n.status = 'sent' AND n.recipient = {:recipient}
			  AND json_extract(n.data, '$.POId') = {:po}
		`).Bind(map[string]any{"recipient": reviseApprover, "po": revisePOID}).One(&notified); err != nil {
			t.Fatal(err)
		}
		if notified.Count != 1 {
			t.Fatalf("approver notifications = %d, want 1", notified.Count)
		}
		email := app.TestMailer.LastMessage()
		for _, want := range []string{"revised purchase order", revisePONumber, "150.00 CAD", "Previous total: 132.10 CAD", "Stuff Unlimited"} {
			if !strings.Contains(email.Text, want) {
				t.Fatalf("approval email is missing %q: %s", want, email.Text)
			}
		}
		if strings.Contains(email.Text, "has created") {
			t.Fatalf("approval email says the purchase order was created: %s", email.Text)
		}

		reviseTestRequest(t, app, http.MethodPost, "/api/purchase_orders/"+revisePOID+"/approve", "", approver, http.StatusOK)
		if got := reviseTestPOField(t, app, revisePOID, "status"); got != "Active" {
			t.Fatalf("status after approval = %q, want Active", got)
		}
		if got := reviseTestPOField(t, app, revisePOID, "po_number"); got != revisePONumber {
			t.Fatalf("po_number after approval = %q, want %q", got, revisePONumber)
		}
	})

	t.Run("a creator who can approve the revision approves it themselves", func(t *testing.T) {
		app := testutils.SetupTestApp(t)
		t.Cleanup(app.Cleanup)

		// The form sends the creator as the approver when they qualify. PO 2025-0007-01
		// was created and approved by author@soup.com.
		const ownPOID, ownPONumber = "25046ft47x49cc2", "2025-0007-01"
		reviseTestRequest(t, app, http.MethodPost, "/api/purchase_orders/"+ownPOID+"/revise",
			`{"total": 150, "approver": "`+reviseApprover+`"}`, approver, http.StatusOK)
		if got := reviseTestPOField(t, app, ownPOID, "status"); got != "Unapproved" {
			t.Fatalf("status after revising = %q, want Unapproved until the creator approves it", got)
		}

		pending := reviseTestRequest(t, app, http.MethodGet, "/api/purchase_orders/pending", "", approver, http.StatusOK)
		if !strings.Contains(pending, ownPOID) {
			t.Fatalf("the revision is not in the creator's pending queue: %s", pending)
		}

		reviseTestRequest(t, app, http.MethodPost, "/api/purchase_orders/"+ownPOID+"/approve", "", approver, http.StatusOK)
		if got := reviseTestPOField(t, app, ownPOID, "status"); got != "Active" {
			t.Fatalf("status after approval = %q, want Active", got)
		}
		if got := reviseTestPOField(t, app, ownPOID, "po_number"); got != ownPONumber {
			t.Fatalf("po_number after approval = %q, want %q", got, ownPONumber)
		}
	})

	t.Run("revision is refused and leaves the PO unchanged", func(t *testing.T) {
		for _, tc := range []struct {
			name     string
			poID     string
			body     string
			asAuthor bool
			setup    string
			status   int
			code     string
		}{
			{name: "caller is not the creator", poID: revisePOID, body: `{"total": 150}`, asAuthor: true, status: http.StatusForbidden, code: "unauthorized_revision"},
			{name: "PO is not Active", poID: "338568325487lo2", body: `{"total": 150}`, status: http.StatusBadRequest, code: "po_not_active"},
			{name: "PO has expenses", poID: "ly8xyzpuj79upq1", body: `{"total": 5000}`, status: http.StatusBadRequest, code: "po_has_expenses"},
			{name: "PO has child POs", poID: "xhkt5lx8cl64nj3", body: `{"total": 5000}`, asAuthor: true, setup: `DELETE FROM expenses WHERE purchase_order = 'xhkt5lx8cl64nj3'`, status: http.StatusBadRequest, code: "po_has_child_pos"},
			{name: "PO is a legacy entry", poID: "po_export_no_exp_1", body: `{"total": 200}`, asAuthor: true, status: http.StatusBadRequest, code: "legacy_po_not_revisable"},
			{name: "total is unchanged", poID: revisePOID, body: `{"total": 132.1}`, status: http.StatusBadRequest, code: "total_unchanged"},
			{name: "total and currency are unchanged", poID: revisePOID, body: `{"total": 132.1, "currency": "cadcurr00000001"}`, status: http.StatusBadRequest, code: "total_unchanged"},
			{name: "currency does not exist", poID: revisePOID, body: `{"total": 150, "currency": "nosuchcurrency0"}`, status: http.StatusBadRequest, code: `"currency":{"code":"not_found"`},
			{name: "currency has no exchange rate", poID: revisePOID, body: `{"total": 150, "currency": "wipcurr00000001"}`, status: http.StatusBadRequest, code: `"currency":{`},
			{name: "child PO leaves its parent's currency", poID: "25046ft47x49cc2", body: `{"total": 140, "currency": "usdcurr00000001"}`, asAuthor: true, status: http.StatusBadRequest, code: `"currency":{"code":"value_mismatch"`},
			{name: "total is not positive", poID: revisePOID, body: `{"total": 0}`, status: http.StatusBadRequest, code: `"total":{"code":"required"`},
			{name: "approver is not valid for the new total", poID: revisePOID, body: `{"total": 150, "approver": "tqqf7q0f3378rvp"}`, status: http.StatusBadRequest, code: "invalid_approver_for_stage"},
		} {
			t.Run(tc.name, func(t *testing.T) {
				app := testutils.SetupTestApp(t)
				t.Cleanup(app.Cleanup)
				if tc.setup != "" {
					if _, err := app.DB().NewQuery(tc.setup).Execute(); err != nil {
						t.Fatal(err)
					}
				}
				before := reviseTestPOField(t, app, tc.poID, "status")
				beforeCurrency := reviseTestPOField(t, app, tc.poID, "currency")
				token := creator
				if tc.asAuthor {
					token = approver
				}
				body := reviseTestRequest(t, app, http.MethodPost, "/api/purchase_orders/"+tc.poID+"/revise", tc.body, token, tc.status)
				if !strings.Contains(body, tc.code) {
					t.Fatalf("response missing %s: %s", tc.code, body)
				}
				if got := reviseTestPOField(t, app, tc.poID, "status"); got != before {
					t.Fatalf("status = %q, want unchanged %q", got, before)
				}
				if got := reviseTestPOField(t, app, tc.poID, "currency"); got != beforeCurrency {
					t.Fatalf("currency = %q, want unchanged %q", got, beforeCurrency)
				}
			})
		}
	})

	t.Run("a revision can change only the currency and keeps it after approval", func(t *testing.T) {
		app := testutils.SetupTestApp(t)
		t.Cleanup(app.Cleanup)

		reviseTestRequest(t, app, http.MethodPost, revisePath, `{"total": 132.1, "currency": "usdcurr00000001"}`, creator, http.StatusOK)
		po, err := app.FindRecordById("purchase_orders", revisePOID)
		if err != nil {
			t.Fatal(err)
		}
		if po.GetString("currency") != "usdcurr00000001" || po.GetString("status") != "Unapproved" {
			t.Fatalf("currency = %q, status = %q", po.GetString("currency"), po.GetString("status"))
		}
		// The seeded rate is 1.35 CAD per USD.
		if got := po.GetFloat("approval_total_home"); math.Abs(got-132.1*1.35) > 0.000001 {
			t.Fatalf("approval_total_home = %v, want %v", got, 132.1*1.35)
		}
		email := app.TestMailer.LastMessage()
		for _, want := range []string{"Revised total: 132.10 USD", "Previous total: 132.10 CAD"} {
			if !strings.Contains(email.Text, want) {
				t.Fatalf("approval email is missing %q: %s", want, email.Text)
			}
		}

		reviseTestRequest(t, app, http.MethodPost, "/api/purchase_orders/"+revisePOID+"/approve", `{}`, approver, http.StatusOK)
		if got := reviseTestPOField(t, app, revisePOID, "status"); got != "Active" {
			t.Fatalf("status = %q, want Active", got)
		}
		if got := reviseTestPOField(t, app, revisePOID, "currency"); got != "usdcurr00000001" {
			t.Fatalf("currency = %q after approval", got)
		}
		if got := reviseTestPOField(t, app, revisePOID, "po_number"); got != revisePONumber {
			t.Fatalf("po_number = %q, want %q", got, revisePONumber)
		}
	})

	t.Run("a revision with no currency in the request keeps its currency", func(t *testing.T) {
		app := testutils.SetupTestApp(t)
		t.Cleanup(app.Cleanup)
		reviseTestRequest(t, app, http.MethodPost, revisePath, `{"total": 150, "currency": "usdcurr00000001"}`, creator, http.StatusOK)
		reviseTestRequest(t, app, http.MethodPost, "/api/purchase_orders/"+revisePOID+"/approve", `{}`, approver, http.StatusOK)
		reviseTestRequest(t, app, http.MethodPost, revisePath, `{"total": 160}`, creator, http.StatusOK)
		if got := reviseTestPOField(t, app, revisePOID, "currency"); got != "usdcurr00000001" {
			t.Fatalf("currency = %q, want it kept", got)
		}
	})

	t.Run("a child PO can be revised and keeps its child number", func(t *testing.T) {
		app := testutils.SetupTestApp(t)
		t.Cleanup(app.Cleanup)
		// PO 2025-0007-01 is the Active child of 2025-0007, owned by author@soup.com.
		body := reviseTestRequest(t, app, http.MethodPost, "/api/purchase_orders/25046ft47x49cc2/revise", `{"total": 140}`, approver, http.StatusOK)
		if !strings.Contains(body, `"status":"Unapproved"`) || !strings.Contains(body, `"po_number":"2025-0007-01"`) {
			t.Fatalf("unexpected child revision response: %s", body)
		}
	})

	t.Run("a revised PO only accepts changes to the total, currency and approvers", func(t *testing.T) {
		app := testutils.SetupTestApp(t)
		t.Cleanup(app.Cleanup)
		reviseTestRequest(t, app, http.MethodPost, revisePath, `{"total": 150}`, creator, http.StatusOK)

		recordPath := "/api/collections/purchase_orders/records/" + revisePOID
		body := reviseTestRequest(t, app, http.MethodPatch, recordPath, `{"description": "a different description"}`, creator, http.StatusBadRequest)
		if !strings.Contains(body, `"description":{"code":"revision_field_locked"`) {
			t.Fatalf("expected revision_field_locked on description: %s", body)
		}

		reviseTestRequest(t, app, http.MethodPatch, recordPath, `{"total": 175}`, creator, http.StatusOK)
		if got := reviseTestPOField(t, app, revisePOID, "total"); got != "175" {
			t.Fatalf("total = %q, want 175", got)
		}

		reviseTestRequest(t, app, http.MethodPatch, recordPath, `{"currency": "usdcurr00000001"}`, creator, http.StatusOK)
		if got := reviseTestPOField(t, app, revisePOID, "currency"); got != "usdcurr00000001" {
			t.Fatalf("currency = %q, want usdcurr00000001", got)
		}
	})

	t.Run("a revised PO cannot be deleted", func(t *testing.T) {
		app := testutils.SetupTestApp(t)
		t.Cleanup(app.Cleanup)
		reviseTestRequest(t, app, http.MethodPost, revisePath, `{"total": 150}`, creator, http.StatusOK)

		reviseTestRequest(t, app, http.MethodDelete, "/api/collections/purchase_orders/records/"+revisePOID, "", creator, http.StatusNotFound)
		if got := reviseTestPOField(t, app, revisePOID, "po_number"); got != revisePONumber {
			t.Fatalf("po_number = %q, want %q", got, revisePONumber)
		}
	})

	t.Run("only the creator can cancel a revision, and only once it is rejected", func(t *testing.T) {
		app := testutils.SetupTestApp(t)
		t.Cleanup(app.Cleanup)
		cancelPath := "/api/purchase_orders/" + revisePOID + "/cancel"
		reviseTestRequest(t, app, http.MethodPost, revisePath, `{"total": 150}`, creator, http.StatusOK)

		body := reviseTestRequest(t, app, http.MethodPost, cancelPath, "", creator, http.StatusBadRequest)
		if !strings.Contains(body, "po_not_active") {
			t.Fatalf("pending revision cancel: %s", body)
		}

		reviseTestRequest(t, app, http.MethodPost, "/api/purchase_orders/"+revisePOID+"/reject",
			`{"rejection_reason": "too expensive"}`, approver, http.StatusOK)

		body = reviseTestRequest(t, app, http.MethodPost, cancelPath, "", payablesAdmin, http.StatusForbidden)
		if !strings.Contains(body, "unauthorized_cancellation") {
			t.Fatalf("payables_admin cancel of rejected revision: %s", body)
		}

		reviseTestRequest(t, app, http.MethodPost, cancelPath, "", creator, http.StatusOK)
		if got := reviseTestPOField(t, app, revisePOID, "status"); got != "Cancelled" {
			t.Fatalf("status = %q, want Cancelled", got)
		}
		if got := reviseTestPOField(t, app, revisePOID, "canceller"); got != reviseCreator {
			t.Fatalf("canceller = %q, want %q", got, reviseCreator)
		}
	})

	t.Run("a revision stays visible to everyone but only its approver can act on it", func(t *testing.T) {
		app := testutils.SetupTestApp(t)
		t.Cleanup(app.Cleanup)
		other, err := testutils.GenerateRecordToken("users", "noclaims@example.com")
		if err != nil {
			t.Fatal(err)
		}
		visiblePath := "/api/purchase_orders/visible/" + revisePOID
		recordPath := "/api/collections/purchase_orders/records/" + revisePOID
		jobPath := "/api/jobs/" + reviseTestPOField(t, app, revisePOID, "job") + "/pos/list"
		reviseTestRequest(t, app, http.MethodPost, revisePath, `{"total": 150}`, creator, http.StatusOK)

		for _, path := range []string{
			visiblePath,
			recordPath,
			"/api/purchase_orders/visible?scope=active",
			"/api/purchase_orders/search",
			jobPath,
		} {
			body := reviseTestRequest(t, app, http.MethodGet, path, "", other, http.StatusOK)
			if !strings.Contains(body, revisePONumber) || !strings.Contains(body, `"Unapproved"`) {
				t.Fatalf("GET %s does not show the revision: %s", path, body)
			}
		}

		reviseTestRequest(t, app, http.MethodGet, "/api/purchase_orders/pending/"+revisePOID, "", other, http.StatusNotFound)
		reviseTestRequest(t, app, http.MethodPost, "/api/purchase_orders/"+revisePOID+"/approve", "", other, http.StatusForbidden)
		reviseTestRequest(t, app, http.MethodPost, "/api/purchase_orders/"+revisePOID+"/reject",
			`{"rejection_reason": "not my purchase order"}`, other, http.StatusForbidden)
		if got := reviseTestPOField(t, app, revisePOID, "status"); got != "Unapproved" {
			t.Fatalf("status = %q, want Unapproved", got)
		}
	})

	t.Run("only the people involved can read why a revision was rejected", func(t *testing.T) {
		app := testutils.SetupTestApp(t)
		t.Cleanup(app.Cleanup)
		other, err := testutils.GenerateRecordToken("users", "noclaims@example.com")
		if err != nil {
			t.Fatal(err)
		}
		reason := "far too expensive"
		visiblePath := "/api/purchase_orders/visible/" + revisePOID
		recordPath := "/api/collections/purchase_orders/records/" + revisePOID
		reviseTestRequest(t, app, http.MethodPost, revisePath, `{"total": 150}`, creator, http.StatusOK)
		reviseTestRequest(t, app, http.MethodPost, "/api/purchase_orders/"+revisePOID+"/reject",
			`{"rejection_reason": "`+reason+`"}`, approver, http.StatusOK)

		// Someone uninvolved still finds it through the visible route, without the
		// reason. The collection API can't hide a field from filters, so it leaves
		// the record out for them, which also stops them probing the reason.
		body := reviseTestRequest(t, app, http.MethodGet, visiblePath, "", other, http.StatusOK)
		if strings.Contains(body, reason) || !strings.Contains(body, revisePONumber) {
			t.Fatalf("visible route for an unrelated user: %s", body)
		}
		reviseTestRequest(t, app, http.MethodGet, recordPath, "", other, http.StatusNotFound)
		filter := url.QueryEscape(`rejection_reason ~ "expensive"`)
		body = reviseTestRequest(t, app, http.MethodGet,
			"/api/collections/purchase_orders/records?filter="+filter, "", other, http.StatusOK)
		if strings.Contains(body, revisePOID) {
			t.Fatalf("filter on the rejection reason found the revision: %s", body)
		}

		for _, path := range []string{visiblePath, recordPath} {
			for _, token := range []string{creator, approver} {
				if body := reviseTestRequest(t, app, http.MethodGet, path, "", token, http.StatusOK); !strings.Contains(body, reason) {
					t.Fatalf("GET %s hides the rejection reason from someone involved: %s", path, body)
				}
			}
		}
	})

	t.Run("the vendor page filter finds a pending revision", func(t *testing.T) {
		app := testutils.SetupTestApp(t)
		t.Cleanup(app.Cleanup)
		other, err := testutils.GenerateRecordToken("users", "noclaims@example.com")
		if err != nil {
			t.Fatal(err)
		}
		vendorID := reviseTestPOField(t, app, revisePOID, "vendor")
		reviseTestRequest(t, app, http.MethodPost, revisePath, `{"total": 150}`, creator, http.StatusOK)
		filter := url.QueryEscape(`vendor='` + vendorID + `' && (status='Unapproved' && po_number!='' && legacy_manual_entry=false)`)
		body := reviseTestRequest(t, app, http.MethodGet,
			"/api/collections/purchase_orders/records?perPage=100&filter="+filter, "", other, http.StatusOK)
		if !strings.Contains(body, revisePOID) {
			t.Fatalf("vendor filter did not find the revision: %s", body)
		}
	})

	t.Run("the job POs tab leaves a revision out of its total but counts it", func(t *testing.T) {
		app := testutils.SetupTestApp(t)
		t.Cleanup(app.Cleanup)
		other, err := testutils.GenerateRecordToken("users", "noclaims@example.com")
		if err != nil {
			t.Fatal(err)
		}
		summaryPath := "/api/jobs/" + reviseTestPOField(t, app, revisePOID, "job") + "/pos/summary"
		var before, after struct {
			TotalAmount   float64 `json:"total_amount"`
			RevisionCount int     `json:"revision_count"`
		}
		body := reviseTestRequest(t, app, http.MethodGet, summaryPath, "", other, http.StatusOK)
		if err := json.Unmarshal([]byte(body), &before); err != nil {
			t.Fatal(err)
		}
		reviseTestRequest(t, app, http.MethodPost, revisePath, `{"total": 150}`, creator, http.StatusOK)
		body = reviseTestRequest(t, app, http.MethodGet, summaryPath, "", other, http.StatusOK)
		if err := json.Unmarshal([]byte(body), &after); err != nil {
			t.Fatal(err)
		}
		if after.RevisionCount != before.RevisionCount+1 {
			t.Fatalf("revision_count = %d, want %d", after.RevisionCount, before.RevisionCount+1)
		}
		if diff := before.TotalAmount - after.TotalAmount; diff < 132.09 || diff > 132.11 {
			t.Fatalf("total fell by %v, want the PO's approved 132.1", diff)
		}
	})

	t.Run("revisable says whether the revise route can accept a PO", func(t *testing.T) {
		app := testutils.SetupTestApp(t)
		t.Cleanup(app.Cleanup)
		for id, want := range map[string]bool{
			revisePOID:           true,
			"ly8xyzpuj79upq1":    false, // has expenses
			"po_export_no_exp_1": false, // legacy
		} {
			body := reviseTestRequest(t, app, http.MethodGet, "/api/purchase_orders/visible/"+id, "", payablesAdmin, http.StatusOK)
			var row struct {
				Revisable bool `json:"revisable"`
			}
			if err := json.Unmarshal([]byte(body), &row); err != nil {
				t.Fatal(err)
			}
			if row.Revisable != want {
				t.Errorf("%s revisable = %v, want %v", id, row.Revisable, want)
			}
		}
	})

	t.Run("a revision with an open child PO is not revisable", func(t *testing.T) {
		app := testutils.SetupTestApp(t)
		t.Cleanup(app.Cleanup)
		// Deleting the fixture's expenses leaves its child PO as the only reason.
		// This edits fixture data at test time, as the refusal table above does.
		expenses, err := app.FindRecordsByFilter("expenses", "purchase_order = 'xhkt5lx8cl64nj3'", "", 0, 0)
		if err != nil {
			t.Fatal(err)
		}
		for _, expense := range expenses {
			if err := app.Delete(expense); err != nil {
				t.Fatal(err)
			}
		}
		body := reviseTestRequest(t, app, http.MethodGet, "/api/purchase_orders/visible/xhkt5lx8cl64nj3", "", payablesAdmin, http.StatusOK)
		if !strings.Contains(body, `"revisable":false`) {
			t.Fatalf("PO with a child is revisable: %s", body)
		}
	})

	t.Run("computed and import fields can't be set on a revision", func(t *testing.T) {
		app := testutils.SetupTestApp(t)
		t.Cleanup(app.Cleanup)
		recordPath := "/api/collections/purchase_orders/records/" + revisePOID
		reviseTestRequest(t, app, http.MethodPost, revisePath, `{"total": 150}`, creator, http.StatusOK)
		for _, patch := range []string{`{"_imported": true}`, `{"attachment_hash": "abc"}`} {
			body := reviseTestRequest(t, app, http.MethodPatch, recordPath, patch, creator, http.StatusBadRequest)
			if !strings.Contains(body, "revision_field_locked") {
				t.Fatalf("PATCH %s: %s", patch, body)
			}
		}
	})

	t.Run("a revision succeeds when its notification can't be created", func(t *testing.T) {
		app := testutils.SetupTestApp(t)
		t.Cleanup(app.Cleanup)
		// Removing the template at test time is how this test makes dispatch fail.
		template, err := app.FindFirstRecordByData("notification_templates", "code", "po_revision_approval_required")
		if err != nil {
			t.Fatal(err)
		}
		if err := app.Delete(template); err != nil {
			t.Fatal(err)
		}
		reviseTestRequest(t, app, http.MethodPost, revisePath, `{"total": 150}`, creator, http.StatusOK)
		if got := reviseTestPOField(t, app, revisePOID, "status"); got != "Unapproved" {
			t.Fatalf("status = %q, want Unapproved", got)
		}
	})

	t.Run("the approver email reports what a recurring PO needs approved", func(t *testing.T) {
		app := testutils.SetupTestApp(t)
		t.Cleanup(app.Cleanup)
		const recurringPO = "rec5e5la2fa4rpn"
		reviseTestRequest(t, app, http.MethodPost, "/api/purchase_orders/"+recurringPO+"/revise", `{"total": 1}`, creator, http.StatusOK)
		po, err := app.FindRecordById("purchase_orders", recurringPO)
		if err != nil {
			t.Fatal(err)
		}
		if po.GetFloat("approval_total") <= po.GetFloat("total") {
			t.Fatalf("approval_total = %v, want more than one occurrence of %v", po.GetFloat("approval_total"), po.GetFloat("total"))
		}
		want := utilities.FormatCurrencyAmount(po.GetFloat("approval_total"), "CAD")
		if email := app.TestMailer.LastMessage(); !strings.Contains(email.Text, want) {
			t.Fatalf("email does not show %s: %s", want, email.Text)
		}
	})

	t.Run("the approver is told again when a revision changes for them", func(t *testing.T) {
		app := testutils.SetupTestApp(t)
		t.Cleanup(app.Cleanup)
		recordPath := "/api/collections/purchase_orders/records/" + revisePOID
		sentEmails := func(code string) int {
			t.Helper()
			var result struct{ Count int }
			if err := app.DB().NewQuery(`
				SELECT COUNT(*) AS count FROM notifications n
				JOIN notification_templates t ON t.id = n.template
				WHERE t.code = {:code} AND json_extract(n.data, '$.POId') = {:po}
			`).Bind(map[string]any{"code": code, "po": revisePOID}).One(&result); err != nil {
				t.Fatal(err)
			}
			return result.Count
		}
		expectEmails := func(step string, want int) {
			t.Helper()
			if got := sentEmails("po_revision_approval_required"); got != want {
				t.Fatalf("%s: revision emails = %d, want %d", step, got, want)
			}
			if got := sentEmails("po_approval_required"); got != 0 {
				t.Fatalf("%s: %d emails say the PO was created", step, got)
			}
		}

		reviseTestRequest(t, app, http.MethodPost, revisePath, `{"total": 150}`, creator, http.StatusOK)
		expectEmails("revised", 1)

		reviseTestRequest(t, app, http.MethodPatch, recordPath, `{"total": 150}`, creator, http.StatusOK)
		expectEmails("saved with no change", 1)

		reviseTestRequest(t, app, http.MethodPatch, recordPath, `{"total": 160}`, creator, http.StatusOK)
		expectEmails("new total", 2)

		reviseTestRequest(t, app, http.MethodPatch, recordPath, `{"currency": "cadcurr00000001"}`, creator, http.StatusOK)
		expectEmails("saved with the same currency", 2)

		reviseTestRequest(t, app, http.MethodPatch, recordPath, `{"currency": "usdcurr00000001"}`, creator, http.StatusOK)
		expectEmails("new currency", 3)
		email := app.TestMailer.LastMessage()
		for _, want := range []string{"Revised total: 160.00 USD", "Previous total: 160.00 CAD"} {
			if !strings.Contains(email.Text, want) {
				t.Fatalf("currency change email is missing %q: %s", want, email.Text)
			}
		}

		reviseTestRequest(t, app, http.MethodPost, "/api/purchase_orders/"+revisePOID+"/reject",
			`{"rejection_reason": "not this total"}`, approver, http.StatusOK)
		reviseTestRequest(t, app, http.MethodPatch, recordPath, `{"total": 160}`, creator, http.StatusOK)
		expectEmails("resubmitted after rejection", 4)
		if email := app.TestMailer.LastMessage(); strings.Contains(email.Text, "Previous total") {
			t.Fatalf("resubmitting unchanged reports a previous total: %s", email.Text)
		}

		// Setting the first approval directly stands in for a stage 1 approval,
		// because this fixture's total doesn't need a second approval.
		po, err := app.FindRecordById("purchase_orders", revisePOID)
		if err != nil {
			t.Fatal(err)
		}
		po.Set("approved", "2026-11-05 12:00:00.000Z")
		if err := app.Save(po); err != nil {
			t.Fatal(err)
		}
		reviseTestRequest(t, app, http.MethodPatch, recordPath, `{"total": 170}`, creator, http.StatusOK)
		expectEmails("approval reset", 5)
		if got := reviseTestPOField(t, app, revisePOID, "approved"); got != "" {
			t.Fatalf("approved = %q after the reset, want it cleared", got)
		}
		if email := app.TestMailer.LastMessage(); strings.Contains(email.Text, "has created") {
			t.Fatalf("reset email says the PO was created: %s", email.Text)
		}
	})

	t.Run("resubmitting a revision rejected at second approval starts over", func(t *testing.T) {
		app := testutils.SetupTestApp(t)
		t.Cleanup(app.Cleanup)
		recordPath := "/api/collections/purchase_orders/records/" + revisePOID
		reviseTestRequest(t, app, http.MethodPost, revisePath, `{"total": 150}`, creator, http.StatusOK)
		// Setting the first approval directly stands in for a stage 1 approval, as
		// in the test above.
		po, err := app.FindRecordById("purchase_orders", revisePOID)
		if err != nil {
			t.Fatal(err)
		}
		po.Set("approved", "2026-11-05 12:00:00.000Z")
		if err := app.Save(po); err != nil {
			t.Fatal(err)
		}
		reviseTestRequest(t, app, http.MethodPost, "/api/purchase_orders/"+revisePOID+"/reject",
			`{"rejection_reason": "not this total"}`, approver, http.StatusOK)
		reviseTestRequest(t, app, http.MethodPatch, recordPath, `{"total": 150}`, creator, http.StatusOK)
		if got := reviseTestPOField(t, app, revisePOID, "approved"); got != "" {
			t.Fatalf("approved = %q after resubmitting, want it cleared", got)
		}
		if got := reviseTestPOField(t, app, revisePOID, "rejected"); got != "" {
			t.Fatalf("rejected = %q after resubmitting, want it cleared", got)
		}
	})
}
