<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_rate_sheets.go","sha256":"e5f9970b0e7bd90997f7c4a4d40188e85767aa1b097e3c38b40270f47a70532b","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_rate_sheets.go

## Opening period ID stands for a project with no written history

Lines 22-24

Asked as: opening period id; project has no history; first period missing; no history written yet

The opening period (id "opening") is the first rate sheet period of a project. A project whose history has not been written yet uses this ID, and Turbo creates its row on the first change (openingRateSheetPeriodID).

## Rate sheet period reason is limited to 500 characters

Lines 26-26

Asked as: reason too long; reason character limit; rate reason too many characters; 500 character limit note

A reason for a rate sheet period change can have at most 500 characters (rateSheetReasonMaxLength).

## Latest period has no end date and opening period has no start date

Lines 36-43

Asked as: why no end date; opening period no start; period end date blank; latest period empty end

In the history, the opening period has an empty start (it covers all earlier time). The latest period has an empty end (it applies from its start with no limit).

## Rate sheet error messages are tied to the field

Lines 64-82

Asked as: error message next to field; where rate error shows; rate sheet form error placement; error has no status

Rate sheet refusals name a field (reason, rate_sheet, effective_date) with a code, so the form can show the message next to that field. When an error has no status, Turbo treats it as a bad request.

## Who can revise or set a first rate sheet

Lines 84-98

Asked as: who can edit rate sheet; permission to change rates; inactive user rate access; who can set first rate sheet

Both permissions need an active user with the job claim. Revising periods also needs the rate_sheet_revise claim. A user who is inactive or lacks the job claim gets neither permission.

## Proposals have no rate sheets

Lines 100-109

Asked as: proposal rate sheet; P job number rates; why no rates on proposal; proposal cannot have rates

Every rate sheet history action loads the job first. When the job number starts with P (a proposal), Turbo refuses with "Proposals do not have rate sheets." When the job does not exist, it says "Project not found."

## Project not found when loading rate sheet history

Lines 100-104

Asked as: project not found rate history; can't find project rates; rate history missing project; project missing for rates

When the project does not exist, Turbo answers "Project not found." for the rate sheet history and for every change to it.

## Period start dates must be valid dates to give an end date

Lines 122-128

Asked as: end date empty; period end date blank; bad start date period; start date not readable

Turbo shows the last day of a period (through) as the day before the next period's start date. When that start date cannot be read as a date (YYYY-MM-DD), the end date is empty.

## Implied opening period when no history exists

Lines 136-144

Asked as: history empty shows period; project has rate sheet no periods; why opening period showing; no rows but rate sheet

A project that has a rate sheet but no history rows shows one opening period (id "opening") on its current rate sheet. A project with no rate sheet and no rows shows no periods.

## How period end dates are shown

Lines 145-161

Asked as: when does period end; period end date day before; period runs through; why latest period has no end

Each period runs through the day before the next period starts. The latest period has no end date. The opening period has an empty start and covers all earlier time.

## Viewing rate sheet history: permission flags

Lines 176-195

Asked as: can revise flag; can set first flag; why cannot edit rate history; rate history permission flags

The history response lists the periods and two flags. can_revise is true for an active user with the job and rate_sheet_revise claims. can_set_first is true when the user has the job claim and the project has no rate sheet.

## History is validated and the job's rate sheet synced after each change

Lines 198-256

Asked as: rate sheet not saved after change; job rate sheet sync; change rolled back rates; rate history save failed

Each change runs in one transaction. Turbo then validates the whole history and sets the job's rate sheet to match the latest period. If any step fails, nothing is saved.

## Rate sheet changes blocked when job editing is off

Lines 200-208

Asked as: job editing disabled rates; rates locked jobs; rate sheet editing turned off; cannot change rates jobs locked

Adding, changing or deleting a period is refused when job editing is disabled (utilities.ErrJobsEditingDisabled). The check runs before anything else.

## Rate sheet changes need an active user with the job claim

Lines 213-219

Asked as: inactive user cannot change rates; need job permission change rates; user not active rate change; job claim needed rates

Every change to a project's rate sheet history is refused with "An active user with the job claim is required." unless the user is active and has the job claim.

## Pending rate sheet absorb blocks history changes

Lines 258-290

Asked as: absorb waiting commit rate sheet; pending absorb blocks rates; undo absorb rate history; rate sheet absorb unfinished

Turbo refuses to change a project's history while a rate sheet absorb that moved the project's rate sheet, or one of its periods, is still waiting to be committed or undone. The message is "A rate sheet absorb that moved this project's rate sheet is waiting to be committed or undone. Finish it before changing the rate sheet history."

## Changing history requires rate_sheet_revise claim

Lines 292-297

Asked as: need revise permission period; cannot add period permission; rate_sheet_revise missing; no permission edit rate periods

Dated additions, edits and deletions of periods are refused with "Changing a project's rate sheet history requires the rate_sheet_revise claim." unless the user holds that claim.

## Reason is required and limited to 500 characters

Lines 299-308

Asked as: empty reason rate change; explain why rate changing; reason field required; reason required rate period

A period change needs a reason. An empty reason gives "Explain why the rate sheet is changing." A reason over 500 characters gives "Keep the reason to 500 characters or fewer." Turbo trims the reason first and counts characters, not bytes.

## Two periods cannot start on the same date

Lines 310-322

Asked as: same start date twice; date already used period; another period starts this date; two periods same date

Turbo refuses a period start date that another period of the same project already uses, with "Another period already starts on this date." A period being edited ignores itself in this check.

## Adding a period requires a rate sheet

Lines 334-336

Asked as: add period no rate sheet; select a rate sheet error; period without rate sheet; rate sheet not chosen

Adding a period without a rate sheet is refused with "Select a rate sheet."

## Setting a project's first rate sheet needs only the job claim

Lines 343-351

Asked as: set first rate sheet; first rate sheet no date; first rate sheet permission; add rate sheet without start

An add with no start date sets the first rate sheet of a project that has none. It needs only the job claim, not rate_sheet_revise. The reason is optional but still limited to 500 characters.

## Start date required once project has a rate sheet

Lines 343-346

Asked as: enter start date rate sheet; date required new rate sheet; new rate sheet no date; start date needed existing rates

When the project already has a rate sheet, an add without a start date is refused with "Enter the date the new rate sheet starts."

## Adding a dated period

Lines 353-378

Asked as: add dated period; set first rate sheet before dated; add period with date; new dated rate period

A dated add needs the rate_sheet_revise claim and a reason. The project must already have a first rate sheet, otherwise Turbo says "Set the project's first rate sheet before adding a dated one." Turbo writes the opening row first if the history does not exist yet, then checks that the date is free.

## New periods record who made them

Lines 382-395

Asked as: created by period; who changed period; updated by rate sheet; who made rate period

A new period stores the user as both creator (created_by) and last updater (updated_by).

## Finding a period to change or delete

Lines 397-416

Asked as: period not found reload; history changed reload; cannot find period edit; rate period missing

Turbo writes the opening row first if the history does not exist. The id "opening" addresses the period with no start date. Otherwise the period must belong to this project, or Turbo says "Rate sheet period not found. The history may have changed; reload it."

## Updating a period

Lines 418-458

Asked as: edit rate period; change period start date; update rate sheet period; changing opening rate reprices

Updating a period needs rate_sheet_revise and a reason. It can change the rate sheet or the start date. An empty rate sheet gives "Select a rate sheet." A new start date must be free. The reason and updater are always overwritten. Changing the opening period's sheet reprices all time before the first dated period.

## Deleting a period

Lines 460-478

Asked as: delete rate period; cannot delete opening period; remove rate sheet period; delete period opening

Deleting a period needs rate_sheet_revise. The deleted period's time is priced by the period before it. The opening period cannot be deleted: "The opening period cannot be deleted. Change its rate sheet instead."
