package hooks

import (
	"context"
	"net/http"
	"sort"
	"time"

	"tybalt/errs"
	"tybalt/utilities"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

// A project's rate sheet history lives in job_rate_sheets. Each row prices the
// project's time from its effective_date until the next row's. The opening row
// has an empty effective_date and covers all time before the first dated row.
// jobs.rate_sheet is kept equal to the latest row's sheet. See
// descriptions/job_rate_sheet_history.md.

// rateSheetActorKey carries the user who set a project's first rate sheet from
// the job save to the opening row. It is not a field, so it is never stored.
const rateSheetActorKey = "@rateSheetActor"

// rateSheetFromJobKey marks an opening row copied from jobs.rate_sheet. That
// sheet was checked when it was set and may have been deactivated since, so
// the row skips the active check.
const rateSheetFromJobKey = "@rateSheetFromJob"

func rateSheetError(status int, field, code, message string) *errs.HookError {
	err := &errs.HookError{Status: status, Message: message}
	if field != "" {
		err.Data = map[string]errs.CodeError{field: {Code: code, Message: message}}
	}
	return err
}

// RateSheetToday is the latest effective date a row may have. WIP uses the UTC
// date as today, so the history does too.
func RateSheetToday() string {
	return time.Now().UTC().Format("2006-01-02")
}

// JobRateSheetRows returns a project's history, opening row first.
func JobRateSheetRows(app core.App, jobID string) ([]*core.Record, error) {
	rows, err := app.FindRecordsByFilter("job_rate_sheets", "job = {:job}", "", 0, 0, dbx.Params{"job": jobID})
	if err != nil {
		return nil, err
	}
	// The empty opening date sorts first.
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].GetString("effective_date") < rows[j].GetString("effective_date")
	})
	return rows, nil
}

// EnsureOpeningRateSheetRow creates the opening row for a project that has a
// rate sheet but no history. Imported projects and projects created before the
// history existed have none until their history is first written.
func EnsureOpeningRateSheetRow(app core.App, ctx context.Context, job *core.Record, actorID string) error {
	if isProposalRecord(job) || job.GetString("rate_sheet") == "" {
		return nil
	}
	count, err := app.CountRecords("job_rate_sheets", dbx.HashExp{"job": job.Id})
	if err != nil || count > 0 {
		return err
	}
	collection, err := app.FindCollectionByNameOrId("job_rate_sheets")
	if err != nil {
		return err
	}
	row := core.NewRecord(collection)
	row.Set("job", job.Id)
	row.Set("rate_sheet", job.GetString("rate_sheet"))
	row.Set("effective_date", "")
	row.Set("created_by", actorID)
	row.Set("updated_by", actorID)
	row.Set(rateSheetFromJobKey, true)
	return app.SaveWithContext(WithProjectWorkflowMutation(ctx, ProjectWorkflowMutationRateSheet), row)
}

// JobRateSheetSamePairs returns the neighbouring rows of a project's history
// that use the same sheet, as "earlier id|later id". A history change never
// makes such a pair, but a rate sheet absorb can: it moves a period from an
// absorbed sheet to the target, and the period next to it may already be on
// the target. The rates are identical, so pricing does not change.
func JobRateSheetSamePairs(app core.App, jobID string) (map[string]bool, error) {
	rows, err := JobRateSheetRows(app, jobID)
	if err != nil {
		return nil, err
	}
	return sameSheetPairs(rows), nil
}

func sameSheetPairs(rows []*core.Record) map[string]bool {
	pairs := map[string]bool{}
	for i := 1; i < len(rows); i++ {
		if rows[i].GetString("rate_sheet") == rows[i-1].GetString("rate_sheet") {
			pairs[rows[i-1].Id+"|"+rows[i].Id] = true
		}
	}
	return pairs
}

// ValidateJobRateSheetHistory checks the rules that depend on a project's whole
// history, after a change to it. The unique indexes already allow only one
// opening row and one row per date. existing holds the pairs of neighbouring
// rows on the same sheet from before the change (JobRateSheetSamePairs). Only
// a pair that the change makes is refused, so a pair that an absorb left does
// not block unrelated changes.
func ValidateJobRateSheetHistory(app core.App, jobID string, existing map[string]bool) error {
	rows, err := JobRateSheetRows(app, jobID)
	if err != nil {
		return err
	}
	if len(rows) == 0 {
		return nil
	}
	if rows[0].GetString("effective_date") != "" {
		return rateSheetError(http.StatusBadRequest, "effective_date", "opening_row_required",
			"The rate sheet history must keep its opening row.")
	}
	for pair := range sameSheetPairs(rows) {
		if !existing[pair] {
			return rateSheetError(http.StatusBadRequest, "rate_sheet", "same_as_previous",
				"Two neighbouring periods would use the same rate sheet. Delete the later one instead.")
		}
	}
	return nil
}

// SyncJobRateSheet sets jobs.rate_sheet to the sheet on the project's latest
// row. It saves the job only when the sheet changed.
func SyncJobRateSheet(app core.App, ctx context.Context, job *core.Record) error {
	rows, err := JobRateSheetRows(app, job.Id)
	if err != nil {
		return err
	}
	if len(rows) == 0 {
		return nil
	}
	latest := rows[len(rows)-1].GetString("rate_sheet")
	if job.GetString("rate_sheet") == latest {
		return nil
	}
	job.Set("rate_sheet", latest)
	utilities.MarkImportedFalseIfChanged(job)
	// Only rate_sheet changes. Like the change-client action, skip validating
	// fields this action does not touch: a legacy project may not meet today's
	// rules for them. The project workflow guard still runs.
	return app.SaveNoValidateWithContext(WithProjectWorkflowMutation(ctx, ProjectWorkflowMutationRateSheet), job)
}

// validateJobRateSheetRow checks one row before it is saved.
func validateJobRateSheetRow(app core.App, row *core.Record) error {
	job, err := app.FindRecordById("jobs", row.GetString("job"))
	if err != nil {
		return rateSheetError(http.StatusBadRequest, "job", "invalid_reference", "The project was not found.")
	}
	if isProposalRecord(job) {
		return rateSheetError(http.StatusBadRequest, "job", "proposal", "Proposals do not have rate sheets.")
	}
	date := row.GetString("effective_date")
	if !row.IsNew() {
		original := row.Original()
		if original.GetString("job") != row.GetString("job") {
			return rateSheetError(http.StatusBadRequest, "job", "immutable", "A rate sheet period cannot move to another project.")
		}
		if (original.GetString("effective_date") == "") != (date == "") {
			return rateSheetError(http.StatusBadRequest, "effective_date", "opening_row_fixed",
				"The opening period has no start date, and other periods must keep one.")
		}
	}
	if date != "" {
		if err := CheckRateSheetPeriodDate(date); err != nil {
			return err
		}
	}

	sheetChanged := row.IsNew() || row.Original().GetString("rate_sheet") != row.GetString("rate_sheet")
	dateChanged := row.IsNew() || row.Original().GetString("effective_date") != date
	if !sheetChanged && !dateChanged {
		return nil
	}
	sheet, err := app.FindRecordById("rate_sheets", row.GetString("rate_sheet"))
	if err != nil {
		return rateSheetError(http.StatusBadRequest, "rate_sheet", "invalid_reference", "The rate sheet was not found.")
	}
	// A row keeps its sheet after the sheet is deactivated. Only a newly
	// chosen sheet must be active.
	fromJob, _ := row.Get(rateSheetFromJobKey).(bool)
	if sheetChanged && !fromJob && !sheet.GetBool("active") {
		return rateSheetError(http.StatusBadRequest, "rate_sheet", "rate_sheet_not_active", "The selected rate sheet is not active.")
	}
	// The opening period covers all earlier time, so it is exempt.
	if date != "" {
		return CheckRateSheetPeriodRevision(sheet, date)
	}
	return nil
}

// CheckRateSheetPeriodDate checks the start date of a dated period: a valid
// date that is not in the future.
func CheckRateSheetPeriodDate(date string) error {
	if _, err := time.Parse("2006-01-02", date); err != nil {
		return rateSheetError(http.StatusBadRequest, "effective_date", "invalid", "Enter a valid date.")
	}
	if date > RateSheetToday() {
		return rateSheetError(http.StatusBadRequest, "effective_date", "future", "The start date cannot be in the future.")
	}
	return nil
}

// CheckRateSheetPeriodRevision refuses a dated period that starts before its
// rate sheet revision took effect.
func CheckRateSheetPeriodRevision(sheet *core.Record, date string) error {
	if date < sheet.GetString("effective_date") {
		return rateSheetError(http.StatusBadRequest, "effective_date", "before_revision",
			"The start date is before "+sheet.GetString("name")+" took effect on "+sheet.GetString("effective_date")+".")
	}
	return nil
}

// AddJobRateSheetHooks keeps the history consistent for every writer.
func AddJobRateSheetHooks(app core.App) {
	// Only the history actions write rows, so the job mirror and the history
	// rules always run with them.
	guard := func(e *core.RecordEvent) error {
		if ProjectWorkflowMutationFromContext(e.Context) != ProjectWorkflowMutationRateSheet {
			return rateSheetError(http.StatusBadRequest, "", "", "use the rate sheet history actions")
		}
		if e.Type != core.ModelEventTypeDelete {
			if err := validateJobRateSheetRow(e.App, e.Record); err != nil {
				return err
			}
		}
		return e.Next()
	}
	app.OnRecordCreate("job_rate_sheets").BindFunc(guard)
	app.OnRecordUpdate("job_rate_sheets").BindFunc(guard)
	app.OnRecordDelete("job_rate_sheets").BindFunc(guard)

	// A project gets its opening row when its first rate sheet is saved: on
	// create, or when a project with no rate sheet is given one. The row is
	// written after the job, in the same transaction.
	opening := func(e *core.RecordEvent) error {
		set := e.Record.IsNew() || e.Record.Original().GetString("rate_sheet") == ""
		if err := e.Next(); err != nil || !set {
			return err
		}
		actor, _ := e.Record.Get(rateSheetActorKey).(string)
		return EnsureOpeningRateSheetRow(e.App, e.Context, e.Record, actor)
	}
	app.OnRecordCreate("jobs").BindFunc(opening)
	app.OnRecordUpdate("jobs").BindFunc(opening)
}
