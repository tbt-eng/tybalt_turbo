<!-- A model wrote this file from the source. Do not edit it: run `go run ./cmd/helprules` in app/. -->
<!-- helprules {"source":"app/routes/job_rate_sheets.go","sha256":"f242880c75e9b0262c714b62487325bb994f0fa569feb972ee78f53a9780c2e1","model":"claude-sonnet-5-5"} -->

# Rules in app/routes/job_rate_sheets.go

## Opening period id for a project with no saved history

Lines 22-24

Asked as: opening period id; project with no rate sheet history; history never written; opening row missing

A project whose rate sheet history has not been written yet has its opening period addressed by the id "opening" (openingRateSheetPeriodID). Turbo creates the opening period's row when the history is first changed.

## Who may change or set a project's rate sheet history

Lines 85-95

Asked as: job claim required; who can change rate sheet; missing rate_sheet_revise permission; can't set first rate sheet

Setting a project's first rate sheet needs the job claim (procurement). Adding a dated period, changing a period or deleting a period also needs the rate_sheet_revise claim. Without the job claim, Turbo says "A user with the job claim is required." Without the revise claim, it says "Changing a project's rate sheet history requires the rate_sheet_revise claim."

## Proposals have no rate sheet history

Lines 97-106

Asked as: proposal rate sheet; P job number rate sheet; why no rate sheets on proposal; project not found rate sheet

Turbo refuses to read or change the rate sheet history of a job whose number starts with P (a proposal). The message is "Proposals do not have rate sheets." If the project id does not exist, the message is "Project not found."

## History shows who created and last updated each period

Lines 108-117

Asked as: who last updated period; creator name on rate sheet history; blank name in history; last changed by

Each period in the history shows the given name and surname of the person who created it (created_by) and of the person who last updated it (updated_by). The name is blank when Turbo has no profile for that person or the field is empty. The history also shows when the period was last updated.

## Period end date is the day before the next period starts

Lines 119-125

Asked as: period end date wrong; through date blank; when does rate sheet period end; end date empty

Turbo works out the last day of a period (Through) as the day before the next period's start date. If the next start date can't be read as a date (YYYY-MM-DD), the end date is blank.

## Implied opening period when no history rows exist

Lines 132-141

Asked as: history empty but has rate sheet; one opening period shown; no periods showing; project shows no periods

A project that has a rate sheet but no history rows shows one opening period (id "opening") on its current rate sheet. It has no start date. A project with neither shows no periods.

## Period start and end dates in the history

Lines 142-158

Asked as: when does rate sheet start; period start date empty; opening period dates; latest period no end date

A period starts on its effective date. An opening period has an empty start and covers all earlier time. Each period ends the day before the next period starts. The latest period has no end date.

## History lists details of each period's rate sheet

Lines 162-171

Asked as: rate sheet code in history; which revision is active; rate sheet details per period; active flag history

Each period in the history shows its rate sheet's code, name, revision, effective date and whether it is active.

## Viewing history shows what the user may do

Lines 173-194

Asked as: why cannot edit rate sheet history; can_revise false; set first rate sheet button missing; who can edit history

The history view returns the project's periods. It reports can_revise as true only for holders of both the job claim and rate_sheet_revise. It reports can_set_first as true only for holders of the job claim when the project has no rate sheet.

## History change is validated and syncs the project's rate sheet

Lines 196-254

Asked as: rate sheet not saved after change; jobs rate_sheet out of sync; change rolled back; history validation failed

Each change runs in one transaction. Afterwards Turbo validates the whole history and sets the project's rate sheet (jobs.rate_sheet) to match the latest period. If any step fails, the change is not kept.

## Rate sheet changes blocked when jobs editing is off

Lines 198-206

Asked as: cannot edit rate sheet jobs editing off; rate sheet locked; job editing disabled rate sheet; can't add rate sheet period

Turbo refuses any add, change or delete of a rate sheet period when job editing is disabled. The error is utilities.ErrJobsEditingDisabled; its message is not in this file.

## Pending rate sheet absorb blocks history changes

Lines 256-288

Asked as: absorb waiting to be committed; undo rate sheet absorb; pending absorb blocks rate sheet; can't change history after absorb

Turbo refuses to change a project's rate sheet history while an uncommitted rate sheet absorb that moved the project's rate sheet or one of its periods can still be undone. The message is "A rate sheet absorb that moved this project's rate sheet is waiting to be committed or undone. Finish it before changing the rate sheet history." The user must commit or undo the absorb first.

## Reason is required and limited to 500 characters

Lines 297-306

Asked as: reason required rate sheet; explain why rate sheet changing; reason too long; 500 character limit reason

A dated add or a change to a period needs a reason, which Turbo trims first. An empty reason gives "Explain why the rate sheet is changing." A reason over 500 characters gives "Keep the reason to 500 characters or fewer."

## Only one period may start on a date

Lines 308-320

Asked as: two periods same start date; date already taken rate sheet; another period starts this day; duplicate rate sheet date

Turbo refuses to add or move a period to a date on which another period of the same project already starts. The message is "Another period already starts on this date."

## Adding a period requires selecting a rate sheet

Lines 332-334

Asked as: no rate sheet selected; rate sheet dropdown empty; select a rate sheet error; add period without sheet

Adding a period without a rate sheet fails with "Select a rate sheet."

## Setting a project's first rate sheet (no date)

Lines 341-349

Asked as: add first rate sheet no date; enter start date rate sheet; project already has rate sheet; optional reason first sheet

Adding a period with no start date sets the first rate sheet and needs only the job claim. It is refused if the project already has a rate sheet, with "Enter the date the new rate sheet starts." The reason is optional here, and only a reason over 500 characters is refused. The new period has an empty start date.

## Adding a dated rate sheet period

Lines 350-376

Asked as: add dated rate sheet; set first rate sheet before dated; new rate sheet from date; add rate sheet effective date

Adding a period with a start date needs the revise claim and a project that already has a rate sheet. Without one, Turbo says "Set the project's first rate sheet before adding a dated one." The reason is required. If the history is not written yet, Turbo first creates the opening row from the current sheet. The user's id is saved as creator and updater.

## New period records who created it and why

Lines 380-392

Asked as: who created rate sheet period; recorded creator of period; saved rate sheet period fields; updated by on new period

When Turbo saves a new rate sheet period, it records the project, the chosen rate sheet, the start date, the reason, and the person making the change as both creator (created_by) and last updater (updated_by). The row's own hooks then check the new period.

## Finding a period by id

Lines 395-414

Asked as: period not found reload; rate sheet period lookup; opening period lookup fails; history may have changed

A period id of "opening" finds the period with no start date; if none exists, Turbo says "This project has no rate sheet." Any other id must belong to this project, or Turbo says "Rate sheet period not found. The history may have changed; reload it." Looking up a period first creates the opening row if the history is not written yet.

## Changing a period's rate sheet or start date

Lines 416-456

Asked as: change rate sheet on period; move rate sheet start date; edit existing rate sheet period; reprice before first period

Changing a period needs the revise claim and a reason. The user may change its rate sheet, its start date or both. An empty rate sheet gives "Select a rate sheet." The new date must not match another period. The reason and updater are always overwritten. Changing the opening period's sheet reprices all time before the first dated period.

## Deleting a period

Lines 458-476

Asked as: remove rate sheet period; delete rate sheet period; cannot delete opening period; undo added rate sheet period

Deleting a period needs the revise claim. Time in the deleted period is then priced by the period before it. The opening period cannot be deleted: "The opening period cannot be deleted. Change its rate sheet instead."
