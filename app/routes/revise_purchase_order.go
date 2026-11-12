// This file implements revising the total, currency, job, description,
// division or category of an Active purchase order. The creator submits them,
// and the purchase order returns to Unapproved in the same transaction. It
// keeps its po_number, so approving it again re-activates the same number.

package routes

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"tybalt/errs"
	"tybalt/hooks"
	"tybalt/notifications"
	"tybalt/utilities"

	"github.com/pocketbase/dbx"
	validation "github.com/pocketbase/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

// A nil field keeps the current value. The new total, job or division may need
// different approvers, so the client can submit them with it. The budget
// confirmation belongs to the job, so it can change only with a new job. A new
// job clears the category and the budget confirmation unless the request sets
// them again. The attachment can't change here, because this route takes JSON;
// the creator replaces it while the revision awaits approval.
type revisePurchaseOrderRequest struct {
	Total                      float64 `json:"total"`
	Currency                   *string `json:"currency"`
	Approver                   *string `json:"approver"`
	PrioritySecondApprover     *string `json:"priority_second_approver"`
	Job                        *string `json:"job"`
	Description                *string `json:"description"`
	Division                   *string `json:"division"`
	Category                   *string `json:"category"`
	CoveredWithinProjectBudget *bool   `json:"covered_within_project_budget"`
}

func createRevisePurchaseOrderHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		id := e.Request.PathValue("id")
		userId := e.Auth.Id

		var req revisePurchaseOrderRequest
		if err := e.BindBody(&req); err != nil {
			return e.JSON(http.StatusBadRequest, map[string]any{
				"message": "invalid request body",
				"code":    "invalid_request_body",
			})
		}

		var httpResponseStatusCode int
		var previous *core.Record
		err := app.RunInTransaction(func(txApp core.App) error {
			po, err := txApp.FindRecordById("purchase_orders", id)
			if err != nil {
				httpResponseStatusCode = http.StatusNotFound
				return &CodeError{
					Code:    "po_not_found",
					Message: fmt.Sprintf("error fetching purchase order: %v", err),
				}
			}

			if po.GetString("status") != "Active" {
				httpResponseStatusCode = http.StatusBadRequest
				return &CodeError{
					Code:    "po_not_active",
					Message: "only active purchase orders can be revised",
				}
			}

			if po.GetString("uid") != userId {
				httpResponseStatusCode = http.StatusForbidden
				return &CodeError{
					Code:    "unauthorized_revision",
					Message: "only the creator can revise this purchase order",
				}
			}

			// Legacy purchase orders are edited through the legacy routes and
			// imported ones are owned by the legacy system. Neither went through
			// Turbo approval, so neither can return to it.
			if po.GetBool("legacy_manual_entry") || po.GetBool("_imported") {
				httpResponseStatusCode = http.StatusBadRequest
				return &CodeError{
					Code:    "legacy_po_not_revisable",
					Message: "legacy purchase orders cannot be revised",
				}
			}

			// An approval only matters once an expense uses it. Without expenses,
			// replacing the approval loses nothing.
			expenses, err := txApp.FindRecordsByFilter("expenses", "purchase_order = {:poId}", "", 1, 0, dbx.Params{
				"poId": po.Id,
			})
			if err != nil {
				httpResponseStatusCode = http.StatusInternalServerError
				return &CodeError{
					Code:    "error_fetching_expenses",
					Message: fmt.Sprintf("error fetching expenses: %v", err),
				}
			}
			if len(expenses) > 0 {
				httpResponseStatusCode = http.StatusBadRequest
				return &CodeError{
					Code:    "po_has_expenses",
					Message: "this purchase order has associated expenses and cannot be revised",
				}
			}

			children, err := txApp.FindRecordsByFilter("purchase_orders", "parent_po = {:poId} && status != 'Cancelled'", "", 1, 0, dbx.Params{
				"poId": po.Id,
			})
			if err != nil {
				httpResponseStatusCode = http.StatusInternalServerError
				return &CodeError{
					Code:    "error_fetching_child_pos",
					Message: fmt.Sprintf("error fetching child purchase orders: %v", err),
				}
			}
			if len(children) > 0 {
				httpResponseStatusCode = http.StatusBadRequest
				return &CodeError{
					Code:    "po_has_child_pos",
					Message: "this purchase order has child purchase orders and cannot be revised",
				}
			}

			previous = po.Original()
			previousJob := strings.TrimSpace(po.GetString("job"))
			job := previousJob
			if req.Job != nil {
				job = strings.TrimSpace(*req.Job)
			}
			if code, message := hooks.RevisionJobChangeRefusal(po, job); code != "" {
				httpResponseStatusCode = http.StatusBadRequest
				return &CodeError{Code: code, Message: message}
			}
			jobChanged := job != previousJob
			if !jobChanged && req.CoveredWithinProjectBudget != nil && *req.CoveredWithinProjectBudget != po.GetBool("covered_within_project_budget") {
				httpResponseStatusCode = http.StatusBadRequest
				return &CodeError{
					Code:    "revision_requires_job_change",
					Message: "the budget confirmation can change only with a new job",
				}
			}

			description := po.GetString("description")
			if req.Description != nil {
				description = *req.Description
			}
			division := po.GetString("division")
			if req.Division != nil {
				division = strings.TrimSpace(*req.Division)
			}
			// A new job clears the category unless the request sets one.
			category := po.GetString("category")
			if jobChanged {
				category = ""
			}
			if req.Category != nil {
				category = strings.TrimSpace(*req.Category)
			}
			categoryChanged := category != po.GetString("category")

			currency := po.GetString("currency")
			if req.Currency != nil {
				currency = strings.TrimSpace(*req.Currency)
			}
			currencyChanged := utilities.EffectiveCurrencyCode(txApp, currency) != utilities.EffectiveCurrencyCode(txApp, po.GetString("currency"))
			if req.Total == po.GetFloat("total") && !currencyChanged && !jobChanged &&
				description == po.GetString("description") && division == po.GetString("division") && !categoryChanged {
				httpResponseStatusCode = http.StatusBadRequest
				return &CodeError{
					Code:    "revision_unchanged",
					Message: "a revision must change the total, currency, job, description, division or category",
				}
			}

			// The branch follows the job in CleanPurchaseOrder. po was loaded
			// with the old job and division, so when either changes
			// ValidatePurchaseOrder checks that the division is allocated to the
			// job, and it checks the approvers against the new division.
			if jobChanged {
				po.Set("job", job)
				po.Set("covered_within_project_budget", req.CoveredWithinProjectBudget != nil && *req.CoveredWithinProjectBudget)
			}
			po.Set("description", description)
			po.Set("division", division)
			po.Set("category", category)
			po.Set("total", req.Total)
			po.Set("currency", currency)
			if req.Approver != nil {
				po.Set("approver", strings.TrimSpace(*req.Approver))
			}
			if req.PrioritySecondApprover != nil {
				po.Set("priority_second_approver", strings.TrimSpace(*req.PrioritySecondApprover))
			}
			po.Set("status", "Unapproved")
			po.Set("approved", "")
			po.Set("second_approval", "")
			po.Set("second_approver", "")

			// Recompute the approval totals and check the approvers against the
			// policy for them. CleanPurchaseOrder also clears any
			// rejection fields.
			if err := hooks.CleanPurchaseOrder(txApp, po); err != nil {
				return err
			}
			if err := hooks.EnsureActiveDivision(txApp, po.GetString("division"), "division"); err != nil {
				return err
			}
			if err := hooks.ValidatePurchaseOrder(txApp, po, false); err != nil {
				return err
			}
			if jobChanged || categoryChanged {
				if err := hooks.ValidatePurchaseOrderCategory(txApp, po); err != nil {
					return err
				}
			}

			if err := txApp.Save(po); err != nil {
				httpResponseStatusCode = http.StatusInternalServerError
				return &CodeError{
					Code:    "error_updating_purchase_order",
					Message: fmt.Sprintf("error updating purchase order: %v", err),
				}
			}

			return nil
		})

		if err != nil {
			var codeError *CodeError
			var hookErr *errs.HookError
			var validationErrs validation.Errors
			switch {
			case errors.As(err, &codeError):
				return e.JSON(httpResponseStatusCode, map[string]any{
					"message": codeError.Message,
					"code":    codeError.Code,
				})
			case errors.As(err, &hookErr):
				return e.JSON(hookErr.Status, hookErr)
			case errors.As(err, &validationErrs):
				return apis.NewBadRequestError("Validation error", validationErrs)
			}
			return e.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		revisedPO, err := app.FindRecordById("purchase_orders", id)
		if err != nil {
			return e.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		if strings.TrimSpace(revisedPO.GetString("approver")) != "" {
			if _, err := notifications.DispatchNotification(app, hooks.RevisionApprovalNotification(app, revisedPO, previous, "", userId, notifications.DeliveryImmediate)); err != nil {
				// Log the error but don't fail the request, as the PO was already revised
				app.Logger().Error("notification not sent: error creating revision approval notification", "error", err)
			}
		}

		return e.JSON(http.StatusOK, revisedPO)
	}
}
