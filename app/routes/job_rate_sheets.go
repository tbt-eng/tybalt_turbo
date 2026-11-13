package routes

import (
	"errors"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"tybalt/errs"
	"tybalt/hooks"
	"tybalt/utilities"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

// A project's rate sheet history decides which rate sheet prices its time on
// each date. See descriptions/job_rate_sheet_history.md.

// openingRateSheetPeriodID addresses the opening period of a project whose
// history has not been written yet. Its row is created on its first change.
const openingRateSheetPeriodID = "opening"

const rateSheetReasonMaxLength = 500

type rateSheetPeriodSheet struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Revision      int    `json:"revision"`
	EffectiveDate string `json:"effective_date"`
	Active        bool   `json:"active"`
}

type rateSheetPeriod struct {
	ID        string               `json:"id"`
	RateSheet rateSheetPeriodSheet `json:"rate_sheet"`
	// Starts is empty for the opening period, which covers all earlier time.
	Starts string `json:"starts"`
	// Through is the last day the period applies, or empty for the latest.
	Through   string `json:"through"`
	Reason    string `json:"reason"`
	CreatedBy string `json:"created_by"`
	UpdatedBy string `json:"updated_by"`
	Updated   string `json:"updated"`
}

type rateSheetHistory struct {
	// CanRevise holders may add, change and delete periods.
	CanRevise bool `json:"can_revise"`
	// CanSetFirst is true when the project has no rate sheet and the user may
	// set its first one.
	CanSetFirst bool              `json:"can_set_first"`
	Periods     []rateSheetPeriod `json:"periods"`
}

type rateSheetPeriodRequest struct {
	RateSheet     *string `json:"rate_sheet"`
	EffectiveDate *string `json:"effective_date"`
	Reason        string  `json:"reason"`
}

func rateSheetRouteError(status int, field, code, message string) *errs.HookError {
	err := &errs.HookError{Status: status, Message: message}
	if field != "" {
		err.Data = map[string]errs.CodeError{field: {Code: code, Message: message}}
	}
	return err
}

func writeRateSheetError(e *core.RequestEvent, err error) error {
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

// rateSheetPermissions reports whether the user may revise a project's
// history, and whether they may set a first rate sheet. Both need an active
// user with the job claim; revising also needs rate_sheet_revise.
func rateSheetPermissions(app core.App, auth *core.Record) (job bool, revise bool, err error) {
	active, err := utilities.IsUserActive(app, auth.Id)
	if err != nil || !active {
		return false, false, err
	}
	job, err = utilities.HasClaim(app, auth, "job")
	if err != nil || !job {
		return false, false, err
	}
	revise, err = utilities.HasClaim(app, auth, "rate_sheet_revise")
	return job, revise, err
}

func loadRateSheetProject(app core.App, id string) (*core.Record, error) {
	job, err := app.FindRecordById("jobs", id)
	if err != nil {
		return nil, rateSheetRouteError(http.StatusNotFound, "", "", "Project not found.")
	}
	if strings.HasPrefix(job.GetString("number"), "P") {
		return nil, rateSheetRouteError(http.StatusBadRequest, "", "", "Proposals do not have rate sheets.")
	}
	return job, nil
}

func profileName(app core.App, uid string) string {
	if uid == "" {
		return ""
	}
	profile, err := app.FindFirstRecordByData("profiles", "uid", uid)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(profile.GetString("given_name") + " " + profile.GetString("surname"))
}

func dayBefore(date string) string {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return ""
	}
	return t.AddDate(0, 0, -1).Format("2006-01-02")
}

func buildRateSheetHistory(app core.App, job *core.Record) ([]rateSheetPeriod, error) {
	rows, err := hooks.JobRateSheetRows(app, job.Id)
	if err != nil {
		return nil, err
	}
	periods := []rateSheetPeriod{}
	// A project with a rate sheet but no history has one implied opening
	// period on its current sheet. Pricing treats it the same way.
	if len(rows) == 0 && job.GetString("rate_sheet") != "" {
		period := rateSheetPeriod{ID: openingRateSheetPeriodID}
		if sheet, err := app.FindRecordById("rate_sheets", job.GetString("rate_sheet")); err == nil {
			period.RateSheet = sheetSummary(sheet)
		}
		return append(periods, period), nil
	}
	for i, row := range rows {
		period := rateSheetPeriod{
			ID:        row.Id,
			Starts:    row.GetString("effective_date"),
			Reason:    row.GetString("reason"),
			CreatedBy: profileName(app, row.GetString("created_by")),
			UpdatedBy: profileName(app, row.GetString("updated_by")),
			Updated:   row.GetString("updated"),
		}
		if i+1 < len(rows) {
			period.Through = dayBefore(rows[i+1].GetString("effective_date"))
		}
		if sheet, err := app.FindRecordById("rate_sheets", row.GetString("rate_sheet")); err == nil {
			period.RateSheet = sheetSummary(sheet)
		}
		periods = append(periods, period)
	}
	return periods, nil
}

func sheetSummary(sheet *core.Record) rateSheetPeriodSheet {
	return rateSheetPeriodSheet{
		ID:            sheet.Id,
		Name:          sheet.GetString("name"),
		Revision:      sheet.GetInt("revision"),
		EffectiveDate: sheet.GetString("effective_date"),
		Active:        sheet.GetBool("active"),
	}
}

// createGetJobRateSheetsHandler returns a project's rate sheet history.
func createGetJobRateSheetsHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		job, err := loadRateSheetProject(app, e.Request.PathValue("id"))
		if err != nil {
			return writeRateSheetError(e, err)
		}
		canJob, canRevise, err := rateSheetPermissions(app, e.Auth)
		if err != nil {
			return writeRateSheetError(e, err)
		}
		periods, err := buildRateSheetHistory(app, job)
		if err != nil {
			return writeRateSheetError(e, err)
		}
		return e.JSON(http.StatusOK, rateSheetHistory{
			CanRevise:   canRevise,
			CanSetFirst: canJob && job.GetString("rate_sheet") == "",
			Periods:     periods,
		})
	}
}

// changeJobRateSheets runs one change to a project's history in a transaction,
// then checks the whole history and updates jobs.rate_sheet to match it.
func changeJobRateSheets(app core.App, e *core.RequestEvent, change func(tx core.App, job *core.Record, revise bool) error) error {
	err := app.RunInTransaction(func(tx core.App) error {
		editing, err := utilities.IsJobsEditingEnabled(tx)
		if err != nil {
			return err
		}
		if !editing {
			return utilities.ErrJobsEditingDisabled
		}
		job, err := loadRateSheetProject(tx, e.Request.PathValue("id"))
		if err != nil {
			return err
		}
		canJob, canRevise, err := rateSheetPermissions(tx, e.Auth)
		if err != nil {
			return err
		}
		if !canJob {
			return rateSheetRouteError(http.StatusForbidden, "", "", "An active user with the job claim is required.")
		}
		if err := checkPendingRateSheetAbsorb(tx, job); err != nil {
			return err
		}
		existing, err := hooks.JobRateSheetSamePairs(tx, job.Id)
		if err != nil {
			return err
		}
		if err := change(tx, job, canRevise); err != nil {
			return err
		}
		if err := hooks.ValidateJobRateSheetHistory(tx, job.Id, existing); err != nil {
			return err
		}
		return hooks.SyncJobRateSheet(tx, e.Request.Context(), job)
	})
	if err != nil {
		return writeRateSheetError(e, err)
	}
	job, err := app.FindRecordById("jobs", e.Request.PathValue("id"))
	if err != nil {
		return writeRateSheetError(e, err)
	}
	periods, err := buildRateSheetHistory(app, job)
	if err != nil {
		return writeRateSheetError(e, err)
	}
	return e.JSON(http.StatusOK, map[string]any{"periods": periods})
}

// checkPendingRateSheetAbsorb refuses a change to a project's history while a
// rate sheet absorb that moved the project's rate sheet or one of its periods
// can still be undone. Undo writes back the saved sheet of each moved job and
// period, so a change made in between would be overwritten, or would leave
// jobs.rate_sheet unequal to the latest period. As for a client change
// (checkPendingClientMerges), the user commits or undoes the absorb first. A
// period can't be pointed at an absorbed sheet, because absorb deletes it.
func checkPendingRateSheetAbsorb(tx core.App, job *core.Record) error {
	actions, err := tx.FindAllRecords("absorb_actions", dbx.HashExp{"collection_name": "rate_sheets"})
	if err != nil || len(actions) == 0 {
		return err
	}
	rows, err := hooks.JobRateSheetRows(tx, job.Id)
	if err != nil {
		return err
	}
	for _, action := range actions {
		var updates map[string]map[string]map[string]string
		if err := action.UnmarshalJSONField("updated_references", &updates); err != nil {
			return err
		}
		_, tracked := updates["jobs"]["rate_sheet"][job.Id]
		for _, row := range rows {
			if _, ok := updates["job_rate_sheets"]["rate_sheet"][row.Id]; ok {
				tracked = true
			}
		}
		if tracked {
			return rateSheetRouteError(http.StatusConflict, "", "", "A rate sheet absorb that moved this project's rate sheet is waiting to be committed or undone. Finish it before changing the rate sheet history.")
		}
	}
	return nil
}

func requireRevise(revise bool) error {
	if !revise {
		return rateSheetRouteError(http.StatusForbidden, "", "", "Changing a project's rate sheet history requires the rate_sheet_revise claim.")
	}
	return nil
}

func checkRateSheetReason(reason string) error {
	n := utf8.RuneCountInString(reason)
	if n == 0 {
		return rateSheetRouteError(http.StatusBadRequest, "reason", "required", "Explain why the rate sheet is changing.")
	}
	if n > rateSheetReasonMaxLength {
		return rateSheetRouteError(http.StatusBadRequest, "reason", "too_long", "Keep the reason to 500 characters or fewer.")
	}
	return nil
}

// checkDateFree refuses a second period on the same date, with a clearer
// message than the unique index gives.
func checkDateFree(tx core.App, jobID, date, exceptID string) error {
	rows, err := tx.FindRecordsByFilter("job_rate_sheets", "job = {:job} && effective_date = {:date} && id != {:id}", "", 1, 0,
		dbx.Params{"job": jobID, "date": date, "id": exceptID})
	if err != nil {
		return err
	}
	if len(rows) > 0 {
		return rateSheetRouteError(http.StatusBadRequest, "effective_date", "duplicate", "Another period already starts on this date.")
	}
	return nil
}

// createAddJobRateSheetHandler adds a period. Without a start date it sets the
// first rate sheet of a project that has none, which needs only the job claim.
// With one, it starts a new rate sheet from that date.
func createAddJobRateSheetHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		var req rateSheetPeriodRequest
		if err := e.BindBody(&req); err != nil {
			return e.BadRequestError("invalid request body", err)
		}
		req.Reason = strings.TrimSpace(req.Reason)
		if req.RateSheet == nil || *req.RateSheet == "" {
			return writeRateSheetError(e, rateSheetRouteError(http.StatusBadRequest, "rate_sheet", "required", "Select a rate sheet."))
		}
		date := ""
		if req.EffectiveDate != nil {
			date = strings.TrimSpace(*req.EffectiveDate)
		}
		return changeJobRateSheets(app, e, func(tx core.App, job *core.Record, revise bool) error {
			ctx := hooks.WithProjectWorkflowMutation(e.Request.Context(), hooks.ProjectWorkflowMutationRateSheet)
			collection, err := tx.FindCollectionByNameOrId("job_rate_sheets")
			if err != nil {
				return err
			}
			row := core.NewRecord(collection)
			row.Set("job", job.Id)
			row.Set("rate_sheet", *req.RateSheet)
			row.Set("created_by", e.Auth.Id)
			row.Set("updated_by", e.Auth.Id)

			if date == "" {
				if job.GetString("rate_sheet") != "" {
					return rateSheetRouteError(http.StatusBadRequest, "effective_date", "required", "Enter the date the new rate sheet starts.")
				}
				if utf8.RuneCountInString(req.Reason) > rateSheetReasonMaxLength {
					return checkRateSheetReason(req.Reason)
				}
				row.Set("effective_date", "")
				row.Set("reason", req.Reason)
				return tx.SaveWithContext(ctx, row)
			}

			if err := requireRevise(revise); err != nil {
				return err
			}
			if job.GetString("rate_sheet") == "" {
				return rateSheetRouteError(http.StatusBadRequest, "effective_date", "first_sheet", "Set the project's first rate sheet before adding a dated one.")
			}
			if err := checkRateSheetReason(req.Reason); err != nil {
				return err
			}
			if err := hooks.EnsureOpeningRateSheetRow(tx, e.Request.Context(), job, ""); err != nil {
				return err
			}
			if err := checkDateFree(tx, job.Id, date, ""); err != nil {
				return err
			}
			row.Set("effective_date", date)
			row.Set("reason", req.Reason)
			return tx.SaveWithContext(ctx, row)
		})
	}
}

// findRateSheetPeriod finds a period of the project, creating the opening row
// first when the history has not been written yet.
func findRateSheetPeriod(tx core.App, e *core.RequestEvent, job *core.Record) (*core.Record, error) {
	if err := hooks.EnsureOpeningRateSheetRow(tx, e.Request.Context(), job, ""); err != nil {
		return nil, err
	}
	id := e.Request.PathValue("period")
	if id == openingRateSheetPeriodID {
		row, err := tx.FindFirstRecordByFilter("job_rate_sheets", "job = {:job} && effective_date = ''", dbx.Params{"job": job.Id})
		if err != nil {
			return nil, rateSheetRouteError(http.StatusNotFound, "", "", "This project has no rate sheet.")
		}
		return row, nil
	}
	row, err := tx.FindRecordById("job_rate_sheets", id)
	if err != nil || row.GetString("job") != job.Id {
		return nil, rateSheetRouteError(http.StatusNotFound, "", "", "Rate sheet period not found. The history may have changed; reload it.")
	}
	return row, nil
}

// createUpdateJobRateSheetHandler changes a period's rate sheet or start date.
// Changing the opening period's sheet reprices all time before the first dated
// period.
func createUpdateJobRateSheetHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		var req rateSheetPeriodRequest
		if err := e.BindBody(&req); err != nil {
			return e.BadRequestError("invalid request body", err)
		}
		req.Reason = strings.TrimSpace(req.Reason)
		return changeJobRateSheets(app, e, func(tx core.App, job *core.Record, revise bool) error {
			if err := requireRevise(revise); err != nil {
				return err
			}
			if err := checkRateSheetReason(req.Reason); err != nil {
				return err
			}
			row, err := findRateSheetPeriod(tx, e, job)
			if err != nil {
				return err
			}
			if req.RateSheet != nil {
				if *req.RateSheet == "" {
					return rateSheetRouteError(http.StatusBadRequest, "rate_sheet", "required", "Select a rate sheet.")
				}
				row.Set("rate_sheet", *req.RateSheet)
			}
			if req.EffectiveDate != nil {
				date := strings.TrimSpace(*req.EffectiveDate)
				if err := checkDateFree(tx, job.Id, date, row.Id); err != nil {
					return err
				}
				row.Set("effective_date", date)
			}
			row.Set("reason", req.Reason)
			row.Set("updated_by", e.Auth.Id)
			ctx := hooks.WithProjectWorkflowMutation(e.Request.Context(), hooks.ProjectWorkflowMutationRateSheet)
			return tx.SaveWithContext(ctx, row)
		})
	}
}

// createDeleteJobRateSheetHandler deletes a dated period. Its time is then
// priced by the period before it. The opening period cannot be deleted.
func createDeleteJobRateSheetHandler(app core.App) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		return changeJobRateSheets(app, e, func(tx core.App, job *core.Record, revise bool) error {
			if err := requireRevise(revise); err != nil {
				return err
			}
			row, err := findRateSheetPeriod(tx, e, job)
			if err != nil {
				return err
			}
			if row.GetString("effective_date") == "" {
				return rateSheetRouteError(http.StatusBadRequest, "", "", "The opening period cannot be deleted. Change its rate sheet instead.")
			}
			ctx := hooks.WithProjectWorkflowMutation(e.Request.Context(), hooks.ProjectWorkflowMutationRateSheet)
			return tx.DeleteWithContext(ctx, row)
		})
	}
}
