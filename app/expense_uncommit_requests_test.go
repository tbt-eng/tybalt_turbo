package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"
	"tybalt/internal/testutils"
	"tybalt/notifications"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

// Fixtures: reimbtestuser01 owns a committed expense of every payment type. Their
// commit dates are too old to request, so each test commits the one it uses
// relative to today. The commit holder and admin holder hold only that claim.
const (
	uncommitCommitEmail   = "fakemanager@fakesite.xyz"
	uncommitAdminEmail    = "admin.only@example.com"
	uncommitPayablesEmail = "book@keeper.com"
	uncommitOtherEmail    = "reimbursements@example.com"

	uncommitOnAccountExpense = "reimbonaccount1"
	uncommitPOExpense        = "6569323gg8184uh"
	uncommitPOExpensePO      = "0pia83nnprdlzf8"

	uncommitReasonDuplicate = "uncreason000001"
	uncommitReasonWrongCurr = "uncreason000004"
	uncommitReasonMistake   = "uncreason000010"
)

func uncommitToken(t testing.TB, email string) map[string]string {
	t.Helper()
	token, err := testutils.GenerateRecordToken("users", email)
	if err != nil {
		t.Fatal(err)
	}
	return map[string]string{"Authorization": token, "Content-Type": "application/json"}
}

// setExpenseCommitDates writes the dates directly so the gates can be tested on
// either side of today.
func setExpenseCommitDates(t testing.TB, app core.App, id string, committed time.Time, payPeriodEnding string) {
	t.Helper()
	if _, err := app.DB().NewQuery(`
		UPDATE expenses SET committed = {:committed}, pay_period_ending = {:ending} WHERE id = {:id}
	`).Bind(dbx.Params{
		"committed": committed.UTC().Format("2006-01-02 15:04:05.000Z"),
		"ending":    payPeriodEnding,
		"id":        id,
	}).Execute(); err != nil {
		t.Fatal(err)
	}
}

func commitExpenseToday(t testing.TB, app core.App, id string) {
	t.Helper()
	setExpenseCommitDates(t, app, id, time.Now(), time.Now().Format(time.DateOnly))
}

func uncommitRequestBody(reason, outcome, note string) *strings.Reader {
	body, _ := json.Marshal(map[string]string{"reason": reason, "outcome": outcome, "note": note})
	return strings.NewReader(string(body))
}

func createUncommitRequest(t testing.TB, app *tests.TestApp, expenseID string, reason, outcome, note string) string {
	t.Helper()
	res := performTestAPIRequest(t, app, http.MethodPost, "/api/expenses/"+expenseID+"/uncommit_requests",
		uncommitRequestBody(reason, outcome, note), uncommitToken(t, uncommitCommitEmail))
	mustStatus(t, res, http.StatusOK)
	var created struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(res.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}
	return created.ID
}

type uncommitApprovalResult struct {
	ID   string `json:"id"`
	OK   bool   `json:"ok"`
	Code string `json:"code"`
}

func approveUncommitRequests(t testing.TB, app *tests.TestApp, ids ...string) []uncommitApprovalResult {
	t.Helper()
	body, _ := json.Marshal(map[string][]string{"ids": ids})
	res := performTestAPIRequest(t, app, http.MethodPost, "/api/expense_uncommit_requests/approve",
		strings.NewReader(string(body)), uncommitToken(t, uncommitAdminEmail))
	mustStatus(t, res, http.StatusOK)
	var response struct {
		Results []uncommitApprovalResult `json:"results"`
	}
	if err := json.Unmarshal(res.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	return response.Results
}

func mustUncommitErrorCode(t testing.TB, body []byte, want string) {
	t.Helper()
	var response struct {
		Code string `json:"code"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("failed to decode %s: %v", body, err)
	}
	if response.Code != want {
		t.Fatalf("expected code %q, got %s", want, body)
	}
}

func setExpensesConfig(t testing.TB, app core.App, values map[string]any) {
	t.Helper()
	config, err := app.FindFirstRecordByData("app_config", "key", "expenses")
	if err != nil {
		collection, err := app.FindCollectionByNameOrId("app_config")
		if err != nil {
			t.Fatal(err)
		}
		config = core.NewRecord(collection)
		config.Set("key", "expenses")
	}
	config.Set("value", values)
	if err := app.Save(config); err != nil {
		t.Fatal(err)
	}
}

func uidFor(t testing.TB, app core.App, email string) string {
	t.Helper()
	user, err := app.FindAuthRecordByEmail("users", email)
	if err != nil {
		t.Fatal(err)
	}
	return user.Id
}

func TestExpenseUncommitRequestPermissions(t *testing.T) {
	app := testutils.SetupTestApp(t)
	commitExpenseToday(t, app, uncommitOnAccountExpense)
	url := "/api/expenses/" + uncommitOnAccountExpense + "/uncommit_requests"

	for _, email := range []string{uncommitAdminEmail, uncommitPayablesEmail, uncommitOtherEmail} {
		res := performTestAPIRequest(t, app, http.MethodPost, url,
			uncommitRequestBody(uncommitReasonWrongCurr, "uncommit_only", ""), uncommitToken(t, email))
		mustStatus(t, res, http.StatusForbidden)
	}

	id := createUncommitRequest(t, app, uncommitOnAccountExpense, uncommitReasonWrongCurr, "uncommit_only", "")

	// Only the requester can cancel, so the admin holder can't.
	res := performTestAPIRequest(t, app, http.MethodPost, "/api/expense_uncommit_requests/"+id+"/cancel", nil, uncommitToken(t, uncommitAdminEmail))
	mustStatus(t, res, http.StatusForbidden)

	// Only admin holders decide, so the commit holder can't.
	approveBody := fmt.Sprintf(`{"ids":["%s"]}`, id)
	res = performTestAPIRequest(t, app, http.MethodPost, "/api/expense_uncommit_requests/approve", strings.NewReader(approveBody), uncommitToken(t, uncommitCommitEmail))
	mustStatus(t, res, http.StatusForbidden)
	res = performTestAPIRequest(t, app, http.MethodPost, "/api/expense_uncommit_requests/"+id+"/decline", strings.NewReader(`{"comment":"no thanks"}`), uncommitToken(t, uncommitCommitEmail))
	mustStatus(t, res, http.StatusForbidden)

	// Others can't see the queue or its reasons.
	for _, path := range []string{"/api/expense_uncommit_requests", "/api/expense_uncommit_requests/reasons"} {
		res = performTestAPIRequest(t, app, http.MethodGet, path, nil, uncommitToken(t, uncommitOtherEmail))
		mustStatus(t, res, http.StatusForbidden)
		for _, email := range []string{uncommitCommitEmail, uncommitAdminEmail} {
			res = performTestAPIRequest(t, app, http.MethodGet, path, nil, uncommitToken(t, email))
			mustStatus(t, res, http.StatusOK)
		}
	}

	res = performTestAPIRequest(t, app, http.MethodPost, "/api/expense_uncommit_requests/"+id+"/cancel", nil, uncommitToken(t, uncommitCommitEmail))
	mustStatus(t, res, http.StatusOK)
}

func TestExpenseUncommitPayrollGate(t *testing.T) {
	expenses := map[string]string{
		"Expense":               "reimbexpense001",
		"PersonalReimbursement": "reimbpersonal01",
		"Allowance":             "reimballowance1",
		"Mileage":               "reimbmileage002",
	}
	today := time.Now()

	for paymentType, id := range expenses {
		t.Run(paymentType, func(t *testing.T) {
			app := testutils.SetupTestApp(t)
			url := "/api/expenses/" + id + "/uncommit_requests"
			token := uncommitToken(t, uncommitCommitEmail)

			// One day past the window. A recent commit date doesn't help.
			setExpenseCommitDates(t, app, id, today, today.AddDate(0, 0, -8).Format(time.DateOnly))
			res := performTestAPIRequest(t, app, http.MethodPost, url, uncommitRequestBody(uncommitReasonWrongCurr, "uncommit_only", ""), token)
			mustStatus(t, res, http.StatusBadRequest)
			mustUncommitErrorCode(t, res.Body.Bytes(), "payroll_closed")

			// The window's last day. An old commit date doesn't hurt.
			setExpenseCommitDates(t, app, id, today.AddDate(0, 0, -200), today.AddDate(0, 0, -7).Format(time.DateOnly))
			res = performTestAPIRequest(t, app, http.MethodPost, url, uncommitRequestBody(uncommitReasonWrongCurr, "uncommit_only", ""), token)
			mustStatus(t, res, http.StatusOK)
		})
	}

	t.Run("a blank pay period is blocked", func(t *testing.T) {
		app := testutils.SetupTestApp(t)
		setExpenseCommitDates(t, app, "reimbexpense001", today, "")
		res := performTestAPIRequest(t, app, http.MethodPost, "/api/expenses/reimbexpense001/uncommit_requests",
			uncommitRequestBody(uncommitReasonWrongCurr, "uncommit_only", ""), uncommitToken(t, uncommitCommitEmail))
		mustStatus(t, res, http.StatusBadRequest)
		mustUncommitErrorCode(t, res.Body.Bytes(), "payroll_closed")
	})

	t.Run("the window comes from app_config", func(t *testing.T) {
		app := testutils.SetupTestApp(t)
		setExpensesConfig(t, app, map[string]any{"uncommit_payroll_lock_days": 2})
		setExpenseCommitDates(t, app, "reimbexpense001", today, today.AddDate(0, 0, -3).Format(time.DateOnly))
		res := performTestAPIRequest(t, app, http.MethodPost, "/api/expenses/reimbexpense001/uncommit_requests",
			uncommitRequestBody(uncommitReasonWrongCurr, "uncommit_only", ""), uncommitToken(t, uncommitCommitEmail))
		mustStatus(t, res, http.StatusBadRequest)
		mustUncommitErrorCode(t, res.Body.Bytes(), "payroll_closed")
	})
}

func TestExpenseUncommitNonPayrollGate(t *testing.T) {
	expenses := map[string]string{
		"OnAccount":           "reimbonaccount1",
		"CorporateCreditCard": "reimbcorpcard01",
		"FuelCard":            "reimbfuelcard01",
	}
	today := time.Now()

	for paymentType, id := range expenses {
		t.Run(paymentType, func(t *testing.T) {
			app := testutils.SetupTestApp(t)
			url := "/api/expenses/" + id + "/uncommit_requests"
			token := uncommitToken(t, uncommitCommitEmail)

			// A closed pay period doesn't matter to these types.
			setExpenseCommitDates(t, app, id, today.AddDate(0, 0, -91), today.AddDate(0, 0, -100).Format(time.DateOnly))
			res := performTestAPIRequest(t, app, http.MethodPost, url, uncommitRequestBody(uncommitReasonWrongCurr, "uncommit_only", ""), token)
			mustStatus(t, res, http.StatusBadRequest)
			mustUncommitErrorCode(t, res.Body.Bytes(), "too_old")

			setExpenseCommitDates(t, app, id, today.AddDate(0, 0, -90), today.AddDate(0, 0, -100).Format(time.DateOnly))
			res = performTestAPIRequest(t, app, http.MethodPost, url, uncommitRequestBody(uncommitReasonWrongCurr, "uncommit_only", ""), token)
			mustStatus(t, res, http.StatusOK)
		})
	}

	t.Run("the limit comes from app_config", func(t *testing.T) {
		app := testutils.SetupTestApp(t)
		setExpensesConfig(t, app, map[string]any{"uncommit_non_payroll_lock_days": 10})
		setExpenseCommitDates(t, app, uncommitOnAccountExpense, today.AddDate(0, 0, -11), "")
		res := performTestAPIRequest(t, app, http.MethodPost, "/api/expenses/"+uncommitOnAccountExpense+"/uncommit_requests",
			uncommitRequestBody(uncommitReasonWrongCurr, "uncommit_only", ""), uncommitToken(t, uncommitCommitEmail))
		mustStatus(t, res, http.StatusBadRequest)
		mustUncommitErrorCode(t, res.Body.Bytes(), "too_old")
	})
}

func TestExpenseUncommitRequestValidation(t *testing.T) {
	app := testutils.SetupTestApp(t)
	commitExpenseToday(t, app, uncommitOnAccountExpense)
	url := "/api/expenses/" + uncommitOnAccountExpense + "/uncommit_requests"
	token := uncommitToken(t, uncommitCommitEmail)

	cases := []struct {
		reason, outcome, note, code string
	}{
		{"nosuchreason000", "uncommit_only", "", "invalid_reason"},
		{uncommitReasonWrongCurr, "nonsense", "", "invalid_outcome"},
		{uncommitReasonDuplicate, "uncommit_and_reject", "  ", "note_required"},
	}
	for _, c := range cases {
		res := performTestAPIRequest(t, app, http.MethodPost, url, uncommitRequestBody(c.reason, c.outcome, c.note), token)
		mustStatus(t, res, http.StatusBadRequest)
		mustUncommitErrorCode(t, res.Body.Bytes(), c.code)
	}

	// An uncommitted expense can't be requested.
	res := performTestAPIRequest(t, app, http.MethodPost, "/api/expenses/eqhozipupteogp8/uncommit_requests",
		uncommitRequestBody(uncommitReasonWrongCurr, "uncommit_only", ""), token)
	mustStatus(t, res, http.StatusBadRequest)
	mustUncommitErrorCode(t, res.Body.Bytes(), "record_not_committed")

	// Neither can an imported one.
	if _, err := app.DB().NewQuery(`UPDATE expenses SET _imported = 1 WHERE id = 'reimbcorpcard01'`).Execute(); err != nil {
		t.Fatal(err)
	}
	commitExpenseToday(t, app, "reimbcorpcard01")
	res = performTestAPIRequest(t, app, http.MethodPost, "/api/expenses/reimbcorpcard01/uncommit_requests",
		uncommitRequestBody(uncommitReasonWrongCurr, "uncommit_only", ""), token)
	mustStatus(t, res, http.StatusBadRequest)
	mustUncommitErrorCode(t, res.Body.Bytes(), "record_imported")
}

func TestExpenseUncommitOnePendingRequestPerExpense(t *testing.T) {
	app := testutils.SetupTestApp(t)
	commitExpenseToday(t, app, uncommitOnAccountExpense)

	id := createUncommitRequest(t, app, uncommitOnAccountExpense, uncommitReasonWrongCurr, "uncommit_only", "")
	res := performTestAPIRequest(t, app, http.MethodPost, "/api/expenses/"+uncommitOnAccountExpense+"/uncommit_requests",
		uncommitRequestBody(uncommitReasonWrongCurr, "uncommit_only", ""), uncommitToken(t, uncommitCommitEmail))
	mustStatus(t, res, http.StatusConflict)
	mustUncommitErrorCode(t, res.Body.Bytes(), "request_pending")

	res = performTestAPIRequest(t, app, http.MethodPost, "/api/expense_uncommit_requests/"+id+"/decline",
		strings.NewReader(`{"comment":"already fixed"}`), uncommitToken(t, uncommitAdminEmail))
	mustStatus(t, res, http.StatusOK)

	// Declining leaves the expense committed and allows a new request.
	expense, err := app.FindRecordById("expenses", uncommitOnAccountExpense)
	if err != nil {
		t.Fatal(err)
	}
	if expense.GetDateTime("committed").IsZero() {
		t.Fatal("declining uncommitted the expense")
	}
	createUncommitRequest(t, app, uncommitOnAccountExpense, uncommitReasonWrongCurr, "uncommit_only", "")
}

func TestExpenseUncommitDeclineRequiresComment(t *testing.T) {
	app := testutils.SetupTestApp(t)
	commitExpenseToday(t, app, uncommitOnAccountExpense)
	id := createUncommitRequest(t, app, uncommitOnAccountExpense, uncommitReasonWrongCurr, "uncommit_only", "")

	res := performTestAPIRequest(t, app, http.MethodPost, "/api/expense_uncommit_requests/"+id+"/decline",
		strings.NewReader(`{"comment":" "}`), uncommitToken(t, uncommitAdminEmail))
	mustStatus(t, res, http.StatusBadRequest)
	mustUncommitErrorCode(t, res.Body.Bytes(), "comment_too_short")
}

func TestExpenseUncommitApproveAndReject(t *testing.T) {
	app := setupClosedPurchaseOrderForUncommit(t, uncommitPOExpensePO)
	commitExpenseToday(t, app, uncommitPOExpense)
	requester := uidFor(t, app, uncommitCommitEmail)

	before, err := app.FindRecordById("expenses", uncommitPOExpense)
	if err != nil {
		t.Fatal(err)
	}
	committedBefore := before.GetString("committed")
	committerBefore := before.GetString("committer")

	id := createUncommitRequest(t, app, uncommitPOExpense, uncommitReasonDuplicate, "uncommit_and_reject", "It is on 2610-0001")

	results := approveUncommitRequests(t, app, id)
	if len(results) != 1 || !results[0].OK {
		t.Fatalf("expected one successful approval, got %+v", results)
	}

	expense, err := app.FindRecordById("expenses", uncommitPOExpense)
	if err != nil {
		t.Fatal(err)
	}
	if !expense.GetDateTime("committed").IsZero() || expense.GetString("committer") != "" ||
		expense.GetString("committed_week_ending") != "" || expense.GetString("pay_period_ending") != "" {
		t.Fatalf("commit fields not cleared: %v", expense.PublicExport())
	}
	if expense.GetDateTime("rejected").IsZero() {
		t.Fatal("expected the expense to be rejected")
	}
	if got := expense.GetString("rejector"); got != requester {
		t.Fatalf("expected rejector %q, got %q", requester, got)
	}
	if got, want := expense.GetString("rejection_reason"), "Duplicate PO, invoice is on another PO: It is on 2610-0001"; got != want {
		t.Fatalf("expected rejection reason %q, got %q", want, got)
	}
	if expense.GetDateTime("approved").IsZero() {
		t.Fatal("expected the expense to stay approved")
	}

	po, err := app.FindRecordById("purchase_orders", uncommitPOExpensePO)
	if err != nil {
		t.Fatal(err)
	}
	if po.GetString("status") != "Active" {
		t.Fatalf("expected the purchase order to reopen, got %q", po.GetString("status"))
	}

	request, err := app.FindRecordById("expense_uncommit_requests", id)
	if err != nil {
		t.Fatal(err)
	}
	if request.GetString("status") != "approved" {
		t.Fatalf("expected approved, got %q", request.GetString("status"))
	}
	if request.GetString("decided_by") != uidFor(t, app, uncommitAdminEmail) || request.GetDateTime("decided_at").IsZero() {
		t.Fatalf("decision not recorded: %v", request.PublicExport())
	}
	if request.GetString("committed") != committedBefore || request.GetString("committer") != committerBefore {
		t.Fatalf("commit state not kept: %v", request.PublicExport())
	}
	if request.GetString("payment_type") != "OnAccount" || request.GetFloat("total") != before.GetFloat("total") {
		t.Fatalf("snapshot not kept: %v", request.PublicExport())
	}

	if count := countNotifications(t, app, "expense_rejected", expense.GetString("uid")); count != 1 {
		t.Fatalf("expected the employee to be notified of the rejection, got %d", count)
	}
	if count := countNotifications(t, app, "expense_uncommit_approved", requester); count != 1 {
		t.Fatalf("expected the requester to be notified, got %d", count)
	}

	// A second approval of the same request is refused.
	again := approveUncommitRequests(t, app, id)
	if again[0].OK || again[0].Code != "request_not_pending" {
		t.Fatalf("expected request_not_pending, got %+v", again[0])
	}
}

func countNotifications(t testing.TB, app core.App, templateCode string, recipient string) int {
	t.Helper()
	var count int
	if err := app.DB().NewQuery(`
		SELECT COUNT(*)
		FROM notifications n
		INNER JOIN notification_templates t ON t.id = n.template
		WHERE t.code = {:code} AND n.recipient = {:recipient}
	`).Bind(dbx.Params{"code": templateCode, "recipient": recipient}).Row(&count); err != nil {
		t.Fatal(err)
	}
	return count
}

func TestExpenseUncommitApproveOnlyReturnsToCommitQueue(t *testing.T) {
	app := testutils.SetupTestApp(t)
	commitExpenseToday(t, app, uncommitOnAccountExpense)
	id := createUncommitRequest(t, app, uncommitOnAccountExpense, uncommitReasonMistake, "uncommit_only", "")

	if results := approveUncommitRequests(t, app, id); !results[0].OK {
		t.Fatalf("expected approval, got %+v", results[0])
	}

	expense, err := app.FindRecordById("expenses", uncommitOnAccountExpense)
	if err != nil {
		t.Fatal(err)
	}
	if !expense.GetDateTime("committed").IsZero() {
		t.Fatal("expected the expense to be uncommitted")
	}
	if !expense.GetDateTime("rejected").IsZero() {
		t.Fatal("expected the expense not to be rejected")
	}

	res := performTestAPIRequest(t, app, http.MethodGet, "/api/expenses/commit_queue", nil, uncommitToken(t, uncommitCommitEmail))
	mustStatus(t, res, http.StatusOK)
	if !strings.Contains(res.Body.String(), uncommitOnAccountExpense) {
		t.Fatal("expected the expense back in the commit queue")
	}
}

func TestExpenseUncommitBulkApproveReportsEachRequest(t *testing.T) {
	app := testutils.SetupTestApp(t)
	today := time.Now()
	commitExpenseToday(t, app, "reimbonaccount1")
	commitExpenseToday(t, app, "reimbcorpcard01")
	commitExpenseToday(t, app, "reimbfuelcard01")

	good := createUncommitRequest(t, app, "reimbonaccount1", uncommitReasonWrongCurr, "uncommit_only", "")
	closed := createUncommitRequest(t, app, "reimbcorpcard01", uncommitReasonWrongCurr, "uncommit_only", "")
	alsoGood := createUncommitRequest(t, app, "reimbfuelcard01", uncommitReasonWrongCurr, "uncommit_only", "")

	// The gate closes while this one waits.
	setExpenseCommitDates(t, app, "reimbcorpcard01", today.AddDate(0, 0, -91), "")

	results := approveUncommitRequests(t, app, good, closed, "nosuchrequest00", alsoGood)
	if len(results) != 4 {
		t.Fatalf("expected 4 results, got %+v", results)
	}
	if !results[0].OK || !results[3].OK {
		t.Fatalf("expected the valid requests to apply, got %+v", results)
	}
	if results[1].OK || results[1].Code != "too_old" {
		t.Fatalf("expected too_old, got %+v", results[1])
	}
	if results[2].OK || results[2].Code != "record_not_found" {
		t.Fatalf("expected record_not_found, got %+v", results[2])
	}

	expense, err := app.FindRecordById("expenses", "reimbcorpcard01")
	if err != nil {
		t.Fatal(err)
	}
	if expense.GetDateTime("committed").IsZero() {
		t.Fatal("the refused request should have left the expense committed")
	}
}

func TestExpenseUncommitQueueFlagsClosedGate(t *testing.T) {
	app := testutils.SetupTestApp(t)
	today := time.Now()
	commitExpenseToday(t, app, "reimbexpense001")
	createUncommitRequest(t, app, "reimbexpense001", uncommitReasonWrongCurr, "uncommit_only", "")

	// The snapshot is what the queue reads, so age it directly.
	if _, err := app.DB().NewQuery(`UPDATE expense_uncommit_requests SET pay_period_ending = {:ending} WHERE expense = 'reimbexpense001'`).
		Bind(dbx.Params{"ending": today.AddDate(0, 0, -8).Format(time.DateOnly)}).Execute(); err != nil {
		t.Fatal(err)
	}

	res := performTestAPIRequest(t, app, http.MethodGet, "/api/expense_uncommit_requests", nil, uncommitToken(t, uncommitAdminEmail))
	mustStatus(t, res, http.StatusOK)
	var rows []struct {
		Blocked *struct {
			Code string `json:"code"`
		} `json:"blocked"`
	}
	if err := json.Unmarshal(res.Body.Bytes(), &rows); err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].Blocked == nil || rows[0].Blocked.Code != "payroll_closed" {
		t.Fatalf("expected one row flagged payroll_closed, got %s", res.Body.String())
	}
}

func TestExpenseUncommitDirectUncommitKeepsHistory(t *testing.T) {
	t.Run("closes a pending request as uncommit only", func(t *testing.T) {
		app := testutils.SetupTestApp(t)
		commitExpenseToday(t, app, uncommitOnAccountExpense)
		id := createUncommitRequest(t, app, uncommitOnAccountExpense, uncommitReasonWrongCurr, "uncommit_and_reject", "")

		res := performTestAPIRequest(t, app, http.MethodPost, "/api/expenses/"+uncommitOnAccountExpense+"/uncommit", nil, uncommitToken(t, uncommitAdminEmail))
		mustStatus(t, res, http.StatusOK)

		request, err := app.FindRecordById("expense_uncommit_requests", id)
		if err != nil {
			t.Fatal(err)
		}
		if request.GetString("status") != "approved" || request.GetString("outcome") != "uncommit_only" {
			t.Fatalf("expected approved and uncommit_only, got %v", request.PublicExport())
		}
		if count := countNotifications(t, app, "expense_uncommit_approved", request.GetString("requester")); count != 1 {
			t.Fatalf("expected the requester to be told, got %d", count)
		}
	})

	t.Run("writes a row when there is no request", func(t *testing.T) {
		app := testutils.SetupTestApp(t)
		res := performTestAPIRequest(t, app, http.MethodPost, "/api/expenses/xg2yeucklhgbs3n/uncommit", nil, uncommitToken(t, uncommitAdminEmail))
		mustStatus(t, res, http.StatusOK)

		request, err := app.FindFirstRecordByData("expense_uncommit_requests", "expense", "xg2yeucklhgbs3n")
		if err != nil {
			t.Fatal(err)
		}
		admin := uidFor(t, app, uncommitAdminEmail)
		if request.GetString("status") != "approved" || request.GetString("requester") != admin || request.GetString("decided_by") != admin {
			t.Fatalf("unexpected row: %v", request.PublicExport())
		}
		if request.GetString("committed") == "" {
			t.Fatal("expected the commit state in the snapshot")
		}
	})
}

// A rejected expense on a duplicate PO has to be deleted before the PO can be
// cancelled, so a request must not stop that.
func TestExpenseUncommitHistorySurvivesDeletingTheExpense(t *testing.T) {
	app := testutils.SetupTestApp(t)
	commitExpenseToday(t, app, uncommitOnAccountExpense)
	id := createUncommitRequest(t, app, uncommitOnAccountExpense, uncommitReasonWrongCurr, "uncommit_and_reject", "")
	if results := approveUncommitRequests(t, app, id); !results[0].OK {
		t.Fatalf("expected approval, got %+v", results[0])
	}

	expense, err := app.FindRecordById("expenses", uncommitOnAccountExpense)
	if err != nil {
		t.Fatal(err)
	}
	if err := app.Delete(expense); err != nil {
		t.Fatalf("expected the expense to be deletable, got %v", err)
	}

	res := performTestAPIRequest(t, app, http.MethodGet, "/api/expense_uncommit_requests?status=decided", nil, uncommitToken(t, uncommitAdminEmail))
	mustStatus(t, res, http.StatusOK)
	var rows []struct {
		ID          string `json:"id"`
		Expense     string `json:"expense"`
		PaymentType string `json:"payment_type"`
		ReasonName  string `json:"reason_name"`
	}
	if err := json.Unmarshal(res.Body.Bytes(), &rows); err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].ID != id || rows[0].Expense != "" || rows[0].PaymentType != "OnAccount" || rows[0].ReasonName != "Wrong currency" {
		t.Fatalf("unexpected history: %s", res.Body.String())
	}
}

func TestExpenseUncommitEditingDisabledBlocksMutations(t *testing.T) {
	app := testutils.SetupTestApp(t)
	commitExpenseToday(t, app, uncommitOnAccountExpense)
	id := createUncommitRequest(t, app, uncommitOnAccountExpense, uncommitReasonWrongCurr, "uncommit_only", "")
	commitExpenseToday(t, app, "reimbcorpcard01")

	setExpensesConfig(t, app, map[string]any{"create_edit_absorb": false})

	commitToken := uncommitToken(t, uncommitCommitEmail)
	adminToken := uncommitToken(t, uncommitAdminEmail)
	cases := []struct {
		path    string
		body    string
		headers map[string]string
	}{
		{"/api/expenses/reimbcorpcard01/uncommit_requests", `{"reason":"uncreason000004","outcome":"uncommit_only"}`, commitToken},
		{"/api/expense_uncommit_requests/approve", fmt.Sprintf(`{"ids":["%s"]}`, id), adminToken},
		{"/api/expense_uncommit_requests/" + id + "/decline", `{"comment":"not now"}`, adminToken},
		{"/api/expense_uncommit_requests/" + id + "/cancel", `{}`, commitToken},
	}
	for _, c := range cases {
		res := performTestAPIRequest(t, app, http.MethodPost, c.path, strings.NewReader(c.body), c.headers)
		mustStatus(t, res, http.StatusForbidden)
	}
}

func TestExpenseUncommitReasonsAndDetailsState(t *testing.T) {
	app := testutils.SetupTestApp(t)
	commitExpenseToday(t, app, uncommitOnAccountExpense)

	res := performTestAPIRequest(t, app, http.MethodGet, "/api/expense_uncommit_requests/reasons", nil, uncommitToken(t, uncommitCommitEmail))
	mustStatus(t, res, http.StatusOK)
	var reasons []struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(res.Body.Bytes(), &reasons); err != nil {
		t.Fatal(err)
	}
	if len(reasons) != 12 || reasons[len(reasons)-1].Name != "Other" {
		t.Fatalf("unexpected reasons: %s", res.Body.String())
	}
	// make_cumulative handles this without an uncommit.
	if strings.Contains(strings.ToLower(res.Body.String()), "cumulative") {
		t.Fatalf("Cumulative should not be a reason: %s", res.Body.String())
	}

	type detailsState struct {
		UncommitRequest *struct {
			CanRequest bool `json:"can_request"`
			Pending    *struct {
				ReasonName string `json:"reason_name"`
			} `json:"pending"`
			Blocked *struct {
				Code string `json:"code"`
			} `json:"blocked"`
		} `json:"uncommit_request"`
	}
	details := func(email string, id string) detailsState {
		t.Helper()
		res := performTestAPIRequest(t, app, http.MethodGet, "/api/expenses/details/"+id, nil, uncommitToken(t, email))
		mustStatus(t, res, http.StatusOK)
		var state detailsState
		if err := json.Unmarshal(res.Body.Bytes(), &state); err != nil {
			t.Fatal(err)
		}
		return state
	}

	state := details(uncommitCommitEmail, uncommitOnAccountExpense)
	if state.UncommitRequest == nil || !state.UncommitRequest.CanRequest {
		t.Fatalf("expected the commit holder to be able to request, got %+v", state.UncommitRequest)
	}

	// The admin holder sees the state but can't request.
	state = details(uncommitAdminEmail, uncommitOnAccountExpense)
	if state.UncommitRequest == nil || state.UncommitRequest.CanRequest {
		t.Fatalf("expected the admin holder not to be able to request, got %+v", state.UncommitRequest)
	}

	// The expense's owner holds neither claim.
	state = details(uncommitOtherEmail, uncommitOnAccountExpense)
	if state.UncommitRequest != nil {
		t.Fatalf("expected no request state for the owner, got %+v", state.UncommitRequest)
	}

	createUncommitRequest(t, app, uncommitOnAccountExpense, uncommitReasonWrongCurr, "uncommit_only", "")
	state = details(uncommitCommitEmail, uncommitOnAccountExpense)
	if state.UncommitRequest.CanRequest || state.UncommitRequest.Pending == nil || state.UncommitRequest.Pending.ReasonName != "Wrong currency" {
		t.Fatalf("expected the pending request, got %+v", state.UncommitRequest)
	}

	// Fixtures committed long ago are blocked and say why.
	state = details(uncommitCommitEmail, "xg2yeucklhgbs3n")
	if state.UncommitRequest.CanRequest || state.UncommitRequest.Blocked == nil || state.UncommitRequest.Blocked.Code != "too_old" {
		t.Fatalf("expected too_old, got %+v", state.UncommitRequest)
	}
}

func TestExpenseUncommitNavBadgeIsForAdminHolders(t *testing.T) {
	app := testutils.SetupTestApp(t)
	commitExpenseToday(t, app, uncommitOnAccountExpense)
	createUncommitRequest(t, app, uncommitOnAccountExpense, uncommitReasonWrongCurr, "uncommit_only", "")

	counts := func(email string) map[string]int {
		t.Helper()
		res := performTestAPIRequest(t, app, http.MethodGet, "/api/nav/badges", nil, uncommitToken(t, email))
		mustStatus(t, res, http.StatusOK)
		out := map[string]int{}
		if err := json.Unmarshal(res.Body.Bytes(), &out); err != nil {
			t.Fatal(err)
		}
		return out
	}
	if got := counts(uncommitAdminEmail)["/expenses/uncommit-requests"]; got != 1 {
		t.Fatalf("expected a count of 1 for the admin holder, got %d", got)
	}
	if _, found := counts(uncommitCommitEmail)["/expenses/uncommit-requests"]; found {
		t.Fatal("expected no count for the commit holder")
	}
}

func TestConvertSystemClosedOneTimePurchaseOrderToCumulative(t *testing.T) {
	url := "/api/purchase_orders/" + uncommitPOExpensePO + "/make_cumulative"

	setup := func(t *testing.T, total float64, closedBySystem bool, imported bool) *tests.TestApp {
		t.Helper()
		app := setupClosedPurchaseOrderForUncommit(t, uncommitPOExpensePO)
		if _, err := app.DB().NewQuery(`
			UPDATE purchase_orders SET total = {:total}, closed_by_system = {:system}, _imported = {:imported} WHERE id = {:id}
		`).Bind(dbx.Params{"total": total, "system": closedBySystem, "imported": imported, "id": uncommitPOExpensePO}).Execute(); err != nil {
			t.Fatal(err)
		}
		return app
	}
	expectUnchanged := func(t *testing.T, app *tests.TestApp) {
		t.Helper()
		po, err := app.FindRecordById("purchase_orders", uncommitPOExpensePO)
		if err != nil {
			t.Fatal(err)
		}
		if po.GetString("type") != "One-Time" || po.GetString("status") != "Closed" {
			t.Fatalf("expected the purchase order unchanged, got %v", po.PublicExport())
		}
	}

	t.Run("converts and reopens when there is room", func(t *testing.T) {
		app := setup(t, 1000, true, false)
		res := performTestAPIRequest(t, app, http.MethodPost, url, nil, uncommitToken(t, uncommitPayablesEmail))
		mustStatus(t, res, http.StatusOK)

		po, err := app.FindRecordById("purchase_orders", uncommitPOExpensePO)
		if err != nil {
			t.Fatal(err)
		}
		if po.GetString("type") != "Cumulative" || po.GetString("status") != "Active" {
			t.Fatalf("expected an Active Cumulative purchase order, got %v", po.PublicExport())
		}
		if !po.GetDateTime("closed").IsZero() || po.GetString("closer") != "" || po.GetBool("closed_by_system") {
			t.Fatalf("expected the closure cleared, got %v", po.PublicExport())
		}

		// The expense itself is untouched.
		expense, err := app.FindRecordById("expenses", uncommitPOExpense)
		if err != nil {
			t.Fatal(err)
		}
		if expense.GetDateTime("committed").IsZero() {
			t.Fatal("expected the expense to stay committed")
		}
	})

	t.Run("refuses when the committed total fills the purchase order", func(t *testing.T) {
		app := setup(t, 485.23, true, false)
		res := performTestAPIRequest(t, app, http.MethodPost, url, nil, uncommitToken(t, uncommitPayablesEmail))
		mustStatus(t, res, http.StatusBadRequest)
		mustUncommitErrorCode(t, res.Body.Bytes(), "po_total_used")
		expectUnchanged(t, app)
	})

	t.Run("refuses a purchase order closed by a person", func(t *testing.T) {
		app := setup(t, 1000, false, false)
		res := performTestAPIRequest(t, app, http.MethodPost, url, nil, uncommitToken(t, uncommitPayablesEmail))
		mustStatus(t, res, http.StatusBadRequest)
		mustUncommitErrorCode(t, res.Body.Bytes(), "po_not_active")
		expectUnchanged(t, app)
	})

	t.Run("refuses an imported purchase order", func(t *testing.T) {
		app := setup(t, 1000, true, true)
		res := performTestAPIRequest(t, app, http.MethodPost, url, nil, uncommitToken(t, uncommitPayablesEmail))
		mustStatus(t, res, http.StatusBadRequest)
		mustUncommitErrorCode(t, res.Body.Bytes(), "po_not_active")
		expectUnchanged(t, app)
	})

	t.Run("still requires the payables_admin claim", func(t *testing.T) {
		app := setup(t, 1000, true, false)
		res := performTestAPIRequest(t, app, http.MethodPost, url, nil, uncommitToken(t, uncommitCommitEmail))
		mustStatus(t, res, http.StatusForbidden)
		expectUnchanged(t, app)
	})
}

func TestExpenseUncommitReminderGoesToAdminHoldersWhileRequestsWait(t *testing.T) {
	app := testutils.SetupTestApp(t)
	admin := uidFor(t, app, uncommitAdminEmail)

	if err := notifications.QueueExpenseUncommitRequestReminders(app, false); err != nil {
		t.Fatal(err)
	}
	if count := countNotifications(t, app, "expense_uncommit_reminder", admin); count != 0 {
		t.Fatalf("expected no reminder with nothing waiting, got %d", count)
	}

	commitExpenseToday(t, app, uncommitOnAccountExpense)
	createUncommitRequest(t, app, uncommitOnAccountExpense, uncommitReasonWrongCurr, "uncommit_only", "")

	// The second run is within a day of the first, so it adds nothing.
	for range 2 {
		if err := notifications.QueueExpenseUncommitRequestReminders(app, false); err != nil {
			t.Fatal(err)
		}
	}
	if count := countNotifications(t, app, "expense_uncommit_reminder", admin); count != 1 {
		t.Fatalf("expected one reminder, got %d", count)
	}
	if count := countNotifications(t, app, "expense_uncommit_reminder", uidFor(t, app, uncommitCommitEmail)); count != 0 {
		t.Fatalf("expected no reminder for the commit holder, got %d", count)
	}
}

// Decided rows link to the expense, which is no longer committed by then.
func TestExpenseUncommitAdminCanOpenExpensesInTheHistory(t *testing.T) {
	app := testutils.SetupTestApp(t)
	adminToken := uncommitToken(t, uncommitAdminEmail)

	t.Run("after approving a request", func(t *testing.T) {
		commitExpenseToday(t, app, uncommitOnAccountExpense)
		id := createUncommitRequest(t, app, uncommitOnAccountExpense, uncommitReasonWrongCurr, "uncommit_and_reject", "")
		if results := approveUncommitRequests(t, app, id); !results[0].OK {
			t.Fatalf("expected approval, got %+v", results[0])
		}
		res := performTestAPIRequest(t, app, http.MethodGet, "/api/expenses/details/"+uncommitOnAccountExpense, nil, adminToken)
		mustStatus(t, res, http.StatusOK)
	})

	t.Run("after a direct uncommit", func(t *testing.T) {
		res := performTestAPIRequest(t, app, http.MethodPost, "/api/expenses/xg2yeucklhgbs3n/uncommit", nil, adminToken)
		mustStatus(t, res, http.StatusOK)
		res = performTestAPIRequest(t, app, http.MethodGet, "/api/expenses/details/xg2yeucklhgbs3n", nil, adminToken)
		mustStatus(t, res, http.StatusOK)
	})

	t.Run("not for other uncommitted expenses", func(t *testing.T) {
		res := performTestAPIRequest(t, app, http.MethodGet, "/api/expenses/details/eqhozipupteogp8", nil, adminToken)
		mustStatus(t, res, http.StatusNotFound)
	})
}
