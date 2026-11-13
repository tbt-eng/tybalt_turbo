package routes

import (
	"errors"
	"net/http"
	"strings"

	"tybalt/hooks"
	"tybalt/utilities"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

// A revision rollout moves every active project on an earlier revision of a
// rate sheet to its latest revision, from one date. Each project gets a new
// period in the same way as a new period on the job details page, so the
// history rules and the jobs.rate_sheet mirror apply to each one. Time from
// the date is repriced, so only admins can do it. See the Revision rollout
// section of descriptions/job_rate_sheet_history.md.

// rateSheetRollout counts the projects on an earlier revision of a sheet.
type rateSheetRollout struct {
	EffectiveDate string `json:"effective_date"`
	// Move is the number of projects that move, or that moved.
	Move int `json:"move"`
	// The projects that do not move, by reason.
	Closed      int `json:"closed"`
	Cancelled   int `json:"cancelled"`
	LaterPeriod int `json:"later_period"`
	jobIDs      []string
}

type rateSheetRolloutRequest struct {
	EffectiveDate string `json:"effective_date"`
	Reason        string `json:"reason"`
	// ExpectedCount is the Move count that the admin confirmed.
	ExpectedCount *int `json:"expected_count"`
}

// errRolloutChanged is returned from the rollout transaction when the
// projects to move are not the number that the admin confirmed.
var errRolloutChanged = errors.New("the projects to move have changed")

func checkRolloutPermission(app core.App, auth *core.Record) error {
	active, err := utilities.IsUserActive(app, auth.Id)
	if err != nil {
		return err
	}
	admin, err := utilities.HasClaim(app, auth, "admin")
	if err != nil {
		return err
	}
	if !active || !admin {
		return rateSheetRouteError(http.StatusForbidden, "", "", "Only an active admin can move projects to a revision.")
	}
	return nil
}

// loadRolloutSheet finds the revision that projects move to. It must be the
// latest revision of its name, and active, because a new period needs an
// active sheet.
func loadRolloutSheet(app core.App, id string) (*core.Record, error) {
	sheet, err := app.FindRecordById("rate_sheets", id)
	if err != nil {
		return nil, rateSheetRouteError(http.StatusNotFound, "", "", "Rate sheet not found.")
	}
	newer, err := hooks.CheckNewerRevisionExists(app, sheet.GetString("name"), sheet.GetInt("revision"))
	if err != nil {
		return nil, err
	}
	if newer {
		return nil, rateSheetRouteError(http.StatusBadRequest, "", "", "Projects can move only to the latest revision of a rate sheet.")
	}
	if !sheet.GetBool("active") {
		return nil, rateSheetRouteError(http.StatusBadRequest, "", "", "Activate this revision before you move projects to it.")
	}
	return sheet, nil
}

// checkRolloutDate applies the date rules of a new period.
func checkRolloutDate(sheet *core.Record, date string) error {
	if err := hooks.CheckRateSheetPeriodDate(date); err != nil {
		return err
	}
	return hooks.CheckRateSheetPeriodRevision(sheet, date)
}

// findRolloutProjects finds the projects whose current sheet is an earlier
// revision of the sheet, and sorts them into those that move from date and
// those that are skipped.
func findRolloutProjects(app core.App, sheet *core.Record, date string) (rateSheetRollout, error) {
	var rows []struct {
		ID     string `db:"id"`
		Status string `db:"status"`
		// LatestStart is the start of the latest dated period, or empty.
		LatestStart string `db:"latest_start"`
	}
	err := app.DB().NewQuery(`
		SELECT j.id, j.status,
			COALESCE((SELECT MAX(p.effective_date) FROM job_rate_sheets p WHERE p.job = j.id), '') AS latest_start
		FROM jobs j
		JOIN rate_sheets rs ON rs.id = j.rate_sheet
		WHERE rs.name = {:name} AND rs.revision < {:revision}
			AND j.number NOT LIKE 'P%'
		ORDER BY j.number
	`).Bind(dbx.Params{
		"name":     sheet.GetString("name"),
		"revision": sheet.GetInt("revision"),
	}).All(&rows)
	if err != nil {
		return rateSheetRollout{}, err
	}

	result := rateSheetRollout{EffectiveDate: date}
	for _, row := range rows {
		switch {
		case row.Status == "Cancelled":
			result.Cancelled++
		case row.Status != "Active":
			// A project is Active, Closed or Cancelled.
			result.Closed++
		case row.LatestStart >= date:
			// A period from date would not be the latest one.
			result.LaterPeriod++
		default:
			result.Move++
			result.jobIDs = append(result.jobIDs, row.ID)
		}
	}
	return result, nil
}

// checkNoPendingRateSheetAbsorb refuses a rollout while any rate sheet absorb
// can still be undone. Undo writes back the saved sheet of each job and
// period that the absorb moved, and would overwrite a moved project.
func checkNoPendingRateSheetAbsorb(app core.App) error {
	pending, err := app.CountRecords("absorb_actions", dbx.HashExp{"collection_name": "rate_sheets"})
	if err != nil {
		return err
	}
	if pending > 0 {
		return rateSheetRouteError(http.StatusConflict, "", "", "A rate sheet absorb is waiting to be committed or undone. Finish it before moving projects.")
	}
	return nil
}

// createGetRateSheetRolloutHandler counts the projects that would move to the
// revision, and those skipped. The date defaults to the revision's effective
// date.
func createGetRateSheetRolloutHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		if err := checkRolloutPermission(app, e.Auth); err != nil {
			return writeRateSheetError(e, err)
		}
		sheet, err := loadRolloutSheet(app, e.Request.PathValue("id"))
		if err != nil {
			return writeRateSheetError(e, err)
		}
		date := strings.TrimSpace(e.Request.URL.Query().Get("effective_date"))
		if date == "" {
			date = sheet.GetString("effective_date")
		}
		if err := checkRolloutDate(sheet, date); err != nil {
			return writeRateSheetError(e, err)
		}
		result, err := findRolloutProjects(app, sheet, date)
		if err != nil {
			return writeRateSheetError(e, err)
		}
		return e.JSON(http.StatusOK, result)
	}
}

// createRateSheetRolloutHandler moves the projects to the revision, all or
// none, in one transaction. It refuses when the number of projects to move is
// not the number that the admin confirmed, and returns the new counts.
func createRateSheetRolloutHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		if err := checkRolloutPermission(app, e.Auth); err != nil {
			return writeRateSheetError(e, err)
		}
		var req rateSheetRolloutRequest
		if err := e.BindBody(&req); err != nil {
			return e.BadRequestError("invalid request body", err)
		}
		date := strings.TrimSpace(req.EffectiveDate)
		reason := strings.TrimSpace(req.Reason)
		if req.ExpectedCount == nil {
			return writeRateSheetError(e, rateSheetRouteError(http.StatusBadRequest, "expected_count", "required", "Confirm the number of projects to move."))
		}

		var result rateSheetRollout
		err := app.RunInTransaction(func(tx core.App) error {
			editing, err := utilities.IsJobsEditingEnabled(tx)
			if err != nil {
				return err
			}
			if !editing {
				return utilities.ErrJobsEditingDisabled
			}
			sheet, err := loadRolloutSheet(tx, e.Request.PathValue("id"))
			if err != nil {
				return err
			}
			if err := checkRolloutDate(sheet, date); err != nil {
				return err
			}
			if err := checkRateSheetReason(reason); err != nil {
				return err
			}
			if err := checkNoPendingRateSheetAbsorb(tx); err != nil {
				return err
			}
			result, err = findRolloutProjects(tx, sheet, date)
			if err != nil {
				return err
			}
			if result.Move != *req.ExpectedCount {
				return errRolloutChanged
			}
			jobs, err := tx.FindRecordsByIds("jobs", result.jobIDs)
			if err != nil {
				return err
			}
			ctx := e.Request.Context()
			for _, job := range jobs {
				err := changeRateSheetHistory(tx, ctx, job, func() error {
					return addDatedRateSheetPeriod(tx, ctx, job, sheet.Id, date, reason, e.Auth.Id)
				})
				if err != nil {
					return err
				}
			}
			return nil
		})
		if errors.Is(err, errRolloutChanged) {
			return e.JSON(http.StatusConflict, struct {
				Message string `json:"message"`
				rateSheetRollout
			}{"The projects to move have changed. Check the new count and try again.", result})
		}
		if err != nil {
			return writeRateSheetError(e, err)
		}
		return e.JSON(http.StatusOK, result)
	}
}
