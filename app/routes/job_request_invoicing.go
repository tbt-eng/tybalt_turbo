package routes

import (
	"context"
	"slices"

	"tybalt/utilities"

	"github.com/pocketbase/pocketbase/core"
)

// requestInvoicingWrite sets the invoicing profile of a request. An empty
// profile clears it.
type requestInvoicingWrite struct {
	Version              int    `json:"version"`
	InvoicingInformation string `json:"invoicing_information"`
}

// canSetRequestInvoicing reports whether auth may set the invoicing profile of
// a request that waits for job creation: a holder of the job claim, while
// project requests and job editing are both on.
func canSetRequestInvoicing(app core.App, r, auth *core.Record) (bool, error) {
	if auth == nil || r.GetString("status") != "Awaiting creation" {
		return false, nil
	}
	enabled, err := utilities.IsJobRequestApprovalRequired(app)
	if err != nil || !enabled {
		return false, err
	}
	editing, err := utilities.IsJobsEditingEnabled(app)
	if err != nil || !editing {
		return false, err
	}
	return utilities.HasClaim(app, auth, "job")
}

// requestInvoicingHandler lets procurement set, change or clear the invoicing
// profile of a request in the Job Create Queue, on anyone's request. Only the
// profile changes. The status, the version and the approval stay: the review
// digest leaves the profile out (requestReviewDigest), so the approval still
// covers the request. The profile must belong to the request's client, and
// its billing contact needs an email address or the no-email mark, as job
// creation checks.
func requestInvoicingHandler(app core.App) func(*core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		var input requestInvoicingWrite
		if err := e.BindBody(&input); err != nil {
			return workflowError("Invalid request data.")
		}
		err := app.RunInTransaction(func(tx core.App) error {
			if err := requestEnabled(tx); err != nil {
				return err
			}
			if err := workflowClaim(tx, e.Auth, "job"); err != nil {
				return err
			}
			// The profile becomes the job's, so the change follows the jobs
			// editing switch, as the other job writes do.
			editing, err := utilities.IsJobsEditingEnabled(tx)
			if err != nil {
				return err
			}
			if !editing {
				return utilities.ErrJobsEditingDisabled
			}
			r, err := loadVisibleRequest(tx, e.Request.PathValue("id"), e.Auth)
			if err != nil {
				return err
			}
			if r.GetString("status") != "Awaiting creation" {
				return workflowError("Procurement can change the invoicing profile only while the request is awaiting job creation.")
			}
			if err = requestVersion(r, input.Version); err != nil {
				return err
			}
			if input.InvoicingInformation == r.GetString("invoicing_information") {
				return nil
			}
			// An approval saved before the profile left the digest holds the
			// profile. Give it the current digest first, so the change does not
			// make it stale.
			if err = upgradeLegacyRequestApproval(tx, r); err != nil {
				return err
			}
			r.Set("invoicing_information", input.InvoicingInformation)
			if err = utilities.ValidateClientReferences(tx, r); err != nil {
				return err
			}
			if err = utilities.ValidateRequestContactEmails(tx, r); err != nil {
				return err
			}
			return tx.SaveWithContext(workflowContext(e.Request.Context()), r)
		})
		if err != nil {
			return err
		}
		return respondRequestDetail(app, e, e.Request.PathValue("id"))
	}
}

// upgradeLegacyRequestApproval gives the request's saved approval the current
// form of the review digest, when it holds an older form that still matches
// the request's data. An older form holds more than the current one, so the
// match proves that nothing the current digest covers has changed since the
// approval. A saved digest that matches no form is left alone, so a stale
// approval stays stale.
func upgradeLegacyRequestApproval(tx core.App, r *core.Record) error {
	a, err := requestAttachment(tx, r.Id)
	if err != nil || a == nil || a.GetString("status") != "Approved" {
		return err
	}
	saved := a.GetString("review_token")
	current, err := requestReviewToken(tx, r)
	if err != nil || saved == "" || saved == current {
		return err
	}
	legacy, err := requestLegacyReviewTokens(tx, r)
	if err != nil || !slices.Contains(legacy, saved) {
		return err
	}
	a.Set("review_token", current)
	return tx.SaveWithContext(workflowContext(context.Background()), a)
}

// upgradeLegacyRequestApprovals upgrades the saved approval of every request
// that waits for job creation (upgradeLegacyRequestApproval). It runs when the
// server starts, so a profile edit or replace after a deploy does not make an
// approval saved before the deploy stale. With nothing to upgrade it only
// reads.
func upgradeLegacyRequestApprovals(app core.App) error {
	return app.RunInTransaction(func(tx core.App) error {
		rows, err := tx.FindRecordsByFilter("job_requests", creationQueueFilter, "", 0, 0)
		if err != nil {
			return err
		}
		for _, r := range rows {
			if err := upgradeLegacyRequestApproval(tx, r); err != nil {
				return err
			}
		}
		return nil
	})
}

// upgradeLegacyApprovalsOnServe runs the upgrade once each time the server
// starts. A failure is logged and does not stop the server: job creation
// still accepts an older digest that matches (approvedDigestMatches).
func upgradeLegacyApprovalsOnServe(app core.App) {
	if err := upgradeLegacyRequestApprovals(app); err != nil {
		app.Logger().Error("upgrade of saved project request approvals failed", "error", err)
	}
}

// UpgradeLegacyRequestApprovalsForTest runs the upgrade that the server runs
// when it starts.
func UpgradeLegacyRequestApprovalsForTest(app core.App) error {
	return upgradeLegacyRequestApprovals(app)
}
