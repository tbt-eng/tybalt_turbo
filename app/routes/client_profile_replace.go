package routes

import (
	"errors"
	"net/http"

	"tybalt/errs"
	"tybalt/utilities"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type replaceInvoicingProfileRequest struct {
	Replacement string `json:"replacement"`
}

// replaceInvoicingProfileResult counts the jobs and job requests that moved to
// the replacement profile.
type replaceInvoicingProfileResult struct {
	Jobs     int `json:"jobs"`
	Requests int `json:"requests"`
}

func replaceProfileError(status int, field, code, message string) *errs.HookError {
	err := &errs.HookError{Status: status, Message: message}
	if field != "" {
		err.Data = map[string]errs.CodeError{field: {Code: code, Message: message}}
	}
	return err
}

// createReplaceInvoicingProfileHandler retires an invoicing profile that was
// made in error: it moves every job and job request that uses the profile to
// another profile of the same client, then deletes the profile. Holders of job
// may do it. It runs in one transaction, so a failure changes nothing.
func createReplaceInvoicingProfileHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		if e.Auth == nil {
			return e.Error(http.StatusUnauthorized, "unauthorized", nil)
		}
		var req replaceInvoicingProfileRequest
		if err := e.BindBody(&req); err != nil {
			return e.Error(http.StatusBadRequest, "invalid request body", err)
		}
		var result replaceInvoicingProfileResult
		err := app.RunInTransaction(func(tx core.App) error {
			var err error
			result, err = replaceInvoicingProfile(tx, e, e.Request.PathValue("id"), e.Request.PathValue("record"), req.Replacement)
			return err
		})
		if err != nil {
			var hookErr *errs.HookError
			if errors.As(err, &hookErr) {
				status := hookErr.Status
				if status == 0 {
					status = http.StatusBadRequest
				}
				return e.JSON(status, hookErr)
			}
			return e.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}
		return e.JSON(http.StatusOK, result)
	}
}

// replaceInvoicingProfile repoints the jobs and requests of client clientID
// that use profileID to replacementID, and deletes profileID.
//
// The jobs and requests are updated with direct SQL, as a client or contact
// merge updates them (absorb_records.go), not by saving each record:
//
//   - A record save runs the job and request hooks, which are built for edits
//     by a person. The project workflow guard refuses a project or request
//     write without a workflow operation, the request actions let only the
//     creator edit a Draft, Recalled or Rejected request, and full job
//     validation can refuse a legacy or older Closed job for fields this
//     action does not touch. Only invoicing_information changes, and both
//     profiles are checked here, so those hooks would add nothing but
//     refusals.
//   - The jobs get what a normal edit of the profile gives them: _imported is
//     set to false and updated to now, as MarkReferencingJobsNotImported does,
//     so the jobs writeback picks them up. The writeback does not send the
//     profile, so the legacy system sees no other change.
//   - The profile is not part of what an approval covers (the review digest
//     leaves it out), so a request awaiting approval or job creation keeps its
//     version and its approval. A request its creator can edit (Draft,
//     Recalled or Rejected) gets a new version and last_edited, as a merge
//     gives them, so an edit form opened before the replace must reload
//     instead of saving the deleted profile.
//
// This works the same for every status, so Closed jobs and requests in any
// status move too.
func replaceInvoicingProfile(tx core.App, e *core.RequestEvent, clientID, profileID, replacementID string) (replaceInvoicingProfileResult, error) {
	result := replaceInvoicingProfileResult{}
	hasJob, err := utilities.HasClaim(tx, e.Auth, "job")
	if err != nil {
		return result, err
	}
	if !hasJob {
		return result, replaceProfileError(http.StatusForbidden, "", "", "A user with the job claim is required to replace an invoicing profile.")
	}
	// The action edits jobs, so it follows the jobs editing switch.
	editing, err := utilities.IsJobsEditingEnabled(tx)
	if err != nil {
		return result, err
	}
	if !editing {
		return result, utilities.ErrJobsEditingDisabled
	}

	profile, err := tx.FindRecordById("client_invoicing_information", profileID)
	if err != nil || profile.GetString("client") != clientID {
		return result, replaceProfileError(http.StatusNotFound, "", "", "Invoicing profile not found for this client.")
	}
	if replacementID == "" {
		return result, replaceProfileError(http.StatusBadRequest, "replacement", "required", "Select the invoicing profile to use instead.")
	}
	if replacementID == profile.Id {
		return result, replaceProfileError(http.StatusBadRequest, "replacement", "unchanged", "Select a different invoicing profile.")
	}
	replacement, err := tx.FindRecordById("client_invoicing_information", replacementID)
	if err != nil {
		return result, replaceProfileError(http.StatusBadRequest, "replacement", "invalid_reference", "The selected invoicing profile was not found.")
	}
	if replacement.GetString("client") != clientID {
		return result, replaceProfileError(http.StatusBadRequest, "replacement", "client_mismatch", "The selected invoicing profile must belong to this client.")
	}

	var jobIDs, requestIDs []string
	if err := tx.DB().Select("id").From("jobs").Where(dbx.HashExp{"invoicing_information": profile.Id}).Column(&jobIDs); err != nil {
		return result, err
	}
	if err := tx.DB().Select("id").From("job_requests").Where(dbx.HashExp{"invoicing_information": profile.Id}).Column(&requestIDs); err != nil {
		return result, err
	}

	// Undo of a pending merge would restore references on these records, or
	// move a profile back to another client and then refuse to finish.
	tracked := map[string]bool{profile.Id: true, replacement.Id: true}
	for _, id := range append(append([]string{}, jobIDs...), requestIDs...) {
		tracked[id] = true
	}
	pending, err := pendingClientMergeTouches(tx, tracked)
	if err != nil {
		return result, err
	}
	if pending {
		return result, replaceProfileError(http.StatusConflict, "", "", "A client or contact merge that involves these invoicing profiles or their jobs or requests is waiting to be committed or undone. Finish it before replacing the profile.")
	}

	// An approval saved before the profile left the digest holds the profile.
	// Give it the current digest first, so the move does not make it stale.
	for _, id := range requestIDs {
		r, err := tx.FindRecordById("job_requests", id)
		if err != nil {
			return result, err
		}
		if err := upgradeLegacyRequestApproval(tx, r); err != nil {
			return result, err
		}
	}

	params := dbx.Params{"profile": profile.Id, "replacement": replacement.Id}
	if _, err := tx.DB().NewQuery(`
		UPDATE jobs
		SET invoicing_information = {:replacement},
		    _imported = false,
		    updated = strftime('%Y-%m-%d %H:%M:%fZ', 'now')
		WHERE invoicing_information = {:profile}
	`).Bind(params).Execute(); err != nil {
		return result, err
	}
	if _, err := tx.DB().NewQuery(`
		UPDATE job_requests
		SET invoicing_information = {:replacement},
		    version = CASE WHEN status IN ('Draft', 'Recalled', 'Rejected') THEN version + 1 ELSE version END,
		    last_edited = CASE WHEN status IN ('Draft', 'Recalled', 'Rejected') THEN strftime('%Y-%m-%d %H:%M:%fZ', 'now') ELSE last_edited END,
		    updated = strftime('%Y-%m-%d %H:%M:%fZ', 'now')
		WHERE invoicing_information = {:profile}
	`).Bind(params).Execute(); err != nil {
		return result, err
	}
	// Every record that uses the replacement must still share its client.
	if err := utilities.ValidateClientReferences(tx, replacement); err != nil {
		return result, err
	}
	// The delete hook refuses a profile that is still used.
	if err := tx.Delete(profile); err != nil {
		return result, err
	}
	result.Jobs = len(jobIDs)
	result.Requests = len(requestIDs)
	return result, nil
}
