<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/job_rate_sheets.go","sha256":"ae944fd397fa9e45fd3ecf12df2924741b32e1ee0dbbbec813c3f40d28e93b9d","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/job_rate_sheets.go

## Rate sheet history errors show a message on the field

Lines 31-37

Asked as: error message on field; rate sheet history error; which field has the error; refused history change message

When Turbo refuses a rate sheet history change, it shows the message and ties it to the field concerned, such as the start date (effective_date), the rate sheet (rate_sheet) or the project (job). Only the refusal for using the wrong action has no field.

## Latest start date allowed for a rate sheet period is today in UTC

Lines 39-43

Asked as: start date tomorrow; date in the future; wrong date near midnight; UTC date versus local date

Turbo takes today's date from the UTC clock (RateSheetToday) for the rate sheet history. A period's start date cannot be later than that date. Near midnight, the UTC date can differ from the local date.

## Rate sheet history is ordered with the opening period first

Lines 45-56

Asked as: order of rate sheet periods; history list sorting; which period is latest; opening period shown first

Turbo lists a project's rate sheet history (job_rate_sheets) sorted by start date (effective_date). The opening period has an empty start date, so it sorts first. The latest period is the one with the last start date.

## Opening row is backfilled for projects without history

Lines 58-81

Asked as: missing opening row; project has no history; imported project no history; backfill rate sheet history

Turbo creates an opening row for a project that has a rate sheet but no history, such as an imported or older project. Nothing is created if the project is a proposal, has no rate sheet, or already has history rows.

## Neighbouring periods cannot use the same rate sheet

Lines 88-129

Asked as: same rate sheet twice; two periods same sheet; neighbouring periods identical sheet; delete the later period

Turbo refuses a history change that leaves two neighbouring periods on the same rate sheet: 'Two neighbouring periods would use the same rate sheet. Delete the later one instead.' Pairs that already existed before the change, for example after a rate sheet absorb, do not block it.

## History must keep its opening row

Lines 117-123

Asked as: opening row missing; earliest period has date; cannot remove first row; history lost its opening row

After a change, if a project has history rows, the earliest must have no start date. Otherwise Turbo refuses with 'The rate sheet history must keep its opening row.'

## Project rate sheet follows the latest period

Lines 133-152

Asked as: project rate sheet not updating; jobs rate sheet out of date; project sheet follows history; imported flag cleared

Turbo sets the project's rate sheet (jobs.rate_sheet) to the sheet of its latest history row. It saves the project only if the sheet differs. It skips validation of other fields, and marks an imported project as no longer imported (MarkImportedFalseIfChanged).

## Proposals cannot have rate sheet periods

Lines 157-163

Asked as: rate sheet on proposal; proposal rate sheet period; proposal has no rate sheet; project not found error

Turbo refuses a rate sheet period for a proposal with 'Proposals do not have rate sheets.' If the project is not found, it says 'The project was not found.' Turbo also creates no opening row for a proposal.

## Rate sheet period cannot move to another project

Lines 165-169

Asked as: move period to other project; change project of period; reassign rate sheet period; period switched project

When a person edits an existing rate sheet period, Turbo refuses a change of its project with 'A rate sheet period cannot move to another project.'

## Opening period keeps no start date; others keep one

Lines 170-173

Asked as: opening period has start date; remove start date from period; no start date error; opening period date

Turbo refuses an edit that gives the opening period a start date, or removes the start date from another period. The message is 'The opening period has no start date, and other periods must keep one.'

## Start date must be a valid date not in the future

Lines 175-182

Asked as: invalid date entered; start date not accepted; future start date refused; date format wrong

A period's start date (effective_date) must be a valid YYYY-MM-DD date, or Turbo says 'Enter a valid date.' A date after today in UTC is refused with 'The start date cannot be in the future.' Today itself is allowed. An empty date (opening row) skips both checks.

## Newly chosen rate sheet must be active

Lines 184-198

Asked as: inactive rate sheet chosen; cannot pick deactivated sheet; rate sheet unavailable; sheet not active error

When a period is created or its rate sheet changes, the chosen rate sheet must be active. Otherwise Turbo says 'The selected rate sheet is not active.' A period that keeps its sheet is not checked, so it stays valid after the sheet is deactivated. An opening row copied from the project's rate sheet also skips the check. A missing rate sheet gives 'The rate sheet was not found.'

## Period cannot start before the rate sheet took effect

Lines 186-204

Asked as: start date before sheet; period earlier than rate sheet; date before rate sheet took effect; rate sheet effective date clash

A dated period cannot start before its rate sheet's effective date. Turbo refuses with 'The start date is before <sheet name> took effect on <date>.' The same date is allowed. The opening period has no date and is exempt. Turbo skips this check when neither the sheet nor the date changed.

## Rate sheet history rows only change through history actions

Lines 212-225

Asked as: cannot edit history table; direct database write refused; change history row directly; use history actions instead

Turbo refuses to create, update or delete a project's rate sheet history row (job_rate_sheets) unless the change comes from the rate sheet history actions. Any other write is refused with 'use the rate sheet history actions'.

## Opening row is created with a project's first rate sheet

Lines 227-237

Asked as: new project first rate sheet; opening row on project create; project given rate sheet; first history row created

When a project is created with a rate sheet, or a project with no rate sheet is given one, Turbo creates the opening history row with no start date. It uses the project's sheet and records the acting user. This runs in the same save.
