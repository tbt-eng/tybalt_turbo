package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"
	"tybalt/internal/testutils"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

// Fixture IDs used by the purchase order group tests. The super_po holder is
// superpo@test.com. The approver, fatt@mac.com, can approve capital purchase
// orders up to 500 in the CI division, which is also the capital second
// approval threshold.
const (
	poGroupCapitalKind = "l3vtlbqg529m52j"
	poGroupDivision    = "vccd5fo56ctbigh"
	poGroupVendor      = "2zqxtsmymf670ha"
	poGroupApprover    = "etysnrlup2f6bak"
	poGroupCAD         = "cadcurr00000001"
	poGroupTimeUID     = "rzr98oadsp9qc11"
	poGroupAuthorUID   = "f2j5a8vk006baub"
	poGroupHalUID      = "66ct66w380ob6w8"
	poGroupCreatorUID  = "superpouser0001"
	poGroupThunderBay  = "80875lm27v8wgi4"
	poGroupCorporate   = "kpj5jijh0if8kx8"
)

func poGroupToken(t testing.TB, email string) string {
	t.Helper()
	token, err := testutils.GenerateRecordToken("users", email)
	if err != nil {
		t.Fatal(err)
	}
	return token
}

func poGroupBody(t testing.TB, groupID string, uploadFields []string, rows []map[string]string, overrides map[string]any) string {
	t.Helper()
	template := map[string]any{
		"type":         "Cumulative",
		"kind":         poGroupCapitalKind,
		"payment_type": "OnAccount",
		"date":         "2026-11-01",
		"division":     poGroupDivision,
		"branch":       poGroupThunderBay,
		"description":  "Workwear allowance",
		"total":        100,
		"currency":     poGroupCAD,
		"vendor":       poGroupVendor,
		"approver":     poGroupApprover,
	}
	for key, value := range overrides {
		template[key] = value
	}
	body, err := json.Marshal(map[string]any{
		"group_id":      groupID,
		"template":      template,
		"upload_fields": uploadFields,
		"rows":          rows,
	})
	if err != nil {
		t.Fatal(err)
	}
	return string(body)
}

func poGroupValidRows() []map[string]string {
	return []map[string]string{
		{"email": "time@test.com", "total": "150"},
		{"email": "Author@Soup.com", "total": "$225.50"},
		{"email": "hal@2005.com", "total": "80"},
	}
}

type poGroupPreviewResponse struct {
	Rows []struct {
		Index        int               `json:"index"`
		Email        string            `json:"email"`
		UID          string            `json:"uid"`
		Name         string            `json:"name"`
		Total        float64           `json:"total"`
		DivisionCode string            `json:"division_code"`
		BranchName   string            `json:"branch_name"`
		Errors       map[string]string `json:"errors"`
	} `json:"rows"`
	Count              int     `json:"count"`
	ErrorCount         int     `json:"error_count"`
	Total              float64 `json:"total"`
	CurrencyCode       string  `json:"currency_code"`
	MissingActiveStaff []struct {
		UID string `json:"uid"`
	} `json:"missing_active_staff"`
}

func decodePOGroupPreview(t testing.TB, body []byte) poGroupPreviewResponse {
	t.Helper()
	var preview poGroupPreviewResponse
	if err := json.Unmarshal(body, &preview); err != nil {
		t.Fatalf("failed to decode preview: %v; body=%s", err, body)
	}
	return preview
}

func countPOGroupRows(t testing.TB, app *tests.TestApp, query string, params dbx.Params) int {
	t.Helper()
	var result struct {
		Count int `db:"count"`
	}
	if err := app.DB().NewQuery(query).Bind(params).One(&result); err != nil {
		t.Fatalf("count query failed: %v", err)
	}
	return result.Count
}

func TestPurchaseOrderGroupPreview(t *testing.T) {
	superPOToken := poGroupToken(t, "superpo@test.com")
	timeToken := poGroupToken(t, "time@test.com")

	scenarios := []tests.ApiScenario{
		{
			Name:           "preview requires the super_po claim",
			Method:         http.MethodPost,
			URL:            "/api/purchase_orders/groups/preview",
			Body:           strings.NewReader(poGroupBody(t, "", []string{"total"}, poGroupValidRows(), nil)),
			Headers:        map[string]string{"Authorization": timeToken},
			ExpectedStatus: http.StatusForbidden,
			ExpectedContent: []string{
				`"code":"super_po_claim_required"`,
			},
			TestAppFactory: testutils.SetupTestApp,
		},
		{
			Name:           "preview reports template errors before looking at rows",
			Method:         http.MethodPost,
			URL:            "/api/purchase_orders/groups/preview",
			Body:           strings.NewReader(poGroupBody(t, "", []string{"total", "job"}, poGroupValidRows(), map[string]any{"approver": "", "type": "Recurring", "branch": ""})),
			Headers:        map[string]string{"Authorization": superPOToken},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedContent: []string{
				`"type":{"code":"invalid_type"`,
				`"upload_fields":{"code":"invalid_upload_field"`,
				`"branch":{"code":"required"`,
			},
			// The approver is chosen from the preview, so only creating needs one.
			NotExpectedContent: []string{
				`"approver":`,
			},
			TestAppFactory: testutils.SetupTestApp,
		},
		{
			// Every row can be approved by fatt@mac.com, and none needs a
			// second approver. Without an approver no row is faulted for one.
			Name:           "preview offers the approvers who can approve every row",
			Method:         http.MethodPost,
			URL:            "/api/purchase_orders/groups/preview",
			Body:           strings.NewReader(poGroupBody(t, "", []string{"total"}, poGroupValidRows(), map[string]any{"approver": ""})),
			Headers:        map[string]string{"Authorization": superPOToken},
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"error_count":0`,
				`"approvers":[`,
				`"id":"etysnrlup2f6bak"`,
				`"second_approvers":[]`,
				`"second_approval_required":false`,
			},
			TestAppFactory: testutils.SetupTestApp,
		},
		{
			// 600 is over the capital second approval threshold, so a second
			// approver is needed, and fatt@mac.com's 500 limit makes them a
			// first-stage approver only.
			Name:   "preview offers second approvers when a row needs one",
			Method: http.MethodPost,
			URL:    "/api/purchase_orders/groups/preview",
			Body: strings.NewReader(poGroupBody(t, "", []string{"total"},
				[]map[string]string{{"email": "time@test.com", "total": "600"}, {"email": "hal@2005.com", "total": "80"}},
				map[string]any{"approver": ""})),
			Headers:        map[string]string{"Authorization": superPOToken},
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"error_count":0`,
				`"second_approval_required":true`,
				`"id":"6bq4j0eb26631dy"`,
			},
			TestAppFactory: testutils.SetupTestApp,
		},
		{
			// book@keeper.com has no default branch. The creator's branch is used
			// as given rather than replaced by the owner's default.
			Name:   "preview uses the creator's branch even when the owner has no default branch",
			Method: http.MethodPost,
			URL:    "/api/purchase_orders/groups/preview",
			Body: strings.NewReader(poGroupBody(t, "", []string{"total"},
				[]map[string]string{{"email": "book@keeper.com", "total": "150"}}, nil)),
			Headers:        map[string]string{"Authorization": superPOToken},
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"error_count":0`,
				`"branch_name":"Thunder Bay"`,
			},
			TestAppFactory: testutils.SetupTestApp,
		},
		{
			// corp.claim@example.com may use the Corporate branch, but the
			// creator, who sets the branch, may not.
			Name:   "preview checks the branch against the creator's claims, not the owner's",
			Method: http.MethodPost,
			URL:    "/api/purchase_orders/groups/preview",
			Body: strings.NewReader(poGroupBody(t, "", []string{"total"},
				[]map[string]string{{"email": "corp.claim@example.com", "total": "150"}}, map[string]any{"branch": poGroupCorporate})),
			Headers:        map[string]string{"Authorization": superPOToken},
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"error_count":1`,
				`"branch":"Corporate requires at least one of these claims`,
			},
			TestAppFactory: testutils.SetupTestApp,
		},
		{
			Name:           "preview resolves each row and lists active staff missing from the upload",
			Method:         http.MethodPost,
			URL:            "/api/purchase_orders/groups/preview",
			Body:           strings.NewReader(poGroupBody(t, "", []string{"total"}, poGroupValidRows(), nil)),
			Headers:        map[string]string{"Authorization": superPOToken},
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"count":3`,
				`"error_count":0`,
				`"total":455.5`,
				`"currency_code":"CAD"`,
				`"uid":"rzr98oadsp9qc11"`,
				`"name":"Tester Time"`,
				`"email":"author@soup.com"`,
				`"total":225.5`,
				`"division_code":"CI"`,
				`"branch_name":"Thunder Bay"`,
			},
			NotExpectedContent: []string{
				`"errors"`,
			},
			TestAppFactory: testutils.SetupTestApp,
		},
		{
			Name:   "preview reports row errors without stopping at the first",
			Method: http.MethodPost,
			URL:    "/api/purchase_orders/groups/preview",
			Body: strings.NewReader(poGroupBody(t, "", []string{"total", "division", "branch"}, []map[string]string{
				{"email": "time@test.com", "total": "150", "division": "CI", "branch": "Thunder Bay"},
				{"email": "nobody@example.com", "total": "150", "division": "CI"},
				{"email": "author@soup.com", "total": "abc", "division": "CI"},
				{"email": "hal@2005.com", "total": "80", "division": "ZZ"},
				{"email": "u_no_claims@example.com", "total": "80", "division": "CI", "branch": "Atlantis"},
				{"email": "orphan@poapprover.com", "total": "600", "division": "CI", "branch": "Thunder Bay"},
				{"email": "self_apv_yes@test.com", "total": "90", "division": "CI"},
				{"email": "SELF_APV_YES@test.com", "total": "90", "division": "CI"},
				{"email": "fatt@mac.com", "total": "", "division": "CI", "branch": "Thunder Bay"},
			}, nil)),
			Headers:        map[string]string{"Authorization": superPOToken},
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"count":9`,
				`"error_count":8`,
				`"Total is required."`,
				`"No staff member has this email."`,
				`"Total \"abc\" is not a number."`,
				`"No division has code \"ZZ\"."`,
				`"No branch has the name or code \"Atlantis\"."`,
				`"This amount needs a second approver. Choose one in the template."`,
				`"This email appears 2 times in the upload."`,
				`"Branch is required."`,
			},
			TestAppFactory: testutils.SetupTestApp,
		},
		{
			Name:   "preview requires a description in each row when it comes from the upload",
			Method: http.MethodPost,
			URL:    "/api/purchase_orders/groups/preview",
			Body: strings.NewReader(poGroupBody(t, "", []string{"total", "description"}, []map[string]string{
				{"email": "time@test.com", "total": "150", "description": "Boots and gloves"},
				{"email": "hal@2005.com", "total": "80", "description": " "},
			}, map[string]any{"description": ""})),
			Headers:        map[string]string{"Authorization": superPOToken},
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"error_count":1`,
				`"description":"Description is required."`,
			},
			TestAppFactory: testutils.SetupTestApp,
		},
		{
			Name:           "preview flags amounts the approver cannot approve",
			Method:         http.MethodPost,
			URL:            "/api/purchase_orders/groups/preview",
			Body:           strings.NewReader(poGroupBody(t, "", []string{"total"}, []map[string]string{{"email": "time@test.com", "total": "150"}}, map[string]any{"division": "2n7htnnf4bokbnh"})),
			Headers:        map[string]string{"Authorization": superPOToken},
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"error_count":1`,
				`"approver":"The approver can't approve this amount for this kind and division."`,
			},
			TestAppFactory: testutils.SetupTestApp,
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestPurchaseOrderGroupPreviewMissingStaff(t *testing.T) {
	app := testutils.SetupTestApp(t)
	t.Cleanup(app.Cleanup)

	res := performTestAPIRequest(t, app, http.MethodPost, "/api/purchase_orders/groups/preview",
		strings.NewReader(poGroupBody(t, "", []string{"total"}, poGroupValidRows(), nil)),
		map[string]string{"Authorization": poGroupToken(t, "superpo@test.com")})
	mustStatus(t, res, http.StatusOK)
	preview := decodePOGroupPreview(t, res.Body.Bytes())

	activeCount := countPOGroupRows(t, app, `
		SELECT COUNT(*) AS count FROM users u JOIN admin_profiles ap ON ap.uid = u.id WHERE ap.active = 1
	`, dbx.Params{})
	if len(preview.MissingActiveStaff) != activeCount-3 {
		t.Fatalf("expected %d active staff missing from the upload, got %d", activeCount-3, len(preview.MissingActiveStaff))
	}
	for _, staff := range preview.MissingActiveStaff {
		if staff.UID == poGroupTimeUID || staff.UID == poGroupAuthorUID || staff.UID == poGroupHalUID {
			t.Fatalf("staff member %s is in the upload but listed as missing", staff.UID)
		}
	}
}

func TestPurchaseOrderGroupCreate(t *testing.T) {
	app := testutils.SetupTestApp(t)
	t.Cleanup(app.Cleanup)
	superPOHeaders := map[string]string{"Authorization": poGroupToken(t, "superpo@test.com")}
	const groupID = "grouptest000001"

	notificationsBefore := countPOGroupRows(t, app, `SELECT COUNT(*) AS count FROM notifications`, dbx.Params{})

	res := performTestAPIRequest(t, app, http.MethodPost, "/api/purchase_orders/groups",
		strings.NewReader(poGroupBody(t, groupID, []string{"total"}, poGroupValidRows(), nil)), superPOHeaders)
	mustStatus(t, res, http.StatusOK)
	if body := mustReadBody(t, res); !strings.Contains(body, `"count":3`) || !strings.Contains(body, `"group_id":"grouptest000001"`) {
		t.Fatalf("unexpected create response: %s", body)
	}

	type createdPO struct {
		UID         string  `db:"uid"`
		Total       float64 `db:"total"`
		Status      string  `db:"status"`
		PONumber    string  `db:"po_number"`
		Approver    string  `db:"approver"`
		Branch      string  `db:"branch"`
		Type        string  `db:"type"`
		Description string  `db:"description"`
	}
	created := []createdPO{}
	if err := app.DB().NewQuery(`
		SELECT uid, total, status, po_number, approver, branch, type, description
		FROM purchase_orders WHERE po_group = {:group} ORDER BY total
	`).Bind(dbx.Params{"group": groupID}).All(&created); err != nil {
		t.Fatal(err)
	}
	want := []struct {
		uid   string
		total float64
	}{{poGroupHalUID, 80}, {poGroupTimeUID, 150}, {poGroupAuthorUID, 225.5}}
	if len(created) != len(want) {
		t.Fatalf("expected %d purchase orders in the group, got %d", len(want), len(created))
	}
	for i, po := range created {
		if po.UID != want[i].uid || po.Total != want[i].total {
			t.Fatalf("row %d: got uid %s total %v, want uid %s total %v", i, po.UID, po.Total, want[i].uid, want[i].total)
		}
		if po.Status != "Unapproved" || po.PONumber != "" || po.Approver != poGroupApprover {
			t.Fatalf("row %d should be Unapproved with no number and the template approver: %+v", i, po)
		}
		if po.Branch != poGroupThunderBay || po.Type != "Cumulative" || po.Description != "Workwear allowance" {
			t.Fatalf("row %d did not take the template: %+v", i, po)
		}
	}

	if creator := countPOGroupRows(t, app, `SELECT COUNT(*) AS count FROM purchase_order_groups WHERE id = {:id} AND creator = {:creator}`,
		dbx.Params{"id": groupID, "creator": poGroupCreatorUID}); creator != 1 {
		t.Fatalf("expected the group to record its creator")
	}

	// The approver gets one email for the whole group, not one per purchase
	// order.
	if added := countPOGroupRows(t, app, `SELECT COUNT(*) AS count FROM notifications`, dbx.Params{}) - notificationsBefore; added != 1 {
		t.Fatalf("expected 1 approval notification for the group, got %d", added)
	}
	if count := countPOGroupRows(t, app, `
		SELECT COUNT(*) AS count
		FROM notifications n JOIN notification_templates nt ON nt.id = n.template
		WHERE nt.code = 'po_group_approval_required' AND n.recipient = {:approver} AND n.status = 'sent'
	`, dbx.Params{"approver": poGroupApprover}); count != 1 {
		t.Fatalf("expected one sent po_group_approval_required notification to the approver")
	}
	messages := app.TestMailer.Messages()
	if len(messages) != 1 || messages[0].Subject != "Approval needed: 3 purchase orders" ||
		!strings.Contains(messages[0].Text, "Total of all 3: 455.50 CAD") ||
		!strings.Contains(messages[0].Text, "Vendor: Big Vendor Industries") {
		t.Fatalf("unexpected group approval email: %+v", messages)
	}

	// Repeating the request creates nothing.
	res = performTestAPIRequest(t, app, http.MethodPost, "/api/purchase_orders/groups",
		strings.NewReader(poGroupBody(t, groupID, []string{"total"}, poGroupValidRows(), nil)), superPOHeaders)
	mustStatus(t, res, http.StatusConflict)
	if !strings.Contains(mustReadBody(t, res), `"code":"group_exists"`) {
		t.Fatalf("expected group_exists, got %s", mustReadBody(t, res))
	}
	if count := countPOGroupRows(t, app, `SELECT COUNT(*) AS count FROM purchase_orders WHERE po_group = {:group}`,
		dbx.Params{"group": groupID}); count != 3 {
		t.Fatalf("repeat request changed the group size to %d", count)
	}

	// A super_po holder can list the group, including purchase orders they could
	// not otherwise see.
	res = performTestAPIRequest(t, app, http.MethodGet, "/api/purchase_orders/visible?scope=group&group="+groupID, nil, superPOHeaders)
	mustStatus(t, res, http.StatusOK)
	var visible []map[string]any
	if err := json.Unmarshal(res.Body.Bytes(), &visible); err != nil {
		t.Fatal(err)
	}
	if len(visible) != 3 {
		t.Fatalf("expected the creator to see 3 group purchase orders, got %d", len(visible))
	}
	if visible[0]["po_group"] != groupID {
		t.Fatalf("expected po_group in visible rows, got %v", visible[0]["po_group"])
	}

	res = performTestAPIRequest(t, app, http.MethodGet, "/api/purchase_orders/groups/"+groupID, nil, superPOHeaders)
	mustStatus(t, res, http.StatusOK)
	if body := mustReadBody(t, res); !strings.Contains(body, `"creator_name":"Group Buyer"`) ||
		!strings.Contains(body, `"count":3`) || !strings.Contains(body, `"unapproved_count":3`) || !strings.Contains(body, `"total":455.5`) ||
		!strings.Contains(body, `"description":"Workwear allowance"`) || !strings.Contains(body, `"vendor_name":"Big Vendor Industries"`) {
		t.Fatalf("unexpected group summary: %s", body)
	}
	res = performTestAPIRequest(t, app, http.MethodGet, "/api/purchase_orders/groups", nil, superPOHeaders)
	mustStatus(t, res, http.StatusOK)
	if body := mustReadBody(t, res); !strings.Contains(body, `"id":"grouptest000001"`) {
		t.Fatalf("expected the group in the list: %s", body)
	}

	// An admin holder can look at groups and their purchase orders but can't
	// create one.
	adminHeaders := map[string]string{"Authorization": poGroupToken(t, "admin.only@example.com")}
	res = performTestAPIRequest(t, app, http.MethodGet, "/api/purchase_orders/groups", nil, adminHeaders)
	mustStatus(t, res, http.StatusOK)
	res = performTestAPIRequest(t, app, http.MethodGet, "/api/purchase_orders/visible?scope=group&group="+groupID, nil, adminHeaders)
	mustStatus(t, res, http.StatusOK)
	visible = nil
	if err := json.Unmarshal(res.Body.Bytes(), &visible); err != nil {
		t.Fatal(err)
	}
	if len(visible) != 3 {
		t.Fatalf("expected an admin to see 3 group purchase orders, got %d", len(visible))
	}
	res = performTestAPIRequest(t, app, http.MethodPost, "/api/purchase_orders/groups/preview",
		strings.NewReader(poGroupBody(t, "", []string{"total"}, poGroupValidRows(), nil)), adminHeaders)
	mustStatus(t, res, http.StatusForbidden)

	// Each owner sees only their own purchase order in the group.
	timeHeaders := map[string]string{"Authorization": poGroupToken(t, "time@test.com")}
	res = performTestAPIRequest(t, app, http.MethodGet, "/api/purchase_orders/visible?scope=group&group="+groupID, nil, timeHeaders)
	mustStatus(t, res, http.StatusOK)
	visible = nil
	if err := json.Unmarshal(res.Body.Bytes(), &visible); err != nil {
		t.Fatal(err)
	}
	if len(visible) != 1 || visible[0]["uid"] != poGroupTimeUID {
		t.Fatalf("expected time@test.com to see only their own group purchase order, got %v", visible)
	}
	res = performTestAPIRequest(t, app, http.MethodGet, "/api/purchase_orders/groups", nil, timeHeaders)
	mustStatus(t, res, http.StatusForbidden)

	// The collection API agrees with the visibility routes. A super_po holder
	// can read every purchase order in the group, which also sends the group
	// page realtime events, while another owner can't. The groups collection
	// is read only through the group routes.
	var authorPOID string
	if err := app.DB().NewQuery(`SELECT id FROM purchase_orders WHERE po_group = {:group} AND uid = {:uid}`).
		Bind(dbx.Params{"group": groupID, "uid": poGroupAuthorUID}).Row(&authorPOID); err != nil {
		t.Fatal(err)
	}
	res = performTestAPIRequest(t, app, http.MethodGet, "/api/collections/purchase_orders/records/"+authorPOID, nil, superPOHeaders)
	mustStatus(t, res, http.StatusOK)
	res = performTestAPIRequest(t, app, http.MethodGet, "/api/collections/purchase_orders/records/"+authorPOID, nil, timeHeaders)
	mustStatus(t, res, http.StatusNotFound)
	res = performTestAPIRequest(t, app, http.MethodGet, "/api/collections/purchase_order_groups/records", nil, superPOHeaders)
	mustStatus(t, res, http.StatusForbidden)

	// The approver approves one purchase order at a time, which numbers it.
	var timePOID string
	if err := app.DB().NewQuery(`SELECT id FROM purchase_orders WHERE po_group = {:group} AND uid = {:uid}`).
		Bind(dbx.Params{"group": groupID, "uid": poGroupTimeUID}).Row(&timePOID); err != nil {
		t.Fatal(err)
	}
	res = performTestAPIRequest(t, app, http.MethodPost, "/api/purchase_orders/"+timePOID+"/approve", strings.NewReader(`{}`),
		map[string]string{"Authorization": poGroupToken(t, "fatt@mac.com")})
	mustStatus(t, res, http.StatusOK)
	if body := mustReadBody(t, res); !strings.Contains(body, `"status":"Active"`) || strings.Contains(body, `"po_number":""`) {
		t.Fatalf("expected an Active, numbered purchase order, got %s", body)
	}

	// The owner can still edit a group purchase order, but can't take it out of
	// the group.
	authorHeaders := map[string]string{"Authorization": poGroupToken(t, "author@soup.com")}
	res = performTestAPIRequest(t, app, http.MethodPatch, "/api/collections/purchase_orders/records/"+authorPOID,
		strings.NewReader(`{"po_group":""}`), authorHeaders)
	mustStatus(t, res, http.StatusBadRequest)
	if !strings.Contains(mustReadBody(t, res), `"po_group":{"code":"immutable_field"`) {
		t.Fatalf("expected po_group immutable_field, got %s", mustReadBody(t, res))
	}
	res = performTestAPIRequest(t, app, http.MethodPatch, "/api/collections/purchase_orders/records/"+authorPOID,
		strings.NewReader(fmt.Sprintf(`{"description":"Workwear allowance and boots","po_group":"%s"}`, groupID)), authorHeaders)
	mustStatus(t, res, http.StatusOK)
}

func TestPurchaseOrderGroupCreateRejectsInvalidRows(t *testing.T) {
	app := testutils.SetupTestApp(t)
	t.Cleanup(app.Cleanup)
	superPOHeaders := map[string]string{"Authorization": poGroupToken(t, "superpo@test.com")}

	rows := append(poGroupValidRows(), map[string]string{"email": "nobody@example.com", "total": "40"})
	res := performTestAPIRequest(t, app, http.MethodPost, "/api/purchase_orders/groups",
		strings.NewReader(poGroupBody(t, "grouptest000002", []string{"total"}, rows, nil)), superPOHeaders)
	mustStatus(t, res, http.StatusBadRequest)
	body := mustReadBody(t, res)
	if !strings.Contains(body, `"code":"rows_invalid"`) || !strings.Contains(body, `"error_count":1`) {
		t.Fatalf("expected rows_invalid with the preview, got %s", body)
	}
	if count := countPOGroupRows(t, app, `SELECT COUNT(*) AS count FROM purchase_orders WHERE po_group != ''`, dbx.Params{}); count != 0 {
		t.Fatalf("an invalid upload created %d purchase orders", count)
	}
	if count := countPOGroupRows(t, app, `SELECT COUNT(*) AS count FROM purchase_order_groups`, dbx.Params{}); count != 0 {
		t.Fatalf("an invalid upload created %d groups", count)
	}

	res = performTestAPIRequest(t, app, http.MethodPost, "/api/purchase_orders/groups",
		strings.NewReader(poGroupBody(t, "NOT-VALID", []string{"total"}, poGroupValidRows(), nil)), superPOHeaders)
	mustStatus(t, res, http.StatusBadRequest)
	if !strings.Contains(mustReadBody(t, res), `"group_id":{"code":"invalid_group_id"`) {
		t.Fatalf("expected invalid_group_id, got %s", mustReadBody(t, res))
	}

	res = performTestAPIRequest(t, app, http.MethodPost, "/api/purchase_orders/groups",
		strings.NewReader(poGroupBody(t, "grouptest000007", []string{"total"}, poGroupValidRows(), map[string]any{"approver": ""})), superPOHeaders)
	mustStatus(t, res, http.StatusBadRequest)
	if !strings.Contains(mustReadBody(t, res), `"approver":{"code":"required"`) {
		t.Fatalf("expected the approver to be required, got %s", mustReadBody(t, res))
	}

	res = performTestAPIRequest(t, app, http.MethodPost, "/api/purchase_orders/groups",
		strings.NewReader(poGroupBody(t, "grouptest000003", []string{"total"}, poGroupValidRows(), nil)),
		map[string]string{"Authorization": poGroupToken(t, "time@test.com")})
	mustStatus(t, res, http.StatusForbidden)
}

func TestPurchaseOrderCreateRejectsPOGroup(t *testing.T) {
	scenarios := []tests.ApiScenario{
		{
			Name:   "a collection create can't join a purchase order group",
			Method: http.MethodPost,
			URL:    "/api/collections/purchase_orders/records",
			Body: strings.NewReader(fmt.Sprintf(`{
				"uid": "rzr98oadsp9qc11",
				"date": "2026-11-01",
				"division": "%s",
				"description": "test purchase order",
				"payment_type": "OnAccount",
				"total": 100,
				"vendor": "%s",
				"approver": "%s",
				"status": "Unapproved",
				"type": "One-Time",
				"kind": "%s",
				"po_group": "grouptest000004"
			}`, poGroupDivision, poGroupVendor, poGroupApprover, poGroupCapitalKind)),
			Headers:        map[string]string{"Authorization": poGroupToken(t, "time@test.com"), "Content-Type": "application/json"},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedContent: []string{
				`"po_group":{"code":"immutable_field"`,
			},
			TestAppFactory: testutils.SetupTestApp,
		},
	}
	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

// The priority second approver gets one email for a group, however many of its
// purchase orders need second approval, instead of one email per purchase
// order.
func TestPurchaseOrderGroupSecondApprovalEmail(t *testing.T) {
	app := testutils.SetupTestApp(t)
	t.Cleanup(app.Cleanup)
	const groupID = "grouptest000005"
	const secondApprover = "6bq4j0eb26631dy"

	rows := []map[string]string{
		{"email": "time@test.com", "total": "600"},
		{"email": "author@soup.com", "total": "650"},
		{"email": "hal@2005.com", "total": "80"},
	}
	res := performTestAPIRequest(t, app, http.MethodPost, "/api/purchase_orders/groups",
		strings.NewReader(poGroupBody(t, groupID, []string{"total"}, rows, map[string]any{"priority_second_approver": secondApprover})),
		map[string]string{"Authorization": poGroupToken(t, "superpo@test.com")})
	mustStatus(t, res, http.StatusOK)

	ids := []string{}
	if err := app.DB().NewQuery(`SELECT id FROM purchase_orders WHERE po_group = {:group} ORDER BY total DESC`).
		Bind(dbx.Params{"group": groupID}).Column(&ids); err != nil {
		t.Fatal(err)
	}
	approverHeaders := map[string]string{"Authorization": poGroupToken(t, "fatt@mac.com")}
	for _, id := range ids {
		res = performTestAPIRequest(t, app, http.MethodPost, "/api/purchase_orders/"+id+"/approve", strings.NewReader(`{}`), approverHeaders)
		mustStatus(t, res, http.StatusOK)
	}

	countSent := func(code string) int {
		return countPOGroupRows(t, app, `
			SELECT COUNT(*) AS count
			FROM notifications n JOIN notification_templates nt ON nt.id = n.template
			WHERE nt.code = {:code} AND n.recipient = {:recipient}
		`, dbx.Params{"code": code, "recipient": secondApprover})
	}
	if count := countSent("po_group_second_approval_required"); count != 1 {
		t.Fatalf("expected one group second approval email, got %d", count)
	}
	if count := countSent("po_priority_second_approval_required"); count != 0 {
		t.Fatalf("expected no per purchase order second approval emails, got %d", count)
	}
	found := false
	for _, message := range app.TestMailer.Messages() {
		if message.Subject == "Second approval needed: purchase orders for 3 staff" {
			found = true
			if !strings.Contains(message.Text, "Group Buyer created 3 purchase orders") ||
				!strings.Contains(message.Text, "Description: Workwear allowance") ||
				!strings.Contains(message.Text, "Vendor: Big Vendor Industries") {
				t.Fatalf("unexpected group second approval email: %s", message.Text)
			}
		}
	}
	if !found {
		t.Fatalf("the group second approval email was not sent: %+v", app.TestMailer.Messages())
	}
}

// A super_po holder with the corporate_branch claim may put a group on the
// Corporate branch. The owners lack the claim, but they chose neither the
// branch nor the purchase order, so they can still edit their purchase orders
// and charge expenses to them.
func TestPurchaseOrderGroupOwnersUseTheCreatorsBranch(t *testing.T) {
	app := testutils.SetupTestApp(t)
	t.Cleanup(app.Cleanup)
	const groupID = "grouptest000006"

	claim, err := app.FindFirstRecordByData("claims", "name", "corporate_branch")
	if err != nil {
		t.Fatal(err)
	}
	userClaims, err := app.FindCollectionByNameOrId("user_claims")
	if err != nil {
		t.Fatal(err)
	}
	grant := core.NewRecord(userClaims)
	grant.Set("uid", poGroupCreatorUID)
	grant.Set("cid", claim.Id)
	if err := app.Save(grant); err != nil {
		t.Fatal(err)
	}

	rows := []map[string]string{
		{"email": "corp.noclaim@example.com", "total": "150"},
		{"email": "u_no_claims@example.com", "total": "120"},
	}
	res := performTestAPIRequest(t, app, http.MethodPost, "/api/purchase_orders/groups",
		strings.NewReader(poGroupBody(t, groupID, []string{"total"}, rows, map[string]any{"branch": poGroupCorporate})),
		map[string]string{"Authorization": poGroupToken(t, "superpo@test.com")})
	mustStatus(t, res, http.StatusOK)

	poFor := func(email string) string {
		var id string
		if err := app.DB().NewQuery(`
			SELECT po.id FROM purchase_orders po JOIN users u ON u.id = po.uid
			WHERE po.po_group = {:group} AND u.email = {:email}
		`).Bind(dbx.Params{"group": groupID, "email": email}).Row(&id); err != nil {
			t.Fatal(err)
		}
		return id
	}

	// The owner edits their unapproved purchase order and keeps its branch.
	res = performTestAPIRequest(t, app, http.MethodPatch, "/api/collections/purchase_orders/records/"+poFor("u_no_claims@example.com"),
		strings.NewReader(`{"description":"Workwear allowance and boots"}`),
		map[string]string{"Authorization": poGroupToken(t, "u_no_claims@example.com")})
	mustStatus(t, res, http.StatusOK)
	if body := mustReadBody(t, res); !strings.Contains(body, `"branch":"`+poGroupCorporate+`"`) {
		t.Fatalf("expected the edit to keep the Corporate branch, got %s", body)
	}

	// The owner charges an expense to their approved purchase order. The
	// expense takes the purchase order's branch.
	poID := poFor("corp.noclaim@example.com")
	res = performTestAPIRequest(t, app, http.MethodPost, "/api/purchase_orders/"+poID+"/approve", strings.NewReader(`{}`),
		map[string]string{"Authorization": poGroupToken(t, "fatt@mac.com")})
	mustStatus(t, res, http.StatusOK)
	expense := &bytes.Buffer{}
	form := multipart.NewWriter(expense)
	for field, value := range map[string]string{
		"uid":            "u_corp_noclaim",
		"date":           "2026-11-01",
		"division":       poGroupDivision,
		"description":    "work boots",
		"payment_type":   "OnAccount",
		"total":          "100",
		"vendor":         poGroupVendor,
		"purchase_order": poID,
	} {
		if err := form.WriteField(field, value); err != nil {
			t.Fatal(err)
		}
	}
	receipt, err := form.CreateFormFile("attachment", "receipt.png")
	if err != nil {
		t.Fatal(err)
	}
	// A PNG header is enough for the file type check.
	if _, err := receipt.Write([]byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}); err != nil {
		t.Fatal(err)
	}
	if err := form.Close(); err != nil {
		t.Fatal(err)
	}
	res = performTestAPIRequest(t, app, http.MethodPost, "/api/collections/expenses/records", expense,
		map[string]string{"Authorization": poGroupToken(t, "corp.noclaim@example.com"), "Content-Type": form.FormDataContentType()})
	mustStatus(t, res, http.StatusOK)
	if body := mustReadBody(t, res); !strings.Contains(body, `"branch":"`+poGroupCorporate+`"`) {
		t.Fatalf("expected the expense to take the Corporate branch, got %s", body)
	}
}
