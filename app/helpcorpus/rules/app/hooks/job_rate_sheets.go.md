<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/hooks/job_rate_sheets.go","sha256":"9211225d92b0bc61f559c45126a70059578f47c407fc7db35a2a84f7cb6f7b33","model":"claude-sonnet-5-5"} -->

# Rules in app/hooks/job_rate_sheets.go

## Rate sheet history rows are priced from their start date

Lines 16-20

Asked as: rate sheet history pricing; which rate applies from date; rate changes over time; time priced by start date

Each row in a project's rate sheet history (job_rate_sheets) prices the project's time from its start date (effective_date) until the next row's start date. The opening row has an empty start date and covers all time before the first dated row. The project's rate sheet (jobs.rate_sheet) always equals the latest row's sheet.

## Opening row skips the active rate sheet check

Lines 26-29

Asked as: old rate sheet deactivated; opening row inactive sheet; copied rate sheet not active; current sheet was disabled

An opening row copied from the project's current rate sheet does not need that sheet to be active. The sheet was checked when it was set and may have been deactivated since.

## Today for rate sheet history is the UTC date

Lines 39-43

Asked as: what is today's date; UTC date for rates; rate sheet date cutoff; latest allowed rate date

The latest start date a rate sheet period may have is today's date in UTC, the same date WIP uses as today.

## Opening row is created for a project with a rate sheet but no history

Lines 58-81

Asked as: project missing rate history; imported project no rates; no opening row; old project rate sheet blank

Turbo creates an opening row (empty start date) for a project that has a rate sheet and no history rows. This covers imported and older projects. Proposals and projects with no rate sheet get no row. The opening row copies the project's rate sheet and records the acting user as creator and updater.

## Neighbouring periods cannot use the same rate sheet

Lines 83-129

Asked as: same rate sheet twice; duplicate rate sheet periods; two periods same rates; rate sheet already used next

Turbo refuses a history change that leaves two neighbouring periods on the same rate sheet, with "Two neighbouring periods would use the same rate sheet. Delete the later one instead." Pairs that already existed before the change (for example after a rate sheet absorb) do not block unrelated changes.

## History must keep its opening row

Lines 120-123

Asked as: first rate period has date; cannot remove opening row; earliest rate row dated; history lost opening period

After a change to a project's history, the earliest row must have no start date. Otherwise Turbo refuses with "The rate sheet history must keep its opening row."

## Project rate sheet follows the latest history row

Lines 133-153

Asked as: project rate sheet not updating; job rate sheet wrong; project rate differs from history; latest rate not showing

Turbo sets the project's rate sheet to the sheet of its latest history row. It saves the project only when the sheet differs, skips validation of other fields, and marks the project as no longer imported if it changed.

## Proposals cannot have rate sheet periods

Lines 157-163

Asked as: rate sheet on proposal; proposal rates not allowed; add rates to proposal; proposal rate period error

Turbo refuses a rate sheet period on a proposal with "Proposals do not have rate sheets." If the project is not found it says "The project was not found."

## A rate sheet period cannot move to another project

Lines 165-169

Asked as: move rate period to other job; change project of rate period; reassign rate sheet period; shift rates to another project

When a period is edited, its project cannot change. Turbo refuses with "A rate sheet period cannot move to another project."

## Opening period cannot gain a date and dated periods cannot lose theirs

Lines 170-173

Asked as: add date to first rate; remove date from rate period; clear start date rates; opening rate date change

Editing a period cannot add a start date to the opening row or clear it on a dated row. Turbo refuses with "The opening period has no start date, and other periods must keep one."

## Period start date must be valid and not in the future

Lines 175-213

Asked as: rate start date invalid; future rate date; rate date tomorrow; wrong date format rates

A dated period needs a valid date (YYYY-MM-DD), or Turbo says "Enter a valid date." A date after today (UTC) is refused with "The start date cannot be in the future." Today itself is allowed.

## Newly chosen rate sheet must be active

Lines 181-195

Asked as: inactive rate sheet chosen; cannot use disabled rate sheet; deactivated rate sheet selected; rate sheet not found

When a period is created or its rate sheet changes, the sheet must be active, or Turbo says "The selected rate sheet is not active." A row that keeps its sheet is not rechecked if the sheet was deactivated later. A missing sheet gives "The rate sheet was not found."

## Unchanged sheet and date skip sheet checks

Lines 181-185

Asked as: editing rate period no change; save rate period without changes; rate sheet check unnecessary; resave same rate period

When an edit changes neither the rate sheet nor the start date of an existing period, Turbo skips the sheet checks.

## Period cannot start before its rate sheet revision took effect

Lines 196-223

Asked as: rate date before sheet existed; start before rate sheet effective; rate period too early; sheet effective date conflict

A dated period whose start date is before the rate sheet's own effective date is refused with "The start date is before <sheet> took effect on <date>." A start date equal to that date is allowed. The opening period is exempt.

## Rate sheet history can only change through its actions

Lines 227-243

Asked as: cannot edit rate history directly; delete rate row manually; rate history import blocked; change rate row outside actions

Writing, editing or deleting a history row outside the rate sheet history actions is refused with "use the rate sheet history actions". Row validation runs on create and update, not on delete.

## First rate sheet on a project creates the opening row

Lines 244-257

Asked as: new project rate sheet setup; first rate sheet added; project created with rates; give project a rate sheet

When a project is created with a rate sheet, or a project without one is given one, Turbo creates the opening row after saving the project, in the same transaction. It records the acting user.
