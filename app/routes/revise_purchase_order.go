// This file implements revising the total of an Active purchase order. The
// creator submits a new total, and the purchase order returns to Unapproved
// in the same transaction. It keeps its po_number, so approving it again
// re-activates the same number.

package routes

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"tybalt/errs"
	"tybalt/hooks"
	"tybalt/notifications"

	"github.com/pocketbase/dbx"
	validation "github.com/pocketbase/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

// A nil approver keeps the assigned approver. The new total may need
// different approvers, so the client can submit them with the total.
type revisePurchaseOrderRequest struct {
	Total                  float64 `json:"total"`
	Approver               *string `json:"approver"`
	PrioritySecondApprover *string `json:"priority_second_approver"`
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

			if req.Total == po.GetFloat("total") {
				httpResponseStatusCode = http.StatusBadRequest
				return &CodeError{
					Code:    "total_unchanged",
					Message: "the revised total must differ from the current total",
				}
			}

			po.Set("total", req.Total)
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

			// Recompute approval_total for the new total and check the approvers
			// against the policy for it. CleanPurchaseOrder also clears any
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
			if _, err := notifications.DispatchNotification(app, hooks.RevisionApprovalNotification(app, revisedPO, userId, notifications.DeliveryImmediate)); err != nil {
				// Log the error but don't fail the request, as the PO was already revised
				app.Logger().Error("notification not sent: error creating revision approval notification", "error", err)
			}
		}

		return e.JSON(http.StatusOK, revisedPO)
	}
}
