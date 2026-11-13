<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/job_rate_sheets.go","sha256":"0657bab2e24411d5f776717fb8984b3d6e02c4612cbd1daa9440dac8c12ac38b","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/job_rate_sheets.go

## Rate sheet history: opening row and dated periods

Lines 16-20

Asked as: opening row; rate sheet history; price before first date; time priced by start date

A project's rate sheet history (job_rate_sheets) prices its time from each row's start date (effective_date) until the next row's start date. The opening row has no start date and covers all time before the first dated row. The project's current rate sheet (jobs.rate_sheet) always equals the sheet on the latest row.

## Start date cannot be in the future (UTC)

Lines 39-43

Asked as: future start date; start date in the future; date later than today; rate sheet date tomorrow

A dated rate sheet period must start on or before today, and Turbo uses the UTC date as today. A later date gives "The start date cannot be in the future." A today date passes.

## Opening row backfilled for projects without history

Lines 58-81

Asked as: missing history rows; imported project no history; old project rate sheet history; backfill opening row

If a project has a rate sheet but no history rows, Turbo creates an opening row from the project's rate sheet. This covers imported projects and projects created before the history existed. Proposals and projects without a rate sheet get no row. A project that already has rows gets none.

## Neighbouring periods cannot use the same rate sheet

Lines 83-131

Asked as: same rate sheet twice; two periods same sheet; neighbouring periods duplicate; delete the later one

A change cannot leave two neighbouring periods on the same rate sheet. Turbo refuses with "Two neighbouring periods would use the same rate sheet. Delete the later one instead." Pairs that already existed before the change, such as those left by a rate sheet absorb, do not block unrelated changes.

## History must keep its opening row

Lines 117-123

Asked as: earliest row has date; first row needs no date; cannot remove opening row; history lost its opening row

After a change, the earliest row in a project's history must have no start date. Otherwise Turbo refuses with "The rate sheet history must keep its opening row."

## Project rate sheet follows the latest history row

Lines 133-153

Asked as: project rate sheet not updated; project shows wrong sheet; current sheet out of sync; imported flag cleared

After the history changes, Turbo sets the project's rate sheet to the sheet on its latest row. It saves the project only when the sheet differs. It skips validation of other fields, so a legacy project is not blocked. It marks the project as no longer imported if it changed.

## Project must exist for a rate sheet period

Lines 157-160

Asked as: project not found; rate sheet period no project; missing project for period; deleted project rate sheet

A rate sheet period must point at an existing project. Otherwise Turbo refuses with "The project was not found."

## Proposals have no rate sheets

Lines 161-163

Asked as: proposal rate sheet; add rate sheet to proposal; proposal pricing sheet; rate sheet on proposal

Turbo refuses a rate sheet period for a proposal with "Proposals do not have rate sheets."

## Period cannot move to another project

Lines 165-169

Asked as: move period to other project; change project of rate period; reassign rate sheet period; wrong project on period

When a person edits a rate sheet period, its project cannot change. Turbo refuses with "A rate sheet period cannot move to another project."

## Opening period cannot gain a date; dated periods cannot lose theirs

Lines 170-173

Asked as: remove start date; clear date from period; give opening row a date; cannot blank start date

When a person edits a period, Turbo refuses to give the opening period a start date or to clear a dated period's start date. The message is "The opening period has no start date, and other periods must keep one."

## Rate sheet must exist

Lines 181-189

Asked as: rate sheet not found; missing rate sheet; chosen sheet does not exist; rate sheet gone

When a period is created, or its sheet or date changes, the chosen rate sheet must exist. Otherwise Turbo refuses with "The rate sheet was not found." Turbo skips this check on an edit that changes neither sheet nor date.

## Newly chosen rate sheet must be active

Lines 190-195

Asked as: inactive rate sheet; deactivated sheet chosen; cannot pick archived sheet; sheet no longer active

A period that newly picks a rate sheet needs an active sheet. Otherwise Turbo refuses with "The selected rate sheet is not active." A row keeps its sheet after the sheet is deactivated. An opening row copied from the project's existing rate sheet skips this check.

## Opening period is exempt from the revision date check

Lines 196-200

Asked as: opening row effective date; opening period revision date error; no date opening rate sheet; effective date check opening

The opening row has no start date, so Turbo does not compare it with the rate sheet's effective date. Only dated periods get that check.

## Start date must be a valid date

Lines 205-208

Asked as: invalid start date; date format wrong; enter date YYYY-MM-DD; bad date rate period

A dated rate sheet period needs a valid date in the form YYYY-MM-DD. Otherwise Turbo refuses with "Enter a valid date."

## Period cannot start before the rate sheet took effect

Lines 217-223

Asked as: start before rate sheet effective; date earlier than rate sheet; rate sheet not yet in effect; period starts too early

A dated period cannot start before its rate sheet revision took effect. Turbo refuses with "The start date is before <sheet name> took effect on <date>." A start date equal to the effective date passes.

## Rate sheet rows only change through history actions

Lines 227-242

Asked as: cannot edit rate sheet rows; edit history directly; rate period changes not saved; use history actions

Turbo refuses to create, update or delete a rate sheet period outside the rate sheet history actions. The message is "use the rate sheet history actions". Deletes skip the row checks.

## Opening row created automatically for a project's first rate sheet

Lines 244-257

Asked as: new project rate sheet row; first rate sheet on project; project created with rate sheet; no opening row after create

When a project is created with a rate sheet, or a project without one is given one, Turbo creates the opening row in the same save. The row uses that rate sheet and has an empty start date. It records the acting user as creator and updater.
