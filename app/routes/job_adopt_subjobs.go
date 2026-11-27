package routes

import (
	"errors"
	"fmt"
	"maps"
	"net/http"
	"regexp"
	"strings"

	"tybalt/errs"
	"tybalt/utilities"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

// A project number that can have subjobs, such as "24-410" or "25-0123".
var jobAdoptParentNumberRegex = regexp.MustCompile(`^\d{2}-\d{3,4}$`)

type jobAdoptSubjobsRequest struct {
	// Subjobs lists the subjobs the user confirmed will be linked to the project.
	Subjobs []string `json:"subjobs"`
}

// unlinkedSubjob is a job numbered as a subjob of a project but not linked to
// it, as legacy subjobs were imported.
type unlinkedSubjob struct {
	ID          string `json:"id"`
	Number      string `json:"number"`
	Description string `json:"description"`
	Client      string `json:"client"`
	// Blocker says why the job cannot be adopted. It is empty when it can.
	Blocker string `json:"blocker,omitempty"`
}

type unlinkedSubjobRow struct {
	ID          string `db:"id"`
	Number      string `db:"number"`
	Description string `db:"description"`
	Client      string `db:"client"`
	Children    int    `db:"children"`
}

// findUnlinkedSubjobs returns the jobs numbered as first-level subjobs of the
// project, such as "24-410-1" for "24-410", that have no parent. A subjob must
// use its parent's client, so one with another client is listed with a blocker
// rather than adopted, and so is one with subjobs of its own.
func findUnlinkedSubjobs(app core.App, project *core.Record) ([]unlinkedSubjob, error) {
	base := project.GetString("number")
	if !jobAdoptParentNumberRegex.MatchString(base) || project.GetString("parent") != "" {
		return nil, nil
	}
	var rows []unlinkedSubjobRow
	// GLOB keeps "24-4100-1" out of "24-410"; the child number regex below keeps
	// out nested numbers such as "24-410-1-2".
	err := app.DB().NewQuery(`
		SELECT j.id, j.number, j.description, j.client,
		  (SELECT COUNT(*) FROM jobs c WHERE c.parent = j.id) AS children
		FROM jobs j
		WHERE j.parent = '' AND j.number GLOB {:glob}
		ORDER BY j.number
	`).Bind(dbx.Params{"glob": base + "-[0-9]*"}).All(&rows)
	if err != nil {
		return nil, err
	}

	projectClient := project.GetString("client")
	projectLabel := noticeClientLabel(app, projectClient)
	subjobs := []unlinkedSubjob{}
	for _, row := range rows {
		if !jobChildNumberRegex.MatchString(row.Number) {
			continue
		}
		subjob := unlinkedSubjob{ID: row.ID, Number: row.Number, Description: row.Description, Client: noticeClientLabel(app, row.Client)}
		switch {
		case projectClient == "":
			subjob.Blocker = "This project has no client. Set its client first."
		case row.Client != projectClient:
			current := subjob.Client
			if current == "" {
				current = "no client"
			}
			subjob.Blocker = fmt.Sprintf("Its client is %s. Change its client to %s first.", current, projectLabel)
		case row.Children > 0:
			subjob.Blocker = "It has subjobs of its own."
		}
		subjobs = append(subjobs, subjob)
	}
	return subjobs, nil
}

// createAdoptSubjobsHandler links a project to the unlinked jobs numbered as
// its subjobs, which legacy subjobs were imported as. Only admins may do it.
//
// Like the set number action, it updates the subjobs directly instead of saving
// them, so the jobs hooks and record validation do not run: only the parent
// changes, and a legacy job may not meet today's rules for fields this action
// does not touch. The subjobs are marked as locally changed so they are written
// back to the legacy system, and a note on the project records the change.
func createAdoptSubjobsHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		if e.Auth == nil {
			return e.Error(http.StatusUnauthorized, "unauthorized", nil)
		}
		jobID := e.Request.PathValue("id")
		var req jobAdoptSubjobsRequest
		if err := e.BindBody(&req); err != nil {
			return e.Error(http.StatusBadRequest, "invalid request body", err)
		}

		var adopted []string
		err := app.RunInTransaction(func(tx core.App) error {
			var err error
			adopted, err = adoptSubjobs(tx, e.Auth, jobID, req)
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
		return e.JSON(http.StatusOK, map[string]any{"adopted": adopted})
	}
}

func adoptSubjobs(tx core.App, auth *core.Record, jobID string, req jobAdoptSubjobsRequest) ([]string, error) {
	isAdmin, err := utilities.HasClaim(tx, auth, "admin")
	if err != nil {
		return nil, err
	}
	if !isAdmin {
		return nil, jobChangeClientError(http.StatusForbidden, "", "", "Only admins can adopt subjobs.")
	}
	editing, err := utilities.IsJobsEditingEnabled(tx)
	if err != nil {
		return nil, err
	}
	if !editing {
		return nil, utilities.ErrJobsEditingDisabled
	}

	project, err := tx.FindRecordById("jobs", jobID)
	if err != nil {
		return nil, jobChangeClientError(http.StatusNotFound, "", "", "Project not found.")
	}
	if !jobAdoptParentNumberRegex.MatchString(project.GetString("number")) || project.GetString("parent") != "" {
		return nil, jobChangeClientError(http.StatusBadRequest, "", "", "Only a top-level project can adopt subjobs.")
	}

	unlinked, err := findUnlinkedSubjobs(tx, project)
	if err != nil {
		return nil, err
	}
	ids := []string{}
	numbers := []string{}
	for _, subjob := range unlinked {
		if subjob.Blocker == "" {
			ids = append(ids, subjob.ID)
			numbers = append(numbers, subjob.Number)
		}
	}
	if len(ids) == 0 {
		return nil, jobChangeClientError(http.StatusBadRequest, "", "", "No unlinked subjobs of this project can be adopted.")
	}
	if !sameIDSet(ids, req.Subjobs) {
		return nil, jobChangeClientError(http.StatusConflict, "subjobs", "unconfirmed", "The subjobs confirmed do not match the subjobs this project can adopt, which may have changed. Review the current list and confirm again.")
	}

	inIDs := make([]any, len(ids))
	for i, id := range ids {
		inIDs[i] = id
	}
	if _, err := tx.DB().Update("jobs", dbx.Params{
		"parent":    project.Id,
		"_imported": false,
		"updated":   types.NowDateTime().String(),
	}, dbx.In("id", inIDs...)).Execute(); err != nil {
		return nil, err
	}

	note := fmt.Sprintf("Adopted unlinked subjobs: %s.", strings.Join(numbers, ", "))
	if err := createJobAuditNote(tx, project, auth.Id, note, ""); err != nil {
		return nil, err
	}
	return numbers, nil
}

// sameIDSet reports whether two lists hold the same ids, ignoring order and
// repeats.
func sameIDSet(want, got []string) bool {
	wantSet := map[string]bool{}
	for _, id := range want {
		wantSet[id] = true
	}
	gotSet := map[string]bool{}
	for _, id := range got {
		gotSet[id] = true
	}
	return maps.Equal(wantSet, gotSet)
}
