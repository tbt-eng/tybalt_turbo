package routes

import (
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
	"tybalt/notifications"
	"tybalt/utilities"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

//go:embed expense_uncommit_requests.sql
var expenseUncommitRequestsQuery string

const (
	expenseUncommitRequestsCollection = "expense_uncommit_requests"
	expenseUncommitReasonsCollection  = "expense_uncommit_reasons"

	expenseUncommitOutcomeReject = "uncommit_and_reject"
	expenseUncommitOutcomeOnly   = "uncommit_only"

	expenseUncommitStatusPending   = "pending"
	expenseUncommitStatusApproved  = "approved"
	expenseUncommitStatusDeclined  = "declined"
	expenseUncommitStatusCancelled = "cancelled"

	expenseUncommitDecidedListLimit = 200
	expenseUncommitNoteMaxLength    = 1000
	expenseUncommitMinCommentLength = 4
)

type expenseUncommitReasonRow struct {
	ID             string `db:"id" json:"id"`
	Name           string `db:"name" json:"name"`
	DefaultOutcome string `db:"default_outcome" json:"default_outcome"`
	NotePrompt     string `db:"note_prompt" json:"note_prompt"`
	NoteRequired   bool   `db:"note_required" json:"note_required"`
}

type expenseUncommitRequestRow struct {
	ID              string  `db:"id" json:"id"`
	Expense         string  `db:"expense" json:"expense"`
	Status          string  `db:"status" json:"status"`
	Outcome         string  `db:"outcome" json:"outcome"`
	Note            string  `db:"note" json:"note"`
	Created         string  `db:"created" json:"created"`
	Requester       string  `db:"requester" json:"requester"`
	DecidedAt       string  `db:"decided_at" json:"decided_at"`
	DecisionComment string  `db:"decision_comment" json:"decision_comment"`
	EmployeeName    string  `db:"employee_name" json:"employee_name"`
	PONumber        string  `db:"po_number" json:"po_number"`
	VendorName      string  `db:"vendor_name" json:"vendor_name"`
	Description     string  `db:"description" json:"description"`
	CurrencyCode    string  `db:"currency_code" json:"currency_code"`
	PaymentType     string  `db:"payment_type" json:"payment_type"`
	Total           float64 `db:"total" json:"total"`
	Committed       string  `db:"committed" json:"committed"`
	PayPeriodEnding string  `db:"pay_period_ending" json:"pay_period_ending"`
	ReasonName      string  `db:"reason_name" json:"reason_name"`
	RequesterName   string  `db:"requester_name" json:"requester_name"`
	DecidedByName   string  `db:"decided_by_name" json:"decided_by_name"`
	// Set on pending rows whose gate closed while they waited.
	Blocked *utilities.ExpenseUncommitBlock `db:"-" json:"blocked,omitempty"`
}

// expenseUncommitRequestState is what the expense details page needs to offer,
// explain or show a request.
type expenseUncommitRequestState struct {
	CanRequest bool                            `json:"can_request"`
	Blocked    *utilities.ExpenseUncommitBlock `json:"blocked,omitempty"`
	Pending    *expenseUncommitRequestRow      `json:"pending,omitempty"`
}

func expenseUncommitError(status int, code string, message string) *expenseUncommitFailure {
	return &expenseUncommitFailure{status: status, code: code, message: message}
}

// expenseUncommitFailure carries the status a route should answer with out of
// a transaction.
type expenseUncommitFailure struct {
	status  int
	code    string
	message string
}

func (f *expenseUncommitFailure) Error() string {
	return f.message
}

func respondExpenseUncommitError(e *core.RequestEvent, err error) error {
	var failure *expenseUncommitFailure
	if errors.As(err, &failure) {
		return e.JSON(failure.status, map[string]string{"code": failure.code, "message": failure.message})
	}
	return e.Error(http.StatusInternalServerError, "failed to update uncommit request", err)
}

// callerHoldsAnyClaim answers true when the caller holds any of claims.
func callerHoldsAnyClaim(app core.App, auth *core.Record, claims ...string) (bool, error) {
	for _, claim := range claims {
		has, err := utilities.HasClaim(app, auth, claim)
		if err != nil {
			return false, err
		}
		if has {
			return true, nil
		}
	}
	return false, nil
}

func newExpenseUncommitForbidden(action string) *expenseUncommitFailure {
	return expenseUncommitError(http.StatusForbidden, "unauthorized", fmt.Sprintf("you are not authorized to %s", action))
}

// findPendingExpenseUncommitRequest returns the expense's pending request, or
// nil when it has none.
func findPendingExpenseUncommitRequest(app core.App, expenseID string) (*core.Record, error) {
	record, err := app.FindFirstRecordByFilter(
		expenseUncommitRequestsCollection,
		"expense = {:expense} && status = {:status}",
		dbx.Params{"expense": expenseID, "status": expenseUncommitStatusPending},
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return record, err
}

// expenseUncommitBlockForRecord applies the eligibility rules that depend only
// on the expense itself.
func expenseUncommitBlockForRecord(app core.App, expense *core.Record) *utilities.ExpenseUncommitBlock {
	return utilities.FindExpenseUncommitBlock(
		app,
		expense.GetString("payment_type"),
		expense.GetString("pay_period_ending"),
		expense.GetDateTime("committed"),
		time.Now(),
	)
}

// setExpenseUncommitSnapshot copies what the queue and history show onto the
// request, because uncommit clears the commit fields and the expense can later
// be deleted.
func setExpenseUncommitSnapshot(app core.App, request *core.Record, expense *core.Record) error {
	var names struct {
		EmployeeName string `db:"employee_name"`
		PONumber     string `db:"po_number"`
		VendorName   string `db:"vendor_name"`
	}
	err := app.DB().NewQuery(`
		SELECT
			COALESCE((SELECT TRIM(p.given_name || ' ' || p.surname) FROM profiles p WHERE p.uid = e.uid), '') AS employee_name,
			COALESCE((SELECT po.po_number FROM purchase_orders po WHERE po.id = e.purchase_order), '') AS po_number,
			COALESCE((SELECT v.name FROM vendors v WHERE v.id = e.vendor), '') AS vendor_name
		FROM expenses e
		WHERE e.id = {:id}
	`).Bind(dbx.Params{"id": expense.Id}).One(&names)
	if err != nil {
		return err
	}

	request.Set("employee_name", names.EmployeeName)
	request.Set("po_number", names.PONumber)
	request.Set("vendor_name", names.VendorName)
	request.Set("description", expense.GetString("description"))
	request.Set("currency_code", utilities.EffectiveCurrencyCode(app, expense.GetString("currency")))
	request.Set("payment_type", expense.GetString("payment_type"))
	request.Set("total", expense.GetFloat("total"))
	request.Set("committed", expense.GetString("committed"))
	request.Set("committer", expense.GetString("committer"))
	request.Set("committed_week_ending", expense.GetString("committed_week_ending"))
	request.Set("pay_period_ending", expense.GetString("pay_period_ending"))
	return nil
}

func findExpenseUncommitRequestRows(app core.App, where string, orderBy string, params dbx.Params) ([]expenseUncommitRequestRow, error) {
	rows := []expenseUncommitRequestRow{}
	query := expenseUncommitRequestsQuery + "\nWHERE " + where + "\nORDER BY " + orderBy
	err := app.DB().NewQuery(query).Bind(params).All(&rows)
	return rows, err
}

// loadExpenseUncommitRequestState returns the request state shown on a
// committed expense's details page, or nil when the caller has no part in
// uncommit requests or the expense isn't committed.
func loadExpenseUncommitRequestState(app core.App, auth *core.Record, expenseID string) (*expenseUncommitRequestState, error) {
	hasCommit, err := utilities.HasClaim(app, auth, "commit")
	if err != nil {
		return nil, err
	}
	hasAdmin, err := utilities.HasClaim(app, auth, "admin")
	if err != nil {
		return nil, err
	}
	if !hasCommit && !hasAdmin {
		return nil, nil
	}

	expense, err := app.FindRecordById("expenses", expenseID)
	if err != nil {
		return nil, err
	}
	if expense.GetDateTime("committed").IsZero() {
		return nil, nil
	}

	state := &expenseUncommitRequestState{}

	rows, err := findExpenseUncommitRequestRows(app, "r.expense = {:expense} AND r.status = {:status}", "r.created", dbx.Params{
		"expense": expenseID,
		"status":  expenseUncommitStatusPending,
	})
	if err != nil {
		return nil, err
	}
	if len(rows) > 0 {
		state.Pending = &rows[0]
	}

	if !expense.GetBool("_imported") {
		state.Blocked = expenseUncommitBlockForRecord(app, expense)
	}
	state.CanRequest = hasCommit && state.Pending == nil && state.Blocked == nil && !expense.GetBool("_imported")
	return state, nil
}

func createGetExpenseUncommitReasonsHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		allowed, err := callerHoldsAnyClaim(app, e.Auth, "commit", "admin")
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to check claims", err)
		}
		if !allowed {
			return respondExpenseUncommitError(e, newExpenseUncommitForbidden("view uncommit reasons"))
		}

		rows := []expenseUncommitReasonRow{}
		if err := app.DB().NewQuery(`
			SELECT id, name, default_outcome, note_prompt, note_required
			FROM expense_uncommit_reasons
			WHERE active = 1
			ORDER BY sort, name
		`).All(&rows); err != nil {
			return e.Error(http.StatusInternalServerError, "failed to load uncommit reasons", err)
		}
		return e.JSON(http.StatusOK, rows)
	}
}

func createGetExpenseUncommitRequestsHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		allowed, err := callerHoldsAnyClaim(app, e.Auth, "commit", "admin")
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to check claims", err)
		}
		if !allowed {
			return respondExpenseUncommitError(e, newExpenseUncommitForbidden("view uncommit requests"))
		}

		var rows []expenseUncommitRequestRow
		status := e.Request.URL.Query().Get("status")
		switch status {
		case "", expenseUncommitStatusPending:
			rows, err = findExpenseUncommitRequestRows(app, "r.status = {:status}", "r.created", dbx.Params{"status": expenseUncommitStatusPending})
		case "decided":
			rows, err = findExpenseUncommitRequestRows(
				app,
				"r.status != {:status}",
				fmt.Sprintf("r.decided_at DESC LIMIT %d", expenseUncommitDecidedListLimit),
				dbx.Params{"status": expenseUncommitStatusPending},
			)
		default:
			return respondExpenseUncommitError(e, expenseUncommitError(http.StatusBadRequest, "invalid_status", "status must be pending or decided"))
		}
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to load uncommit requests", err)
		}

		now := time.Now()
		for i := range rows {
			if rows[i].Status != expenseUncommitStatusPending {
				continue
			}
			committed, err := types.ParseDateTime(rows[i].Committed)
			if err != nil {
				return e.Error(http.StatusInternalServerError, "failed to read a committed date", err)
			}
			rows[i].Blocked = utilities.FindExpenseUncommitBlock(app, rows[i].PaymentType, rows[i].PayPeriodEnding, committed, now)
		}
		return e.JSON(http.StatusOK, rows)
	}
}

func createCreateExpenseUncommitRequestHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		if err := requireExpensesEditing(app, "expenses"); err != nil {
			return err
		}

		allowed, err := callerHoldsAnyClaim(app, e.Auth, "commit")
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to check claims", err)
		}
		if !allowed {
			return respondExpenseUncommitError(e, newExpenseUncommitForbidden("request an uncommit"))
		}

		var body struct {
			Reason  string `json:"reason"`
			Outcome string `json:"outcome"`
			Note    string `json:"note"`
		}
		if err := e.BindBody(&body); err != nil {
			return respondExpenseUncommitError(e, expenseUncommitError(http.StatusBadRequest, "invalid_request_body", "invalid request body"))
		}
		note := strings.TrimSpace(body.Note)

		var requestID string
		err = app.RunInTransaction(func(txApp core.App) error {
			expense, err := txApp.FindRecordById("expenses", e.Request.PathValue("id"))
			if err != nil {
				return expenseUncommitError(http.StatusNotFound, "record_not_found", "expense not found")
			}
			if expense.GetDateTime("committed").IsZero() {
				return expenseUncommitError(http.StatusBadRequest, "record_not_committed", "this expense is not committed")
			}
			if expense.GetBool("_imported") {
				return expenseUncommitError(http.StatusBadRequest, "record_imported", "imported expenses can't be requested")
			}

			pending, err := findPendingExpenseUncommitRequest(txApp, expense.Id)
			if err != nil {
				return err
			}
			if pending != nil {
				return expenseUncommitError(http.StatusConflict, "request_pending", "this expense already has a pending uncommit request")
			}

			if block := expenseUncommitBlockForRecord(txApp, expense); block != nil {
				return expenseUncommitError(http.StatusBadRequest, block.Code, block.Message)
			}

			reason, err := txApp.FindRecordById(expenseUncommitReasonsCollection, body.Reason)
			if err != nil || !reason.GetBool("active") {
				return expenseUncommitError(http.StatusBadRequest, "invalid_reason", "choose a reason")
			}
			if body.Outcome != expenseUncommitOutcomeReject && body.Outcome != expenseUncommitOutcomeOnly {
				return expenseUncommitError(http.StatusBadRequest, "invalid_outcome", "choose whether the expense should also be rejected")
			}
			if note == "" && reason.GetBool("note_required") {
				return expenseUncommitError(http.StatusBadRequest, "note_required", "this reason needs a note")
			}
			if len(note) > expenseUncommitNoteMaxLength {
				return expenseUncommitError(http.StatusBadRequest, "note_too_long", "the note is too long")
			}

			collection, err := txApp.FindCollectionByNameOrId(expenseUncommitRequestsCollection)
			if err != nil {
				return err
			}
			request := core.NewRecord(collection)
			request.Set("expense", expense.Id)
			request.Set("requester", e.Auth.Id)
			request.Set("reason", reason.Id)
			request.Set("note", note)
			request.Set("outcome", body.Outcome)
			request.Set("status", expenseUncommitStatusPending)
			if err := setExpenseUncommitSnapshot(txApp, request, expense); err != nil {
				return err
			}
			if err := txApp.Save(request); err != nil {
				return err
			}
			requestID = request.Id
			return nil
		})
		if err != nil {
			return respondExpenseUncommitError(e, err)
		}

		return e.JSON(http.StatusOK, map[string]string{"id": requestID})
	}
}

func createCancelExpenseUncommitRequestHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		if err := requireExpensesEditing(app, "expenses"); err != nil {
			return err
		}

		err := app.RunInTransaction(func(txApp core.App) error {
			request, err := txApp.FindRecordById(expenseUncommitRequestsCollection, e.Request.PathValue("id"))
			if err != nil {
				return expenseUncommitError(http.StatusNotFound, "record_not_found", "uncommit request not found")
			}
			if request.GetString("requester") != e.Auth.Id {
				return newExpenseUncommitForbidden("cancel this uncommit request")
			}
			if request.GetString("status") != expenseUncommitStatusPending {
				return expenseUncommitError(http.StatusBadRequest, "request_not_pending", "this uncommit request has already been decided")
			}
			request.Set("status", expenseUncommitStatusCancelled)
			request.Set("decided_by", e.Auth.Id)
			request.Set("decided_at", time.Now())
			return txApp.Save(request)
		})
		if err != nil {
			return respondExpenseUncommitError(e, err)
		}
		return e.JSON(http.StatusOK, map[string]string{"message": "uncommit request cancelled"})
	}
}

func createDeclineExpenseUncommitRequestHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		if err := requireExpensesEditing(app, "expenses"); err != nil {
			return err
		}

		allowed, err := callerHoldsAnyClaim(app, e.Auth, "admin")
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to check claims", err)
		}
		if !allowed {
			return respondExpenseUncommitError(e, newExpenseUncommitForbidden("decline this uncommit request"))
		}

		var body struct {
			Comment string `json:"comment"`
		}
		if err := e.BindBody(&body); err != nil {
			return respondExpenseUncommitError(e, expenseUncommitError(http.StatusBadRequest, "invalid_request_body", "invalid request body"))
		}
		comment := strings.TrimSpace(body.Comment)

		var request *core.Record
		err = app.RunInTransaction(func(txApp core.App) error {
			var err error
			request, err = txApp.FindRecordById(expenseUncommitRequestsCollection, e.Request.PathValue("id"))
			if err != nil {
				return expenseUncommitError(http.StatusNotFound, "record_not_found", "uncommit request not found")
			}
			if request.GetString("status") != expenseUncommitStatusPending {
				return expenseUncommitError(http.StatusBadRequest, "request_not_pending", "this uncommit request has already been decided")
			}
			if len(comment) < expenseUncommitMinCommentLength {
				return expenseUncommitError(http.StatusBadRequest, "comment_too_short", "say why the request is being declined")
			}
			if len(comment) > expenseUncommitNoteMaxLength {
				return expenseUncommitError(http.StatusBadRequest, "comment_too_long", "the comment is too long")
			}
			request.Set("status", expenseUncommitStatusDeclined)
			request.Set("decided_by", e.Auth.Id)
			request.Set("decided_at", time.Now())
			request.Set("decision_comment", comment)
			return txApp.Save(request)
		})
		if err != nil {
			return respondExpenseUncommitError(e, err)
		}

		queueExpenseUncommitDecision(app, request)
		return e.JSON(http.StatusOK, map[string]string{"message": "uncommit request declined"})
	}
}

// expenseUncommitApprovalResult is one request's outcome in a bulk approval.
type expenseUncommitApprovalResult struct {
	ID      string `json:"id"`
	OK      bool   `json:"ok"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// createApproveExpenseUncommitRequestsHandler approves one or many requests.
// Each runs in its own transaction, so one failure doesn't block the rest, and
// the response reports each. It answers 200 even when some fail.
func createApproveExpenseUncommitRequestsHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		if err := requireExpensesEditing(app, "expenses"); err != nil {
			return err
		}

		allowed, err := callerHoldsAnyClaim(app, e.Auth, "admin")
		if err != nil {
			return e.Error(http.StatusInternalServerError, "failed to check claims", err)
		}
		if !allowed {
			return respondExpenseUncommitError(e, newExpenseUncommitForbidden("approve uncommit requests"))
		}

		var body struct {
			IDs []string `json:"ids"`
		}
		if err := e.BindBody(&body); err != nil || len(body.IDs) == 0 {
			return respondExpenseUncommitError(e, expenseUncommitError(http.StatusBadRequest, "invalid_request_body", "ids is required"))
		}

		results := make([]expenseUncommitApprovalResult, 0, len(body.IDs))
		for _, id := range body.IDs {
			result := expenseUncommitApprovalResult{ID: id, OK: true}
			if err := approveExpenseUncommitRequest(app, id, e.Auth.Id); err != nil {
				result.OK = false
				result.Message = err.Error()
				var failure *expenseUncommitFailure
				if errors.As(err, &failure) {
					result.Code = failure.code
				} else {
					result.Code = "approval_failed"
					app.Logger().Error("error approving expense uncommit request", "request_id", id, "error", err)
				}
			}
			results = append(results, result)
		}

		return e.JSON(http.StatusOK, map[string]any{"results": results})
	}
}

func approveExpenseUncommitRequest(app core.App, requestID string, deciderID string) error {
	var request, expense *core.Record
	var rejectionReason string

	err := app.RunInTransaction(func(txApp core.App) error {
		var err error
		request, err = txApp.FindRecordById(expenseUncommitRequestsCollection, requestID)
		if err != nil {
			return expenseUncommitError(http.StatusNotFound, "record_not_found", "uncommit request not found")
		}
		if request.GetString("status") != expenseUncommitStatusPending {
			return expenseUncommitError(http.StatusBadRequest, "request_not_pending", "this uncommit request has already been decided")
		}

		expense, err = txApp.FindRecordById("expenses", request.GetString("expense"))
		if err != nil {
			return expenseUncommitError(http.StatusNotFound, "record_not_found", "the expense no longer exists")
		}
		if expense.GetDateTime("committed").IsZero() {
			return expenseUncommitError(http.StatusBadRequest, "record_not_committed", "this expense is not committed")
		}
		if expense.GetBool("_imported") {
			return expenseUncommitError(http.StatusBadRequest, "record_imported", "imported expenses can't be requested")
		}
		if block := expenseUncommitBlockForRecord(txApp, expense); block != nil {
			return expenseUncommitError(http.StatusBadRequest, block.Code, block.Message)
		}

		rejectionReason, _, err = applyExpenseUncommitRequest(txApp, request, expense, deciderID)
		return err
	})
	if err != nil {
		return err
	}

	queueExpenseUncommitApproval(app, request, expense, rejectionReason)
	return nil
}

// applyExpenseUncommitRequest does what a request asks for and marks it
// approved: it uncommits the expense and, for uncommit_and_reject, rejects it
// too. Both approving from the queue and the direct uncommit button come here,
// so a pending request gets the same result either way. It returns the rejection
// reason, blank if there was no rejection, and on failure the status to answer
// with.
func applyExpenseUncommitRequest(txApp core.App, request *core.Record, expense *core.Record, deciderID string) (string, int, error) {
	if err := setExpenseUncommitSnapshot(txApp, request, expense); err != nil {
		return "", http.StatusInternalServerError, &CodeError{
			Code:    "error_recording_uncommit",
			Message: fmt.Sprintf("error recording uncommit: %v", err),
		}
	}

	if status, err := uncommitRecordInTransaction(txApp, expense); err != nil {
		return "", status, err
	}

	rejectionReason := ""
	if request.GetString("outcome") == expenseUncommitOutcomeReject {
		rejectionReason = expenseUncommitRejectionReason(txApp, request)
		// The requester holds the commit claim and made this call. The admin
		// holder approving it may hold neither reject permission.
		if status, err := rejectRecordInTransaction(txApp, expense, rejectionReason, request.GetString("requester")); err != nil {
			return "", status, err
		}
	}

	request.Set("status", expenseUncommitStatusApproved)
	request.Set("decided_by", deciderID)
	request.Set("decided_at", time.Now())
	if err := txApp.Save(request); err != nil {
		return "", http.StatusInternalServerError, &CodeError{
			Code:    "error_recording_uncommit",
			Message: fmt.Sprintf("error recording uncommit: %v", err),
		}
	}
	return rejectionReason, 0, nil
}

// queueExpenseUncommitApproval tells the requester, and the employee when the
// expense was rejected. A rejection reason is blank when there was no rejection.
func queueExpenseUncommitApproval(app core.App, request *core.Record, expense *core.Record, rejectionReason string) {
	queueExpenseUncommitDecision(app, request)
	if rejectionReason == "" {
		return
	}
	if err := notifications.QueueExpenseRejectedNotifications(app, expense, request.GetString("requester"), rejectionReason); err != nil {
		app.Logger().Error("error queueing expense rejection notifications", "expense_id", expense.Id, "error", err)
	}
}

func expenseUncommitReasonName(app core.App, request *core.Record) string {
	reason, err := app.FindRecordById(expenseUncommitReasonsCollection, request.GetString("reason"))
	if err != nil {
		return ""
	}
	return reason.GetString("name")
}

func expenseUncommitRejectionReason(app core.App, request *core.Record) string {
	name := expenseUncommitReasonName(app, request)
	note := request.GetString("note")
	switch {
	case name != "" && note != "":
		return name + ": " + note
	case name != "":
		return name
	case note != "":
		return note
	}
	// The reason was deleted while the request waited and there is no note.
	return "Uncommit requested"
}

func queueExpenseUncommitDecision(app core.App, request *core.Record) {
	if err := notifications.QueueExpenseUncommitDecisionNotification(app, request, expenseUncommitReasonName(app, request)); err != nil {
		app.Logger().Error("error queueing expense uncommit decision notification", "request_id", request.Id, "error", err)
	}
}

// findOrNewDirectExpenseUncommitRequest returns the request a direct uncommit
// should apply: the expense's pending request, so its outcome is honoured, or
// otherwise a new row recording an uncommit_only by the caller. A new row is
// saved when it is applied.
func findOrNewDirectExpenseUncommitRequest(txApp core.App, expense *core.Record, uid string) (*core.Record, error) {
	request, err := findPendingExpenseUncommitRequest(txApp, expense.Id)
	if err != nil || request != nil {
		return request, err
	}

	collection, err := txApp.FindCollectionByNameOrId(expenseUncommitRequestsCollection)
	if err != nil {
		return nil, err
	}
	request = core.NewRecord(collection)
	request.Set("expense", expense.Id)
	request.Set("requester", uid)
	request.Set("outcome", expenseUncommitOutcomeOnly)
	return request, nil
}
